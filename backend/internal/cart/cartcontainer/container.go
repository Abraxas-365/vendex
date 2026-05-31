package cartcontainer

import (
	"github.com/Abraxas-365/vendex/internal/cart/cartapi"
	"github.com/Abraxas-365/vendex/internal/cart/cartinfra"
	"github.com/Abraxas-365/vendex/internal/cart/cartsrv"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all cart domain dependencies.
type Container struct {
	Service *cartsrv.Service
	Handler *cartapi.Handler
}

// New creates a fully-wired cart container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := cartinfra.NewPostgresRepo(db)
	svc := cartsrv.New(repo, bus)
	handler := cartapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers cart admin HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
