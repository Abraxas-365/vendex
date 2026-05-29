package catalogapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/catalog/catalogsrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

// Handler exposes HTTP endpoints for the catalog domain.
type Handler struct {
	svc *catalogsrv.Service
}

// NewHandler creates a new catalog API handler.
func NewHandler(svc *catalogsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all catalog routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Categories
	mux.HandleFunc("POST /categories", h.CreateCategory)
	mux.HandleFunc("GET /categories/{id}", h.GetCategory)
	mux.HandleFunc("GET /categories", h.ListCategories)
	mux.HandleFunc("DELETE /categories/{id}", h.DeleteCategory)

	// Collections
	mux.HandleFunc("POST /collections", h.CreateCollection)
	mux.HandleFunc("GET /collections/{id}", h.GetCollection)
	mux.HandleFunc("GET /collections", h.ListCollections)
	mux.HandleFunc("DELETE /collections/{id}", h.DeleteCollection)
}

// --- Category handlers ---

type createCategoryRequest struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	ParentID    *string `json:"parent_id,omitempty"`
	Description string  `json:"description"`
}

// CreateCategory handles POST /categories.
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	var req createCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var parentID *kernel.CategoryID
	if req.ParentID != nil {
		pid := kernel.CategoryID(*req.ParentID)
		parentID = &pid
	}

	c, err := h.svc.CreateCategory(r.Context(), tenantID, catalogsrv.CreateCategoryInput{
		Name:        req.Name,
		Slug:        req.Slug,
		ParentID:    parentID,
		Description: req.Description,
	})
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, c)
}

// GetCategory handles GET /categories/{id}.
func (h *Handler) GetCategory(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.CategoryID(r.PathValue("id"))

	c, err := h.svc.GetCategoryByID(r.Context(), tenantID, id)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, c)
}

// ListCategories handles GET /categories.
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	pg := paginationFromQuery(r)

	result, err := h.svc.ListCategories(r.Context(), tenantID, pg)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// DeleteCategory handles DELETE /categories/{id}.
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.CategoryID(r.PathValue("id"))

	if err := h.svc.DeleteCategory(r.Context(), tenantID, id); err != nil {
		writeErrx(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- Collection handlers ---

type createCollectionRequest struct {
	Name        string   `json:"name"`
	Slug        string   `json:"slug"`
	Description string   `json:"description"`
	ProductIDs  []string `json:"product_ids"`
	IsAutomatic bool     `json:"is_automatic"`
	Rules       map[string]any `json:"rules"`
}

// CreateCollection handles POST /collections.
func (h *Handler) CreateCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	var req createCollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	productIDs := make([]kernel.ProductID, len(req.ProductIDs))
	for i, id := range req.ProductIDs {
		productIDs[i] = kernel.ProductID(id)
	}

	c, err := h.svc.CreateCollection(r.Context(), tenantID, catalogsrv.CreateCollectionInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		ProductIDs:  productIDs,
		IsAutomatic: req.IsAutomatic,
		Rules:       req.Rules,
	})
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, c)
}

// GetCollection handles GET /collections/{id}.
func (h *Handler) GetCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.CollectionID(r.PathValue("id"))

	c, err := h.svc.GetCollectionByID(r.Context(), tenantID, id)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, c)
}

// ListCollections handles GET /collections.
func (h *Handler) ListCollections(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	pg := paginationFromQuery(r)

	result, err := h.svc.ListCollections(r.Context(), tenantID, pg)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// DeleteCollection handles DELETE /collections/{id}.
func (h *Handler) DeleteCollection(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.CollectionID(r.PathValue("id"))

	if err := h.svc.DeleteCollection(r.Context(), tenantID, id); err != nil {
		writeErrx(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// --- helpers ---

func tenantFromRequest(r *http.Request) kernel.TenantID {
	return kernel.TenantID(r.Header.Get("X-Tenant-ID"))
}

func paginationFromQuery(r *http.Request) kernel.Pagination {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	return kernel.NewPagination(page, pageSize)
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func writeErrx(w http.ResponseWriter, err error) {
	status := errx.HTTPStatus(err)
	msg := errx.Message(err)
	writeJSON(w, status, map[string]string{"error": msg, "code": errx.Code(err)})
}
