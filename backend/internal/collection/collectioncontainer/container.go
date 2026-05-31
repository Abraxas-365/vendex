package collectioncontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/collection/collectionapi"
	"github.com/Abraxas-365/hada-commerce/internal/collection/collectioninfra"
	"github.com/Abraxas-365/hada-commerce/internal/collection/collectionsrv"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
)

// Container wires the collection domain dependencies together.
type Container struct {
	Service *collectionsrv.Service
	Handler *collectionapi.Handler
}

// New creates a fully-wired collection container.
func New(db *sqlx.DB, bus eventbus.Bus) *Container {
	repo := collectioninfra.NewPostgresRepo(db)
	svc := collectionsrv.New(repo, bus)
	handler := collectionapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers admin (authenticated) collection routes.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers unauthenticated collection routes.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
