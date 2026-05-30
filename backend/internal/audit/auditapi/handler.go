package auditapi

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/audit"
	"github.com/Abraxas-365/hada-commerce/internal/audit/auditsrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Handler exposes HTTP endpoints for the audit-log domain.
type Handler struct {
	svc *auditsrv.Service
}

// NewHandler creates a new audit HTTP handler.
func NewHandler(svc *auditsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all audit routes (protected — require auth).
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/audit")
	g.Get("/", h.List)
	g.Get("/stats", h.Stats)
	g.Get("/:id", h.GetByID)
	g.Post("/", h.Create)
}

// List handles GET /audit?user_id=&action=&resource_type=&resource_id=&from=&to=&page=&page_size=
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	filter := audit.AuditFilter{
		UserID:       c.Query("user_id"),
		Action:       c.Query("action"),
		ResourceType: c.Query("resource_type"),
		ResourceID:   c.Query("resource_id"),
	}

	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return errx.New("invalid 'from' timestamp; use RFC3339 format", errx.TypeValidation)
		}
		filter.From = &t
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return errx.New("invalid 'to' timestamp; use RFC3339 format", errx.TypeValidation)
		}
		filter.To = &t
	}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	result, err := h.svc.List(c.Context(), authCtx.TenantID, filter, page, pageSize)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// GetByID handles GET /audit/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	id := kernel.AuditEntryID(c.Params("id"))
	if id.IsEmpty() {
		return errx.New("audit entry id is required", errx.TypeValidation)
	}

	entry, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(entry)
}

// Stats handles GET /audit/stats?from=&to=
func (h *Handler) Stats(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	// Default: last 30 days.
	to := time.Now().UTC()
	from := to.AddDate(0, 0, -30)

	if fromStr := c.Query("from"); fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err != nil {
			return errx.New("invalid 'from' timestamp; use RFC3339 format", errx.TypeValidation)
		}
		from = t
	}
	if toStr := c.Query("to"); toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err != nil {
			return errx.New("invalid 'to' timestamp; use RFC3339 format", errx.TypeValidation)
		}
		to = t
	}

	stats, err := h.svc.GetStats(c.Context(), authCtx.TenantID, from, to)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"stats": stats,
		"from":  from.Format(time.RFC3339),
		"to":    to.Format(time.RFC3339),
	})
}

// Create handles POST /audit — allows services to log audit events directly via the API.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	var input audit.CreateAuditInput
	if err := c.BodyParser(&input); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	// Always use the authenticated tenant and user — do not allow caller override.
	input.TenantID = authCtx.TenantID
	if authCtx.UserID != nil {
		input.UserID = string(*authCtx.UserID)
	}
	if input.UserEmail == "" {
		input.UserEmail = authCtx.Email
	}

	// Capture IP and User-Agent from request if not explicitly set.
	if input.IPAddress == "" {
		input.IPAddress = c.IP()
	}
	if input.UserAgent == "" {
		input.UserAgent = c.Get("User-Agent")
	}

	entry, err := h.svc.Log(c.Context(), input)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(entry)
}
