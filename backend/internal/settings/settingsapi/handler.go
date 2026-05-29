package settingsapi

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/settings"
	"github.com/Abraxas-365/hada-commerce/internal/settings/settingssrv"
)

// Handler exposes HTTP endpoints for the settings domain.
type Handler struct {
	svc *settingssrv.Service
}

// NewHandler creates a new settings API handler.
func NewHandler(svc *settingssrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all settings routes on the given Fiber router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/settings")
	g.Get("/", h.Get)
	g.Put("/", h.Update)
}

// updateRequest is the JSON body for updating store settings.
type updateRequest struct {
	StoreName      string                  `json:"store_name"`
	StoreEmail     string                  `json:"store_email"`
	StorePhone     string                  `json:"store_phone"`
	Currency       string                  `json:"currency"`
	Timezone       string                  `json:"timezone"`
	Address        settings.StoreAddress   `json:"address"`
	LogoURL        string                  `json:"logo_url"`
	FaviconURL     string                  `json:"favicon_url"`
	SocialLinks    settings.SocialLinks    `json:"social_links"`
	CheckoutConfig settings.CheckoutConfig `json:"checkout_config"`
}

// Get handles GET /settings.
// Returns the current settings for the tenant, creating defaults if none exist.
func (h *Handler) Get(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	ss, err := h.svc.Get(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ss)
}

// Update handles PUT /settings.
// Upserts settings for the tenant from the request body.
func (h *Handler) Update(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	ss, err := h.svc.Update(c.Context(), authCtx.TenantID, settingssrv.UpdateInput{
		StoreName:      req.StoreName,
		StoreEmail:     req.StoreEmail,
		StorePhone:     req.StorePhone,
		Currency:       req.Currency,
		Timezone:       req.Timezone,
		Address:        req.Address,
		LogoURL:        req.LogoURL,
		FaviconURL:     req.FaviconURL,
		SocialLinks:    req.SocialLinks,
		CheckoutConfig: req.CheckoutConfig,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(ss)
}
