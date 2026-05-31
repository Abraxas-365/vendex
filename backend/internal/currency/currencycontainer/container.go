package currencycontainer

import (
	"github.com/Abraxas-365/vendex/internal/currency/currencyapi"
	"github.com/Abraxas-365/vendex/internal/currency/currencyinfra"
	"github.com/Abraxas-365/vendex/internal/currency/currencysrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all currency domain dependencies.
type Container struct {
	Service *currencysrv.Service
	Handler *currencyapi.Handler
}

// New creates a fully-wired currency container.
func New(db *sqlx.DB) *Container {
	repo := currencyinfra.NewPostgresRepo(db)
	svc := currencysrv.New(repo)
	handler := currencyapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers protected currency admin HTTP routes.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers public currency HTTP routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
