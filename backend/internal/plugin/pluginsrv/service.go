package pluginsrv

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/plugin"
	"github.com/google/uuid"
)

// Service handles plugin domain business logic.
type Service struct {
	pluginRepo      plugin.PluginRepository
	versionRepo     plugin.PluginVersionRepository
	installRepo     plugin.InstallationRepository
	bus             eventbus.Bus
}

// New creates a new plugin service.
func New(
	pluginRepo plugin.PluginRepository,
	versionRepo plugin.PluginVersionRepository,
	installRepo plugin.InstallationRepository,
	bus eventbus.Bus,
) *Service {
	return &Service{
		pluginRepo:  pluginRepo,
		versionRepo: versionRepo,
		installRepo: installRepo,
		bus:         bus,
	}
}

// -----------------------------------------------------------------------
// Plugin catalogue operations (global, not tenant-scoped)
// -----------------------------------------------------------------------

// ListPlugins returns a paginated list of all plugins in the global catalogue.
func (s *Service) ListPlugins(ctx context.Context, pg kernel.PaginationOptions) (kernel.Paginated[plugin.Plugin], error) {
	return s.pluginRepo.List(ctx, pg)
}

// GetPlugin retrieves a single plugin by ID.
func (s *Service) GetPlugin(ctx context.Context, id kernel.PluginID) (*plugin.Plugin, error) {
	return s.pluginRepo.GetByID(ctx, id)
}

// CreatePlugin creates a new plugin in the global catalogue (admin operation).
func (s *Service) CreatePlugin(ctx context.Context, p *plugin.Plugin) (*plugin.Plugin, error) {
	if p.Name == "" {
		return nil, plugin.ErrInvalidInput
	}
	now := time.Now()
	p.ID = kernel.PluginID(uuid.NewString())
	p.CreatedAt = now
	p.UpdatedAt = now
	if p.Tags == nil {
		p.Tags = json.RawMessage("[]")
	}
	if err := s.pluginRepo.Create(ctx, p); err != nil {
		return nil, errx.Wrap(err, "creating plugin", errx.TypeInternal)
	}
	return p, nil
}

// CreateVersion publishes a new version for an existing plugin.
func (s *Service) CreateVersion(ctx context.Context, v *plugin.PluginVersion) (*plugin.PluginVersion, error) {
	if v.PluginID.IsEmpty() || v.Version == "" {
		return nil, plugin.ErrInvalidInput
	}
	v.ID = kernel.PluginVersionID(uuid.NewString())
	v.CreatedAt = time.Now()
	if v.Permissions == nil {
		v.Permissions = json.RawMessage("[]")
	}
	if err := s.versionRepo.Create(ctx, v); err != nil {
		return nil, errx.Wrap(err, "creating plugin version", errx.TypeInternal)
	}
	return v, nil
}

// GetManifest returns the manifest JSON for the latest version of a plugin.
func (s *Service) GetManifest(ctx context.Context, pluginID kernel.PluginID) (string, error) {
	v, err := s.versionRepo.GetLatest(ctx, pluginID)
	if err != nil {
		return "", err
	}
	return v.ManifestJSON, nil
}

// -----------------------------------------------------------------------
// Installation operations (tenant-scoped)
// -----------------------------------------------------------------------

// ListInstalled returns all plugin installations for a tenant.
func (s *Service) ListInstalled(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[plugin.PluginInstallation], error) {
	return s.installRepo.ListByTenant(ctx, tenantID, pg)
}

// Install installs a specific version of a plugin for a tenant.
func (s *Service) Install(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID, versionID kernel.PluginVersionID) (*plugin.PluginInstallation, error) {
	// Check plugin exists
	p, err := s.pluginRepo.GetByID(ctx, pluginID)
	if err != nil {
		return nil, err
	}

	// Resolve version: if empty, use latest
	var v *plugin.PluginVersion
	if versionID.IsEmpty() {
		v, err = s.versionRepo.GetLatest(ctx, pluginID)
	} else {
		v, err = s.versionRepo.GetByID(ctx, versionID)
	}
	if err != nil {
		return nil, err
	}

	// Check not already installed
	existing, err := s.installRepo.GetByTenantAndPlugin(ctx, tenantID, pluginID)
	if err != nil && !errx.Is(err, plugin.ErrInstallationNotFound) {
		return nil, errx.Wrap(err, "checking existing installation", errx.TypeInternal)
	}
	if existing != nil {
		return nil, plugin.ErrAlreadyInstalled
	}

	now := time.Now()
	installation := &plugin.PluginInstallation{
		ID:          kernel.InstallationID(uuid.NewString()),
		TenantID:    tenantID,
		PluginID:    pluginID,
		VersionID:   v.ID,
		Status:      plugin.StatusActive,
		Settings:    json.RawMessage("{}"),
		InstalledAt: now,
		UpdatedAt:   now,
	}

	if err := s.installRepo.Create(ctx, installation); err != nil {
		return nil, errx.Wrap(err, "installing plugin", errx.TypeInternal)
	}

	// Publish event
	if evt, evtErr := eventbus.NewEvent(eventbus.PluginInstalled, tenantID, eventbus.PluginPayload{
		PluginID:   string(pluginID),
		PluginName: p.Name,
		Version:    v.Version,
	}); evtErr == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return installation, nil
}

// Uninstall removes a plugin installation for a tenant.
func (s *Service) Uninstall(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) error {
	// Check it exists
	installation, err := s.installRepo.GetByTenantAndPlugin(ctx, tenantID, pluginID)
	if err != nil {
		return err
	}

	p, err := s.pluginRepo.GetByID(ctx, pluginID)
	if err != nil {
		return err
	}

	v, err := s.versionRepo.GetByID(ctx, installation.VersionID)
	if err != nil {
		return err
	}

	if err := s.installRepo.Delete(ctx, tenantID, pluginID); err != nil {
		return errx.Wrap(err, "uninstalling plugin", errx.TypeInternal)
	}

	// Publish event
	if evt, evtErr := eventbus.NewEvent(eventbus.PluginUninstalled, tenantID, eventbus.PluginPayload{
		PluginID:   string(pluginID),
		PluginName: p.Name,
		Version:    v.Version,
	}); evtErr == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return nil
}

// UpdateSettings updates the settings for a plugin installation.
func (s *Service) UpdateSettings(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID, settings json.RawMessage) (*plugin.PluginInstallation, error) {
	installation, err := s.installRepo.GetByTenantAndPlugin(ctx, tenantID, pluginID)
	if err != nil {
		return nil, err
	}

	installation.Settings = settings
	installation.UpdatedAt = time.Now()

	if err := s.installRepo.Update(ctx, installation); err != nil {
		return nil, errx.Wrap(err, "updating plugin settings", errx.TypeInternal)
	}

	return installation, nil
}

// Enable activates a previously disabled plugin installation.
func (s *Service) Enable(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*plugin.PluginInstallation, error) {
	return s.setStatus(ctx, tenantID, pluginID, plugin.StatusActive)
}

// Disable deactivates a plugin installation.
func (s *Service) Disable(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID) (*plugin.PluginInstallation, error) {
	return s.setStatus(ctx, tenantID, pluginID, plugin.StatusInactive)
}

func (s *Service) setStatus(ctx context.Context, tenantID kernel.TenantID, pluginID kernel.PluginID, status plugin.InstallationStatus) (*plugin.PluginInstallation, error) {
	installation, err := s.installRepo.GetByTenantAndPlugin(ctx, tenantID, pluginID)
	if err != nil {
		return nil, err
	}

	installation.Status = status
	installation.UpdatedAt = time.Now()

	if err := s.installRepo.Update(ctx, installation); err != nil {
		return nil, errx.Wrap(err, "updating plugin status", errx.TypeInternal)
	}

	return installation, nil
}
