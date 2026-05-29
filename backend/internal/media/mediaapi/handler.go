package mediaapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/media"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediasrv"
)

// Handler exposes HTTP endpoints for the media domain.
type Handler struct {
	svc *mediasrv.Service
}

// NewHandler creates a new media API handler.
func NewHandler(svc *mediasrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all media routes on the given Fiber router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/media")
	g.Post("/upload", h.Upload)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Delete("/:id", h.Delete)
}

// Upload handles POST /media/upload.
// Accepts a multipart form with a "file" field.
func (h *Handler) Upload(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return media.ErrInvalidFile
	}

	src, err := file.Open()
	if err != nil {
		return media.ErrUploadFailed
	}
	defer src.Close()

	altText := c.FormValue("alt")
	result, err := h.svc.Upload(c.Context(), mediasrv.UploadInput{
		TenantID:    authCtx.TenantID,
		Filename:    file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        file.Size,
		Alt:         altText,
		UploadedBy:  string(authCtx.UserID),
		Data:        src,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// List handles GET /media.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	p := kernel.PaginationOptions{Page: page, PageSize: pageSize}

	result, err := h.svc.List(c.Context(), authCtx.TenantID, p)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// GetByID handles GET /media/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	id := kernel.MediaID(c.Params("id"))
	m, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(m)
}

// Delete handles DELETE /media/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx, ok := c.Locals("auth").(*kernel.AuthContext)
	if !ok || authCtx == nil {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	id := kernel.MediaID(c.Params("id"))
	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}
