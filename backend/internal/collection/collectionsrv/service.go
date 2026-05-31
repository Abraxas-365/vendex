package collectionsrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/collection"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Service implements all business logic for the collection domain.
type Service struct {
	repo collection.Repository
	bus  eventbus.Bus
}

// New creates a new Service.
func New(repo collection.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// --------------------------------------------------------------------------
// Collection CRUD
// --------------------------------------------------------------------------

// Create validates the input and persists a new collection.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, in collection.CreateInput) (*collection.Collection, error) {
	if in.Name == "" {
		return nil, collection.ErrNameRequired
	}
	if in.Slug == "" {
		return nil, collection.ErrSlugRequired
	}

	colType := in.Type
	if colType == "" {
		colType = collection.CollectionManual
	}
	if colType != collection.CollectionManual && colType != collection.CollectionAuto {
		return nil, collection.ErrInvalidType
	}

	// Slug uniqueness check.
	existing, err := s.repo.GetBySlug(ctx, tenantID, in.Slug)
	if err != nil && !errx.Is(err, collection.ErrNotFound) {
		return nil, errx.Wrap(err, "checking slug uniqueness", errx.TypeInternal)
	}
	if existing != nil {
		return nil, collection.ErrDuplicateSlug
	}

	rules := in.Rules
	if rules == nil {
		rules = []collection.CollectionRule{}
	}

	now := time.Now().UTC()
	c := &collection.Collection{
		ID:              kernel.CollectionID(uuid.NewString()),
		TenantID:        tenantID,
		Name:            in.Name,
		Slug:            in.Slug,
		Description:     in.Description,
		ImageURL:        in.ImageURL,
		Type:            colType,
		Rules:           rules,
		IsActive:        in.IsActive,
		SortOrder:       in.SortOrder,
		MetaTitle:       in.MetaTitle,
		MetaDescription: in.MetaDescription,
		PublishedAt:     in.PublishedAt,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, eventbus.CollectionCreated, tenantID, collectionPayload(c))

	return c, nil
}

// GetByID retrieves a collection by its ID and populates ProductCount.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) (*collection.Collection, error) {
	c, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}
	c.ProductCount, _ = s.repo.CountProducts(ctx, tenantID, id)
	return c, nil
}

// GetBySlug retrieves a collection by its URL slug and populates ProductCount.
func (s *Service) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*collection.Collection, error) {
	c, err := s.repo.GetBySlug(ctx, tenantID, slug)
	if err != nil {
		return nil, err
	}
	c.ProductCount, _ = s.repo.CountProducts(ctx, tenantID, c.ID)
	return c, nil
}

// List returns a paginated list of collections.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, activeOnly bool, page, pageSize int) (kernel.Paginated[collection.Collection], error) {
	pg := kernel.NewPaginationOptions(page, pageSize)
	result, err := s.repo.List(ctx, tenantID, activeOnly, pg)
	if err != nil {
		return kernel.Paginated[collection.Collection]{}, err
	}
	return result, nil
}

// Update applies partial changes to a collection and persists them.
func (s *Service) Update(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID, in collection.UpdateInput) (*collection.Collection, error) {
	c, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		if *in.Name == "" {
			return nil, collection.ErrNameRequired
		}
		c.Name = *in.Name
	}
	if in.Slug != nil {
		if *in.Slug == "" {
			return nil, collection.ErrSlugRequired
		}
		if *in.Slug != string(c.Slug) {
			// Slug uniqueness check.
			existing, err := s.repo.GetBySlug(ctx, tenantID, *in.Slug)
			if err != nil && !errx.Is(err, collection.ErrNotFound) {
				return nil, errx.Wrap(err, "checking slug uniqueness", errx.TypeInternal)
			}
			if existing != nil {
				return nil, collection.ErrDuplicateSlug
			}
		}
		c.Slug = *in.Slug
	}
	if in.Description != nil {
		c.Description = *in.Description
	}
	if in.ImageURL != nil {
		c.ImageURL = *in.ImageURL
	}
	if in.Type != nil {
		t := *in.Type
		if t != collection.CollectionManual && t != collection.CollectionAuto {
			return nil, collection.ErrInvalidType
		}
		c.Type = t
	}
	if in.Rules != nil {
		c.Rules = in.Rules
	}
	if in.IsActive != nil {
		c.IsActive = *in.IsActive
	}
	if in.SortOrder != nil {
		c.SortOrder = *in.SortOrder
	}
	if in.MetaTitle != nil {
		c.MetaTitle = *in.MetaTitle
	}
	if in.MetaDescription != nil {
		c.MetaDescription = *in.MetaDescription
	}
	if in.ClearPublishedAt {
		c.PublishedAt = nil
	} else if in.PublishedAt != nil {
		c.PublishedAt = in.PublishedAt
	}

	c.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}

	s.publishEvent(ctx, eventbus.CollectionUpdated, tenantID, collectionPayload(c))

	return c, nil
}

// Delete removes a collection and all its product memberships (via FK cascade).
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// --------------------------------------------------------------------------
// Product membership
// --------------------------------------------------------------------------

// AddProduct adds a product to a manual collection.
func (s *Service) AddProduct(ctx context.Context, tenantID kernel.TenantID, collectionID, productID string, sortOrder int) (*collection.CollectionProduct, error) {
	cp := &collection.CollectionProduct{
		ID:           kernel.CollectionProductID(uuid.NewString()),
		TenantID:     tenantID,
		CollectionID: kernel.CollectionID(collectionID),
		ProductID:    productID,
		SortOrder:    sortOrder,
		AddedAt:      time.Now().UTC(),
	}

	if err := s.repo.AddProduct(ctx, cp); err != nil {
		return nil, err
	}
	return cp, nil
}

// RemoveProduct removes a product from a collection.
func (s *Service) RemoveProduct(ctx context.Context, tenantID kernel.TenantID, collectionID, productID string) error {
	return s.repo.RemoveProduct(ctx, tenantID, kernel.CollectionID(collectionID), productID)
}

// ListProducts returns a paginated list of products in a collection.
func (s *Service) ListProducts(ctx context.Context, tenantID kernel.TenantID, collectionID string, page, pageSize int) (kernel.Paginated[collection.CollectionProduct], error) {
	pg := kernel.NewPaginationOptions(page, pageSize)
	return s.repo.ListProducts(ctx, tenantID, kernel.CollectionID(collectionID), pg)
}

// ReorderProducts updates the sort_order for products in the given order.
func (s *Service) ReorderProducts(ctx context.Context, tenantID kernel.TenantID, collectionID string, productIDs []string) error {
	return s.repo.ReorderProducts(ctx, tenantID, kernel.CollectionID(collectionID), productIDs)
}

// --------------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------------

func collectionPayload(c *collection.Collection) eventbus.CollectionPayload {
	return eventbus.CollectionPayload{
		CollectionID: string(c.ID),
		Name:         c.Name,
		Slug:         c.Slug,
	}
}

func (s *Service) publishEvent(ctx context.Context, evtType eventbus.EventType, tenantID kernel.TenantID, payload any) {
	if evt, err := eventbus.NewEvent(evtType, tenantID, payload); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}
}
