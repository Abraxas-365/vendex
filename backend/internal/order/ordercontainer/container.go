package ordercontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/order/orderapi"
	"github.com/Abraxas-365/hada-commerce/internal/order/orderinfra"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordersrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all order domain dependencies.
type Container struct {
	Service *ordersrv.Service
	Handler *orderapi.Handler
}

// New creates a fully-wired order container.
func New(db *sqlx.DB) *Container {
	repo := orderinfra.NewPostgresRepo(db)
	svc := ordersrv.New(repo)
	handler := orderapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers order HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
