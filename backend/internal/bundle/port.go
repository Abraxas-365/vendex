package bundle

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines persistence operations for the bundle domain.
type Repository interface {
	// Bundle CRUD
	Create(ctx context.Context, b *Bundle) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BundleID) (*Bundle, error)
	GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Bundle, error)
	Update(ctx context.Context, b *Bundle) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.BundleID) error
	List(ctx context.Context, tenantID kernel.TenantID, activeOnly bool, pg kernel.PaginationOptions) (kernel.Paginated[Bundle], error)

	// Bundle items
	AddItem(ctx context.Context, item *BundleItem) error
	GetItemByID(ctx context.Context, tenantID kernel.TenantID, itemID kernel.BundleItemID) (*BundleItem, error)
	ListItems(ctx context.Context, tenantID kernel.TenantID, bundleID kernel.BundleID) ([]BundleItem, error)
	RemoveItem(ctx context.Context, tenantID kernel.TenantID, itemID kernel.BundleItemID) error
}
