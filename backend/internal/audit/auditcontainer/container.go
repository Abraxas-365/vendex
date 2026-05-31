package auditcontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/audit/auditapi"
	"github.com/Abraxas-365/vendex/internal/audit/auditinfra"
	"github.com/Abraxas-365/vendex/internal/audit/auditsrv"
)

// Container wires together all audit domain dependencies.
type Container struct {
	Service *auditsrv.Service
	Handler *auditapi.Handler
}

// New creates a fully-wired audit container.
func New(db *sqlx.DB) *Container {
	repo := auditinfra.NewPostgresRepo(db)
	svc := auditsrv.New(repo)
	handler := auditapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers audit HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
