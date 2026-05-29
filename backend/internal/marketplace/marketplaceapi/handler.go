package marketplaceapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace/marketplacesrv"
)

// Handler exposes HTTP endpoints for the marketplace domain.
type Handler struct {
	svc *marketplacesrv.Service
}

// NewHandler creates a new marketplace API handler.
func NewHandler(svc *marketplacesrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all marketplace routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /marketplace/plugins", h.ListPlugins)
	mux.HandleFunc("GET /marketplace/plugins/{id}", h.GetPlugin)
	mux.HandleFunc("POST /marketplace/install", h.Install)
	mux.HandleFunc("POST /marketplace/uninstall", h.Uninstall)
	mux.HandleFunc("GET /marketplace/installed", h.ListInstalled)
	mux.HandleFunc("PUT /marketplace/plugins/{id}/settings", h.UpdateSettings)
}

// ListPlugins handles GET /marketplace/plugins.
func (h *Handler) ListPlugins(w http.ResponseWriter, r *http.Request) {
	pg := paginationFromQuery(r)

	result, err := h.svc.ListAvailable(r.Context(), pg)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// GetPlugin handles GET /marketplace/plugins/{id}.
func (h *Handler) GetPlugin(w http.ResponseWriter, r *http.Request) {
	id := kernel.PluginID(r.PathValue("id"))

	detail, err := h.svc.GetPlugin(r.Context(), id)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, detail)
}

// installRequest is the JSON body for install/uninstall.
type installRequest struct {
	PluginID string `json:"plugin_id"`
}

// Install handles POST /marketplace/install.
func (h *Handler) Install(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	var req installRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	inst, err := h.svc.Install(r.Context(), tenantID, kernel.PluginID(req.PluginID))
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, inst)
}

// Uninstall handles POST /marketplace/uninstall.
func (h *Handler) Uninstall(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	var req installRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.svc.Uninstall(r.Context(), tenantID, kernel.PluginID(req.PluginID)); err != nil {
		writeErrx(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListInstalled handles GET /marketplace/installed.
func (h *Handler) ListInstalled(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	installations, err := h.svc.ListInstalled(r.Context(), tenantID)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, installations)
}

// settingsRequest is the JSON body for updating plugin settings.
type settingsRequest struct {
	Settings map[string]any `json:"settings"`
}

// UpdateSettings handles PUT /marketplace/plugins/{id}/settings.
func (h *Handler) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	pluginID := kernel.PluginID(r.PathValue("id"))

	var req settingsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	inst, err := h.svc.UpdateSettings(r.Context(), tenantID, pluginID, req.Settings)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, inst)
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
