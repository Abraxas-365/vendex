package productinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/product"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements product.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed product repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, p *product.Product) error {
	imagesJSON, err := json.Marshal(p.Images)
	if err != nil {
		return errx.Wrap(err, "marshaling images", errx.TypeInternal)
	}
	tagsJSON, err := json.Marshal(p.Tags)
	if err != nil {
		return errx.Wrap(err, "marshaling tags", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO products (id, tenant_id, name, description, price_amount, price_currency, sku, images, category_id, tags, status, stock, has_variants, meta_title, meta_description, slug, canonical_url, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)`,
		string(p.ID), string(p.TenantID), p.Name, p.Description,
		p.Price.Amount, p.Price.Currency, p.SKU,
		string(imagesJSON), string(p.CategoryID), string(tagsJSON),
		string(p.Status), p.Stock, p.HasVariants,
		p.MetaTitle, p.MetaDescription, p.Slug, p.CanonicalURL,
		p.CreatedAt, p.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting product", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) (*product.Product, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, description, price_amount, price_currency, sku, images, category_id, tags, status, stock, has_variants, meta_title, meta_description, slug, canonical_url, created_at, updated_at
		FROM products WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	p, err := scanProduct(row)
	if err == sql.ErrNoRows {
		return nil, product.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning product", errx.TypeInternal)
	}
	return p, nil
}

func (r *PostgresRepo) GetBySKU(ctx context.Context, tenantID kernel.TenantID, sku string) (*product.Product, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, description, price_amount, price_currency, sku, images, category_id, tags, status, stock, has_variants, meta_title, meta_description, slug, canonical_url, created_at, updated_at
		FROM products WHERE sku = $1 AND tenant_id = $2`,
		sku, string(tenantID),
	)
	p, err := scanProduct(row)
	if err == sql.ErrNoRows {
		return nil, product.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning product by SKU", errx.TypeInternal)
	}
	return p, nil
}

func (r *PostgresRepo) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*product.Product, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, description, price_amount, price_currency, sku, images, category_id, tags, status, stock, has_variants, meta_title, meta_description, slug, canonical_url, created_at, updated_at
		FROM products WHERE slug = $1 AND tenant_id = $2`,
		slug, string(tenantID),
	)
	p, err := scanProduct(row)
	if err == sql.ErrNoRows {
		return nil, product.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning product by slug", errx.TypeInternal)
	}
	return p, nil
}

func (r *PostgresRepo) Update(ctx context.Context, p *product.Product) error {
	imagesJSON, err := json.Marshal(p.Images)
	if err != nil {
		return errx.Wrap(err, "marshaling images", errx.TypeInternal)
	}
	tagsJSON, err := json.Marshal(p.Tags)
	if err != nil {
		return errx.Wrap(err, "marshaling tags", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE products SET name=$1, description=$2, price_amount=$3, price_currency=$4, sku=$5, images=$6, category_id=$7, tags=$8, status=$9, stock=$10, has_variants=$11, meta_title=$12, meta_description=$13, slug=$14, canonical_url=$15, updated_at=$16
		WHERE id=$17 AND tenant_id=$18`,
		p.Name, p.Description, p.Price.Amount, p.Price.Currency,
		p.SKU, string(imagesJSON), string(p.CategoryID), string(tagsJSON),
		string(p.Status), p.Stock, p.HasVariants,
		p.MetaTitle, p.MetaDescription, p.Slug, p.CanonicalURL, p.UpdatedAt,
		string(p.ID), string(p.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating product", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ProductID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM products WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting product", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[product.Product], error) {
	return r.queryProducts(ctx, tenantID, pg, "", nil)
}

func (r *PostgresRepo) ListByCategory(ctx context.Context, tenantID kernel.TenantID, categoryID kernel.CategoryID, pg kernel.PaginationOptions) (kernel.Paginated[product.Product], error) {
	return r.queryProducts(ctx, tenantID, pg, "AND category_id = $3", []any{string(categoryID)})
}

func (r *PostgresRepo) queryProducts(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions, extraWhere string, extraArgs []any) (kernel.Paginated[product.Product], error) {
	var zero kernel.Paginated[product.Product]

	baseArgs := []any{string(tenantID)}
	countQuery := "SELECT COUNT(*) FROM products WHERE tenant_id = $1 " + extraWhere
	args := append(baseArgs, extraArgs...)

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting products", errx.TypeInternal)
	}

	nextParam := len(args) + 1
	dataQuery := fmt.Sprintf(`
		SELECT id, tenant_id, name, description, price_amount, price_currency, sku, images, category_id, tags, status, stock, has_variants, meta_title, meta_description, slug, canonical_url, created_at, updated_at
		FROM products WHERE tenant_id = $1 %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, extraWhere, nextParam, nextParam+1)
	args = append(args, pg.Limit(), pg.Offset())

	rows, err := r.db.QueryContext(ctx, dataQuery, args...)
	if err != nil {
		return zero, errx.Wrap(err, "querying products", errx.TypeInternal)
	}
	defer rows.Close()

	var products []product.Product
	for rows.Next() {
		p, err := scanProductRow(rows)
		if err != nil {
			return zero, errx.Wrap(err, "scanning product row", errx.TypeInternal)
		}
		products = append(products, *p)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating products", errx.TypeInternal)
	}

	return kernel.NewPaginated(products, pg.Page, pg.PageSize, total), nil
}

// scanner is implemented by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanProduct(row *sql.Row) (*product.Product, error) {
	return scanFields(row)
}

func scanProductRow(rows interface{ Scan(dest ...any) error }) (*product.Product, error) {
	return scanFields(rows)
}

func scanFields(s scanner) (*product.Product, error) {
	var p product.Product
	var id, tenantID, categoryID, status string
	var imagesJSON, tagsJSON string

	err := s.Scan(
		&id, &tenantID, &p.Name, &p.Description,
		&p.Price.Amount, &p.Price.Currency, &p.SKU,
		&imagesJSON, &categoryID, &tagsJSON,
		&status, &p.Stock, &p.HasVariants,
		&p.MetaTitle, &p.MetaDescription, &p.Slug, &p.CanonicalURL,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	p.ID = kernel.ProductID(id)
	p.TenantID = kernel.TenantID(tenantID)
	p.CategoryID = kernel.CategoryID(categoryID)
	p.Status = product.Status(status)

	_ = json.Unmarshal([]byte(imagesJSON), &p.Images)
	_ = json.Unmarshal([]byte(tagsJSON), &p.Tags)

	// Ensure nil slices become empty slices.
	if p.Images == nil {
		p.Images = []string{}
	}
	if p.Tags == nil {
		p.Tags = []string{}
	}

	return &p, nil
}

// Ensure interface compliance.
var _ product.Repository = (*PostgresRepo)(nil)
