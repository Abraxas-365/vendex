package plugin

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PluginRepository defines persistence operations for plugins (global catalogue).
type PluginRepository interface {
	Create(ctx context.Context, p *Plugin) error
	GetByID(ctx context.Context, id kernel.PluginID) (*Plugin, error)
	Update(ctx context.Context, p *Plugin) error
	Delete(ctx context.Context, id kernel.PluginID) error
	List(ctx context.Context, pg kernel.PaginationOptions) (kernel.Paginated[Plugin], error)
}

// PluginVersionRepository defines persistence operations for plugin versions.
type PluginVersionRepository interface {
	Create(ctx context.Context, v *PluginVersion) error
	GetByID(ctx context.Context, id kernel.PluginVersionID) (*PluginVersion, error)
	ListByPlugin(ctx context.Context, pluginID kernel.PluginID, pg kernel.PaginationOptions) (kernel.Paginated[PluginVersion], error)
	GetLatest(ctx context.Context, pluginID kernel.PluginID) (*PluginVersion, error)
}

// InstallationRepository defines persistence operations for plugin installations.
type InstallationRepository interface {
	Create(ctx context.Context, i *PluginInstallation) error
	GetByTenantAndPlugin(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*PluginInstallation, error)
	ListByTenant(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[PluginInstallation], error)
	ListActiveByTenant(ctx context.Context, tenantID kernel.TenantID) ([]PluginInstallation, error)
	// GetJSManifestData returns script entries for all active installations that have
	// a non-empty frontend_url. It joins plugin_installations with plugins and
	// plugin_versions in a single query for efficiency.
	GetJSManifestData(ctx context.Context, tenantID kernel.TenantID) ([]PluginScript, error)
	Update(ctx context.Context, i *PluginInstallation) error
	Delete(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) error
}
