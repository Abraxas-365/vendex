package inventory

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// MovementType categorizes stock movement direction/cause.
type MovementType string

const (
	MovementReceived    MovementType = "received"
	MovementSold        MovementType = "sold"
	MovementReturned    MovementType = "returned"
	MovementAdjusted    MovementType = "adjusted"
	MovementTransferred MovementType = "transferred"
)

// Warehouse represents a physical or logical storage location.
type Warehouse struct {
	ID        kernel.WarehouseID `json:"id" db:"id"`
	TenantID  kernel.TenantID    `json:"tenant_id" db:"tenant_id"`
	Name      string             `json:"name" db:"name"`
	Address   string             `json:"address" db:"address"`
	IsDefault bool               `json:"is_default" db:"is_default"`
	Active    bool               `json:"active" db:"active"`
	CreatedAt time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" db:"updated_at"`
}

// StockLevel tracks how many units of a product/variant are in a warehouse.
type StockLevel struct {
	ID                kernel.StockLevelID `json:"id" db:"id"`
	TenantID          kernel.TenantID     `json:"tenant_id" db:"tenant_id"`
	ProductID         kernel.ProductID    `json:"product_id" db:"product_id"`
	VariantID         *kernel.VariantID   `json:"variant_id,omitempty" db:"variant_id"`
	WarehouseID       kernel.WarehouseID  `json:"warehouse_id" db:"warehouse_id"`
	Quantity          int                 `json:"quantity" db:"quantity"`
	Reserved          int                 `json:"reserved" db:"reserved"`
	LowStockThreshold int                 `json:"low_stock_threshold" db:"low_stock_threshold"`
	CreatedAt         time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at" db:"updated_at"`
}

// Available returns the quantity that can still be sold (quantity minus reserved).
func (s *StockLevel) Available() int {
	avail := s.Quantity - s.Reserved
	if avail < 0 {
		return 0
	}
	return avail
}

// IsLow returns true when available stock is at or below the low stock threshold.
func (s *StockLevel) IsLow() bool {
	return s.Available() <= s.LowStockThreshold
}

// StockMovement records every change to a stock level for auditability.
type StockMovement struct {
	ID          kernel.StockMovementID `json:"id" db:"id"`
	TenantID    kernel.TenantID        `json:"tenant_id" db:"tenant_id"`
	ProductID   kernel.ProductID       `json:"product_id" db:"product_id"`
	VariantID   *kernel.VariantID      `json:"variant_id,omitempty" db:"variant_id"`
	WarehouseID kernel.WarehouseID     `json:"warehouse_id" db:"warehouse_id"`
	Type        MovementType           `json:"type" db:"type"`
	Quantity    int                    `json:"quantity" db:"quantity"`
	Reference   string                 `json:"reference" db:"reference"`
	Note        string                 `json:"note" db:"note"`
	CreatedBy   string                 `json:"created_by" db:"created_by"`
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
}

// ─── Input DTOs ───────────────────────────────────────────────────────────────

// CreateWarehouseInput holds the data needed to create a warehouse.
type CreateWarehouseInput struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	IsDefault bool   `json:"is_default"`
}

// UpdateWarehouseInput holds updatable warehouse fields.
type UpdateWarehouseInput struct {
	Name      string `json:"name"`
	Address   string `json:"address"`
	IsDefault bool   `json:"is_default"`
	Active    bool   `json:"active"`
}

// AdjustStockInput describes a stock adjustment operation.
type AdjustStockInput struct {
	ProductID   kernel.ProductID  `json:"product_id"`
	VariantID   *kernel.VariantID `json:"variant_id,omitempty"`
	WarehouseID kernel.WarehouseID `json:"warehouse_id"`
	Quantity    int               `json:"quantity"`
	Type        MovementType      `json:"type"`
	Reference   string            `json:"reference"`
	Note        string            `json:"note"`
	CreatedBy   string            `json:"created_by"`
}
