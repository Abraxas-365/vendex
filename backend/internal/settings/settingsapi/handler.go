package settingsapi

import (
	"encoding/json"
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
	"github.com/Abraxas-365/hada-commerce/internal/settings"
	"github.com/Abraxas-365/hada-commerce/internal/settings/settingssrv"
)

// Handler exposes HTTP endpoints for the settings domain.
type Handler struct {
	svc *settingssrv.Service
}

// NewHandler creates a new settings API handler.
func NewHandler(svc *settingssrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all settings routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /settings", h.Get)
	mux.HandleFunc("PUT /settings", h.Update)
}

// updateRequest is the JSON body for updating store settings.
type updateRequest struct {
	StoreName      string                  `json:"store_name"`
	StoreEmail     string                  `json:"store_email"`
	StorePhone     string                  `json:"store_phone"`
	Currency       string                  `json:"currency"`
	Timezone       string                  `json:"timezone"`
	Address        settings.StoreAddress   `json:"address"`
	LogoURL        string                  `json:"logo_url"`
	FaviconURL     string                  `json:"favicon_url"`
	SocialLinks    settings.SocialLinks    `json:"social_links"`
	CheckoutConfig settings.CheckoutConfig `json:"checkout_config"`
}

// Get handles GET /settings.
// Returns the current settings for the tenant, creating defaults if none exist.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	ss, err := h.svc.Get(r.Context(), tenantID)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, ss)
}

// Update handles PUT /settings.
// Upserts settings for the tenant from the request body.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	var req updateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	ss, err := h.svc.Update(r.Context(), tenantID, settingssrv.UpdateInput{
		StoreName:      req.StoreName,
		StoreEmail:     req.StoreEmail,
		StorePhone:     req.StorePhone,
		Currency:       req.Currency,
		Timezone:       req.Timezone,
		Address:        req.Address,
		LogoURL:        req.LogoURL,
		FaviconURL:     req.FaviconURL,
		SocialLinks:    req.SocialLinks,
		CheckoutConfig: req.CheckoutConfig,
	})
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, ss)
}

// --- helpers ---

func tenantFromRequest(r *http.Request) kernel.TenantID {
	return kernel.TenantID(r.Header.Get("X-Tenant-ID"))
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
