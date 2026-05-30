package i18napi

import (
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/i18n/i18nsrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the i18n domain.
type Handler struct {
	svc *i18nsrv.Service
}

// NewHandler creates a new i18n API handler.
func NewHandler(svc *i18nsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected admin routes for translation management.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/i18n")

	// Supported locales (static, no tenant needed)
	g.Get("/supported-locales", h.ListSupportedLocales)

	// Entity translation management
	g.Put("/:entityType/:entityId/:locale", h.SetTranslations)
	g.Get("/:entityType/:entityId/:locale", h.GetTranslations)
	g.Get("/:entityType/:entityId/locales", h.ListLocales)
	g.Delete("/:entityType/:entityId/:locale/:field", h.DeleteTranslation)
	g.Delete("/:entityType/:entityId", h.DeleteAllTranslations)
}

// RegisterPublicRoutes registers public read-only routes for translation retrieval.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/i18n")
	g.Get("/supported-locales", h.ListSupportedLocalesPublic)
	g.Get("/:entityType/:entityId/:locale", h.GetTranslationsPublic)
	g.Get("/:entityType/:entityId/locales", h.ListLocalesPublic)
}

// ============================================================================
// Request/Response types
// ============================================================================

type setTranslationsRequest struct {
	Fields map[string]string `json:"fields"`
}

// ============================================================================
// Admin handlers
// ============================================================================

// SetTranslations handles PUT /i18n/:entityType/:entityId/:locale.
// Sets (upserts) all provided field translations for the given entity+locale.
func (h *Handler) SetTranslations(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	entityType := c.Params("entityType")
	entityID := c.Params("entityId")
	locale := c.Params("locale")

	var req setTranslationsRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if len(req.Fields) == 0 {
		return errx.New("fields map must not be empty", errx.TypeValidation)
	}

	if err := h.svc.SetTranslations(c.Context(), authCtx.TenantID, entityType, entityID, locale, req.Fields); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "translations saved"})
}

// GetTranslations handles GET /i18n/:entityType/:entityId/:locale.
// Returns the translation bundle for the given entity+locale.
func (h *Handler) GetTranslations(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	entityType := c.Params("entityType")
	entityID := c.Params("entityId")
	locale := c.Params("locale")

	bundle, err := h.svc.GetTranslations(c.Context(), authCtx.TenantID, entityType, entityID, locale)
	if err != nil {
		return err
	}

	return c.JSON(bundle)
}

// ListLocales handles GET /i18n/:entityType/:entityId/locales.
// Returns available locales for the given entity.
func (h *Handler) ListLocales(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	entityType := c.Params("entityType")
	entityID := c.Params("entityId")

	locales, err := h.svc.ListLocales(c.Context(), authCtx.TenantID, entityType, entityID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"locales": locales})
}

// DeleteTranslation handles DELETE /i18n/:entityType/:entityId/:locale/:field.
// Removes a single translated field.
func (h *Handler) DeleteTranslation(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	entityType := c.Params("entityType")
	entityID := c.Params("entityId")
	locale := c.Params("locale")
	field := c.Params("field")

	if err := h.svc.DeleteTranslation(c.Context(), authCtx.TenantID, entityType, entityID, locale, field); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// DeleteAllTranslations handles DELETE /i18n/:entityType/:entityId.
// Removes all translations for an entity (all locales, all fields).
func (h *Handler) DeleteAllTranslations(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	entityType := c.Params("entityType")
	entityID := c.Params("entityId")

	if err := h.svc.DeleteAllTranslations(c.Context(), authCtx.TenantID, entityType, entityID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListSupportedLocales handles GET /i18n/supported-locales (admin).
func (h *Handler) ListSupportedLocales(c *fiber.Ctx) error {
	locales := h.svc.ListSupportedLocales()
	return c.JSON(locales)
}

// ============================================================================
// Public handlers
// ============================================================================

// GetTranslationsPublic handles GET /i18n/:entityType/:entityId/:locale (public).
// Tenant identified via X-Tenant-ID header.
func (h *Handler) GetTranslationsPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	entityType := c.Params("entityType")
	entityID := c.Params("entityId")
	locale := c.Params("locale")

	bundle, err := h.svc.GetTranslations(c.Context(), tenantID, entityType, entityID, locale)
	if err != nil {
		return err
	}

	return c.JSON(bundle)
}

// ListLocalesPublic handles GET /i18n/:entityType/:entityId/locales (public).
// Tenant identified via X-Tenant-ID header.
func (h *Handler) ListLocalesPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	entityType := c.Params("entityType")
	entityID := c.Params("entityId")

	locales, err := h.svc.ListLocales(c.Context(), tenantID, entityType, entityID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"locales": locales})
}

// ListSupportedLocalesPublic handles GET /i18n/supported-locales (public).
func (h *Handler) ListSupportedLocalesPublic(c *fiber.Ctx) error {
	locales := h.svc.ListSupportedLocales()
	return c.JSON(locales)
}
