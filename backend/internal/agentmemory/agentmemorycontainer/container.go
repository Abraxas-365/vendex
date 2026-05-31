// Package agentmemorycontainer wires together the agent memory domain.
package agentmemorycontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/agentmemory/agentmemoryapi"
	"github.com/Abraxas-365/vendex/internal/agentmemory/agentmemoryinfra"
	"github.com/Abraxas-365/vendex/internal/agentmemory/agentmemorysrv"
)

// Container holds the wired agent memory domain components.
type Container struct {
	Service *agentmemorysrv.Service
	Handler *agentmemoryapi.Handler
}

// New creates a fully-wired agentmemory container.
func New(db *sqlx.DB) *Container {
	repo := agentmemoryinfra.NewPostgresRepository(db)
	svc := agentmemorysrv.NewService(repo)
	handler := agentmemoryapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers agent memory HTTP routes.
func (c *Container) RegisterRoutes(r fiber.Router) {
	c.Handler.RegisterRoutes(r.Group("/agent/memories"))
}
