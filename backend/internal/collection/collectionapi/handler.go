package collectionapi

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/collection"
	"github.com/Abraxas-365/hada-commerce/internal/collection/collectionsrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Handler exposes HTTP endpoints for the collection domain.
type Handler struct {
	svc *collectionsrv.Service
}

// NewHandler creates a new Handler.
func NewHandler(svc *collectionsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers admin (authenticated) routes.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/collections")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)

	// Product membership
	g.Post("/:id/products", h.AddProduct)
	g.Delete("/:id/products/:productId", h.RemoveProduct)
	g.Get("/:id/products", h.ListProducts)
	g.Put("/:id/products/reorder", h.ReorderProducts)
}

// RegisterPublicRoutes registers unauthenticated routes (tenant via X-Tenant-ID header).
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/collections")
	g.Get("/", h.ListPublic)
	g.Get("/:slug", h.GetBySlugPublic)
	g.Get("/:slug/products", h.ListProductsPublic)
}

// --------------------------------------------------------------------------
// Admin handlers
// --------------------------------------------------------------------------

type createRequest struct {
	Name            string                      `json:"name"`
	Slug            string                      `json:"slug"`
	Description     string                      `json:"description"`
	ImageURL        string                      `json:"image_url"`
	Type            string                      `json:"type"`
	Rules           []collection.CollectionRule `json:"rules"`
	IsActive        bool                        `json:"is_active"`
	SortOrder       int                         `json:"sort_order"`
	MetaTitle       string                      `json:"meta_title"`
	MetaDescription string                      `json:"meta_description"`
	PublishedAt     *time.Time                  `json:"published_at"`
}

// Create handles POST /collections.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	colType := collection.CollectionType(req.Type)
	if colType == "" {
		colType = collection.CollectionManual
	}

	result, err := h.svc.Create(c.Context(), authCtx.TenantID, collection.CreateInput{
		Name:            req.Name,
		Slug:            req.Slug,
		Description:     req.Description,
		ImageURL:        req.ImageURL,
		Type:            colType,
		Rules:           req.Rules,
		IsActive:        req.IsActive,
		SortOrder:       req.SortOrder,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		PublishedAt:     req.PublishedAt,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// List handles GET /collections (admin — returns all, including inactive).
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	page, pageSize := paginationFromCtx(c)

	activeOnly := c.Query("active_only") == "true"

	result, err := h.svc.List(c.Context(), authCtx.TenantID, activeOnly, page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// GetByID handles GET /collections/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CollectionID(c.Params("id"))

	col, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(col)
}

type updateRequest struct {
	Name            *string                     `json:"name"`
	Slug            *string                     `json:"slug"`
	Description     *string                     `json:"description"`
	ImageURL        *string                     `json:"image_url"`
	Type            *string                     `json:"type"`
	Rules           []collection.CollectionRule `json:"rules"`
	IsActive        *bool                       `json:"is_active"`
	SortOrder       *int                        `json:"sort_order"`
	MetaTitle       *string                     `json:"meta_title"`
	MetaDescription *string                     `json:"meta_description"`
	PublishedAt     *time.Time                  `json:"published_at"`
	ClearPublishedAt bool                       `json:"clear_published_at"`
}

// Update handles PUT /collections/:id.
func (h *Handler) Update(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CollectionID(c.Params("id"))

	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	var colType *collection.CollectionType
	if req.Type != nil {
		t := collection.CollectionType(*req.Type)
		colType = &t
	}

	in := collection.UpdateInput{
		Name:             req.Name,
		Slug:             req.Slug,
		Description:      req.Description,
		ImageURL:         req.ImageURL,
		Type:             colType,
		Rules:            req.Rules,
		IsActive:         req.IsActive,
		SortOrder:        req.SortOrder,
		MetaTitle:        req.MetaTitle,
		MetaDescription:  req.MetaDescription,
		PublishedAt:      req.PublishedAt,
		ClearPublishedAt: req.ClearPublishedAt,
	}

	result, err := h.svc.Update(c.Context(), authCtx.TenantID, id, in)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// Delete handles DELETE /collections/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CollectionID(c.Params("id"))

	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// AddProduct handles POST /collections/:id/products.
func (h *Handler) AddProduct(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	collectionID := c.Params("id")

	var req struct {
		ProductID string `json:"product_id"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.ProductID == "" {
		return errx.New("product_id is required", errx.TypeValidation)
	}

	cp, err := h.svc.AddProduct(c.Context(), authCtx.TenantID, collectionID, req.ProductID, req.SortOrder)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(cp)
}

// RemoveProduct handles DELETE /collections/:id/products/:productId.
func (h *Handler) RemoveProduct(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	collectionID := c.Params("id")
	productID := c.Params("productId")

	if err := h.svc.RemoveProduct(c.Context(), authCtx.TenantID, collectionID, productID); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ListProducts handles GET /collections/:id/products.
func (h *Handler) ListProducts(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	collectionID := c.Params("id")
	page, pageSize := paginationFromCtx(c)

	result, err := h.svc.ListProducts(c.Context(), authCtx.TenantID, collectionID, page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// ReorderProducts handles PUT /collections/:id/products/reorder.
func (h *Handler) ReorderProducts(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	collectionID := c.Params("id")

	var req struct {
		ProductIDs []string `json:"product_ids"`
	}
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if err := h.svc.ReorderProducts(c.Context(), authCtx.TenantID, collectionID, req.ProductIDs); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// --------------------------------------------------------------------------
// Public handlers
// --------------------------------------------------------------------------

// ListPublic handles GET /collections (public — active only).
func (h *Handler) ListPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	page, pageSize := paginationFromCtx(c)
	result, err := h.svc.List(c.Context(), tenantID, true, page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// GetBySlugPublic handles GET /collections/:slug (public).
func (h *Handler) GetBySlugPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	slug := c.Params("slug")
	col, err := h.svc.GetBySlug(c.Context(), tenantID, slug)
	if err != nil {
		return err
	}
	return c.JSON(col)
}

// ListProductsPublic handles GET /collections/:slug/products (public).
func (h *Handler) ListProductsPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	slug := c.Params("slug")
	col, err := h.svc.GetBySlug(c.Context(), tenantID, slug)
	if err != nil {
		return err
	}

	page, pageSize := paginationFromCtx(c)
	result, err := h.svc.ListProducts(c.Context(), tenantID, string(col.ID), page, pageSize)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// --------------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------------

func paginationFromCtx(c *fiber.Ctx) (page, pageSize int) {
	page, _ = strconv.Atoi(c.Query("page"))
	pageSize, _ = strconv.Atoi(c.Query("page_size"))
	return
}
