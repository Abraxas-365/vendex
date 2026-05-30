package productinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/product"
	"github.com/jmoiron/sqlx"
)

// VariantPostgresRepo implements product.VariantRepository using sqlx.
type VariantPostgresRepo struct {
	db *sqlx.DB
}

// NewVariantPostgresRepo creates a new PostgreSQL-backed variant repository.
func NewVariantPostgresRepo(db *sqlx.DB) *VariantPostgresRepo {
	return &VariantPostgresRepo{db: db}
}

// ─── Options ─────────────────────────────────────────────────────────────────

func (r *VariantPostgresRepo) CreateOption(ctx context.Context, opt *product.ProductOption) error {
	valuesJSON, err := json.Marshal(opt.Values)
	if err != nil {
		return errx.Wrap(err, "marshaling option values", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO product_options (id, product_id, tenant_id, name, position, values, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(opt.ID), string(opt.ProductID), string(opt.TenantID),
		opt.Name, opt.Position, string(valuesJSON),
		opt.CreatedAt, opt.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting product option", errx.TypeInternal)
	}
	return nil
}

func (r *VariantPostgresRepo) ListOptions(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) ([]product.ProductOption, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, product_id, tenant_id, name, position, values, created_at, updated_at
		FROM product_options
		WHERE tenant_id = $1 AND product_id = $2
		ORDER BY position ASC, created_at ASC`,
		string(tenantID), string(productID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying product options", errx.TypeInternal)
	}
	defer rows.Close()

	var opts []product.ProductOption
	for rows.Next() {
		opt, err := scanOption(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning product option row", errx.TypeInternal)
		}
		opts = append(opts, *opt)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating product options", errx.TypeInternal)
	}
	if opts == nil {
		opts = []product.ProductOption{}
	}
	return opts, nil
}

func (r *VariantPostgresRepo) UpdateOption(ctx context.Context, opt *product.ProductOption) error {
	valuesJSON, err := json.Marshal(opt.Values)
	if err != nil {
		return errx.Wrap(err, "marshaling option values", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE product_options
		SET name=$1, position=$2, values=$3, updated_at=$4
		WHERE id=$5 AND tenant_id=$6`,
		opt.Name, opt.Position, string(valuesJSON), opt.UpdatedAt,
		string(opt.ID), string(opt.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating product option", errx.TypeInternal)
	}
	return nil
}

func (r *VariantPostgresRepo) DeleteOption(ctx context.Context, tenantID kernel.TenantID, id kernel.OptionID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM product_options WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting product option", errx.TypeInternal)
	}
	return nil
}

// ─── Variants ─────────────────────────────────────────────────────────────────

func (r *VariantPostgresRepo) CreateVariant(ctx context.Context, v *product.ProductVariant) error {
	optionsJSON, err := json.Marshal(v.Options)
	if err != nil {
		return errx.Wrap(err, "marshaling variant options", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO product_variants (id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		string(v.ID), string(v.ProductID), string(v.TenantID),
		v.SKU, v.Price.Amount, v.Price.Currency,
		v.Stock, string(optionsJSON), v.Active,
		v.CreatedAt, v.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting product variant", errx.TypeInternal)
	}
	return nil
}

func (r *VariantPostgresRepo) GetVariantByID(ctx context.Context, tenantID kernel.TenantID, id kernel.VariantID) (*product.ProductVariant, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active, created_at, updated_at
		FROM product_variants WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	v, err := scanVariantRow(row)
	if err == sql.ErrNoRows {
		return nil, product.ErrVariantNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning product variant", errx.TypeInternal)
	}
	return v, nil
}

func (r *VariantPostgresRepo) ListVariants(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) ([]product.ProductVariant, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active, created_at, updated_at
		FROM product_variants
		WHERE tenant_id=$1 AND product_id=$2
		ORDER BY created_at ASC`,
		string(tenantID), string(productID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying product variants", errx.TypeInternal)
	}
	defer rows.Close()

	var variants []product.ProductVariant
	for rows.Next() {
		v, err := scanVariantRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning product variant row", errx.TypeInternal)
		}
		variants = append(variants, *v)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating product variants", errx.TypeInternal)
	}
	if variants == nil {
		variants = []product.ProductVariant{}
	}
	return variants, nil
}

func (r *VariantPostgresRepo) UpdateVariant(ctx context.Context, v *product.ProductVariant) error {
	optionsJSON, err := json.Marshal(v.Options)
	if err != nil {
		return errx.Wrap(err, "marshaling variant options", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE product_variants
		SET sku=$1, price_amount=$2, price_currency=$3, stock=$4, options=$5, active=$6, updated_at=$7
		WHERE id=$8 AND tenant_id=$9`,
		v.SKU, v.Price.Amount, v.Price.Currency,
		v.Stock, string(optionsJSON), v.Active, v.UpdatedAt,
		string(v.ID), string(v.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating product variant", errx.TypeInternal)
	}
	return nil
}

func (r *VariantPostgresRepo) DeleteVariant(ctx context.Context, tenantID kernel.TenantID, id kernel.VariantID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM product_variants WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting product variant", errx.TypeInternal)
	}
	return nil
}

func (r *VariantPostgresRepo) GetVariantBySKU(ctx context.Context, tenantID kernel.TenantID, sku string) (*product.ProductVariant, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active, created_at, updated_at
		FROM product_variants WHERE sku=$1 AND tenant_id=$2`,
		sku, string(tenantID),
	)
	v, err := scanVariantRow(row)
	if err == sql.ErrNoRows {
		return nil, product.ErrVariantNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning product variant by SKU", errx.TypeInternal)
	}
	return v, nil
}

// ─── Scan helpers ─────────────────────────────────────────────────────────────

type variantScanner interface {
	Scan(dest ...any) error
}

func scanOption(s variantScanner) (*product.ProductOption, error) {
	var opt product.ProductOption
	var id, productID, tenantID string
	var valuesJSON string

	err := s.Scan(
		&id, &productID, &tenantID,
		&opt.Name, &opt.Position, &valuesJSON,
		&opt.CreatedAt, &opt.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	opt.ID = kernel.OptionID(id)
	opt.ProductID = kernel.ProductID(productID)
	opt.TenantID = kernel.TenantID(tenantID)

	_ = json.Unmarshal([]byte(valuesJSON), &opt.Values)
	if opt.Values == nil {
		opt.Values = []string{}
	}

	return &opt, nil
}

func scanVariantRow(s variantScanner) (*product.ProductVariant, error) {
	var v product.ProductVariant
	var id, productID, tenantID string
	var optionsJSON string

	// Use placeholder for created_at and updated_at to map to time.Time
	var createdAt, updatedAt time.Time

	err := s.Scan(
		&id, &productID, &tenantID,
		&v.SKU, &v.Price.Amount, &v.Price.Currency,
		&v.Stock, &optionsJSON, &v.Active,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}

	v.ID = kernel.VariantID(id)
	v.ProductID = kernel.ProductID(productID)
	v.TenantID = kernel.TenantID(tenantID)
	v.CreatedAt = createdAt
	v.UpdatedAt = updatedAt

	_ = json.Unmarshal([]byte(optionsJSON), &v.Options)
	if v.Options == nil {
		v.Options = map[string]string{}
	}

	return &v, nil
}

// Ensure interface compliance.
var _ product.VariantRepository = (*VariantPostgresRepo)(nil)
