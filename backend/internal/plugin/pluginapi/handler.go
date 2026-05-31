package pluginapi

import (
	"encoding/json"
	"strconv"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/plugin"
	"github.com/Abraxas-365/vendex/internal/plugin/pluginsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the plugin domain.
type Handler struct {
	svc *pluginsrv.Service
}

// NewHandler creates a new plugin API handler.
func NewHandler(svc *pluginsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all plugin routes on the given router.
// NOTE: /plugins/installed and /plugins/js-manifest must be registered BEFORE
// /plugins/:id to avoid the path collision where those names are matched as IDs.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/plugins")

	// Static sub-paths MUST come before the parametric /:id route.
	g.Get("/installed", h.ListInstalled)
	g.Get("/js-manifest", h.JSManifest)
	g.Get("/", h.ListPlugins)
	g.Post("/", h.CreatePlugin)
	g.Post("/install", h.Install)

	// Parametric routes
	g.Get("/:id/manifest", h.GetManifest)
	g.Get("/:id", h.GetPlugin)
	g.Post("/:id/uninstall", h.Uninstall)
	g.Put("/:id/settings", h.UpdateSettings)
	g.Post("/:id/enable", h.Enable)
	g.Post("/:id/disable", h.Disable)
	g.Post("/:id/versions", h.CreateVersion)
}

// RegisterPublicRoutes registers unauthenticated plugin routes.
// These routes read the tenant from the X-Tenant-ID header directly.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/plugins")
	g.Get("/js-manifest", h.JSManifest)
}

// ListPlugins handles GET /plugins — returns the global plugin catalogue.
func (h *Handler) ListPlugins(c *fiber.Ctx) error {
	pg := paginationFromQuery(c)
	result, err := h.svc.ListPlugins(c.Context(), pg)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// GetPlugin handles GET /plugins/:id
func (h *Handler) GetPlugin(c *fiber.Ctx) error {
	id := kernel.PluginID(c.Params("id"))
	p, err := h.svc.GetPlugin(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(p)
}

// createPluginRequest is the JSON body for creating a plugin.
type createPluginRequest struct {
	Name        string          `json:"name"`
	DisplayName string          `json:"display_name"`
	Description string          `json:"description"`
	Author      string          `json:"author"`
	Icon        string          `json:"icon"`
	Category    string          `json:"category"`
	Tags        json.RawMessage `json:"tags"`
}

// CreatePlugin handles POST /plugins (admin: adds plugin to global catalogue).
func (h *Handler) CreatePlugin(c *fiber.Ctx) error {
	var req createPluginRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}

	p := &plugin.Plugin{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Author:      req.Author,
		Icon:        req.Icon,
		Category:    req.Category,
		Tags:        req.Tags,
	}

	result, err := h.svc.CreatePlugin(c.Context(), p)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

// ListInstalled handles GET /plugins/installed — tenant-scoped.
func (h *Handler) ListInstalled(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromQuery(c)

	result, err := h.svc.ListInstalled(c.Context(), authCtx.TenantID, pg)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// jsManifestResponse is the response body for the JS manifest endpoint.
type jsManifestResponse struct {
	Scripts []plugin.PluginScript `json:"scripts"`
}

// JSManifest handles GET /plugins/js-manifest.
// It returns an aggregated list of JS bundle URLs for all active plugin
// installations belonging to the tenant. The tenant is resolved from the
// authenticated context when available; otherwise it falls back to the
// X-Tenant-ID header so that public storefront renderers can call this endpoint
// without authentication.
func (h *Handler) JSManifest(c *fiber.Ctx) error {
	var tenantID kernel.TenantID

	// Prefer the auth context (authenticated routes), fall back to header (public routes).
	if auth, ok := c.Locals("auth").(*kernel.AuthContext); ok && auth != nil {
		tenantID = auth.TenantID
	} else {
		raw := c.Get("X-Tenant-ID")
		if raw == "" {
			return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
		}
		tenantID = kernel.TenantID(raw)
	}

	scripts, err := h.svc.GetJSManifest(c.Context(), tenantID)
	if err != nil {
		return err
	}

	return c.JSON(jsManifestResponse{Scripts: scripts})
}

// installRequest is the JSON body for installing a plugin.
type installRequest struct {
	PluginID  string `json:"plugin_id"`
	VersionID string `json:"version_id"`
}

// Install handles POST /plugins/install
func (h *Handler) Install(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req installRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.PluginID == "" {
		return errx.New("plugin_id is required", errx.TypeValidation)
	}

	installation, err := h.svc.Install(
		c.Context(),
		authCtx.TenantID,
		kernel.PluginID(req.PluginID),
		kernel.PluginVersionID(req.VersionID),
	)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(installation)
}

// Uninstall handles POST /plugins/:id/uninstall
func (h *Handler) Uninstall(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PluginID(c.Params("id"))

	if err := h.svc.Uninstall(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}

// UpdateSettings handles PUT /plugins/:id/settings
func (h *Handler) UpdateSettings(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PluginID(c.Params("id"))

	var settings json.RawMessage
	if err := c.BodyParser(&settings); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	installation, err := h.svc.UpdateSettings(c.Context(), authCtx.TenantID, id, settings)
	if err != nil {
		return err
	}
	return c.JSON(installation)
}

// Enable handles POST /plugins/:id/enable
func (h *Handler) Enable(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PluginID(c.Params("id"))

	installation, err := h.svc.Enable(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(installation)
}

// Disable handles POST /plugins/:id/disable
func (h *Handler) Disable(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.PluginID(c.Params("id"))

	installation, err := h.svc.Disable(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(installation)
}

// GetManifest handles GET /plugins/:id/manifest
func (h *Handler) GetManifest(c *fiber.Ctx) error {
	id := kernel.PluginID(c.Params("id"))

	manifest, err := h.svc.GetManifest(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"manifest": manifest})
}

// createVersionRequest is the JSON body for creating a plugin version.
type createVersionRequest struct {
	Version        string          `json:"version"`
	Changelog      string          `json:"changelog"`
	Permissions    json.RawMessage `json:"permissions"`
	ManifestJSON   string          `json:"manifest_json"`
	FrontendURL    string          `json:"frontend_url"`
	BackendEntry   string          `json:"backend_entry"`
	MinPlatformVer string          `json:"min_platform_ver"`
}

// CreateVersion handles POST /plugins/:id/versions
func (h *Handler) CreateVersion(c *fiber.Ctx) error {
	id := kernel.PluginID(c.Params("id"))

	var req createVersionRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Version == "" {
		return errx.New("version is required", errx.TypeValidation)
	}

	v := &plugin.PluginVersion{
		PluginID:       id,
		Version:        req.Version,
		Changelog:      req.Changelog,
		Permissions:    req.Permissions,
		ManifestJSON:   req.ManifestJSON,
		FrontendURL:    req.FrontendURL,
		BackendEntry:   req.BackendEntry,
		MinPlatformVer: req.MinPlatformVer,
	}

	result, err := h.svc.CreateVersion(c.Context(), v)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(result)
}

// --- helpers ---

func paginationFromQuery(c *fiber.Ctx) kernel.PaginationOptions {
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	opts := kernel.PaginationOptions{Page: page, PageSize: pageSize}
	if opts.Page < 1 {
		opts.Page = 1
	}
	if opts.PageSize < 1 {
		opts.PageSize = 20
	}
	if opts.PageSize > 100 {
		opts.PageSize = 100
	}
	return opts
}
