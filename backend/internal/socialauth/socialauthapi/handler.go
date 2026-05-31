package socialauthapi

import (
	"strconv"

	customerauth "github.com/Abraxas-365/vendex/internal/customer/auth"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/socialauth"
	"github.com/Abraxas-365/vendex/internal/socialauth/socialauthsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the social auth domain.
type Handler struct {
	svc *socialauthsrv.Service
}

// NewHandler creates a new social auth API handler.
func NewHandler(svc *socialauthsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// ============================================================================
// Route registration
// ============================================================================

// RegisterPublicRoutes registers public (storefront) social auth routes.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/auth/social")
	g.Post("/link", h.LinkAccount)
	g.Get("/providers", h.ListProviders)
}

// RegisterCustomerRoutes registers customer-authenticated social account routes.
// The router must already have CustomerMiddleware applied.
func (h *Handler) RegisterCustomerRoutes(router fiber.Router) {
	g := router.Group("/account/social-accounts")
	g.Get("/", h.ListMyAccounts)
	g.Delete("/:id", h.UnlinkMyAccount)
}

// RegisterRoutes registers admin-protected social auth routes.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/social-accounts")
	g.Get("/", h.AdminList)
	g.Get("/customer/:customerId", h.AdminListByCustomer)
}

// ============================================================================
// Public handlers
// ============================================================================

type linkRequest struct {
	CustomerID     string `json:"customer_id"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url"`
	AccessToken    string `json:"access_token"`
	RefreshToken   string `json:"refresh_token"`
}

// LinkAccount handles POST /auth/social/link.
// Links an OAuth social account to a customer. Tenant is provided via X-Tenant-ID header.
func (h *Handler) LinkAccount(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	var req linkRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	if req.CustomerID == "" {
		return errx.New("customer_id is required", errx.TypeValidation)
	}
	if req.Provider == "" {
		return errx.New("provider is required", errx.TypeValidation)
	}
	if req.ProviderUserID == "" {
		return errx.New("provider_user_id is required", errx.TypeValidation)
	}

	input := socialauth.LinkInput{
		CustomerID:     kernel.CustomerID(req.CustomerID),
		Provider:       req.Provider,
		ProviderUserID: req.ProviderUserID,
		Email:          req.Email,
		Name:           req.Name,
		AvatarURL:      req.AvatarURL,
		AccessToken:    req.AccessToken,
		RefreshToken:   req.RefreshToken,
	}

	sa, err := h.svc.LinkAccount(c.Context(), tenantID, input)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(sa)
}

// ListProviders handles GET /auth/social/providers.
// Returns the list of supported OAuth providers.
func (h *Handler) ListProviders(c *fiber.Ctx) error {
	providers := h.svc.GetProviderConfig()
	return c.JSON(fiber.Map{"providers": providers})
}

// ============================================================================
// Customer-protected handlers
// ============================================================================

// ListMyAccounts handles GET /account/social-accounts.
func (h *Handler) ListMyAccounts(c *fiber.Ctx) error {
	auth, ok := customerauth.GetCustomerAuthContext(c)
	if !ok {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	accounts, err := h.svc.ListByCustomer(c.Context(), auth.TenantID, auth.CustomerID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"accounts": accounts})
}

// UnlinkMyAccount handles DELETE /account/social-accounts/:id.
func (h *Handler) UnlinkMyAccount(c *fiber.Ctx) error {
	auth, ok := customerauth.GetCustomerAuthContext(c)
	if !ok {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	id := kernel.SocialAccountID(c.Params("id"))
	if id.IsEmpty() {
		return errx.New("id is required", errx.TypeValidation)
	}

	// Verify the account belongs to this customer before unlinking.
	accounts, err := h.svc.ListByCustomer(c.Context(), auth.TenantID, auth.CustomerID)
	if err != nil {
		return err
	}
	found := false
	for _, a := range accounts {
		if a.ID == id {
			found = true
			break
		}
	}
	if !found {
		return socialauth.ErrNotFound()
	}

	if err := h.svc.UnlinkAccount(c.Context(), auth.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ============================================================================
// Admin-protected handlers
// ============================================================================

// AdminList handles GET /social-accounts (admin, paginated).
func (h *Handler) AdminList(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	pg := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, pg)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// AdminListByCustomer handles GET /social-accounts/customer/:customerId (admin).
func (h *Handler) AdminListByCustomer(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	customerID := kernel.CustomerID(c.Params("customerId"))

	if customerID.IsEmpty() {
		return errx.New("customerId is required", errx.TypeValidation)
	}

	accounts, err := h.svc.ListByCustomer(c.Context(), authCtx.TenantID, customerID)
	if err != nil {
		return err
	}

	return c.JSON(fiber.Map{"accounts": accounts})
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
