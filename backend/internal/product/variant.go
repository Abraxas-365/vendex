package product

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// ProductOption represents a configurable option for a product, e.g. "Size" or "Color".
// Each option has a set of predefined values (e.g. ["S","M","L","XL"]).
type ProductOption struct {
	ID        kernel.OptionID  `json:"id" db:"id"`
	ProductID kernel.ProductID `json:"product_id" db:"product_id"`
	TenantID  kernel.TenantID  `json:"tenant_id" db:"tenant_id"`
	Name      string           `json:"name" db:"name"`     // e.g. "Size", "Color"
	Position  int              `json:"position" db:"position"` // display order
	Values    []string         `json:"values"`              // e.g. ["S","M","L","XL"] stored as JSONB
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt time.Time        `json:"updated_at" db:"updated_at"`
}

// ProductVariant represents a specific combination of option values with its own price/SKU/stock.
type ProductVariant struct {
	ID        kernel.VariantID  `json:"id" db:"id"`
	ProductID kernel.ProductID  `json:"product_id" db:"product_id"`
	TenantID  kernel.TenantID   `json:"tenant_id" db:"tenant_id"`
	SKU       string            `json:"sku" db:"sku"`
	Price     kernel.Money      `json:"price"`
	Stock     int               `json:"stock" db:"stock"`
	Options   map[string]string `json:"options"`  // e.g. {"Size":"M","Color":"Red"} stored as JSONB
	Active    bool              `json:"active" db:"active"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
}
