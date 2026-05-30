package returnsapi

import (
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/returns"
	"github.com/Abraxas-365/hada-commerce/internal/returns/returnssrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the returns domain.
type Handler struct {
	svc *returnssrv.Service
}

// NewHandler creates a new returns API handler.
func NewHandler(svc *returnssrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected (admin) return routes.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/returns")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id/approve", h.Approve)
	g.Put("/:id/reject", h.Reject)
	g.Put("/:id/received", h.MarkReceived)
	g.Put("/:id/refunded", h.MarkRefunded)
	g.Put("/:id/close", h.Close)
}

// RegisterPublicRoutes registers public (customer-facing) return routes.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/returns")
	g.Post("/request", h.CreatePublic)
	g.Get("/order/:orderId", h.ListByOrder)
}

// -----------------------------------------------------------------------
// Protected handlers
// -----------------------------------------------------------------------

type createReturnItemRequest struct {
	ProductID string `json:"product_id"`
	VariantID string `json:"variant_id,omitempty"`
	Quantity  int    `json:"quantity"`
	Reason    string `json:"reason,omitempty"`
	Condition string `json:"condition,omitempty"`
}

type createReturnRequest struct {
	OrderID    string                    `json:"order_id"`
	CustomerID string                    `json:"customer_id"`
	Reason     string                    `json:"reason"`
	Notes      string                    `json:"notes,omitempty"`
	Items      []createReturnItemRequest `json:"items"`
}

// Create handles POST /returns (admin-facing).
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	in, err := buildCreateInput(req)
	if err != nil {
		return err
	}

	result, err := h.svc.CreateReturn(c.Context(), authCtx.TenantID, in)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// GetByID handles GET /returns/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReturnID(c.Params("id"))

	result, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// List handles GET /returns.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	status := c.Query("status")

	result, err := h.svc.List(c.Context(), authCtx.TenantID, status, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

type approveRequest struct {
	AdminNotes     string `json:"admin_notes,omitempty"`
	Resolution     string `json:"resolution"`
	RefundCents    int64  `json:"refund_amount_cents"`
	RefundCurrency string `json:"refund_currency,omitempty"`
}

// Approve handles PUT /returns/:id/approve.
func (h *Handler) Approve(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReturnID(c.Params("id"))

	var req approveRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	in := returns.ApproveInput{
		AdminNotes:     req.AdminNotes,
		Resolution:     returns.Resolution(req.Resolution),
		RefundCents:    req.RefundCents,
		RefundCurrency: req.RefundCurrency,
	}

	result, err := h.svc.Approve(c.Context(), authCtx.TenantID, id, in)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

type rejectRequest struct {
	AdminNotes string `json:"admin_notes,omitempty"`
}

// Reject handles PUT /returns/:id/reject.
func (h *Handler) Reject(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReturnID(c.Params("id"))

	var req rejectRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	result, err := h.svc.Reject(c.Context(), authCtx.TenantID, id, req.AdminNotes)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// MarkReceived handles PUT /returns/:id/received.
func (h *Handler) MarkReceived(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReturnID(c.Params("id"))

	result, err := h.svc.MarkReceived(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// MarkRefunded handles PUT /returns/:id/refunded.
func (h *Handler) MarkRefunded(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReturnID(c.Params("id"))

	result, err := h.svc.MarkRefunded(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// Close handles PUT /returns/:id/close.
func (h *Handler) Close(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReturnID(c.Params("id"))

	result, err := h.svc.Close(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// -----------------------------------------------------------------------
// Public handlers
// -----------------------------------------------------------------------

type publicCreateReturnRequest struct {
	OrderID    string                    `json:"order_id"`
	CustomerID string                    `json:"customer_id"`
	Reason     string                    `json:"reason"`
	Notes      string                    `json:"notes,omitempty"`
	Items      []createReturnItemRequest `json:"items"`
}

// CreatePublic handles POST /returns/request (customer-facing).
func (h *Handler) CreatePublic(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID is required", errx.TypeAuthorization)
	}

	var req publicCreateReturnRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	in, err := buildCreateInput(createReturnRequest(req))
	if err != nil {
		return err
	}

	result, err := h.svc.CreateReturn(c.Context(), tenantID, in)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(result)
}

// ListByOrder handles GET /returns/order/:orderId (public).
func (h *Handler) ListByOrder(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID is required", errx.TypeAuthorization)
	}

	orderID := kernel.OrderID(c.Params("orderId"))
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	result, err := h.svc.ListByOrder(c.Context(), tenantID, orderID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// -----------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------

func buildCreateInput(req createReturnRequest) (returns.CreateReturnInput, error) {
	if req.Reason == "" {
		return returns.CreateReturnInput{}, errx.New("reason is required", errx.TypeValidation)
	}
	if len(req.Items) == 0 {
		return returns.CreateReturnInput{}, errx.New("at least one item is required", errx.TypeValidation)
	}

	items := make([]returns.ReturnItemInput, len(req.Items))
	for i, it := range req.Items {
		if it.Quantity < 1 {
			it.Quantity = 1
		}
		cond := returns.ItemCondition(it.Condition)
		if cond == "" {
			cond = returns.ConditionUnopened
		}
		items[i] = returns.ReturnItemInput{
			ProductID: kernel.ProductID(it.ProductID),
			VariantID: kernel.VariantID(it.VariantID),
			Quantity:  it.Quantity,
			Reason:    it.Reason,
			Condition: cond,
		}
	}

	return returns.CreateReturnInput{
		OrderID:    kernel.OrderID(req.OrderID),
		CustomerID: kernel.CustomerID(req.CustomerID),
		Reason:     req.Reason,
		Notes:      req.Notes,
		Items:      items,
	}, nil
}
