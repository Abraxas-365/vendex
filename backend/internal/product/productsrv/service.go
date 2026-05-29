package productsrv

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
	"github.com/Abraxas-365/hada-commerce/internal/product"
)

// Service handles product business logic.
type Service struct {
	repo product.Repository
}

// New creates a new product service.
func New(repo product.Repository) *Service {
	return &Service{repo: repo}
}

// CreateInput holds the data needed to create a product.
type CreateInput struct {
	Name        string
	Description string
	Price       kernel.Money
	SKU         string
	Images      []string
	CategoryID  kernel.CategoryID
	Tags        []string
	Stock       int
}

// Create creates a new product for the given tenant.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, in CreateInput) (*product.Product, error) {
	if in.Price.Amount <= 0 {
		return nil, product.ErrInvalidPrice
	}

	// Check for duplicate SKU.
	if in.SKU != "" {
		existing, err := s.repo.GetBySKU(ctx, tenantID, in.SKU)
		if err != nil && !errx.Is(err, product.ErrNotFound) {
			return nil, fmt.Errorf("checking SKU uniqueness: %w", err)
		}
		if existing != nil {
			return nil, product.ErrDuplicateSKU
		}
	}

	now := time.Now()
	p := &product.Product{
		ID:          kernel.ProductID(generateID()),
		TenantID:    tenantID,
		Name:        in.Name,
		Description: in.Description,
		Price:       in.Price,
		SKU:         in.SKU,
		Images:      in.Images,
		CategoryID:  in.CategoryID,
		Tags:        in.Tags,
		Status:      product.StatusDraft,
		Stock:       in.Stock,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, fmt.Errorf("creating product: %w", err)
	}
	return p, nil
}

// GetByID retrieves a product by ID, scoped to tenant.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) (*product.Product, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}

// Update persists changes to a product.
func (s *Service) Update(ctx context.Context, p *product.Product) error {
	p.UpdatedAt = time.Now()
	return s.repo.Update(ctx, p)
}

// Delete removes a product by ID, scoped to tenant.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// List returns a paginated list of products for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.Pagination) (kernel.PaginatedResult[product.Product], error) {
	return s.repo.List(ctx, tenantID, pg)
}

// ListByCategory returns products in a specific category.
func (s *Service) ListByCategory(ctx context.Context, tenantID kernel.TenantID, categoryID kernel.CategoryID, pg kernel.Pagination) (kernel.PaginatedResult[product.Product], error) {
	return s.repo.ListByCategory(ctx, tenantID, categoryID, pg)
}

func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40 // version 4
	b[8] = (b[8] & 0x3f) | 0x80 // variant 1
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:16])
}
