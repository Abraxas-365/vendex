package customergroupinfra

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Abraxas-365/hada-commerce/internal/customergroup"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements customergroup.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed customer group repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance.
var _ customergroup.Repository = (*PostgresRepo)(nil)

// ─────────────────────────────────────────────────────────────────────────────
// CustomerGroup CRUD
// ─────────────────────────────────────────────────────────────────────────────

func (r *PostgresRepo) Create(ctx context.Context, g *customergroup.CustomerGroup) error {
	rulesJSON, err := json.Marshal(g.Rules)
	if err != nil {
		return errx.Wrap(err, "marshaling group rules", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO customer_groups (id, tenant_id, name, description, rules, auto_assign, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(g.ID), string(g.TenantID), g.Name, g.Description,
		string(rulesJSON), g.AutoAssign, g.CreatedAt, g.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting customer group", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerGroupID) (*customergroup.CustomerGroup, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, description, rules, auto_assign, created_at, updated_at
		FROM customer_groups
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	g, err := scanGroup(row)
	if err == sql.ErrNoRows {
		return nil, customergroup.ErrGroupNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning customer group", errx.TypeInternal)
	}
	return g, nil
}

func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID) ([]customergroup.CustomerGroup, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, description, rules, auto_assign, created_at, updated_at
		FROM customer_groups
		WHERE tenant_id = $1
		ORDER BY created_at DESC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying customer groups", errx.TypeInternal)
	}
	defer rows.Close()

	var groups []customergroup.CustomerGroup
	for rows.Next() {
		g, err := scanGroupRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning customer group row", errx.TypeInternal)
		}
		groups = append(groups, *g)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating customer groups", errx.TypeInternal)
	}
	if groups == nil {
		groups = []customergroup.CustomerGroup{}
	}
	return groups, nil
}

func (r *PostgresRepo) Update(ctx context.Context, g *customergroup.CustomerGroup) error {
	rulesJSON, err := json.Marshal(g.Rules)
	if err != nil {
		return errx.Wrap(err, "marshaling group rules", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE customer_groups
		SET name=$1, description=$2, rules=$3, auto_assign=$4, updated_at=$5
		WHERE id=$6 AND tenant_id=$7`,
		g.Name, g.Description, string(rulesJSON), g.AutoAssign, g.UpdatedAt,
		string(g.ID), string(g.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating customer group", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.CustomerGroupID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM customer_groups WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting customer group", errx.TypeInternal)
	}
	return nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Membership
// ─────────────────────────────────────────────────────────────────────────────

func (r *PostgresRepo) AddMember(ctx context.Context, m *customergroup.GroupMembership) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO customer_group_memberships (id, group_id, customer_id, tenant_id, assigned_at)
		VALUES ($1, $2, $3, $4, $5)`,
		string(m.ID), string(m.GroupID), string(m.CustomerID), string(m.TenantID), m.AssignedAt,
	)
	if err != nil {
		// Detect unique constraint violation (already a member).
		if isUniqueViolation(err) {
			return customergroup.ErrAlreadyMember
		}
		return errx.Wrap(err, "adding group member", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) RemoveMember(ctx context.Context, tenantID kernel.TenantID, groupID kernel.CustomerGroupID, customerID kernel.CustomerID) error {
	res, err := r.db.ExecContext(ctx, `
		DELETE FROM customer_group_memberships
		WHERE group_id=$1 AND customer_id=$2 AND tenant_id=$3`,
		string(groupID), string(customerID), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "removing group member", errx.TypeInternal)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return customergroup.ErrMemberNotFound
	}
	return nil
}

func (r *PostgresRepo) ListMembers(ctx context.Context, tenantID kernel.TenantID, groupID kernel.CustomerGroupID) ([]customergroup.GroupMembership, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, group_id, customer_id, tenant_id, assigned_at
		FROM customer_group_memberships
		WHERE group_id=$1 AND tenant_id=$2
		ORDER BY assigned_at DESC`,
		string(groupID), string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying group members", errx.TypeInternal)
	}
	defer rows.Close()

	var members []customergroup.GroupMembership
	for rows.Next() {
		var m customergroup.GroupMembership
		var id, groupID, customerID, tenantID string
		if err := rows.Scan(&id, &groupID, &customerID, &tenantID, &m.AssignedAt); err != nil {
			return nil, errx.Wrap(err, "scanning group member", errx.TypeInternal)
		}
		m.ID = kernel.CustomerGroupMembershipID(id)
		m.GroupID = kernel.CustomerGroupID(groupID)
		m.CustomerID = kernel.CustomerID(customerID)
		m.TenantID = kernel.TenantID(tenantID)
		members = append(members, m)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating group members", errx.TypeInternal)
	}
	if members == nil {
		members = []customergroup.GroupMembership{}
	}
	return members, nil
}

func (r *PostgresRepo) GetCustomerGroups(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) ([]customergroup.CustomerGroup, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT cg.id, cg.tenant_id, cg.name, cg.description, cg.rules, cg.auto_assign, cg.created_at, cg.updated_at
		FROM customer_groups cg
		INNER JOIN customer_group_memberships cgm ON cg.id = cgm.group_id
		WHERE cgm.customer_id=$1 AND cgm.tenant_id=$2
		ORDER BY cg.created_at DESC`,
		string(customerID), string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying customer's groups", errx.TypeInternal)
	}
	defer rows.Close()

	var groups []customergroup.CustomerGroup
	for rows.Next() {
		g, err := scanGroupRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning customer group row", errx.TypeInternal)
		}
		groups = append(groups, *g)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating customer groups", errx.TypeInternal)
	}
	if groups == nil {
		groups = []customergroup.CustomerGroup{}
	}
	return groups, nil
}

// ─────────────────────────────────────────────────────────────────────────────
// Helpers
// ─────────────────────────────────────────────────────────────────────────────

type rowScanner interface {
	Scan(dest ...any) error
}

func scanGroup(row *sql.Row) (*customergroup.CustomerGroup, error) {
	return scanGroupFields(row)
}

func scanGroupRow(rows *sql.Rows) (*customergroup.CustomerGroup, error) {
	return scanGroupFields(rows)
}

func scanGroupFields(s rowScanner) (*customergroup.CustomerGroup, error) {
	var g customergroup.CustomerGroup
	var id, tenantID, rulesJSON string

	err := s.Scan(&id, &tenantID, &g.Name, &g.Description, &rulesJSON, &g.AutoAssign, &g.CreatedAt, &g.UpdatedAt)
	if err != nil {
		return nil, err
	}

	g.ID = kernel.CustomerGroupID(id)
	g.TenantID = kernel.TenantID(tenantID)
	_ = json.Unmarshal([]byte(rulesJSON), &g.Rules)

	return &g, nil
}

// isUniqueViolation detects PostgreSQL unique-constraint violations (error code 23505).
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	return len(err.Error()) > 0 && contains(err.Error(), "23505")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsRune(s, substr))
}

func containsRune(s, substr string) bool {
	for i := range s {
		if i+len(substr) <= len(s) && s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
