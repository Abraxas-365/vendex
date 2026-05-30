package inventory

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines all persistence operations for the inventory domain.
type Repository interface {
	// ─── Warehouse ────────────────────────────────────────────────────────────

	CreateWarehouse(ctx context.Context, w Warehouse) (Warehouse, error)
	GetWarehouse(ctx context.Context, tenantID kernel.TenantID, id kernel.WarehouseID) (Warehouse, error)
	ListWarehouses(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[Warehouse], error)
	UpdateWarehouse(ctx context.Context, w Warehouse) (Warehouse, error)
	DeleteWarehouse(ctx context.Context, tenantID kernel.TenantID, id kernel.WarehouseID) error

	// ─── Stock levels ─────────────────────────────────────────────────────────

	// GetStockLevel fetches the stock level for a specific product/variant/warehouse combo.
	GetStockLevel(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, variantID *kernel.VariantID, warehouseID kernel.WarehouseID) (StockLevel, error)
	// ListStockLevels returns all stock levels for a product across all warehouses.
	ListStockLevels(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) ([]StockLevel, error)
	// UpsertStockLevel creates or updates a stock level record.
	UpsertStockLevel(ctx context.Context, sl StockLevel) (StockLevel, error)
	// GetLowStockItems returns all stock levels whose available quantity is at or below the threshold.
	GetLowStockItems(ctx context.Context, tenantID kernel.TenantID) ([]StockLevel, error)

	// ─── Stock movements ──────────────────────────────────────────────────────

	CreateMovement(ctx context.Context, m StockMovement) (StockMovement, error)
	ListMovements(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, pg kernel.PaginationOptions) (kernel.Paginated[StockMovement], error)
}
