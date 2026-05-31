package analyticscontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/analytics/analyticsapi"
	"github.com/Abraxas-365/vendex/internal/analytics/analyticsinfra"
	"github.com/Abraxas-365/vendex/internal/analytics/analyticssrv"
)

// Container wires together all analytics domain dependencies.
type Container struct {
	Service *analyticssrv.Service
	Handler *analyticsapi.Handler
}

// New creates a fully-wired analytics container.
func New(db *sqlx.DB) *Container {
	repo := analyticsinfra.NewPostgresRepo(db)
	svc := analyticssrv.New(repo)
	handler := analyticsapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers analytics HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
