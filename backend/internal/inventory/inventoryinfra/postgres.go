package inventoryinfra

import (
	"context"
	"database/sql"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/inventory"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements inventory.Repository using sqlx / PostgreSQL.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed inventory repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// ─── Warehouse ────────────────────────────────────────────────────────────────

func (r *PostgresRepo) CreateWarehouse(ctx context.Context, w inventory.Warehouse) (inventory.Warehouse, error) {
	const q = `
		INSERT INTO warehouses (id, tenant_id, name, address, is_default, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := r.db.ExecContext(ctx, q,
		string(w.ID), string(w.TenantID), w.Name, w.Address,
		w.IsDefault, w.Active, w.CreatedAt, w.UpdatedAt,
	)
	if err != nil {
		return inventory.Warehouse{}, errx.Wrap(err, "inserting warehouse", errx.TypeInternal)
	}
	return w, nil
}

func (r *PostgresRepo) GetWarehouse(ctx context.Context, tenantID kernel.TenantID, id kernel.WarehouseID) (inventory.Warehouse, error) {
	const q = `
		SELECT id, tenant_id, name, address, is_default, active, created_at, updated_at
		FROM warehouses WHERE id = $1 AND tenant_id = $2`

	var w inventory.Warehouse
	var wID, wTenantID string

	err := r.db.QueryRowContext(ctx, q, string(id), string(tenantID)).Scan(
		&wID, &wTenantID, &w.Name, &w.Address,
		&w.IsDefault, &w.Active, &w.CreatedAt, &w.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return inventory.Warehouse{}, inventory.ErrWarehouseNotFound
	}
	if err != nil {
		return inventory.Warehouse{}, errx.Wrap(err, "scanning warehouse", errx.TypeInternal)
	}

	w.ID = kernel.WarehouseID(wID)
	w.TenantID = kernel.TenantID(wTenantID)
	return w, nil
}

func (r *PostgresRepo) ListWarehouses(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[inventory.Warehouse], error) {
	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM warehouses WHERE tenant_id = $1`,
		string(tenantID),
	).Scan(&total); err != nil {
		return kernel.Paginated[inventory.Warehouse]{}, errx.Wrap(err, "counting warehouses", errx.TypeInternal)
	}

	const q = `
		SELECT id, tenant_id, name, address, is_default, active, created_at, updated_at
		FROM warehouses WHERE tenant_id = $1
		ORDER BY is_default DESC, name ASC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.QueryContext(ctx, q, string(tenantID), pg.Limit(), pg.Offset())
	if err != nil {
		return kernel.Paginated[inventory.Warehouse]{}, errx.Wrap(err, "querying warehouses", errx.TypeInternal)
	}
	defer rows.Close()

	var items []inventory.Warehouse
	for rows.Next() {
		var w inventory.Warehouse
		var wID, wTenantID string
		if err := rows.Scan(&wID, &wTenantID, &w.Name, &w.Address, &w.IsDefault, &w.Active, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return kernel.Paginated[inventory.Warehouse]{}, errx.Wrap(err, "scanning warehouse row", errx.TypeInternal)
		}
		w.ID = kernel.WarehouseID(wID)
		w.TenantID = kernel.TenantID(wTenantID)
		items = append(items, w)
	}
	if err := rows.Err(); err != nil {
		return kernel.Paginated[inventory.Warehouse]{}, errx.Wrap(err, "iterating warehouses", errx.TypeInternal)
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

func (r *PostgresRepo) UpdateWarehouse(ctx context.Context, w inventory.Warehouse) (inventory.Warehouse, error) {
	const q = `
		UPDATE warehouses
		SET name=$1, address=$2, is_default=$3, active=$4, updated_at=$5
		WHERE id=$6 AND tenant_id=$7`

	_, err := r.db.ExecContext(ctx, q,
		w.Name, w.Address, w.IsDefault, w.Active, w.UpdatedAt,
		string(w.ID), string(w.TenantID),
	)
	if err != nil {
		return inventory.Warehouse{}, errx.Wrap(err, "updating warehouse", errx.TypeInternal)
	}
	return w, nil
}

func (r *PostgresRepo) DeleteWarehouse(ctx context.Context, tenantID kernel.TenantID, id kernel.WarehouseID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM warehouses WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting warehouse", errx.TypeInternal)
	}
	return nil
}

// ─── Stock levels ─────────────────────────────────────────────────────────────

func (r *PostgresRepo) GetStockLevel(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, variantID *kernel.VariantID, warehouseID kernel.WarehouseID) (inventory.StockLevel, error) {
	const q = `
		SELECT id, tenant_id, product_id, variant_id, warehouse_id,
		       quantity, reserved, low_stock_threshold, created_at, updated_at
		FROM stock_levels
		WHERE tenant_id = $1 AND product_id = $2 AND warehouse_id = $3
		  AND ($4::uuid IS NULL AND variant_id IS NULL OR variant_id = $4::uuid)`

	var varIDStr *string
	if variantID != nil {
		s := string(*variantID)
		varIDStr = &s
	}

	var sl inventory.StockLevel
	var slID, slTenantID, slProductID, slWarehouseID string
	var slVariantID *string

	err := r.db.QueryRowContext(ctx, q, string(tenantID), string(productID), string(warehouseID), varIDStr).Scan(
		&slID, &slTenantID, &slProductID, &slVariantID, &slWarehouseID,
		&sl.Quantity, &sl.Reserved, &sl.LowStockThreshold, &sl.CreatedAt, &sl.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return inventory.StockLevel{}, inventory.ErrStockLevelNotFound
	}
	if err != nil {
		return inventory.StockLevel{}, errx.Wrap(err, "scanning stock level", errx.TypeInternal)
	}

	sl.ID = kernel.StockLevelID(slID)
	sl.TenantID = kernel.TenantID(slTenantID)
	sl.ProductID = kernel.ProductID(slProductID)
	sl.WarehouseID = kernel.WarehouseID(slWarehouseID)
	if slVariantID != nil {
		vid := kernel.VariantID(*slVariantID)
		sl.VariantID = &vid
	}
	return sl, nil
}

func (r *PostgresRepo) ListStockLevels(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) ([]inventory.StockLevel, error) {
	const q = `
		SELECT id, tenant_id, product_id, variant_id, warehouse_id,
		       quantity, reserved, low_stock_threshold, created_at, updated_at
		FROM stock_levels
		WHERE tenant_id = $1 AND product_id = $2
		ORDER BY warehouse_id`

	rows, err := r.db.QueryContext(ctx, q, string(tenantID), string(productID))
	if err != nil {
		return nil, errx.Wrap(err, "querying stock levels", errx.TypeInternal)
	}
	defer rows.Close()

	var items []inventory.StockLevel
	for rows.Next() {
		sl, err := scanStockLevel(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning stock level row", errx.TypeInternal)
		}
		items = append(items, sl)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating stock levels", errx.TypeInternal)
	}
	return items, nil
}

func (r *PostgresRepo) UpsertStockLevel(ctx context.Context, sl inventory.StockLevel) (inventory.StockLevel, error) {
	const q = `
		INSERT INTO stock_levels (id, tenant_id, product_id, variant_id, warehouse_id,
		                          quantity, reserved, low_stock_threshold, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (tenant_id, product_id, COALESCE(variant_id, '00000000-0000-0000-0000-000000000000'::uuid), warehouse_id)
		DO UPDATE SET quantity = EXCLUDED.quantity,
		              reserved = EXCLUDED.reserved,
		              low_stock_threshold = EXCLUDED.low_stock_threshold,
		              updated_at = EXCLUDED.updated_at`

	var varIDStr *string
	if sl.VariantID != nil {
		s := string(*sl.VariantID)
		varIDStr = &s
	}

	_, err := r.db.ExecContext(ctx, q,
		string(sl.ID), string(sl.TenantID), string(sl.ProductID), varIDStr, string(sl.WarehouseID),
		sl.Quantity, sl.Reserved, sl.LowStockThreshold, sl.CreatedAt, sl.UpdatedAt,
	)
	if err != nil {
		return inventory.StockLevel{}, errx.Wrap(err, "upserting stock level", errx.TypeInternal)
	}
	return sl, nil
}

func (r *PostgresRepo) GetLowStockItems(ctx context.Context, tenantID kernel.TenantID) ([]inventory.StockLevel, error) {
	const q = `
		SELECT id, tenant_id, product_id, variant_id, warehouse_id,
		       quantity, reserved, low_stock_threshold, created_at, updated_at
		FROM stock_levels
		WHERE tenant_id = $1
		  AND (quantity - reserved) <= low_stock_threshold
		ORDER BY (quantity - reserved) ASC`

	rows, err := r.db.QueryContext(ctx, q, string(tenantID))
	if err != nil {
		return nil, errx.Wrap(err, "querying low stock items", errx.TypeInternal)
	}
	defer rows.Close()

	var items []inventory.StockLevel
	for rows.Next() {
		sl, err := scanStockLevel(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning low stock row", errx.TypeInternal)
		}
		items = append(items, sl)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating low stock items", errx.TypeInternal)
	}
	return items, nil
}

// ─── Stock movements ──────────────────────────────────────────────────────────

func (r *PostgresRepo) CreateMovement(ctx context.Context, m inventory.StockMovement) (inventory.StockMovement, error) {
	const q = `
		INSERT INTO stock_movements (id, tenant_id, product_id, variant_id, warehouse_id,
		                             type, quantity, reference, note, created_by, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`

	var varIDStr *string
	if m.VariantID != nil {
		s := string(*m.VariantID)
		varIDStr = &s
	}

	_, err := r.db.ExecContext(ctx, q,
		string(m.ID), string(m.TenantID), string(m.ProductID), varIDStr, string(m.WarehouseID),
		string(m.Type), m.Quantity, m.Reference, m.Note, m.CreatedBy, m.CreatedAt,
	)
	if err != nil {
		return inventory.StockMovement{}, errx.Wrap(err, "inserting stock movement", errx.TypeInternal)
	}
	return m, nil
}

func (r *PostgresRepo) ListMovements(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, pg kernel.PaginationOptions) (kernel.Paginated[inventory.StockMovement], error) {
	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM stock_movements WHERE tenant_id = $1 AND product_id = $2`,
		string(tenantID), string(productID),
	).Scan(&total); err != nil {
		return kernel.Paginated[inventory.StockMovement]{}, errx.Wrap(err, "counting stock movements", errx.TypeInternal)
	}

	const q = `
		SELECT id, tenant_id, product_id, variant_id, warehouse_id,
		       type, quantity, reference, note, created_by, created_at
		FROM stock_movements
		WHERE tenant_id = $1 AND product_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, q, string(tenantID), string(productID), pg.Limit(), pg.Offset())
	if err != nil {
		return kernel.Paginated[inventory.StockMovement]{}, errx.Wrap(err, "querying stock movements", errx.TypeInternal)
	}
	defer rows.Close()

	var items []inventory.StockMovement
	for rows.Next() {
		m, err := scanMovement(rows)
		if err != nil {
			return kernel.Paginated[inventory.StockMovement]{}, errx.Wrap(err, "scanning stock movement row", errx.TypeInternal)
		}
		items = append(items, m)
	}
	if err := rows.Err(); err != nil {
		return kernel.Paginated[inventory.StockMovement]{}, errx.Wrap(err, "iterating stock movements", errx.TypeInternal)
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

// ─── scan helpers ─────────────────────────────────────────────────────────────

type rowScanner interface {
	Scan(dest ...any) error
}

func scanStockLevel(s rowScanner) (inventory.StockLevel, error) {
	var sl inventory.StockLevel
	var slID, slTenantID, slProductID, slWarehouseID string
	var slVariantID *string

	if err := s.Scan(
		&slID, &slTenantID, &slProductID, &slVariantID, &slWarehouseID,
		&sl.Quantity, &sl.Reserved, &sl.LowStockThreshold, &sl.CreatedAt, &sl.UpdatedAt,
	); err != nil {
		return inventory.StockLevel{}, err
	}

	sl.ID = kernel.StockLevelID(slID)
	sl.TenantID = kernel.TenantID(slTenantID)
	sl.ProductID = kernel.ProductID(slProductID)
	sl.WarehouseID = kernel.WarehouseID(slWarehouseID)
	if slVariantID != nil {
		vid := kernel.VariantID(*slVariantID)
		sl.VariantID = &vid
	}
	return sl, nil
}

func scanMovement(s rowScanner) (inventory.StockMovement, error) {
	var m inventory.StockMovement
	var mID, mTenantID, mProductID, mWarehouseID, mType string
	var mVariantID *string

	if err := s.Scan(
		&mID, &mTenantID, &mProductID, &mVariantID, &mWarehouseID,
		&mType, &m.Quantity, &m.Reference, &m.Note, &m.CreatedBy, &m.CreatedAt,
	); err != nil {
		return inventory.StockMovement{}, err
	}

	m.ID = kernel.StockMovementID(mID)
	m.TenantID = kernel.TenantID(mTenantID)
	m.ProductID = kernel.ProductID(mProductID)
	m.WarehouseID = kernel.WarehouseID(mWarehouseID)
	m.Type = inventory.MovementType(mType)
	if mVariantID != nil {
		vid := kernel.VariantID(*mVariantID)
		m.VariantID = &vid
	}
	return m, nil
}

// Compile-time check that PostgresRepo satisfies the repository interface.
var _ inventory.Repository = (*PostgresRepo)(nil)
