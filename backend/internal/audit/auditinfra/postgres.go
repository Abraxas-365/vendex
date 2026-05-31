package auditinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/audit"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// PostgresRepo implements audit.Repository using PostgreSQL.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo returns a new PostgreSQL-backed audit repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Create inserts a new audit entry and returns it.
func (r *PostgresRepo) Create(ctx context.Context, entry audit.AuditEntry) (audit.AuditEntry, error) {
	changesJSON, err := marshalJSON(entry.Changes)
	if err != nil {
		return audit.AuditEntry{}, errx.Wrap(err, "marshaling audit changes", errx.TypeInternal)
	}
	metaJSON, err := marshalJSON(entry.Metadata)
	if err != nil {
		return audit.AuditEntry{}, errx.Wrap(err, "marshaling audit metadata", errx.TypeInternal)
	}

	const q = `
		INSERT INTO audit_logs
			(id, tenant_id, user_id, user_email, action, resource_type, resource_id,
			 changes, metadata, ip_address, user_agent, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err = r.db.ExecContext(ctx, q,
		string(entry.ID),
		string(entry.TenantID),
		entry.UserID,
		nullableString(entry.UserEmail),
		entry.Action,
		entry.ResourceType,
		nullableString(entry.ResourceID),
		changesJSON,
		metaJSON,
		nullableString(entry.IPAddress),
		nullableString(entry.UserAgent),
		entry.CreatedAt,
	)
	if err != nil {
		return audit.AuditEntry{}, errx.Wrap(err, "inserting audit entry", errx.TypeInternal)
	}

	return entry, nil
}

// GetByID retrieves a single audit entry scoped to the tenant.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.AuditEntryID) (audit.AuditEntry, error) {
	const q = `
		SELECT id, tenant_id, user_id, COALESCE(user_email,''), action, resource_type,
		       COALESCE(resource_id,''), changes, metadata,
		       COALESCE(ip_address,''), COALESCE(user_agent,''), created_at
		FROM audit_logs
		WHERE id = $1 AND tenant_id = $2
	`

	var (
		idStr, tidStr, userID, userEmail, action, resType, resID, ipAddr, userAgent string
		changesRaw, metaRaw                                                          []byte
		createdAt                                                                    time.Time
	)

	err := r.db.QueryRowContext(ctx, q, string(id), string(tenantID)).Scan(
		&idStr, &tidStr, &userID, &userEmail, &action, &resType, &resID,
		&changesRaw, &metaRaw,
		&ipAddr, &userAgent, &createdAt,
	)
	if err == sql.ErrNoRows {
		return audit.AuditEntry{}, errx.Wrap(audit.ErrNotFound, string(id), errx.TypeNotFound)
	}
	if err != nil {
		return audit.AuditEntry{}, errx.Wrap(err, "querying audit entry", errx.TypeInternal)
	}

	changes, err := unmarshalJSON(changesRaw)
	if err != nil {
		return audit.AuditEntry{}, errx.Wrap(err, "unmarshaling audit changes", errx.TypeInternal)
	}
	meta, err := unmarshalJSON(metaRaw)
	if err != nil {
		return audit.AuditEntry{}, errx.Wrap(err, "unmarshaling audit metadata", errx.TypeInternal)
	}

	return audit.AuditEntry{
		ID:           kernel.AuditEntryID(idStr),
		TenantID:     kernel.TenantID(tidStr),
		UserID:       userID,
		UserEmail:    userEmail,
		Action:       action,
		ResourceType: resType,
		ResourceID:   resID,
		Changes:      changes,
		Metadata:     meta,
		IPAddress:    ipAddr,
		UserAgent:    userAgent,
		CreatedAt:    createdAt.UTC(),
	}, nil
}

// List returns a paginated, filtered list of audit entries for the tenant.
func (r *PostgresRepo) List(
	ctx context.Context,
	tenantID kernel.TenantID,
	filter audit.AuditFilter,
	p kernel.PaginationOptions,
) (kernel.Paginated[audit.AuditEntry], error) {
	// Build dynamic WHERE clause.
	args := []any{string(tenantID)}
	where := "WHERE tenant_id = $1"
	idx := 2

	if filter.UserID != "" {
		where += fmt.Sprintf(" AND user_id = $%d", idx)
		args = append(args, filter.UserID)
		idx++
	}
	if filter.Action != "" {
		where += fmt.Sprintf(" AND action = $%d", idx)
		args = append(args, filter.Action)
		idx++
	}
	if filter.ResourceType != "" {
		where += fmt.Sprintf(" AND resource_type = $%d", idx)
		args = append(args, filter.ResourceType)
		idx++
	}
	if filter.ResourceID != "" {
		where += fmt.Sprintf(" AND resource_id = $%d", idx)
		args = append(args, filter.ResourceID)
		idx++
	}
	if filter.From != nil {
		where += fmt.Sprintf(" AND created_at >= $%d", idx)
		args = append(args, *filter.From)
		idx++
	}
	if filter.To != nil {
		where += fmt.Sprintf(" AND created_at <= $%d", idx)
		args = append(args, *filter.To)
		idx++
	}

	// Count.
	var total int
	if err := r.db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM audit_logs "+where, args...,
	).Scan(&total); err != nil {
		return kernel.Paginated[audit.AuditEntry]{}, errx.Wrap(err, "counting audit entries", errx.TypeInternal)
	}

	// Data.
	dataArgs := append(args, p.Limit(), p.Offset())
	orderClause := fmt.Sprintf(" ORDER BY created_at DESC LIMIT $%d OFFSET $%d", idx, idx+1)

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, tenant_id, user_id, COALESCE(user_email,''), action, resource_type,
		        COALESCE(resource_id,''), changes, metadata,
		        COALESCE(ip_address,''), COALESCE(user_agent,''), created_at
		 FROM audit_logs `+where+orderClause, dataArgs...)
	if err != nil {
		return kernel.Paginated[audit.AuditEntry]{}, errx.Wrap(err, "listing audit entries", errx.TypeInternal)
	}
	defer rows.Close()

	var items []audit.AuditEntry
	for rows.Next() {
		var (
			idStr, tidStr, userID, userEmail, action, resType, resID, ipAddr, userAgent string
			changesRaw, metaRaw                                                          []byte
			createdAt                                                                    time.Time
		)
		if err := rows.Scan(
			&idStr, &tidStr, &userID, &userEmail, &action, &resType, &resID,
			&changesRaw, &metaRaw,
			&ipAddr, &userAgent, &createdAt,
		); err != nil {
			return kernel.Paginated[audit.AuditEntry]{}, errx.Wrap(err, "scanning audit entry", errx.TypeInternal)
		}

		changes, err := unmarshalJSON(changesRaw)
		if err != nil {
			return kernel.Paginated[audit.AuditEntry]{}, errx.Wrap(err, "unmarshaling changes", errx.TypeInternal)
		}
		meta, err := unmarshalJSON(metaRaw)
		if err != nil {
			return kernel.Paginated[audit.AuditEntry]{}, errx.Wrap(err, "unmarshaling metadata", errx.TypeInternal)
		}

		items = append(items, audit.AuditEntry{
			ID:           kernel.AuditEntryID(idStr),
			TenantID:     kernel.TenantID(tidStr),
			UserID:       userID,
			UserEmail:    userEmail,
			Action:       action,
			ResourceType: resType,
			ResourceID:   resID,
			Changes:      changes,
			Metadata:     meta,
			IPAddress:    ipAddr,
			UserAgent:    userAgent,
			CreatedAt:    createdAt.UTC(),
		})
	}
	if err := rows.Err(); err != nil {
		return kernel.Paginated[audit.AuditEntry]{}, errx.Wrap(err, "iterating audit entries", errx.TypeInternal)
	}

	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// CountByAction returns the count of audit entries grouped by action within a time range.
func (r *PostgresRepo) CountByAction(
	ctx context.Context,
	tenantID kernel.TenantID,
	from, to time.Time,
) ([]audit.ActionStats, error) {
	const q = `
		SELECT action, COUNT(*) AS count
		FROM audit_logs
		WHERE tenant_id = $1 AND created_at >= $2 AND created_at <= $3
		GROUP BY action
		ORDER BY count DESC
	`
	rows, err := r.db.QueryContext(ctx, q, string(tenantID), from, to)
	if err != nil {
		return nil, errx.Wrap(err, "counting audit by action", errx.TypeInternal)
	}
	defer rows.Close()

	var stats []audit.ActionStats
	for rows.Next() {
		var s audit.ActionStats
		if err := rows.Scan(&s.Action, &s.Count); err != nil {
			return nil, errx.Wrap(err, "scanning action stats", errx.TypeInternal)
		}
		stats = append(stats, s)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating action stats", errx.TypeInternal)
	}
	return stats, nil
}

// Compile-time interface check.
var _ audit.Repository = (*PostgresRepo)(nil)

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

func marshalJSON(v map[string]any) ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	return json.Marshal(v)
}

func unmarshalJSON(b []byte) (map[string]any, error) {
	if len(b) == 0 || string(b) == "null" {
		return nil, nil
	}
	var m map[string]any
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}
