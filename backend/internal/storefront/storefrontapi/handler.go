package storefrontapi

import (
	"context"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/storefront"
	"github.com/Abraxas-365/vendex/internal/storefront/storefrontsrv"
)

// PageRenderer renders a storefront page as a full HTML5 document.
type PageRenderer interface {
	RenderPage(ctx context.Context, page *storefront.Page) (string, error)
}

// Handler exposes HTTP endpoints for the storefront domain.
type Handler struct {
	svc      *storefrontsrv.Service
	renderer PageRenderer
}

// NewHandler creates a new storefront API handler.
// renderer is optional — pass nil to disable HTML rendering (JSON-only mode).
func NewHandler(svc *storefrontsrv.Service, renderer PageRenderer) *Handler {
	return &Handler{svc: svc, renderer: renderer}
}

// RegisterRoutes registers all storefront admin routes on the given Fiber router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Admin routes
	pages := router.Group("/storefront/pages")
	pages.Post("/", h.CreatePage)
	pages.Get("/", h.ListPages)
	pages.Get("/:id", h.GetPage)
	pages.Put("/:id", h.UpdatePage)
	pages.Post("/:id/publish", h.Publish)
	pages.Post("/:id/unpublish", h.Unpublish)
	pages.Post("/:id/archive", h.Archive)
	pages.Get("/:id/versions", h.ListVersions)
	pages.Get("/:id/versions/:version", h.GetVersion)
	pages.Delete("/:id", h.DeletePage)
	pages.Get("/by-slug/:slug", h.GetPageBySlug)

	// Block type routes (admin)
	blockTypes := router.Group("/storefront/block-types")
	blockTypes.Get("/", h.ListBlockTypes)
	blockTypes.Post("/", h.CreateBlockType)
	blockTypes.Get("/:id", h.GetBlockType)
	blockTypes.Put("/:id", h.UpdateBlockType)
	blockTypes.Delete("/:id", h.DeleteBlockType)
}

// RegisterPublicRoutes registers unauthenticated storefront routes.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Get("/pages/slug/:slug", h.GetPublishedPage)
	router.Get("/pages", h.ListPublishedPages)
}

// --- Public handler ---

// GetPublishedPage handles GET /storefront/pages/by-slug/:slug — public page serving.
//
// If the request Accept header contains "text/html" and a renderer is configured,
// the page is returned as a fully rendered HTML5 document. Otherwise JSON is returned.
func (h *Handler) GetPublishedPage(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID required", errx.TypeAuthorization)
	}

	page, err := h.svc.GetPublishedPage(c.Context(), tenantID, c.Params("slug"))
	if err != nil {
		return err
	}

	// If client wants HTML and we have a renderer, render the page.
	accept := c.Get("Accept")
	if h.renderer != nil && strings.Contains(accept, "text/html") {
		html, err := h.renderer.RenderPage(context.Background(), page)
		if err != nil {
			return errx.Wrap(err, "failed to render page", errx.TypeInternal)
		}
		c.Set("Content-Type", "text/html; charset=utf-8")
		return c.SendString(html)
	}

	return c.JSON(page)
}

// pageNavItem is the minimal page data returned by the public list endpoint.
type pageNavItem struct {
	Slug  string `json:"slug"`
	Title string `json:"title"`
}

// ListPublishedPages handles GET /pages — returns slug+title for all published pages.
// Used by the frontend footer/nav to know which CMS pages exist.
func (h *Handler) ListPublishedPages(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID required", errx.TypeAuthorization)
	}

	status := storefront.PageStatusPublished
	result, err := h.svc.ListPages(c.Context(), tenantID, &status, kernel.PaginationOptions{Page: 1, PageSize: 100})
	if err != nil {
		return err
	}

	items := make([]pageNavItem, 0, len(result.Items))
	for _, p := range result.Items {
		// Skip template pages (underscore-prefixed slugs like _plp, _pdp, _home)
		if strings.HasPrefix(p.Slug, "_") {
			continue
		}
		items = append(items, pageNavItem{Slug: p.Slug, Title: p.Title})
	}
	return c.JSON(items)
}

// --- Admin handlers ---

type createPageRequest struct {
	Slug        string              `json:"slug"`
	Title       string              `json:"title"`
	HTML        string              `json:"html"`
	CSS         string              `json:"css"`
	Meta        storefront.PageMeta `json:"meta"`
	ContentType storefront.ContentType `json:"content_type"`
	Sections    []storefront.Section   `json:"sections"`
	ByAgent     bool                `json:"by_agent"`
}

// CreatePage handles POST /pages.
func (h *Handler) CreatePage(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createPageRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	page, err := h.svc.CreatePage(c.Context(), storefrontsrv.CreatePageInput{
		TenantID:    authCtx.TenantID,
		Slug:        req.Slug,
		Title:       req.Title,
		HTML:        req.HTML,
		CSS:         req.CSS,
		Meta:        req.Meta,
		ContentType: req.ContentType,
		Sections:    req.Sections,
		CreatedBy:   authCtx.UserID.String(),
		ByAgent:     req.ByAgent,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(page)
}

// GetPage handles GET /pages/:id.
func (h *Handler) GetPage(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PageID(c.Params("id"))

	page, err := h.svc.GetPage(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(page)
}

// ListPages handles GET /pages with optional ?status= filter.
func (h *Handler) ListPages(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromCtx(c)

	var status *storefront.PageStatus
	if s := c.Query("status"); s != "" {
		ps := storefront.PageStatus(s)
		status = &ps
	}

	result, err := h.svc.ListPages(c.Context(), authCtx.TenantID, status, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

type updatePageRequest struct {
	Title   *string              `json:"title,omitempty"`
	HTML    *string              `json:"html,omitempty"`
	CSS     *string              `json:"css,omitempty"`
	Meta    *storefront.PageMeta `json:"meta,omitempty"`
	Comment string               `json:"comment"`
}

// UpdatePage handles PUT /pages/:id.
func (h *Handler) UpdatePage(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PageID(c.Params("id"))

	var req updatePageRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	page, err := h.svc.UpdatePage(c.Context(), storefrontsrv.UpdatePageInput{
		TenantID: authCtx.TenantID,
		ID:       id,
		Title:    req.Title,
		HTML:     req.HTML,
		CSS:      req.CSS,
		Meta:     req.Meta,
		EditedBy: authCtx.UserID.String(),
		Comment:  req.Comment,
	})
	if err != nil {
		return err
	}

	return c.JSON(page)
}

// Publish handles POST /pages/:id/publish.
func (h *Handler) Publish(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PageID(c.Params("id"))

	page, err := h.svc.Publish(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(page)
}

// Unpublish handles POST /pages/:id/unpublish.
func (h *Handler) Unpublish(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PageID(c.Params("id"))

	page, err := h.svc.Unpublish(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(page)
}

// Archive handles POST /pages/:id/archive.
func (h *Handler) Archive(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PageID(c.Params("id"))

	page, err := h.svc.Archive(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(page)
}

// ListVersions handles GET /pages/:id/versions.
func (h *Handler) ListVersions(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PageID(c.Params("id"))

	versions, err := h.svc.ListVersions(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(versions)
}

// GetVersion handles GET /pages/:id/versions/:version.
func (h *Handler) GetVersion(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PageID(c.Params("id"))
	version, err := strconv.Atoi(c.Params("version"))
	if err != nil {
		return errx.New("invalid version number", errx.TypeValidation)
	}

	v, err := h.svc.GetVersion(c.Context(), authCtx.TenantID, id, version)
	if err != nil {
		return err
	}

	return c.JSON(v)
}

// DeletePage handles DELETE /pages/:id.
func (h *Handler) DeletePage(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PageID(c.Params("id"))

	if err := h.svc.DeletePage(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetPageBySlug handles GET /pages/by-slug/:slug — admin version (any status).
func (h *Handler) GetPageBySlug(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	slug := c.Params("slug")

	page, err := h.svc.GetPageBySlug(c.Context(), authCtx.TenantID, slug)
	if err != nil {
		return err
	}

	return c.JSON(page)
}

// --- Block Type handlers ---

// ListBlockTypes handles GET /block-types — returns all block types, optionally filtered by ?category=.
func (h *Handler) ListBlockTypes(c *fiber.Ctx) error {
	category := c.Query("category")

	blockTypes, err := h.svc.ListBlockTypes(c.Context(), category)
	if err != nil {
		return err
	}

	return c.JSON(blockTypes)
}

// GetBlockType handles GET /block-types/:id.
func (h *Handler) GetBlockType(c *fiber.Ctx) error {
	id := kernel.BlockTypeID(c.Params("id"))

	bt, err := h.svc.GetBlockType(c.Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(bt)
}

type createBlockTypeRequest struct {
	Name            string          `json:"name"`
	DisplayName     string          `json:"display_name"`
	Category        string          `json:"category"`
	Schema          json.RawMessage `json:"schema"`
	DefaultSettings json.RawMessage `json:"default_settings"`
	Icon            string          `json:"icon"`
}

// CreateBlockType handles POST /block-types.
func (h *Handler) CreateBlockType(c *fiber.Ctx) error {
	var req createBlockTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}
	if req.DisplayName == "" {
		return errx.New("display_name is required", errx.TypeValidation)
	}
	if req.Category == "" {
		return errx.New("category is required", errx.TypeValidation)
	}

	schema := []byte(req.Schema)
	if len(schema) == 0 {
		schema = []byte("{}")
	}
	defaultSettings := []byte(req.DefaultSettings)
	if len(defaultSettings) == 0 {
		defaultSettings = []byte("{}")
	}

	bt, err := h.svc.CreateBlockType(c.Context(), storefrontsrv.CreateBlockTypeInput{
		Name:            req.Name,
		DisplayName:     req.DisplayName,
		Category:        req.Category,
		Schema:          schema,
		DefaultSettings: defaultSettings,
		Icon:            req.Icon,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(bt)
}

type updateBlockTypeRequest struct {
	Name            string          `json:"name"`
	DisplayName     string          `json:"display_name"`
	Category        string          `json:"category"`
	Schema          json.RawMessage `json:"schema"`
	DefaultSettings json.RawMessage `json:"default_settings"`
	Icon            string          `json:"icon"`
}

// UpdateBlockType handles PUT /block-types/:id.
func (h *Handler) UpdateBlockType(c *fiber.Ctx) error {
	id := kernel.BlockTypeID(c.Params("id"))

	var req updateBlockTypeRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	schema := []byte(req.Schema)
	if len(schema) == 0 {
		schema = []byte("{}")
	}
	defaultSettings := []byte(req.DefaultSettings)
	if len(defaultSettings) == 0 {
		defaultSettings = []byte("{}")
	}

	bt, err := h.svc.UpdateBlockType(c.Context(), storefrontsrv.UpdateBlockTypeInput{
		ID:              id,
		Name:            req.Name,
		DisplayName:     req.DisplayName,
		Category:        req.Category,
		Schema:          schema,
		DefaultSettings: defaultSettings,
		Icon:            req.Icon,
	})
	if err != nil {
		return err
	}

	return c.JSON(bt)
}

// DeleteBlockType handles DELETE /block-types/:id.
func (h *Handler) DeleteBlockType(c *fiber.Ctx) error {
	id := kernel.BlockTypeID(c.Params("id"))

	if err := h.svc.DeleteBlockType(c.Context(), id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// --- helpers ---

func paginationFromCtx(c *fiber.Ctx) kernel.PaginationOptions {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	return kernel.PaginationOptions{Page: page, PageSize: pageSize}
}
