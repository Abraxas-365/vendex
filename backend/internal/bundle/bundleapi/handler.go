package bundleapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/bundle"
	"github.com/Abraxas-365/hada-commerce/internal/bundle/bundlesrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Handler exposes HTTP endpoints for the bundle domain.
type Handler struct {
	svc *bundlesrv.Service
}

// NewHandler creates a new bundle API handler.
func NewHandler(svc *bundlesrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected admin bundle routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/bundles")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
	g.Post("/:id/items", h.AddItem)
	g.Delete("/:id/items/:itemId", h.RemoveItem)
}

// RegisterPublicRoutes registers unauthenticated read-only bundle routes.
// Tenant is identified via the X-Tenant-ID header.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/bundles")
	g.Get("/", h.ListPublic)
	g.Get("/slug/:slug", h.GetBySlugPublic)
	g.Get("/:id/price", h.GetPricePublic)
}

// ─── Admin handlers ───────────────────────────────────────────────────────────

type createBundleRequest struct {
	Name          string `json:"name"`
	Slug          string `json:"slug"`
	Description   string `json:"description"`
	DiscountType  string `json:"discount_type"`
	DiscountValue int    `json:"discount_value"`
	Active        bool   `json:"active"`
}

// Create handles POST /bundles.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createBundleRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}

	dt := bundle.DiscountType(req.DiscountType)
	if dt == "" {
		dt = bundle.DiscountPercentage
	}

	b, err := h.svc.Create(c.Context(), authCtx.TenantID, bundle.CreateBundleInput{
		Name:          req.Name,
		Slug:          req.Slug,
		Description:   req.Description,
		DiscountType:  dt,
		DiscountValue: req.DiscountValue,
		Active:        req.Active,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(b)
}

// List handles GET /bundles (admin — shows all bundles).
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, false, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetByID handles GET /bundles/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BundleID(c.Params("id"))

	b, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(b)
}

type updateBundleRequest struct {
	Name          *string `json:"name"`
	Description   *string `json:"description"`
	DiscountType  *string `json:"discount_type"`
	DiscountValue *int    `json:"discount_value"`
	Active        *bool   `json:"active"`
}

// Update handles PUT /bundles/:id.
func (h *Handler) Update(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BundleID(c.Params("id"))

	var req updateBundleRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	in := bundle.UpdateBundleInput{
		Name:          req.Name,
		Description:   req.Description,
		DiscountValue: req.DiscountValue,
		Active:        req.Active,
	}
	if req.DiscountType != nil {
		dt := bundle.DiscountType(*req.DiscountType)
		in.DiscountType = &dt
	}

	b, err := h.svc.Update(c.Context(), authCtx.TenantID, id, in)
	if err != nil {
		return err
	}

	return c.JSON(b)
}

// Delete handles DELETE /bundles/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BundleID(c.Params("id"))

	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

type addItemRequest struct {
	ProductID string  `json:"product_id"`
	VariantID *string `json:"variant_id"`
	Quantity  int     `json:"quantity"`
}

// AddItem handles POST /bundles/:id/items.
func (h *Handler) AddItem(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	bundleID := kernel.BundleID(c.Params("id"))

	var req addItemRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.ProductID == "" {
		return errx.New("product_id is required", errx.TypeValidation)
	}
	if req.Quantity < 1 {
		req.Quantity = 1
	}

	item, err := h.svc.AddItem(c.Context(), authCtx.TenantID, bundleID, bundle.AddBundleItemInput{
		ProductID: req.ProductID,
		VariantID: req.VariantID,
		Quantity:  req.Quantity,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(item)
}

// RemoveItem handles DELETE /bundles/:id/items/:itemId.
func (h *Handler) RemoveItem(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	bundleID := kernel.BundleID(c.Params("id"))
	itemID := kernel.BundleItemID(c.Params("itemId"))

	if err := h.svc.RemoveItem(c.Context(), authCtx.TenantID, bundleID, itemID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ─── Public handlers ──────────────────────────────────────────────────────────

// ListPublic handles GET /bundles (public — active bundles only).
func (h *Handler) ListPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	pg := paginationFromQuery(c)
	result, err := h.svc.List(c.Context(), tenantID, true, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetBySlugPublic handles GET /bundles/slug/:slug with calculated price.
func (h *Handler) GetBySlugPublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	slug := c.Params("slug")
	b, err := h.svc.GetBySlug(c.Context(), tenantID, slug)
	if err != nil {
		return err
	}

	return c.JSON(b)
}

// GetPricePublic handles GET /bundles/:id/price.
// Accepts optional query params: product prices are resolved server-side if
// a price provider is available; here we return a zero-base-total result
// indicating the caller must supply product prices via the admin route or
// integrate the product service.
func (h *Handler) GetPricePublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	id := kernel.BundleID(c.Params("id"))
	// Calculate price with an empty price map — items without a matching price
	// will contribute 0 to the base total, which is expected when no external
	// product price feed is wired in this endpoint.
	result, err := h.svc.CalculatePrice(c.Context(), tenantID, id, map[kernel.ProductID]kernel.Money{})
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

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
