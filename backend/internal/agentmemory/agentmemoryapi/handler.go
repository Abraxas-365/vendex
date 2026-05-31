// Package agentmemoryapi provides HTTP handlers for the agent memory domain.
package agentmemoryapi

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/agentmemory"
	"github.com/Abraxas-365/vendex/internal/agentmemory/agentmemorysrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Handler exposes agent memory CRUD and search over HTTP.
type Handler struct {
	svc *agentmemorysrv.Service
}

// NewHandler creates a new agentmemory Handler.
func NewHandler(svc *agentmemorysrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers agent memory routes on the given router group.
// Expected to be called with a group rooted at "/agent/memories".
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/search", h.Search)
	r.Get("/", h.List)
	r.Post("/", h.Create)
	r.Get("/:id", h.GetByID)
	r.Put("/:id", h.Update)
	r.Delete("/:id", h.Delete)
}

func tenantID(c *fiber.Ctx) kernel.TenantID {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	return authCtx.TenantID
}

// List returns paginated memories for the tenant.
func (h *Handler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	p := kernel.NewPaginationOptions(page, pageSize)

	result, err := h.svc.List(c.Context(), tenantID(c), p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// Search queries memories by full-text, category, and/or tags.
func (h *Handler) Search(c *fiber.Ctx) error {
	query := c.Query("q")
	category := c.Query("category")
	tagsParam := c.Query("tags")

	var tags []string
	if tagsParam != "" {
		for _, t := range strings.Split(tagsParam, ",") {
			if t = strings.TrimSpace(t); t != "" {
				tags = append(tags, t)
			}
		}
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	p := kernel.NewPaginationOptions(page, pageSize)

	opts := agentmemory.MemorySearchOptions{
		Query:    query,
		Category: category,
		Tags:     tags,
	}

	result, err := h.svc.Search(c.Context(), tenantID(c), opts, p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// Create creates a new memory entry.
func (h *Handler) Create(c *fiber.Ctx) error {
	var req agentmemory.CreateMemoryRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	m, err := h.svc.Create(c.Context(), tenantID(c), req)
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(m)
}

// GetByID returns a single memory entry by ID.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := kernel.AgentMemoryID(c.Params("id"))

	m, err := h.svc.Get(c.Context(), tenantID(c), id)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Update applies a partial update to an existing memory entry.
func (h *Handler) Update(c *fiber.Ctx) error {
	id := kernel.AgentMemoryID(c.Params("id"))

	var req agentmemory.UpdateMemoryRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	m, err := h.svc.Update(c.Context(), tenantID(c), id, req)
	if err != nil {
		return err
	}
	return c.JSON(m)
}

// Delete removes a memory entry by ID.
func (h *Handler) Delete(c *fiber.Ctx) error {
	id := kernel.AgentMemoryID(c.Params("id"))

	if err := h.svc.Delete(c.Context(), tenantID(c), id); err != nil {
		return err
	}
	return c.SendStatus(fiber.StatusNoContent)
}
