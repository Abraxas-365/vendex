package themeapi

import (
	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/theme"
	"github.com/Abraxas-365/hada-commerce/internal/theme/themesrv"
)

// Handler exposes HTTP endpoints for the theme domain.
type Handler struct {
	svc *themesrv.Service
}

// NewHandler creates a new theme API handler.
func NewHandler(svc *themesrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all theme routes on the given Fiber router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	themes := router.Group("/themes")
	themes.Get("/", h.ListThemes)
	themes.Post("/", h.CreateTheme)
	themes.Get("/active", h.GetActiveTheme)
	themes.Get("/:id", h.GetTheme)
	themes.Put("/:id", h.UpdateTheme)
	themes.Post("/:id/activate", h.ActivateTheme)
	themes.Post("/:id/duplicate", h.DuplicateTheme)
	themes.Delete("/:id", h.DeleteTheme)
}

// RegisterPublicRoutes registers unauthenticated theme routes.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Get("/themes/active", h.GetPublicActiveTheme)
}

// --- request/response types ---

type createThemeRequest struct {
	Name   string             `json:"name"`
	Tokens *theme.ThemeTokens `json:"tokens,omitempty"`
}

type updateThemeRequest struct {
	Name   *string            `json:"name,omitempty"`
	Tokens *theme.ThemeTokens `json:"tokens,omitempty"`
}

type duplicateThemeRequest struct {
	Name string `json:"name"`
}

// --- handlers ---

// ListThemes handles GET /themes.
func (h *Handler) ListThemes(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	themes, err := h.svc.ListThemes(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(themes)
}

// CreateTheme handles POST /themes.
func (h *Handler) CreateTheme(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createThemeRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("theme name is required", errx.TypeValidation)
	}

	t, err := h.svc.CreateTheme(c.Context(), themesrv.CreateThemeInput{
		TenantID: authCtx.TenantID,
		Name:     req.Name,
		Tokens:   req.Tokens,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(t)
}

// GetActiveTheme handles GET /themes/active (authenticated).
func (h *Handler) GetActiveTheme(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	t, err := h.svc.GetActiveTheme(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(t)
}

// GetPublicActiveTheme handles GET /themes/active (public, reads tenant from X-Tenant-ID header).
func (h *Handler) GetPublicActiveTheme(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID required", errx.TypeValidation)
	}

	t, err := h.svc.GetActiveTheme(c.Context(), tenantID)
	if err != nil {
		return err
	}

	return c.JSON(t)
}

// GetTheme handles GET /themes/:id.
func (h *Handler) GetTheme(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ThemeID(c.Params("id"))

	t, err := h.svc.GetTheme(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(t)
}

// UpdateTheme handles PUT /themes/:id.
func (h *Handler) UpdateTheme(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ThemeID(c.Params("id"))

	var req updateThemeRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	t, err := h.svc.UpdateTheme(c.Context(), themesrv.UpdateThemeInput{
		TenantID: authCtx.TenantID,
		ID:       id,
		Name:     req.Name,
		Tokens:   req.Tokens,
	})
	if err != nil {
		return err
	}

	return c.JSON(t)
}

// ActivateTheme handles POST /themes/:id/activate.
func (h *Handler) ActivateTheme(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ThemeID(c.Params("id"))

	t, err := h.svc.ActivateTheme(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(t)
}

// DuplicateTheme handles POST /themes/:id/duplicate.
func (h *Handler) DuplicateTheme(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ThemeID(c.Params("id"))

	var req duplicateThemeRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("new theme name is required", errx.TypeValidation)
	}

	t, err := h.svc.DuplicateTheme(c.Context(), authCtx.TenantID, id, req.Name)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(t)
}

// DeleteTheme handles DELETE /themes/:id.
func (h *Handler) DeleteTheme(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ThemeID(c.Params("id"))

	if err := h.svc.DeleteTheme(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
