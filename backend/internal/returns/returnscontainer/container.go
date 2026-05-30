package returnscontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/returns/returnsapi"
	"github.com/Abraxas-365/hada-commerce/internal/returns/returnsinfra"
	"github.com/Abraxas-365/hada-commerce/internal/returns/returnssrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all returns domain dependencies.
type Container struct {
	Service *returnssrv.Service
	Handler *returnsapi.Handler
}

// New creates a fully-wired returns container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := returnsinfra.NewPostgresRepo(db)
	svc := returnssrv.New(repo, bus)
	handler := returnsapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected (admin) return routes.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers public (customer-facing) return routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
