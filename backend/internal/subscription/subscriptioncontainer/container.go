package subscriptioncontainer

import (
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/subscription/subscriptionapi"
	"github.com/Abraxas-365/vendex/internal/subscription/subscriptioninfra"
	"github.com/Abraxas-365/vendex/internal/subscription/subscriptionsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all subscription domain dependencies.
type Container struct {
	Service *subscriptionsrv.Service
	Handler *subscriptionapi.Handler
}

// New creates a fully-wired subscription container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := subscriptioninfra.NewPostgresRepo(db)
	svc := subscriptionsrv.New(repo, bus)
	handler := subscriptionapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers admin subscription routes on the given (protected) router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterCustomerRoutes registers storefront subscription routes.
func (c *Container) RegisterCustomerRoutes(router fiber.Router) {
	c.Handler.RegisterCustomerRoutes(router)
}
