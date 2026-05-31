package taxcontainer

import (
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/tax/taxapi"
	"github.com/Abraxas-365/vendex/internal/tax/taxinfra"
	"github.com/Abraxas-365/vendex/internal/tax/taxsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all tax domain dependencies.
type Container struct {
	Service *taxsrv.Service
	Handler *taxapi.Handler
}

// New creates a fully-wired tax container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := taxinfra.NewPostgresRepo(db)
	svc := taxsrv.New(repo, bus)
	handler := taxapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected tax admin HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
