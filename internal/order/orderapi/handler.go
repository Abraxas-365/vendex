package orderapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
	"github.com/Abraxas-365/hada-commerce/internal/order"
	"github.com/Abraxas-365/hada-commerce/internal/order/ordersrv"
)

// Handler exposes HTTP endpoints for the order domain.
type Handler struct {
	svc *ordersrv.Service
}

// NewHandler creates a new order API handler.
func NewHandler(svc *ordersrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all order routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /orders", h.Create)
	mux.HandleFunc("GET /orders/{id}", h.GetByID)
	mux.HandleFunc("GET /orders", h.List)
	mux.HandleFunc("PUT /orders/{id}/status", h.UpdateStatus)
	mux.HandleFunc("POST /orders/{id}/cancel", h.Cancel)
}

type createItemRequest struct {
	ProductID   string `json:"product_id"`
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	PriceAmount int64  `json:"price_amount"`
	Currency    string `json:"currency"`
}

type createRequest struct {
	CustomerID string              `json:"customer_id"`
	Items      []createItemRequest `json:"items"`
	Address    order.Address       `json:"shipping_address"`
}

// Create handles POST /orders.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	items := make([]ordersrv.CreateItemInput, len(req.Items))
	for i, it := range req.Items {
		items[i] = ordersrv.CreateItemInput{
			ProductID:   kernel.ProductID(it.ProductID),
			ProductName: it.ProductName,
			Quantity:    it.Quantity,
			UnitPrice:   kernel.NewMoney(it.PriceAmount, it.Currency),
		}
	}

	o, err := h.svc.Create(r.Context(), tenantID, ordersrv.CreateInput{
		CustomerID:      kernel.CustomerID(req.CustomerID),
		Items:           items,
		ShippingAddress: req.Address,
	})
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, o)
}

// GetByID handles GET /orders/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.OrderID(r.PathValue("id"))

	o, err := h.svc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, o)
}

// List handles GET /orders.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	pg := paginationFromQuery(r)

	customerID := r.URL.Query().Get("customer_id")
	var result kernel.PaginatedResult[order.Order]
	var err error

	if customerID != "" {
		result, err = h.svc.ListByCustomer(r.Context(), tenantID, kernel.CustomerID(customerID), pg)
	} else {
		result, err = h.svc.List(r.Context(), tenantID, pg)
	}
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

// UpdateStatus handles PUT /orders/{id}/status.
func (h *Handler) UpdateStatus(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.OrderID(r.PathValue("id"))

	var req updateStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	o, err := h.svc.UpdateStatus(r.Context(), tenantID, id, order.OrderStatus(req.Status))
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, o)
}

// Cancel handles POST /orders/{id}/cancel.
func (h *Handler) Cancel(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.OrderID(r.PathValue("id"))

	o, err := h.svc.Cancel(r.Context(), tenantID, id)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, o)
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
