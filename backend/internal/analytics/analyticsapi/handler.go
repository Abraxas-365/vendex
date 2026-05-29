package analyticsapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/analytics/analyticssrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

// Handler exposes HTTP endpoints for the analytics domain.
type Handler struct {
	svc *analyticssrv.Service
}

// NewHandler creates a new analytics API handler.
func NewHandler(svc *analyticssrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all analytics routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /analytics/dashboard", h.GetDashboardStats)
	mux.HandleFunc("GET /analytics/revenue", h.GetRevenueTimeline)
	mux.HandleFunc("GET /analytics/top-products", h.GetTopProducts)
	mux.HandleFunc("GET /analytics/order-status", h.GetOrderStatusBreakdown)
	mux.HandleFunc("GET /analytics/recent-orders", h.GetRecentOrders)
}

// GetDashboardStats handles GET /analytics/dashboard.
func (h *Handler) GetDashboardStats(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	stats, err := h.svc.GetDashboardStats(r.Context(), tenantID)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, stats)
}

// GetRevenueTimeline handles GET /analytics/revenue?days=30.
func (h *Handler) GetRevenueTimeline(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	days := 30
	if d, err := strconv.Atoi(r.URL.Query().Get("days")); err == nil && d > 0 {
		days = d
	}

	points, err := h.svc.GetRevenueTimeline(r.Context(), tenantID, days)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, points)
}

// GetTopProducts handles GET /analytics/top-products?limit=5.
func (h *Handler) GetTopProducts(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	limit := 5
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}

	products, err := h.svc.GetTopProducts(r.Context(), tenantID, limit)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, products)
}

// GetOrderStatusBreakdown handles GET /analytics/order-status.
func (h *Handler) GetOrderStatusBreakdown(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	breakdown, err := h.svc.GetOrderStatusBreakdown(r.Context(), tenantID)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, breakdown)
}

// GetRecentOrders handles GET /analytics/recent-orders?limit=5.
func (h *Handler) GetRecentOrders(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	limit := 5
	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}

	orders, err := h.svc.GetRecentOrders(r.Context(), tenantID, limit)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, orders)
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
