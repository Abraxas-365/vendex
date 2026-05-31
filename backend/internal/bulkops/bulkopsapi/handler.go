package bulkopsapi

import (
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/bulkops"
	"github.com/Abraxas-365/hada-commerce/internal/bulkops/bulkopssrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the bulk operations domain.
type Handler struct {
	svc *bulkopssrv.Service
}

// NewHandler creates a new bulk operations API handler.
func NewHandler(svc *bulkopssrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected bulk operation routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/bulk-operations")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Get("/:id/items", h.ListItems)
	g.Post("/:id/process", h.Process)
	g.Post("/:id/cancel", h.Cancel)
}

// ---------------------------------------------------------------------------
// Handlers
// ---------------------------------------------------------------------------

// Create handles POST /bulk-operations.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var body struct {
		Type         string                 `json:"type"`
		ResourceType string                 `json:"resource_type"`
		ResourceIDs  []string               `json:"resource_ids"`
		Parameters   map[string]interface{} `json:"parameters"`
	}
	if err := c.BodyParser(&body); err != nil {
		return errx.Validation("invalid request body")
	}

	if body.Type == "" {
		return errx.Validation("type is required")
	}
	if body.ResourceType == "" {
		return errx.Validation("resource_type is required")
	}
	if len(body.ResourceIDs) == 0 {
		return errx.Validation("resource_ids must not be empty")
	}

	createdBy := ""
	if authCtx.UserID != nil {
		createdBy = string(*authCtx.UserID)
	}

	op, err := h.svc.Create(c.Context(), authCtx.TenantID, bulkops.CreateInput{
		Type:         bulkops.OperationType(body.Type),
		ResourceType: body.ResourceType,
		ResourceIDs:  body.ResourceIDs,
		Parameters:   body.Parameters,
		CreatedBy:    createdBy,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(op)
}

// List handles GET /bulk-operations.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	page, pageSize := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetByID handles GET /bulk-operations/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BulkOperationID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("id is required")
	}

	op, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(op)
}

// ListItems handles GET /bulk-operations/:id/items.
func (h *Handler) ListItems(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BulkOperationID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("id is required")
	}

	page, pageSize := paginationFromQuery(c)

	result, err := h.svc.ListItems(c.Context(), authCtx.TenantID, id, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// Process handles POST /bulk-operations/:id/process.
func (h *Handler) Process(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BulkOperationID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("id is required")
	}

	op, err := h.svc.Process(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(op)
}

// Cancel handles POST /bulk-operations/:id/cancel.
func (h *Handler) Cancel(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BulkOperationID(c.Params("id"))

	if id.IsEmpty() {
		return errx.Validation("id is required")
	}

	if err := h.svc.Cancel(c.Context(), authCtx.TenantID, id); err != nil {
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
