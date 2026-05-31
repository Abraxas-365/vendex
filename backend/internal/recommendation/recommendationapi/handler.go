package recommendationapi

import (
	"strconv"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/recommendation"
	"github.com/Abraxas-365/vendex/internal/recommendation/recommendationsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the recommendation domain.
type Handler struct {
	svc *recommendationsrv.Service
}

// New creates a new recommendation API handler.
func New(svc *recommendationsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// ---------------------------------------------------------------------------
// Route registration
// ---------------------------------------------------------------------------

// RegisterRoutes registers protected (admin) recommendation routes.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/recommendations")
	g.Get("/rules", h.ListRules)
	g.Post("/rules", h.CreateRule)
	g.Put("/rules/:id", h.UpdateRule)
	g.Delete("/rules/:id", h.DeleteRule)
}

// RegisterPublicRoutes registers public recommendation routes.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/recommendations")
	g.Post("/track/view", h.TrackView)
	g.Post("/track/interaction", h.TrackInteraction)
	g.Get("/product/:productId", h.GetForProduct)
	g.Get("/trending", h.GetTrending)
	g.Get("/recently-viewed", h.GetRecentlyViewed)
	g.Get("/personalized", h.GetPersonalized)
}

// ---------------------------------------------------------------------------
// Tracking handlers
// ---------------------------------------------------------------------------

// TrackView handles POST /recommendations/track/view
func (h *Handler) TrackView(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.Unauthorized("tenant id is required")
	}

	var input recommendation.TrackViewInput
	if err := c.BodyParser(&input); err != nil {
		return errx.Validation("invalid request body")
	}

	if err := h.svc.TrackView(c.Context(), tenantID, input); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// TrackInteraction handles POST /recommendations/track/interaction
func (h *Handler) TrackInteraction(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.Unauthorized("tenant id is required")
	}

	var input recommendation.TrackInteractionInput
	if err := c.BodyParser(&input); err != nil {
		return errx.Validation("invalid request body")
	}

	if err := h.svc.TrackInteraction(c.Context(), tenantID, input); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ---------------------------------------------------------------------------
// Recommendation handlers
// ---------------------------------------------------------------------------

// GetForProduct handles GET /recommendations/product/:productId
func (h *Handler) GetForProduct(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.Unauthorized("tenant id is required")
	}

	productID := c.Params("productId")
	if productID == "" {
		return errx.Validation("product_id is required")
	}

	recommendationType := c.Query("type", "frequently_bought_together")
	limit := queryInt(c, "limit", 10)

	result, err := h.svc.GetRecommendations(c.Context(), tenantID, recommendation.GetRecommendationsInput{
		ProductID: productID,
		VisitorID: c.Query("visitor_id"),
		Type:      recommendationType,
		Limit:     limit,
	})
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": result})
}

// GetTrending handles GET /recommendations/trending
func (h *Handler) GetTrending(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.Unauthorized("tenant id is required")
	}

	limit := queryInt(c, "limit", 10)
	sinceDays := queryInt(c, "since_days", 7)
	since := time.Duration(sinceDays) * 24 * time.Hour

	result, err := h.svc.GetTrending(c.Context(), tenantID, limit, since)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": result})
}

// GetRecentlyViewed handles GET /recommendations/recently-viewed
func (h *Handler) GetRecentlyViewed(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.Unauthorized("tenant id is required")
	}

	visitorID := c.Query("visitor_id")
	if visitorID == "" {
		return errx.Validation("visitor_id query parameter is required")
	}

	limit := queryInt(c, "limit", 10)

	result, err := h.svc.GetRecentlyViewed(c.Context(), tenantID, visitorID, limit)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": result})
}

// GetPersonalized handles GET /recommendations/personalized
func (h *Handler) GetPersonalized(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.Unauthorized("tenant id is required")
	}

	visitorID := c.Query("visitor_id")
	if visitorID == "" {
		return errx.Validation("visitor_id query parameter is required")
	}

	limit := queryInt(c, "limit", 10)

	result, err := h.svc.GetPersonalized(c.Context(), tenantID, visitorID, limit)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": result})
}

// ---------------------------------------------------------------------------
// Rule handlers
// ---------------------------------------------------------------------------

// ListRules handles GET /recommendations/rules
func (h *Handler) ListRules(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	rules, err := h.svc.ListRules(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": rules})
}

// CreateRule handles POST /recommendations/rules
func (h *Handler) CreateRule(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var input recommendation.CreateRuleInput
	if err := c.BodyParser(&input); err != nil {
		return errx.Validation("invalid request body")
	}

	rule, err := h.svc.CreateRule(c.Context(), authCtx.TenantID, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(rule)
}

// UpdateRule handles PUT /recommendations/rules/:id
func (h *Handler) UpdateRule(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewRecommendationRuleID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("rule id is required")
	}

	var input recommendation.UpdateRuleInput
	if err := c.BodyParser(&input); err != nil {
		return errx.Validation("invalid request body")
	}

	rule, err := h.svc.UpdateRule(c.Context(), authCtx.TenantID, id, input)
	if err != nil {
		return err
	}

	return c.JSON(rule)
}

// DeleteRule handles DELETE /recommendations/rules/:id
func (h *Handler) DeleteRule(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewRecommendationRuleID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("rule id is required")
	}

	if err := h.svc.DeleteRule(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func queryInt(c *fiber.Ctx, key string, defaultVal int) int {
	v, err := strconv.Atoi(c.Query(key))
	if err != nil || v <= 0 {
		return defaultVal
	}
	return v
}
