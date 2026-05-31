package cartrecoverycontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/cartrecovery/cartrecoveryapi"
	"github.com/Abraxas-365/vendex/internal/cartrecovery/cartrecoveryinfra"
	"github.com/Abraxas-365/vendex/internal/cartrecovery/cartrecoverysrv"
)

// Container wires together all cart recovery domain dependencies.
type Container struct {
	Service *cartrecoverysrv.Service
	Handler *cartrecoveryapi.Handler
}

// New creates a fully-wired cart recovery container.
func New(db *sqlx.DB) *Container {
	repo := cartrecoveryinfra.NewPostgresRepo(db)
	svc := cartrecoverysrv.New(repo)
	handler := cartrecoveryapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers all cart recovery admin routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
