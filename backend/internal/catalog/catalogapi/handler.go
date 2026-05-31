package catalogapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/catalog/catalogsrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Handler exposes HTTP endpoints for the catalog domain.
type Handler struct {
	svc *catalogsrv.Service
}

// NewHandler creates a new catalog API handler.
func NewHandler(svc *catalogsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all catalog routes on the given Fiber router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// Categories
	cat := router.Group("/categories")
	cat.Post("/", h.CreateCategory)
	cat.Get("/:id", h.GetCategory)
	cat.Get("/", h.ListCategories)
	cat.Delete("/:id", h.DeleteCategory)

	// Collections
	col := router.Group("/collections")
	col.Post("/", h.CreateCollection)
	col.Get("/:id", h.GetCollection)
	col.Get("/", h.ListCollections)
	col.Delete("/:id", h.DeleteCollection)
}

// RegisterPublicRoutes registers unauthenticated, read-only catalog routes.
// Tenant is identified via the X-Tenant-ID header.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Get("/categories", h.ListCategoriesPublic)
	router.Get("/collections", h.ListCollectionsPublic)
	router.Get("/collections/:id", h.GetCollectionPublic)
}

// --- Category handlers ---

type createCategoryRequest struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	ParentID    *string `json:"parent_id,omitempty"`
	Description string  `json:"description"`
}

// CreateCategory handles POST /categories.
func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	var parentID *kernel.CategoryID
	if req.ParentID != nil {
		pid := kernel.CategoryID(*req.ParentID)
		parentID = &pid
	}

	cat, err := h.svc.CreateCategory(c.Context(), authCtx.TenantID, catalogsrv.CreateCategoryInput{
		Name:        req.Name,
		Slug:        req.Slug,
		ParentID:    parentID,
		Description: req.Description,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(cat)
}

// GetCategory handles GET /categories/:id.
func (h *Handler) GetCategory(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CategoryID(c.Params("id"))

	cat, err := h.svc.GetCategoryByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(cat)
}

// ListCategories handles GET /categories.
func (h *Handler) ListCategories(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromCtx(c)

	result, err := h.svc.ListCategories(c.Context(), authCtx.TenantID, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// DeleteCategory handles DELETE /categories/:id.
func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CategoryID(c.Params("id"))

	if err := h.svc.DeleteCategory(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// --- Collection handlers ---

type createCollectionRequest struct {
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	ProductIDs  []string       `json:"product_ids"`
	IsAutomatic bool           `json:"is_automatic"`
	Rules       map[string]any `json:"rules"`
}

// CreateCollection handles POST /collections.
func (h *Handler) CreateCollection(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createCollectionRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	productIDs := make([]kernel.ProductID, len(req.ProductIDs))
	for i, id := range req.ProductIDs {
		productIDs[i] = kernel.ProductID(id)
	}

	col, err := h.svc.CreateCollection(c.Context(), authCtx.TenantID, catalogsrv.CreateCollectionInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ProductIDs:  productIDs,
		IsAutomatic: req.IsAutomatic,
		Rules:       req.Rules,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(col)
}

// GetCollection handles GET /collections/:id.
func (h *Handler) GetCollection(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CollectionID(c.Params("id"))

	col, err := h.svc.GetCollectionByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(col)
}

// ListCollections handles GET /collections.
func (h *Handler) ListCollections(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromCtx(c)

	result, err := h.svc.ListCollections(c.Context(), authCtx.TenantID, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// DeleteCollection handles DELETE /collections/:id.
func (h *Handler) DeleteCollection(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CollectionID(c.Params("id"))

	if err := h.svc.DeleteCollection(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// --- Public handlers ---

// ListCategoriesPublic handles GET /categories (public, no auth).
func (h *Handler) ListCategoriesPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	pg := paginationFromCtx(c)
	result, err := h.svc.ListCategories(c.Context(), tenantID, pg)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// ListCollectionsPublic handles GET /collections (public, no auth).
func (h *Handler) ListCollectionsPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	pg := paginationFromCtx(c)
	result, err := h.svc.ListCollections(c.Context(), tenantID, pg)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// GetCollectionPublic handles GET /collections/:id (public, no auth).
func (h *Handler) GetCollectionPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	id := kernel.CollectionID(c.Params("id"))
	col, err := h.svc.GetCollectionByID(c.Context(), tenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(col)
}

// --- helpers ---

func paginationFromCtx(c *fiber.Ctx) kernel.PaginationOptions {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	return kernel.PaginationOptions{Page: page, PageSize: pageSize}
}
