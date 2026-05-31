package customergroupapi

import (
	"github.com/Abraxas-365/vendex/internal/customergroup"
	"github.com/Abraxas-365/vendex/internal/customergroup/customergroupsrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the customer group domain.
type Handler struct {
	svc *customergroupsrv.Service
}

// NewHandler creates a new customer group API handler.
func NewHandler(svc *customergroupsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all customer-group routes on the given router.
// All routes are admin-protected (caller must pass an authenticated router).
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/customer-groups")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)

	// Membership sub-routes
	g.Post("/:id/members", h.AddMember)
	g.Delete("/:id/members/:customerId", h.RemoveMember)
	g.Get("/:id/members", h.ListMembers)

	// Customer-centric route
	router.Get("/customers/:customerId/groups", h.GetCustomerGroups)
}

// ─────────────────────────────────────────────────────────────────────────────
// Group CRUD
// ─────────────────────────────────────────────────────────────────────────────

type createGroupRequest struct {
	Name        string                    `json:"name"`
	Description string                    `json:"description"`
	Rules       customergroup.GroupRules   `json:"rules"`
	AutoAssign  bool                       `json:"auto_assign"`
}

// Create handles POST /customer-groups.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}

	group, err := h.svc.Create(c.Context(), authCtx.TenantID, customergroup.CreateGroupRequest{
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
		AutoAssign:  req.AutoAssign,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(group)
}

// GetByID handles GET /customer-groups/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CustomerGroupID(c.Params("id"))

	group, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(group)
}

// List handles GET /customer-groups.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	groups, err := h.svc.List(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": groups, "total": len(groups)})
}

type updateGroupRequest struct {
	Name        *string                    `json:"name,omitempty"`
	Description *string                    `json:"description,omitempty"`
	Rules       *customergroup.GroupRules   `json:"rules,omitempty"`
	AutoAssign  *bool                       `json:"auto_assign,omitempty"`
}

// Update handles PUT /customer-groups/:id.
func (h *Handler) Update(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CustomerGroupID(c.Params("id"))

	var req updateGroupRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	group, err := h.svc.Update(c.Context(), authCtx.TenantID, id, customergroup.UpdateGroupRequest{
		Name:        req.Name,
		Description: req.Description,
		Rules:       req.Rules,
		AutoAssign:  req.AutoAssign,
	})
	if err != nil {
		return err
	}

	return c.JSON(group)
}

// Delete handles DELETE /customer-groups/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CustomerGroupID(c.Params("id"))

	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ─────────────────────────────────────────────────────────────────────────────
// Membership
// ─────────────────────────────────────────────────────────────────────────────

type addMemberRequest struct {
	CustomerID string `json:"customer_id"`
}

// AddMember handles POST /customer-groups/:id/members.
func (h *Handler) AddMember(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	groupID := kernel.CustomerGroupID(c.Params("id"))

	var req addMemberRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.CustomerID == "" {
		return errx.New("customer_id is required", errx.TypeValidation)
	}

	membership, err := h.svc.AddMember(c.Context(), authCtx.TenantID, groupID, kernel.CustomerID(req.CustomerID))
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(membership)
}

// RemoveMember handles DELETE /customer-groups/:id/members/:customerId.
func (h *Handler) RemoveMember(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	groupID := kernel.CustomerGroupID(c.Params("id"))
	customerID := kernel.CustomerID(c.Params("customerId"))

	if err := h.svc.RemoveMember(c.Context(), authCtx.TenantID, groupID, customerID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ListMembers handles GET /customer-groups/:id/members.
func (h *Handler) ListMembers(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	groupID := kernel.CustomerGroupID(c.Params("id"))

	members, err := h.svc.ListMembers(c.Context(), authCtx.TenantID, groupID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": members, "total": len(members)})
}

// GetCustomerGroups handles GET /customers/:customerId/groups.
func (h *Handler) GetCustomerGroups(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	customerID := kernel.CustomerID(c.Params("customerId"))

	groups, err := h.svc.GetCustomerGroups(c.Context(), authCtx.TenantID, customerID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": groups, "total": len(groups)})
}
