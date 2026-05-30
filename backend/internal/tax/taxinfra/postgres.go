package taxinfra

import (
	"context"
	"database/sql"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/tax"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements tax.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed tax rate repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, rate *tax.TaxRate) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO tax_rates (id, tenant_id, name, rate, country, state, city, zip_code, priority, compound, includes_shipping, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
		string(rate.ID), string(rate.TenantID), rate.Name, rate.Rate,
		rate.Country, rate.State, rate.City, rate.ZipCode,
		rate.Priority, rate.Compound, rate.IncludesShipping, rate.Active,
		rate.CreatedAt, rate.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting tax rate", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.TaxRateID) (*tax.TaxRate, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, rate, country, state, city, zip_code, priority, compound, includes_shipping, active, created_at, updated_at
		FROM tax_rates WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	rate, err := scanTaxRate(row)
	if err == sql.ErrNoRows {
		return nil, tax.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning tax rate", errx.TypeInternal)
	}
	return rate, nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID) ([]tax.TaxRate, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, rate, country, state, city, zip_code, priority, compound, includes_shipping, active, created_at, updated_at
		FROM tax_rates WHERE tenant_id = $1
		ORDER BY country ASC, state ASC, city ASC, priority ASC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying tax rates", errx.TypeInternal)
	}
	defer rows.Close()

	var rates []tax.TaxRate
	for rows.Next() {
		rate, err := scanTaxRateRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning tax rate row", errx.TypeInternal)
		}
		rates = append(rates, *rate)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating tax rates", errx.TypeInternal)
	}

	if rates == nil {
		rates = []tax.TaxRate{}
	}
	return rates, nil
}

func (r *PostgresRepo) Update(ctx context.Context, rate *tax.TaxRate) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE tax_rates
		SET name=$1, rate=$2, country=$3, state=$4, city=$5, zip_code=$6,
		    priority=$7, compound=$8, includes_shipping=$9, active=$10, updated_at=$11
		WHERE id=$12 AND tenant_id=$13`,
		rate.Name, rate.Rate, rate.Country, rate.State, rate.City, rate.ZipCode,
		rate.Priority, rate.Compound, rate.IncludesShipping, rate.Active, rate.UpdatedAt,
		string(rate.ID), string(rate.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating tax rate", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.TaxRateID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tax_rates WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting tax rate", errx.TypeInternal)
	}
	return nil
}

// FindByLocation finds active tax rates matching the given location.
// Matches at multiple jurisdiction levels (country, state, city, zip).
func (r *PostgresRepo) FindByLocation(ctx context.Context, tenantID kernel.TenantID, country, state, city, zipCode string) ([]tax.TaxRate, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, rate, country, state, city, zip_code, priority, compound, includes_shipping, active, created_at, updated_at
		FROM tax_rates
		WHERE tenant_id = $1
		  AND country = $2
		  AND active = true
		  AND (state = '' OR state IS NULL OR state = $3)
		  AND (city = '' OR city IS NULL OR city = $4)
		  AND (zip_code = '' OR zip_code IS NULL OR zip_code = $5)
		ORDER BY priority ASC`,
		string(tenantID), country, state, city, zipCode,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying tax rates by location", errx.TypeInternal)
	}
	defer rows.Close()

	var rates []tax.TaxRate
	for rows.Next() {
		rate, err := scanTaxRateRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning tax rate row", errx.TypeInternal)
		}
		rates = append(rates, *rate)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating tax rates", errx.TypeInternal)
	}

	if rates == nil {
		rates = []tax.TaxRate{}
	}
	return rates, nil
}

// scanner is implemented by both *sql.Row and *sql.Rows.
type scanner interface {
	Scan(dest ...any) error
}

func scanTaxRate(row *sql.Row) (*tax.TaxRate, error) {
	return scanFields(row)
}

func scanTaxRateRow(rows interface{ Scan(dest ...any) error }) (*tax.TaxRate, error) {
	return scanFields(rows)
}

func scanFields(s scanner) (*tax.TaxRate, error) {
	var r tax.TaxRate
	var id, tenantID string

	err := s.Scan(
		&id, &tenantID, &r.Name, &r.Rate,
		&r.Country, &r.State, &r.City, &r.ZipCode,
		&r.Priority, &r.Compound, &r.IncludesShipping, &r.Active,
		&r.CreatedAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	r.ID = kernel.TaxRateID(id)
	r.TenantID = kernel.TenantID(tenantID)

	return &r, nil
}

// Ensure interface compliance.
var _ tax.Repository = (*PostgresRepo)(nil)
