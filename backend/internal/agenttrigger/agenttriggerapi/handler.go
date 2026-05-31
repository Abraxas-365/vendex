// Package agenttriggerapi provides HTTP handlers for event-triggered agent actions.
package agenttriggerapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/agenttrigger/agenttrigger"
	"github.com/Abraxas-365/vendex/internal/agenttrigger/agenttriggersrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// tenantID extracts the tenant ID from the Fiber context (set by auth middleware).
func tenantID(c *fiber.Ctx) kernel.TenantID {
	return kernel.TenantID(c.Locals("tenant_id").(string))
}

// Handler exposes agent trigger CRUD and execution history over HTTP.
type Handler struct {
	svc *agenttriggersrv.Service
}

// NewHandler creates a new agenttrigger Handler.
func NewHandler(svc *agenttriggersrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers agent trigger routes on the given router group.
// Expected to be called with a group like "/agent/triggers".
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/event-types", h.GetEventTypes)
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
	r.Post("/:id/enable", h.Enable)
	r.Post("/:id/disable", h.Disable)
	r.Get("/:id/logs", h.GetLogs)
}

// GetEventTypes returns the list of valid event types for UI dropdowns.
func (h *Handler) GetEventTypes(c *fiber.Ctx) error {
	types := h.svc.GetValidEventTypes()
	return c.JSON(fiber.Map{"event_types": types})
}

// List returns paginated triggers for the tenant.
func (h *Handler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	p := kernel.NewPaginationOptions(page, pageSize)

	result, err := h.svc.List(c.Context(), tenantID(c), p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// Create creates a new event trigger.
func (h *Handler) Create(c *fiber.Ctx) error {
	var req agenttrigger.CreateTriggerRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	trigger, err := h.svc.Create(c.Context(), tenantID(c), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(trigger)
}

// GetByID returns a single trigger by ID.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := kernel.AgentTriggerID(c.Params("id"))

	trigger, err := h.svc.Get(c.Context(), tenantID(c), id)
	if err != nil {
		return err
	}
	return c.JSON(trigger)
}

// Update updates a trigger's mutable fields.
func (h *Handler) Update(c *fiber.Ctx) error {
	id := kernel.AgentTriggerID(c.Params("id"))

	var req agenttrigger.UpdateTriggerRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	trigger, err := h.svc.Update(c.Context(), tenantID(c), id, req)
	if err != nil {
		return err
	}
	return c.JSON(trigger)
}

// Delete removes a trigger.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := kernel.AgentTriggerID(c.Params("id"))

	if err := h.svc.Delete(c.Context(), tenantID(c), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// Enable enables a trigger.
func (h *Handler) Enable(c *fiber.Ctx) error {
	id := kernel.AgentTriggerID(c.Params("id"))

	trigger, err := h.svc.Enable(c.Context(), tenantID(c), id)
	if err != nil {
		return err
	}
	return c.JSON(trigger)
}

// Disable disables a trigger.
func (h *Handler) Disable(c *fiber.Ctx) error {
	id := kernel.AgentTriggerID(c.Params("id"))

	trigger, err := h.svc.Disable(c.Context(), tenantID(c), id)
	if err != nil {
		return err
	}
	return c.JSON(trigger)
}

// GetLogs returns paginated execution logs for a trigger.
func (h *Handler) GetLogs(c *fiber.Ctx) error {
	id := kernel.AgentTriggerID(c.Params("id"))

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	p := kernel.NewPaginationOptions(page, pageSize)

	result, err := h.svc.GetLogs(c.Context(), tenantID(c), id, p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}
