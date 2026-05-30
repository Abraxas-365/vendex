package product

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for products.
type Repository interface {
	Create(ctx context.Context, p *Product) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) (*Product, error)
	Update(ctx context.Context, p *Product) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) error
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[Product], error)
	ListByCategory(ctx context.Context, tenantID kernel.TenantID, categoryID kernel.CategoryID, pg kernel.PaginationOptions) (kernel.Paginated[Product], error)
	GetBySKU(ctx context.Context, tenantID kernel.TenantID, sku string) (*Product, error)
	GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Product, error)
}

// VariantRepository defines persistence operations for product options and variants.
type VariantRepository interface {
	// Options
	CreateOption(ctx context.Context, opt *ProductOption) error
	ListOptions(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) ([]ProductOption, error)
	UpdateOption(ctx context.Context, opt *ProductOption) error
	DeleteOption(ctx context.Context, tenantID kernel.TenantID, id kernel.OptionID) error

	// Variants
	CreateVariant(ctx context.Context, v *ProductVariant) error
	GetVariantByID(ctx context.Context, tenantID kernel.TenantID, id kernel.VariantID) (*ProductVariant, error)
	ListVariants(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) ([]ProductVariant, error)
	UpdateVariant(ctx context.Context, v *ProductVariant) error
	DeleteVariant(ctx context.Context, tenantID kernel.TenantID, id kernel.VariantID) error
	GetVariantBySKU(ctx context.Context, tenantID kernel.TenantID, sku string) (*ProductVariant, error)
}
