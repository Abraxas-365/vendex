package settingssrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/settings"
)

// Service handles store-settings business logic.
type Service struct {
	repo settings.Repository
	bus  eventbus.Bus
}

// New creates a new settings service.
func New(repo settings.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// UpdateInput holds the fields that callers may change.
type UpdateInput struct {
	StoreName      string
	StoreEmail     string
	StorePhone     string
	Currency       string
	Timezone       string
	Address        settings.StoreAddress
	LogoURL        string
	FaviconURL     string
	SocialLinks    settings.SocialLinks
	CheckoutConfig settings.CheckoutConfig
}

// Get returns the settings for a tenant, creating defaults if none exist yet.
func (s *Service) Get(ctx context.Context, tenantID kernel.TenantID) (*settings.StoreSettings, error) {
	ss, err := s.repo.Get(ctx, tenantID)
	if errx.IsNotFound(err) {
		defaults := settings.DefaultSettings(tenantID)
		if upsertErr := s.repo.Upsert(ctx, defaults); upsertErr != nil {
			return nil, errx.Wrap(upsertErr, "creating default settings", errx.TypeInternal)
		}
		return defaults, nil
	}
	if err != nil {
		return nil, errx.Wrap(err, "getting settings", errx.TypeInternal)
	}
	return ss, nil
}

// Update applies the given input to the tenant's settings and persists the result.
// If no settings row exists yet, defaults are created first.
func (s *Service) Update(ctx context.Context, tenantID kernel.TenantID, in UpdateInput) (*settings.StoreSettings, error) {
	ss, err := s.Get(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	ss.StoreName = in.StoreName
	ss.StoreEmail = in.StoreEmail
	ss.StorePhone = in.StorePhone
	ss.Currency = in.Currency
	ss.Timezone = in.Timezone
	ss.Address = in.Address
	ss.LogoURL = in.LogoURL
	ss.FaviconURL = in.FaviconURL
	ss.SocialLinks = in.SocialLinks
	ss.CheckoutConfig = in.CheckoutConfig
	ss.UpdatedAt = time.Now()

	if err := s.repo.Upsert(ctx, ss); err != nil {
		return nil, errx.Wrap(err, "upserting settings", errx.TypeInternal)
	}

	// Build a list of fields that were updated.
	var changedFields []string
	if in.StoreName != "" {
		changedFields = append(changedFields, "store_name")
	}
	if in.StoreEmail != "" {
		changedFields = append(changedFields, "store_email")
	}
	if in.StorePhone != "" {
		changedFields = append(changedFields, "store_phone")
	}
	if in.Currency != "" {
		changedFields = append(changedFields, "currency")
	}
	if in.Timezone != "" {
		changedFields = append(changedFields, "timezone")
	}
	if in.LogoURL != "" {
		changedFields = append(changedFields, "logo_url")
	}
	if in.FaviconURL != "" {
		changedFields = append(changedFields, "favicon_url")
	}

	if evt, err := eventbus.NewEvent(eventbus.SettingsUpdated, tenantID, eventbus.SettingsPayload{
		Fields: changedFields,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return ss, nil
}
