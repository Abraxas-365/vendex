package bundlesrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/hada-commerce/internal/bundle"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/eventbus"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Service handles bundle business logic.
type Service struct {
	repo bundle.Repository
	bus  eventbus.Bus
}

// New creates a new bundle service.
func New(repo bundle.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// ─── Bundle CRUD ──────────────────────────────────────────────────────────────

// Create creates a new bundle for the given tenant.
func (s *Service) Create(ctx context.Context, tenantID kernel.TenantID, in bundle.CreateBundleInput) (*bundle.Bundle, error) {
	if in.Name == "" {
		return nil, errx.New("bundle name is required", errx.TypeValidation)
	}

	// Validate discount.
	if err := validateDiscount(in.DiscountType, in.DiscountValue); err != nil {
		return nil, err
	}

	// Auto-generate slug if not provided.
	slug := in.Slug
	if slug == "" {
		slug = bundle.GenerateSlug(in.Name)
	}

	// Ensure slug is unique within tenant.
	existing, err := s.repo.GetBySlug(ctx, tenantID, slug)
	if err != nil && !errx.Is(err, bundle.ErrNotFound) {
		return nil, errx.Wrap(err, "checking slug uniqueness", errx.TypeInternal)
	}
	if existing != nil {
		return nil, bundle.ErrSlugTaken
	}

	now := time.Now()
	b := &bundle.Bundle{
		ID:            kernel.BundleID(uuid.NewString()),
		TenantID:      tenantID,
		Name:          in.Name,
		Slug:          slug,
		Description:   in.Description,
		DiscountType:  in.DiscountType,
		DiscountValue: in.DiscountValue,
		Active:        in.Active,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	if err := s.repo.Create(ctx, b); err != nil {
		return nil, errx.Wrap(err, "creating bundle", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.BundleCreated, tenantID, eventbus.BundlePayload{
		BundleID:      string(b.ID),
		Name:          b.Name,
		Slug:          b.Slug,
		DiscountType:  string(b.DiscountType),
		DiscountValue: b.DiscountValue,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return b, nil
}

// GetByID retrieves a bundle by ID and eagerly loads its items.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BundleID) (*bundle.Bundle, error) {
	b, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.ListItems(ctx, tenantID, id)
	if err != nil {
		return nil, errx.Wrap(err, "loading bundle items", errx.TypeInternal)
	}
	b.Items = items

	return b, nil
}

// GetBySlug retrieves a bundle by slug and eagerly loads its items.
func (s *Service) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*bundle.Bundle, error) {
	b, err := s.repo.GetBySlug(ctx, tenantID, slug)
	if err != nil {
		return nil, err
	}

	items, err := s.repo.ListItems(ctx, tenantID, b.ID)
	if err != nil {
		return nil, errx.Wrap(err, "loading bundle items", errx.TypeInternal)
	}
	b.Items = items

	return b, nil
}

// List returns a paginated list of bundles for a tenant.
// When activeOnly is true, only active bundles are returned.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, activeOnly bool, pg kernel.PaginationOptions) (kernel.Paginated[bundle.Bundle], error) {
	return s.repo.List(ctx, tenantID, activeOnly, pg)
}

// Update applies changes from UpdateBundleInput to an existing bundle.
func (s *Service) Update(ctx context.Context, tenantID kernel.TenantID, id kernel.BundleID, in bundle.UpdateBundleInput) (*bundle.Bundle, error) {
	b, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		b.Name = *in.Name
	}
	if in.Description != nil {
		b.Description = *in.Description
	}
	if in.DiscountType != nil {
		b.DiscountType = *in.DiscountType
	}
	if in.DiscountValue != nil {
		b.DiscountValue = *in.DiscountValue
	}
	if in.Active != nil {
		b.Active = *in.Active
	}

	if err := validateDiscount(b.DiscountType, b.DiscountValue); err != nil {
		return nil, err
	}

	b.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, b); err != nil {
		return nil, errx.Wrap(err, "updating bundle", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.BundleUpdated, tenantID, eventbus.BundlePayload{
		BundleID:      string(b.ID),
		Name:          b.Name,
		Slug:          b.Slug,
		DiscountType:  string(b.DiscountType),
		DiscountValue: b.DiscountValue,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return b, nil
}

// Delete removes a bundle (and its items via CASCADE) by ID.
func (s *Service) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.BundleID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// ─── Bundle Items ──────────────────────────────────────────────────────────────

// AddItem adds a product (with optional variant) to a bundle.
func (s *Service) AddItem(ctx context.Context, tenantID kernel.TenantID, bundleID kernel.BundleID, in bundle.AddBundleItemInput) (*bundle.BundleItem, error) {
	// Verify the bundle exists.
	if _, err := s.repo.GetByID(ctx, tenantID, bundleID); err != nil {
		return nil, err
	}

	if in.Quantity < 1 {
		in.Quantity = 1
	}

	var variantID *kernel.VariantID
	if in.VariantID != nil && *in.VariantID != "" {
		vid := kernel.VariantID(*in.VariantID)
		variantID = &vid
	}

	item := &bundle.BundleItem{
		ID:        kernel.BundleItemID(uuid.NewString()),
		TenantID:  tenantID,
		BundleID:  bundleID,
		ProductID: kernel.ProductID(in.ProductID),
		VariantID: variantID,
		Quantity:  in.Quantity,
		CreatedAt: time.Now(),
	}

	if err := s.repo.AddItem(ctx, item); err != nil {
		return nil, errx.Wrap(err, "adding bundle item", errx.TypeInternal)
	}

	return item, nil
}

// RemoveItem removes an item from a bundle.
func (s *Service) RemoveItem(ctx context.Context, tenantID kernel.TenantID, bundleID kernel.BundleID, itemID kernel.BundleItemID) error {
	// Verify the bundle exists and belongs to tenant.
	if _, err := s.repo.GetByID(ctx, tenantID, bundleID); err != nil {
		return err
	}

	// Verify the item exists and belongs to the tenant.
	if _, err := s.repo.GetItemByID(ctx, tenantID, itemID); err != nil {
		return err
	}

	return s.repo.RemoveItem(ctx, tenantID, itemID)
}

// ─── Price Calculation ────────────────────────────────────────────────────────

// CalculatePrice computes the bundle pricing given a map of product prices keyed
// by ProductID. The caller is responsible for fetching those prices from the
// product domain.
func (s *Service) CalculatePrice(
	ctx context.Context,
	tenantID kernel.TenantID,
	id kernel.BundleID,
	productPrices map[kernel.ProductID]kernel.Money,
) (bundle.BundlePriceResult, error) {
	b, err := s.GetByID(ctx, tenantID, id)
	if err != nil {
		return bundle.BundlePriceResult{}, err
	}

	if len(b.Items) == 0 {
		return bundle.BundlePriceResult{}, bundle.ErrNoItems
	}

	// Default currency to USD; override with first item's currency if available.
	currency := "USD"
	var baseAmount int64
	for _, item := range b.Items {
		price, ok := productPrices[item.ProductID]
		if !ok {
			continue
		}
		if currency == "USD" && price.Currency != "" {
			currency = price.Currency
		}
		baseAmount += price.Amount * int64(item.Quantity)
	}

	baseTotal := kernel.NewMoney(baseAmount, currency)

	var discountAmount int64
	switch b.DiscountType {
	case bundle.DiscountPercentage:
		discountAmount = baseAmount * int64(b.DiscountValue) / 100
	case bundle.DiscountFixed:
		discountAmount = int64(b.DiscountValue)
		if discountAmount > baseAmount {
			discountAmount = baseAmount
		}
	}

	finalAmount := baseAmount - discountAmount
	if finalAmount < 0 {
		finalAmount = 0
	}

	return bundle.BundlePriceResult{
		BundleID:       b.ID,
		BaseTotal:      baseTotal,
		DiscountAmount: kernel.NewMoney(discountAmount, currency),
		FinalTotal:     kernel.NewMoney(finalAmount, currency),
		DiscountType:   b.DiscountType,
		DiscountValue:  b.DiscountValue,
	}, nil
}

// ─── Helpers ──────────────────────────────────────────────────────────────────

func validateDiscount(dt bundle.DiscountType, value int) error {
	switch dt {
	case bundle.DiscountPercentage:
		if value < 0 || value > 100 {
			return bundle.ErrInvalidDiscount
		}
	case bundle.DiscountFixed:
		if value < 0 {
			return bundle.ErrInvalidDiscount
		}
	case "":
		// default to percentage
	default:
		return errx.New("discount_type must be 'percentage' or 'fixed'", errx.TypeValidation)
	}
	return nil
}
