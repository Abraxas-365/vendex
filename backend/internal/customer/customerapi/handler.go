package customerapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/customer"
	"github.com/Abraxas-365/hada-commerce/internal/customer/customersrv"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

// Handler exposes HTTP endpoints for the customer domain.
type Handler struct {
	svc *customersrv.Service
}

// NewHandler creates a new customer API handler.
func NewHandler(svc *customersrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all customer routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /customers", h.Create)
	mux.HandleFunc("GET /customers/{id}", h.GetByID)
	mux.HandleFunc("GET /customers", h.List)
	mux.HandleFunc("DELETE /customers/{id}", h.Delete)
}

type createRequest struct {
	Email     string             `json:"email"`
	Name      string             `json:"name"`
	Phone     string             `json:"phone"`
	Addresses []customer.Address `json:"addresses"`
}

// Create handles POST /customers.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	c, err := h.svc.Create(r.Context(), tenantID, customersrv.CreateInput{
		Email:     req.Email,
		Name:      req.Name,
		Phone:     req.Phone,
		Addresses: req.Addresses,
	})
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, c)
}

// GetByID handles GET /customers/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.CustomerID(r.PathValue("id"))

	c, err := h.svc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, c)
}

// List handles GET /customers.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	pg := paginationFromQuery(r)

	result, err := h.svc.List(r.Context(), tenantID, pg)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// Delete handles DELETE /customers/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.CustomerID(r.PathValue("id"))

	if err := h.svc.Delete(r.Context(), tenantID, id); err != nil {
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
