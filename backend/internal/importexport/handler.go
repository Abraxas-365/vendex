package importexport

import (
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for CSV import/export.
type Handler struct {
	exporter *ExportService
	importer *ImportService
}

// NewHandler creates a new import/export handler.
func NewHandler(exporter *ExportService, importer *ImportService) *Handler {
	return &Handler{
		exporter: exporter,
		importer: importer,
	}
}

// RegisterRoutes registers all import/export routes on the given (protected) router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/import-export")
	g.Get("/products/export", h.ExportProducts)
	g.Get("/orders/export", h.ExportOrders)
	g.Get("/customers/export", h.ExportCustomers)
	g.Post("/products/import", h.ImportProducts)
}

// ExportProducts handles GET /import-export/products/export
func (h *Handler) ExportProducts(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=products.csv")

	if err := h.exporter.ExportProducts(c.Context(), authCtx.TenantID, c.Response().BodyWriter()); err != nil {
		return err
	}
	return nil
}

// ExportOrders handles GET /import-export/orders/export
func (h *Handler) ExportOrders(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=orders.csv")

	if err := h.exporter.ExportOrders(c.Context(), authCtx.TenantID, c.Response().BodyWriter()); err != nil {
		return err
	}
	return nil
}

// ExportCustomers handles GET /import-export/customers/export
func (h *Handler) ExportCustomers(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=customers.csv")

	if err := h.exporter.ExportCustomers(c.Context(), authCtx.TenantID, c.Response().BodyWriter()); err != nil {
		return err
	}
	return nil
}

// ImportProducts handles POST /import-export/products/import
// Expects a multipart/form-data upload with field name "file".
func (h *Handler) ImportProducts(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return errx.New("missing or invalid file field", errx.TypeValidation)
	}

	f, err := fileHeader.Open()
	if err != nil {
		return errx.Wrap(err, "opening uploaded file", errx.TypeInternal)
	}
	defer f.Close()

	result, err := h.importer.ImportProducts(c.Context(), authCtx.TenantID, f)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}
