// Package agentapi provides the HTTP handler for the AI store assistant chat endpoint.
// It streams agent events to the client using Server-Sent Events (SSE).
package agentapi

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Abraxas-365/hada-commerce/internal/agent"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/harness"
	harnesstools "github.com/Abraxas-365/harness/tools"
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

// ChatRequest is the JSON body accepted by the POST /agent/chat endpoint.
type ChatRequest struct {
	Message   string `json:"message"`
	SessionID string `json:"session_id,omitempty"`
}

// Handler manages the agent chat HTTP endpoint.
// It maintains a per-tenant session cache so conversation history
// persists across multiple HTTP requests.
type Handler struct {
	apiKey       string
	model        string
	systemPrompt string
	domainTools  []harnesstools.Tool

	mu       sync.RWMutex
	sessions map[string]*sessionEntry // key = tenantID + ":" + sessionID
}

type sessionEntry struct {
	session *harness.Session
	h       *harness.Harness
}

// NewHandler creates a new agent chat handler.
//
//   - apiKey: Anthropic API key
//   - model: model identifier, e.g. "claude-sonnet-4-20250514"
//   - systemPrompt: override the default store-assistant system prompt (pass "" for default)
//   - domainTools: pre-adapted harness tools (use agent.AdaptTools(agent.Setup(...)) to create these)
func NewHandler(apiKey, model, systemPrompt string, domainTools []harnesstools.Tool) *Handler {
	if systemPrompt == "" {
		systemPrompt = defaultSystemPrompt
	}
	return &Handler{
		apiKey:       apiKey,
		model:        model,
		systemPrompt: systemPrompt,
		domainTools:  domainTools,
		sessions:     make(map[string]*sessionEntry),
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
	sessionKey := tenantID + ":" + sessionID

	// Obtain or create a harness session for this tenant+session pair.
	sess, err := h.getOrCreateSession(sessionKey)
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
// Sessions are keyed by tenantID:sessionID and reuse conversation history.
func (h *Handler) getOrCreateSession(key string) (*harness.Session, error) {
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

	opts := []harness.Option{
		harness.WithAPIKey(h.apiKey),
		harness.WithModel(h.model),
		harness.WithSystemPrompt(h.systemPrompt),
		harness.WithPermissionMode("headless"),
	}
	if len(h.domainTools) > 0 {
		opts = append(opts, harness.WithTools(h.domainTools...))
	}

	harn, err := harness.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("creating harness: %w", err)
	}

	sess := harn.NewSession()
	h.sessions[key] = &sessionEntry{session: sess, h: harn}
	return sess, nil
}
