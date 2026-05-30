package plugin

import (
	"encoding/json"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Plugin is a global (non-tenant-scoped) plugin in the catalogue.
type Plugin struct {
	ID          kernel.PluginID `json:"id" db:"id"`
	Name        string          `json:"name" db:"name"`
	DisplayName string          `json:"display_name" db:"display_name"`
	Description string          `json:"description" db:"description"`
	Author      string          `json:"author" db:"author"`
	Icon        string          `json:"icon" db:"icon"`
	Category    string          `json:"category" db:"category"`
	Tags        json.RawMessage `json:"tags" db:"tags"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// PluginVersion represents a published version of a plugin.
type PluginVersion struct {
	ID             kernel.PluginVersionID `json:"id" db:"id"`
	PluginID       kernel.PluginID        `json:"plugin_id" db:"plugin_id"`
	Version        string                 `json:"version" db:"version"`
	Changelog      string                 `json:"changelog" db:"changelog"`
	Permissions    json.RawMessage        `json:"permissions" db:"permissions"`
	ManifestJSON   string                 `json:"manifest_json" db:"manifest_json"`
	FrontendURL    string                 `json:"frontend_url" db:"frontend_url"`
	BackendEntry   string                 `json:"backend_entry" db:"backend_entry"`
	MinPlatformVer string                 `json:"min_platform_ver" db:"min_platform_ver"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
}

// InstallationStatus represents the lifecycle state of an installation.
type InstallationStatus string

const (
	StatusActive   InstallationStatus = "active"
	StatusInactive InstallationStatus = "inactive"
	StatusFailed   InstallationStatus = "failed"
)

// PluginInstallation represents a per-tenant plugin activation record.
type PluginInstallation struct {
	ID          kernel.InstallationID  `json:"id" db:"id"`
	TenantID    kernel.TenantID        `json:"tenant_id" db:"tenant_id"`
	PluginID    kernel.PluginID        `json:"plugin_id" db:"plugin_id"`
	VersionID   kernel.PluginVersionID `json:"version_id" db:"version_id"`
	Status      InstallationStatus     `json:"status" db:"status"`
	Settings    json.RawMessage        `json:"settings" db:"settings"`
	InstalledAt time.Time              `json:"installed_at" db:"installed_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
}
