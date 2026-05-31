package customergroupcontainer

import (
	"github.com/Abraxas-365/vendex/internal/customergroup/customergroupapi"
	"github.com/Abraxas-365/vendex/internal/customergroup/customergroupinfra"
	"github.com/Abraxas-365/vendex/internal/customergroup/customergroupsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all customer group domain dependencies.
type Container struct {
	Service *customergroupsrv.Service
	Handler *customergroupapi.Handler
}

// New creates a fully-wired customer group container.
func New(db *sqlx.DB) *Container {
	repo := customergroupinfra.NewPostgresRepo(db)
	svc := customergroupsrv.New(repo)
	handler := customergroupapi.NewHandler(svc)

	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers admin-facing customer group routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
