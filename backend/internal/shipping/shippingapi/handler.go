package shippingapi

import (
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/shipping"
	"github.com/Abraxas-365/vendex/internal/shipping/shippingsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the shipping domain.
type Handler struct {
	svc *shippingsrv.Service
}

// NewHandler creates a new shipping API handler.
func NewHandler(svc *shippingsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected admin routes.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Zone CRUD
	zones := router.Group("/shipping/zones")
	zones.Post("/", h.CreateZone)
	zones.Get("/", h.ListZones)
	zones.Get("/:id", h.GetZone)
	zones.Put("/:id", h.UpdateZone)
	zones.Delete("/:id", h.DeleteZone)

	// Rates nested under zones
	zones.Post("/:zoneId/rates", h.CreateRate)
	zones.Get("/:zoneId/rates", h.ListRates)

	// Rates standalone (for update/delete by rate ID)
	rates := router.Group("/shipping/rates")
	rates.Put("/:id", h.UpdateRate)
	rates.Delete("/:id", h.DeleteRate)
	rates.Get("/:id", h.GetRate)
}

// RegisterPublicRoutes registers public (unauthenticated) routes.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Post("/shipping/calculate", h.CalculateShipping)
}

// ---------------------------------------------------------------------------
// Zone handlers
// ---------------------------------------------------------------------------

type createZoneRequest struct {
	Name      string   `json:"name"`
	Countries []string `json:"countries"`
	States    []string `json:"states"`
}

// CreateZone handles POST /shipping/zones.
func (h *Handler) CreateZone(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createZoneRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	zone, err := h.svc.CreateZone(c.Context(), authCtx.TenantID, req.Name, req.Countries, req.States)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(zone)
}

// GetZone handles GET /shipping/zones/:id.
func (h *Handler) GetZone(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ShippingZoneID(c.Params("id"))

	zone, err := h.svc.GetZone(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(zone)
}

// ListZones handles GET /shipping/zones.
func (h *Handler) ListZones(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	zones, err := h.svc.ListZones(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": zones, "total": len(zones)})
}

type updateZoneRequest struct {
	Name      string   `json:"name"`
	Countries []string `json:"countries"`
	States    []string `json:"states"`
}

// UpdateZone handles PUT /shipping/zones/:id.
func (h *Handler) UpdateZone(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ShippingZoneID(c.Params("id"))

	var req updateZoneRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	zone, err := h.svc.UpdateZone(c.Context(), authCtx.TenantID, id, req.Name, req.Countries, req.States)
	if err != nil {
		return err
	}
	return c.JSON(zone)
}

// DeleteZone handles DELETE /shipping/zones/:id.
func (h *Handler) DeleteZone(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ShippingZoneID(c.Params("id"))

	if err := h.svc.DeleteZone(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ---------------------------------------------------------------------------
// Rate handlers
// ---------------------------------------------------------------------------

type createRateRequest struct {
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	PriceAmount    int64    `json:"price_amount"`
	PriceCurrency  string   `json:"price_currency"`
	MinWeight      *float64 `json:"min_weight,omitempty"`
	MaxWeight      *float64 `json:"max_weight,omitempty"`
	MinOrderAmount *int64   `json:"min_order_amount,omitempty"`
	MaxOrderAmount *int64   `json:"max_order_amount,omitempty"`
	EstDaysMin     *int     `json:"est_days_min,omitempty"`
	EstDaysMax     *int     `json:"est_days_max,omitempty"`
}

// CreateRate handles POST /shipping/zones/:zoneId/rates.
func (h *Handler) CreateRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	zoneID := kernel.ShippingZoneID(c.Params("zoneId"))

	var req createRateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	rate, err := h.svc.CreateRate(c.Context(), authCtx.TenantID, shippingsrv.CreateRateInput{
		ZoneID:         zoneID,
		Name:           req.Name,
		Type:           shipping.RateType(req.Type),
		PriceAmount:    req.PriceAmount,
		PriceCurrency:  req.PriceCurrency,
		MinWeight:      req.MinWeight,
		MaxWeight:      req.MaxWeight,
		MinOrderAmount: req.MinOrderAmount,
		MaxOrderAmount: req.MaxOrderAmount,
		EstDaysMin:     req.EstDaysMin,
		EstDaysMax:     req.EstDaysMax,
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(rate)
}

// GetRate handles GET /shipping/rates/:id.
func (h *Handler) GetRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ShippingRateID(c.Params("id"))

	rate, err := h.svc.GetRate(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(rate)
}

// ListRates handles GET /shipping/zones/:zoneId/rates.
func (h *Handler) ListRates(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	zoneID := kernel.ShippingZoneID(c.Params("zoneId"))

	rates, err := h.svc.ListRates(c.Context(), authCtx.TenantID, zoneID)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"items": rates, "total": len(rates)})
}

type updateRateRequest struct {
	Name           *string  `json:"name,omitempty"`
	Type           *string  `json:"type,omitempty"`
	PriceAmount    *int64   `json:"price_amount,omitempty"`
	PriceCurrency  *string  `json:"price_currency,omitempty"`
	MinWeight      *float64 `json:"min_weight,omitempty"`
	MaxWeight      *float64 `json:"max_weight,omitempty"`
	MinOrderAmount *int64   `json:"min_order_amount,omitempty"`
	MaxOrderAmount *int64   `json:"max_order_amount,omitempty"`
	EstDaysMin     *int     `json:"est_days_min,omitempty"`
	EstDaysMax     *int     `json:"est_days_max,omitempty"`
	Active         *bool    `json:"active,omitempty"`
}

// UpdateRate handles PUT /shipping/rates/:id.
func (h *Handler) UpdateRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ShippingRateID(c.Params("id"))

	var req updateRateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	var rateType *shipping.RateType
	if req.Type != nil {
		t := shipping.RateType(*req.Type)
		rateType = &t
	}

	rate, err := h.svc.UpdateRate(c.Context(), authCtx.TenantID, id, shippingsrv.UpdateRateInput{
		Name:           req.Name,
		Type:           rateType,
		PriceAmount:    req.PriceAmount,
		PriceCurrency:  req.PriceCurrency,
		MinWeight:      req.MinWeight,
		MaxWeight:      req.MaxWeight,
		MinOrderAmount: req.MinOrderAmount,
		MaxOrderAmount: req.MaxOrderAmount,
		EstDaysMin:     req.EstDaysMin,
		EstDaysMax:     req.EstDaysMax,
		Active:         req.Active,
	})
	if err != nil {
		return err
	}
	return c.JSON(rate)
}

// DeleteRate handles DELETE /shipping/rates/:id.
func (h *Handler) DeleteRate(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ShippingRateID(c.Params("id"))

	if err := h.svc.DeleteRate(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ---------------------------------------------------------------------------
// Calculate shipping (public)
// ---------------------------------------------------------------------------

type calculateRequest struct {
	Country     string  `json:"country"`
	State       string  `json:"state"`
	OrderAmount int64   `json:"order_amount"`
	Weight      float64 `json:"weight"`
}

// CalculateShipping handles POST /shipping/calculate.
func (h *Handler) CalculateShipping(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	var req calculateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Country == "" {
		return errx.New("country is required", errx.TypeValidation)
	}

	rates, err := h.svc.CalculateShipping(c.Context(), tenantID, req.Country, req.State, req.OrderAmount, req.Weight)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"rates": rates})
}
