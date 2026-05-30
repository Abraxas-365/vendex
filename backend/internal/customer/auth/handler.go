package auth

import (
	"context"
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/customer"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/order"
	"github.com/gofiber/fiber/v2"
)

// OrderService is the interface for listing a customer's orders.
type OrderService interface {
	ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[order.Order], error)
}

// Handler exposes HTTP endpoints for customer authentication and account management.
type Handler struct {
	svc        *Service
	middleware *CustomerMiddleware
	orderSvc   OrderService
}

// NewHandler creates a new customer auth handler.
func NewHandler(svc *Service, middleware *CustomerMiddleware, orderSvc OrderService) *Handler {
	return &Handler{
		svc:        svc,
		middleware: middleware,
		orderSvc:   orderSvc,
	}
}

// RegisterPublicRoutes registers public (unauthenticated) auth routes.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/storefront/auth")
	g.Post("/register", h.Register)
	g.Post("/login", h.Login)
	g.Post("/refresh", h.Refresh)
}

// RegisterProtectedRoutes registers customer-authenticated account routes.
// The router should already have CustomerMiddleware applied.
func (h *Handler) RegisterProtectedRoutes(router fiber.Router) {
	g := router.Group("/storefront/account")
	g.Get("/profile", h.GetProfile)
	g.Put("/profile", h.UpdateProfile)
	g.Put("/password", h.ChangePassword)
	g.Get("/orders", h.ListOrders)
	g.Get("/addresses", h.ListAddresses)
	g.Post("/addresses", h.AddAddress)
	g.Put("/addresses/:idx/default", h.SetDefaultAddress)
}

// ============================================================================
// Public handlers
// ============================================================================

// Register handles POST /storefront/auth/register.
func (h *Handler) Register(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	resp, err := h.svc.Register(c.Context(), tenantID, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(resp)
}

// Login handles POST /storefront/auth/login.
func (h *Handler) Login(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	var input LoginInput
	if err := c.BodyParser(&input); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	resp, err := h.svc.Login(c.Context(), tenantID, input)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

// Refresh handles POST /storefront/auth/refresh.
func (h *Handler) Refresh(c *fiber.Ctx) error {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := c.BodyParser(&body); err != nil || body.RefreshToken == "" {
		return errx.New("refresh_token is required", errx.TypeValidation)
	}

	resp, err := h.svc.RefreshToken(c.Context(), body.RefreshToken)
	if err != nil {
		return err
	}

	return c.JSON(resp)
}

// ============================================================================
// Protected handlers
// ============================================================================

// GetProfile handles GET /storefront/account/profile.
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	auth, ok := GetCustomerAuthContext(c)
	if !ok {
		return ErrTokenRequired()
	}

	cust, err := h.svc.GetProfile(c.Context(), auth.TenantID, auth.CustomerID)
	if err != nil {
		return err
	}

	return c.JSON(cust)
}

type updateProfileRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

// UpdateProfile handles PUT /storefront/account/profile.
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	auth, ok := GetCustomerAuthContext(c)
	if !ok {
		return ErrTokenRequired()
	}

	var req updateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	cust, err := h.svc.UpdateProfile(c.Context(), auth.TenantID, auth.CustomerID, req.Name, req.Phone)
	if err != nil {
		return err
	}

	return c.JSON(cust)
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePassword handles PUT /storefront/account/password.
func (h *Handler) ChangePassword(c *fiber.Ctx) error {
	auth, ok := GetCustomerAuthContext(c)
	if !ok {
		return ErrTokenRequired()
	}

	var req changePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.OldPassword == "" || req.NewPassword == "" {
		return errx.New("old_password and new_password are required", errx.TypeValidation)
	}

	if err := h.svc.ChangePassword(c.Context(), auth.TenantID, auth.CustomerID, req.OldPassword, req.NewPassword); err != nil {
		return err
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// ListOrders handles GET /storefront/account/orders.
func (h *Handler) ListOrders(c *fiber.Ctx) error {
	auth, ok := GetCustomerAuthContext(c)
	if !ok {
		return ErrTokenRequired()
	}

	pg := paginationFromQuery(c)

	result, err := h.orderSvc.ListByCustomer(c.Context(), auth.TenantID, auth.CustomerID, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// ListAddresses handles GET /storefront/account/addresses.
func (h *Handler) ListAddresses(c *fiber.Ctx) error {
	auth, ok := GetCustomerAuthContext(c)
	if !ok {
		return ErrTokenRequired()
	}

	cust, err := h.svc.GetProfile(c.Context(), auth.TenantID, auth.CustomerID)
	if err != nil {
		return err
	}

	return c.JSON(cust.Addresses)
}

type addAddressRequest struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	IsDefault  bool   `json:"is_default"`
}

// AddAddress handles POST /storefront/account/addresses.
func (h *Handler) AddAddress(c *fiber.Ctx) error {
	auth, ok := GetCustomerAuthContext(c)
	if !ok {
		return ErrTokenRequired()
	}

	var req addAddressRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	cust, err := h.svc.GetProfile(c.Context(), auth.TenantID, auth.CustomerID)
	if err != nil {
		return err
	}

	cust.AddAddress(customer.Address{
		Street:     req.Street,
		City:       req.City,
		State:      req.State,
		Country:    req.Country,
		PostalCode: req.PostalCode,
		IsDefault:  req.IsDefault,
	})

	if err := h.svc.customerSvc.Update(c.Context(), cust); err != nil {
		return errx.Wrap(err, "saving address", errx.TypeInternal)
	}

	return c.Status(fiber.StatusCreated).JSON(cust.Addresses)
}

// SetDefaultAddress handles PUT /storefront/account/addresses/:idx/default.
func (h *Handler) SetDefaultAddress(c *fiber.Ctx) error {
	auth, ok := GetCustomerAuthContext(c)
	if !ok {
		return ErrTokenRequired()
	}

	idx, err := strconv.Atoi(c.Params("idx"))
	if err != nil {
		return errx.New("invalid address index", errx.TypeValidation)
	}

	cust, err := h.svc.GetProfile(c.Context(), auth.TenantID, auth.CustomerID)
	if err != nil {
		return err
	}

	if !cust.SetDefaultAddress(idx) {
		return errx.New("address index out of range", errx.TypeValidation)
	}

	if err := h.svc.customerSvc.Update(c.Context(), cust); err != nil {
		return errx.Wrap(err, "updating default address", errx.TypeInternal)
	}

	return c.JSON(cust.Addresses)
}

// ============================================================================
// Helpers
// ============================================================================

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
