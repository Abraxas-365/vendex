package inventorysrv

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/inventory"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Service implements inventory business logic.
type Service struct {
	repo inventory.Repository
	bus  eventbus.Bus
}

// NewService creates a new inventory service.
func NewService(repo inventory.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// ─── Warehouse ────────────────────────────────────────────────────────────────

// CreateWarehouse creates a new warehouse for the given tenant.
func (s *Service) CreateWarehouse(ctx context.Context, tenantID kernel.TenantID, in inventory.CreateWarehouseInput) (inventory.Warehouse, error) {
	if in.Name == "" {
		return inventory.Warehouse{}, errx.New("warehouse name is required", errx.TypeValidation)
	}

	now := time.Now().UTC()
	w := inventory.Warehouse{
		ID:        kernel.WarehouseID(uuid.NewString()),
		TenantID:  tenantID,
		Name:      in.Name,
		Address:   in.Address,
		IsDefault: in.IsDefault,
		Active:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	return s.repo.CreateWarehouse(ctx, w)
}

// GetWarehouse retrieves a warehouse by ID, scoped to tenant.
func (s *Service) GetWarehouse(ctx context.Context, tenantID kernel.TenantID, id kernel.WarehouseID) (inventory.Warehouse, error) {
	return s.repo.GetWarehouse(ctx, tenantID, id)
}

// ListWarehouses returns a paginated list of warehouses for the tenant.
func (s *Service) ListWarehouses(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[inventory.Warehouse], error) {
	return s.repo.ListWarehouses(ctx, tenantID, pg)
}

// UpdateWarehouse applies the given update to an existing warehouse.
func (s *Service) UpdateWarehouse(ctx context.Context, tenantID kernel.TenantID, id kernel.WarehouseID, in inventory.UpdateWarehouseInput) (inventory.Warehouse, error) {
	w, err := s.repo.GetWarehouse(ctx, tenantID, id)
	if err != nil {
		return inventory.Warehouse{}, err
	}

	w.Name = in.Name
	w.Address = in.Address
	w.IsDefault = in.IsDefault
	w.Active = in.Active
	w.UpdatedAt = time.Now().UTC()

	return s.repo.UpdateWarehouse(ctx, w)
}

// DeleteWarehouse permanently removes a warehouse.
func (s *Service) DeleteWarehouse(ctx context.Context, tenantID kernel.TenantID, id kernel.WarehouseID) error {
	return s.repo.DeleteWarehouse(ctx, tenantID, id)
}

// ─── Stock ────────────────────────────────────────────────────────────────────

// GetStock returns the stock level for a specific product/variant in a warehouse.
// If warehouseID is empty, the first (default) warehouse's stock is returned.
func (s *Service) GetStock(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, variantID *kernel.VariantID, warehouseID kernel.WarehouseID) (inventory.StockLevel, error) {
	return s.repo.GetStockLevel(ctx, tenantID, productID, variantID, warehouseID)
}

// ListStockLevels returns all stock level entries for a product across warehouses.
func (s *Service) ListStockLevels(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) ([]inventory.StockLevel, error) {
	return s.repo.ListStockLevels(ctx, tenantID, productID)
}

// AdjustStock applies a stock adjustment, records the movement, updates the stock level,
// and fires domain events.
func (s *Service) AdjustStock(ctx context.Context, tenantID kernel.TenantID, in inventory.AdjustStockInput) (inventory.StockLevel, error) {
	// Validate movement type.
	switch in.Type {
	case inventory.MovementReceived, inventory.MovementSold, inventory.MovementReturned,
		inventory.MovementAdjusted, inventory.MovementTransferred:
	default:
		return inventory.StockLevel{}, inventory.ErrInvalidMovementType
	}

	// Fetch or build current stock level.
	existing, err := s.repo.GetStockLevel(ctx, tenantID, in.ProductID, in.VariantID, in.WarehouseID)
	if err != nil && !errx.IsNotFound(err) {
		return inventory.StockLevel{}, errx.Wrap(err, "fetching stock level", errx.TypeInternal)
	}

	now := time.Now().UTC()

	// Bootstrap a new stock level if none exists yet.
	if errx.IsNotFound(err) {
		existing = inventory.StockLevel{
			ID:                kernel.StockLevelID(uuid.NewString()),
			TenantID:          tenantID,
			ProductID:         in.ProductID,
			VariantID:         in.VariantID,
			WarehouseID:       in.WarehouseID,
			Quantity:          0,
			Reserved:          0,
			LowStockThreshold: 5,
			CreatedAt:         now,
			UpdatedAt:         now,
		}
	}

	// Apply delta.
	newQty := existing.Quantity + in.Quantity
	if newQty < 0 {
		return inventory.StockLevel{}, inventory.ErrInvalidQuantity
	}
	existing.Quantity = newQty
	existing.UpdatedAt = now

	// Persist stock level.
	updated, err := s.repo.UpsertStockLevel(ctx, existing)
	if err != nil {
		return inventory.StockLevel{}, errx.Wrap(err, "upserting stock level", errx.TypeInternal)
	}

	// Record movement.
	movement := inventory.StockMovement{
		ID:          kernel.StockMovementID(uuid.NewString()),
		TenantID:    tenantID,
		ProductID:   in.ProductID,
		VariantID:   in.VariantID,
		WarehouseID: in.WarehouseID,
		Type:        in.Type,
		Quantity:    in.Quantity,
		Reference:   in.Reference,
		Note:        in.Note,
		CreatedBy:   in.CreatedBy,
		CreatedAt:   now,
	}
	if _, err := s.repo.CreateMovement(ctx, movement); err != nil {
		// Non-fatal — stock is already updated. Log via event bus if desired.
		_ = err
	}

	// Fire stock.updated event.
	if evt, err := eventbus.NewEvent(eventbus.StockUpdated, tenantID, map[string]any{
		"product_id":   string(in.ProductID),
		"warehouse_id": string(in.WarehouseID),
		"quantity":     updated.Quantity,
		"available":    updated.Available(),
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	// Fire low-stock alert if threshold crossed.
	if updated.IsLow() {
		if evt, err := eventbus.NewEvent(eventbus.StockLowAlert, tenantID, map[string]any{
			"product_id":   string(in.ProductID),
			"warehouse_id": string(in.WarehouseID),
			"available":    updated.Available(),
			"threshold":    updated.LowStockThreshold,
		}); err == nil {
			_ = s.bus.Publish(ctx, evt)
		}
	}

	return updated, nil
}

// GetLowStockItems returns all stock levels that are at or below their low-stock threshold.
func (s *Service) GetLowStockItems(ctx context.Context, tenantID kernel.TenantID) ([]inventory.StockLevel, error) {
	return s.repo.GetLowStockItems(ctx, tenantID)
}

// ListMovements returns a paginated history of stock movements for a product.
func (s *Service) ListMovements(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, pg kernel.PaginationOptions) (kernel.Paginated[inventory.StockMovement], error) {
	return s.repo.ListMovements(ctx, tenantID, productID, pg)
}
