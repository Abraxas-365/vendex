package promoinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/promo"
)

// PostgresPromoRepository implements promo.PromoRepository using PostgreSQL.
type PostgresPromoRepository struct {
	db *sqlx.DB
}

// NewPostgresPromoRepository creates a new PostgresPromoRepository.
func NewPostgresPromoRepository(db *sqlx.DB) *PostgresPromoRepository {
	return &PostgresPromoRepository{db: db}
}

// Create inserts a new promo row.
func (r *PostgresPromoRepository) Create(ctx context.Context, p *promo.Promo) error {
	productIDs, err := jsonMarshalStringSlice(p.TargetProductIDs)
	if err != nil {
		return errx.Wrap(err, "marshal target_product_ids", errx.TypeInternal)
	}
	categoryIDs, err := jsonMarshalStringSlice(p.TargetCategoryIDs)
	if err != nil {
		return errx.Wrap(err, "marshal target_category_ids", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO promos
			(id, tenant_id, code, type, value, min_order_amount, max_uses, used_count,
			 starts_at, ends_at, active, created_at,
			 target_product_ids, target_category_ids, customer_group_id, stackable,
			 buy_quantity, get_quantity, get_product_id, get_discount)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16,$17,$18,$19,$20)`,
		string(p.ID), string(p.TenantID), p.Code, string(p.Type),
		p.Value, p.MinOrderAmount, p.MaxUses, p.UsedCount,
		p.StartsAt, p.EndsAt, p.Active, p.CreatedAt,
		productIDs, categoryIDs, nullableString(p.CustomerGroupID), p.Stackable,
		p.BuyQuantity, p.GetQuantity, nullableString(p.GetProductID), p.GetDiscount,
	)
	if err != nil {
		if isDuplicateCode(err) {
			return promo.ErrCodeAlreadyExists
		}
		return errx.Wrap(err, "insert promo", errx.TypeInternal)
	}
	return nil
}

// GetByID retrieves a promo by primary key, scoped to the tenant.
func (r *PostgresPromoRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PromoID) (*promo.Promo, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, code, type, value, min_order_amount, max_uses, used_count,
		       starts_at, ends_at, active, created_at,
		       target_product_ids, target_category_ids, customer_group_id, stackable,
		       buy_quantity, get_quantity, get_product_id, get_discount
		FROM promos WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	return scanPromo(row)
}

// GetByCode retrieves a promo by its code (case-insensitive), scoped to the tenant.
func (r *PostgresPromoRepository) GetByCode(ctx context.Context, tenantID kernel.TenantID, code string) (*promo.Promo, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, code, type, value, min_order_amount, max_uses, used_count,
		       starts_at, ends_at, active, created_at,
		       target_product_ids, target_category_ids, customer_group_id, stackable,
		       buy_quantity, get_quantity, get_product_id, get_discount
		FROM promos WHERE tenant_id=$1 AND UPPER(code)=UPPER($2)`,
		string(tenantID), code,
	)
	return scanPromo(row)
}

// Update persists mutations to an existing promo.
func (r *PostgresPromoRepository) Update(ctx context.Context, p *promo.Promo) error {
	productIDs, err := jsonMarshalStringSlice(p.TargetProductIDs)
	if err != nil {
		return errx.Wrap(err, "marshal target_product_ids", errx.TypeInternal)
	}
	categoryIDs, err := jsonMarshalStringSlice(p.TargetCategoryIDs)
	if err != nil {
		return errx.Wrap(err, "marshal target_category_ids", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE promos
		SET code=$3, type=$4, value=$5, min_order_amount=$6, max_uses=$7,
		    used_count=$8, starts_at=$9, ends_at=$10, active=$11,
		    target_product_ids=$12, target_category_ids=$13, customer_group_id=$14,
		    stackable=$15, buy_quantity=$16, get_quantity=$17, get_product_id=$18, get_discount=$19
		WHERE tenant_id=$1 AND id=$2`,
		string(p.TenantID), string(p.ID),
		p.Code, string(p.Type), p.Value,
		p.MinOrderAmount, p.MaxUses, p.UsedCount,
		p.StartsAt, p.EndsAt, p.Active,
		productIDs, categoryIDs, nullableString(p.CustomerGroupID),
		p.Stackable, p.BuyQuantity, p.GetQuantity, nullableString(p.GetProductID), p.GetDiscount,
	)
	if err != nil {
		return errx.Wrap(err, "update promo", errx.TypeInternal)
	}
	return nil
}

// IncrementUsedCount atomically increments used_count for a promo.
func (r *PostgresPromoRepository) IncrementUsedCount(ctx context.Context, tenantID kernel.TenantID, id kernel.PromoID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE promos SET used_count = used_count + 1
		WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	if err != nil {
		return errx.Wrap(err, "increment used_count", errx.TypeInternal)
	}
	return nil
}

// List returns promos for a tenant with pagination.
func (r *PostgresPromoRepository) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[promo.Promo], error) {
	var zero kernel.Paginated[promo.Promo]

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM promos WHERE tenant_id=$1`, string(tenantID)).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "count promos", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, code, type, value, min_order_amount, max_uses, used_count,
		       starts_at, ends_at, active, created_at,
		       target_product_ids, target_category_ids, customer_group_id, stackable,
		       buy_quantity, get_quantity, get_product_id, get_discount
		FROM promos WHERE tenant_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "list promos", errx.TypeInternal)
	}
	defer rows.Close()

	var promos []promo.Promo
	for rows.Next() {
		pr, err := scanPromoRow(rows)
		if err != nil {
			return zero, err
		}
		promos = append(promos, *pr)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterate promos", errx.TypeInternal)
	}
	return kernel.NewPaginated(promos, p.Page, p.PageSize, total), nil
}

// scanPromo scans a single promo row from sql.Row.
func scanPromo(row *sql.Row) (*promo.Promo, error) {
	var p promo.Promo
	var idStr, tenantStr, typeStr string
	var startsAt, endsAt sql.NullTime
	var productIDsJSON, categoryIDsJSON []byte
	var customerGroupID, getProductID sql.NullString

	err := row.Scan(
		&idStr, &tenantStr, &p.Code, &typeStr,
		&p.Value, &p.MinOrderAmount, &p.MaxUses, &p.UsedCount,
		&startsAt, &endsAt, &p.Active, &p.CreatedAt,
		&productIDsJSON, &categoryIDsJSON, &customerGroupID, &p.Stackable,
		&p.BuyQuantity, &p.GetQuantity, &getProductID, &p.GetDiscount,
	)
	if err == sql.ErrNoRows {
		return nil, promo.ErrPromoNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scan promo", errx.TypeInternal)
	}

	p.ID = kernel.PromoID(idStr)
	p.TenantID = kernel.TenantID(tenantStr)
	p.Type = promo.PromoType(typeStr)
	if startsAt.Valid {
		t := startsAt.Time
		p.StartsAt = &t
	}
	if endsAt.Valid {
		t := endsAt.Time
		p.EndsAt = &t
	}
	if customerGroupID.Valid {
		p.CustomerGroupID = customerGroupID.String
	}
	if getProductID.Valid {
		p.GetProductID = getProductID.String
	}
	p.TargetProductIDs = jsonUnmarshalStringSlice(productIDsJSON)
	p.TargetCategoryIDs = jsonUnmarshalStringSlice(categoryIDsJSON)
	return &p, nil
}

// scanPromoRow scans a single promo from sql.Rows.
func scanPromoRow(rows *sql.Rows) (*promo.Promo, error) {
	var p promo.Promo
	var idStr, tenantStr, typeStr string
	var startsAt, endsAt sql.NullTime
	var productIDsJSON, categoryIDsJSON []byte
	var customerGroupID, getProductID sql.NullString

	err := rows.Scan(
		&idStr, &tenantStr, &p.Code, &typeStr,
		&p.Value, &p.MinOrderAmount, &p.MaxUses, &p.UsedCount,
		&startsAt, &endsAt, &p.Active, &p.CreatedAt,
		&productIDsJSON, &categoryIDsJSON, &customerGroupID, &p.Stackable,
		&p.BuyQuantity, &p.GetQuantity, &getProductID, &p.GetDiscount,
	)
	if err != nil {
		return nil, errx.Wrap(err, "scan promo row", errx.TypeInternal)
	}

	p.ID = kernel.PromoID(idStr)
	p.TenantID = kernel.TenantID(tenantStr)
	p.Type = promo.PromoType(typeStr)
	if startsAt.Valid {
		t := startsAt.Time
		p.StartsAt = &t
	}
	if endsAt.Valid {
		t := endsAt.Time
		p.EndsAt = &t
	}
	if customerGroupID.Valid {
		p.CustomerGroupID = customerGroupID.String
	}
	if getProductID.Valid {
		p.GetProductID = getProductID.String
	}
	p.TargetProductIDs = jsonUnmarshalStringSlice(productIDsJSON)
	p.TargetCategoryIDs = jsonUnmarshalStringSlice(categoryIDsJSON)
	return &p, nil
}

// isDuplicateCode detects unique constraint violations on the code column.
func isDuplicateCode(err error) bool {
	return strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate")
}

// jsonMarshalStringSlice encodes a string slice as a JSON byte slice.
// Returns '[]' for nil/empty slices.
func jsonMarshalStringSlice(s []string) ([]byte, error) {
	if len(s) == 0 {
		return []byte("[]"), nil
	}
	return json.Marshal(s)
}

// jsonUnmarshalStringSlice decodes a JSON byte slice into a string slice.
// Returns nil for NULL or empty arrays.
func jsonUnmarshalStringSlice(b []byte) []string {
	if len(b) == 0 {
		return nil
	}
	var result []string
	if err := json.Unmarshal(b, &result); err != nil {
		return nil
	}
	if len(result) == 0 {
		return nil
	}
	return result
}

// nullableString converts an empty string to sql.NullString (NULL in DB).
func nullableString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

// Ensure interface compliance.
var _ promo.PromoRepository = (*PostgresPromoRepository)(nil)
