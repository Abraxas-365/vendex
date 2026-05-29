package mediaapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
	"github.com/Abraxas-365/hada-commerce/internal/media"
	"github.com/Abraxas-365/hada-commerce/internal/media/mediasrv"
)

const (
	// maxUploadSize is the maximum file size accepted (32 MiB).
	maxUploadSize = 32 << 20
)

// Handler exposes media HTTP endpoints.
type Handler struct {
	svc *mediasrv.Service
}

// New creates a new media Handler.
func New(svc *mediasrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes wires all media routes onto the provided ServeMux.
//
//	POST /admin/media          — multipart file upload
//	GET  /admin/media          — list media
//	GET  /admin/media/:id      — get by ID
//	DELETE /admin/media/:id    — delete
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /admin/media", h.upload)
	mux.HandleFunc("GET /admin/media", h.list)
	mux.HandleFunc("GET /admin/media/{id}", h.getByID)
	mux.HandleFunc("DELETE /admin/media/{id}", h.delete)
}

type contextKey string

const contextKeyTenantID contextKey = "tenant_id"

func tenantFromContext(r *http.Request) kernel.TenantID {
	if v, ok := r.Context().Value(contextKeyTenantID).(string); ok {
		return kernel.TenantID(v)
	}
	return ""
}

// upload handles multipart/form-data file uploads.
// Form fields:
//
//	file       — the file (required)
//	alt        — alt text for images (optional)
//	uploaded_by — identifier of the uploader (optional)
func (h *Handler) upload(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		writeError(w, media.ErrFileTooLarge)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		writeError(w, media.ErrInvalidFile)
		return
	}
	defer file.Close()

	alt := r.FormValue("alt")
	uploadedBy := r.FormValue("uploaded_by")
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	m, err := h.svc.Upload(r.Context(), mediasrv.UploadInput{
		TenantID:    tenantID,
		Filename:    header.Filename,
		ContentType: contentType,
		Size:        header.Size,
		Alt:         alt,
		UploadedBy:  uploadedBy,
		Data:        file,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, m)
}

func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	p := paginationFromRequest(r)

	result, err := h.svc.List(r.Context(), tenantID, p)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) getByID(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	id := kernel.MediaID(r.PathValue("id"))

	m, err := h.svc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, m)
}

func (h *Handler) delete(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	id := kernel.MediaID(r.PathValue("id"))

	if err := h.svc.Delete(r.Context(), tenantID, id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
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
