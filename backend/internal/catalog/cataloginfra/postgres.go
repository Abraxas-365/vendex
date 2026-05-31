package cataloginfra

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/catalog"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// --- Category Repository ---

// CategoryPostgresRepo implements catalog.CategoryRepository.
type CategoryPostgresRepo struct {
	db *sqlx.DB
}

// NewCategoryPostgresRepo creates a new PostgreSQL-backed category repository.
func NewCategoryPostgresRepo(db *sqlx.DB) *CategoryPostgresRepo {
	return &CategoryPostgresRepo{db: db}
}

func (r *CategoryPostgresRepo) Create(ctx context.Context, c *catalog.Category) error {
	var parentID *string
	if c.ParentID != nil {
		s := string(*c.ParentID)
		parentID = &s
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO categories (id, tenant_id, name, slug, parent_id, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(c.ID), string(c.TenantID), c.Name, c.Slug,
		parentID, c.Description, c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting category", errx.TypeInternal)
	}
	return nil
}

func (r *CategoryPostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CategoryID) (*catalog.Category, error) {
	var c catalog.Category
	var parentID *string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, parent_id, description, created_at, updated_at
		FROM categories WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	).Scan(&c.ID, &c.TenantID, &c.Name, &c.Slug, &parentID, &c.Description, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, catalog.ErrCategoryNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "getting category by id", errx.TypeInternal)
	}

	if parentID != nil {
		pid := kernel.CategoryID(*parentID)
		c.ParentID = &pid
	}
	return &c, nil
}

func (r *CategoryPostgresRepo) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*catalog.Category, error) {
	var c catalog.Category
	var parentID *string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, parent_id, description, created_at, updated_at
		FROM categories WHERE slug = $1 AND tenant_id = $2`,
		slug, string(tenantID),
	).Scan(&c.ID, &c.TenantID, &c.Name, &c.Slug, &parentID, &c.Description, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, catalog.ErrCategoryNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "getting category by slug", errx.TypeInternal)
	}

	if parentID != nil {
		pid := kernel.CategoryID(*parentID)
		c.ParentID = &pid
	}
	return &c, nil
}

func (r *CategoryPostgresRepo) Update(ctx context.Context, c *catalog.Category) error {
	var parentID *string
	if c.ParentID != nil {
		s := string(*c.ParentID)
		parentID = &s
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE categories SET name=$1, slug=$2, parent_id=$3, description=$4, updated_at=$5
		WHERE id=$6 AND tenant_id=$7`,
		c.Name, c.Slug, parentID, c.Description, c.UpdatedAt,
		string(c.ID), string(c.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating category", errx.TypeInternal)
	}
	return nil
}

func (r *CategoryPostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CategoryID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting category", errx.TypeInternal)
	}
	return nil
}

func (r *CategoryPostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[catalog.Category], error) {
	var zero kernel.Paginated[catalog.Category]

	var total int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM categories WHERE tenant_id = $1", string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting categories", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, slug, parent_id, description, created_at, updated_at
		FROM categories WHERE tenant_id = $1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying categories", errx.TypeInternal)
	}
	defer rows.Close()

	var categories []catalog.Category
	for rows.Next() {
		c, err := scanCategoryFields(rows)
		if err != nil {
			return zero, err
		}
		categories = append(categories, *c)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating categories", errx.TypeInternal)
	}

	return kernel.NewPaginated(categories, pg.Page, pg.PageSize, total), nil
}

func (r *CategoryPostgresRepo) ListByParent(ctx context.Context, tenantID kernel.TenantID, parentID *kernel.CategoryID, pg kernel.PaginationOptions) (kernel.Paginated[catalog.Category], error) {
	var zero kernel.Paginated[catalog.Category]

	var total int
	var countArgs []any
	var countQuery string

	if parentID == nil {
		countQuery = "SELECT COUNT(*) FROM categories WHERE tenant_id = $1 AND parent_id IS NULL"
		countArgs = []any{string(tenantID)}
	} else {
		countQuery = "SELECT COUNT(*) FROM categories WHERE tenant_id = $1 AND parent_id = $2"
		countArgs = []any{string(tenantID), string(*parentID)}
	}

	if err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting categories by parent", errx.TypeInternal)
	}

	var dataQuery string
	var dataArgs []any

	if parentID == nil {
		dataQuery = `
			SELECT id, tenant_id, name, slug, parent_id, description, created_at, updated_at
			FROM categories WHERE tenant_id = $1 AND parent_id IS NULL
			ORDER BY name ASC
			LIMIT $2 OFFSET $3`
		dataArgs = []any{string(tenantID), pg.Limit(), pg.Offset()}
	} else {
		dataQuery = `
			SELECT id, tenant_id, name, slug, parent_id, description, created_at, updated_at
			FROM categories WHERE tenant_id = $1 AND parent_id = $2
			ORDER BY name ASC
			LIMIT $3 OFFSET $4`
		dataArgs = []any{string(tenantID), string(*parentID), pg.Limit(), pg.Offset()}
	}

	rows, err := r.db.QueryContext(ctx, dataQuery, dataArgs...)
	if err != nil {
		return zero, errx.Wrap(err, "querying categories by parent", errx.TypeInternal)
	}
	defer rows.Close()

	var categories []catalog.Category
	for rows.Next() {
		c, err := scanCategoryFields(rows)
		if err != nil {
			return zero, err
		}
		categories = append(categories, *c)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating categories by parent", errx.TypeInternal)
	}

	return kernel.NewPaginated(categories, pg.Page, pg.PageSize, total), nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanCategoryFields(s rowScanner) (*catalog.Category, error) {
	var c catalog.Category
	var id, tenantID string
	var parentID *string

	err := s.Scan(&id, &tenantID, &c.Name, &c.Slug, &parentID, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, errx.Wrap(err, "scanning category", errx.TypeInternal)
	}

	c.ID = kernel.CategoryID(id)
	c.TenantID = kernel.TenantID(tenantID)
	if parentID != nil {
		pid := kernel.CategoryID(*parentID)
		c.ParentID = &pid
	}

	return &c, nil
}

// Ensure interface compliance.
var _ catalog.CategoryRepository = (*CategoryPostgresRepo)(nil)

// --- Collection Repository ---

// CollectionPostgresRepo implements catalog.CollectionRepository.
type CollectionPostgresRepo struct {
	db *sqlx.DB
}

// NewCollectionPostgresRepo creates a new PostgreSQL-backed collection repository.
func NewCollectionPostgresRepo(db *sqlx.DB) *CollectionPostgresRepo {
	return &CollectionPostgresRepo{db: db}
}

func (r *CollectionPostgresRepo) Create(ctx context.Context, c *catalog.Collection) error {
	productIDsJSON, err := json.Marshal(c.ProductIDs)
	if err != nil {
		return errx.Wrap(err, "marshaling product IDs", errx.TypeInternal)
	}
	rulesJSON, err := json.Marshal(c.Rules)
	if err != nil {
		return errx.Wrap(err, "marshaling rules", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO collections (id, tenant_id, name, slug, description, product_ids, is_automatic, rules, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(c.ID), string(c.TenantID), c.Name, c.Slug, c.Description,
		string(productIDsJSON), c.IsAutomatic, string(rulesJSON),
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting collection", errx.TypeInternal)
	}
	return nil
}

func (r *CollectionPostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) (*catalog.Collection, error) {
	var c catalog.Collection
	var cID, cTenantID string
	var productIDsJSON, rulesJSON string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, description, product_ids, is_automatic, rules, created_at, updated_at
		FROM collections WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	).Scan(&cID, &cTenantID, &c.Name, &c.Slug, &c.Description,
		&productIDsJSON, &c.IsAutomatic, &rulesJSON, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, catalog.ErrCollectionNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "getting collection by id", errx.TypeInternal)
	}

	c.ID = kernel.CollectionID(cID)
	c.TenantID = kernel.TenantID(cTenantID)
	_ = json.Unmarshal([]byte(productIDsJSON), &c.ProductIDs)
	if c.ProductIDs == nil {
		c.ProductIDs = []kernel.ProductID{}
	}
	_ = json.Unmarshal([]byte(rulesJSON), &c.Rules)
	if c.Rules == nil {
		c.Rules = map[string]any{}
	}
	return &c, nil
}

func (r *CollectionPostgresRepo) GetBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*catalog.Collection, error) {
	var c catalog.Collection
	var cID, cTenantID string
	var productIDsJSON, rulesJSON string

	err := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, description, product_ids, is_automatic, rules, created_at, updated_at
		FROM collections WHERE slug = $1 AND tenant_id = $2`,
		slug, string(tenantID),
	).Scan(&cID, &cTenantID, &c.Name, &c.Slug, &c.Description,
		&productIDsJSON, &c.IsAutomatic, &rulesJSON, &c.CreatedAt, &c.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, catalog.ErrCollectionNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "getting collection by slug", errx.TypeInternal)
	}

	c.ID = kernel.CollectionID(cID)
	c.TenantID = kernel.TenantID(cTenantID)
	_ = json.Unmarshal([]byte(productIDsJSON), &c.ProductIDs)
	if c.ProductIDs == nil {
		c.ProductIDs = []kernel.ProductID{}
	}
	_ = json.Unmarshal([]byte(rulesJSON), &c.Rules)
	if c.Rules == nil {
		c.Rules = map[string]any{}
	}
	return &c, nil
}

func (r *CollectionPostgresRepo) Update(ctx context.Context, c *catalog.Collection) error {
	productIDsJSON, err := json.Marshal(c.ProductIDs)
	if err != nil {
		return errx.Wrap(err, "marshaling product IDs", errx.TypeInternal)
	}
	rulesJSON, err := json.Marshal(c.Rules)
	if err != nil {
		return errx.Wrap(err, "marshaling rules", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE collections SET name=$1, slug=$2, description=$3, product_ids=$4, is_automatic=$5, rules=$6, updated_at=$7
		WHERE id=$8 AND tenant_id=$9`,
		c.Name, c.Slug, c.Description, string(productIDsJSON),
		c.IsAutomatic, string(rulesJSON), c.UpdatedAt,
		string(c.ID), string(c.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating collection", errx.TypeInternal)
	}
	return nil
}

func (r *CollectionPostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CollectionID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM collections WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting collection", errx.TypeInternal)
	}
	return nil
}

func (r *CollectionPostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[catalog.Collection], error) {
	var zero kernel.Paginated[catalog.Collection]

	var total int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM collections WHERE tenant_id = $1", string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting collections", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, slug, description, product_ids, is_automatic, rules, created_at, updated_at
		FROM collections WHERE tenant_id = $1
		ORDER BY name ASC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying collections", errx.TypeInternal)
	}
	defer rows.Close()

	var collections []catalog.Collection
	for rows.Next() {
		c, err := scanCollectionFields(rows)
		if err != nil {
			return zero, err
		}
		collections = append(collections, *c)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating collections", errx.TypeInternal)
	}

	return kernel.NewPaginated(collections, pg.Page, pg.PageSize, total), nil
}

func scanCollectionFields(s rowScanner) (*catalog.Collection, error) {
	var c catalog.Collection
	var id, tenantID string
	var productIDsJSON, rulesJSON string

	err := s.Scan(&id, &tenantID, &c.Name, &c.Slug, &c.Description,
		&productIDsJSON, &c.IsAutomatic, &rulesJSON,
		&c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, errx.Wrap(err, "scanning collection", errx.TypeInternal)
	}

	c.ID = kernel.CollectionID(id)
	c.TenantID = kernel.TenantID(tenantID)

	_ = json.Unmarshal([]byte(productIDsJSON), &c.ProductIDs)
	if c.ProductIDs == nil {
		c.ProductIDs = []kernel.ProductID{}
	}

	_ = json.Unmarshal([]byte(rulesJSON), &c.Rules)
	if c.Rules == nil {
		c.Rules = map[string]any{}
	}

	return &c, nil
}

// Ensure interface compliance.
var _ catalog.CollectionRepository = (*CollectionPostgresRepo)(nil)
