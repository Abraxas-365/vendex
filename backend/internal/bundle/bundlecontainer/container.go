package bundlecontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/bundle/bundleapi"
	"github.com/Abraxas-365/hada-commerce/internal/bundle/bundleinfra"
	"github.com/Abraxas-365/hada-commerce/internal/bundle/bundlesrv"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
)

// Container wires together all bundle domain dependencies.
type Container struct {
	Service *bundlesrv.Service
	Handler *bundleapi.Handler
}

// New creates a fully-wired bundle container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := bundleinfra.NewPostgresRepository(db)
	svc := bundlesrv.New(repo, bus)
	handler := bundleapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected bundle HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers public (unauthenticated) bundle HTTP routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
