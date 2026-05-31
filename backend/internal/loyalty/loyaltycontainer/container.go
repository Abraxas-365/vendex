package loyaltycontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/loyalty/loyaltyapi"
	"github.com/Abraxas-365/vendex/internal/loyalty/loyaltyinfra"
	"github.com/Abraxas-365/vendex/internal/loyalty/loyaltysrv"
)

// Container wires together the loyalty domain's repository, service, and handler.
type Container struct {
	Service *loyaltysrv.Service
	Handler *loyaltyapi.Handler
}

// New builds the full loyalty dependency graph.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := loyaltyinfra.NewPostgresRepository(db)
	svc := loyaltysrv.New(repo, bus)
	handler := loyaltyapi.New(svc)
	return &Container{Service: svc, Handler: handler}
}

// RegisterRoutes wires all protected loyalty routes.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes wires all public loyalty routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
