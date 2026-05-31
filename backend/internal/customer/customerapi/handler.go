package customerapi

import (
	"strconv"

	"github.com/Abraxas-365/vendex/internal/customer"
	"github.com/Abraxas-365/vendex/internal/customer/customersrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the customer domain.
type Handler struct {
	svc *customersrv.Service
}

// NewHandler creates a new customer API handler.
func NewHandler(svc *customersrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all customer routes on the given router.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/customers")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
}

type createRequest struct {
	Email     string            `json:"email"`
	Name      string            `json:"name"`
	Phone     string            `json:"phone"`
	Addresses []customer.Address `json:"addresses"`
}

// Create handles POST /customers.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.Email == "" {
		return errx.New("email is required", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}

	cust, err := h.svc.Create(c.Context(), authCtx.TenantID, customersrv.CreateInput{
		Email:     req.Email,
		Name:      req.Name,
		Phone:     req.Phone,
		Addresses: req.Addresses,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(cust)
}

// GetByID handles GET /customers/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CustomerID(c.Params("id"))

	cust, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(cust)
}

type updateRequest struct {
	Name  string             `json:"name"`
	Phone string             `json:"phone"`
	Email string             `json:"email"`
}

// Update handles PUT /customers/:id.
func (h *Handler) Update(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CustomerID(c.Params("id"))

	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	cust, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	if req.Name != "" {
		cust.Name = req.Name
	}
	if req.Phone != "" {
		cust.Phone = req.Phone
	}
	if req.Email != "" {
		email := kernel.NewEmail(req.Email)
		if email.IsEmpty() {
			return customer.ErrInvalidEmail
		}
		cust.Email = email
	}

	if err := h.svc.Update(c.Context(), cust); err != nil {
		return err
	}

	return c.JSON(cust)
}

// Delete handles DELETE /customers/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.CustomerID(c.Params("id"))

	if err := h.svc.Delete(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// List handles GET /customers.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
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
