package paymentapi

import (
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/payment/paymentsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the payment domain.
type Handler struct {
	svc *paymentsrv.Service
}

// NewHandler creates a new payment API handler.
func NewHandler(svc *paymentsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all payment routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/payments")
	g.Post("/", h.Create)
	g.Get("/:id", h.GetByID)
	g.Post("/:id/process", h.Process)
	g.Post("/:id/refund", h.Refund)
	g.Get("/:id/refunds", h.ListRefunds)
	g.Get("/order/:orderId", h.ListByOrder)
}

// ─── Request/response types ──────────────────────────────────────────────────

type createPaymentRequest struct {
	OrderID  string `json:"order_id"`
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
	Provider string `json:"provider"`
	Method   string `json:"method"`
}

type processPaymentRequest struct {
	Token string `json:"token"`
}

type createRefundRequest struct {
	Amount int64  `json:"amount"`
	Reason string `json:"reason"`
}

// ─── Handlers ────────────────────────────────────────────────────────────────

// Create handles POST /payments — creates a pending payment.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createPaymentRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.OrderID == "" {
		return errx.New("order_id is required", errx.TypeValidation)
	}
	if req.Amount <= 0 {
		return errx.New("amount must be greater than zero", errx.TypeValidation)
	}
	if req.Currency == "" {
		return errx.New("currency is required", errx.TypeValidation)
	}
	if req.Provider == "" {
		req.Provider = "manual"
	}

	p, err := h.svc.CreatePayment(
		c.Context(),
		authCtx.TenantID,
		kernel.OrderID(req.OrderID),
		req.Amount,
		req.Currency,
		req.Provider,
		req.Method,
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(p)
}

// GetByID handles GET /payments/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PaymentID(c.Params("id"))

	p, err := h.svc.GetPayment(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(p)
}

// Process handles POST /payments/:id/process — charges via the provider.
func (h *Handler) Process(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PaymentID(c.Params("id"))

	var req processPaymentRequest
	// Token is optional for manual provider; ignore parse errors
	_ = c.BodyParser(&req)

	p, err := h.svc.ProcessPayment(c.Context(), authCtx.TenantID, id, req.Token)
	if err != nil {
		return err
	}

	return c.JSON(p)
}

// Refund handles POST /payments/:id/refund — creates a refund.
func (h *Handler) Refund(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PaymentID(c.Params("id"))

	var req createRefundRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Amount <= 0 {
		return errx.New("amount must be greater than zero", errx.TypeValidation)
	}

	r, err := h.svc.CreateRefund(c.Context(), authCtx.TenantID, id, req.Amount, req.Reason)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(r)
}

// ListRefunds handles GET /payments/:id/refunds.
func (h *Handler) ListRefunds(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PaymentID(c.Params("id"))

	refunds, err := h.svc.ListRefunds(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(refunds)
}

// ListByOrder handles GET /payments/order/:orderId.
func (h *Handler) ListByOrder(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	orderID := kernel.OrderID(c.Params("orderId"))

	payments, err := h.svc.ListPaymentsByOrder(c.Context(), authCtx.TenantID, orderID)
	if err != nil {
		return err
	}

	return c.JSON(payments)
}
