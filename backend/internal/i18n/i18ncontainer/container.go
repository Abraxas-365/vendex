package i18ncontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/i18n/i18napi"
	"github.com/Abraxas-365/hada-commerce/internal/i18n/i18ninfra"
	"github.com/Abraxas-365/hada-commerce/internal/i18n/i18nsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all i18n domain dependencies.
type Container struct {
	Service *i18nsrv.Service
	Handler *i18napi.Handler
}

// New creates a fully-wired i18n container.
func New(db *sqlx.DB) *Container {
	repo := i18ninfra.NewPostgresRepo(db)
	svc := i18nsrv.New(repo)
	handler := i18napi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected admin routes for translation management.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers public read-only translation routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
