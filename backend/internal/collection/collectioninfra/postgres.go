package collectioninfra

import (
	"context"
	stdsql "database/sql"
	"encoding/json"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/Abraxas-365/vendex/internal/collection"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// PostgresRepo implements collection.Repository using sqlx / PostgreSQL.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgresRepo.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Compile-time interface check.
var _ collection.Repository = (*PostgresRepo)(nil)

// --------------------------------------------------------------------------
// Collection CRUD
// --------------------------------------------------------------------------

func (r *PostgresRepo) Create(ctx context.Context, c *collection.Collection) error {
	rulesJSON, err := json.Marshal(c.Rules)
	if err != nil {
		return errx.Wrap(err, "marshaling rules", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO collections (
			id, tenant_id, name, slug, description, image_url,
			type, rules, is_active, sort_order,
			meta_title, meta_description, published_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13,
			$14, $15
		)`,
		string(c.ID), string(c.TenantID), c.Name, c.Slug, c.Description,
		nullableString(c.ImageURL),
		string(c.Type), string(rulesJSON), c.IsActive, c.SortOrder,
		nullableString(c.MetaTitle), nullableString(c.MetaDescription),
		c.PublishedAt,
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return collection.ErrDuplicateSlug
		}
		return errx.Wrap(err, "inserting collection", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) (*collection.Collection, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, description, COALESCE(image_url,''),
		       COALESCE(type,'manual'), COALESCE(rules,'[]'::text), is_active,
		       COALESCE(sort_order,0), COALESCE(meta_title,''), COALESCE(meta_description,''),
		       published_at, created_at, updated_at
		FROM collections
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	return scanCollection(row)
}

func (r *PostgresRepo) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*collection.Collection, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, description, COALESCE(image_url,''),
		       COALESCE(type,'manual'), COALESCE(rules,'[]'::text), is_active,
		       COALESCE(sort_order,0), COALESCE(meta_title,''), COALESCE(meta_description,''),
		       published_at, created_at, updated_at
		FROM collections
		WHERE slug = $1 AND tenant_id = $2`,
		slug, string(tenantID),
	)
	return scanCollection(row)
}

func (r *PostgresRepo) Update(ctx context.Context, c *collection.Collection) error {
	rulesJSON, err := json.Marshal(c.Rules)
	if err != nil {
		return errx.Wrap(err, "marshaling rules", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE collections
		SET name=$1, slug=$2, description=$3, image_url=$4,
		    type=$5, rules=$6, is_active=$7, sort_order=$8,
		    meta_title=$9, meta_description=$10, published_at=$11,
		    updated_at=$12
		WHERE id=$13 AND tenant_id=$14`,
		c.Name, c.Slug, c.Description, nullableString(c.ImageURL),
		string(c.Type), string(rulesJSON), c.IsActive, c.SortOrder,
		nullableString(c.MetaTitle), nullableString(c.MetaDescription),
		c.PublishedAt, c.UpdatedAt,
		string(c.ID), string(c.TenantID),
	)
	if err != nil {
		if isUniqueViolation(err) {
			return collection.ErrDuplicateSlug
		}
		return errx.Wrap(err, "updating collection", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM collections WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting collection", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, activeOnly bool, pg kernel.PaginationOptions) (kernel.Paginated[collection.Collection], error) {
	var zero kernel.Paginated[collection.Collection]

	where := "WHERE tenant_id = $1"
	args := []any{string(tenantID)}
	if activeOnly {
		where += " AND is_active = TRUE"
	}

	var total int
	if err := r.db.QueryRowContext(ctx, fmt.Sprintf("SELECT COUNT(*) FROM collections %s", where), args...).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting collections", errx.TypeInternal)
	}

	args = append(args, pg.Limit(), pg.Offset())
	rows, err := r.db.QueryContext(ctx, fmt.Sprintf(`
		SELECT id, tenant_id, name, slug, description, COALESCE(image_url,''),
		       COALESCE(type,'manual'), COALESCE(rules,'[]'::text), is_active,
		       COALESCE(sort_order,0), COALESCE(meta_title,''), COALESCE(meta_description,''),
		       published_at, created_at, updated_at
		FROM collections
		%s
		ORDER BY sort_order ASC, created_at DESC
		LIMIT $%d OFFSET $%d`, where, len(args)-1, len(args)),
		args...,
	)
	if err != nil {
		return zero, errx.Wrap(err, "listing collections", errx.TypeInternal)
	}
	defer rows.Close()

	var items []collection.Collection
	for rows.Next() {
		c, err := scanCollectionRow(rows)
		if err != nil {
			return zero, err
		}
		items = append(items, *c)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating collections", errx.TypeInternal)
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

func (r *PostgresRepo) CountProducts(ctx context.Context, tenantID kernel.TenantID, collectionID kernel.CollectionID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM collection_products WHERE collection_id = $1 AND tenant_id = $2`,
		string(collectionID), string(tenantID),
	).Scan(&count)
	if err != nil {
		return 0, errx.Wrap(err, "counting collection products", errx.TypeInternal)
	}
	return count, nil
}

// --------------------------------------------------------------------------
// Product membership
// --------------------------------------------------------------------------

func (r *PostgresRepo) AddProduct(ctx context.Context, cp *collection.CollectionProduct) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO collection_products (id, tenant_id, collection_id, product_id, sort_order, added_at)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		string(cp.ID), string(cp.TenantID),
		string(cp.CollectionID), cp.ProductID,
		cp.SortOrder, cp.AddedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return collection.ErrAlreadyInCollection
		}
		return errx.Wrap(err, "adding product to collection", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) RemoveProduct(ctx context.Context, tenantID kernel.TenantID, collectionID kernel.CollectionID, productID string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM collection_products
		WHERE tenant_id = $1 AND collection_id = $2 AND product_id = $3`,
		string(tenantID), string(collectionID), productID,
	)
	if err != nil {
		return errx.Wrap(err, "removing product from collection", errx.TypeInternal)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return collection.ErrProductNotFound
	}
	return nil
}

func (r *PostgresRepo) ListProducts(ctx context.Context, tenantID kernel.TenantID, collectionID kernel.CollectionID, pg kernel.PaginationOptions) (kernel.Paginated[collection.CollectionProduct], error) {
	var zero kernel.Paginated[collection.CollectionProduct]

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM collection_products WHERE collection_id = $1 AND tenant_id = $2`,
		string(collectionID), string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting collection products", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, collection_id, product_id, sort_order, added_at
		FROM collection_products
		WHERE collection_id = $1 AND tenant_id = $2
		ORDER BY sort_order ASC, added_at ASC
		LIMIT $3 OFFSET $4`,
		string(collectionID), string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "listing collection products", errx.TypeInternal)
	}
	defer rows.Close()

	var items []collection.CollectionProduct
	for rows.Next() {
		var cp collection.CollectionProduct
		var id, tenantIDStr, collectionIDStr string
		if err := rows.Scan(&id, &tenantIDStr, &collectionIDStr, &cp.ProductID, &cp.SortOrder, &cp.AddedAt); err != nil {
			return zero, errx.Wrap(err, "scanning collection product", errx.TypeInternal)
		}
		cp.ID = kernel.CollectionProductID(id)
		cp.TenantID = kernel.TenantID(tenantIDStr)
		cp.CollectionID = kernel.CollectionID(collectionIDStr)
		items = append(items, cp)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating collection products", errx.TypeInternal)
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

func (r *PostgresRepo) ReorderProducts(ctx context.Context, tenantID kernel.TenantID, collectionID kernel.CollectionID, productIDs []string) error {
	for i, productID := range productIDs {
		_, err := r.db.ExecContext(ctx, `
			UPDATE collection_products
			SET sort_order = $1
			WHERE tenant_id = $2 AND collection_id = $3 AND product_id = $4`,
			i, string(tenantID), string(collectionID), productID,
		)
		if err != nil {
			return errx.Wrap(err, fmt.Sprintf("reordering product %s", productID), errx.TypeInternal)
		}
	}
	return nil
}

// --------------------------------------------------------------------------
// Helpers
// --------------------------------------------------------------------------

type rowScanner interface {
	Scan(dest ...any) error
}

func scanCollection(row *stdsql.Row) (*collection.Collection, error) {
	var c collection.Collection
	var id, tenantID, colType, rulesJSON string

	err := row.Scan(
		&id, &tenantID, &c.Name, &c.Slug, &c.Description,
		&c.ImageURL, &colType, &rulesJSON,
		&c.IsActive, &c.SortOrder,
		&c.MetaTitle, &c.MetaDescription,
		&c.PublishedAt,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err == stdsql.ErrNoRows {
		return nil, collection.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning collection", errx.TypeInternal)
	}

	c.ID = kernel.CollectionID(id)
	c.TenantID = kernel.TenantID(tenantID)
	c.Type = collection.CollectionType(colType)
	if err := json.Unmarshal([]byte(rulesJSON), &c.Rules); err != nil || c.Rules == nil {
		c.Rules = []collection.CollectionRule{}
	}
	return &c, nil
}

func scanCollectionRow(s rowScanner) (*collection.Collection, error) {
	var c collection.Collection
	var id, tenantID, colType, rulesJSON string

	err := s.Scan(
		&id, &tenantID, &c.Name, &c.Slug, &c.Description,
		&c.ImageURL, &colType, &rulesJSON,
		&c.IsActive, &c.SortOrder,
		&c.MetaTitle, &c.MetaDescription,
		&c.PublishedAt,
		&c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, errx.Wrap(err, "scanning collection row", errx.TypeInternal)
	}

	c.ID = kernel.CollectionID(id)
	c.TenantID = kernel.TenantID(tenantID)
	c.Type = collection.CollectionType(colType)
	if err := json.Unmarshal([]byte(rulesJSON), &c.Rules); err != nil || c.Rules == nil {
		c.Rules = []collection.CollectionRule{}
	}
	return &c, nil
}

// nullableString returns a *string so that empty strings are stored as NULL.
func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// isUniqueViolation returns true when err is a PostgreSQL unique-constraint error (code 23505).
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}
