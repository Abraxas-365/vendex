package productapi

import (
	"strconv"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/product"
	"github.com/Abraxas-365/vendex/internal/product/productsrv"
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
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)

	// Option routes
	g.Post("/:id/options", h.CreateOption)
	g.Get("/:id/options", h.ListOptions)
	g.Put("/options/:optionId", h.UpdateOption)
	g.Delete("/options/:optionId", h.DeleteOption)

	// Variant routes
	g.Post("/:id/variants", h.CreateVariant)
	g.Get("/:id/variants", h.ListVariants)
	g.Get("/variants/:variantId", h.GetVariant)
	g.Put("/variants/:variantId", h.UpdateVariant)
	g.Delete("/variants/:variantId", h.DeleteVariant)
}

// RegisterPublicRoutes registers unauthenticated, read-only product routes.
// Tenant is identified via the X-Tenant-ID header.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/products")
	g.Get("/", h.ListProductsPublic)
	g.Get("/slug/:slug", h.GetProductBySlugPublic)
	g.Get("/:id", h.GetProductPublic)
	g.Get("/:id/options", h.ListOptionsPublic)
	g.Get("/:id/variants", h.ListVariantsPublic)
}

// createRequest is the JSON body for creating a product.
type createRequest struct {
	Name            string   `json:"name"`
	Description     string   `json:"description"`
	PriceAmount     int64    `json:"price_amount"`
	Currency        string   `json:"currency"`
	SKU             string   `json:"sku"`
	Images          []string `json:"images"`
	CategoryID      string   `json:"category_id"`
	Tags            []string `json:"tags"`
	Stock           int      `json:"stock"`
	Slug            string   `json:"slug"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
	CanonicalURL    string   `json:"canonical_url"`
}

// Create handles POST /products.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	p, err := h.svc.Create(c.Context(), authCtx.TenantID, productsrv.CreateInput{
		Name:            req.Name,
		Description:     req.Description,
		Price:           kernel.NewMoney(req.PriceAmount, req.Currency),
		SKU:             req.SKU,
		Images:          req.Images,
		CategoryID:      kernel.CategoryID(req.CategoryID),
		Tags:            req.Tags,
		Stock:           req.Stock,
		Slug:            req.Slug,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
		CanonicalURL:    req.CanonicalURL,
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
	if req.MetaTitle != "" {
		existing.MetaTitle = req.MetaTitle
	}
	if req.MetaDescription != "" {
		existing.MetaDescription = req.MetaDescription
	}
	if req.Slug != "" {
		existing.Slug = req.Slug
	}
	if req.CanonicalURL != "" {
		existing.CanonicalURL = req.CanonicalURL
	}

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

// ─── Option handlers ──────────────────────────────────────────────────────────

type createOptionRequest struct {
	Name     string   `json:"name"`
	Position int      `json:"position"`
	Values   []string `json:"values"`
}

// CreateOption handles POST /products/:id/options.
func (h *Handler) CreateOption(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	productID := kernel.ProductID(c.Params("id"))

	var req createOptionRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("option name is required", errx.TypeValidation)
	}

	opt, err := h.svc.CreateOption(c.Context(), authCtx.TenantID, productID, req.Name, req.Position, req.Values)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(opt)
}

// ListOptions handles GET /products/:id/options.
func (h *Handler) ListOptions(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	productID := kernel.ProductID(c.Params("id"))

	opts, err := h.svc.ListOptions(c.Context(), authCtx.TenantID, productID)
	if err != nil {
		return err
	}

	return c.JSON(opts)
}

// UpdateOption handles PUT /products/options/:optionId.
func (h *Handler) UpdateOption(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	optionID := kernel.OptionID(c.Params("optionId"))

	var req createOptionRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("option name is required", errx.TypeValidation)
	}

	opt, err := h.svc.UpdateOption(c.Context(), authCtx.TenantID, optionID, req.Name, req.Position, req.Values)
	if err != nil {
		return err
	}

	return c.JSON(opt)
}

// DeleteOption handles DELETE /products/options/:optionId.
func (h *Handler) DeleteOption(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	optionID := kernel.OptionID(c.Params("optionId"))

	if err := h.svc.DeleteOption(c.Context(), authCtx.TenantID, optionID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ─── Variant handlers ─────────────────────────────────────────────────────────

type createVariantRequest struct {
	SKU         string            `json:"sku"`
	PriceAmount int64             `json:"price_amount"`
	Currency    string            `json:"currency"`
	Stock       int               `json:"stock"`
	Options     map[string]string `json:"options"`
}

type updateVariantRequest struct {
	SKU         string            `json:"sku"`
	PriceAmount int64             `json:"price_amount"`
	Currency    string            `json:"currency"`
	Stock       int               `json:"stock"`
	Options     map[string]string `json:"options"`
	Active      bool              `json:"active"`
}

// CreateVariant handles POST /products/:id/variants.
func (h *Handler) CreateVariant(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	productID := kernel.ProductID(c.Params("id"))

	var req createVariantRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	v, err := h.svc.CreateVariant(c.Context(), authCtx.TenantID, productID,
		req.SKU, kernel.NewMoney(req.PriceAmount, req.Currency), req.Stock, req.Options,
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(v)
}

// ListVariants handles GET /products/:id/variants.
func (h *Handler) ListVariants(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	productID := kernel.ProductID(c.Params("id"))

	variants, err := h.svc.ListVariants(c.Context(), authCtx.TenantID, productID)
	if err != nil {
		return err
	}

	return c.JSON(variants)
}

// GetVariant handles GET /products/variants/:variantId.
func (h *Handler) GetVariant(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	variantID := kernel.VariantID(c.Params("variantId"))

	v, err := h.svc.GetVariant(c.Context(), authCtx.TenantID, variantID)
	if err != nil {
		return err
	}

	return c.JSON(v)
}

// UpdateVariant handles PUT /products/variants/:variantId.
func (h *Handler) UpdateVariant(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	variantID := kernel.VariantID(c.Params("variantId"))

	var req updateVariantRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	v, err := h.svc.UpdateVariant(c.Context(), authCtx.TenantID, variantID,
		req.SKU, kernel.NewMoney(req.PriceAmount, req.Currency), req.Stock, req.Options, req.Active,
	)
	if err != nil {
		return err
	}

	return c.JSON(v)
}

// DeleteVariant handles DELETE /products/variants/:variantId.
func (h *Handler) DeleteVariant(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	variantID := kernel.VariantID(c.Params("variantId"))

	if err := h.svc.DeleteVariant(c.Context(), authCtx.TenantID, variantID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ─── Public handlers ──────────────────────────────────────────────────────────

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

// GetProductBySlugPublic handles GET /products/slug/:slug (public, no auth).
// Enables SEO-friendly product URLs.
func (h *Handler) GetProductBySlugPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	slug := c.Params("slug")
	p, err := h.svc.GetBySlug(c.Context(), tenantID, slug)
	if err != nil {
		return err
	}
	return c.JSON(p)
}

// ListOptionsPublic handles GET /products/:id/options (public).
func (h *Handler) ListOptionsPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	productID := kernel.ProductID(c.Params("id"))
	opts, err := h.svc.ListOptions(c.Context(), tenantID, productID)
	if err != nil {
		return err
	}
	return c.JSON(opts)
}

// ListVariantsPublic handles GET /products/:id/variants (public).
func (h *Handler) ListVariantsPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	productID := kernel.ProductID(c.Params("id"))
	variants, err := h.svc.ListVariants(c.Context(), tenantID, productID)
	if err != nil {
		return err
	}
	return c.JSON(variants)
}

// ─── helpers ──────────────────────────────────────────────────────────────────

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
