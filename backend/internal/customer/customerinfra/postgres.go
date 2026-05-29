package customerinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/customer"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PostgresRepo implements customer.Repository using database/sql.
type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed customer repository.
func NewPostgresRepo(db *sql.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

func (r *PostgresRepo) Create(ctx context.Context, c *customer.Customer) error {
	addrJSON, err := json.Marshal(c.Addresses)
	if err != nil {
		return fmt.Errorf("marshaling addresses: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO customers (id, tenant_id, email, name, phone, addresses, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(c.ID), string(c.TenantID), c.Email.String(),
		c.Name, c.Phone, string(addrJSON),
		c.CreatedAt, c.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting customer: %w", err)
	}
	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) (*customer.Customer, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, email, name, phone, addresses, created_at, updated_at
		FROM customers WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	return scanCustomer(row)
}

func (r *PostgresRepo) GetByEmail(ctx context.Context, tenantID kernel.TenantID, email kernel.Email) (*customer.Customer, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, email, name, phone, addresses, created_at, updated_at
		FROM customers WHERE email = $1 AND tenant_id = $2`,
		email.String(), string(tenantID),
	)
	return scanCustomer(row)
}

func (r *PostgresRepo) Update(ctx context.Context, c *customer.Customer) error {
	addrJSON, err := json.Marshal(c.Addresses)
	if err != nil {
		return fmt.Errorf("marshaling addresses: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE customers SET email=$1, name=$2, phone=$3, addresses=$4, updated_at=$5
		WHERE id=$6 AND tenant_id=$7`,
		c.Email.String(), c.Name, c.Phone, string(addrJSON), c.UpdatedAt,
		string(c.ID), string(c.TenantID),
	)
	if err != nil {
		return fmt.Errorf("updating customer: %w", err)
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerID) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM customers WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return fmt.Errorf("deleting customer: %w", err)
	}
	return nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.Pagination) (kernel.PaginatedResult[customer.Customer], error) {
	var zero kernel.PaginatedResult[customer.Customer]

	var total int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM customers WHERE tenant_id = $1",
		string(tenantID),
	).Scan(&total); err != nil {
		return zero, fmt.Errorf("counting customers: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, email, name, phone, addresses, created_at, updated_at
		FROM customers WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, fmt.Errorf("querying customers: %w", err)
	}
	defer rows.Close()

	var customers []customer.Customer
	for rows.Next() {
		c, err := scanCustomerRow(rows)
		if err != nil {
			return zero, err
		}
		customers = append(customers, *c)
	}
	if err := rows.Err(); err != nil {
		return zero, fmt.Errorf("iterating customers: %w", err)
	}

	return kernel.NewPaginatedResult(customers, total, pg), nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanCustomer(row *sql.Row) (*customer.Customer, error) {
	c, err := scanFields(row)
	if err == sql.ErrNoRows {
		return nil, customer.ErrNotFound
	}
	return c, err
}

func scanCustomerRow(rows *sql.Rows) (*customer.Customer, error) {
	return scanFields(rows)
}

func scanFields(s scanner) (*customer.Customer, error) {
	var c customer.Customer
	var id, tenantID, email string
	var addrJSON string
	var createdAt, updatedAt time.Time

	err := s.Scan(&id, &tenantID, &email, &c.Name, &c.Phone, &addrJSON, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	c.ID = kernel.CustomerID(id)
	c.TenantID = kernel.TenantID(tenantID)
	c.Email = kernel.Email(email)
	c.CreatedAt = createdAt
	c.UpdatedAt = updatedAt

	_ = json.Unmarshal([]byte(addrJSON), &c.Addresses)
	if c.Addresses == nil {
		c.Addresses = []customer.Address{}
	}

	return &c, nil
}

// Ensure interface compliance.
var _ customer.Repository = (*PostgresRepo)(nil)
