package bulkopscontainer

import (
	"github.com/Abraxas-365/vendex/internal/bulkops/bulkopsapi"
	"github.com/Abraxas-365/vendex/internal/bulkops/bulkopsinfra"
	"github.com/Abraxas-365/vendex/internal/bulkops/bulkopssrv"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all bulk operations domain dependencies.
type Container struct {
	Service *bulkopssrv.Service
	Handler *bulkopsapi.Handler
}

// New creates a fully-wired bulk operations container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := bulkopsinfra.NewPostgresRepo(db)
	svc := bulkopssrv.New(repo, bus)
	handler := bulkopsapi.NewHandler(svc)

	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected bulk operation routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
