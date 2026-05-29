package promoapi

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/promo"
	"github.com/Abraxas-365/hada-commerce/internal/promo/promosrv"
)

// Handler exposes promo HTTP endpoints.
type Handler struct {
	svc *promosrv.Service
}

// New creates a new promo Handler.
func New(svc *promosrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes wires all promo routes onto the provided Fiber router.
//
//	POST   /admin/promos               — create promo
//	GET    /admin/promos               — list promos
//	GET    /admin/promos/:id           — get promo
//	POST   /admin/promos/:id/deactivate — deactivate
//	POST   /promos/validate            — validate a code for an order total (public)
func (h *Handler) RegisterRoutes(router fiber.Router) {
	admin := router.Group("/admin/promos")
	admin.Post("/", h.create)
	admin.Get("/", h.list)
	admin.Get("/:id", h.getByID)
	admin.Post("/:id/deactivate", h.deactivate)

	router.Post("/promos/validate", h.validate)
}

// --- handlers ---

type createPromoRequest struct {
	Code           string     `json:"code"`
	Type           string     `json:"type"`
	Value          int64      `json:"value"`
	MinOrderAmount *int64     `json:"min_order_amount,omitempty"`
	MaxUses        *int       `json:"max_uses,omitempty"`
	StartsAt       *time.Time `json:"starts_at,omitempty"`
	EndsAt         *time.Time `json:"ends_at,omitempty"`
}

func (h *Handler) create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createPromoRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	p, err := h.svc.Create(c.Context(), promosrv.CreateInput{
		TenantID:       authCtx.TenantID,
		Code:           req.Code,
		Type:           promo.PromoType(req.Type),
		Value:          req.Value,
		MinOrderAmount: req.MinOrderAmount,
		MaxUses:        req.MaxUses,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(p)
}

func (h *Handler) list(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	p := paginationFromCtx(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *Handler) getByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PromoID(c.Params("id"))

	p, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(p)
}

func (h *Handler) deactivate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PromoID(c.Params("id"))

	p, err := h.svc.Deactivate(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(p)
}

type validateRequest struct {
	Code            string `json:"code"`
	OrderTotalCents int64  `json:"order_total_cents"`
}

func (h *Handler) validate(c *fiber.Ctx) error {
	// validate is public — use X-Tenant-ID header
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID required", errx.TypeAuthorization)
	}

	var req validateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	result, err := h.svc.Validate(c.Context(), tenantID, req.Code, req.OrderTotalCents)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// --- helpers ---

func paginationFromCtx(c *fiber.Ctx) kernel.PaginationOptions {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("page_size"))
	return kernel.PaginationOptions{Page: page, PageSize: size}
}
