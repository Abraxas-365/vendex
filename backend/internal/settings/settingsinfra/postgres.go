package settingsinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/settings"
)

// PostgresRepo implements settings.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed settings repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Get retrieves the store settings for the given tenant.
// Returns settings.ErrNotFound if no row exists.
func (r *PostgresRepo) Get(ctx context.Context, tenantID kernel.TenantID) (*settings.StoreSettings, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT tenant_id, store_name, store_email, store_phone, currency, timezone,
		       address, logo_url, favicon_url, social_links, checkout_config, updated_at
		FROM store_settings
		WHERE tenant_id = $1`,
		string(tenantID),
	)

	var s settings.StoreSettings
	var tenantStr string
	var addressJSON, socialLinksJSON, checkoutConfigJSON string
	var updatedAt time.Time

	err := row.Scan(
		&tenantStr,
		&s.StoreName,
		&s.StoreEmail,
		&s.StorePhone,
		&s.Currency,
		&s.Timezone,
		&addressJSON,
		&s.LogoURL,
		&s.FaviconURL,
		&socialLinksJSON,
		&checkoutConfigJSON,
		&updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, settings.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning settings row", errx.TypeInternal)
	}

	s.TenantID = kernel.TenantID(tenantStr)
	s.UpdatedAt = updatedAt

	if err := json.Unmarshal([]byte(addressJSON), &s.Address); err != nil {
		return nil, errx.Wrap(err, "unmarshaling address", errx.TypeInternal)
	}
	if err := json.Unmarshal([]byte(socialLinksJSON), &s.SocialLinks); err != nil {
		return nil, errx.Wrap(err, "unmarshaling social_links", errx.TypeInternal)
	}
	if err := json.Unmarshal([]byte(checkoutConfigJSON), &s.CheckoutConfig); err != nil {
		return nil, errx.Wrap(err, "unmarshaling checkout_config", errx.TypeInternal)
	}

	return &s, nil
}

// Upsert creates or updates the settings row for the given tenant.
func (r *PostgresRepo) Upsert(ctx context.Context, s *settings.StoreSettings) error {
	addressJSON, err := json.Marshal(s.Address)
	if err != nil {
		return fmt.Errorf("marshaling address: %w", err)
	}
	socialLinksJSON, err := json.Marshal(s.SocialLinks)
	if err != nil {
		return fmt.Errorf("marshaling social_links: %w", err)
	}
	checkoutConfigJSON, err := json.Marshal(s.CheckoutConfig)
	if err != nil {
		return fmt.Errorf("marshaling checkout_config: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO store_settings
			(tenant_id, store_name, store_email, store_phone, currency, timezone,
			 address, logo_url, favicon_url, social_links, checkout_config, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		ON CONFLICT (tenant_id) DO UPDATE SET
			store_name      = EXCLUDED.store_name,
			store_email     = EXCLUDED.store_email,
			store_phone     = EXCLUDED.store_phone,
			currency        = EXCLUDED.currency,
			timezone        = EXCLUDED.timezone,
			address         = EXCLUDED.address,
			logo_url        = EXCLUDED.logo_url,
			favicon_url     = EXCLUDED.favicon_url,
			social_links    = EXCLUDED.social_links,
			checkout_config = EXCLUDED.checkout_config,
			updated_at      = EXCLUDED.updated_at`,
		string(s.TenantID),
		s.StoreName,
		s.StoreEmail,
		s.StorePhone,
		s.Currency,
		s.Timezone,
		string(addressJSON),
		s.LogoURL,
		s.FaviconURL,
		string(socialLinksJSON),
		string(checkoutConfigJSON),
		s.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "upserting settings", errx.TypeInternal)
	}
	return nil
}

// Ensure interface compliance at compile time.
var _ settings.Repository = (*PostgresRepo)(nil)
