package marketplace

import (
	"encoding/json"
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// PresetStatus represents the lifecycle of a preset.
type PresetStatus string

const (
	PresetStatusDraft     PresetStatus = "draft"
	PresetStatusPublished PresetStatus = "published"
	PresetStatusArchived  PresetStatus = "archived"
)

// PresetVisibility controls who can see/install a preset.
type PresetVisibility string

const (
	PresetVisibilityPublic  PresetVisibility = "public"  // visible in marketplace
	PresetVisibilityPrivate PresetVisibility = "private" // only owner tenant
)

// Preset defines a reusable agent workspace configuration.
// It bundles: Docker image, tools, system prompt, and custom frontend.
type Preset struct {
	ID            kernel.PresetID  `json:"id"`
	TenantID      kernel.TenantID  `json:"tenant_id"`
	Name          string           `json:"name"`
	Slug          string           `json:"slug"`
	Description   string           `json:"description"`
	Version       string           `json:"version"`
	Image         string           `json:"image"`
	FrontendPort  int              `json:"frontend_port"`
	SystemPrompt  string           `json:"system_prompt"`
	ToolsManifest json.RawMessage  `json:"tools_manifest"`
	Status        PresetStatus     `json:"status"`
	Visibility    PresetVisibility `json:"visibility"`
	Icon          string           `json:"icon"`
	Tags          []string         `json:"tags"`
	InstallCount  int              `json:"install_count"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

// PresetInstall tracks which tenants have installed which presets.
type PresetInstall struct {
	ID          string          `json:"id"`
	TenantID    kernel.TenantID `json:"tenant_id"`
	PresetID    kernel.PresetID `json:"preset_id"`
	InstalledAt time.Time       `json:"installed_at"`
	Config      json.RawMessage `json:"config"`
}

// CreatePresetRequest holds input for creating a preset.
type CreatePresetRequest struct {
	Name          string           `json:"name"`
	Slug          string           `json:"slug"`
	Description   string           `json:"description"`
	Version       string           `json:"version"`
	Image         string           `json:"image"`
	FrontendPort  int              `json:"frontend_port"`
	SystemPrompt  string           `json:"system_prompt"`
	ToolsManifest json.RawMessage  `json:"tools_manifest"`
	Visibility    PresetVisibility `json:"visibility"`
	Icon          string           `json:"icon"`
	Tags          []string         `json:"tags"`
}

// UpdatePresetRequest holds partial update fields.
type UpdatePresetRequest struct {
	Name          *string           `json:"name"`
	Description   *string           `json:"description"`
	Version       *string           `json:"version"`
	Image         *string           `json:"image"`
	FrontendPort  *int              `json:"frontend_port"`
	SystemPrompt  *string           `json:"system_prompt"`
	ToolsManifest *json.RawMessage  `json:"tools_manifest"`
	Status        *PresetStatus     `json:"status"`
	Visibility    *PresetVisibility `json:"visibility"`
	Icon          *string           `json:"icon"`
	Tags          *[]string         `json:"tags"`
}

// PresetListOptions configures the preset list query (marketplace browsing).
type PresetListOptions struct {
	kernel.PaginationOptions
	Status     *PresetStatus
	Visibility *PresetVisibility
	Tag        string
	Search     string
}
