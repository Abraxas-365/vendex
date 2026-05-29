package marketplaceapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplacesrv"
)

// Handler exposes HTTP endpoints for the marketplace domain.
type Handler struct {
	svc *marketplacesrv.VendorService
}

// NewHandler creates a new marketplace API handler.
func NewHandler(svc *marketplacesrv.VendorService) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all marketplace routes on the given Fiber router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	vendors := router.Group("/marketplace/vendors")
	vendors.Post("/", h.CreateVendor)
	vendors.Get("/", h.ListVendors)
	vendors.Get("/:id", h.GetVendor)
	vendors.Put("/:id", h.UpdateVendor)
	vendors.Delete("/:id", h.DeleteVendor)

	vendors.Get("/:id/products", h.ListVendorProducts)
	vendors.Post("/:id/products", h.AddVendorProduct)
	vendors.Delete("/:id/products/:pid", h.RemoveVendorProduct)

	vendors.Get("/:id/orders", h.ListVendorOrders)
}

func authFromCtx(c *fiber.Ctx) (*kernel.AuthContext, error) {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return nil, errx.New("unauthorized", errx.TypeAuthorization)
	}
	return authCtx, nil
}

func pagination(c *fiber.Ctx) kernel.PaginationOptions {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	return kernel.PaginationOptions{Page: page, PageSize: pageSize}
}

// CreateVendor handles POST /marketplace/vendors.
func (h *Handler) CreateVendor(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	var req marketplace.CreateVendorRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	vendor, err := h.svc.CreateVendor(c.Context(), auth.TenantID, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(vendor)
}

// GetVendor handles GET /marketplace/vendors/:id.
func (h *Handler) GetVendor(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	id := kernel.VendorID(c.Params("id"))
	vendor, err := h.svc.GetVendor(c.Context(), auth.TenantID, id)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(vendor)
}

// UpdateVendor handles PUT /marketplace/vendors/:id.
func (h *Handler) UpdateVendor(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	id := kernel.VendorID(c.Params("id"))
	var req marketplace.UpdateVendorRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	vendor, err := h.svc.UpdateVendor(c.Context(), auth.TenantID, id, req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(vendor)
}

// DeleteVendor handles DELETE /marketplace/vendors/:id.
func (h *Handler) DeleteVendor(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	id := kernel.VendorID(c.Params("id"))
	if err := h.svc.DeleteVendor(c.Context(), auth.TenantID, id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ListVendors handles GET /marketplace/vendors.
func (h *Handler) ListVendors(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	result, err := h.svc.ListVendors(c.Context(), auth.TenantID, pagination(c))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// AddVendorProduct handles POST /marketplace/vendors/:id/products.
func (h *Handler) AddVendorProduct(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	vendorID := kernel.VendorID(c.Params("id"))
	var vp marketplace.VendorProduct
	if err := c.BodyParser(&vp); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	vp.VendorID = vendorID

	result, err := h.svc.AddVendorProduct(c.Context(), auth.TenantID, vp)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

// ListVendorProducts handles GET /marketplace/vendors/:id/products.
func (h *Handler) ListVendorProducts(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	vendorID := kernel.VendorID(c.Params("id"))
	result, err := h.svc.ListVendorProducts(c.Context(), auth.TenantID, vendorID, pagination(c))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(result)
}

// RemoveVendorProduct handles DELETE /marketplace/vendors/:id/products/:pid.
func (h *Handler) RemoveVendorProduct(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	pid := c.Params("pid")
	if err := h.svc.RemoveVendorProduct(c.Context(), auth.TenantID, pid); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// ListVendorOrders handles GET /marketplace/vendors/:id/orders.
func (h *Handler) ListVendorOrders(c *fiber.Ctx) error {
	auth, err := authFromCtx(c)
	if err != nil {
		return err
	}

	vendorID := kernel.VendorID(c.Params("id"))
	result, err := h.svc.ListVendorOrders(c.Context(), auth.TenantID, vendorID, pagination(c))
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusOK).JSON(result)
}
