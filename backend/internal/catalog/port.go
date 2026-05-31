package catalog

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// CategoryRepository defines persistence operations for categories.
type CategoryRepository interface {
	Create(ctx context.Context, c *Category) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CategoryID) (*Category, error)
	GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Category, error)
	Update(ctx context.Context, c *Category) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CategoryID) error
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[Category], error)
	ListByParent(ctx context.Context, tenantID kernel.TenantID, parentID *kernel.CategoryID, pg kernel.PaginationOptions) (kernel.Paginated[Category], error)
}

// CollectionRepository defines persistence operations for collections.
type CollectionRepository interface {
	Create(ctx context.Context, c *Collection) error
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) (*Collection, error)
	GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Collection, error)
	Update(ctx context.Context, c *Collection) error
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) error
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[Collection], error)
}
