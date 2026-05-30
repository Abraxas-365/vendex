package webhookapi

import (
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/webhook"
	"github.com/Abraxas-365/hada-commerce/internal/webhook/webhooksrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the webhook domain.
type Handler struct {
	svc *webhooksrv.Service
}

// NewHandler creates a new webhook API handler.
func NewHandler(svc *webhooksrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected webhook routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/webhooks")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
	g.Put("/:id/toggle", h.Toggle)
	g.Get("/:id/deliveries", h.ListDeliveries)
	g.Post("/deliveries/:deliveryID/retry", h.RetryDelivery)
	g.Get("/deliveries/:deliveryID", h.GetDelivery)
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

type createRequest struct {
	URL         string   `json:"url"`
	Secret      string   `json:"secret"`
	Events      []string `json:"events"`
	Description string   `json:"description"`
}

type updateRequest struct {
	URL         *string  `json:"url,omitempty"`
	Secret      *string  `json:"secret,omitempty"`
	Events      []string `json:"events,omitempty"`
	Description *string  `json:"description,omitempty"`
}

type toggleRequest struct {
	Active bool `json:"active"`
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// Create handles POST /webhooks.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.URL == "" {
		return errx.New("url is required", errx.TypeValidation)
	}
	if len(req.Events) == 0 {
		return errx.New("at least one event type is required", errx.TypeValidation)
	}

	input := webhook.CreateWebhookInput{
		URL:         req.URL,
		Secret:      req.Secret,
		Events:      req.Events,
		Description: req.Description,
	}

	wh, err := h.svc.Create(c.Context(), authCtx.TenantID, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(wh)
}

// GetByID handles GET /webhooks/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewWebhookID(c.Params("id"))

	wh, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(wh)
}

// List handles GET /webhooks.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	page, pageSize := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// Update handles PUT /webhooks/:id.
func (h *Handler) Update(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewWebhookID(c.Params("id"))

	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	input := webhook.UpdateWebhookInput{
		URL:         req.URL,
		Secret:      req.Secret,
		Events:      req.Events,
		Description: req.Description,
	}

	wh, err := h.svc.Update(c.Context(), authCtx.TenantID, id, input)
	if err != nil {
		return err
	}

	return c.JSON(wh)
}

// Delete handles DELETE /webhooks/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewWebhookID(c.Params("id"))

	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Toggle handles PUT /webhooks/:id/toggle.
func (h *Handler) Toggle(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewWebhookID(c.Params("id"))

	var req toggleRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	wh, err := h.svc.Toggle(c.Context(), authCtx.TenantID, id, req.Active)
	if err != nil {
		return err
	}

	return c.JSON(wh)
}

// ListDeliveries handles GET /webhooks/:id/deliveries.
func (h *Handler) ListDeliveries(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	webhookID := kernel.NewWebhookID(c.Params("id"))
	page, pageSize := paginationFromQuery(c)

	result, err := h.svc.ListDeliveries(c.Context(), authCtx.TenantID, webhookID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetDelivery handles GET /webhooks/deliveries/:deliveryID.
func (h *Handler) GetDelivery(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewWebhookDeliveryID(c.Params("deliveryID"))

	delivery, err := h.svc.GetDelivery(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(delivery)
}

// RetryDelivery handles POST /webhooks/deliveries/:deliveryID/retry.
func (h *Handler) RetryDelivery(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewWebhookDeliveryID(c.Params("deliveryID"))

	delivery, err := h.svc.RetryDelivery(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(delivery)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func paginationFromQuery(c *fiber.Ctx) (page, pageSize int) {
	page, _ = strconv.Atoi(c.Query("page"))
	pageSize, _ = strconv.Atoi(c.Query("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return
}
