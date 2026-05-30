package inventorycontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/inventory/inventoryapi"
	"github.com/Abraxas-365/hada-commerce/internal/inventory/inventoryinfra"
	"github.com/Abraxas-365/hada-commerce/internal/inventory/inventorysrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all inventory domain dependencies.
type Container struct {
	Service *inventorysrv.Service
	Handler *inventoryapi.Handler
}

// New creates a fully-wired inventory container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := inventoryinfra.NewPostgresRepo(db)
	svc := inventorysrv.NewService(repo, bus)
	handler := inventoryapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers inventory HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
