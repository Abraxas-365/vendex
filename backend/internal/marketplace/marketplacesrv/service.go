package marketplacesrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/marketplace"
)

// VendorService implements marketplace vendor business logic.
type VendorService struct {
	vendorRepo  marketplace.VendorRepository
	productRepo marketplace.VendorProductRepository
	orderRepo   marketplace.VendorOrderRepository
}

// New creates a new VendorService.
func New(
	vendorRepo marketplace.VendorRepository,
	productRepo marketplace.VendorProductRepository,
	orderRepo marketplace.VendorOrderRepository,
) *VendorService {
	return &VendorService{
		vendorRepo:  vendorRepo,
		productRepo: productRepo,
		orderRepo:   orderRepo,
	}
}

// CreateVendor creates a new marketplace vendor.
func (s *VendorService) CreateVendor(ctx context.Context, tenantID kernel.TenantID, req marketplace.CreateVendorRequest) (marketplace.Vendor, error) {
	vendor := marketplace.Vendor{
		ID:          kernel.VendorID(uuid.New().String()),
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Email:       req.Email,
		Phone:       req.Phone,
		Status:      marketplace.VendorStatusPending,
		Commission:  req.Commission,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	return s.vendorRepo.Create(ctx, vendor)
}

// GetVendor retrieves a vendor by ID.
func (s *VendorService) GetVendor(ctx context.Context, tenantID kernel.TenantID, id kernel.VendorID) (marketplace.Vendor, error) {
	return s.vendorRepo.GetByID(ctx, tenantID, id)
}

// UpdateVendor applies partial updates to a vendor.
func (s *VendorService) UpdateVendor(ctx context.Context, tenantID kernel.TenantID, id kernel.VendorID, req marketplace.UpdateVendorRequest) (marketplace.Vendor, error) {
	vendor, err := s.vendorRepo.GetByID(ctx, tenantID, id)
	if err != nil {
		return marketplace.Vendor{}, err
	}

	if req.Name != nil {
		vendor.Name = *req.Name
	}
	if req.Description != nil {
		vendor.Description = *req.Description
	}
	if req.Email != nil {
		vendor.Email = *req.Email
	}
	if req.Phone != nil {
		vendor.Phone = *req.Phone
	}
	if req.Status != nil {
		vendor.Status = marketplace.VendorStatus(*req.Status)
	}
	if req.Commission != nil {
		vendor.Commission = *req.Commission
	}
	vendor.UpdatedAt = time.Now()

	return s.vendorRepo.Update(ctx, vendor)
}

// DeleteVendor removes a vendor by ID.
func (s *VendorService) DeleteVendor(ctx context.Context, tenantID kernel.TenantID, id kernel.VendorID) error {
	return s.vendorRepo.Delete(ctx, tenantID, id)
}

// ListVendors returns paginated vendors for a tenant.
func (s *VendorService) ListVendors(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.Vendor], error) {
	return s.vendorRepo.List(ctx, tenantID, p)
}

// AddVendorProduct links a product to a vendor with vendor-specific pricing.
func (s *VendorService) AddVendorProduct(ctx context.Context, tenantID kernel.TenantID, vp marketplace.VendorProduct) (marketplace.VendorProduct, error) {
	vp.ID = uuid.New().String()
	vp.TenantID = tenantID
	vp.CreatedAt = time.Now()
	return s.productRepo.Create(ctx, vp)
}

// ListVendorProducts returns paginated products for a vendor.
func (s *VendorService) ListVendorProducts(ctx context.Context, tenantID kernel.TenantID, vendorID kernel.VendorID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.VendorProduct], error) {
	return s.productRepo.ListByVendor(ctx, tenantID, vendorID, p)
}

// RemoveVendorProduct removes a vendor-product link.
func (s *VendorService) RemoveVendorProduct(ctx context.Context, tenantID kernel.TenantID, id string) error {
	return s.productRepo.Delete(ctx, tenantID, id)
}

// ListVendorOrders returns paginated orders for a vendor.
func (s *VendorService) ListVendorOrders(ctx context.Context, tenantID kernel.TenantID, vendorID kernel.VendorID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.VendorOrder], error) {
	return s.orderRepo.ListByVendor(ctx, tenantID, vendorID, p)
}
