package wishlistapi

import (
	customerauth "github.com/Abraxas-365/hada-commerce/internal/customer/auth"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/wishlist/wishlistsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the wishlist domain.
type Handler struct {
	svc *wishlistsrv.Service
}

// NewHandler creates a new wishlist API handler.
func NewHandler(svc *wishlistsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterCustomerRoutes registers customer-authenticated wishlist routes.
// The router must already have CustomerMiddleware applied.
func (h *Handler) RegisterCustomerRoutes(router fiber.Router) {
	g := router.Group("/storefront/wishlist")
	g.Get("/", h.GetWishlist)
	g.Post("/items", h.AddItem)
	g.Delete("/items/:itemId", h.RemoveItem)
	g.Delete("/", h.ClearWishlist)
}

// --- Request types ---

type addItemRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id"`
}

// --- Handlers ---

// GetWishlist handles GET /storefront/wishlist.
func (h *Handler) GetWishlist(c *fiber.Ctx) error {
	auth, ok := customerauth.GetCustomerAuthContext(c)
	if !ok {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	tenantID := kernel.TenantID(auth.TenantID)
	customerID := kernel.CustomerID(auth.CustomerID)

	result, err := h.svc.GetOrCreateWishlist(c.Context(), tenantID, customerID)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// AddItem handles POST /storefront/wishlist/items.
func (h *Handler) AddItem(c *fiber.Ctx) error {
	auth, ok := customerauth.GetCustomerAuthContext(c)
	if !ok {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	tenantID := kernel.TenantID(auth.TenantID)
	customerID := kernel.CustomerID(auth.CustomerID)

	var req addItemRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.ProductID == "" {
		return errx.New("product_id is required", errx.TypeValidation)
	}

	result, err := h.svc.AddItem(c.Context(), tenantID, customerID, kernel.ProductID(req.ProductID), req.VariantID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// RemoveItem handles DELETE /storefront/wishlist/items/:itemId.
func (h *Handler) RemoveItem(c *fiber.Ctx) error {
	auth, ok := customerauth.GetCustomerAuthContext(c)
	if !ok {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	tenantID := kernel.TenantID(auth.TenantID)
	customerID := kernel.CustomerID(auth.CustomerID)
	itemID := kernel.WishlistItemID(c.Params("itemId"))

	result, err := h.svc.RemoveItem(c.Context(), tenantID, customerID, itemID)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// ClearWishlist handles DELETE /storefront/wishlist.
func (h *Handler) ClearWishlist(c *fiber.Ctx) error {
	auth, ok := customerauth.GetCustomerAuthContext(c)
	if !ok {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	tenantID := kernel.TenantID(auth.TenantID)
	customerID := kernel.CustomerID(auth.CustomerID)

	if err := h.svc.ClearWishlist(c.Context(), tenantID, customerID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
