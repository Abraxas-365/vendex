package settingscontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/settings/settingsapi"
	"github.com/Abraxas-365/hada-commerce/internal/settings/settingsinfra"
	"github.com/Abraxas-365/hada-commerce/internal/settings/settingssrv"
)

// Container wires together all settings domain dependencies.
type Container struct {
	Service *settingssrv.Service
	Handler *settingsapi.Handler
}

// New creates a fully-wired settings container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := settingsinfra.NewPostgresRepo(db)
	svc := settingssrv.New(repo, bus)
	handler := settingsapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers settings HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
