package storefrontapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/storefront"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontsrv"
)

// Handler exposes HTTP endpoints for the storefront domain.
type Handler struct {
	svc *storefrontsrv.Service
}

// NewHandler creates a new storefront API handler.
func NewHandler(svc *storefrontsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all storefront routes on the given Fiber router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Public routes (no auth needed)
	router.Get("/storefront/:slug", h.GetPublishedPage)

	// Admin routes
	pages := router.Group("/pages")
	pages.Post("/", h.CreatePage)
	pages.Get("/", h.ListPages)
	pages.Get("/:id", h.GetPage)
	pages.Put("/:id", h.UpdatePage)
	pages.Post("/:id/publish", h.Publish)
	pages.Post("/:id/unpublish", h.Unpublish)
	pages.Post("/:id/archive", h.Archive)
	pages.Get("/:id/versions", h.ListVersions)
	pages.Get("/:id/versions/:version", h.GetVersion)
}

// --- Public handler ---

// GetPublishedPage handles GET /storefront/:slug — public page serving.
func (h *Handler) GetPublishedPage(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID required", errx.TypeAuthorization)
	}

	page, err := h.svc.GetPublishedPage(c.Context(), tenantID, c.Params("slug"))
	if err != nil {
		return err
	}

	return c.JSON(page)
}

// --- Admin handlers ---

type createPageRequest struct {
	Slug        string              `json:"slug"`
	Title       string              `json:"title"`
	HTML        string              `json:"html"`
	CSS         string              `json:"css"`
	Meta        storefront.PageMeta `json:"meta"`
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
		TenantID:  authCtx.TenantID,
		Slug:      req.Slug,
		Title:     req.Title,
		HTML:      req.HTML,
		CSS:       req.CSS,
		Meta:      req.Meta,
		CreatedBy: string(authCtx.UserID),
		ByAgent:   req.ByAgent,
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
		EditedBy: string(authCtx.UserID),
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

// --- helpers ---

func paginationFromCtx(c *fiber.Ctx) kernel.PaginationOptions {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	return kernel.PaginationOptions{Page: page, PageSize: pageSize}
}
