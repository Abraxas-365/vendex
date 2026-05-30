package searchcontainer

import (
	"github.com/Abraxas-365/hada-commerce/internal/search/searchapi"
	"github.com/Abraxas-365/hada-commerce/internal/search/searchinfra"
	"github.com/Abraxas-365/hada-commerce/internal/search/searchsrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all search domain dependencies.
type Container struct {
	Service *searchsrv.Service
	Handler *searchapi.Handler
}

// New creates a fully-wired search container.
func New(db *sqlx.DB) *Container {
	repo := searchinfra.NewPostgresRepo(db)
	svc := searchsrv.New(repo)
	handler := searchapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterPublicRoutes registers public search routes on the given router.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
