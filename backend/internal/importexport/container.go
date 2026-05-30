package importexport

import (
	"github.com/gofiber/fiber/v2"
)

// Container wires the import/export domain.
type Container struct {
	Handler *Handler
}

// New creates a new import/export container from the provided service dependencies.
func New(products ProductLister, orders OrderLister, customers CustomerLister, creator ProductCreator) *Container {
	exporter := NewExportService(products, orders, customers)
	importer := NewImportService(creator)
	handler := NewHandler(exporter, importer)
	return &Container{Handler: handler}
}

// RegisterRoutes registers all import/export routes on the given (protected) router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
