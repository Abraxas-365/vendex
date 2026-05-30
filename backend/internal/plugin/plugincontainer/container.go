package plugincontainer

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/plugin"
	"github.com/Abraxas-365/hada-commerce/internal/plugin/pluginapi"
	"github.com/Abraxas-365/hada-commerce/internal/plugin/plugininfra"
	"github.com/Abraxas-365/hada-commerce/internal/plugin/pluginsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all plugin domain dependencies.
type Container struct {
	Service *pluginsrv.Service
	Handler *pluginapi.Handler
}

// New creates a fully-wired plugin container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	pluginRepo := plugininfra.NewPluginRepo(db)
	versionRepo := plugininfra.NewVersionRepo(db)
	installRepo := plugininfra.NewInstallationRepo(db)
	svc := pluginsrv.New(pluginRepo, versionRepo, installRepo, bus)
	handler := pluginapi.NewHandler(svc)

	// Wire webhook dispatcher — delivers domain events to installed plugin endpoints
	dispatcher := plugin.NewWebhookDispatcher(installRepo, versionRepo, &http.Client{})
	bus.SubscribeAll(dispatcher.HandleEvent)

	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers plugin HTTP routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
