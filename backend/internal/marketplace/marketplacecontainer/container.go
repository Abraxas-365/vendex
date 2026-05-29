package marketplacecontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplaceapi"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplaceinfra"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplacesrv"
)

// Container wires together all marketplace domain dependencies.
type Container struct {
	Service *marketplacesrv.VendorService
	Handler *marketplaceapi.Handler
}

// New creates a fully-wired marketplace container.
func New(db *sqlx.DB) *Container {
	vendorRepo  := marketplaceinfra.NewPostgresVendorRepo(db)
	productRepo := marketplaceinfra.NewPostgresVendorProductRepo(db)
	orderRepo   := marketplaceinfra.NewPostgresVendorOrderRepo(db)
	svc     := marketplacesrv.New(vendorRepo, productRepo, orderRepo)
	handler := marketplaceapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers marketplace HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
