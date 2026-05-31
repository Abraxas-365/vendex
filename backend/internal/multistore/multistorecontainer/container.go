package multistorecontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/multistore/multistoreapi"
	"github.com/Abraxas-365/hada-commerce/internal/multistore/multistoreinfra"
	"github.com/Abraxas-365/hada-commerce/internal/multistore/multistoresrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all multistore domain dependencies.
type Container struct {
	Service *multistoresrv.Service
	Handler *multistoreapi.Handler
}

// New creates a fully-wired multistore container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := multistoreinfra.NewPostgresRepo(db)
	svc := multistoresrv.New(repo, bus)
	handler := multistoreapi.NewHandler(svc)

	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected admin routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers unauthenticated routes on the given router.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
