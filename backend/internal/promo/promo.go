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
	PromoTypeBuyXGetY     PromoType = "buy_x_get_y"  // Buy X items, get Y free/discounted
)

// Promo represents a promotional discount code.
// Business rules:
//   - A promo is valid only when Active=true, within the StartsAt–EndsAt window,
//     UsedCount < MaxUses (when MaxUses is set), and the order meets MinOrderAmount.
//   - Applying a promo increments UsedCount atomically in the repository.
//   - Deactivated promos cannot be reactivated via normal service methods.
type Promo struct {
	ID       kernel.PromoID  `json:"id" db:"id"`
	TenantID kernel.TenantID `json:"tenant_id" db:"tenant_id"`
	Code     string          `json:"code" db:"code"`
	Type     PromoType       `json:"type" db:"type"`
	// Value meaning depends on Type:
	//   percentage  → integer 0–100
	//   fixed_amount → cents
	//   free_shipping → not used
	//   buy_x_get_y → not used (use BuyQuantity/GetQuantity/GetDiscount)
	Value          int64      `json:"value" db:"value"`
	MinOrderAmount *int64     `json:"min_order_amount,omitempty" db:"min_order_amount"`
	MaxUses        *int       `json:"max_uses,omitempty" db:"max_uses"`
	UsedCount      int        `json:"used_count" db:"used_count"`
	StartsAt       *time.Time `json:"starts_at,omitempty" db:"starts_at"`
	EndsAt         *time.Time `json:"ends_at,omitempty" db:"ends_at"`
	Active         bool       `json:"active" db:"active"`
	CreatedAt      time.Time  `json:"created_at" db:"created_at"`

	// Targeting — all optional; nil/empty means "applies to everything"
	TargetProductIDs  []string `json:"target_product_ids,omitempty" db:"target_product_ids"`
	TargetCategoryIDs []string `json:"target_category_ids,omitempty" db:"target_category_ids"`
	CustomerGroupID   string   `json:"customer_group_id,omitempty" db:"customer_group_id"`
	Stackable         bool     `json:"stackable" db:"stackable"`

	// Buy X Get Y fields (only used when Type == PromoTypeBuyXGetY)
	BuyQuantity *int    `json:"buy_quantity,omitempty" db:"buy_quantity"`
	GetQuantity *int    `json:"get_quantity,omitempty" db:"get_quantity"`
	GetProductID string `json:"get_product_id,omitempty" db:"get_product_id"`
	GetDiscount  *int64 `json:"get_discount,omitempty" db:"get_discount"`
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

// MatchesProduct returns true if this promo has no product targeting, or if
// the given productID is in the target list.
func (p *Promo) MatchesProduct(productID string) bool {
	if len(p.TargetProductIDs) == 0 {
		return true
	}
	for _, id := range p.TargetProductIDs {
		if id == productID {
			return true
		}
	}
	return false
}

// MatchesCategory returns true if this promo has no category targeting, or if
// the given categoryID is in the target list.
func (p *Promo) MatchesCategory(categoryID string) bool {
	if len(p.TargetCategoryIDs) == 0 {
		return true
	}
	for _, id := range p.TargetCategoryIDs {
		if id == categoryID {
			return true
		}
	}
	return false
}

// Discount computes the discount amount (in cents) for the given order total.
// For buy_x_get_y promos, use CalculateBuyXGetYDiscount instead.
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
	case PromoTypeBuyXGetY:
		return 0 // use CalculateBuyXGetYDiscount with item-level context
	default:
		return 0
	}
}

// CalculateBuyXGetYDiscount computes the discount for a buy-X-get-Y promo.
// qualifyingQty is the number of qualifying items in the cart,
// getProductPrice is the unit price (cents) of the "get" product.
// Returns 0 if any required field is missing.
func (p *Promo) CalculateBuyXGetYDiscount(qualifyingQty int, getProductPrice int64) int64 {
	if p.BuyQuantity == nil || p.GetQuantity == nil || p.GetDiscount == nil {
		return 0
	}
	sets := qualifyingQty / *p.BuyQuantity
	freeItems := sets * *p.GetQuantity
	return int64(freeItems) * getProductPrice * *p.GetDiscount / 100
}
