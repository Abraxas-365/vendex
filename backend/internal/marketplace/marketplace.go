package marketplace

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PluginCategory classifies plugins by origin.
type PluginCategory string

const (
	CategoryOfficial  PluginCategory = "official"
	CategoryCommunity PluginCategory = "community"
	CategoryCustom    PluginCategory = "custom"
)

// Plugin represents a plugin available in the marketplace.
type Plugin struct {
	ID          kernel.PluginID `json:"id"`
	Name        string          `json:"name"`         // unique slug
	DisplayName string          `json:"display_name"`
	Description string          `json:"description"`
	Author      string          `json:"author"`
	Icon        string          `json:"icon"`         // URL or icon name
	Category    PluginCategory  `json:"category"`     // "official", "community", "custom"
	Tags        []string        `json:"tags"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// PluginVersion is a specific release of a plugin.
type PluginVersion struct {
	ID             kernel.PluginVersionID `json:"id"`
	PluginID       kernel.PluginID        `json:"plugin_id"`
	Version        string                 `json:"version"`
	Changelog      string                 `json:"changelog"`
	Permissions    []string               `json:"permissions"`
	ManifestJSON   string                 `json:"manifest_json"`
	FrontendURL    string                 `json:"frontend_url"`
	BackendEntry   string                 `json:"backend_entry"`
	MinPlatformVer string                 `json:"min_platform_ver"`
	CreatedAt      time.Time              `json:"created_at"`
}

// Manifest describes plugin capabilities (parsed from ManifestJSON).
type Manifest struct {
	Name        string         `json:"name"`
	DisplayName string         `json:"display_name"`
	Version     string         `json:"version"`
	Description string         `json:"description"`
	Author      string         `json:"author"`
	Permissions []string       `json:"permissions"`
	UI          ManifestUI     `json:"ui"`
	Tools       []ManifestTool `json:"tools"`
	Migrations  []string       `json:"migrations"`
}

// ManifestUI describes the UI extensions a plugin provides.
type ManifestUI struct {
	Tabs    []ManifestTab    `json:"tabs"`
	Widgets []ManifestWidget `json:"widgets"`
}

// ManifestTab describes a navigation tab added by a plugin.
type ManifestTab struct {
	Label string `json:"label"`
	Icon  string `json:"icon"`
	Entry string `json:"entry"`
}

// ManifestWidget describes a UI widget injected into a slot.
type ManifestWidget struct {
	Slot      string `json:"slot"`      // "product-detail", "checkout", "dashboard"
	Component string `json:"component"`
	Entry     string `json:"entry"`
}

// ManifestTool describes an AI tool exposed by a plugin.
type ManifestTool struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ParseManifest parses a Manifest from JSON data.
func ParseManifest(data string) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal([]byte(data), &m); err != nil {
		return nil, fmt.Errorf("parsing manifest: %w", err)
	}
	return &m, nil
}

// HasPermission returns true if the manifest declares the given permission.
func (m *Manifest) HasPermission(perm string) bool {
	for _, p := range m.Permissions {
		if p == perm {
			return true
		}
	}
	return false
}

// InstallationStatus represents the state of an installed plugin.
type InstallationStatus string

const (
	InstallActive   InstallationStatus = "active"
	InstallInactive InstallationStatus = "inactive"
	InstallFailed   InstallationStatus = "failed"
)

// Installation tracks which plugins are installed for a tenant.
type Installation struct {
	ID          kernel.InstallationID
	TenantID    kernel.TenantID
	PluginID    kernel.PluginID
	VersionID   kernel.PluginVersionID
	Status      InstallationStatus
	Settings    map[string]any // tenant-specific plugin config
	InstalledAt time.Time
	UpdatedAt   time.Time
}

// Activate transitions the installation to active status.
func (i *Installation) Activate() {
	i.Status = InstallActive
	i.UpdatedAt = time.Now()
}

// Deactivate transitions the installation to inactive status.
func (i *Installation) Deactivate() {
	i.Status = InstallInactive
	i.UpdatedAt = time.Now()
}
