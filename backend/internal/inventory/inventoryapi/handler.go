package inventoryapi

import (
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/inventory"
	"github.com/Abraxas-365/hada-commerce/internal/inventory/inventorysrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the inventory domain.
type Handler struct {
	svc *inventorysrv.Service
}

// NewHandler creates a new inventory API handler.
func NewHandler(svc *inventorysrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all inventory routes on the given router.
// All routes require authentication (caller is responsible for auth middleware).
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/inventory")

	// Warehouses
	g.Post("/warehouses", h.CreateWarehouse)
	g.Get("/warehouses", h.ListWarehouses)
	g.Get("/warehouses/:id", h.GetWarehouse)
	g.Put("/warehouses/:id", h.UpdateWarehouse)
	g.Delete("/warehouses/:id", h.DeleteWarehouse)

	// Stock levels
	g.Get("/stock/:productId", h.GetStockLevels)
	g.Post("/stock/adjust", h.AdjustStock)
	g.Get("/stock/low", h.GetLowStockItems)

	// Stock movements
	g.Get("/movements/:productId", h.ListMovements)
}

// ─── Request/response bodies ──────────────────────────────────────────────────

type createWarehouseReq struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	IsDefault bool   `json:"is_default"`
}

type updateWarehouseReq struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	IsDefault bool   `json:"is_default"`
	Active    bool   `json:"active"`
}

type adjustStockReq struct {
	ProductID   string  `json:"product_id"`
	VariantID   *string `json:"variant_id,omitempty"`
	WarehouseID string  `json:"warehouse_id"`
	Quantity    int     `json:"quantity"`
	Type        string  `json:"type"`
	Reference   string  `json:"reference"`
	Note        string  `json:"note"`
}

// ─── Warehouse handlers ───────────────────────────────────────────────────────

// CreateWarehouse handles POST /inventory/warehouses.
func (h *Handler) CreateWarehouse(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createWarehouseReq
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("warehouse name is required", errx.TypeValidation)
	}

	w, err := h.svc.CreateWarehouse(c.Context(), authCtx.TenantID, inventory.CreateWarehouseInput{
		Name:      req.Name,
		Address:   req.Address,
		IsDefault: req.IsDefault,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(w)
}

// ListWarehouses handles GET /inventory/warehouses.
func (h *Handler) ListWarehouses(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromQuery(c)

	result, err := h.svc.ListWarehouses(c.Context(), authCtx.TenantID, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetWarehouse handles GET /inventory/warehouses/:id.
func (h *Handler) GetWarehouse(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.WarehouseID(c.Params("id"))

	w, err := h.svc.GetWarehouse(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(w)
}

// UpdateWarehouse handles PUT /inventory/warehouses/:id.
func (h *Handler) UpdateWarehouse(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.WarehouseID(c.Params("id"))

	var req updateWarehouseReq
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("warehouse name is required", errx.TypeValidation)
	}

	w, err := h.svc.UpdateWarehouse(c.Context(), authCtx.TenantID, id, inventory.UpdateWarehouseInput{
		Name:      req.Name,
		Address:   req.Address,
		IsDefault: req.IsDefault,
		Active:    req.Active,
	})
	if err != nil {
		return err
	}

	return c.JSON(w)
}

// DeleteWarehouse handles DELETE /inventory/warehouses/:id.
func (h *Handler) DeleteWarehouse(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.WarehouseID(c.Params("id"))

	if err := h.svc.DeleteWarehouse(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ─── Stock handlers ───────────────────────────────────────────────────────────

// GetStockLevels handles GET /inventory/stock/:productId.
// Returns all stock levels for the product across all warehouses.
func (h *Handler) GetStockLevels(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	productID := kernel.ProductID(c.Params("productId"))

	levels, err := h.svc.ListStockLevels(c.Context(), authCtx.TenantID, productID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": levels})
}

// AdjustStock handles POST /inventory/stock/adjust.
func (h *Handler) AdjustStock(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req adjustStockReq
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.ProductID == "" {
		return errx.New("product_id is required", errx.TypeValidation)
	}
	if req.WarehouseID == "" {
		return errx.New("warehouse_id is required", errx.TypeValidation)
	}
	if req.Type == "" {
		return errx.New("type is required", errx.TypeValidation)
	}
	if req.Quantity == 0 {
		return errx.New("quantity must be non-zero", errx.TypeValidation)
	}

	in := inventory.AdjustStockInput{
		ProductID:   kernel.ProductID(req.ProductID),
		WarehouseID: kernel.WarehouseID(req.WarehouseID),
		Quantity:    req.Quantity,
		Type:        inventory.MovementType(req.Type),
		Reference:   req.Reference,
		Note:        req.Note,
	}
	if req.VariantID != nil {
		vid := kernel.VariantID(*req.VariantID)
		in.VariantID = &vid
	}

	// Optionally carry the authenticated user as "created_by".
	if authCtx.UserID != nil {
		in.CreatedBy = string(*authCtx.UserID)
	}

	sl, err := h.svc.AdjustStock(c.Context(), authCtx.TenantID, in)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(sl)
}

// GetLowStockItems handles GET /inventory/stock/low.
func (h *Handler) GetLowStockItems(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	items, err := h.svc.GetLowStockItems(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": items})
}

// ─── Movement handlers ────────────────────────────────────────────────────────

// ListMovements handles GET /inventory/movements/:productId.
func (h *Handler) ListMovements(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	productID := kernel.ProductID(c.Params("productId"))
	pg := paginationFromQuery(c)

	result, err := h.svc.ListMovements(c.Context(), authCtx.TenantID, productID, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
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
