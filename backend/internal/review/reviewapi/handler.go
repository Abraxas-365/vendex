package reviewapi

import (
	"strconv"

	customerauth "github.com/Abraxas-365/vendex/internal/customer/auth"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/review"
	"github.com/Abraxas-365/vendex/internal/review/reviewsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the review domain.
type Handler struct {
	svc *reviewsrv.Service
}

// NewHandler creates a new review API handler.
func NewHandler(svc *reviewsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers admin-authenticated review routes.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/reviews")
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id/approve", h.Approve)
	g.Put("/:id/reject", h.Reject)
	g.Post("/:id/respond", h.RespondAsAdmin)
}

// RegisterPublicRoutes registers unauthenticated, read-only review routes.
// Tenant is identified via the X-Tenant-ID header.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/products/:productId/reviews")
	g.Get("/", h.ListByProduct)
	g.Get("/stats", h.GetStats)
}

// RegisterCustomerRoutes registers customer-authenticated review routes.
// The router must already have CustomerMiddleware applied.
func (h *Handler) RegisterCustomerRoutes(router fiber.Router) {
	router.Post("/products/:productId/reviews", h.Create)
}

// ─── Request types ─────────────────────────────────────────────────────────────

type createReviewRequest struct {
	Rating int      `json:"rating"`
	Title  string   `json:"title"`
	Body   string   `json:"body"`
	Images []string `json:"images"`
}

type adminRespondRequest struct {
	Response string `json:"response"`
}

// ─── Admin handlers ─────────────────────────────────────────────────────────────

// List handles GET /reviews (admin).
// Accepts optional ?status= and pagination params.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	status := c.Query("status")
	pg := paginationFromQuery(c)

	result, err := h.svc.List(c.Context(), authCtx.TenantID, status, pg)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// GetByID handles GET /reviews/:id (admin).
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReviewID(c.Params("id"))

	rv, err := h.svc.GetByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(rv)
}

// Approve handles PUT /reviews/:id/approve (admin).
func (h *Handler) Approve(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReviewID(c.Params("id"))

	rv, err := h.svc.Approve(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(rv)
}

// Reject handles PUT /reviews/:id/reject (admin).
func (h *Handler) Reject(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReviewID(c.Params("id"))

	rv, err := h.svc.Reject(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}
	return c.JSON(rv)
}

// RespondAsAdmin handles POST /reviews/:id/respond (admin).
func (h *Handler) RespondAsAdmin(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ReviewID(c.Params("id"))

	var req adminRespondRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Response == "" {
		return errx.New("response is required", errx.TypeValidation)
	}

	rv, err := h.svc.RespondAsAdmin(c.Context(), authCtx.TenantID, id, req.Response)
	if err != nil {
		return err
	}
	return c.JSON(rv)
}

// ─── Public / Customer handlers ─────────────────────────────────────────────────

// ListByProduct handles GET /products/:productId/reviews (public).
// Returns only approved reviews unless admin header is set.
func (h *Handler) ListByProduct(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	productID := kernel.ProductID(c.Params("productId"))
	if productID.IsEmpty() {
		return errx.New("productId is required", errx.TypeValidation)
	}

	// Default to approved reviews for public access.
	status := c.Query("status", string(review.StatusApproved))
	pg := paginationFromQuery(c)

	result, err := h.svc.ListByProduct(c.Context(), tenantID, productID, status, pg)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// GetStats handles GET /products/:productId/reviews/stats (public).
func (h *Handler) GetStats(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	productID := kernel.ProductID(c.Params("productId"))
	if productID.IsEmpty() {
		return errx.New("productId is required", errx.TypeValidation)
	}

	stats, err := h.svc.GetStats(c.Context(), tenantID, productID)
	if err != nil {
		return err
	}
	return c.JSON(stats)
}

// Create handles POST /products/:productId/reviews (customer-authenticated).
func (h *Handler) Create(c *fiber.Ctx) error {
	auth, ok := customerauth.GetCustomerAuthContext(c)
	if !ok {
		return errx.New("unauthorized", errx.TypeAuthorization)
	}

	tenantID := kernel.TenantID(auth.TenantID)
	productID := kernel.ProductID(c.Params("productId"))
	if productID.IsEmpty() {
		return errx.New("productId is required", errx.TypeValidation)
	}

	var req createReviewRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Rating < 1 || req.Rating > 5 {
		return errx.New("rating must be between 1 and 5", errx.TypeValidation)
	}

	rv, err := h.svc.Create(c.Context(), tenantID, review.CreateReviewInput{
		ProductID:  productID,
		CustomerID: kernel.CustomerID(auth.CustomerID),
		Rating:     req.Rating,
		Title:      req.Title,
		Body:       req.Body,
		Images:     req.Images,
	})
	if err != nil {
		return err
	}
	return c.Status(fiber.StatusCreated).JSON(rv)
}

// ─── helpers ───────────────────────────────────────────────────────────────────

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
