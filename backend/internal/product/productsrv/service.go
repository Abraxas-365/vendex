package productsrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/product"
)

// Service handles product business logic.
type Service struct {
	repo        product.Repository
	variantRepo product.VariantRepository
	bus         eventbus.Bus
}

// New creates a new product service.
func New(repo product.Repository, variantRepo product.VariantRepository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, variantRepo: variantRepo, bus: bus}
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
			return nil, errx.Wrap(err, "checking SKU uniqueness", errx.TypeInternal)
		}
		if existing != nil {
			return nil, product.ErrDuplicateSKU
		}
	}

	now := time.Now()
	p := &product.Product{
		ID:          kernel.ProductID(uuid.NewString()),
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
		return nil, errx.Wrap(err, "creating product", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.ProductCreated, tenantID, eventbus.ProductPayload{
		ProductID: string(p.ID),
		Name:      p.Name,
		SKU:       p.SKU,
		Price:     int(p.Price.Amount),
		Currency:  p.Price.Currency,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return p, nil
}

// GetByID retrieves a product by ID, scoped to tenant.
// It also loads the product's options and variants.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) (*product.Product, error) {
	p, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	// Load options and variants eagerly.
	opts, err := s.variantRepo.ListOptions(ctx, tenantID, id)
	if err != nil {
		return nil, errx.Wrap(err, "loading product options", errx.TypeInternal)
	}
	p.Options = opts

	variants, err := s.variantRepo.ListVariants(ctx, tenantID, id)
	if err != nil {
		return nil, errx.Wrap(err, "loading product variants", errx.TypeInternal)
	}
	p.Variants = variants

	return p, nil
}

// Update persists changes to a product.
func (s *Service) Update(ctx context.Context, p *product.Product) error {
	p.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, p); err != nil {
		return err
	}

	if evt, err := eventbus.NewEvent(eventbus.ProductUpdated, p.TenantID, eventbus.ProductPayload{
		ProductID: string(p.ID),
		Name:      p.Name,
		SKU:       p.SKU,
		Price:     int(p.Price.Amount),
		Currency:  p.Price.Currency,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return nil
}

// Delete removes a product by ID, scoped to tenant.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) error {
	if err := s.repo.Delete(ctx, tenantID, id); err != nil {
		return err
	}

	if evt, err := eventbus.NewEvent(eventbus.ProductDeleted, tenantID, eventbus.ProductPayload{
		ProductID: string(id),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return nil
}

// List returns a paginated list of products for a tenant.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[product.Product], error) {
	return s.repo.List(ctx, tenantID, pg)
}

// ListByCategory returns products in a specific category.
func (s *Service) ListByCategory(ctx context.Context, tenantID kernel.TenantID, categoryID kernel.CategoryID, pg kernel.PaginationOptions) (kernel.Paginated[product.Product], error) {
	return s.repo.ListByCategory(ctx, tenantID, categoryID, pg)
}

// ─── Option methods ───────────────────────────────────────────────────────────

// CreateOption creates a new product option.
func (s *Service) CreateOption(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, name string, position int, values []string) (*product.ProductOption, error) {
	now := time.Now()
	opt := &product.ProductOption{
		ID:        kernel.OptionID(uuid.NewString()),
		ProductID: productID,
		TenantID:  tenantID,
		Name:      name,
		Position:  position,
		Values:    values,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.variantRepo.CreateOption(ctx, opt); err != nil {
		return nil, err
	}
	return opt, nil
}

// ListOptions returns all options for a product.
func (s *Service) ListOptions(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) ([]product.ProductOption, error) {
	return s.variantRepo.ListOptions(ctx, tenantID, productID)
}

// UpdateOption updates a product option's name, position, and values.
func (s *Service) UpdateOption(ctx context.Context, tenantID kernel.TenantID, optionID kernel.OptionID, name string, position int, values []string) (*product.ProductOption, error) {
	now := time.Now()
	opt := &product.ProductOption{
		ID:        optionID,
		TenantID:  tenantID,
		Name:      name,
		Position:  position,
		Values:    values,
		UpdatedAt: now,
	}
	if err := s.variantRepo.UpdateOption(ctx, opt); err != nil {
		return nil, err
	}
	return opt, nil
}

// DeleteOption removes a product option.
func (s *Service) DeleteOption(ctx context.Context, tenantID kernel.TenantID, optionID kernel.OptionID) error {
	return s.variantRepo.DeleteOption(ctx, tenantID, optionID)
}

// ─── Variant methods ──────────────────────────────────────────────────────────

// CreateVariant creates a new product variant.
func (s *Service) CreateVariant(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, sku string, price kernel.Money, stock int, options map[string]string) (*product.ProductVariant, error) {
	if price.Amount <= 0 {
		return nil, product.ErrInvalidPrice
	}

	// Check for duplicate SKU within tenant.
	if sku != "" {
		existing, err := s.variantRepo.GetVariantBySKU(ctx, tenantID, sku)
		if err != nil && !errx.Is(err, product.ErrVariantNotFound) {
			return nil, errx.Wrap(err, "checking variant SKU uniqueness", errx.TypeInternal)
		}
		if existing != nil {
			return nil, product.ErrDuplicateVariantSKU
		}
	}

	now := time.Now()
	v := &product.ProductVariant{
		ID:        kernel.VariantID(uuid.NewString()),
		ProductID: productID,
		TenantID:  tenantID,
		SKU:       sku,
		Price:     price,
		Stock:     stock,
		Options:   options,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if v.Options == nil {
		v.Options = map[string]string{}
	}

	if err := s.variantRepo.CreateVariant(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

// GetVariant retrieves a specific variant by ID.
func (s *Service) GetVariant(ctx context.Context, tenantID kernel.TenantID, variantID kernel.VariantID) (*product.ProductVariant, error) {
	return s.variantRepo.GetVariantByID(ctx, tenantID, variantID)
}

// ListVariants returns all variants for a product.
func (s *Service) ListVariants(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) ([]product.ProductVariant, error) {
	return s.variantRepo.ListVariants(ctx, tenantID, productID)
}

// UpdateVariant updates a product variant's fields.
func (s *Service) UpdateVariant(ctx context.Context, tenantID kernel.TenantID, variantID kernel.VariantID, sku string, price kernel.Money, stock int, options map[string]string, active bool) (*product.ProductVariant, error) {
	now := time.Now()
	v := &product.ProductVariant{
		ID:        variantID,
		TenantID:  tenantID,
		SKU:       sku,
		Price:     price,
		Stock:     stock,
		Options:   options,
		Active:    active,
		UpdatedAt: now,
	}
	if v.Options == nil {
		v.Options = map[string]string{}
	}

	if err := s.variantRepo.UpdateVariant(ctx, v); err != nil {
		return nil, err
	}
	return v, nil
}

// DeleteVariant removes a product variant.
func (s *Service) DeleteVariant(ctx context.Context, tenantID kernel.TenantID, variantID kernel.VariantID) error {
	return s.variantRepo.DeleteVariant(ctx, tenantID, variantID)
}
