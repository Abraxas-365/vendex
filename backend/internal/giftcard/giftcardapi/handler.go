package giftcardapi

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/giftcard"
	"github.com/Abraxas-365/hada-commerce/internal/giftcard/giftcardsrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Handler exposes gift card HTTP endpoints.
type Handler struct {
	svc *giftcardsrv.Service
}

// New creates a new gift card Handler.
func New(svc *giftcardsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes wires all admin (protected) gift card routes onto the provided Fiber router.
//
//	POST   /admin/gift-cards                       — create gift card
//	GET    /admin/gift-cards                       — list gift cards
//	GET    /admin/gift-cards/:id                   — get gift card by ID
//	PUT    /admin/gift-cards/:id                   — update gift card
//	DELETE /admin/gift-cards/:id                   — delete gift card
//	POST   /admin/gift-cards/:id/disable           — disable gift card
//	GET    /admin/gift-cards/:id/transactions      — list transactions
func (h *Handler) RegisterRoutes(router fiber.Router) {
	admin := router.Group("/admin/gift-cards")
	admin.Post("/", h.create)
	admin.Get("/", h.list)
	admin.Get("/:id", h.getByID)
	admin.Put("/:id", h.update)
	admin.Delete("/:id", h.delete)
	admin.Post("/:id/disable", h.disable)
	admin.Get("/:id/transactions", h.listTransactions)
}

// RegisterPublicRoutes wires all public gift card routes.
//
//	POST /gift-cards/check-balance — check balance by code
//	POST /gift-cards/redeem        — redeem a gift card
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Post("/gift-cards/check-balance", h.checkBalance)
	router.Post("/gift-cards/redeem", h.redeem)
}

// ---------------------------------------------------------------------------
// Admin handlers
// ---------------------------------------------------------------------------

type createGiftCardRequest struct {
	Code          string     `json:"code"`
	AmountCents   int64      `json:"amount_cents"`
	Currency      string     `json:"currency"`
	ExpiresAt     *time.Time `json:"expires_at,omitempty"`
	CreatedBy     string     `json:"created_by"`
}

func (h *Handler) create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createGiftCardRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Code == "" {
		return errx.New("code is required", errx.TypeValidation)
	}
	if req.AmountCents <= 0 {
		return errx.New("amount_cents must be greater than zero", errx.TypeValidation)
	}
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	gc, err := h.svc.Create(c.Context(), authCtx.TenantID, giftcard.CreateInput{
		Code:          req.Code,
		InitialAmount: kernel.Money{Amount: req.AmountCents, Currency: currency},
		ExpiresAt:     req.ExpiresAt,
		CreatedBy:     req.CreatedBy,
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(gc)
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
	id := kernel.GiftCardID(c.Params("id"))

	gc, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(gc)
}

type updateGiftCardRequest struct {
	Code      *string    `json:"code,omitempty"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
	Active    *bool      `json:"active,omitempty"`
}

func (h *Handler) update(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.GiftCardID(c.Params("id"))

	var req updateGiftCardRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	gc, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	if req.Code != nil {
		gc.Code = *req.Code
	}
	if req.ExpiresAt != nil {
		gc.ExpiresAt = req.ExpiresAt
	}
	if req.Active != nil {
		gc.Active = *req.Active
	}

	// Use the repo via service would require an Update method on service;
	// for admin updates we go through the repo indirectly via Disable or by
	// returning errors if code/expiry changes are needed. Here we just expose
	// the disable flow and GetByID. For a complete update we build a lightweight
	// helper on the service.
	if err := h.svc.UpdateCard(c.Context(), authCtx.TenantID, gc); err != nil {
		return err
	}
	return c.JSON(gc)
}

func (h *Handler) delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.GiftCardID(c.Params("id"))

	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) disable(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.GiftCardID(c.Params("id"))

	gc, err := h.svc.Disable(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(gc)
}

func (h *Handler) listTransactions(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.GiftCardID(c.Params("id"))

	txs, err := h.svc.ListTransactions(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(txs)
}

// ---------------------------------------------------------------------------
// Public handlers
// ---------------------------------------------------------------------------

type checkBalanceRequest struct {
	Code string `json:"code"`
}

func (h *Handler) checkBalance(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID required", errx.TypeAuthorization)
	}

	var req checkBalanceRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Code == "" {
		return errx.New("code is required", errx.TypeValidation)
	}

	gc, err := h.svc.CheckBalance(c.Context(), tenantID, req.Code)
	if err != nil {
		return err
	}
	return c.JSON(gc)
}

type redeemRequest struct {
	Code     string  `json:"code"`
	Amount   int64   `json:"amount"`
	Currency string  `json:"currency"`
	OrderID  *string `json:"order_id,omitempty"`
	Note     string  `json:"note,omitempty"`
}

func (h *Handler) redeem(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID required", errx.TypeAuthorization)
	}

	var req redeemRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Code == "" {
		return errx.New("code is required", errx.TypeValidation)
	}
	if req.Amount <= 0 {
		return errx.New("amount must be greater than zero", errx.TypeValidation)
	}
	currency := req.Currency
	if currency == "" {
		currency = "USD"
	}

	input := giftcardsrv.RedeemInput{
		Code:   req.Code,
		Amount: kernel.Money{Amount: req.Amount, Currency: currency},
		Note:   req.Note,
	}
	if req.OrderID != nil {
		oid := kernel.OrderID(*req.OrderID)
		input.OrderID = &oid
	}

	gc, err := h.svc.Redeem(c.Context(), tenantID, input)
	if err != nil {
		return err
	}
	return c.JSON(gc)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func paginationFromCtx(c *fiber.Ctx) kernel.PaginationOptions {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("page_size"))
	return kernel.PaginationOptions{Page: page, PageSize: size}
}
