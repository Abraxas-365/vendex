// Package approvalinfra provides a PostgreSQL implementation of the approval Repository.
package approvalinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/approval"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PostgresRepository implements approval.Repository backed by PostgreSQL.
type PostgresRepository struct{ db *sqlx.DB }

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// dbApprovalRequest is the sqlx-scannable row type for approval_requests.
type dbApprovalRequest struct {
	ID          string       `db:"id"`
	TenantID    string       `db:"tenant_id"`
	SessionID   string       `db:"session_id"`
	ToolName    string       `db:"tool_name"`
	ToolInput   []byte       `db:"tool_input"`
	Status      string       `db:"status"`
	Reason      string       `db:"reason"`
	RequestedBy string       `db:"requested_by"`
	ReviewedBy  string       `db:"reviewed_by"`
	CreatedAt   time.Time    `db:"created_at"`
	ReviewedAt  sql.NullTime `db:"reviewed_at"`
}

func fromDB(row dbApprovalRequest) approval.ApprovalRequest {
	toolInput := json.RawMessage(row.ToolInput)
	if len(toolInput) == 0 {
		toolInput = json.RawMessage("{}")
	}

	var reviewedAt *time.Time
	if row.ReviewedAt.Valid {
		t := row.ReviewedAt.Time
		reviewedAt = &t
	}

	return approval.ApprovalRequest{
		ID:          kernel.ApprovalRequestID(row.ID),
		TenantID:    kernel.TenantID(row.TenantID),
		SessionID:   row.SessionID,
		ToolName:    row.ToolName,
		ToolInput:   toolInput,
		Status:      row.Status,
		Reason:      row.Reason,
		RequestedBy: row.RequestedBy,
		ReviewedBy:  row.ReviewedBy,
		CreatedAt:   row.CreatedAt,
		ReviewedAt:  reviewedAt,
	}
}

// Create persists a new approval request.
func (r *PostgresRepository) Create(ctx context.Context, req approval.ApprovalRequest) (approval.ApprovalRequest, error) {
	toolInput := req.ToolInput
	if len(toolInput) == 0 {
		toolInput = json.RawMessage("{}")
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO approval_requests
			(id, tenant_id, session_id, tool_name, tool_input, status, reason,
			 requested_by, reviewed_by, created_at, reviewed_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		string(req.ID),
		string(req.TenantID),
		req.SessionID,
		req.ToolName,
		[]byte(toolInput),
		req.Status,
		req.Reason,
		req.RequestedBy,
		req.ReviewedBy,
		req.CreatedAt,
		req.ReviewedAt,
	)
	if err != nil {
		return approval.ApprovalRequest{}, errx.Wrap(err, "create approval request", errx.TypeInternal)
	}
	return req, nil
}

// GetByID retrieves a single approval request by ID, scoped to the tenant.
func (r *PostgresRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ApprovalRequestID) (approval.ApprovalRequest, error) {
	var row dbApprovalRequest
	err := r.db.GetContext(ctx, &row, `
		SELECT id, tenant_id, session_id, tool_name, tool_input, status, reason,
		       requested_by, reviewed_by, created_at, reviewed_at
		FROM approval_requests
		WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err == sql.ErrNoRows {
		return approval.ApprovalRequest{}, approval.ErrNotFound
	}
	if err != nil {
		return approval.ApprovalRequest{}, errx.Wrap(err, "get approval request", errx.TypeInternal)
	}
	return fromDB(row), nil
}

// List returns paginated approval requests, optionally filtered by status.
// Pass an empty status string to return all statuses.
func (r *PostgresRepository) List(ctx context.Context, tenantID kernel.TenantID, status string, p kernel.PaginationOptions) (kernel.Paginated[approval.ApprovalRequest], error) {
	var total int
	var countErr error
	if status != "" {
		countErr = r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM approval_requests WHERE tenant_id=$1 AND status=$2`,
			string(tenantID), status,
		).Scan(&total)
	} else {
		countErr = r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM approval_requests WHERE tenant_id=$1`,
			string(tenantID),
		).Scan(&total)
	}
	if countErr != nil {
		return kernel.Paginated[approval.ApprovalRequest]{}, errx.Wrap(countErr, "count approval requests", errx.TypeInternal)
	}

	var rows []dbApprovalRequest
	var queryErr error
	if status != "" {
		queryErr = r.db.SelectContext(ctx, &rows, `
			SELECT id, tenant_id, session_id, tool_name, tool_input, status, reason,
			       requested_by, reviewed_by, created_at, reviewed_at
			FROM approval_requests
			WHERE tenant_id=$1 AND status=$2
			ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
			string(tenantID), status, p.Limit(), p.Offset(),
		)
	} else {
		queryErr = r.db.SelectContext(ctx, &rows, `
			SELECT id, tenant_id, session_id, tool_name, tool_input, status, reason,
			       requested_by, reviewed_by, created_at, reviewed_at
			FROM approval_requests
			WHERE tenant_id=$1
			ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
			string(tenantID), p.Limit(), p.Offset(),
		)
	}
	if queryErr != nil {
		return kernel.Paginated[approval.ApprovalRequest]{}, errx.Wrap(queryErr, "list approval requests", errx.TypeInternal)
	}

	items := make([]approval.ApprovalRequest, len(rows))
	for i, row := range rows {
		items[i] = fromDB(row)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// UpdateStatus sets the status, reason, reviewedBy, and reviewedAt on an approval request.
func (r *PostgresRepository) UpdateStatus(
	ctx context.Context,
	tenantID kernel.TenantID,
	id kernel.ApprovalRequestID,
	status, reason, reviewedBy string,
) (approval.ApprovalRequest, error) {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE approval_requests
		SET status=$3, reason=$4, reviewed_by=$5, reviewed_at=$6
		WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
		status, reason, reviewedBy, now,
	)
	if err != nil {
		return approval.ApprovalRequest{}, errx.Wrap(err, "update approval request status", errx.TypeInternal)
	}
	return r.GetByID(ctx, tenantID, id)
}

// CountPending returns the number of pending approval requests for a tenant.
func (r *PostgresRepository) CountPending(ctx context.Context, tenantID kernel.TenantID) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM approval_requests WHERE tenant_id=$1 AND status='pending'`,
		string(tenantID),
	).Scan(&count)
	if err != nil {
		return 0, errx.Wrap(err, "count pending approval requests", errx.TypeInternal)
	}
	return count, nil
}

// Ensure interface compliance.
var _ approval.Repository = (*PostgresRepository)(nil)
