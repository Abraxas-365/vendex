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

// Deps holds all optional resolver dependencies for the storefront renderer.
// All fields are optional — the renderer degrades gracefully when nil.
type Deps struct {
	ProductLister    renderer.ProductLister
	CollectionGetter renderer.CollectionGetter
	SettingsGetter   renderer.SettingsGetter
}

// New creates a fully-wired storefront container.
//
// themeGetter is used by the renderer to resolve the active theme for HTML rendering.
// Pass nil to disable HTML rendering (JSON-only mode).
//
// deps provides optional data resolvers for live product/collection/settings data.
// Use an empty Deps{} if no data resolvers are available yet.
func New(db *sqlx.DB, bus eventbus.Bus, themeGetter renderer.ThemeGetter, deps Deps) *Container {
	pagesRepo := storefrontinfra.NewPagePostgresRepo(db)
	versionsRepo := storefrontinfra.NewPageVersionPostgresRepo(db)
	blockTypesRepo := storefrontinfra.NewBlockTypePostgresRepo(db)
	svc := storefrontsrv.New(pagesRepo, versionsRepo, blockTypesRepo, bus)

	var r *renderer.Renderer
	if themeGetter != nil {
		navRepo := storefrontinfra.NewNavMenuPostgresRepo(db)
		overrideRepo := storefrontinfra.NewTemplateOverridePostgresRepo(db)

		r = renderer.NewWithConfig(themeGetter, renderer.Config{
			ProductLister:    deps.ProductLister,
			CollectionGetter: deps.CollectionGetter,
			SettingsGetter:   deps.SettingsGetter,
			NavRepo:          navRepo,
			OverrideRepo:     overrideRepo,
		})
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
