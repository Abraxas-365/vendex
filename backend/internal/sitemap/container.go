package sitemap

import "github.com/gofiber/fiber/v2"

// Container wires together all sitemap dependencies.
type Container struct {
	Handler *Handler
}

// New creates a fully-wired sitemap Container.
func New(products ProductLister, catalog CategoryLister, storefront PageLister) *Container {
	svc := NewService(products, catalog, storefront)
	handler := NewHandler(svc)
	return &Container{Handler: handler}
}

// RegisterPublicRoutes registers the public sitemap route on the given router.
func (c *Container) RegisterPublicRoutes(router fiber.Router) {
	c.Handler.RegisterPublicRoutes(router)
}
