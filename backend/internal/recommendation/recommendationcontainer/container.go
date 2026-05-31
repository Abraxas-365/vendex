package recommendationcontainer

import (
	"github.com/Abraxas-365/vendex/internal/recommendation/recommendationapi"
	"github.com/Abraxas-365/vendex/internal/recommendation/recommendationinfra"
	"github.com/Abraxas-365/vendex/internal/recommendation/recommendationsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all recommendation domain dependencies.
type Container struct {
	Service *recommendationsrv.Service
	Handler *recommendationapi.Handler
}

// New builds the full recommendation dependency graph.
func New(db *sqlx.DB) *Container {
	repo := recommendationinfra.NewPostgresRepository(db)
	svc := recommendationsrv.New(repo)
	handler := recommendationapi.New(svc)
	return &Container{Service: svc, Handler: handler}
}

// RegisterRoutes wires all protected (admin) recommendation routes.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes wires all public recommendation routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
