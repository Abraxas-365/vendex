package socialauthinfra

import (
	"context"
	"database/sql"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/socialauth"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements socialauth.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed social auth repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Create persists a new social account record.
func (r *PostgresRepo) Create(ctx context.Context, sa socialauth.SocialAccount) (socialauth.SocialAccount, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO social_accounts
			(id, tenant_id, customer_id, provider, provider_user_id, email, name, avatar_url, access_token, refresh_token, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		string(sa.ID),
		string(sa.TenantID),
		string(sa.CustomerID),
		sa.Provider,
		sa.ProviderUserID,
		sa.Email,
		sa.Name,
		sa.AvatarURL,
		sa.AccessToken,
		sa.RefreshToken,
		sa.CreatedAt,
		sa.UpdatedAt,
	)
	if err != nil {
		return socialauth.SocialAccount{}, errx.Wrap(err, "inserting social account", errx.TypeInternal)
	}
	return sa, nil
}

// GetByID retrieves a social account by its ID, scoped to the tenant.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.SocialAccountID) (socialauth.SocialAccount, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, customer_id, provider, provider_user_id,
		       email, name, avatar_url, access_token, refresh_token, created_at, updated_at
		FROM social_accounts
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	sa, err := scanSocialAccount(row)
	if err == sql.ErrNoRows {
		return socialauth.SocialAccount{}, socialauth.ErrNotFound()
	}
	if err != nil {
		return socialauth.SocialAccount{}, errx.Wrap(err, "scanning social account by id", errx.TypeInternal)
	}
	return sa, nil
}

// GetByProvider retrieves a social account by provider + provider_user_id, scoped to the tenant.
func (r *PostgresRepo) GetByProvider(ctx context.Context, tenantID kernel.TenantID, provider, providerUserID string) (socialauth.SocialAccount, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, customer_id, provider, provider_user_id,
		       email, name, avatar_url, access_token, refresh_token, created_at, updated_at
		FROM social_accounts
		WHERE tenant_id = $1 AND provider = $2 AND provider_user_id = $3`,
		string(tenantID), provider, providerUserID,
	)
	sa, err := scanSocialAccount(row)
	if err == sql.ErrNoRows {
		return socialauth.SocialAccount{}, socialauth.ErrNotFound()
	}
	if err != nil {
		return socialauth.SocialAccount{}, errx.Wrap(err, "scanning social account by provider", errx.TypeInternal)
	}
	return sa, nil
}

// ListByCustomer returns all social accounts linked to a customer, scoped to the tenant.
func (r *PostgresRepo) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) ([]socialauth.SocialAccount, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, customer_id, provider, provider_user_id,
		       email, name, avatar_url, access_token, refresh_token, created_at, updated_at
		FROM social_accounts
		WHERE tenant_id = $1 AND customer_id = $2
		ORDER BY created_at ASC`,
		string(tenantID), string(customerID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying social accounts by customer", errx.TypeInternal)
	}
	defer rows.Close()

	var accounts []socialauth.SocialAccount
	for rows.Next() {
		sa, err := scanSocialAccountRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning social account row", errx.TypeInternal)
		}
		accounts = append(accounts, sa)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating social accounts", errx.TypeInternal)
	}
	if accounts == nil {
		accounts = []socialauth.SocialAccount{}
	}
	return accounts, nil
}

// List returns a paginated list of all social accounts for the tenant.
func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[socialauth.SocialAccount], error) {
	var zero kernel.Paginated[socialauth.SocialAccount]

	var total int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM social_accounts WHERE tenant_id = $1",
		string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting social accounts", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, customer_id, provider, provider_user_id,
		       email, name, avatar_url, access_token, refresh_token, created_at, updated_at
		FROM social_accounts
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying social accounts", errx.TypeInternal)
	}
	defer rows.Close()

	var accounts []socialauth.SocialAccount
	for rows.Next() {
		sa, err := scanSocialAccountRow(rows)
		if err != nil {
			return zero, errx.Wrap(err, "scanning social account row", errx.TypeInternal)
		}
		accounts = append(accounts, sa)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating social accounts", errx.TypeInternal)
	}
	if accounts == nil {
		accounts = []socialauth.SocialAccount{}
	}

	return kernel.NewPaginated(accounts, pg.Page, pg.PageSize, total), nil
}

// Delete removes a social account by ID, scoped to the tenant.
func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.SocialAccountID) error {
	result, err := r.db.ExecContext(ctx,
		"DELETE FROM social_accounts WHERE id = $1 AND tenant_id = $2",
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting social account", errx.TypeInternal)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return errx.Wrap(err, "checking rows affected", errx.TypeInternal)
	}
	if n == 0 {
		return socialauth.ErrNotFound()
	}
	return nil
}

// -------------------------------------------------------------------------
// Scan helpers
// -------------------------------------------------------------------------

type scanner interface {
	Scan(dest ...any) error
}

func scanSocialAccount(s *sql.Row) (socialauth.SocialAccount, error) {
	return scanFields(s)
}

func scanSocialAccountRow(rows *sql.Rows) (socialauth.SocialAccount, error) {
	return scanFields(rows)
}

func scanFields(s scanner) (socialauth.SocialAccount, error) {
	var sa socialauth.SocialAccount
	var id, tenantID, customerID string

	err := s.Scan(
		&id, &tenantID, &customerID,
		&sa.Provider, &sa.ProviderUserID,
		&sa.Email, &sa.Name, &sa.AvatarURL,
		&sa.AccessToken, &sa.RefreshToken,
		&sa.CreatedAt, &sa.UpdatedAt,
	)
	if err != nil {
		return socialauth.SocialAccount{}, err
	}

	sa.ID = kernel.SocialAccountID(id)
	sa.TenantID = kernel.TenantID(tenantID)
	sa.CustomerID = kernel.CustomerID(customerID)

	return sa, nil
}

// Compile-time interface check.
var _ socialauth.Repository = (*PostgresRepo)(nil)
