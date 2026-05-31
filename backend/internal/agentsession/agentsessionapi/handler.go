// Package agentsessionapi provides HTTP handlers for agent session management.
package agentsessionapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/agentsession"
	"github.com/Abraxas-365/vendex/internal/agentsession/agentsessionsrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// tenantID extracts the tenant ID from the Fiber context (set by auth middleware).
func tenantID(c *fiber.Ctx) kernel.TenantID {
	return kernel.TenantID(c.Locals("tenant_id").(string))
}

// Handler exposes agent session CRUD over HTTP.
type Handler struct {
	svc *agentsessionsrv.Service
}

// NewHandler creates a new agentsession Handler.
func NewHandler(svc *agentsessionsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers agent session routes on the given router group.
// Expected to be called with a group like "/agent/sessions".
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Post("/", h.CreateSession)
	r.Get("/", h.ListSessions)
	r.Get("/:id", h.GetSession)
	r.Post("/:id/stop", h.StopSession)
	r.Post("/:id/messages", h.SendMessage)
	r.Get("/:id/messages", h.GetHistory)
}

// CreateSession creates a new agent workspace session from a preset.
func (h *Handler) CreateSession(c *fiber.Ctx) error {
	var req agentsession.CreateSessionRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	sess, err := h.svc.CreateSession(c.Context(), tenantID(c), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(sess)
}

// ListSessions returns paginated sessions for the tenant.
func (h *Handler) ListSessions(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	p := kernel.NewPaginationOptions(page, pageSize)

	result, err := h.svc.ListSessions(c.Context(), tenantID(c), p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// GetSession returns a single session by ID.
func (h *Handler) GetSession(c *fiber.Ctx) error {
	id := kernel.AgentSessionID(c.Params("id"))

	sess, err := h.svc.GetSession(c.Context(), tenantID(c), id)
	if err != nil {
		return err
	}
	return c.JSON(sess)
}

// StopSession stops a running session's container.
func (h *Handler) StopSession(c *fiber.Ctx) error {
	id := kernel.AgentSessionID(c.Params("id"))

	sess, err := h.svc.StopSession(c.Context(), tenantID(c), id)
	if err != nil {
		return err
	}
	return c.JSON(sess)
}

// sendMessageRequest is the JSON body for sending a chat message.
type sendMessageRequest struct {
	Content  string `json:"content"`
	Role     string `json:"role"`
	ToolName string `json:"tool_name,omitempty"`
}

// SendMessage appends a message to the session's chat history.
func (h *Handler) SendMessage(c *fiber.Ctx) error {
	sessionID := kernel.AgentSessionID(c.Params("id"))

	var req sendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Role == "" {
		req.Role = "user"
	}

	msg, err := h.svc.SendMessage(c.Context(), tenantID(c), sessionID, req.Role, req.Content, req.ToolName)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(msg)
}

// GetHistory returns paginated chat history for a session.
func (h *Handler) GetHistory(c *fiber.Ctx) error {
	sessionID := kernel.AgentSessionID(c.Params("id"))

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "50"))
	p := kernel.NewPaginationOptions(page, pageSize)

	result, err := h.svc.GetHistory(c.Context(), tenantID(c), sessionID, p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}
