// Package agentsessioncontainer wires together the agent session domain.
package agentsessioncontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/agentsession/agentsessionapi"
	"github.com/Abraxas-365/hada-commerce/internal/agentsession/agentsessioninfra"
	"github.com/Abraxas-365/hada-commerce/internal/agentsession/agentsessionsrv"
	"github.com/Abraxas-365/hada-commerce/internal/containerx"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplacesrv"
)

// Container holds the wired agentsession domain components.
type Container struct {
	Service *agentsessionsrv.Service
	Handler *agentsessionapi.Handler
}

// Deps holds external dependencies needed by the agent session domain.
type Deps struct {
	DB        *sqlx.DB
	Manager   containerx.Manager
	PresetSvc *marketplacesrv.PresetService
}

// New creates a fully-wired agentsession container.
func New(deps Deps) *Container {
	sessionRepo := agentsessioninfra.NewPostgresSessionRepo(deps.DB)
	chatRepo := agentsessioninfra.NewPostgresChatRepo(deps.DB)

	svc := agentsessionsrv.NewService(sessionRepo, chatRepo, deps.Manager, deps.PresetSvc)
	handler := agentsessionapi.NewHandler(svc)

	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers agent session HTTP routes.
func (c *Container) RegisterRoutes(r fiber.Router) {
	c.Handler.RegisterRoutes(r.Group("/agent/sessions"))
}
