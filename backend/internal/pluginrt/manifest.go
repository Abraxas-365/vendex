package pluginrt

import (
	"encoding/json"
	"fmt"
)

// Manifest describes a plugin's capabilities.
// Parsed from manifest.json stored in PluginVersion.ManifestJSON.
type Manifest struct {
	Name        string       `json:"name"`
	DisplayName string       `json:"display_name"`
	Version     string       `json:"version"`
	Description string       `json:"description"`
	Author      string       `json:"author"`
	Permissions []string     `json:"permissions"`
	UI          UIConfig     `json:"ui"`
	Tools       []ToolConfig `json:"tools"`
	Migrations  []string     `json:"migrations"`
}

// UIConfig describes the UI extension points a plugin provides.
type UIConfig struct {
	Tabs    []TabConfig    `json:"tabs"`
	Widgets []WidgetConfig `json:"widgets"`
}

// TabConfig describes a plugin-provided navigation tab.
type TabConfig struct {
	Label string `json:"label"`
	Icon  string `json:"icon"`
	Entry string `json:"entry"` // relative path to plugin frontend HTML
}

// WidgetConfig describes a plugin-provided UI widget injected into a platform slot.
type WidgetConfig struct {
	Slot      string `json:"slot"`      // "product-detail", "checkout", "dashboard"
	Component string `json:"component"`
	Entry     string `json:"entry"` // relative path to widget JS bundle
}

// ToolConfig describes an agent tool exposed by the plugin.
type ToolConfig struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ParseManifest parses a plugin manifest from raw JSON bytes.
func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("parsing plugin manifest: %w", err)
	}
	return &m, nil
}
