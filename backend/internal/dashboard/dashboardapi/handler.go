package dashboardapi

import (
	"strconv"
	"time"

	"github.com/Abraxas-365/vendex/internal/dashboard"
	"github.com/Abraxas-365/vendex/internal/dashboard/dashboardsrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the dashboard reporting domain.
type Handler struct {
	svc *dashboardsrv.Service
}

// NewHandler creates a new dashboard API handler.
func NewHandler(svc *dashboardsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all dashboard routes on the given router.
// All routes require authentication (caller is responsible for auth middleware).
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/dashboard")
	g.Get("/sales", h.GetSalesOverview)
	g.Get("/top-products", h.GetTopProducts)
	g.Get("/revenue", h.GetRevenueByDay)
	g.Get("/customers", h.GetCustomerStats)
	g.Get("/funnel", h.GetConversionFunnel)
}

// GetSalesOverview handles GET /dashboard/sales?from=2024-01-01&to=2024-12-31
func (h *Handler) GetSalesOverview(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	dr, err := parseDateRange(c)
	if err != nil {
		return err
	}

	overview, err := h.svc.GetSalesOverview(c.Context(), authCtx.TenantID, dr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(overview)
}

// GetTopProducts handles GET /dashboard/top-products?limit=10&from=...&to=...
func (h *Handler) GetTopProducts(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	dr, err := parseDateRange(c)
	if err != nil {
		return err
	}

	limit := 10
	if l, e := strconv.Atoi(c.Query("limit")); e == nil && l > 0 {
		limit = l
	}

	products, err := h.svc.GetTopProducts(c.Context(), authCtx.TenantID, dr, limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

// GetRevenueByDay handles GET /dashboard/revenue?from=...&to=...
func (h *Handler) GetRevenueByDay(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	dr, err := parseDateRange(c)
	if err != nil {
		return err
	}

	points, err := h.svc.GetRevenueByDay(c.Context(), authCtx.TenantID, dr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(points)
}

// GetCustomerStats handles GET /dashboard/customers?from=...&to=...
func (h *Handler) GetCustomerStats(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	dr, err := parseDateRange(c)
	if err != nil {
		return err
	}

	stats, err := h.svc.GetCustomerStats(c.Context(), authCtx.TenantID, dr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

// GetConversionFunnel handles GET /dashboard/funnel?from=...&to=...
func (h *Handler) GetConversionFunnel(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	dr, err := parseDateRange(c)
	if err != nil {
		return err
	}

	funnel, err := h.svc.GetConversionFunnel(c.Context(), authCtx.TenantID, dr)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(funnel)
}

// ── helpers ───────────────────────────────────────────────────────────────────

// parseDateRange extracts optional ?from= and ?to= query parameters.
// Both values are in YYYY-MM-DD format. If omitted, zero values are returned
// and the service layer will apply sensible defaults.
func parseDateRange(c *fiber.Ctx) (dashboard.DateRange, error) {
	var dr dashboard.DateRange

	if raw := c.Query("from"); raw != "" {
		t, err := time.Parse("2006-01-02", raw)
		if err != nil {
			return dr, errx.New("invalid 'from' date, expected YYYY-MM-DD", errx.TypeValidation)
		}
		dr.From = t.UTC()
	}

	if raw := c.Query("to"); raw != "" {
		t, err := time.Parse("2006-01-02", raw)
		if err != nil {
			return dr, errx.New("invalid 'to' date, expected YYYY-MM-DD", errx.TypeValidation)
		}
		// Include the entire "to" day.
		dr.To = t.UTC().Add(24*time.Hour - time.Nanosecond)
	}

	return dr, nil
}
