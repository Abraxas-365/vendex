package loyaltyapi

import (
	"strconv"

	"github.com/gofiber/fiber/v2"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/loyalty"
	"github.com/Abraxas-365/vendex/internal/loyalty/loyaltysrv"
)

// Handler exposes loyalty HTTP endpoints via Fiber v2.
type Handler struct {
	svc *loyaltysrv.Service
}

// New creates a new loyalty Handler.
func New(svc *loyaltysrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes wires all admin (protected) loyalty routes.
//
//	POST   /admin/loyalty/rewards                   — create reward
//	GET    /admin/loyalty/rewards                   — list rewards
//	PUT    /admin/loyalty/rewards/:id               — update reward
//	GET    /admin/loyalty/accounts                  — list all accounts
//	GET    /admin/loyalty/accounts/:id              — get account detail
//	POST   /admin/loyalty/accounts/:id/adjust       — manual point adjustment
//	GET    /admin/loyalty/accounts/:id/transactions — list account transactions
func (h *Handler) RegisterRoutes(router fiber.Router) {
	admin := router.Group("/admin/loyalty")

	// Rewards
	admin.Post("/rewards", h.createReward)
	admin.Get("/rewards", h.listRewards)
	admin.Put("/rewards/:id", h.updateReward)

	// Accounts (admin)
	admin.Get("/accounts", h.listAccounts)
	admin.Get("/accounts/:id", h.getAccount)
	admin.Post("/accounts/:id/adjust", h.adjustPoints)
	admin.Get("/accounts/:id/transactions", h.listTransactions)
}

// RegisterPublicRoutes wires all public (customer-facing) loyalty routes.
//
//	GET  /loyalty/account — get own account (customer identified by X-Customer-ID header)
//	POST /loyalty/redeem  — redeem a reward
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	router.Get("/loyalty/account", h.getOwnAccount)
	router.Post("/loyalty/redeem", h.redeemPoints)
}

// ============================================================================
// Admin handlers
// ============================================================================

type createRewardRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	PointsCost  int    `json:"points_cost"`
	RewardType  string `json:"reward_type"`
	ValueCents  int    `json:"value_cents"`
}

func (h *Handler) createReward(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRewardRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}
	if req.PointsCost <= 0 {
		return errx.New("points_cost must be greater than zero", errx.TypeValidation)
	}
	if req.RewardType == "" {
		return errx.New("reward_type is required", errx.TypeValidation)
	}

	reward, err := h.svc.CreateReward(c.Context(), authCtx.TenantID, loyalty.CreateRewardInput{
		Name:        req.Name,
		Description: req.Description,
		PointsCost:  req.PointsCost,
		RewardType:  req.RewardType,
		ValueCents:  req.ValueCents,
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(reward)
}

func (h *Handler) listRewards(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	p := paginationFromCtx(c)

	result, err := h.svc.ListRewards(c.Context(), authCtx.TenantID, p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

type updateRewardRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	PointsCost  *int    `json:"points_cost,omitempty"`
	RewardType  *string `json:"reward_type,omitempty"`
	ValueCents  *int    `json:"value_cents,omitempty"`
	Active      *bool   `json:"active,omitempty"`
}

func (h *Handler) updateReward(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.RewardID(c.Params("id"))

	var req updateRewardRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	reward, err := h.svc.UpdateReward(c.Context(), authCtx.TenantID, id, loyalty.UpdateRewardInput{
		Name:        req.Name,
		Description: req.Description,
		PointsCost:  req.PointsCost,
		RewardType:  req.RewardType,
		ValueCents:  req.ValueCents,
		Active:      req.Active,
	})
	if err != nil {
		return err
	}
	return c.JSON(reward)
}

func (h *Handler) listAccounts(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	p := paginationFromCtx(c)

	result, err := h.svc.ListAccounts(c.Context(), authCtx.TenantID, p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

func (h *Handler) getAccount(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.LoyaltyAccountID(c.Params("id"))

	account, err := h.svc.GetAccount(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(account)
}

type adjustPointsRequest struct {
	Points    int    `json:"points"`
	Reference string `json:"reference"`
	Note      string `json:"note"`
}

func (h *Handler) adjustPoints(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.LoyaltyAccountID(c.Params("id"))

	var req adjustPointsRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Points == 0 {
		return errx.New("points must not be zero", errx.TypeValidation)
	}

	account, err := h.svc.AdjustPoints(c.Context(), authCtx.TenantID, id, loyalty.AdjustPointsInput{
		Points:    req.Points,
		Reference: req.Reference,
		Note:      req.Note,
	})
	if err != nil {
		return err
	}
	return c.JSON(account)
}

func (h *Handler) listTransactions(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.LoyaltyAccountID(c.Params("id"))
	p := paginationFromCtx(c)

	result, err := h.svc.ListTransactions(c.Context(), authCtx.TenantID, id, p)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// ============================================================================
// Public (customer-facing) handlers
// ============================================================================

func (h *Handler) getOwnAccount(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID required", errx.TypeAuthorization)
	}

	customerID := kernel.CustomerID(c.Get("X-Customer-ID"))
	if customerID == "" {
		return errx.New("customer ID required", errx.TypeAuthorization)
	}

	account, err := h.svc.GetOrCreateAccount(c.Context(), tenantID, customerID)
	if err != nil {
		return err
	}
	return c.JSON(account)
}

type redeemRequest struct {
	RewardID string `json:"reward_id"`
	Note     string `json:"note"`
}

func (h *Handler) redeemPoints(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID required", errx.TypeAuthorization)
	}

	customerID := kernel.CustomerID(c.Get("X-Customer-ID"))
	if customerID == "" {
		return errx.New("customer ID required", errx.TypeAuthorization)
	}

	var req redeemRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.RewardID == "" {
		return errx.New("reward_id is required", errx.TypeValidation)
	}

	account, err := h.svc.RedeemPoints(c.Context(), tenantID, loyalty.RedeemPointsInput{
		CustomerID: customerID,
		RewardID:   kernel.RewardID(req.RewardID),
		Note:       req.Note,
	})
	if err != nil {
		return err
	}
	return c.JSON(account)
}

// ============================================================================
// Helpers
// ============================================================================

func paginationFromCtx(c *fiber.Ctx) kernel.PaginationOptions {
	page, _ := strconv.Atoi(c.Query("page"))
	size, _ := strconv.Atoi(c.Query("page_size"))
	return kernel.NewPaginationOptions(page, size)
}
