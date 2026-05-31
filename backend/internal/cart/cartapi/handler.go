package cartapi

import (
	"github.com/Abraxas-365/vendex/internal/cart/cartsrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the cart domain.
type Handler struct {
	svc *cartsrv.Service
}

// NewHandler creates a new cart API handler.
func NewHandler(svc *cartsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers (empty) admin cart routes.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	// No admin-only cart routes at this time.
}

// RegisterPublicRoutes registers unauthenticated cart routes (tenant via X-Tenant-ID header).
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/cart")
	g.Post("/", h.GetOrCreate)
	g.Get("/:id", h.GetByID)
	g.Post("/:id/items", h.AddItem)
	g.Put("/:id/items/:itemId", h.UpdateItem)
	g.Delete("/:id/items/:itemId", h.RemoveItem)
	g.Delete("/:id", h.DeleteCart)
}

// --- Request/response types ---

type getOrCreateRequest struct {
	SessionID  string `json:"session_id"`
	CustomerID string `json:"customer_id"`
	Currency   string `json:"currency"`
}

type addItemRequest struct {
	ProductID           string `json:"product_id"`
	VariantID           string `json:"variant_id"`
	Quantity            int    `json:"quantity"`
	UnitPriceAmount     int64  `json:"unit_price_amount"`
	UnitPriceCurrency   string `json:"unit_price_currency"`
}

type updateItemRequest struct {
	Quantity int `json:"quantity"`
}

// --- Handlers ---

// GetOrCreate handles POST /cart — creates or retrieves a cart.
func (h *Handler) GetOrCreate(c *fiber.Ctx) error {
	tenantID, err := tenantFromHeader(c)
	if err != nil {
		return err
	}

	var req getOrCreateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.SessionID == "" && req.CustomerID == "" {
		return errx.New("session_id or customer_id is required", errx.TypeValidation)
	}

	result, err := h.svc.GetOrCreateCart(
		c.Context(),
		tenantID,
		req.SessionID,
		kernel.CustomerID(req.CustomerID),
		req.Currency,
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// GetByID handles GET /cart/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	tenantID, err := tenantFromHeader(c)
	if err != nil {
		return err
	}

	cartID := kernel.CartID(c.Params("id"))
	result, err := h.svc.GetCart(c.Context(), tenantID, cartID)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// AddItem handles POST /cart/:id/items.
func (h *Handler) AddItem(c *fiber.Ctx) error {
	tenantID, err := tenantFromHeader(c)
	if err != nil {
		return err
	}

	cartID := kernel.CartID(c.Params("id"))

	var req addItemRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.ProductID == "" {
		return errx.New("product_id is required", errx.TypeValidation)
	}
	if req.Quantity <= 0 {
		return errx.New("quantity must be greater than 0", errx.TypeValidation)
	}
	currency := req.UnitPriceCurrency
	if currency == "" {
		currency = "USD"
	}

	result, err := h.svc.AddItem(
		c.Context(),
		tenantID,
		cartID,
		kernel.ProductID(req.ProductID),
		req.VariantID,
		req.Quantity,
		kernel.NewMoney(req.UnitPriceAmount, currency),
	)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// UpdateItem handles PUT /cart/:id/items/:itemId.
func (h *Handler) UpdateItem(c *fiber.Ctx) error {
	tenantID, err := tenantFromHeader(c)
	if err != nil {
		return err
	}

	cartID := kernel.CartID(c.Params("id"))
	itemID := kernel.CartItemID(c.Params("itemId"))

	var req updateItemRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.Quantity <= 0 {
		return errx.New("quantity must be greater than 0", errx.TypeValidation)
	}

	result, err := h.svc.UpdateItemQuantity(c.Context(), tenantID, cartID, itemID, req.Quantity)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// RemoveItem handles DELETE /cart/:id/items/:itemId.
func (h *Handler) RemoveItem(c *fiber.Ctx) error {
	tenantID, err := tenantFromHeader(c)
	if err != nil {
		return err
	}

	cartID := kernel.CartID(c.Params("id"))
	itemID := kernel.CartItemID(c.Params("itemId"))

	result, err := h.svc.RemoveItem(c.Context(), tenantID, cartID, itemID)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// DeleteCart handles DELETE /cart/:id — clears and deletes the cart.
func (h *Handler) DeleteCart(c *fiber.Ctx) error {
	tenantID, err := tenantFromHeader(c)
	if err != nil {
		return err
	}

	cartID := kernel.CartID(c.Params("id"))

	if err := h.svc.DeleteCart(c.Context(), tenantID, cartID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// --- helpers ---

func tenantFromHeader(c *fiber.Ctx) (kernel.TenantID, error) {
	tid := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tid == "" {
		return "", errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}
	return tid, nil
}
