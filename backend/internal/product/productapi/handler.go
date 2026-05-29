package productapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
	"github.com/Abraxas-365/hada-commerce/internal/product/productsrv"
)

// Handler exposes HTTP endpoints for the product domain.
type Handler struct {
	svc *productsrv.Service
}

// NewHandler creates a new product API handler.
func NewHandler(svc *productsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers all product routes on the given mux.
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /products", h.Create)
	mux.HandleFunc("GET /products/{id}", h.GetByID)
	mux.HandleFunc("GET /products", h.List)
	mux.HandleFunc("PUT /products/{id}", h.Update)
	mux.HandleFunc("DELETE /products/{id}", h.Delete)
}

// createRequest is the JSON body for creating a product.
type createRequest struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	PriceAmount int64        `json:"price_amount"`
	Currency    string       `json:"currency"`
	SKU         string       `json:"sku"`
	Images      []string     `json:"images"`
	CategoryID  string       `json:"category_id"`
	Tags        []string     `json:"tags"`
	Stock       int          `json:"stock"`
}

// Create handles POST /products.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	p, err := h.svc.Create(r.Context(), tenantID, productsrv.CreateInput{
		Name:        req.Name,
		Description: req.Description,
		Price:       kernel.NewMoney(req.PriceAmount, req.Currency),
		SKU:         req.SKU,
		Images:      req.Images,
		CategoryID:  kernel.CategoryID(req.CategoryID),
		Tags:        req.Tags,
		Stock:       req.Stock,
	})
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, p)
}

// GetByID handles GET /products/{id}.
func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.ProductID(r.PathValue("id"))

	p, err := h.svc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// List handles GET /products.
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

// Update handles PUT /products/{id}.
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.ProductID(r.PathValue("id"))

	existing, err := h.svc.GetByID(r.Context(), tenantID, id)
	if err != nil {
		writeErrx(w, err)
		return
	}

	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	existing.Name = req.Name
	existing.Description = req.Description
	existing.Price = kernel.NewMoney(req.PriceAmount, req.Currency)
	existing.SKU = req.SKU
	existing.Images = req.Images
	existing.CategoryID = kernel.CategoryID(req.CategoryID)
	existing.Tags = req.Tags
	existing.Stock = req.Stock

	if err := h.svc.Update(r.Context(), existing); err != nil {
		writeErrx(w, err)
		return
	}

	writeJSON(w, http.StatusOK, existing)
}

// Delete handles DELETE /products/{id}.
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	tenantID := tenantFromRequest(r)
	id := kernel.ProductID(r.PathValue("id"))

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
