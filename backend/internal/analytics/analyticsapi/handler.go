package analyticsapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/analytics/analyticssrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Handler exposes HTTP endpoints for the analytics domain.
type Handler struct {
	svc *analyticssrv.Service
}

// NewHandler creates a new analytics API handler.
func NewHandler(svc *analyticssrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all analytics routes on the given Fiber router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/analytics")
	g.Get("/dashboard", h.GetDashboardStats)
	g.Get("/revenue", h.GetRevenueTimeline)
	g.Get("/top-products", h.GetTopProducts)
	g.Get("/order-status", h.GetOrderStatusBreakdown)
	g.Get("/recent-orders", h.GetRecentOrders)
}

// GetDashboardStats handles GET /analytics/dashboard.
func (h *Handler) GetDashboardStats(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	stats, err := h.svc.GetDashboardStats(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(stats)
}

// GetRevenueTimeline handles GET /analytics/revenue?days=30.
func (h *Handler) GetRevenueTimeline(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	days := 30
	if d, err := strconv.Atoi(c.Query("days")); err == nil && d > 0 {
		days = d
	}

	points, err := h.svc.GetRevenueTimeline(c.Context(), authCtx.TenantID, days)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(points)
}

// GetTopProducts handles GET /analytics/top-products?limit=5.
func (h *Handler) GetTopProducts(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	limit := 5
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	products, err := h.svc.GetTopProducts(c.Context(), authCtx.TenantID, limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

// GetOrderStatusBreakdown handles GET /analytics/order-status.
func (h *Handler) GetOrderStatusBreakdown(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	breakdown, err := h.svc.GetOrderStatusBreakdown(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(breakdown)
}

// GetRecentOrders handles GET /analytics/recent-orders?limit=5.
func (h *Handler) GetRecentOrders(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	limit := 5
	if l, err := strconv.Atoi(c.Query("limit")); err == nil && l > 0 {
		limit = l
	}

	orders, err := h.svc.GetRecentOrders(c.Context(), authCtx.TenantID, limit)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(orders)
}
