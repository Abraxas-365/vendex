package catalogsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/catalog"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Service handles catalog business logic for both categories and collections.
type Service struct {
	categories  catalog.CategoryRepository
	collections catalog.CollectionRepository
}

// New creates a new catalog service.
func New(categories catalog.CategoryRepository, collections catalog.CollectionRepository) *Service {
	return &Service{categories: categories, collections: collections}
}

// --- Category operations ---

// CreateCategoryInput holds the data needed to create a category.
type CreateCategoryInput struct {
	Name        string
	Slug        string
	ParentID    *kernel.CategoryID
	Description string
}

// CreateCategory creates a new category for the given tenant.
func (s *Service) CreateCategory(ctx context.Context, tenantID kernel.TenantID, in CreateCategoryInput) (*catalog.Category, error) {
	// Check for duplicate slug.
	existing, err := s.categories.GetBySlug(ctx, tenantID, in.Slug)
	if err != nil && !errx.Is(err, catalog.ErrCategoryNotFound) {
		return nil, errx.Wrap(err, "checking slug uniqueness", errx.TypeInternal)
	}
	if existing != nil {
		return nil, catalog.ErrCategoryDuplicateSlug
	}

	now := time.Now()
	c := &catalog.Category{
		ID:          kernel.CategoryID(generateID()),
		TenantID:    tenantID,
		Name:        in.Name,
		Slug:        in.Slug,
		ParentID:    in.ParentID,
		Description: in.Description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.categories.Create(ctx, c); err != nil {
		return nil, errx.Wrap(err, "creating category", errx.TypeInternal)
	}
	return c, nil
}

// GetCategoryByID retrieves a category by ID, scoped to tenant.
func (s *Service) GetCategoryByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CategoryID) (*catalog.Category, error) {
	return s.categories.GetByID(ctx, tenantID, id)
}

// UpdateCategory persists changes to a category.
func (s *Service) UpdateCategory(ctx context.Context, c *catalog.Category) error {
	c.UpdatedAt = time.Now()
	return s.categories.Update(ctx, c)
}

// DeleteCategory removes a category by ID, scoped to tenant.
func (s *Service) DeleteCategory(ctx context.Context, tenantID kernel.TenantID, id kernel.CategoryID) error {
	return s.categories.Delete(ctx, tenantID, id)
}

// ListCategories returns a paginated list of categories for a tenant.
func (s *Service) ListCategories(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[catalog.Category], error) {
	return s.categories.List(ctx, tenantID, pg)
}

// ListCategoriesByParent returns child categories of a parent (nil for root categories).
func (s *Service) ListCategoriesByParent(ctx context.Context, tenantID kernel.TenantID, parentID *kernel.CategoryID, pg kernel.PaginationOptions) (kernel.Paginated[catalog.Category], error) {
	return s.categories.ListByParent(ctx, tenantID, parentID, pg)
}

// --- Collection operations ---

// CreateCollectionInput holds the data needed to create a collection.
type CreateCollectionInput struct {
	Name        string
	Slug        string
	Description string
	ProductIDs  []kernel.ProductID
	IsAutomatic bool
	Rules       map[string]any
}

// CreateCollection creates a new collection for the given tenant.
func (s *Service) CreateCollection(ctx context.Context, tenantID kernel.TenantID, in CreateCollectionInput) (*catalog.Collection, error) {
	// Check for duplicate slug.
	existing, err := s.collections.GetBySlug(ctx, tenantID, in.Slug)
	if err != nil && !errx.Is(err, catalog.ErrCollectionNotFound) {
		return nil, errx.Wrap(err, "checking slug uniqueness", errx.TypeInternal)
	}
	if existing != nil {
		return nil, catalog.ErrCollectionDupSlug
	}

	now := time.Now()
	productIDs := in.ProductIDs
	if productIDs == nil {
		productIDs = []kernel.ProductID{}
	}
	rules := in.Rules
	if rules == nil {
		rules = map[string]any{}
	}

	c := &catalog.Collection{
		ID:          kernel.CollectionID(generateID()),
		TenantID:    tenantID,
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		ProductIDs:  productIDs,
		IsAutomatic: in.IsAutomatic,
		Rules:       rules,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.collections.Create(ctx, c); err != nil {
		return nil, errx.Wrap(err, "creating collection", errx.TypeInternal)
	}
	return c, nil
}

// GetCollectionByID retrieves a collection by ID, scoped to tenant.
func (s *Service) GetCollectionByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) (*catalog.Collection, error) {
	return s.collections.GetByID(ctx, tenantID, id)
}

// UpdateCollection persists changes to a collection.
func (s *Service) UpdateCollection(ctx context.Context, c *catalog.Collection) error {
	c.UpdatedAt = time.Now()
	return s.collections.Update(ctx, c)
}

// DeleteCollection removes a collection by ID, scoped to tenant.
func (s *Service) DeleteCollection(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) error {
	return s.collections.Delete(ctx, tenantID, id)
}

// ListCollections returns a paginated list of collections for a tenant.
func (s *Service) ListCollections(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[catalog.Collection], error) {
	return s.collections.List(ctx, tenantID, pg)
}


