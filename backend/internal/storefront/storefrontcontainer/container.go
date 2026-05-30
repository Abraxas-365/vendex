package storefrontcontainer

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/renderer"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontapi"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontinfra"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontsrv"
)

// Container wires together all storefront domain dependencies.
type Container struct {
	Service  *storefrontsrv.Service
	Handler  *storefrontapi.Handler
	Renderer *renderer.Renderer
}

// New creates a fully-wired storefront container.
// themeGetter is used by the renderer to resolve the active theme for HTML rendering.
// Pass nil to disable HTML rendering (JSON-only mode).
func New(db *sqlx.DB, bus eventbus.Bus, themeGetter renderer.ThemeGetter) *Container {
	pagesRepo := storefrontinfra.NewPagePostgresRepo(db)
	versionsRepo := storefrontinfra.NewPageVersionPostgresRepo(db)
	blockTypesRepo := storefrontinfra.NewBlockTypePostgresRepo(db)
	svc := storefrontsrv.New(pagesRepo, versionsRepo, blockTypesRepo, bus)

	var r *renderer.Renderer
	if themeGetter != nil {
		r = renderer.New(themeGetter)
	}

	handler := storefrontapi.NewHandler(svc, r)
	return &Container{
		Service:  svc,
		Handler:  handler,
		Renderer: r,
	}
}

// RegisterRoutes registers storefront HTTP routes on the given Fiber router.
func (c *Container) RegisterRoutes(router fiber.Router) {
	c.Handler.RegisterRoutes(router)
}
