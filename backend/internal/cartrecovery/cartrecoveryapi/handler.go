package cartrecoveryapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/cartrecovery"
	"github.com/Abraxas-365/hada-commerce/internal/cartrecovery/cartrecoverysrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Handler exposes HTTP endpoints for the cart recovery admin domain.
type Handler struct {
	svc *cartrecoverysrv.Service
}

// NewHandler creates a new cart recovery API handler.
func NewHandler(svc *cartrecoverysrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all cart-recovery admin routes on the given router.
// The router must already have authentication middleware applied.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/cart-recovery")
	g.Get("/", h.List)
	g.Get("/stats", h.GetStats)
	g.Put("/:id/status", h.UpdateStatus)
}

// ─────────────────────────────────────────────────────────────────────────────
// Request types
// ─────────────────────────────────────────────────────────────────────────────

type updateStatusRequest struct {
	Status string `json:"status"`
}

// ─────────────────────────────────────────────────────────────────────────────
// Handlers
// ─────────────────────────────────────────────────────────────────────────────

// List handles GET /cart-recovery — returns a paginated list of recovery emails.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	result, err := h.svc.List(c.Context(), authCtx.TenantID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetStats handles GET /cart-recovery/stats — returns recovery stats for the tenant.
func (h *Handler) GetStats(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	stats, err := h.svc.GetStats(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(stats)
}

// UpdateStatus handles PUT /cart-recovery/:id/status — transitions the recovery email status.
func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	id := kernel.NewRecoveryID(c.Params("id"))
	if id.IsEmpty() {
		return errx.New("id is required", errx.TypeValidation)
	}

	var req updateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	var (
		rec *cartrecovery.RecoveryEmail
		err error
	)

	switch req.Status {
	case cartrecovery.StatusSent:
		rec, err = h.svc.MarkSent(c.Context(), authCtx.TenantID, id)
	case cartrecovery.StatusClicked:
		rec, err = h.svc.MarkClicked(c.Context(), authCtx.TenantID, id)
	case cartrecovery.StatusConverted:
		rec, err = h.svc.MarkConverted(c.Context(), authCtx.TenantID, id)
	default:
		return errx.New("status must be one of: sent, clicked, converted", errx.TypeValidation)
	}

	if err != nil {
		return err
	}

	return c.JSON(rec)
}
