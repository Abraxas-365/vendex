package currencyinfra

import (
	"context"
	"database/sql"

	"github.com/Abraxas-365/hada-commerce/internal/currency"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements currency.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed currency rate repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

const selectFields = `id, tenant_id, base_currency, target_currency, rate, auto_update, updated_at, created_at`

func (r *PostgresRepo) Create(ctx context.Context, rate *currency.CurrencyRate) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO currency_rates (id, tenant_id, base_currency, target_currency, rate, auto_update, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(rate.ID), string(rate.TenantID),
		rate.BaseCurrency, rate.TargetCurrency,
		rate.Rate, rate.AutoUpdate,
		rate.CreatedAt, rate.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting currency rate", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetByPair(ctx context.Context, tenantID kernel.TenantID, base, target string) (*currency.CurrencyRate, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+selectFields+`
		FROM currency_rates
		WHERE tenant_id = $1 AND base_currency = $2 AND target_currency = $3`,
		string(tenantID), base, target,
	)
	rate, err := scanRate(row)
	if err == sql.ErrNoRows {
		return nil, currency.ErrRateNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning currency rate", errx.TypeInternal)
	}
	return rate, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CurrencyRateID) (*currency.CurrencyRate, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT `+selectFields+`
		FROM currency_rates
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	rate, err := scanRate(row)
	if err == sql.ErrNoRows {
		return nil, currency.ErrRateNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning currency rate by id", errx.TypeInternal)
	}
	return rate, nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID) ([]currency.CurrencyRate, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT `+selectFields+`
		FROM currency_rates
		WHERE tenant_id = $1
		ORDER BY base_currency ASC, target_currency ASC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying currency rates", errx.TypeInternal)
	}
	defer rows.Close()

	var rates []currency.CurrencyRate
	for rows.Next() {
		rate, err := scanRateRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning currency rate row", errx.TypeInternal)
		}
		rates = append(rates, *rate)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating currency rates", errx.TypeInternal)
	}

	if rates == nil {
		rates = []currency.CurrencyRate{}
	}
	return rates, nil
}

func (r *PostgresRepo) Update(ctx context.Context, rate *currency.CurrencyRate) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE currency_rates
		SET rate = $1, auto_update = $2, updated_at = $3
		WHERE id = $4 AND tenant_id = $5`,
		rate.Rate, rate.AutoUpdate, rate.UpdatedAt,
		string(rate.ID), string(rate.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating currency rate", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CurrencyRateID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM currency_rates WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting currency rate", errx.TypeInternal)
	}
	return nil
}

// scanner is implemented by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanRate(row *sql.Row) (*currency.CurrencyRate, error) {
	return scanFields(row)
}

func scanRateRow(rows interface{ Scan(dest ...any) error }) (*currency.CurrencyRate, error) {
	return scanFields(rows)
}

func scanFields(s scanner) (*currency.CurrencyRate, error) {
	var cr currency.CurrencyRate
	var id, tenantID string

	err := s.Scan(
		&id, &tenantID,
		&cr.BaseCurrency, &cr.TargetCurrency,
		&cr.Rate, &cr.AutoUpdate,
		&cr.UpdatedAt, &cr.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	cr.ID = kernel.CurrencyRateID(id)
	cr.TenantID = kernel.TenantID(tenantID)

	return &cr, nil
}

// Ensure interface compliance.
var _ currency.Repository = (*PostgresRepo)(nil)
