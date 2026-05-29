package catalogsrv

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/catalog"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
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
		return nil, fmt.Errorf("checking slug uniqueness: %w", err)
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
		return nil, fmt.Errorf("creating category: %w", err)
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
func (s *Service) ListCategories(ctx context.Context, tenantID kernel.TenantID, pg kernel.Pagination) (kernel.PaginatedResult[catalog.Category], error) {
	return s.categories.List(ctx, tenantID, pg)
}

// ListCategoriesByParent returns child categories of a parent (nil for root categories).
func (s *Service) ListCategoriesByParent(ctx context.Context, tenantID kernel.TenantID, parentID *kernel.CategoryID, pg kernel.Pagination) (kernel.PaginatedResult[catalog.Category], error) {
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
		return nil, fmt.Errorf("checking slug uniqueness: %w", err)
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
		return nil, fmt.Errorf("creating collection: %w", err)
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
func (s *Service) ListCollections(ctx context.Context, tenantID kernel.TenantID, pg kernel.Pagination) (kernel.PaginatedResult[catalog.Collection], error) {
	return s.collections.List(ctx, tenantID, pg)
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
