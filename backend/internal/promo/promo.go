package promo

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PromoType determines how a discount is applied.
type PromoType string

const (
	PromoTypePercentage   PromoType = "percentage"    // Value is a percentage (0–100)
	PromoTypeFixedAmount  PromoType = "fixed_amount"  // Value is cents subtracted from total
	PromoTypeFreeShipping PromoType = "free_shipping" // Free shipping; Value is ignored
)

// Promo represents a promotional discount code.
// Business rules:
//   - A promo is valid only when Active=true, within the StartsAt–EndsAt window,
//     UsedCount < MaxUses (when MaxUses is set), and the order meets MinOrderAmount.
//   - Applying a promo increments UsedCount atomically in the repository.
//   - Deactivated promos cannot be reactivated via normal service methods.
type Promo struct {
	ID             kernel.PromoID  `json:"id" db:"id"`
	TenantID       kernel.TenantID `json:"tenant_id" db:"tenant_id"`
	Code           string          `json:"code" db:"code"`
	Type           PromoType       `json:"type" db:"type"`
	// Value meaning depends on Type:
	//   percentage  → integer 0–100
	//   fixed_amount → cents
	//   free_shipping → not used
	Value          int64           `json:"value" db:"value"`
	MinOrderAmount *int64          `json:"min_order_amount,omitempty" db:"min_order_amount"`
	MaxUses        *int            `json:"max_uses,omitempty" db:"max_uses"`
	UsedCount      int             `json:"used_count" db:"used_count"`
	StartsAt       *time.Time      `json:"starts_at,omitempty" db:"starts_at"`
	EndsAt         *time.Time      `json:"ends_at,omitempty" db:"ends_at"`
	Active         bool            `json:"active" db:"active"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
}

// IsExpired returns true when the promo is past its EndsAt date.
func (p *Promo) IsExpired(now time.Time) bool {
	return p.EndsAt != nil && now.After(*p.EndsAt)
}

// IsStarted returns true when the promo's start window has begun.
func (p *Promo) IsStarted(now time.Time) bool {
	return p.StartsAt == nil || !now.Before(*p.StartsAt)
}

// IsMaxUsesReached returns true when the promo has no remaining uses.
func (p *Promo) IsMaxUsesReached() bool {
	return p.MaxUses != nil && p.UsedCount >= *p.MaxUses
}

// MeetsMinOrder returns true when the order total satisfies the minimum order amount.
func (p *Promo) MeetsMinOrder(orderTotalCents int64) bool {
	return p.MinOrderAmount == nil || orderTotalCents >= *p.MinOrderAmount
}

// Discount computes the discount amount (in cents) for the given order total.
func (p *Promo) Discount(orderTotalCents int64) int64 {
	switch p.Type {
	case PromoTypePercentage:
		return orderTotalCents * p.Value / 100
	case PromoTypeFixedAmount:
		if p.Value > orderTotalCents {
			return orderTotalCents
		}
		return p.Value
	case PromoTypeFreeShipping:
		return 0 // shipping discount handled separately
	default:
		return 0
	}
}
