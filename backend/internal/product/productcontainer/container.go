package productcontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/product/productapi"
	"github.com/Abraxas-365/hada-commerce/internal/product/productinfra"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all product domain dependencies.
type Container struct {
	Service *productsrv.Service
	Handler *productapi.Handler
}

// New creates a fully-wired product container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := productinfra.NewPostgresRepo(db)
	variantRepo := productinfra.NewVariantPostgresRepo(db)
	svc := productsrv.New(repo, variantRepo, bus)
	handler := productapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers product HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
