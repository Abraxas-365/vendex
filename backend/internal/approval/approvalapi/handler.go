// Package approvalapi provides HTTP handlers for the approval workflow domain.
package approvalapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/approval/approvalsrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// tenantID extracts the tenant ID from Fiber context (set by auth middleware).
func tenantID(c *fiber.Ctx) kernel.TenantID {
	return kernel.TenantID(c.Locals("tenant_id").(string))
}

// reviewedBy extracts the acting user identifier from Fiber context.
// Falls back to empty string if not present.
func reviewedBy(c *fiber.Ctx) string {
	v := c.Locals("user_id")
	if v == nil {
		return ""
	}
	s, ok := v.(string)
	if !ok {
		return ""
	}
	return s
}

// Handler exposes approval workflow endpoints over HTTP.
type Handler struct {
	svc *approvalsrv.Service
}

// NewHandler creates a new approval Handler.
func NewHandler(svc *approvalsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers approval routes on the given router group.
// Expected to be called with a group like "/approvals".
func (h *Handler) RegisterRoutes(r fiber.Router) {
	r.Get("/count", h.CountPending) // must be before /:id to avoid shadowing
	r.Get("/", h.List)
	r.Get("/:id", h.GetByID)
	r.Post("/:id/approve", h.Approve)
	r.Post("/:id/reject", h.Reject)
}

// List returns paginated approval requests, optionally filtered by status query param.
// Query params: status (optional), page (default 1), page_size (default 20).
func (h *Handler) List(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))
	status := c.Query("status", "")
	p := kernel.NewPaginationOptions(page, pageSize)

	result, err := h.svc.List(c.Context(), tenantID(c), status, p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// GetByID returns a single approval request with full tool input.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := kernel.ApprovalRequestID(c.Params("id"))

	req, err := h.svc.GetByID(c.Context(), tenantID(c), id)
	if err != nil {
		return err
	}
	return c.JSON(req)
}

// approveRejectBody is the JSON body accepted by Approve and Reject.
type approveRejectBody struct {
	Reason string `json:"reason"`
}

// Approve marks a pending request as approved.
func (h *Handler) Approve(c *fiber.Ctx) error {
	id := kernel.ApprovalRequestID(c.Params("id"))

	var body approveRejectBody
	if err := c.BodyParser(&body); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	req, err := h.svc.Approve(c.Context(), tenantID(c), id, reviewedBy(c), body.Reason)
	if err != nil {
		return err
	}
	return c.JSON(req)
}

// Reject marks a pending request as rejected.
func (h *Handler) Reject(c *fiber.Ctx) error {
	id := kernel.ApprovalRequestID(c.Params("id"))

	var body approveRejectBody
	if err := c.BodyParser(&body); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	req, err := h.svc.Reject(c.Context(), tenantID(c), id, reviewedBy(c), body.Reason)
	if err != nil {
		return err
	}
	return c.JSON(req)
}

// CountPending returns the count of pending approval requests for the tenant.
// Response: {"count": N}
func (h *Handler) CountPending(c *fiber.Ctx) error {
	count, err := h.svc.CountPending(c.Context(), tenantID(c))
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"count": count})
}
