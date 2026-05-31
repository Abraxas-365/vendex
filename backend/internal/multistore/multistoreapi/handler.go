package multistoreapi

import (
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/multistore"
	"github.com/Abraxas-365/hada-commerce/internal/multistore/multistoresrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the multistore domain.
type Handler struct {
	svc *multistoresrv.Service
}

// NewHandler creates a new multistore API handler.
func NewHandler(svc *multistoresrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected admin routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/storefronts")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
	g.Put("/:id/default", h.SetDefault)
	g.Post("/:id/catalogs", h.AddCatalog)
	g.Delete("/:id/catalogs/:catalogId", h.RemoveCatalog)
	g.Get("/:id/catalogs", h.ListCatalogs)
}

// RegisterPublicRoutes registers unauthenticated routes (slug/domain resolution).
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/storefronts")
	g.Get("/by-slug/:slug", h.GetBySlug)
	g.Get("/by-domain/:domain", h.GetByDomain)
}

// ---------------------------------------------------------------------------
// Admin handlers
// ---------------------------------------------------------------------------

// Create handles POST /storefronts.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var body struct {
		Name            string                 `json:"name"`
		Slug            string                 `json:"slug"`
		Domain          *string                `json:"domain"`
		Description     string                 `json:"description"`
		ThemeID         string                 `json:"theme_id"`
		LogoURL         string                 `json:"logo_url"`
		DefaultLocale   string                 `json:"default_locale"`
		DefaultCurrency string                 `json:"default_currency"`
		IsActive        bool                   `json:"is_active"`
		Settings        map[string]interface{} `json:"settings"`
	}
	if err := c.BodyParser(&body); err != nil {
		return errx.Validation("invalid request body")
	}

	input := multistore.CreateInput{
		Name:            body.Name,
		Slug:            body.Slug,
		Domain:          body.Domain,
		Description:     body.Description,
		ThemeID:         body.ThemeID,
		LogoURL:         body.LogoURL,
		DefaultLocale:   body.DefaultLocale,
		DefaultCurrency: body.DefaultCurrency,
		IsActive:        body.IsActive,
		Settings:        body.Settings,
	}

	sf, err := h.svc.Create(c.Context(), authCtx.TenantID, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(sf)
}

// List handles GET /storefronts.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	page, pageSize := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetByID handles GET /storefronts/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewStorefrontEntryID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("storefront id is required")
	}

	sf, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(sf)
}

// Update handles PUT /storefronts/:id.
func (h *Handler) Update(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewStorefrontEntryID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("storefront id is required")
	}

	var body struct {
		Name            *string                `json:"name"`
		Domain          *string                `json:"domain"`
		Description     *string                `json:"description"`
		ThemeID         *string                `json:"theme_id"`
		LogoURL         *string                `json:"logo_url"`
		DefaultLocale   *string                `json:"default_locale"`
		DefaultCurrency *string                `json:"default_currency"`
		IsActive        *bool                  `json:"is_active"`
		Settings        map[string]interface{} `json:"settings"`
	}
	if err := c.BodyParser(&body); err != nil {
		return errx.Validation("invalid request body")
	}

	input := multistore.UpdateInput{
		Name:            body.Name,
		Domain:          body.Domain,
		Description:     body.Description,
		ThemeID:         body.ThemeID,
		LogoURL:         body.LogoURL,
		DefaultLocale:   body.DefaultLocale,
		DefaultCurrency: body.DefaultCurrency,
		IsActive:        body.IsActive,
		Settings:        body.Settings,
	}

	sf, err := h.svc.Update(c.Context(), authCtx.TenantID, id, input)
	if err != nil {
		return err
	}

	return c.JSON(sf)
}

// Delete handles DELETE /storefronts/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewStorefrontEntryID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("storefront id is required")
	}

	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// SetDefault handles PUT /storefronts/:id/default.
func (h *Handler) SetDefault(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewStorefrontEntryID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("storefront id is required")
	}

	if err := h.svc.SetDefault(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// AddCatalog handles POST /storefronts/:id/catalogs.
func (h *Handler) AddCatalog(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	storefrontID := kernel.NewStorefrontEntryID(c.Params("id"))

	if storefrontID.IsEmpty() {
		return errx.Validation("storefront id is required")
	}

	var body struct {
		CatalogID string `json:"catalog_id"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.BodyParser(&body); err != nil {
		return errx.Validation("invalid request body")
	}
	if body.CatalogID == "" {
		return errx.Validation("catalog_id is required")
	}

	sc, err := h.svc.AddCatalog(c.Context(), authCtx.TenantID, storefrontID, body.CatalogID, body.SortOrder)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(sc)
}

// RemoveCatalog handles DELETE /storefronts/:id/catalogs/:catalogId.
func (h *Handler) RemoveCatalog(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	storefrontID := kernel.NewStorefrontEntryID(c.Params("id"))
	catalogID := c.Params("catalogId")

	if storefrontID.IsEmpty() {
		return errx.Validation("storefront id is required")
	}
	if catalogID == "" {
		return errx.Validation("catalog id is required")
	}

	if err := h.svc.RemoveCatalog(c.Context(), authCtx.TenantID, storefrontID, catalogID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListCatalogs handles GET /storefronts/:id/catalogs.
func (h *Handler) ListCatalogs(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	storefrontID := kernel.NewStorefrontEntryID(c.Params("id"))

	if storefrontID.IsEmpty() {
		return errx.Validation("storefront id is required")
	}

	catalogs, err := h.svc.ListCatalogs(c.Context(), authCtx.TenantID, storefrontID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": catalogs, "total": len(catalogs)})
}

// ---------------------------------------------------------------------------
// Public handlers
// ---------------------------------------------------------------------------

// GetBySlug handles GET /storefronts/by-slug/:slug.
func (h *Handler) GetBySlug(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.Validation("X-Tenant-ID header is required")
	}

	slug := c.Params("slug")
	if slug == "" {
		return errx.Validation("slug is required")
	}

	sf, err := h.svc.GetBySlug(c.Context(), tenantID, slug)
	if err != nil {
		return err
	}

	return c.JSON(sf)
}

// GetByDomain handles GET /storefronts/by-domain/:domain.
func (h *Handler) GetByDomain(c *fiber.Ctx) error {
	domain := c.Params("domain")
	if domain == "" {
		return errx.Validation("domain is required")
	}

	sf, err := h.svc.GetByDomain(c.Context(), domain)
	if err != nil {
		return err
	}

	return c.JSON(sf)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func paginationFromQuery(c *fiber.Ctx) (page, pageSize int) {
	page, _ = strconv.Atoi(c.Query("page"))
	pageSize, _ = strconv.Atoi(c.Query("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return
}
