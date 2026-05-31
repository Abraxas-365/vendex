// Package agenttriggercontainer wires together the agenttrigger domain.
package agenttriggercontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/agenttrigger/agenttriggerapi"
	"github.com/Abraxas-365/hada-commerce/internal/agenttrigger/agenttriggerinfra"
	"github.com/Abraxas-365/hada-commerce/internal/agenttrigger/agenttriggersrv"
)

// Container holds the wired agenttrigger domain components.
type Container struct {
	Service *agenttriggersrv.Service
	Handler *agenttriggerapi.Handler
}

// New creates a fully-wired agenttrigger container.
func New(db *sqlx.DB) *Container {
	triggerRepo := agenttriggerinfra.NewPostgresTriggerRepo(db)
	logRepo := agenttriggerinfra.NewPostgresTriggerLogRepo(db)

	svc := agenttriggersrv.NewService(triggerRepo, logRepo)
	handler := agenttriggerapi.NewHandler(svc)

	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers agent trigger HTTP routes.
func (c *Container) RegisterRoutes(r fiber.Router) {
	c.Handler.RegisterRoutes(r.Group("/agent/triggers"))
}
