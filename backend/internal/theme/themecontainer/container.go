package themecontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/theme/themeapi"
	"github.com/Abraxas-365/hada-commerce/internal/theme/themeinfra"
	"github.com/Abraxas-365/hada-commerce/internal/theme/themesrv"
)

// Container wires together all theme domain dependencies.
type Container struct {
	Service *themesrv.Service
	Handler *themeapi.Handler
}

// New creates a fully-wired theme container.
func New(db *sqlx.DB) *Container {
	repo := themeinfra.NewPostgresRepo(db)
	svc := themesrv.New(repo)
	handler := themeapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers theme HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
