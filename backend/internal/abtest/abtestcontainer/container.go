package abtestcontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/abtest/abtestapi"
	"github.com/Abraxas-365/hada-commerce/internal/abtest/abtestinfra"
	"github.com/Abraxas-365/hada-commerce/internal/abtest/abtestsrv"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all A/B testing domain dependencies.
type Container struct {
	Service *abtestsrv.Service
	Handler *abtestapi.Handler
}

// New creates a fully-wired A/B testing container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := abtestinfra.NewPostgresRepo(db)
	svc := abtestsrv.New(repo, bus)
	handler := abtestapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected A/B test admin routes.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers unauthenticated A/B test routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
