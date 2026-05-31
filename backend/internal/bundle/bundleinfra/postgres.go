package bundleinfra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/bundle"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// PostgresRepository implements bundle.Repository using sqlx.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgreSQL-backed bundle repository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// ─── Bundle CRUD ──────────────────────────────────────────────────────────────

func (r *PostgresRepository) Create(ctx context.Context, b *bundle.Bundle) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO bundles
			(id, tenant_id, name, slug, description, discount_type, discount_value, active, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(b.ID), string(b.TenantID),
		b.Name, b.Slug, b.Description,
		string(b.DiscountType), b.DiscountValue, b.Active,
		b.CreatedAt, b.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting bundle", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BundleID) (*bundle.Bundle, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, description, discount_type, discount_value, active, created_at, updated_at
		FROM bundles
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	b, err := scanBundle(row)
	if err == sql.ErrNoRows {
		return nil, bundle.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning bundle", errx.TypeInternal)
	}
	return b, nil
}

func (r *PostgresRepository) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*bundle.Bundle, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, description, discount_type, discount_value, active, created_at, updated_at
		FROM bundles
		WHERE slug = $1 AND tenant_id = $2`,
		slug, string(tenantID),
	)
	b, err := scanBundle(row)
	if err == sql.ErrNoRows {
		return nil, bundle.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning bundle by slug", errx.TypeInternal)
	}
	return b, nil
}

func (r *PostgresRepository) Update(ctx context.Context, b *bundle.Bundle) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE bundles
		SET name=$1, slug=$2, description=$3, discount_type=$4, discount_value=$5, active=$6, updated_at=$7
		WHERE id=$8 AND tenant_id=$9`,
		b.Name, b.Slug, b.Description,
		string(b.DiscountType), b.DiscountValue, b.Active, b.UpdatedAt,
		string(b.ID), string(b.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating bundle", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.BundleID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM bundles WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting bundle", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepository) List(ctx context.Context, tenantID kernel.TenantID, activeOnly bool, pg kernel.PaginationOptions) (kernel.Paginated[bundle.Bundle], error) {
	var zero kernel.Paginated[bundle.Bundle]

	whereExtra := ""
	args := []any{string(tenantID)}
	if activeOnly {
		whereExtra = " AND active = true"
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM bundles WHERE tenant_id = $1%s", whereExtra)
	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting bundles", errx.TypeInternal)
	}

	dataQuery := fmt.Sprintf(`
		SELECT id, tenant_id, name, slug, description, discount_type, discount_value, active, created_at, updated_at
		FROM bundles
		WHERE tenant_id = $1%s
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`, whereExtra)

	rows, err := r.db.QueryContext(ctx, dataQuery, string(tenantID), pg.Limit(), pg.Offset())
	if err != nil {
		return zero, errx.Wrap(err, "querying bundles", errx.TypeInternal)
	}
	defer rows.Close()

	var bundles []bundle.Bundle
	for rows.Next() {
		b, err := scanBundleRow(rows)
		if err != nil {
			return zero, errx.Wrap(err, "scanning bundle row", errx.TypeInternal)
		}
		bundles = append(bundles, *b)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating bundles", errx.TypeInternal)
	}

	return kernel.NewPaginated(bundles, pg.Page, pg.PageSize, total), nil
}

// ─── Bundle Items ──────────────────────────────────────────────────────────────

func (r *PostgresRepository) AddItem(ctx context.Context, item *bundle.BundleItem) error {
	var variantID *string
	if item.VariantID != nil {
		s := string(*item.VariantID)
		variantID = &s
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO bundle_items (id, tenant_id, bundle_id, product_id, variant_id, quantity, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		string(item.ID), string(item.TenantID),
		string(item.BundleID), string(item.ProductID),
		variantID, item.Quantity, item.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting bundle item", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepository) GetItemByID(ctx context.Context, tenantID kernel.TenantID, itemID kernel.BundleItemID) (*bundle.BundleItem, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, bundle_id, product_id, variant_id, quantity, created_at
		FROM bundle_items
		WHERE id = $1 AND tenant_id = $2`,
		string(itemID), string(tenantID),
	)
	item, err := scanBundleItem(row)
	if err == sql.ErrNoRows {
		return nil, bundle.ErrItemNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning bundle item", errx.TypeInternal)
	}
	return item, nil
}

func (r *PostgresRepository) ListItems(ctx context.Context, tenantID kernel.TenantID, bundleID kernel.BundleID) ([]bundle.BundleItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, bundle_id, product_id, variant_id, quantity, created_at
		FROM bundle_items
		WHERE bundle_id = $1 AND tenant_id = $2
		ORDER BY created_at ASC`,
		string(bundleID), string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying bundle items", errx.TypeInternal)
	}
	defer rows.Close()

	var items []bundle.BundleItem
	for rows.Next() {
		item, err := scanBundleItemRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning bundle item row", errx.TypeInternal)
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating bundle items", errx.TypeInternal)
	}

	if items == nil {
		items = []bundle.BundleItem{}
	}

	return items, nil
}

func (r *PostgresRepository) RemoveItem(ctx context.Context, tenantID kernel.TenantID, itemID kernel.BundleItemID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM bundle_items WHERE id = $1 AND tenant_id = $2`,
		string(itemID), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting bundle item", errx.TypeInternal)
	}
	return nil
}

// ─── Scan helpers ─────────────────────────────────────────────────────────────

type rowScanner interface {
	Scan(dest ...any) error
}

func scanBundle(row *sql.Row) (*bundle.Bundle, error) {
	return scanBundleFields(row)
}

func scanBundleRow(rows interface{ Scan(dest ...any) error }) (*bundle.Bundle, error) {
	return scanBundleFields(rows)
}

func scanBundleFields(s rowScanner) (*bundle.Bundle, error) {
	var b bundle.Bundle
	var id, tenantID, discountType string
	err := s.Scan(
		&id, &tenantID, &b.Name, &b.Slug, &b.Description,
		&discountType, &b.DiscountValue, &b.Active,
		&b.CreatedAt, &b.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	b.ID = kernel.BundleID(id)
	b.TenantID = kernel.TenantID(tenantID)
	b.DiscountType = bundle.DiscountType(discountType)
	return &b, nil
}

func scanBundleItem(row *sql.Row) (*bundle.BundleItem, error) {
	return scanBundleItemFields(row)
}

func scanBundleItemRow(rows interface{ Scan(dest ...any) error }) (*bundle.BundleItem, error) {
	return scanBundleItemFields(rows)
}

func scanBundleItemFields(s rowScanner) (*bundle.BundleItem, error) {
	var item bundle.BundleItem
	var id, tenantID, bundleID, productID string
	var variantID *string

	err := s.Scan(&id, &tenantID, &bundleID, &productID, &variantID, &item.Quantity, &item.CreatedAt)
	if err != nil {
		return nil, err
	}

	item.ID = kernel.BundleItemID(id)
	item.TenantID = kernel.TenantID(tenantID)
	item.BundleID = kernel.BundleID(bundleID)
	item.ProductID = kernel.ProductID(productID)

	if variantID != nil {
		vid := kernel.VariantID(*variantID)
		item.VariantID = &vid
	}

	return &item, nil
}

// Compile-time interface check.
var _ bundle.Repository = (*PostgresRepository)(nil)
