// Package agentapi provides the HTTP handler for the AI store assistant chat endpoint.
// It streams agent events to the client using Server-Sent Events (SSE).
package agentapi

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/Abraxas-365/hada-commerce/internal/agent"
	"github.com/Abraxas-365/hada-commerce/internal/containerx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/harness"
	"github.com/gofiber/fiber/v2"
)

const defaultSystemPrompt = `You are an AI store assistant for an e-commerce platform. You help merchants manage their store by using the available tools.

You can:
- Create and manage products, collections, and bundles
- Query and manage orders, payments, and returns
- Manage shipping zones, tax rates, and currencies
- Handle customer groups, loyalty programs, and gift cards
- Manage blog posts, promotions, and A/B tests
- View dashboard analytics and audit logs
- Manage inventory, reviews, and recommendations

Always be concise and helpful. When asked to perform an action, use the appropriate tool. When displaying results, format them clearly.`

// PresetProvider retrieves preset configuration for the agent.
// This decouples agentapi from the marketplace package.
type PresetProvider interface {
	GetPresetSystemPrompt(ctx context.Context, presetID string) (string, error)
	GetPresetToolsManifest(ctx context.Context, presetID string) (json.RawMessage, error)
}

// WorkspaceProvider checks whether a session has an active workspace container.
// This decouples agentapi from the agentsession package.
type WorkspaceProvider interface {
	// GetActiveWorkspace returns the container ID and preview URL for a running session
	// workspace. Returns ("", "", nil) when no workspace is active for the session.
	GetActiveWorkspace(ctx context.Context, tenantID, sessionID string) (containerID string, previewURL string, err error)
}

// ChatRequest is the JSON body accepted by the POST /agent/chat endpoint.
type ChatRequest struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
	PresetID  string `json:"preset_id,omitempty"`
}

// Handler manages the agent chat HTTP endpoint.
// It maintains a per-tenant session cache so conversation history
// persists across multiple HTTP requests. Domain tools are created
// per-tenant dynamically using agent.Services, ensuring proper
// tenant scoping for multi-tenant deployments.
type Handler struct {
	apiKey            string
	model             string
	systemPrompt      string
	services          agent.Services
	presetProvider    PresetProvider    // may be nil — presets disabled
	workspaceProvider WorkspaceProvider // may be nil — workspace tools disabled
	containerMgr      containerx.Manager // may be nil — workspace tools disabled

	mu       sync.RWMutex
	sessions map[string]*sessionEntry // key = tenantID + ":" + presetID + ":" + sessionID
}

// staticAccessor implements workspace.ContainerAccessor with fixed container info.
type staticAccessor struct {
	containerID containerx.ID
	previewURL  string
}

func (a *staticAccessor) ContainerID() containerx.ID { return a.containerID }
func (a *staticAccessor) PreviewBaseURL() string      { return a.previewURL }

type sessionEntry struct {
	session *harness.Session
	h       *harness.Harness
}

// NewHandler creates a new agent chat handler.
//
//   - apiKey: Anthropic API key
//   - model: model identifier, e.g. "claude-sonnet-4-20250514"
//   - systemPrompt: override the default store-assistant system prompt (pass "" for default)
//   - services: domain services used to create tenant-scoped tools per session
//   - presetProvider: optional preset config provider (pass nil to disable presets)
//   - workspaceProvider: optional workspace provider; pass nil to disable workspace tools
//   - containerMgr: optional container manager; pass nil to disable workspace tools
func NewHandler(apiKey, model, systemPrompt string, services agent.Services, presetProvider PresetProvider, workspaceProvider WorkspaceProvider, containerMgr containerx.Manager) *Handler {
	if systemPrompt == "" {
		systemPrompt = defaultSystemPrompt
	}
	return &Handler{
		apiKey:            apiKey,
		model:             model,
		systemPrompt:      systemPrompt,
		services:          services,
		presetProvider:    presetProvider,
		workspaceProvider: workspaceProvider,
		containerMgr:      containerMgr,
		sessions:          make(map[string]*sessionEntry),
	}
}

// RegisterRoutes mounts the agent chat route onto the given Fiber router.
// Expected to be called with a protected (authenticated) router group.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	router.Post("/agent/chat", h.Chat)
}

// Chat handles POST /agent/chat.
//
// Request body:
//
//	{ "message": "...", "session_id": "optional-id-for-history" }
//
// Response: SSE stream of agent.Event JSON objects:
//
//	data: {"kind":"text_delta","text":"...","timestamp":"..."}\n\n
//	data: {"kind":"tool_start","tool_name":"...","tool_input":"...",...}\n\n
//	data: {"kind":"tool_end","tool_name":"...","result":"...",...}\n\n
//	data: {"kind":"turn_end","timestamp":"..."}\n\n
//
// The stream ends after the turn_end or error event.
func (h *Handler) Chat(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil || !authCtx.IsValid() {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "UNAUTHORIZED",
			"message": "authentication required",
		})
	}

	var req ChatRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "INVALID_BODY",
			"message": "invalid request body",
		})
	}
	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "MISSING_MESSAGE",
			"message": "message is required",
		})
	}

	tenantID := string(authCtx.TenantID)
	sessionID := req.SessionID
	if sessionID == "" {
		sessionID = "default"
	}
	presetID := req.PresetID
	sessionKey := tenantID + ":" + presetID + ":" + sessionID

	// Obtain or create a harness session for this tenant+session+preset triple.
	sess, err := h.getOrCreateSession(sessionKey, presetID)
	if err != nil {
		log.Printf("[agentapi] failed to create harness session key=%q err=%v", sessionKey, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "SESSION_ERROR",
			"message": "failed to initialise agent session",
		})
	}

	// Set SSE headers before streaming.
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")
	c.Set("X-Accel-Buffering", "no") // disable nginx buffering for SSE

	// Capture the message and context before entering the stream writer.
	message := req.Message
	ctx := c.Context()

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		// eventCh is a buffered channel that decouples event production
		// (harness agent loop goroutine) from event consumption (this writer).
		eventCh := make(chan agent.Event, 64)

		// Create a per-request EventHandler that forwards events to eventCh.
		evtHandler := agent.NewEventHandler(sessionKey)
		evtHandler.OnEvent = func(e agent.Event) {
			select {
			case eventCh <- e:
			default:
				// Drop event if channel is full to avoid blocking the agent loop.
				log.Printf("[agentapi] event channel full for session=%q, dropping event kind=%q", sessionKey, e.Kind)
			}
		}

		// Replace the session's event handler for this request.
		sess.SetHandler(evtHandler)

		// Run the agent loop in a background goroutine.
		done := make(chan error, 1)
		go func() {
			done <- sess.Send(context.WithoutCancel(ctx), message)
		}()

		// writeEvent serialises an agent.Event and writes it as an SSE line.
		writeEvent := func(e agent.Event) {
			data, err := json.Marshal(e)
			if err != nil {
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", data)
			if err := w.Flush(); err != nil {
				log.Printf("[agentapi] flush error session=%q err=%v", sessionKey, err)
			}
		}

		// Drain events until the turn ends or an error is received.
		for {
			select {
			case evt := <-eventCh:
				writeEvent(evt)
				if evt.Kind == agent.EventTurnEnd || evt.Kind == agent.EventError {
					// Drain any remaining queued events before returning.
					for {
						select {
						case remaining := <-eventCh:
							writeEvent(remaining)
						default:
							return
						}
					}
				}
			case err := <-done:
				if err != nil {
					errEvt := agent.Event{Kind: agent.EventError, Error: err.Error()}
					writeEvent(errEvt)
				}
				// Drain any final events flushed before done fired.
				for {
					select {
					case remaining := <-eventCh:
						writeEvent(remaining)
					default:
						return
					}
				}
			}
		}
	})

	return nil
}

// getOrCreateSession returns an existing harness session or creates a new one.
// Sessions are keyed by tenantID:presetID:sessionID and reuse conversation history.
// Tools are created per-tenant to ensure proper multi-tenant scoping.
// When presetID is non-empty, the preset's system prompt is used instead of the default.
func (h *Handler) getOrCreateSession(key string, presetID string) (*harness.Session, error) {
	h.mu.RLock()
	entry, ok := h.sessions[key]
	h.mu.RUnlock()
	if ok {
		return entry.session, nil
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	// Double-check after acquiring write lock.
	if entry, ok = h.sessions[key]; ok {
		return entry.session, nil
	}

	// Extract tenantID from the session key (format: "tenantID:presetID:sessionID").
	tenantID := key
	if idx := strings.IndexByte(key, ':'); idx >= 0 {
		tenantID = key[:idx]
	}

	// Determine the system prompt: use preset's if available, else default.
	sysPrompt := h.systemPrompt
	if presetID != "" && h.presetProvider != nil {
		if presetPrompt, err := h.presetProvider.GetPresetSystemPrompt(context.Background(), presetID); err == nil && presetPrompt != "" {
			sysPrompt = presetPrompt
		}
	}

	// Create tenant-scoped domain tools.
	domainTools := agent.AdaptTools(agent.Setup(kernel.TenantID(tenantID), h.services))

	// Inject workspace tools if the session has an active container.
	if h.workspaceProvider != nil && h.containerMgr != nil {
		// Extract sessionID — key format: "tenantID:presetID:sessionID"
		sessionIDPart := key
		if idx := strings.LastIndexByte(key, ':'); idx >= 0 {
			sessionIDPart = key[idx+1:]
		}
		if containerID, previewURL, err := h.workspaceProvider.GetActiveWorkspace(context.Background(), tenantID, sessionIDPart); err == nil && containerID != "" {
			accessor := &staticAccessor{
				containerID: containerx.ID(containerID),
				previewURL:  previewURL,
			}
			wsTools := agent.AdaptTools(agent.WorkspaceTools(h.containerMgr, accessor))
			domainTools = append(domainTools, wsTools...)
		}
	}

	opts := []harness.Option{
		harness.WithAPIKey(h.apiKey),
		harness.WithModel(h.model),
		harness.WithSystemPrompt(sysPrompt),
		harness.WithPermissionMode("headless"),
	}
	if len(domainTools) > 0 {
		opts = append(opts, harness.WithTools(domainTools...))
	}

	harn, err := harness.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("creating harness: %w", err)
	}

	sess := harn.NewSession()
	h.sessions[key] = &sessionEntry{session: sess, h: harn}
	return sess, nil
}
