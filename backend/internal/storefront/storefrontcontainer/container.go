package storefrontcontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontapi"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontinfra"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontsrv"
)

// Container wires together all storefront domain dependencies.
type Container struct {
	Service *storefrontsrv.Service
	Handler *storefrontapi.Handler
}

// New creates a fully-wired storefront container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	pagesRepo := storefrontinfra.NewPagePostgresRepo(db)
	versionsRepo := storefrontinfra.NewPageVersionPostgresRepo(db)
	blockTypesRepo := storefrontinfra.NewBlockTypePostgresRepo(db)
	svc := storefrontsrv.New(pagesRepo, versionsRepo, blockTypesRepo, bus)
	handler := storefrontapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers storefront HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
