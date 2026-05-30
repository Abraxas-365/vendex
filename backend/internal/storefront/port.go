package storefront

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PageRepository defines persistence operations for Page entities.
// All operations are scoped by TenantID.
type PageRepository interface {
	// Create persists a new page.
	Create(ctx context.Context, page *Page) error
	// GetByID retrieves a page by its ID within the tenant.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PageID) (*Page, error)
	// GetBySlug retrieves a page by its URL slug within the tenant.
	GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Page, error)
	// GetPublished retrieves a published page by slug — used for public serving.
	GetPublished(ctx context.Context, tenantID kernel.TenantID, slug string) (*Page, error)
	// Update persists changes to an existing page.
	Update(ctx context.Context, page *Page) error
	// ListByStatus returns all pages matching the given status with pagination.
	ListByStatus(ctx context.Context, tenantID kernel.TenantID, status PageStatus, p kernel.PaginationOptions) (kernel.Paginated[Page], error)
	// List returns all pages for a tenant with pagination.
	List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[Page], error)
}

// PageVersionRepository defines persistence for immutable page version snapshots.
type PageVersionRepository interface {
	// Create persists a new version snapshot. Versions are never deleted.
	Create(ctx context.Context, version *PageVersion) error
	// GetByVersion retrieves a specific version of a page.
	GetByVersion(ctx context.Context, tenantID kernel.TenantID, pageID kernel.PageID, version int) (*PageVersion, error)
	// ListByPage returns all versions for a page ordered by version desc.
	ListByPage(ctx context.Context, tenantID kernel.TenantID, pageID kernel.PageID) ([]PageVersion, error)
}

// BlockTypeRepository defines persistence operations for BlockType entities.
// BlockTypes are global (not per-tenant) — they define the available block vocabulary.
type BlockTypeRepository interface {
	// Create persists a new block type.
	Create(ctx context.Context, bt *BlockType) error
	// GetByID retrieves a block type by its ID.
	GetByID(ctx context.Context, id kernel.BlockTypeID) (*BlockType, error)
	// GetByName retrieves a block type by its unique name.
	GetByName(ctx context.Context, name string) (*BlockType, error)
	// List returns all block types, optionally filtered by category (empty string = all).
	List(ctx context.Context, category string) ([]BlockType, error)
	// Update persists changes to an existing block type.
	Update(ctx context.Context, bt *BlockType) error
	// Delete removes a block type by ID.
	Delete(ctx context.Context, id kernel.BlockTypeID) error
}
