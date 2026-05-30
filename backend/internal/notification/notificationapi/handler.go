package notificationapi

import (
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/notification/notificationsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the notification domain.
type Handler struct {
	svc *notificationsrv.Service
}

// NewHandler creates a new notification API handler.
func NewHandler(svc *notificationsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected notification routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/notifications")
	g.Get("/", h.List)
	g.Get("/unread-count", h.UnreadCount)
	g.Put("/read-all", h.MarkAllRead)
	g.Put("/:id/read", h.MarkRead)
	g.Delete("/:id", h.Delete)
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// List handles GET /notifications — lists notifications for the authenticated user.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	if authCtx.UserID == nil {
		return errx.Unauthorized("user identity required")
	}

	unreadOnly := c.Query("unread_only") == "true" || c.Query("unread_only") == "1"
	page, pageSize := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, *authCtx.UserID, unreadOnly, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// UnreadCount handles GET /notifications/unread-count.
func (h *Handler) UnreadCount(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	if authCtx.UserID == nil {
		return errx.Unauthorized("user identity required")
	}

	count, err := h.svc.GetUnreadCount(c.Context(), authCtx.TenantID, *authCtx.UserID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"unread_count": count})
}

// MarkRead handles PUT /notifications/:id/read.
func (h *Handler) MarkRead(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewNotificationID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("notification id is required")
	}

	if err := h.svc.MarkRead(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// MarkAllRead handles PUT /notifications/read-all.
func (h *Handler) MarkAllRead(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	if authCtx.UserID == nil {
		return errx.Unauthorized("user identity required")
	}

	if err := h.svc.MarkAllRead(c.Context(), authCtx.TenantID, *authCtx.UserID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Delete handles DELETE /notifications/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewNotificationID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("notification id is required")
	}

	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func paginationFromQuery(c *fiber.Ctx) (page, pageSize int) {
	page, _ = strconv.Atoi(c.Query("page"))
	pageSize, _ = strconv.Atoi(c.Query("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return
}
