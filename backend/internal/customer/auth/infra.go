package auth

import (
	"context"
	"database/sql"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresCredentialsRepo implements CredentialsRepository using sqlx.
type PostgresCredentialsRepo struct {
	db *sqlx.DB
}

// NewPostgresCredentialsRepo creates a new Postgres-backed credentials repository.
func NewPostgresCredentialsRepo(db *sqlx.DB) *PostgresCredentialsRepo {
	return &PostgresCredentialsRepo{db: db}
}

// Create inserts new customer credentials.
func (r *PostgresCredentialsRepo) Create(ctx context.Context, creds *CustomerCredentials) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO customer_credentials (id, customer_id, tenant_id, email, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		creds.ID,
		string(creds.CustomerID),
		string(creds.TenantID),
		creds.Email,
		creds.PasswordHash,
		creds.CreatedAt,
		creds.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting customer credentials", errx.TypeInternal)
	}
	return nil
}

// GetByEmail retrieves credentials by email scoped to tenant.
func (r *PostgresCredentialsRepo) GetByEmail(ctx context.Context, tenantID kernel.TenantID, email string) (*CustomerCredentials, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, customer_id, tenant_id, email, password_hash, created_at, updated_at
		FROM customer_credentials
		WHERE email = $1 AND tenant_id = $2`,
		email, string(tenantID),
	)
	creds, err := scanCredentials(row)
	if err == sql.ErrNoRows {
		return nil, ErrRegistry.New(CodeCustomerNotFound)
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning customer credentials by email", errx.TypeInternal)
	}
	return creds, nil
}

// GetByCustomerID retrieves credentials by customer ID scoped to tenant.
func (r *PostgresCredentialsRepo) GetByCustomerID(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (*CustomerCredentials, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, customer_id, tenant_id, email, password_hash, created_at, updated_at
		FROM customer_credentials
		WHERE customer_id = $1 AND tenant_id = $2`,
		string(customerID), string(tenantID),
	)
	creds, err := scanCredentials(row)
	if err == sql.ErrNoRows {
		return nil, ErrRegistry.New(CodeCustomerNotFound)
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning customer credentials by customer ID", errx.TypeInternal)
	}
	return creds, nil
}

// UpdatePassword updates the password hash for a customer.
func (r *PostgresCredentialsRepo) UpdatePassword(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, passwordHash string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE customer_credentials
		SET password_hash = $1, updated_at = $2
		WHERE customer_id = $3 AND tenant_id = $4`,
		passwordHash, time.Now(), string(customerID), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating customer password", errx.TypeInternal)
	}
	return nil
}

// scanCredentials scans a single row into a CustomerCredentials.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanCredentials(row rowScanner) (*CustomerCredentials, error) {
	var c CustomerCredentials
	var customerID, tenantID string

	err := row.Scan(
		&c.ID,
		&customerID,
		&tenantID,
		&c.Email,
		&c.PasswordHash,
		&c.CreatedAt,
		&c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	c.CustomerID = kernel.CustomerID(customerID)
	c.TenantID = kernel.TenantID(tenantID)
	return &c, nil
}

// Ensure interface compliance.
var _ CredentialsRepository = (*PostgresCredentialsRepo)(nil)
