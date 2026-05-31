package collection

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines all persistence operations for the collection domain.
type Repository interface {
	// --- Collection CRUD ---

	// Create persists a new collection.
	Create(ctx context.Context, c *Collection) error

	// GetByID returns a collection by its ID, scoped to the tenant.
	// Returns ErrNotFound when no matching row exists.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) (*Collection, error)

	// GetBySlug returns a collection by its URL slug, scoped to the tenant.
	// Returns ErrNotFound when no matching row exists.
	GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Collection, error)

	// Update persists changes to an existing collection.
	Update(ctx context.Context, c *Collection) error

	// Delete removes a collection by ID, scoped to the tenant.
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) error

	// List returns a paginated slice of collections for the given tenant.
	// When activeOnly is true only collections with is_active=true are returned.
	List(ctx context.Context, tenantID kernel.TenantID, activeOnly bool, pg kernel.PaginationOptions) (kernel.Paginated[Collection], error)

	// CountProducts returns the number of products belonging to the collection.
	CountProducts(ctx context.Context, tenantID kernel.TenantID, collectionID kernel.CollectionID) (int, error)

	// --- Collection product membership ---

	// AddProduct creates a CollectionProduct record.
	// Returns ErrAlreadyInCollection when the product is already a member.
	AddProduct(ctx context.Context, cp *CollectionProduct) error

	// RemoveProduct deletes the membership record.
	// Returns ErrProductNotFound when the product is not a member.
	RemoveProduct(ctx context.Context, tenantID kernel.TenantID, collectionID kernel.CollectionID, productID string) error

	// ListProducts returns a paginated list of products in the collection.
	ListProducts(ctx context.Context, tenantID kernel.TenantID, collectionID kernel.CollectionID, pg kernel.PaginationOptions) (kernel.Paginated[CollectionProduct], error)

	// ReorderProducts sets the sort_order for every provided productID.
	// productIDs are applied in order (index 0 → sort_order 0, index 1 → sort_order 1, …).
	ReorderProducts(ctx context.Context, tenantID kernel.TenantID, collectionID kernel.CollectionID, productIDs []string) error
}
