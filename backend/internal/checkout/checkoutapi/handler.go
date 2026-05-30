package checkoutapi

import (
	"github.com/Abraxas-365/hada-commerce/internal/checkout"
	"github.com/Abraxas-365/hada-commerce/internal/checkout/checkoutsrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/order"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the checkout domain.
type Handler struct {
	svc *checkoutsrv.Service
}

// NewHandler creates a new checkout API handler.
func NewHandler(svc *checkoutsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterPublicRoutes registers unauthenticated checkout routes (tenant via X-Tenant-ID).
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/checkout")
	g.Post("/preview", h.Preview)
	g.Post("/process", h.Process)
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

type addressInput struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

type checkoutRequest struct {
	CartID          string        `json:"cart_id"`
	ShippingAddress addressInput  `json:"shipping_address"`
	BillingAddress  *addressInput `json:"billing_address"`
	ShippingRateID  string        `json:"shipping_rate_id"`
	PromoCode       string        `json:"promo_code"`
	PaymentProvider string        `json:"payment_provider"`
	PaymentMethod   string        `json:"payment_method"`
	PaymentToken    string        `json:"payment_token"`
	CustomerID      string        `json:"customer_id"`
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// Preview handles POST /checkout/preview — returns cost breakdown without committing.
func (h *Handler) Preview(c *fiber.Ctx) error {
	tenantID, err := tenantFromHeader(c)
	if err != nil {
		return err
	}

	var req checkoutRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.CartID == "" {
		return errx.New("cart_id is required", errx.TypeValidation)
	}
	if req.ShippingRateID == "" {
		return errx.New("shipping_rate_id is required", errx.TypeValidation)
	}

	input := buildCheckoutInput(req)
	summary, err := h.svc.Preview(c.Context(), tenantID, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(summary)
}

// Process handles POST /checkout/process — performs the full checkout.
func (h *Handler) Process(c *fiber.Ctx) error {
	tenantID, err := tenantFromHeader(c)
	if err != nil {
		return err
	}

	var req checkoutRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.CartID == "" {
		return errx.New("cart_id is required", errx.TypeValidation)
	}
	if req.ShippingRateID == "" {
		return errx.New("shipping_rate_id is required", errx.TypeValidation)
	}
	if req.PaymentProvider == "" {
		return errx.New("payment_provider is required", errx.TypeValidation)
	}

	input := buildCheckoutInput(req)
	result, err := h.svc.Process(c.Context(), tenantID, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func buildCheckoutInput(req checkoutRequest) checkout.CheckoutInput {
	input := checkout.CheckoutInput{
		CartID: kernel.CartID(req.CartID),
		ShippingAddress: order.Address{
			Street:     req.ShippingAddress.Street,
			City:       req.ShippingAddress.City,
			State:      req.ShippingAddress.State,
			Country:    req.ShippingAddress.Country,
			PostalCode: req.ShippingAddress.PostalCode,
		},
		ShippingRateID:  kernel.ShippingRateID(req.ShippingRateID),
		PromoCode:       req.PromoCode,
		PaymentProvider: req.PaymentProvider,
		PaymentMethod:   req.PaymentMethod,
		PaymentToken:    req.PaymentToken,
		CustomerID:      kernel.CustomerID(req.CustomerID),
	}

	if req.BillingAddress != nil {
		ba := order.Address{
			Street:     req.BillingAddress.Street,
			City:       req.BillingAddress.City,
			State:      req.BillingAddress.State,
			Country:    req.BillingAddress.Country,
			PostalCode: req.BillingAddress.PostalCode,
		}
		input.BillingAddress = &ba
	}

	return input
}

func tenantFromHeader(c *fiber.Ctx) (kernel.TenantID, error) {
	tid := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tid == "" {
		return "", errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}
	return tid, nil
}
