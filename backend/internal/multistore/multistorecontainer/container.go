package multistorecontainer

import (
	"database/sql"

	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/multistore"
	"github.com/Abraxas-365/vendex/internal/multistore/multistoreapi"
	"github.com/Abraxas-365/vendex/internal/multistore/multistoreinfra"
	"github.com/Abraxas-365/vendex/internal/multistore/multistoresrv"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

// Container wires together all multistore domain dependencies.
type Container struct {
	Service          *multistoresrv.Service
	Handler          *multistoreapi.Handler
	TenantResolver   *multistore.TenantResolver
	TenantMiddleware *multistoreapi.TenantMiddleware
}

// New creates a fully-wired multistore container.
// baseDomain is the platform domain (e.g. "vendex.ai") used for subdomain resolution.
func New(db *sqlx.DB, rawDB *sql.DB, bus eventbus.Bus, baseDomain string) *Container {
	repo := multistoreinfra.NewPostgresRepo(db)
	domainRepo := multistoreinfra.NewDomainRepo(rawDB)
	svc := multistoresrv.New(repo, bus)
	handler := multistoreapi.NewHandler(svc)
	resolver := multistore.NewTenantResolver(baseDomain, domainRepo)
	tenantMw := multistoreapi.NewTenantMiddleware(resolver)

	return &Container{
		Service:          svc,
		Handler:          handler,
		TenantResolver:   resolver,
		TenantMiddleware: tenantMw,
	}
}

// RegisterRoutes registers protected admin routes on the given router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}

// RegisterPublicRoutes registers unauthenticated routes on the given router.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
