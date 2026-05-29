package storefrontapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
	"github.com/Abraxas-365/hada-commerce/internal/storefront"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/storefrontsrv"
)

// Handler exposes storefront HTTP endpoints.
type Handler struct {
	svc *storefrontsrv.Service
}

// New creates a new storefront Handler.
func New(svc *storefrontsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes wires all routes onto the provided ServeMux.
// Public:
//
//	GET /pages/:slug           — serve published page HTML
//
// Admin (prefix /admin/pages):
//
//	GET    /admin/pages           — list all pages
//	POST   /admin/pages           — create page
//	GET    /admin/pages/:id       — get page by ID
//	PUT    /admin/pages/:id       — update page content
//	POST   /admin/pages/:id/publish   — publish
//	POST   /admin/pages/:id/unpublish — unpublish
//	POST   /admin/pages/:id/archive   — archive
//	GET    /admin/pages/:id/versions  — list versions
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	// Public
	mux.HandleFunc("GET /pages/{slug}", h.getPublishedPage)

	// Admin CRUD
	mux.HandleFunc("GET /admin/pages", h.listPages)
	mux.HandleFunc("POST /admin/pages", h.createPage)
	mux.HandleFunc("GET /admin/pages/{id}", h.getPage)
	mux.HandleFunc("PUT /admin/pages/{id}", h.updatePage)
	mux.HandleFunc("POST /admin/pages/{id}/publish", h.publish)
	mux.HandleFunc("POST /admin/pages/{id}/unpublish", h.unpublish)
	mux.HandleFunc("POST /admin/pages/{id}/archive", h.archive)
	mux.HandleFunc("GET /admin/pages/{id}/versions", h.listVersions)
}

// tenantFromContext extracts the TenantID stored in the request context.
// Convention: middleware sets context value with key "tenant_id".
func tenantFromContext(r *http.Request) kernel.TenantID {
	if v, ok := r.Context().Value(contextKeyTenantID).(string); ok {
		return kernel.TenantID(v)
	}
	return ""
}

type contextKey string

const contextKeyTenantID contextKey = "tenant_id"

// --- Public endpoints ---

func (h *Handler) getPublishedPage(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	slug := r.PathValue("slug")

	page, err := h.svc.GetPublishedPage(r.Context(), tenantID, slug)
	if err != nil {
		writeError(w, err)
		return
	}

	// Serve raw HTML for published pages.
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(page.HTML))
}

// --- Admin endpoints ---

func (h *Handler) listPages(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	p := paginationFromRequest(r)

	var status *storefront.PageStatus
	if s := r.URL.Query().Get("status"); s != "" {
		ps := storefront.PageStatus(s)
		status = &ps
	}

	result, err := h.svc.ListPages(r.Context(), tenantID, status, p)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

type createPageRequest struct {
	Slug      string              `json:"slug"`
	Title     string              `json:"title"`
	HTML      string              `json:"html"`
	CSS       string              `json:"css"`
	Meta      storefront.PageMeta `json:"meta"`
	CreatedBy string              `json:"created_by"`
	ByAgent   bool                `json:"by_agent"`
}

func (h *Handler) createPage(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)

	var req createPageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errx.New("INVALID_REQUEST", "invalid request body", http.StatusBadRequest))
		return
	}

	page, err := h.svc.CreatePage(r.Context(), storefrontsrv.CreatePageInput{
		TenantID:  tenantID,
		Slug:      req.Slug,
		Title:     req.Title,
		HTML:      req.HTML,
		CSS:       req.CSS,
		Meta:      req.Meta,
		CreatedBy: req.CreatedBy,
		ByAgent:   req.ByAgent,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, page)
}

func (h *Handler) getPage(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	id := kernel.PageID(r.PathValue("id"))

	page, err := h.svc.GetPage(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

type updatePageRequest struct {
	Title    *string              `json:"title,omitempty"`
	HTML     *string              `json:"html,omitempty"`
	CSS      *string              `json:"css,omitempty"`
	Meta     *storefront.PageMeta `json:"meta,omitempty"`
	EditedBy string               `json:"edited_by"`
	Comment  string               `json:"comment"`
}

func (h *Handler) updatePage(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	id := kernel.PageID(r.PathValue("id"))

	var req updatePageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errx.New("INVALID_REQUEST", "invalid request body", http.StatusBadRequest))
		return
	}

	page, err := h.svc.UpdatePage(r.Context(), storefrontsrv.UpdatePageInput{
		TenantID: tenantID,
		ID:       id,
		Title:    req.Title,
		HTML:     req.HTML,
		CSS:      req.CSS,
		Meta:     req.Meta,
		EditedBy: req.EditedBy,
		Comment:  req.Comment,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (h *Handler) publish(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	id := kernel.PageID(r.PathValue("id"))

	page, err := h.svc.Publish(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (h *Handler) unpublish(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	id := kernel.PageID(r.PathValue("id"))

	page, err := h.svc.Unpublish(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (h *Handler) archive(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	id := kernel.PageID(r.PathValue("id"))

	page, err := h.svc.Archive(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, page)
}

func (h *Handler) listVersions(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	id := kernel.PageID(r.PathValue("id"))

	versions, err := h.svc.ListVersions(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, versions)
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, err error) {
	status := errx.HTTPStatus(err)
	body := map[string]string{
		"code":    errx.Code(err),
		"message": errx.Message(err),
	}
	writeJSON(w, status, body)
}

func paginationFromRequest(r *http.Request) kernel.Pagination {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	size, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	return kernel.NewPagination(page, size)
}
