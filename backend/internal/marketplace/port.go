package marketplace

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// VendorRepository defines persistence operations for Vendor records.
type VendorRepository interface {
	Create(ctx context.Context, v Vendor) (Vendor, error)
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.VendorID) (Vendor, error)
	Update(ctx context.Context, v Vendor) (Vendor, error)
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.VendorID) error
	List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[Vendor], error)
}

// VendorProductRepository defines persistence for VendorProduct link records.
type VendorProductRepository interface {
	Create(ctx context.Context, vp VendorProduct) (VendorProduct, error)
	GetByID(ctx context.Context, tenantID kernel.TenantID, id string) (VendorProduct, error)
	Delete(ctx context.Context, tenantID kernel.TenantID, id string) error
	ListByVendor(ctx context.Context, tenantID kernel.TenantID, vendorID kernel.VendorID, p kernel.PaginationOptions) (kernel.Paginated[VendorProduct], error)
}

// VendorOrderRepository defines persistence for vendor-order accounting records.
type VendorOrderRepository interface {
	Create(ctx context.Context, vo VendorOrder) (VendorOrder, error)
	GetByID(ctx context.Context, tenantID kernel.TenantID, id string) (VendorOrder, error)
	ListByVendor(ctx context.Context, tenantID kernel.TenantID, vendorID kernel.VendorID, p kernel.PaginationOptions) (kernel.Paginated[VendorOrder], error)
}
