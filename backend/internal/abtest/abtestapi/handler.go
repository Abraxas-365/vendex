package abtestapi

import (
	"strconv"

	"github.com/Abraxas-365/vendex/internal/abtest"
	"github.com/Abraxas-365/vendex/internal/abtest/abtestsrv"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the A/B testing domain.
type Handler struct {
	svc *abtestsrv.Service
}

// NewHandler creates a new A/B testing API handler.
func NewHandler(svc *abtestsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected (admin) A/B test routes.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/experiments")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get("/:id", h.GetByID)
	g.Put("/:id", h.Update)
	g.Delete("/:id", h.Delete)
	g.Put("/:id/start", h.Start)
	g.Put("/:id/pause", h.Pause)
	g.Put("/:id/complete", h.Complete)
	g.Post("/:id/variants", h.AddVariant)
	g.Delete("/:id/variants/:variantId", h.RemoveVariant)
	g.Get("/:id/results", h.GetResults)
}

// RegisterPublicRoutes registers unauthenticated A/B testing routes.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/experiments")
	g.Post("/assign", h.Assign)
	g.Post("/convert", h.Convert)
}

// ─── Admin handlers ───────────────────────────────────────────────────────────

type createRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	Type           string `json:"type"`
	TrafficPercent int    `json:"traffic_percent"`
}

// Create handles POST /experiments.
func (h *Handler) Create(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}

	expType := abtest.ExperimentType(req.Type)
	if expType == "" {
		expType = abtest.TypePage
	}
	trafficPercent := req.TrafficPercent
	if trafficPercent == 0 {
		trafficPercent = 100
	}

	e, err := h.svc.CreateExperiment(c.Context(), authCtx.TenantID, abtest.CreateExperimentInput{
		Name:           req.Name,
		Description:    req.Description,
		Type:           expType,
		TrafficPercent: trafficPercent,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(e)
}

// List handles GET /experiments.
func (h *Handler) List(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	status := c.Query("status", "")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "20"))

	result, err := h.svc.ListExperiments(c.Context(), authCtx.TenantID, status, page, pageSize)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetByID handles GET /experiments/:id.
func (h *Handler) GetByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ExperimentID(c.Params("id"))

	e, err := h.svc.GetExperiment(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(e)
}

type updateRequest struct {
	Name           *string `json:"name"`
	Description    *string `json:"description"`
	TrafficPercent *int    `json:"traffic_percent"`
}

// Update handles PUT /experiments/:id.
func (h *Handler) Update(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ExperimentID(c.Params("id"))

	var req updateRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	e, err := h.svc.UpdateExperiment(c.Context(), authCtx.TenantID, id, abtest.UpdateExperimentInput{
		Name:           req.Name,
		Description:    req.Description,
		TrafficPercent: req.TrafficPercent,
	})
	if err != nil {
		return err
	}

	return c.JSON(e)
}

// Delete handles DELETE /experiments/:id.
func (h *Handler) Delete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ExperimentID(c.Params("id"))

	if err := h.svc.DeleteExperiment(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// Start handles PUT /experiments/:id/start.
func (h *Handler) Start(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ExperimentID(c.Params("id"))

	e, err := h.svc.StartExperiment(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(e)
}

// Pause handles PUT /experiments/:id/pause.
func (h *Handler) Pause(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ExperimentID(c.Params("id"))

	e, err := h.svc.PauseExperiment(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(e)
}

type completeRequest struct {
	WinnerVariantID string `json:"winner_variant_id"`
}

// Complete handles PUT /experiments/:id/complete.
func (h *Handler) Complete(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ExperimentID(c.Params("id"))

	var req completeRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.WinnerVariantID == "" {
		return errx.New("winner_variant_id is required", errx.TypeValidation)
	}

	e, err := h.svc.CompleteExperiment(c.Context(), authCtx.TenantID, id, kernel.ExperimentVariantID(req.WinnerVariantID))
	if err != nil {
		return err
	}

	return c.JSON(e)
}

type addVariantRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Weight      int                    `json:"weight"`
	IsControl   bool                   `json:"is_control"`
	Config      map[string]interface{} `json:"config"`
}

// AddVariant handles POST /experiments/:id/variants.
func (h *Handler) AddVariant(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ExperimentID(c.Params("id"))

	var req addVariantRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.Name == "" {
		return errx.New("name is required", errx.TypeValidation)
	}

	v, err := h.svc.AddVariant(c.Context(), authCtx.TenantID, id, abtest.CreateVariantInput{
		Name:        req.Name,
		Description: req.Description,
		Weight:      req.Weight,
		IsControl:   req.IsControl,
		Config:      req.Config,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(v)
}

// RemoveVariant handles DELETE /experiments/:id/variants/:variantId.
func (h *Handler) RemoveVariant(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ExperimentID(c.Params("id"))
	variantID := kernel.ExperimentVariantID(c.Params("variantId"))

	if err := h.svc.RemoveVariant(c.Context(), authCtx.TenantID, id, variantID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GetResults handles GET /experiments/:id/results.
func (h *Handler) GetResults(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.ExperimentID(c.Params("id"))

	results, err := h.svc.GetResults(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(results)
}

// ─── Public handlers ─────────────────────────────────────────────────────────

type assignRequest struct {
	ExperimentID string `json:"experiment_id"`
	VisitorID    string `json:"visitor_id"`
}

// Assign handles POST /experiments/assign (public).
func (h *Handler) Assign(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID.IsEmpty() {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	var req assignRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.ExperimentID == "" || req.VisitorID == "" {
		return errx.New("experiment_id and visitor_id are required", errx.TypeValidation)
	}

	assignment, err := h.svc.AssignVisitor(c.Context(), tenantID,
		kernel.ExperimentID(req.ExperimentID), req.VisitorID)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(assignment)
}

type convertRequest struct {
	ExperimentID string `json:"experiment_id"`
	VisitorID    string `json:"visitor_id"`
	RevenueCents int64  `json:"revenue_cents"`
}

// Convert handles POST /experiments/convert (public).
func (h *Handler) Convert(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID.IsEmpty() {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	var req convertRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}
	if req.ExperimentID == "" || req.VisitorID == "" {
		return errx.New("experiment_id and visitor_id are required", errx.TypeValidation)
	}

	if err := h.svc.RecordConversion(c.Context(), tenantID,
		kernel.ExperimentID(req.ExperimentID), req.VisitorID, req.RevenueCents); err != nil {
		return err
	}

	return c.JSON(fiber.Map{"status": "ok"})
}
