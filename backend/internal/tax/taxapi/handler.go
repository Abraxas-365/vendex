package taxapi

import (
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/tax/taxsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the tax domain.
type Handler struct {
	svc *taxsrv.Service
}

// NewHandler creates a new tax API handler.
func NewHandler(svc *taxsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all protected tax admin routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/tax/rates")
	g.Post("/", h.CreateRate)
	g.Get("/", h.ListRates)
	g.Get("/:id", h.GetRate)
	g.Put("/:id", h.UpdateRate)
	g.Delete("/:id", h.DeleteRate)
}

// RegisterPublicRoutes registers the public tax calculation endpoint.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Post("/tax/calculate", h.CalculateTax)
}

// createRateRequest is the JSON body for creating a tax rate.
type createRateRequest struct {
	Name             string  `json:"name"`
	Rate             float64 `json:"rate"`
	Country          string  `json:"country"`
	State            string  `json:"state"`
	City             string  `json:"city"`
	ZipCode          string  `json:"zip_code"`
	Priority         int     `json:"priority"`
	Compound         bool    `json:"compound"`
	IncludesShipping bool    `json:"includes_shipping"`
	Active           bool    `json:"active"`
}

// calculateTaxRequest is the JSON body for calculating tax.
type calculateTaxRequest struct {
	Subtotal int64  `json:"subtotal"`
	Shipping int64  `json:"shipping"`
	Country  string `json:"country"`
	State    string `json:"state"`
	City     string `json:"city"`
	ZipCode  string `json:"zip_code"`
}

// CreateRate handles POST /tax/rates.
func (h *Handler) CreateRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}
	if req.Country == "" {
		return errx.New("country is required", errx.TypeValidation)
	}

	rate, err := h.svc.CreateRate(c.Context(), authCtx.TenantID, taxsrv.CreateRateInput{
		Name:             req.Name,
		Rate:             req.Rate,
		Country:          req.Country,
		State:            req.State,
		City:             req.City,
		ZipCode:          req.ZipCode,
		Priority:         req.Priority,
		Compound:         req.Compound,
		IncludesShipping: req.IncludesShipping,
		Active:           req.Active,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(rate)
}

// ListRates handles GET /tax/rates.
func (h *Handler) ListRates(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	rates, err := h.svc.ListRates(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(rates)
}

// GetRate handles GET /tax/rates/:id.
func (h *Handler) GetRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.TaxRateID(c.Params("id"))

	rate, err := h.svc.GetRate(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(rate)
}

// UpdateRate handles PUT /tax/rates/:id.
func (h *Handler) UpdateRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.TaxRateID(c.Params("id"))

	var req createRateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}
	if req.Country == "" {
		return errx.New("country is required", errx.TypeValidation)
	}

	rate, err := h.svc.UpdateRate(c.Context(), authCtx.TenantID, id, taxsrv.UpdateRateInput{
		Name:             req.Name,
		Rate:             req.Rate,
		Country:          req.Country,
		State:            req.State,
		City:             req.City,
		ZipCode:          req.ZipCode,
		Priority:         req.Priority,
		Compound:         req.Compound,
		IncludesShipping: req.IncludesShipping,
		Active:           req.Active,
	})
	if err != nil {
		return err
	}

	return c.JSON(rate)
}

// DeleteRate handles DELETE /tax/rates/:id.
func (h *Handler) DeleteRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.TaxRateID(c.Params("id"))

	if err := h.svc.DeleteRate(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// CalculateTax handles POST /tax/calculate.
// This is a public endpoint; tenant is identified via the X-Tenant-ID header.
func (h *Handler) CalculateTax(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	var req calculateTaxRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.Country == "" {
		return errx.New("country is required", errx.TypeValidation)
	}
	if req.Subtotal < 0 {
		return errx.New("subtotal must be non-negative", errx.TypeValidation)
	}

	result, err := h.svc.CalculateTax(
		c.Context(),
		tenantID,
		req.Subtotal,
		req.Shipping,
		req.Country,
		req.State,
		req.City,
		req.ZipCode,
	)
	if err != nil {
		return err
	}

	return c.JSON(result)
}


