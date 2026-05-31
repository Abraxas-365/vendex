package shippingcontainer

import (
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/shipping/shippingapi"
	"github.com/Abraxas-365/vendex/internal/shipping/shippinginfra"
	"github.com/Abraxas-365/vendex/internal/shipping/shippingsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all shipping domain dependencies.
type Container struct {
	Service *shippingsrv.Service
	Handler *shippingapi.Handler
}

// New creates a fully-wired shipping container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	zoneRepo := shippinginfra.NewZonePostgresRepo(db)
	rateRepo := shippinginfra.NewRatePostgresRepo(db)
	svc := shippingsrv.New(zoneRepo, rateRepo, bus)
	handler := shippingapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected shipping routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
