package multistore

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines the persistence contract for the multistore domain.
type Repository interface {
	// Create persists a new storefront.
	Create(ctx context.Context, sf *Storefront) error

	// GetByID returns a storefront scoped to the tenant.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID) (*Storefront, error)

	// GetBySlug returns a storefront by its slug within a tenant.
	GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*Storefront, error)

	// GetByDomain returns a storefront by its custom domain.
	// Domain lookups are global (not tenant-scoped) since domains are globally unique.
	GetByDomain(ctx context.Context, domain string) (*Storefront, error)

	// List returns paginated storefronts for a tenant.
	List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[Storefront], error)

	// Update persists changes to an existing storefront.
	Update(ctx context.Context, sf *Storefront) error

	// Delete removes a storefront.
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID) error

	// ClearDefault clears the is_default flag on all storefronts for the tenant.
	ClearDefault(ctx context.Context, tenantID kernel.TenantID) error

	// SetDefault marks a storefront as the default and clears all others.
	SetDefault(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID) error

	// AddCatalog links a catalog to a storefront.
	AddCatalog(ctx context.Context, sc *StorefrontCatalog) error

	// RemoveCatalog removes a catalog link from a storefront.
	RemoveCatalog(ctx context.Context, tenantID kernel.TenantID, storefrontID kernel.StorefrontEntryID, catalogID string) error

	// ListCatalogs returns all catalog links for a storefront.
	ListCatalogs(ctx context.Context, tenantID kernel.TenantID, storefrontID kernel.StorefrontEntryID) ([]StorefrontCatalog, error)
}
