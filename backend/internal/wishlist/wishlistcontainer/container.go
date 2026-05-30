package wishlistcontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/wishlist/wishlistapi"
	"github.com/Abraxas-365/hada-commerce/internal/wishlist/wishlistinfra"
	"github.com/Abraxas-365/hada-commerce/internal/wishlist/wishlistsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all wishlist domain dependencies.
type Container struct {
	Service *wishlistsrv.Service
	Handler *wishlistapi.Handler
}

// New creates a fully-wired wishlist container.
func New(db *sqlx.DB) *Container {
	repo := wishlistinfra.NewPostgresRepo(db)
	svc := wishlistsrv.New(repo)
	handler := wishlistapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterCustomerRoutes registers customer-authenticated wishlist routes.
// The router should already have customer auth middleware applied.
func (c *Container) RegisterCustomerRoutes(router fiber.Router) {
	c.Handler.RegisterCustomerRoutes(router)
}
