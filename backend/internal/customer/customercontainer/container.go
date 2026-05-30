package customercontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/customer/customerapi"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customerinfra"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customersrv"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all customer domain dependencies.
type Container struct {
	Service *customersrv.Service
	Handler *customerapi.Handler
}

// New creates a fully-wired customer container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := customerinfra.NewPostgresRepo(db)
	svc := customersrv.New(repo, bus)
	handler := customerapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers customer HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
