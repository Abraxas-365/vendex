package sitemap

import (
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes the sitemap.xml HTTP endpoint.
type Handler struct {
	svc *Service
}

// NewHandler creates a new sitemap Handler.
func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterPublicRoutes registers the public sitemap route.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Get("/sitemap.xml", h.HandleSitemap)
}

// HandleSitemap handles GET /sitemap.xml.
// Required headers:
//   - X-Tenant-ID  — identifies the store tenant
//   - X-Base-URL   — base URL of the storefront (default: https://store.example.com)
func (h *Handler) HandleSitemap(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	baseURL := c.Get("X-Base-URL")
	if baseURL == "" {
		baseURL = "https://store.example.com"
	}

	sm, err := h.svc.Generate(c.Context(), tenantID, baseURL)
	if err != nil {
		return errx.Wrap(err, "failed to generate sitemap", errx.TypeInternal)
	}

	c.Set(fiber.HeaderContentType, "application/xml; charset=utf-8")
	return c.SendString(sm.ToXML())
}
