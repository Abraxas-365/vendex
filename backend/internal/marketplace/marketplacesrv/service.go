package marketplacesrv

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
	"github.com/Abraxas-365/hada-commerce/internal/marketplace"
)

// Service handles marketplace business logic.
type Service struct {
	repo marketplace.Repository
}

// New creates a new marketplace service.
func New(repo marketplace.Repository) *Service {
	return &Service{repo: repo}
}

// ListAvailable returns a paginated list of all plugins in the marketplace.
func (s *Service) ListAvailable(ctx context.Context, pg kernel.Pagination) (kernel.PaginatedResult[marketplace.Plugin], error) {
	return s.repo.ListPlugins(ctx, pg)
}

// PluginDetail holds a plugin together with its latest version.
type PluginDetail struct {
	Plugin        marketplace.Plugin        `json:"plugin"`
	LatestVersion *marketplace.PluginVersion `json:"latest_version,omitempty"`
}

// GetPlugin returns a plugin with its latest version.
func (s *Service) GetPlugin(ctx context.Context, id kernel.PluginID) (*PluginDetail, error) {
	p, err := s.repo.GetPlugin(ctx, id)
	if err != nil {
		return nil, err
	}

	latest, err := s.repo.GetLatestVersion(ctx, id)
	if err != nil && !errx.Is(err, marketplace.ErrVersionNotFound) {
		return nil, fmt.Errorf("getting latest version: %w", err)
	}

	return &PluginDetail{Plugin: *p, LatestVersion: latest}, nil
}

// Install installs the latest version of a plugin for a tenant.
func (s *Service) Install(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*marketplace.Installation, error) {
	// Check plugin exists.
	if _, err := s.repo.GetPlugin(ctx, pluginID); err != nil {
		return nil, err
	}

	// Check not already installed.
	existing, err := s.repo.GetInstallation(ctx, tenantID, pluginID)
	if err != nil && !errx.Is(err, marketplace.ErrNotInstalled) {
		return nil, fmt.Errorf("checking existing installation: %w", err)
	}
	if existing != nil {
		return nil, marketplace.ErrAlreadyInstalled
	}

	// Get latest version.
	latest, err := s.repo.GetLatestVersion(ctx, pluginID)
	if err != nil {
		return nil, fmt.Errorf("getting latest version: %w", err)
	}

	now := time.Now()
	inst := &marketplace.Installation{
		ID:          kernel.InstallationID(generateID()),
		TenantID:    tenantID,
		PluginID:    pluginID,
		VersionID:   latest.ID,
		Status:      marketplace.InstallActive,
		Settings:    map[string]any{},
		InstalledAt: now,
		UpdatedAt:   now,
	}

	if err := s.repo.CreateInstallation(ctx, inst); err != nil {
		return nil, fmt.Errorf("creating installation: %w", err)
	}
	return inst, nil
}

// Uninstall removes a plugin installation for a tenant.
func (s *Service) Uninstall(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) error {
	// Verify it's installed.
	if _, err := s.repo.GetInstallation(ctx, tenantID, pluginID); err != nil {
		return err
	}
	return s.repo.DeleteInstallation(ctx, tenantID, pluginID)
}

// Activate sets an installed plugin to active status.
func (s *Service) Activate(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*marketplace.Installation, error) {
	inst, err := s.repo.GetInstallation(ctx, tenantID, pluginID)
	if err != nil {
		return nil, err
	}

	inst.Activate()
	if err := s.repo.UpdateInstallation(ctx, inst); err != nil {
		return nil, fmt.Errorf("updating installation: %w", err)
	}
	return inst, nil
}

// Deactivate sets an installed plugin to inactive status.
func (s *Service) Deactivate(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*marketplace.Installation, error) {
	inst, err := s.repo.GetInstallation(ctx, tenantID, pluginID)
	if err != nil {
		return nil, err
	}

	inst.Deactivate()
	if err := s.repo.UpdateInstallation(ctx, inst); err != nil {
		return nil, fmt.Errorf("updating installation: %w", err)
	}
	return inst, nil
}

// ListInstalled returns all installations for a tenant.
func (s *Service) ListInstalled(ctx context.Context, tenantID kernel.TenantID) ([]marketplace.Installation, error) {
	return s.repo.ListInstallations(ctx, tenantID)
}

// UpdateSettings updates the tenant-specific configuration for an installed plugin.
func (s *Service) UpdateSettings(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID, settings map[string]any) (*marketplace.Installation, error) {
	inst, err := s.repo.GetInstallation(ctx, tenantID, pluginID)
	if err != nil {
		return nil, err
	}

	inst.Settings = settings
	inst.UpdatedAt = time.Now()
	if err := s.repo.UpdateInstallation(ctx, inst); err != nil {
		return nil, fmt.Errorf("updating settings: %w", err)
	}
	return inst, nil
}

// PublishPluginInput holds the data needed to publish a new plugin.
type PublishPluginInput struct {
	Name        string
	DisplayName string
	Description string
	Author      string
	Icon        string
	Category    marketplace.PluginCategory
	Tags        []string
}

// PublishVersionInput holds the data needed for the first version.
type PublishVersionInput struct {
	Version        string
	Changelog      string
	Permissions    []string
	ManifestJSON   string
	FrontendURL    string
	BackendEntry   string
	MinPlatformVer string
}

// PublishPlugin creates a new plugin with its first version.
func (s *Service) PublishPlugin(ctx context.Context, pluginIn PublishPluginInput, versionIn PublishVersionInput) (*marketplace.Plugin, *marketplace.PluginVersion, error) {
	// Check name uniqueness.
	existing, err := s.repo.GetPluginByName(ctx, pluginIn.Name)
	if err != nil && !errx.Is(err, marketplace.ErrPluginNotFound) {
		return nil, nil, fmt.Errorf("checking plugin name: %w", err)
	}
	if existing != nil {
		return nil, nil, marketplace.ErrPluginNameTaken
	}

	now := time.Now()
	pluginID := kernel.PluginID(generateID())

	p := &marketplace.Plugin{
		ID:          pluginID,
		Name:        pluginIn.Name,
		DisplayName: pluginIn.DisplayName,
		Description: pluginIn.Description,
		Author:      pluginIn.Author,
		Icon:        pluginIn.Icon,
		Category:    pluginIn.Category,
		Tags:        pluginIn.Tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if p.Tags == nil {
		p.Tags = []string{}
	}

	if err := s.repo.CreatePlugin(ctx, p); err != nil {
		return nil, nil, fmt.Errorf("creating plugin: %w", err)
	}

	v := &marketplace.PluginVersion{
		ID:             kernel.PluginVersionID(generateID()),
		PluginID:       pluginID,
		Version:        versionIn.Version,
		Changelog:      versionIn.Changelog,
		Permissions:    versionIn.Permissions,
		ManifestJSON:   versionIn.ManifestJSON,
		FrontendURL:    versionIn.FrontendURL,
		BackendEntry:   versionIn.BackendEntry,
		MinPlatformVer: versionIn.MinPlatformVer,
		CreatedAt:      now,
	}
	if v.Permissions == nil {
		v.Permissions = []string{}
	}

	if err := s.repo.CreateVersion(ctx, v); err != nil {
		return nil, nil, fmt.Errorf("creating version: %w", err)
	}

	return p, v, nil
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 1
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
