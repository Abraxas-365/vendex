package productapi

import (
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/product"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the product domain.
type Handler struct {
	svc *productsrv.Service
}

// NewHandler creates a new product API handler.
func NewHandler(svc *productsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all product routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/products")
	g.Post("/", h.Create)
	g.Get("/:id", h.GetByID)
	g.Get("/", h.List)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
}

// RegisterPublicRoutes registers unauthenticated, read-only product routes.
// Tenant is identified via the X-Tenant-ID header.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/products")
	g.Get("/", h.ListProductsPublic)
	g.Get("/:id", h.GetProductPublic)
}

// createRequest is the JSON body for creating a product.
type createRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	PriceAmount int64    `json:"price_amount"`
	Currency    string   `json:"currency"`
	SKU         string   `json:"sku"`
	Images      []string `json:"images"`
	CategoryID  string   `json:"category_id"`
	Tags        []string `json:"tags"`
	Stock       int      `json:"stock"`
}

// Create handles POST /products.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	p, err := h.svc.Create(c.Context(), authCtx.TenantID, productsrv.CreateInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       kernel.NewMoney(req.PriceAmount, req.Currency),
		SKU:         req.SKU,
		Images:      req.Images,
		CategoryID:  kernel.CategoryID(req.CategoryID),
		Tags:        req.Tags,
		Stock:       req.Stock,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(p)
}

// GetByID handles GET /products/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ProductID(c.Params("id"))

	p, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(p)
}

// List handles GET /products.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// Update handles PUT /products/:id.
func (h *Handler) Update(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ProductID(c.Params("id"))

	existing, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Price = kernel.NewMoney(req.PriceAmount, req.Currency)
	existing.SKU = req.SKU
	existing.Images = req.Images
	existing.CategoryID = kernel.CategoryID(req.CategoryID)
	existing.Tags = req.Tags
	existing.Stock = req.Stock

	if err := h.svc.Update(c.Context(), existing); err != nil {
		return err
	}

	return c.JSON(existing)
}

// Delete handles DELETE /products/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ProductID(c.Params("id"))

	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// --- Public handlers ---

// ListProductsPublic handles GET /products (public, no auth).
// Accepts optional query params: page, page_size, category_id.
func (h *Handler) ListProductsPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	pg := paginationFromQuery(c)
	categoryIDStr := c.Query("category_id")

	if categoryIDStr != "" {
		result, err := h.svc.ListByCategory(c.Context(), tenantID, kernel.CategoryID(categoryIDStr), pg)
		if err != nil {
			return err
		}
		return c.JSON(result)
	}

	result, err := h.svc.List(c.Context(), tenantID, pg)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// GetProductPublic handles GET /products/:id (public, no auth).
func (h *Handler) GetProductPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	id := kernel.ProductID(c.Params("id"))
	p, err := h.svc.GetByID(c.Context(), tenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(p)
}

// --- helpers ---

func paginationFromQuery(c *fiber.Ctx) kernel.PaginationOptions {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	opts := kernel.PaginationOptions{Page: page, PageSize: pageSize}
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	if opts.PageSize > 100 {
		opts.PageSize = 100
	}
	return opts
}

// Ensure product is imported for any future use.
var _ = product.StatusActive
