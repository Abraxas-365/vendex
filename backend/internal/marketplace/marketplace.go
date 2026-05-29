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
	ID          kernel.PluginID
	Name        string // unique slug: "reviews", "loyalty", etc.
	DisplayName string
	Description string
	Author      string
	Icon        string         // URL or icon name
	Category    PluginCategory // "official", "community", "custom"
	Tags        []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// PluginVersion is a specific release of a plugin.
type PluginVersion struct {
	ID             kernel.PluginVersionID
	PluginID       kernel.PluginID
	Version        string // semver: "1.0.0"
	Changelog      string
	Permissions    []string // e.g. ["products:read", "orders:read"]
	ManifestJSON   string   // the full manifest.json content
	FrontendURL    string   // URL to the plugin's frontend bundle
	BackendEntry   string   // Go plugin entry point or binary path
	MinPlatformVer string   // minimum hada-commerce version
	CreatedAt      time.Time
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
