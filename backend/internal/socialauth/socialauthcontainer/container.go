package socialauthcontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/socialauth/socialauthapi"
	"github.com/Abraxas-365/hada-commerce/internal/socialauth/socialauthinfra"
	"github.com/Abraxas-365/hada-commerce/internal/socialauth/socialauthsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all social auth domain dependencies.
type Container struct {
	Service *socialauthsrv.Service
	Handler *socialauthapi.Handler
}

// New creates a fully-wired social auth container.
func New(db *sqlx.DB) *Container {
	repo := socialauthinfra.NewPostgresRepo(db)
	svc := socialauthsrv.New(repo)
	handler := socialauthapi.NewHandler(svc)

	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterPublicRoutes registers public storefront social auth routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}

// RegisterCustomerRoutes registers customer-authenticated social account routes.
// The router should already have CustomerMiddleware applied.
func (c *Container) RegisterCustomerRoutes(router fiber.Router) {
	c.Handler.RegisterCustomerRoutes(router)
}

// RegisterRoutes registers admin-protected social auth routes.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
