package dashboardcontainer

import (
	"github.com/Abraxas-365/vendex/internal/dashboard/dashboardapi"
	"github.com/Abraxas-365/vendex/internal/dashboard/dashboardinfra"
	"github.com/Abraxas-365/vendex/internal/dashboard/dashboardsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all dashboard domain dependencies.
type Container struct {
	Service *dashboardsrv.Service
	Handler *dashboardapi.Handler
}

// New creates a fully-wired dashboard container.
func New(db *sqlx.DB) *Container {
	repo := dashboardinfra.NewPostgresRepo(db)
	svc := dashboardsrv.New(repo)
	handler := dashboardapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers dashboard HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
