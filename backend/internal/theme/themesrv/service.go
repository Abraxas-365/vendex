package themesrv

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/theme"
)

// Service implements the theme business logic.
type Service struct {
	repo theme.ThemeRepository
}

// New creates a new theme Service.
func New(repo theme.ThemeRepository) *Service {
	return &Service{repo: repo}
}

// CreateThemeInput holds the data needed to create a new theme.
type CreateThemeInput struct {
	TenantID kernel.TenantID
	Name     string
	Tokens   *theme.ThemeTokens
}

// UpdateThemeInput holds the fields that can change on a theme update.
type UpdateThemeInput struct {
	TenantID kernel.TenantID
	ID       kernel.ThemeID
	Name     *string
	Tokens   *theme.ThemeTokens
}

// CreateTheme creates a new theme with DefaultTokens merged with any provided tokens.
func (s *Service) CreateTheme(ctx context.Context, input CreateThemeInput) (*theme.Theme, error) {
	if input.Name == "" {
		return nil, errx.Validation("theme name is required")
	}

	tokens := theme.DefaultTokens()
	if input.Tokens != nil {
		tokens = mergeTokens(tokens, *input.Tokens)
	}

	now := time.Now().UTC()
	t := &theme.Theme{
		ID:        kernel.ThemeID(generateID()),
		TenantID:  input.TenantID,
		Name:      input.Name,
		IsActive:  false,
		Tokens:    tokens,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// GetTheme retrieves a theme by ID for a tenant.
func (s *Service) GetTheme(ctx context.Context, tenantID kernel.TenantID, id kernel.ThemeID) (*theme.Theme, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// GetActiveTheme returns the active theme for a tenant.
// If no active theme exists, it creates a default one and activates it.
func (s *Service) GetActiveTheme(ctx context.Context, tenantID kernel.TenantID) (*theme.Theme, error) {
	t, err := s.repo.GetActive(ctx, tenantID)
	if err == nil {
		return t, nil
	}
	if !errx.IsNotFound(err) {
		return nil, err
	}

	// No active theme — create a default one and activate it.
	now := time.Now().UTC()
	defaultTheme := &theme.Theme{
		ID:        kernel.ThemeID(generateID()),
		TenantID:  tenantID,
		Name:      "Default",
		IsActive:  true,
		Tokens:    theme.DefaultTokens(),
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, defaultTheme); err != nil {
		return nil, errx.Wrap(err, "creating default theme", errx.TypeInternal)
	}
	return defaultTheme, nil
}

// ListThemes returns all themes for a tenant.
func (s *Service) ListThemes(ctx context.Context, tenantID kernel.TenantID) ([]theme.Theme, error) {
	return s.repo.List(ctx, tenantID)
}

// UpdateTheme updates the name and/or tokens of an existing theme.
func (s *Service) UpdateTheme(ctx context.Context, input UpdateThemeInput) (*theme.Theme, error) {
	t, err := s.repo.GetByID(ctx, input.TenantID, input.ID)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		if *input.Name == "" {
			return nil, errx.Validation("theme name cannot be empty")
		}
		t.Name = *input.Name
	}
	if input.Tokens != nil {
		t.Tokens = *input.Tokens
	}
	t.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	return t, nil
}

// ActivateTheme deactivates all other themes for the tenant, then activates the given one.
func (s *Service) ActivateTheme(ctx context.Context, tenantID kernel.TenantID, id kernel.ThemeID) (*theme.Theme, error) {
	t, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Deactivate all themes for the tenant.
	if err := s.repo.DeactivateAll(ctx, tenantID); err != nil {
		return nil, errx.Wrap(err, "deactivating themes", errx.TypeInternal)
	}

	t.IsActive = true
	t.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, t); err != nil {
		return nil, errx.Wrap(err, "activating theme", errx.TypeInternal)
	}
	return t, nil
}

// DeleteTheme deletes a theme. The active theme cannot be deleted.
func (s *Service) DeleteTheme(ctx context.Context, tenantID kernel.TenantID, id kernel.ThemeID) error {
	t, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if t.IsActive {
		return theme.ErrThemeActiveDelete
	}
	return s.repo.Delete(ctx, tenantID, id)
}

// DuplicateTheme copies an existing theme under a new name.
func (s *Service) DuplicateTheme(ctx context.Context, tenantID kernel.TenantID, id kernel.ThemeID, newName string) (*theme.Theme, error) {
	if newName == "" {
		return nil, errx.Validation("new theme name is required")
	}

	source, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	duplicate := &theme.Theme{
		ID:        kernel.ThemeID(generateID()),
		TenantID:  tenantID,
		Name:      newName,
		IsActive:  false,
		Tokens:    source.Tokens,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.repo.Create(ctx, duplicate); err != nil {
		return nil, err
	}
	return duplicate, nil
}

// mergeTokens overlays non-zero fields from override onto base.
// Only string fields that are non-empty and float64 fields that are non-zero override the base.
func mergeTokens(base, override theme.ThemeTokens) theme.ThemeTokens {
	merged := base

	// Colors
	if override.Colors.Primary != "" {
		merged.Colors.Primary = override.Colors.Primary
	}
	if override.Colors.Secondary != "" {
		merged.Colors.Secondary = override.Colors.Secondary
	}
	if override.Colors.Background != "" {
		merged.Colors.Background = override.Colors.Background
	}
	if override.Colors.Surface != "" {
		merged.Colors.Surface = override.Colors.Surface
	}
	if override.Colors.Text != "" {
		merged.Colors.Text = override.Colors.Text
	}
	if override.Colors.TextMuted != "" {
		merged.Colors.TextMuted = override.Colors.TextMuted
	}
	if override.Colors.Border != "" {
		merged.Colors.Border = override.Colors.Border
	}
	if override.Colors.Error != "" {
		merged.Colors.Error = override.Colors.Error
	}
	if override.Colors.Success != "" {
		merged.Colors.Success = override.Colors.Success
	}
	if override.Colors.Warning != "" {
		merged.Colors.Warning = override.Colors.Warning
	}
	if override.Colors.Info != "" {
		merged.Colors.Info = override.Colors.Info
	}

	// Typography
	if override.Typography.FontHeading != "" {
		merged.Typography.FontHeading = override.Typography.FontHeading
	}
	if override.Typography.FontBody != "" {
		merged.Typography.FontBody = override.Typography.FontBody
	}
	if override.Typography.BaseSize != "" {
		merged.Typography.BaseSize = override.Typography.BaseSize
	}
	if override.Typography.ScaleRatio != 0 {
		merged.Typography.ScaleRatio = override.Typography.ScaleRatio
	}

	// Spacing
	if override.Spacing.Unit != "" {
		merged.Spacing.Unit = override.Spacing.Unit
	}
	if override.Spacing.SectionPadding != "" {
		merged.Spacing.SectionPadding = override.Spacing.SectionPadding
	}

	// Borders
	if override.Borders.RadiusSm != "" {
		merged.Borders.RadiusSm = override.Borders.RadiusSm
	}
	if override.Borders.RadiusMd != "" {
		merged.Borders.RadiusMd = override.Borders.RadiusMd
	}
	if override.Borders.RadiusLg != "" {
		merged.Borders.RadiusLg = override.Borders.RadiusLg
	}
	if override.Borders.RadiusFull != "" {
		merged.Borders.RadiusFull = override.Borders.RadiusFull
	}

	// Shadows
	if override.Shadows.Sm != "" {
		merged.Shadows.Sm = override.Shadows.Sm
	}
	if override.Shadows.Md != "" {
		merged.Shadows.Md = override.Shadows.Md
	}
	if override.Shadows.Lg != "" {
		merged.Shadows.Lg = override.Shadows.Lg
	}

	return merged
}

// generateID creates a new UUID string.
func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
