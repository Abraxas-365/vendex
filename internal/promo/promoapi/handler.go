package promoapi

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
	"github.com/Abraxas-365/hada-commerce/internal/promo"
	"github.com/Abraxas-365/hada-commerce/internal/promo/promosrv"
)

// Handler exposes promo HTTP endpoints.
type Handler struct {
	svc *promosrv.Service
}

// New creates a new promo Handler.
func New(svc *promosrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes wires all promo routes onto the provided ServeMux.
//
//	POST   /admin/promos               — create promo
//	GET    /admin/promos               — list promos
//	GET    /admin/promos/:id           — get promo
//	POST   /admin/promos/:id/deactivate — deactivate
//	POST   /promos/validate            — validate a code for an order total (public)
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /admin/promos", h.create)
	mux.HandleFunc("GET /admin/promos", h.list)
	mux.HandleFunc("GET /admin/promos/{id}", h.getByID)
	mux.HandleFunc("POST /admin/promos/{id}/deactivate", h.deactivate)
	mux.HandleFunc("POST /promos/validate", h.validate)
}

type contextKey string

const contextKeyTenantID contextKey = "tenant_id"

func tenantFromContext(r *http.Request) kernel.TenantID {
	if v, ok := r.Context().Value(contextKeyTenantID).(string); ok {
		return kernel.TenantID(v)
	}
	return ""
}

// --- handlers ---

type createPromoRequest struct {
	Code           string     `json:"code"`
	Type           string     `json:"type"`
	Value          int64      `json:"value"`
	MinOrderAmount *int64     `json:"min_order_amount,omitempty"`
	MaxUses        *int       `json:"max_uses,omitempty"`
	StartsAt       *time.Time `json:"starts_at,omitempty"`
	EndsAt         *time.Time `json:"ends_at,omitempty"`
}

func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)

	var req createPromoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errx.New("INVALID_REQUEST", "invalid request body", http.StatusBadRequest))
		return
	}

	p, err := h.svc.Create(r.Context(), promosrv.CreateInput{
		TenantID:       tenantID,
		Code:           req.Code,
		Type:           promo.PromoType(req.Type),
		Value:          req.Value,
		MinOrderAmount: req.MinOrderAmount,
		MaxUses:        req.MaxUses,
		StartsAt:       req.StartsAt,
		EndsAt:         req.EndsAt,
	})
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, p)
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
	id := kernel.PromoID(r.PathValue("id"))

	p, err := h.svc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

func (h *Handler) deactivate(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)
	id := kernel.PromoID(r.PathValue("id"))

	p, err := h.svc.Deactivate(r.Context(), tenantID, id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, p)
}

type validateRequest struct {
	Code            string `json:"code"`
	OrderTotalCents int64  `json:"order_total_cents"`
}

func (h *Handler) validate(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromContext(r)

	var req validateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, errx.New("INVALID_REQUEST", "invalid request body", http.StatusBadRequest))
		return
	}

	result, err := h.svc.Validate(r.Context(), tenantID, req.Code, req.OrderTotalCents)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, result)
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
