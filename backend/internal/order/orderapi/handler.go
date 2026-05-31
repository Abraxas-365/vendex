package orderapi

import (
	"strconv"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/order"
	"github.com/Abraxas-365/vendex/internal/order/ordersrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the order domain.
type Handler struct {
	svc *ordersrv.Service
}

// NewHandler creates a new order API handler.
func NewHandler(svc *ordersrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all order routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/orders")
	g.Post("/", h.Create)
	g.Get("/:id", h.GetByID)
	g.Get("/", h.List)
	g.Put("/:id/status", h.UpdateStatus)
	g.Post("/:id/cancel", h.Cancel)
}

type createItemRequest struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	PriceAmount int64  `json:"price_amount"`
	Currency    string `json:"currency"`
}

type createRequest struct {
	CustomerID string              `json:"customer_id"`
	Items      []createItemRequest `json:"items"`
	Address    order.Address       `json:"shipping_address"`
}

// Create handles POST /orders.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	items := make([]ordersrv.CreateItemInput, len(req.Items))
	for i, it := range req.Items {
		items[i] = ordersrv.CreateItemInput{
			ProductID:   kernel.ProductID(it.ProductID),
			ProductName: it.ProductName,
			Quantity:    it.Quantity,
			UnitPrice:   kernel.NewMoney(it.PriceAmount, it.Currency),
		}
	}

	o, err := h.svc.Create(c.Context(), authCtx.TenantID, ordersrv.CreateInput{
		CustomerID:      kernel.CustomerID(req.CustomerID),
		Items:           items,
		ShippingAddress: req.Address,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(o)
}

// GetByID handles GET /orders/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.OrderID(c.Params("id"))

	o, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(o)
}

// List handles GET /orders.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromQuery(c)

	customerID := c.Query("customer_id")
	var result kernel.Paginated[order.Order]
	var err error

	if customerID != "" {
		result, err = h.svc.ListByCustomer(c.Context(), authCtx.TenantID, kernel.CustomerID(customerID), pg)
	} else {
		result, err = h.svc.List(c.Context(), authCtx.TenantID, pg)
	}
	if err != nil {
		return err
	}

	return c.JSON(result)
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

// UpdateStatus handles PUT /orders/:id/status.
func (h *Handler) UpdateStatus(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.OrderID(c.Params("id"))

	var req updateStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	o, err := h.svc.UpdateStatus(c.Context(), authCtx.TenantID, id, order.OrderStatus(req.Status))
	if err != nil {
		return err
	}

	return c.JSON(o)
}

// Cancel handles POST /orders/:id/cancel.
func (h *Handler) Cancel(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.OrderID(c.Params("id"))

	o, err := h.svc.Cancel(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(o)
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
