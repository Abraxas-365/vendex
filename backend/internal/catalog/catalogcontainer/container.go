package catalogcontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/catalog/catalogapi"
	"github.com/Abraxas-365/hada-commerce/internal/catalog/cataloginfra"
	"github.com/Abraxas-365/hada-commerce/internal/catalog/catalogsrv"
)

// Container wires together all catalog domain dependencies.
type Container struct {
	Service *catalogsrv.Service
	Handler *catalogapi.Handler
}

// New creates a fully-wired catalog container.
func New(db *sqlx.DB) *Container {
	categoryRepo := cataloginfra.NewCategoryPostgresRepo(db)
	collectionRepo := cataloginfra.NewCollectionPostgresRepo(db)
	svc := catalogsrv.New(categoryRepo, collectionRepo)
	handler := catalogapi.NewHandler(svc)
	return &Container{
		Service: svc,
		Handler: handler,
	}
}

// RegisterRoutes registers catalog HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
