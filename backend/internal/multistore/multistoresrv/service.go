package multistoresrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/multistore"
	"github.com/google/uuid"
)

// Service implements multistore business logic.
type Service struct {
	repo multistore.Repository
	bus  eventbus.Bus
}

// New creates a new multistore Service.
func New(repo multistore.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// Create creates a new storefront for a tenant.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, input multistore.CreateInput) (*multistore.Storefront, error) {
	if input.Name == "" {
		return nil, multistore.ErrNameRequired
	}
	if input.Slug == "" {
		return nil, multistore.ErrSlugRequired
	}

	locale := input.DefaultLocale
	if locale == "" {
		locale = "en"
	}
	currency := input.DefaultCurrency
	if currency == "" {
		currency = "USD"
	}
	settings := input.Settings
	if settings == nil {
		settings = map[string]interface{}{}
	}

	now := time.Now().UTC()
	sf := &multistore.Storefront{
		ID:              kernel.NewStorefrontEntryID(uuid.NewString()),
		TenantID:        tenantID,
		Name:            input.Name,
		Slug:            input.Slug,
		Domain:          input.Domain,
		Description:     input.Description,
		ThemeID:         input.ThemeID,
		LogoURL:         input.LogoURL,
		DefaultLocale:   locale,
		DefaultCurrency: currency,
		IsActive:        input.IsActive,
		IsDefault:       false,
		Settings:        settings,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Create(ctx, sf); err != nil {
		return nil, err
	}

	event, err := eventbus.NewEvent(eventbus.StorefrontCreated, tenantID, eventbus.StorefrontPayload{
		StorefrontID: string(sf.ID),
		Name:         sf.Name,
		Slug:         sf.Slug,
	})
	if err == nil {
		_ = s.bus.Publish(ctx, event)
	}

	return sf, nil
}

// GetByID returns a storefront by ID, scoped to the tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID) (*multistore.Storefront, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// GetBySlug returns a storefront by slug, scoped to the tenant.
func (s *Service) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*multistore.Storefront, error) {
	if slug == "" {
		return nil, multistore.ErrSlugRequired
	}
	return s.repo.GetBySlug(ctx, tenantID, slug)
}

// GetByDomain returns a storefront by its custom domain (global lookup).
func (s *Service) GetByDomain(ctx context.Context, domain string) (*multistore.Storefront, error) {
	if domain == "" {
		return nil, multistore.ErrNotFound
	}
	return s.repo.GetByDomain(ctx, domain)
}

// List returns paginated storefronts for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[multistore.Storefront], error) {
	return s.repo.List(ctx, tenantID, page, pageSize)
}

// Update applies partial updates to a storefront.
func (s *Service) Update(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID, input multistore.UpdateInput) (*multistore.Storefront, error) {
	sf, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if input.Name != nil {
		if *input.Name == "" {
			return nil, multistore.ErrNameRequired
		}
		sf.Name = *input.Name
	}
	if input.Domain != nil {
		sf.Domain = input.Domain
	}
	if input.Description != nil {
		sf.Description = *input.Description
	}
	if input.ThemeID != nil {
		sf.ThemeID = *input.ThemeID
	}
	if input.LogoURL != nil {
		sf.LogoURL = *input.LogoURL
	}
	if input.DefaultLocale != nil {
		sf.DefaultLocale = *input.DefaultLocale
	}
	if input.DefaultCurrency != nil {
		sf.DefaultCurrency = *input.DefaultCurrency
	}
	if input.IsActive != nil {
		sf.IsActive = *input.IsActive
	}
	if input.Settings != nil {
		sf.Settings = input.Settings
	}

	sf.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, sf); err != nil {
		return nil, err
	}

	event, err := eventbus.NewEvent(eventbus.StorefrontUpdated, tenantID, eventbus.StorefrontPayload{
		StorefrontID: string(sf.ID),
		Name:         sf.Name,
		Slug:         sf.Slug,
	})
	if err == nil {
		_ = s.bus.Publish(ctx, event)
	}

	return sf, nil
}

// Delete removes a storefront. The default storefront cannot be deleted.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID) error {
	sf, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	if sf.IsDefault {
		return multistore.ErrDeleteDefault
	}

	if err := s.repo.Delete(ctx, tenantID, id); err != nil {
		return err
	}

	event, err := eventbus.NewEvent(eventbus.StorefrontDeleted, tenantID, eventbus.StorefrontPayload{
		StorefrontID: string(sf.ID),
		Name:         sf.Name,
		Slug:         sf.Slug,
	})
	if err == nil {
		_ = s.bus.Publish(ctx, event)
	}

	return nil
}

// SetDefault marks the given storefront as the default and clears all others.
func (s *Service) SetDefault(ctx context.Context, tenantID kernel.TenantID, id kernel.StorefrontEntryID) error {
	// Verify the storefront exists and belongs to the tenant.
	_, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return s.repo.SetDefault(ctx, tenantID, id)
}

// AddCatalog links a catalog to a storefront.
func (s *Service) AddCatalog(ctx context.Context, tenantID kernel.TenantID, storefrontID kernel.StorefrontEntryID, catalogID string, sortOrder int) (*multistore.StorefrontCatalog, error) {
	// Verify the storefront belongs to the tenant.
	if _, err := s.repo.GetByID(ctx, tenantID, storefrontID); err != nil {
		return nil, err
	}

	sc := &multistore.StorefrontCatalog{
		ID:           kernel.NewStorefrontCatalogID(uuid.NewString()),
		TenantID:     tenantID,
		StorefrontID: storefrontID,
		CatalogID:    catalogID,
		SortOrder:    sortOrder,
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.repo.AddCatalog(ctx, sc); err != nil {
		return nil, err
	}

	return sc, nil
}

// RemoveCatalog removes a catalog link from a storefront.
func (s *Service) RemoveCatalog(ctx context.Context, tenantID kernel.TenantID, storefrontID kernel.StorefrontEntryID, catalogID string) error {
	return s.repo.RemoveCatalog(ctx, tenantID, storefrontID, catalogID)
}

// ListCatalogs returns all catalog links for a storefront.
func (s *Service) ListCatalogs(ctx context.Context, tenantID kernel.TenantID, storefrontID kernel.StorefrontEntryID) ([]multistore.StorefrontCatalog, error) {
	// Verify the storefront belongs to the tenant.
	if _, err := s.repo.GetByID(ctx, tenantID, storefrontID); err != nil {
		return nil, err
	}
	return s.repo.ListCatalogs(ctx, tenantID, storefrontID)
}
