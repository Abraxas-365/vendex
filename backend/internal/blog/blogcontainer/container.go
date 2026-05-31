package blogcontainer

import (
	"github.com/Abraxas-365/vendex/internal/blog/blogapi"
	"github.com/Abraxas-365/vendex/internal/blog/bloginfra"
	"github.com/Abraxas-365/vendex/internal/blog/blogsrv"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all blog domain dependencies.
type Container struct {
	Service *blogsrv.Service
	Handler *blogapi.Handler
}

// New creates a fully-wired blog container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := bloginfra.NewPostgresRepo(db)
	svc := blogsrv.New(repo, bus)
	handler := blogapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected (admin) blog HTTP routes.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers public blog HTTP routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
