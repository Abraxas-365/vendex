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

// PresetRepository defines persistence operations for Preset records.
type PresetRepository interface {
	Create(ctx context.Context, p Preset) (Preset, error)
	GetByID(ctx context.Context, id kernel.PresetID) (Preset, error)
	GetBySlug(ctx context.Context, slug string) (Preset, error)
	Update(ctx context.Context, p Preset) (Preset, error)
	Delete(ctx context.Context, id kernel.PresetID) error
	List(ctx context.Context, opts PresetListOptions) (kernel.Paginated[Preset], error)
	ListByTenant(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[Preset], error)
}

// PresetInstallRepository tracks preset installations per tenant.
type PresetInstallRepository interface {
	Install(ctx context.Context, install PresetInstall) (PresetInstall, error)
	Uninstall(ctx context.Context, tenantID kernel.TenantID, presetID kernel.PresetID) error
	ListByTenant(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[PresetInstall], error)
	IsInstalled(ctx context.Context, tenantID kernel.TenantID, presetID kernel.PresetID) (bool, error)
}
