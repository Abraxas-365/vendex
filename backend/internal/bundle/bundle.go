package bundle

import (
	"strings"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// DiscountType defines how the bundle discount is applied.
type DiscountType string

const (
	// DiscountPercentage applies a percentage discount (0–100) to the bundle total.
	DiscountPercentage DiscountType = "percentage"
	// DiscountFixed deducts a fixed amount (in cents) from the bundle total.
	DiscountFixed DiscountType = "fixed"
)

// Bundle is a grouped set of products sold together at a discounted price.
type Bundle struct {
	ID            kernel.BundleID `json:"id"             db:"id"`
	TenantID      kernel.TenantID `json:"tenant_id"      db:"tenant_id"`
	Name          string          `json:"name"           db:"name"`
	Slug          string          `json:"slug"           db:"slug"`
	Description   string          `json:"description"    db:"description"`
	DiscountType  DiscountType    `json:"discount_type"  db:"discount_type"`
	DiscountValue int             `json:"discount_value" db:"discount_value"`
	Active        bool            `json:"active"         db:"active"`
	Items         []BundleItem    `json:"items,omitempty"`
	CreatedAt     time.Time       `json:"created_at"     db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"     db:"updated_at"`
}

// BundleItem is a single product (with optional variant) inside a bundle.
type BundleItem struct {
	ID        kernel.BundleItemID `json:"id"         db:"id"`
	TenantID  kernel.TenantID     `json:"tenant_id"  db:"tenant_id"`
	BundleID  kernel.BundleID     `json:"bundle_id"  db:"bundle_id"`
	ProductID kernel.ProductID    `json:"product_id" db:"product_id"`
	VariantID *kernel.VariantID   `json:"variant_id" db:"variant_id"`
	Quantity  int                 `json:"quantity"   db:"quantity"`
	CreatedAt time.Time           `json:"created_at" db:"created_at"`
}

// BundlePriceResult holds the pricing breakdown for a bundle.
type BundlePriceResult struct {
	BundleID      kernel.BundleID `json:"bundle_id"`
	BaseTotal     kernel.Money    `json:"base_total"`      // sum of item prices × quantities
	DiscountAmount kernel.Money   `json:"discount_amount"` // amount saved
	FinalTotal    kernel.Money    `json:"final_total"`     // price to pay
	DiscountType  DiscountType    `json:"discount_type"`
	DiscountValue int             `json:"discount_value"`
}

// CreateBundleInput contains the data required to create a bundle.
type CreateBundleInput struct {
	Name          string
	Slug          string
	Description   string
	DiscountType  DiscountType
	DiscountValue int
	Active        bool
}

// UpdateBundleInput contains the fields that may be updated on a bundle.
type UpdateBundleInput struct {
	Name          *string
	Description   *string
	DiscountType  *DiscountType
	DiscountValue *int
	Active        *bool
}

// AddBundleItemInput contains the data needed to add an item to a bundle.
type AddBundleItemInput struct {
	ProductID string
	VariantID *string
	Quantity  int
}

// GenerateSlug creates a URL-safe slug from a bundle name.
func GenerateSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteByte('-')
		}
	}
	result := b.String()
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}
