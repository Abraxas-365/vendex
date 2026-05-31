package subscriptionapi

import (
	"strconv"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/subscription"
	"github.com/Abraxas-365/vendex/internal/subscription/subscriptionsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the subscription domain.
type Handler struct {
	svc *subscriptionsrv.Service
}

// NewHandler creates a new subscription API handler.
func NewHandler(svc *subscriptionsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers admin subscription routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/subscriptions")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/due", h.ListDue)
	g.Get("/:id", h.GetByID)
	g.Post("/:id/cancel", h.Cancel)
	g.Post("/:id/pause", h.Pause)
	g.Post("/:id/resume", h.Resume)
	g.Get("/:id/billing", h.ListBillingRecords)
}

// RegisterCustomerRoutes registers customer-facing subscription routes.
func (h *Handler) RegisterCustomerRoutes(router fiber.Router) {
	g := router.Group("/storefront/subscriptions")
	g.Get("/", h.ListByCustomer)
}

// ---------------------------------------------------------------------------
// Request types
// ---------------------------------------------------------------------------

type createRequest struct {
	CustomerID  string            `json:"customer_id"`
	ProductID   string            `json:"product_id"`
	VariantID   *string           `json:"variant_id,omitempty"`
	PriceAmount int64             `json:"price_amount"`
	Currency    string            `json:"currency"`
	Interval    string            `json:"interval"`
	TrialEndsAt *time.Time        `json:"trial_ends_at,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

type recordBillingRequest struct {
	PriceAmount   int64   `json:"price_amount"`
	Currency      string  `json:"currency"`
	Status        string  `json:"status"`
	OrderID       *string `json:"order_id,omitempty"`
	FailureReason *string `json:"failure_reason,omitempty"`
}

// ---------------------------------------------------------------------------
// Admin handlers
// ---------------------------------------------------------------------------

// Create handles POST /subscriptions.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.CustomerID == "" || req.ProductID == "" || req.Interval == "" || req.Currency == "" {
		return errx.New("customer_id, product_id, interval, and currency are required", errx.TypeValidation)
	}

	input := subscriptionsrv.CreateInput{
		CustomerID:  kernel.NewCustomerID(req.CustomerID),
		ProductID:   kernel.NewProductID(req.ProductID),
		Price:       kernel.NewMoney(req.PriceAmount, req.Currency),
		Interval:    subscription.BillingInterval(req.Interval),
		TrialEndsAt: req.TrialEndsAt,
		Metadata:    req.Metadata,
	}

	if req.VariantID != nil {
		v := kernel.NewVariantID(*req.VariantID)
		input.VariantID = &v
	}

	sub, err := h.svc.Create(c.Context(), authCtx.TenantID, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(sub)
}

// GetByID handles GET /subscriptions/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewSubscriptionID(c.Params("id"))

	sub, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(sub)
}

// List handles GET /subscriptions.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	page, pageSize := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// ListDue handles GET /subscriptions/due.
func (h *Handler) ListDue(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	subs, err := h.svc.ListDueBilling(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": subs})
}

// Cancel handles POST /subscriptions/:id/cancel.
func (h *Handler) Cancel(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewSubscriptionID(c.Params("id"))

	sub, err := h.svc.Cancel(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(sub)
}

// Pause handles POST /subscriptions/:id/pause.
func (h *Handler) Pause(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewSubscriptionID(c.Params("id"))

	sub, err := h.svc.Pause(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(sub)
}

// Resume handles POST /subscriptions/:id/resume.
func (h *Handler) Resume(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewSubscriptionID(c.Params("id"))

	sub, err := h.svc.Resume(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(sub)
}

// ListBillingRecords handles GET /subscriptions/:id/billing.
func (h *Handler) ListBillingRecords(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.NewSubscriptionID(c.Params("id"))
	page, pageSize := paginationFromQuery(c)

	result, err := h.svc.ListBillingRecords(c.Context(), authCtx.TenantID, id, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// ---------------------------------------------------------------------------
// Customer (storefront) handlers
// ---------------------------------------------------------------------------

// ListByCustomer handles GET /storefront/subscriptions.
// Extracts customerID from the "customer_id" query param (typically set by customer auth middleware).
func (h *Handler) ListByCustomer(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	customerID := c.Query("customer_id")
	if customerID == "" {
		return errx.New("customer_id query parameter is required", errx.TypeValidation)
	}

	subs, err := h.svc.ListByCustomer(c.Context(), authCtx.TenantID, kernel.NewCustomerID(customerID))
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"items": subs})
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
