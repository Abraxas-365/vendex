package searchapi

import (
	"strconv"
	"strings"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/search"
	"github.com/Abraxas-365/hada-commerce/internal/search/searchsrv"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the search domain.
type Handler struct {
	svc *searchsrv.Service
}

// NewHandler creates a new search HTTP handler.
func NewHandler(svc *searchsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterPublicRoutes registers unauthenticated search routes.
// Tenant is identified via the X-Tenant-ID header.
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/search")
	g.Get("/", h.Search)
	g.Get("/suggest", h.Suggest)
}

// Search handles GET /search.
// Query params: q, category_id, tags (comma-separated), status, min_price, max_price,
// sort_by, page, page_size.
func (h *Handler) Search(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	q := buildSearchQuery(c)
	result, err := h.svc.Search(c.Context(), tenantID, q)
	if err != nil {
		return err
	}
	return c.JSON(result)
}

// Suggest handles GET /search/suggest.
// Query params: q, limit.
func (h *Handler) Suggest(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("X-Tenant-ID header is required", errx.TypeValidation)
	}

	prefix := c.Query("q")
	if prefix == "" {
		return errx.New("query parameter 'q' is required", errx.TypeValidation)
	}

	limit, _ := strconv.Atoi(c.Query("limit"))

	suggestions, err := h.svc.Suggest(c.Context(), tenantID, prefix, limit)
	if err != nil {
		return err
	}
	return c.JSON(fiber.Map{"suggestions": suggestions})
}

// buildSearchQuery extracts a SearchQuery from the Fiber context's query params.
func buildSearchQuery(c *fiber.Ctx) search.SearchQuery {
	q := search.SearchQuery{
		Query:      c.Query("q"),
		CategoryID: c.Query("category_id"),
		Status:     c.Query("status"),
		SortBy:     c.Query("sort_by"),
	}

	// Default status to "active" for public storefront searches.
	if q.Status == "" {
		q.Status = "active"
	}

	// Tags: comma-separated.
	if tagsParam := c.Query("tags"); tagsParam != "" {
		for _, t := range strings.Split(tagsParam, ",") {
			t = strings.TrimSpace(t)
			if t != "" {
				q.Tags = append(q.Tags, t)
			}
		}
	}

	// Price range.
	if minStr := c.Query("min_price"); minStr != "" {
		if v, err := strconv.ParseInt(minStr, 10, 64); err == nil {
			q.MinPrice = &v
		}
	}
	if maxStr := c.Query("max_price"); maxStr != "" {
		if v, err := strconv.ParseInt(maxStr, 10, 64); err == nil {
			q.MaxPrice = &v
		}
	}

	// Pagination.
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))
	q.Page = page
	q.PageSize = pageSize

	return q
}
