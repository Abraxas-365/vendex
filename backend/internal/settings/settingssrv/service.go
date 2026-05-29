package settingssrv

import (
	"context"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/settings"
)

// Service handles store-settings business logic.
type Service struct {
	repo settings.Repository
}

// New creates a new settings service.
func New(repo settings.Repository) *Service {
	return &Service{repo: repo}
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
			return nil, fmt.Errorf("creating default settings: %w", upsertErr)
		}
		return defaults, nil
	}
	if err != nil {
		return nil, fmt.Errorf("getting settings: %w", err)
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
		return nil, fmt.Errorf("upserting settings: %w", err)
	}
	return ss, nil
}
