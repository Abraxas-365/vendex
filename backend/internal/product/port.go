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
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.Pagination) (kernel.PaginatedResult[Product], error)
	ListByCategory(ctx context.Context, tenantID kernel.TenantID, categoryID kernel.CategoryID, pg kernel.Pagination) (kernel.PaginatedResult[Product], error)
	GetBySKU(ctx context.Context, tenantID kernel.TenantID, sku string) (*Product, error)
}
