package reviewcontainer

import (
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/review/reviewapi"
	"github.com/Abraxas-365/vendex/internal/review/reviewinfra"
	"github.com/Abraxas-365/vendex/internal/review/reviewsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all review domain dependencies.
type Container struct {
	Service *reviewsrv.Service
	Handler *reviewapi.Handler
}

// New creates a fully-wired review container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := reviewinfra.NewPostgresRepo(db)
	svc := reviewsrv.New(repo, bus)
	handler := reviewapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers admin-authenticated review routes.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers public, read-only review routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}

// RegisterCustomerRoutes registers customer-authenticated review routes.
// The router should already have customer auth middleware applied.
func (c *Container) RegisterCustomerRoutes(router fiber.Router) {
	c.Handler.RegisterCustomerRoutes(router)
}
