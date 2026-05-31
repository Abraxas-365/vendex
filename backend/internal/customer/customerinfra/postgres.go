package customerinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Abraxas-365/vendex/internal/customer"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements customer.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed customer repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, c *customer.Customer) error {
	addrJSON, err := json.Marshal(c.Addresses)
	if err != nil {
		return errx.Wrap(err, "marshaling addresses", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO customers (id, tenant_id, email, name, phone, addresses, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(c.ID), string(c.TenantID), c.Email.String(),
		c.Name, c.Phone, string(addrJSON),
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting customer", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) (*customer.Customer, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, email, name, phone, addresses, created_at, updated_at
		FROM customers WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	c, err := scanCustomer(row)
	if err == sql.ErrNoRows {
		return nil, customer.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning customer", errx.TypeInternal)
	}
	return c, nil
}

func (r *PostgresRepo) GetByEmail(ctx context.Context, tenantID kernel.TenantID, email kernel.Email) (*customer.Customer, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, email, name, phone, addresses, created_at, updated_at
		FROM customers WHERE email = $1 AND tenant_id = $2`,
		email.String(), string(tenantID),
	)
	c, err := scanCustomer(row)
	if err == sql.ErrNoRows {
		return nil, customer.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning customer by email", errx.TypeInternal)
	}
	return c, nil
}

func (r *PostgresRepo) Update(ctx context.Context, c *customer.Customer) error {
	addrJSON, err := json.Marshal(c.Addresses)
	if err != nil {
		return errx.Wrap(err, "marshaling addresses", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE customers SET email=$1, name=$2, phone=$3, addresses=$4, updated_at=$5
		WHERE id=$6 AND tenant_id=$7`,
		c.Email.String(), c.Name, c.Phone, string(addrJSON), c.UpdatedAt,
		string(c.ID), string(c.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating customer", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM customers WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting customer", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[customer.Customer], error) {
	var zero kernel.Paginated[customer.Customer]

	var total int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM customers WHERE tenant_id = $1",
		string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting customers", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, email, name, phone, addresses, created_at, updated_at
		FROM customers WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying customers", errx.TypeInternal)
	}
	defer rows.Close()

	var customers []customer.Customer
	for rows.Next() {
		c, err := scanCustomerRow(rows)
		if err != nil {
			return zero, errx.Wrap(err, "scanning customer row", errx.TypeInternal)
		}
		customers = append(customers, *c)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating customers", errx.TypeInternal)
	}

	return kernel.NewPaginated(customers, pg.Page, pg.PageSize, total), nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCustomer(row *sql.Row) (*customer.Customer, error) {
	return scanFields(row)
}

func scanCustomerRow(rows *sql.Rows) (*customer.Customer, error) {
	return scanFields(rows)
}

func scanFields(s scanner) (*customer.Customer, error) {
	var c customer.Customer
	var id, tenantID, email string
	var addrJSON string

	err := s.Scan(&id, &tenantID, &email, &c.Name, &c.Phone, &addrJSON, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	c.ID = kernel.CustomerID(id)
	c.TenantID = kernel.TenantID(tenantID)
	c.Email = kernel.Email(email)

	_ = json.Unmarshal([]byte(addrJSON), &c.Addresses)
	if c.Addresses == nil {
		c.Addresses = []customer.Address{}
	}

	return &c, nil
}

// Ensure interface compliance.
var _ customer.Repository = (*PostgresRepo)(nil)

// scanCustomerFields is an alias for fmt package usage.
var _ = fmt.Sprintf
