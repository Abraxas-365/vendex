package promosrv

import (
	"context"
	"fmt"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/promo"
)

// Service implements promo business logic.
type Service struct {
	repo promo.PromoRepository
}

// New creates a new promo Service.
func New(repo promo.PromoRepository) *Service {
	return &Service{repo: repo}
}

// CreateInput holds the data needed to create a new promo code.
type CreateInput struct {
	TenantID       kernel.TenantID
	Code           string
	Type           promo.PromoType
	Value          int64
	MinOrderAmount *int64
	MaxUses        *int
	StartsAt       *time.Time
	EndsAt         *time.Time

	// Targeting — all optional
	TargetProductIDs  []string
	TargetCategoryIDs []string
	CustomerGroupID   string
	Stackable         bool

	// Buy X Get Y — only used when Type == PromoTypeBuyXGetY
	BuyQuantity  *int
	GetQuantity  *int
	GetProductID string
	GetDiscount  *int64
}

// Create persists a new promo code.
func (s *Service) Create(ctx context.Context, input CreateInput) (*promo.Promo, error) {
	p := &promo.Promo{
		ID:             kernel.PromoID(generateUUID()),
		TenantID:       input.TenantID,
		Code:           input.Code,
		Type:           input.Type,
		Value:          input.Value,
		MinOrderAmount: input.MinOrderAmount,
		MaxUses:        input.MaxUses,
		StartsAt:       input.StartsAt,
		EndsAt:         input.EndsAt,
		Active:         true,
		CreatedAt:      time.Now().UTC(),

		TargetProductIDs:  input.TargetProductIDs,
		TargetCategoryIDs: input.TargetCategoryIDs,
		CustomerGroupID:   input.CustomerGroupID,
		Stackable:         input.Stackable,

		BuyQuantity:  input.BuyQuantity,
		GetQuantity:  input.GetQuantity,
		GetProductID: input.GetProductID,
		GetDiscount:  input.GetDiscount,
	}

	if err := s.repo.Create(ctx, p); err != nil {
		return nil, errx.Wrap(err, "create promo", errx.TypeInternal)
	}
	return p, nil
}

// ValidationResult holds the outcome of a promo validation.
type ValidationResult struct {
	Valid          bool   `json:"valid"`
	DiscountCents  int64  `json:"discount_cents"`
	IsFreeShipping bool   `json:"is_free_shipping"`
	IsBuyXGetY     bool   `json:"is_buy_x_get_y"`
	BuyQuantity    *int   `json:"buy_quantity,omitempty"`
	GetQuantity    *int   `json:"get_quantity,omitempty"`
	GetProductID   string `json:"get_product_id,omitempty"`
	GetDiscount    *int64 `json:"get_discount,omitempty"`
	Reason         string `json:"reason,omitempty"`
}

// Validate checks whether a promo code is applicable to the given order total.
// It does NOT increment the used_count — call Apply separately.
func (s *Service) Validate(ctx context.Context, tenantID kernel.TenantID, code string, orderTotalCents int64) (ValidationResult, error) {
	p, err := s.repo.GetByCode(ctx, tenantID, code)
	if err != nil {
		return ValidationResult{}, err
	}

	now := time.Now().UTC()

	if !p.Active {
		return ValidationResult{Valid: false, Reason: promo.ErrPromoInactive.Error()}, nil
	}
	if !p.IsStarted(now) {
		return ValidationResult{Valid: false, Reason: promo.ErrPromoNotStarted.Error()}, nil
	}
	if p.IsExpired(now) {
		return ValidationResult{Valid: false, Reason: promo.ErrPromoExpired.Error()}, nil
	}
	if p.IsMaxUsesReached() {
		return ValidationResult{Valid: false, Reason: promo.ErrPromoMaxUses.Error()}, nil
	}
	if !p.MeetsMinOrder(orderTotalCents) {
		return ValidationResult{Valid: false, Reason: promo.ErrPromoMinOrder.Error()}, nil
	}

	result := ValidationResult{
		Valid:          true,
		DiscountCents:  p.Discount(orderTotalCents),
		IsFreeShipping: p.Type == promo.PromoTypeFreeShipping,
		IsBuyXGetY:     p.Type == promo.PromoTypeBuyXGetY,
	}

	// Expose buy-X-get-Y metadata so callers can compute item-level discounts.
	if p.Type == promo.PromoTypeBuyXGetY {
		result.BuyQuantity = p.BuyQuantity
		result.GetQuantity = p.GetQuantity
		result.GetProductID = p.GetProductID
		result.GetDiscount = p.GetDiscount
	}

	return result, nil
}

// Apply validates a promo and increments its usage counter atomically.
// Returns the discount amount in cents.
func (s *Service) Apply(ctx context.Context, tenantID kernel.TenantID, code string, orderTotalCents int64) (int64, error) {
	result, err := s.Validate(ctx, tenantID, code, orderTotalCents)
	if err != nil {
		return 0, err
	}
	if !result.Valid {
		return 0, errx.Business(fmt.Sprintf("promo not applicable: %s", result.Reason))
	}

	p, err := s.repo.GetByCode(ctx, tenantID, code)
	if err != nil {
		return 0, err
	}

	if err := s.repo.IncrementUsedCount(ctx, tenantID, p.ID); err != nil {
		return 0, errx.Wrap(err, "increment promo usage", errx.TypeInternal)
	}

	return result.DiscountCents, nil
}

// Deactivate marks a promo as inactive so it can no longer be applied.
func (s *Service) Deactivate(ctx context.Context, tenantID kernel.TenantID, id kernel.PromoID) (*promo.Promo, error) {
	p, err := s.repo.GetByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	p.Active = false
	if err := s.repo.Update(ctx, p); err != nil {
		return nil, errx.Wrap(err, "deactivate promo", errx.TypeInternal)
	}
	return p, nil
}

// List returns all promos for a tenant with pagination.
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[promo.Promo], error) {
	return s.repo.List(ctx, tenantID, p)
}

// GetByID retrieves a promo by ID.
func (s *Service) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PromoID) (*promo.Promo, error) {
	return s.repo.GetByID(ctx, tenantID, id)
}
