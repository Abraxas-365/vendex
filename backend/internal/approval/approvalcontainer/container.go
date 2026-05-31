// Package approvalcontainer wires together the approval workflow domain.
package approvalcontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/approval/approvalapi"
	"github.com/Abraxas-365/vendex/internal/approval/approvalinfra"
	"github.com/Abraxas-365/vendex/internal/approval/approvalsrv"
)

// Container holds the wired approval domain components.
type Container struct {
	Service *approvalsrv.Service
	Handler *approvalapi.Handler
}

// New creates a fully-wired approval container.
func New(db *sqlx.DB) *Container {
	repo := approvalinfra.NewPostgresRepository(db)
	svc := approvalsrv.NewService(repo)
	handler := approvalapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers approval HTTP routes on the given router.
func (c *Container) RegisterRoutes(r fiber.Router) {
	c.Handler.RegisterRoutes(r.Group("/approvals"))
}
