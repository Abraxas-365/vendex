// Package renderer — resolver.go defines the data-access interfaces the renderer
// depends on and the lightweight value types for navigation and template overrides.
package renderer

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/catalog"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/product"
	"github.com/Abraxas-365/hada-commerce/internal/settings"
)

// ──────────────────────────────────────────────────────────────────────────────
// Data-source interfaces
// ──────────────────────────────────────────────────────────────────────────────

// ProductLister lists products for a tenant.
type ProductLister interface {
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[product.Product], error)
	ListByCategory(ctx context.Context, tenantID kernel.TenantID, categoryID kernel.CategoryID, pg kernel.PaginationOptions) (kernel.Paginated[product.Product], error)
}

// CollectionGetter retrieves a single collection by ID for a tenant.
type CollectionGetter interface {
	GetCollectionByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) (*catalog.Collection, error)
}

// SettingsGetter retrieves store-wide settings (name, logo, social links, etc.)
type SettingsGetter interface {
	Get(ctx context.Context, tenantID kernel.TenantID) (*settings.StoreSettings, error)
}

// ──────────────────────────────────────────────────────────────────────────────
// Navigation
// ──────────────────────────────────────────────────────────────────────────────

// NavLocation indicates where in the page layout a nav menu item appears.
type NavLocation string

const (
	NavLocationHeader NavLocation = "header"
	NavLocationFooter NavLocation = "footer"
)

// NavMenuItem is a single entry in a storefront navigation menu.
type NavMenuItem struct {
	ID       string
	Label    string
	URL      string
	Position int
	ParentID string // empty if root-level item
}

// NavMenuRepository loads navigation menu items for a tenant.
type NavMenuRepository interface {
	// ListByLocation returns all items at a given location (header/footer),
	// ordered by position ascending.
	ListByLocation(ctx context.Context, tenantID kernel.TenantID, location NavLocation) ([]NavMenuItem, error)
}

// ──────────────────────────────────────────────────────────────────────────────
// Template overrides
// ──────────────────────────────────────────────────────────────────────────────

// TemplateOverride holds a tenant-specific Go template string for a block type.
type TemplateOverride struct {
	BlockType string
	Template  string
}

// TemplateOverrideRepository loads per-tenant block template overrides.
type TemplateOverrideRepository interface {
	// GetByBlockType returns the override for the given block type,
	// or nil (no error) if none exists.
	GetByBlockType(ctx context.Context, tenantID kernel.TenantID, blockType string) (*TemplateOverride, error)
}
