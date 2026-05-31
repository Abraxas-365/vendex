package marketplaceapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/marketplace"
	"github.com/Abraxas-365/vendex/internal/marketplace/marketplacesrv"
)

// PresetHandler exposes HTTP endpoints for the preset domain.
type PresetHandler struct {
	svc *marketplacesrv.PresetService
}

// NewPresetHandler creates a new PresetHandler.
func NewPresetHandler(svc *marketplacesrv.PresetService) *PresetHandler {
	return &PresetHandler{svc: svc}
}

// RegisterPresetRoutes registers all preset routes on the given Fiber router.
func (h *PresetHandler) RegisterPresetRoutes(router fiber.Router) {
	// Marketplace browsing (public presets)
	marketplace := router.Group("/marketplace/presets")
	marketplace.Get("/", h.ListMarketplace)
	marketplace.Get("/slug/:slug", h.GetBySlug)
	marketplace.Get("/:id", h.GetPreset)

	// Tenant-owned presets (CRUD)
	my := router.Group("/marketplace/my-presets")
	my.Post("/", h.CreatePreset)
	my.Get("/", h.ListMyPresets)
	my.Put("/:id", h.UpdatePreset)
	my.Delete("/:id", h.DeletePreset)
	my.Post("/:id/publish", h.PublishPreset)
	my.Post("/:id/archive", h.ArchivePreset)

	// Installations
	installs := router.Group("/marketplace/installations")
	installs.Get("/", h.ListInstalled)
	installs.Post("/", h.InstallPreset)
	installs.Delete("/:preset_id", h.UninstallPreset)
}

// ListMarketplace handles GET /marketplace/presets — browse public presets.
func (h *PresetHandler) ListMarketplace(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	search := c.Query("search", "")

	opts := marketplace.PresetListOptions{
		PaginationOptions: kernel.PaginationOptions{Page: page, PageSize: pageSize},
		Search:            search,
	}

	result, err := h.svc.ListMarketplace(c.Context(), opts)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// GetPreset handles GET /marketplace/presets/:id.
func (h *PresetHandler) GetPreset(c *fiber.Ctx) error {
	id := kernel.PresetID(c.Params("id"))
	p, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(p)
}

// GetBySlug handles GET /marketplace/presets/slug/:slug.
func (h *PresetHandler) GetBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")
	p, err := h.svc.GetBySlug(c.Context(), slug)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(p)
}

// CreatePreset handles POST /marketplace/my-presets.
func (h *PresetHandler) CreatePreset(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	var req marketplace.CreatePresetRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	p, err := h.svc.Create(c.Context(), auth.TenantID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(p)
}

// ListMyPresets handles GET /marketplace/my-presets.
func (h *PresetHandler) ListMyPresets(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	result, err := h.svc.ListByTenant(c.Context(), auth.TenantID, pagination(c))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// UpdatePreset handles PUT /marketplace/my-presets/:id.
func (h *PresetHandler) UpdatePreset(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	id := kernel.PresetID(c.Params("id"))
	var req marketplace.UpdatePresetRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	p, err := h.svc.Update(c.Context(), auth.TenantID, id, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(p)
}

// DeletePreset handles DELETE /marketplace/my-presets/:id.
func (h *PresetHandler) DeletePreset(c *fiber.Ctx) error {
	id := kernel.PresetID(c.Params("id"))
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// PublishPreset handles POST /marketplace/my-presets/:id/publish.
func (h *PresetHandler) PublishPreset(c *fiber.Ctx) error {
	id := kernel.PresetID(c.Params("id"))
	p, err := h.svc.Publish(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(p)
}

// ArchivePreset handles POST /marketplace/my-presets/:id/archive.
func (h *PresetHandler) ArchivePreset(c *fiber.Ctx) error {
	id := kernel.PresetID(c.Params("id"))
	p, err := h.svc.Archive(c.Context(), id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(p)
}

// InstallPreset handles POST /marketplace/installations.
func (h *PresetHandler) InstallPreset(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	var body struct {
		PresetID string `json:"preset_id"`
		Config   []byte `json:"config"`
	}
	if err := c.BodyParser(&body); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	install, err := h.svc.Install(c.Context(), auth.TenantID, kernel.PresetID(body.PresetID), body.Config)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(install)
}

// UninstallPreset handles DELETE /marketplace/installations/:preset_id.
func (h *PresetHandler) UninstallPreset(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	presetID := kernel.PresetID(c.Params("preset_id"))
	if err := h.svc.Uninstall(c.Context(), auth.TenantID, presetID); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ListInstalled handles GET /marketplace/installations.
func (h *PresetHandler) ListInstalled(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	result, err := h.svc.ListInstalled(c.Context(), auth.TenantID, pagination(c))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(result)
}
