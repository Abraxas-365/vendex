package currencyapi

import (
	"github.com/Abraxas-365/vendex/internal/currency/currencysrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the currency domain.
type Handler struct {
	svc *currencysrv.Service
}

// NewHandler creates a new currency API handler.
func NewHandler(svc *currencysrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected admin routes for exchange rate management.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/currency-rates")
	g.Post("/", h.SetRate)
	g.Get("/", h.ListRates)
	g.Delete("/:id", h.DeleteRate)

	router.Get("/currencies", h.ListSupportedCurrencies)
}

// RegisterPublicRoutes registers the public currency conversion endpoint.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Post("/currency/convert", h.Convert)
}

// setRateRequest is the JSON body for setting an exchange rate.
type setRateRequest struct {
	BaseCurrency   string  `json:"base_currency"`
	TargetCurrency string  `json:"target_currency"`
	Rate           float64 `json:"rate"`
	AutoUpdate     bool    `json:"auto_update"`
}

// convertRequest is the JSON body for converting a currency amount.
type convertRequest struct {
	Amount         int64  `json:"amount"`
	Currency       string `json:"currency"`
	TargetCurrency string `json:"target_currency"`
}

// SetRate handles POST /currency-rates.
// Creates or updates an exchange rate for the authenticated tenant.
func (h *Handler) SetRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req setRateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.BaseCurrency == "" {
		return errx.New("base_currency is required", errx.TypeValidation)
	}
	if req.TargetCurrency == "" {
		return errx.New("target_currency is required", errx.TypeValidation)
	}

	rate, err := h.svc.SetRate(c.Context(), authCtx.TenantID, currencysrv.SetRateInput{
		BaseCurrency:   req.BaseCurrency,
		TargetCurrency: req.TargetCurrency,
		Rate:           req.Rate,
		AutoUpdate:     req.AutoUpdate,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(rate)
}

// ListRates handles GET /currency-rates.
func (h *Handler) ListRates(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	rates, err := h.svc.ListRates(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(rates)
}

// DeleteRate handles DELETE /currency-rates/:id.
func (h *Handler) DeleteRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CurrencyRateID(c.Params("id"))

	if err := h.svc.DeleteRate(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListSupportedCurrencies handles GET /currencies.
// Returns the static list of currencies supported by the system.
func (h *Handler) ListSupportedCurrencies(c *fiber.Ctx) error {
	currencies := h.svc.ListSupportedCurrencies()
	return c.JSON(currencies)
}

// Convert handles POST /currency/convert.
// Public endpoint: tenant identified via X-Tenant-ID header.
func (h *Handler) Convert(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	var req convertRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.Currency == "" {
		return errx.New("currency is required", errx.TypeValidation)
	}
	if req.TargetCurrency == "" {
		return errx.New("target_currency is required", errx.TypeValidation)
	}
	if req.Amount < 0 {
		return errx.New("amount must be non-negative", errx.TypeValidation)
	}

	amount := kernel.Money{Amount: req.Amount, Currency: req.Currency}

	result, err := h.svc.Convert(c.Context(), tenantID, amount, req.TargetCurrency)
	if err != nil {
		return err
	}

	return c.JSON(result)
}
