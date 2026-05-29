package promoinfra

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/promo"
)

// PostgresPromoRepository implements promo.PromoRepository using PostgreSQL.
type PostgresPromoRepository struct {
	db *sql.DB
}

// NewPostgresPromoRepository creates a new PostgresPromoRepository.
func NewPostgresPromoRepository(db *sql.DB) *PostgresPromoRepository {
	return &PostgresPromoRepository{db: db}
}

// Create inserts a new promo row.
func (r *PostgresPromoRepository) Create(ctx context.Context, p *promo.Promo) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO promos
			(id, tenant_id, code, type, value, min_order_amount, max_uses, used_count, starts_at, ends_at, active, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		string(p.ID), string(p.TenantID), p.Code, string(p.Type),
		p.Value, p.MinOrderAmount, p.MaxUses, p.UsedCount,
		p.StartsAt, p.EndsAt, p.Active, p.CreatedAt,
	)
	if err != nil {
		if isDuplicateCode(err) {
			return promo.ErrCodeAlreadyExists
		}
		return fmt.Errorf("insert promo: %w", err)
	}
	return nil
}

// GetByID retrieves a promo by primary key, scoped to the tenant.
func (r *PostgresPromoRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.PromoID) (*promo.Promo, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, code, type, value, min_order_amount, max_uses, used_count, starts_at, ends_at, active, created_at
		FROM promos WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	return scanPromo(row)
}

// GetByCode retrieves a promo by its code (case-insensitive), scoped to the tenant.
func (r *PostgresPromoRepository) GetByCode(ctx context.Context, tenantID kernel.TenantID, code string) (*promo.Promo, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, code, type, value, min_order_amount, max_uses, used_count, starts_at, ends_at, active, created_at
		FROM promos WHERE tenant_id=$1 AND UPPER(code)=UPPER($2)`,
		string(tenantID), code,
	)
	return scanPromo(row)
}

// Update persists mutations to an existing promo.
func (r *PostgresPromoRepository) Update(ctx context.Context, p *promo.Promo) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE promos
		SET code=$3, type=$4, value=$5, min_order_amount=$6, max_uses=$7,
		    used_count=$8, starts_at=$9, ends_at=$10, active=$11
		WHERE tenant_id=$1 AND id=$2`,
		string(p.TenantID), string(p.ID),
		p.Code, string(p.Type), p.Value,
		p.MinOrderAmount, p.MaxUses, p.UsedCount,
		p.StartsAt, p.EndsAt, p.Active,
	)
	if err != nil {
		return fmt.Errorf("update promo: %w", err)
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
		return fmt.Errorf("increment used_count: %w", err)
	}
	return nil
}

// List returns promos for a tenant with pagination.
func (r *PostgresPromoRepository) List(ctx context.Context, tenantID kernel.TenantID, p kernel.Pagination) (kernel.PaginatedResult[promo.Promo], error) {
	var total int
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM promos WHERE tenant_id=$1`, string(tenantID)).Scan(&total)
	if err != nil {
		return kernel.PaginatedResult[promo.Promo]{}, fmt.Errorf("count promos: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, code, type, value, min_order_amount, max_uses, used_count, starts_at, ends_at, active, created_at
		FROM promos WHERE tenant_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.PaginatedResult[promo.Promo]{}, fmt.Errorf("list promos: %w", err)
	}
	defer rows.Close()

	var promos []promo.Promo
	for rows.Next() {
		pr, err := scanPromoRow(rows)
		if err != nil {
			return kernel.PaginatedResult[promo.Promo]{}, err
		}
		promos = append(promos, *pr)
	}
	if err := rows.Err(); err != nil {
		return kernel.PaginatedResult[promo.Promo]{}, fmt.Errorf("iterate promos: %w", err)
	}
	return kernel.NewPaginatedResult(promos, total, p), nil
}

// scanPromo scans a single promo row from sql.Row.
func scanPromo(row *sql.Row) (*promo.Promo, error) {
	var p promo.Promo
	var idStr, tenantStr, typeStr string
	var startsAt, endsAt sql.NullTime

	err := row.Scan(
		&idStr, &tenantStr, &p.Code, &typeStr,
		&p.Value, &p.MinOrderAmount, &p.MaxUses, &p.UsedCount,
		&startsAt, &endsAt, &p.Active, &p.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, promo.ErrPromoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan promo: %w", err)
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
	return &p, nil
}

// scanPromoRow scans a single promo from sql.Rows.
func scanPromoRow(rows *sql.Rows) (*promo.Promo, error) {
	var p promo.Promo
	var idStr, tenantStr, typeStr string
	var startsAt, endsAt sql.NullTime

	err := rows.Scan(
		&idStr, &tenantStr, &p.Code, &typeStr,
		&p.Value, &p.MinOrderAmount, &p.MaxUses, &p.UsedCount,
		&startsAt, &endsAt, &p.Active, &p.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("scan promo row: %w", err)
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
	return &p, nil
}

// isDuplicateCode detects unique constraint violations on the code column.
func isDuplicateCode(err error) bool {
	return strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate")
}

// Ensure time import is used.
var _ = time.Time{}
