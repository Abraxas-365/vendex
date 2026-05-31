package shippinginfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/shipping"
	"github.com/jmoiron/sqlx"
)

// ---------------------------------------------------------------------------
// Zone repository
// ---------------------------------------------------------------------------

// ZonePostgresRepo implements shipping.ZoneRepository using sqlx.
type ZonePostgresRepo struct {
	db *sqlx.DB
}

// NewZonePostgresRepo creates a new PostgreSQL-backed zone repository.
func NewZonePostgresRepo(db *sqlx.DB) *ZonePostgresRepo {
	return &ZonePostgresRepo{db: db}
}

func (r *ZonePostgresRepo) Create(ctx context.Context, zone *shipping.ShippingZone) error {
	countriesJSON, err := json.Marshal(zone.Countries)
	if err != nil {
		return errx.Wrap(err, "marshaling countries", errx.TypeInternal)
	}
	statesJSON, err := json.Marshal(zone.States)
	if err != nil {
		return errx.Wrap(err, "marshaling states", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO shipping_zones (id, tenant_id, name, countries, states, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		string(zone.ID), string(zone.TenantID), zone.Name,
		string(countriesJSON), string(statesJSON),
		zone.CreatedAt, zone.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting shipping zone", errx.TypeInternal)
	}
	return nil
}

func (r *ZonePostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingZoneID) (*shipping.ShippingZone, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, countries, states, created_at, updated_at
		FROM shipping_zones WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	return r.scanZone(row)
}

func (r *ZonePostgresRepo) List(ctx context.Context, tenantID kernel.TenantID) ([]shipping.ShippingZone, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, countries, states, created_at, updated_at
		FROM shipping_zones WHERE tenant_id = $1 ORDER BY created_at DESC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying shipping zones", errx.TypeInternal)
	}
	defer rows.Close()

	var zones []shipping.ShippingZone
	for rows.Next() {
		z, err := r.scanZoneRow(rows)
		if err != nil {
			return nil, err
		}
		zones = append(zones, *z)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating shipping zones", errx.TypeInternal)
	}
	if zones == nil {
		zones = []shipping.ShippingZone{}
	}
	return zones, nil
}

func (r *ZonePostgresRepo) Update(ctx context.Context, zone *shipping.ShippingZone) error {
	countriesJSON, err := json.Marshal(zone.Countries)
	if err != nil {
		return errx.Wrap(err, "marshaling countries", errx.TypeInternal)
	}
	statesJSON, err := json.Marshal(zone.States)
	if err != nil {
		return errx.Wrap(err, "marshaling states", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE shipping_zones SET name=$1, countries=$2, states=$3, updated_at=$4
		WHERE id=$5 AND tenant_id=$6`,
		zone.Name, string(countriesJSON), string(statesJSON), zone.UpdatedAt,
		string(zone.ID), string(zone.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating shipping zone", errx.TypeInternal)
	}
	return nil
}

func (r *ZonePostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingZoneID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM shipping_zones WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting shipping zone", errx.TypeInternal)
	}
	return nil
}

// FindByAddress finds zones that cover the given country (and optionally state).
// It matches zones where the countries JSONB array contains the given country
// OR where the countries array is empty (covers everywhere).
// For state matching, it additionally filters on states that contain the given
// state OR where states is empty (covers all states in matched country).
func (r *ZonePostgresRepo) FindByAddress(ctx context.Context, tenantID kernel.TenantID, country, state string) ([]shipping.ShippingZone, error) {
	// Build JSONB arguments
	countryJSON := fmt.Sprintf(`[%q]`, country)

	var rows *sql.Rows
	var err error

	if state != "" {
		stateJSON := fmt.Sprintf(`[%q]`, state)
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, tenant_id, name, countries, states, created_at, updated_at
			FROM shipping_zones
			WHERE tenant_id = $1
			  AND (countries = '[]'::jsonb OR countries @> $2::jsonb)
			  AND (states = '[]'::jsonb OR states @> $3::jsonb)
			ORDER BY created_at DESC`,
			string(tenantID), countryJSON, stateJSON,
		)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, tenant_id, name, countries, states, created_at, updated_at
			FROM shipping_zones
			WHERE tenant_id = $1
			  AND (countries = '[]'::jsonb OR countries @> $2::jsonb)
			  AND states = '[]'::jsonb
			ORDER BY created_at DESC`,
			string(tenantID), countryJSON,
		)
	}

	if err != nil {
		return nil, errx.Wrap(err, "querying zones by address", errx.TypeInternal)
	}
	defer rows.Close()

	var zones []shipping.ShippingZone
	for rows.Next() {
		z, err := r.scanZoneRow(rows)
		if err != nil {
			return nil, err
		}
		zones = append(zones, *z)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating zones by address", errx.TypeInternal)
	}
	if zones == nil {
		zones = []shipping.ShippingZone{}
	}
	return zones, nil
}

// scanZone scans a single zone from a QueryRow result.
func (r *ZonePostgresRepo) scanZone(row *sql.Row) (*shipping.ShippingZone, error) {
	var z shipping.ShippingZone
	var id, tenantID, countriesJSON, statesJSON string

	err := row.Scan(&id, &tenantID, &z.Name, &countriesJSON, &statesJSON, &z.CreatedAt, &z.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, shipping.ErrZoneNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning shipping zone", errx.TypeInternal)
	}

	z.ID = kernel.ShippingZoneID(id)
	z.TenantID = kernel.TenantID(tenantID)
	_ = json.Unmarshal([]byte(countriesJSON), &z.Countries)
	_ = json.Unmarshal([]byte(statesJSON), &z.States)
	if z.Countries == nil {
		z.Countries = []string{}
	}
	if z.States == nil {
		z.States = []string{}
	}
	return &z, nil
}

// scanZoneRow scans a single zone from a Rows cursor.
func (r *ZonePostgresRepo) scanZoneRow(rows *sql.Rows) (*shipping.ShippingZone, error) {
	var z shipping.ShippingZone
	var id, tenantID, countriesJSON, statesJSON string

	err := rows.Scan(&id, &tenantID, &z.Name, &countriesJSON, &statesJSON, &z.CreatedAt, &z.UpdatedAt)
	if err != nil {
		return nil, errx.Wrap(err, "scanning shipping zone row", errx.TypeInternal)
	}

	z.ID = kernel.ShippingZoneID(id)
	z.TenantID = kernel.TenantID(tenantID)
	_ = json.Unmarshal([]byte(countriesJSON), &z.Countries)
	_ = json.Unmarshal([]byte(statesJSON), &z.States)
	if z.Countries == nil {
		z.Countries = []string{}
	}
	if z.States == nil {
		z.States = []string{}
	}
	return &z, nil
}

// Ensure interface compliance.
var _ shipping.ZoneRepository = (*ZonePostgresRepo)(nil)

// ---------------------------------------------------------------------------
// Rate repository
// ---------------------------------------------------------------------------

// RatePostgresRepo implements shipping.RateRepository using sqlx.
type RatePostgresRepo struct {
	db *sqlx.DB
}

// NewRatePostgresRepo creates a new PostgreSQL-backed rate repository.
func NewRatePostgresRepo(db *sqlx.DB) *RatePostgresRepo {
	return &RatePostgresRepo{db: db}
}

func (r *RatePostgresRepo) Create(ctx context.Context, rate *shipping.ShippingRate) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO shipping_rates
		  (id, zone_id, tenant_id, name, type, price_amount, price_currency,
		   min_weight, max_weight, min_order_amount, max_order_amount,
		   est_days_min, est_days_max, active, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16)`,
		string(rate.ID), string(rate.ZoneID), string(rate.TenantID),
		rate.Name, string(rate.Type),
		rate.Price.Amount, rate.Price.Currency,
		rate.MinWeight, rate.MaxWeight,
		rate.MinOrderAmount, rate.MaxOrderAmount,
		rate.EstDaysMin, rate.EstDaysMax,
		rate.Active, rate.CreatedAt, rate.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting shipping rate", errx.TypeInternal)
	}
	return nil
}

func (r *RatePostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingRateID) (*shipping.ShippingRate, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, zone_id, tenant_id, name, type, price_amount, price_currency,
		       min_weight, max_weight, min_order_amount, max_order_amount,
		       est_days_min, est_days_max, active, created_at, updated_at
		FROM shipping_rates WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	return r.scanRate(row)
}

func (r *RatePostgresRepo) ListByZone(ctx context.Context, tenantID kernel.TenantID, zoneID kernel.ShippingZoneID) ([]shipping.ShippingRate, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, zone_id, tenant_id, name, type, price_amount, price_currency,
		       min_weight, max_weight, min_order_amount, max_order_amount,
		       est_days_min, est_days_max, active, created_at, updated_at
		FROM shipping_rates WHERE zone_id = $1 AND tenant_id = $2 ORDER BY created_at DESC`,
		string(zoneID), string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying shipping rates", errx.TypeInternal)
	}
	defer rows.Close()

	var rates []shipping.ShippingRate
	for rows.Next() {
		rate, err := r.scanRateRow(rows)
		if err != nil {
			return nil, err
		}
		rates = append(rates, *rate)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating shipping rates", errx.TypeInternal)
	}
	if rates == nil {
		rates = []shipping.ShippingRate{}
	}
	return rates, nil
}

func (r *RatePostgresRepo) Update(ctx context.Context, rate *shipping.ShippingRate) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE shipping_rates SET
		  name=$1, type=$2, price_amount=$3, price_currency=$4,
		  min_weight=$5, max_weight=$6, min_order_amount=$7, max_order_amount=$8,
		  est_days_min=$9, est_days_max=$10, active=$11, updated_at=$12
		WHERE id=$13 AND tenant_id=$14`,
		rate.Name, string(rate.Type),
		rate.Price.Amount, rate.Price.Currency,
		rate.MinWeight, rate.MaxWeight,
		rate.MinOrderAmount, rate.MaxOrderAmount,
		rate.EstDaysMin, rate.EstDaysMax,
		rate.Active, rate.UpdatedAt,
		string(rate.ID), string(rate.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating shipping rate", errx.TypeInternal)
	}
	return nil
}

func (r *RatePostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.ShippingRateID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM shipping_rates WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting shipping rate", errx.TypeInternal)
	}
	return nil
}

// scanRate scans a single rate from a QueryRow result.
func (r *RatePostgresRepo) scanRate(row *sql.Row) (*shipping.ShippingRate, error) {
	var rate shipping.ShippingRate
	var id, zoneID, tenantID, rateType string

	err := row.Scan(
		&id, &zoneID, &tenantID,
		&rate.Name, &rateType,
		&rate.Price.Amount, &rate.Price.Currency,
		&rate.MinWeight, &rate.MaxWeight,
		&rate.MinOrderAmount, &rate.MaxOrderAmount,
		&rate.EstDaysMin, &rate.EstDaysMax,
		&rate.Active, &rate.CreatedAt, &rate.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, shipping.ErrRateNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning shipping rate", errx.TypeInternal)
	}

	rate.ID = kernel.ShippingRateID(id)
	rate.ZoneID = kernel.ShippingZoneID(zoneID)
	rate.TenantID = kernel.TenantID(tenantID)
	rate.Type = shipping.RateType(rateType)
	return &rate, nil
}

// scanRateRow scans a single rate from a Rows cursor.
func (r *RatePostgresRepo) scanRateRow(rows *sql.Rows) (*shipping.ShippingRate, error) {
	var rate shipping.ShippingRate
	var id, zoneID, tenantID, rateType string

	err := rows.Scan(
		&id, &zoneID, &tenantID,
		&rate.Name, &rateType,
		&rate.Price.Amount, &rate.Price.Currency,
		&rate.MinWeight, &rate.MaxWeight,
		&rate.MinOrderAmount, &rate.MaxOrderAmount,
		&rate.EstDaysMin, &rate.EstDaysMax,
		&rate.Active, &rate.CreatedAt, &rate.UpdatedAt,
	)
	if err != nil {
		return nil, errx.Wrap(err, "scanning shipping rate row", errx.TypeInternal)
	}

	rate.ID = kernel.ShippingRateID(id)
	rate.ZoneID = kernel.ShippingZoneID(zoneID)
	rate.TenantID = kernel.TenantID(tenantID)
	rate.Type = shipping.RateType(rateType)
	return &rate, nil
}

// Ensure interface compliance.
var _ shipping.RateRepository = (*RatePostgresRepo)(nil)
