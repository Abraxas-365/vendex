package marketplacesrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/marketplace"
)

// PresetService implements preset business logic.
type PresetService struct {
	presetRepo  marketplace.PresetRepository
	installRepo marketplace.PresetInstallRepository
}

// NewPresetService creates a new PresetService.
func NewPresetService(
	presetRepo marketplace.PresetRepository,
	installRepo marketplace.PresetInstallRepository,
) *PresetService {
	return &PresetService{
		presetRepo:  presetRepo,
		installRepo: installRepo,
	}
}

// Create creates a new preset for a tenant.
func (s *PresetService) Create(ctx context.Context, tenantID kernel.TenantID, req marketplace.CreatePresetRequest) (marketplace.Preset, error) {
	now := time.Now()
	p := marketplace.Preset{
		ID:            kernel.PresetID(uuid.New().String()),
		TenantID:      tenantID,
		Name:          req.Name,
		Slug:          req.Slug,
		Description:   req.Description,
		Version:       req.Version,
		Image:         req.Image,
		FrontendPort:  req.FrontendPort,
		SystemPrompt:  req.SystemPrompt,
		ToolsManifest: req.ToolsManifest,
		Status:        marketplace.PresetStatusDraft,
		Visibility:    req.Visibility,
		Icon:          req.Icon,
		Tags:          req.Tags,
		InstallCount:  0,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if p.Visibility == "" {
		p.Visibility = marketplace.PresetVisibilityPrivate
	}
	if p.Version == "" {
		p.Version = "1.0.0"
	}
	if p.FrontendPort == 0 {
		p.FrontendPort = 8080
	}
	return s.presetRepo.Create(ctx, p)
}

// Get retrieves a preset by ID.
func (s *PresetService) Get(ctx context.Context, id kernel.PresetID) (marketplace.Preset, error) {
	return s.presetRepo.GetByID(ctx, id)
}

// GetBySlug retrieves a preset by its slug.
func (s *PresetService) GetBySlug(ctx context.Context, slug string) (marketplace.Preset, error) {
	return s.presetRepo.GetBySlug(ctx, slug)
}

// Update applies partial updates to a preset.
func (s *PresetService) Update(ctx context.Context, tenantID kernel.TenantID, id kernel.PresetID, req marketplace.UpdatePresetRequest) (marketplace.Preset, error) {
	p, err := s.presetRepo.GetByID(ctx, id)
	if err != nil {
		return marketplace.Preset{}, err
	}

	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Version != nil {
		p.Version = *req.Version
	}
	if req.Image != nil {
		p.Image = *req.Image
	}
	if req.FrontendPort != nil {
		p.FrontendPort = *req.FrontendPort
	}
	if req.SystemPrompt != nil {
		p.SystemPrompt = *req.SystemPrompt
	}
	if req.ToolsManifest != nil {
		p.ToolsManifest = *req.ToolsManifest
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if req.Visibility != nil {
		p.Visibility = *req.Visibility
	}
	if req.Icon != nil {
		p.Icon = *req.Icon
	}
	if req.Tags != nil {
		p.Tags = *req.Tags
	}
	p.UpdatedAt = time.Now()

	return s.presetRepo.Update(ctx, p)
}

// Publish transitions a preset to published status.
func (s *PresetService) Publish(ctx context.Context, id kernel.PresetID) (marketplace.Preset, error) {
	p, err := s.presetRepo.GetByID(ctx, id)
	if err != nil {
		return marketplace.Preset{}, err
	}
	p.Status = marketplace.PresetStatusPublished
	p.UpdatedAt = time.Now()
	return s.presetRepo.Update(ctx, p)
}

// Archive transitions a preset to archived status.
func (s *PresetService) Archive(ctx context.Context, id kernel.PresetID) (marketplace.Preset, error) {
	p, err := s.presetRepo.GetByID(ctx, id)
	if err != nil {
		return marketplace.Preset{}, err
	}
	p.Status = marketplace.PresetStatusArchived
	p.UpdatedAt = time.Now()
	return s.presetRepo.Update(ctx, p)
}

// Delete removes a preset.
func (s *PresetService) Delete(ctx context.Context, id kernel.PresetID) error {
	return s.presetRepo.Delete(ctx, id)
}

// ListMarketplace returns paginated public/published presets for marketplace browsing.
func (s *PresetService) ListMarketplace(ctx context.Context, opts marketplace.PresetListOptions) (kernel.Paginated[marketplace.Preset], error) {
	published := marketplace.PresetStatusPublished
	public := marketplace.PresetVisibilityPublic
	opts.Status = &published
	opts.Visibility = &public
	return s.presetRepo.List(ctx, opts)
}

// ListByTenant returns paginated presets owned by a specific tenant.
func (s *PresetService) ListByTenant(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.Preset], error) {
	return s.presetRepo.ListByTenant(ctx, tenantID, p)
}

// Install installs a preset for a tenant.
func (s *PresetService) Install(ctx context.Context, tenantID kernel.TenantID, presetID kernel.PresetID, config []byte) (marketplace.PresetInstall, error) {
	install := marketplace.PresetInstall{
		ID:          uuid.New().String(),
		TenantID:    tenantID,
		PresetID:    presetID,
		InstalledAt: time.Now(),
		Config:      config,
	}
	return s.installRepo.Install(ctx, install)
}

// Uninstall removes a preset installation for a tenant.
func (s *PresetService) Uninstall(ctx context.Context, tenantID kernel.TenantID, presetID kernel.PresetID) error {
	return s.installRepo.Uninstall(ctx, tenantID, presetID)
}

// ListInstalled returns paginated preset installations for a tenant.
func (s *PresetService) ListInstalled(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.PresetInstall], error) {
	return s.installRepo.ListByTenant(ctx, tenantID, p)
}
