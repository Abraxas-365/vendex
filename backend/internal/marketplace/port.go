package marketplace

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for the marketplace domain.
type Repository interface {
	// Plugins
	CreatePlugin(ctx context.Context, p *Plugin) error
	GetPlugin(ctx context.Context, id kernel.PluginID) (*Plugin, error)
	GetPluginByName(ctx context.Context, name string) (*Plugin, error)
	ListPlugins(ctx context.Context, pg kernel.Pagination) (kernel.PaginatedResult[Plugin], error)
	ListPluginsByCategory(ctx context.Context, cat PluginCategory, pg kernel.Pagination) (kernel.PaginatedResult[Plugin], error)

	// Versions
	CreateVersion(ctx context.Context, v *PluginVersion) error
	GetVersion(ctx context.Context, id kernel.PluginVersionID) (*PluginVersion, error)
	GetLatestVersion(ctx context.Context, pluginID kernel.PluginID) (*PluginVersion, error)
	ListVersions(ctx context.Context, pluginID kernel.PluginID) ([]PluginVersion, error)

	// Installations
	CreateInstallation(ctx context.Context, inst *Installation) error
	GetInstallation(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*Installation, error)
	UpdateInstallation(ctx context.Context, inst *Installation) error
	DeleteInstallation(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) error
	ListInstallations(ctx context.Context, tenantID kernel.TenantID) ([]Installation, error)
}
