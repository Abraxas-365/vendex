package bulkopsinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Abraxas-365/vendex/internal/bulkops"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements bulkops.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed bulk operations repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Compile-time interface check.
var _ bulkops.Repository = (*PostgresRepo)(nil)

// ---------------------------------------------------------------------------
// Create
// ---------------------------------------------------------------------------

// Create persists a new bulk operation and its items atomically.
func (r *PostgresRepo) Create(ctx context.Context, op *bulkops.BulkOperation, items []bulkops.BulkOperationItem) error {
	paramsJSON, err := json.Marshal(op.Parameters)
	if err != nil {
		return errx.Wrap(err, "marshalling parameters", errx.TypeInternal)
	}
	errorsJSON, err := json.Marshal(op.Errors)
	if err != nil {
		return errx.Wrap(err, "marshalling errors", errx.TypeInternal)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return errx.Wrap(err, "beginning transaction", errx.TypeInternal)
	}
	defer func() { _ = tx.Rollback() }()

	_, err = tx.ExecContext(ctx, `
		INSERT INTO bulk_operations (
			id, tenant_id, type, resource_type, status,
			total_items, processed_items, failed_items,
			parameters, errors, created_by, started_at, completed_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8,
			$9, $10, $11, $12, $13, $14
		)`,
		string(op.ID), string(op.TenantID), string(op.Type), op.ResourceType, string(op.Status),
		op.TotalItems, op.ProcessedItems, op.FailedItems,
		paramsJSON, errorsJSON, op.CreatedBy, op.StartedAt, op.CompletedAt, op.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting bulk operation", errx.TypeInternal)
	}

	for _, item := range items {
		_, err = tx.ExecContext(ctx, `
			INSERT INTO bulk_operation_items (
				id, tenant_id, operation_id, resource_id, status,
				error_message, processed_at, created_at
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
			string(item.ID), string(item.TenantID), string(item.OperationID),
			item.ResourceID, string(item.Status),
			nullString(item.ErrorMessage), item.ProcessedAt, item.CreatedAt,
		)
		if err != nil {
			return errx.Wrap(err, "inserting bulk operation item", errx.TypeInternal)
		}
	}

	if err := tx.Commit(); err != nil {
		return errx.Wrap(err, "committing transaction", errx.TypeInternal)
	}
	return nil
}

// ---------------------------------------------------------------------------
// GetByID
// ---------------------------------------------------------------------------

// GetByID returns a bulk operation scoped to the tenant.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BulkOperationID) (*bulkops.BulkOperation, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, type, resource_type, status,
		       total_items, processed_items, failed_items,
		       parameters, errors, created_by, started_at, completed_at, created_at
		FROM bulk_operations
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)

	op, err := scanOperation(row.Scan)
	if err == sql.ErrNoRows {
		return nil, bulkops.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning bulk operation", errx.TypeInternal)
	}
	return op, nil
}

// ---------------------------------------------------------------------------
// List
// ---------------------------------------------------------------------------

// List returns paginated bulk operations for a tenant.
func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, page, pageSize int) (kernel.Paginated[bulkops.BulkOperation], error) {
	var zero kernel.Paginated[bulkops.BulkOperation]
	pg := kernel.NewPaginationOptions(page, pageSize)

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM bulk_operations WHERE tenant_id = $1`,
		string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting bulk operations", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, type, resource_type, status,
		       total_items, processed_items, failed_items,
		       parameters, errors, created_by, started_at, completed_at, created_at
		FROM bulk_operations
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying bulk operations", errx.TypeInternal)
	}
	defer rows.Close()

	items := make([]bulkops.BulkOperation, 0)
	for rows.Next() {
		op, err := scanOperation(rows.Scan)
		if err != nil {
			return zero, errx.Wrap(err, "scanning bulk operation row", errx.TypeInternal)
		}
		items = append(items, *op)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating bulk operations", errx.TypeInternal)
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

// ---------------------------------------------------------------------------
// ListItems
// ---------------------------------------------------------------------------

// ListItems returns paginated items for a specific bulk operation.
func (r *PostgresRepo) ListItems(ctx context.Context, tenantID kernel.TenantID, operationID kernel.BulkOperationID, page, pageSize int) (kernel.Paginated[bulkops.BulkOperationItem], error) {
	var zero kernel.Paginated[bulkops.BulkOperationItem]
	pg := kernel.NewPaginationOptions(page, pageSize)

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM bulk_operation_items WHERE operation_id = $1 AND tenant_id = $2`,
		string(operationID), string(tenantID),
	).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting bulk operation items", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, operation_id, resource_id, status,
		       error_message, processed_at, created_at
		FROM bulk_operation_items
		WHERE operation_id = $1 AND tenant_id = $2
		ORDER BY created_at ASC
		LIMIT $3 OFFSET $4`,
		string(operationID), string(tenantID), pg.Limit(), pg.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "querying bulk operation items", errx.TypeInternal)
	}
	defer rows.Close()

	items := make([]bulkops.BulkOperationItem, 0)
	for rows.Next() {
		item, err := scanItem(rows.Scan)
		if err != nil {
			return zero, errx.Wrap(err, "scanning bulk operation item row", errx.TypeInternal)
		}
		items = append(items, *item)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating bulk operation items", errx.TypeInternal)
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

// ---------------------------------------------------------------------------
// UpdateStatus
// ---------------------------------------------------------------------------

// UpdateStatus updates only the status of a bulk operation.
func (r *PostgresRepo) UpdateStatus(ctx context.Context, tenantID kernel.TenantID, id kernel.BulkOperationID, status bulkops.OperationStatus) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE bulk_operations SET status = $1 WHERE id = $2 AND tenant_id = $3`,
		string(status), string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating bulk operation status", errx.TypeInternal)
	}
	return nil
}

// ---------------------------------------------------------------------------
// UpdateOperation
// ---------------------------------------------------------------------------

// UpdateOperation persists the full mutable state of a bulk operation.
func (r *PostgresRepo) UpdateOperation(ctx context.Context, op *bulkops.BulkOperation) error {
	paramsJSON, err := json.Marshal(op.Parameters)
	if err != nil {
		return errx.Wrap(err, "marshalling parameters", errx.TypeInternal)
	}
	opErrors := op.Errors
	if opErrors == nil {
		opErrors = []bulkops.OperationError{}
	}
	errorsJSON, err := json.Marshal(opErrors)
	if err != nil {
		return errx.Wrap(err, "marshalling errors", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE bulk_operations
		SET status          = $1,
		    processed_items = $2,
		    failed_items    = $3,
		    parameters      = $4,
		    errors          = $5,
		    started_at      = $6,
		    completed_at    = $7
		WHERE id = $8 AND tenant_id = $9`,
		string(op.Status), op.ProcessedItems, op.FailedItems,
		paramsJSON, errorsJSON,
		op.StartedAt, op.CompletedAt,
		string(op.ID), string(op.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating bulk operation", errx.TypeInternal)
	}
	return nil
}

// ---------------------------------------------------------------------------
// UpdateItem
// ---------------------------------------------------------------------------

// UpdateItem persists the result of a single bulk operation item.
func (r *PostgresRepo) UpdateItem(ctx context.Context, item *bulkops.BulkOperationItem) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE bulk_operation_items
		SET status        = $1,
		    error_message = $2,
		    processed_at  = $3
		WHERE id = $4 AND tenant_id = $5`,
		string(item.Status), nullString(item.ErrorMessage), item.ProcessedAt,
		string(item.ID), string(item.TenantID),
	)
	if err != nil {
		return errx.Wrap(err, "updating bulk operation item", errx.TypeInternal)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Scan helpers
// ---------------------------------------------------------------------------

type scanFunc func(dest ...any) error

func scanOperation(scan scanFunc) (*bulkops.BulkOperation, error) {
	var (
		op           bulkops.BulkOperation
		id, tenantID, opType, resourceType, status string
		paramsJSON, errorsJSON                      []byte
		startedAt, completedAt                      sql.NullTime
	)

	err := scan(
		&id, &tenantID, &opType, &resourceType, &status,
		&op.TotalItems, &op.ProcessedItems, &op.FailedItems,
		&paramsJSON, &errorsJSON, &op.CreatedBy,
		&startedAt, &completedAt, &op.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	op.ID = kernel.BulkOperationID(id)
	op.TenantID = kernel.TenantID(tenantID)
	op.Type = bulkops.OperationType(opType)
	op.ResourceType = resourceType
	op.Status = bulkops.OperationStatus(status)

	if startedAt.Valid {
		t := startedAt.Time
		op.StartedAt = &t
	}
	if completedAt.Valid {
		t := completedAt.Time
		op.CompletedAt = &t
	}

	if err := json.Unmarshal(paramsJSON, &op.Parameters); err != nil {
		op.Parameters = map[string]interface{}{}
	}
	if len(errorsJSON) > 0 {
		if err := json.Unmarshal(errorsJSON, &op.Errors); err != nil {
			op.Errors = []bulkops.OperationError{}
		}
	} else {
		op.Errors = []bulkops.OperationError{}
	}

	return &op, nil
}

func scanItem(scan scanFunc) (*bulkops.BulkOperationItem, error) {
	var (
		item         bulkops.BulkOperationItem
		id, tenantID, operationID, resourceID, status string
		errorMessage sql.NullString
		processedAt  sql.NullTime
	)

	err := scan(
		&id, &tenantID, &operationID, &resourceID, &status,
		&errorMessage, &processedAt, &item.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.ID = kernel.BulkOperationItemID(id)
	item.TenantID = kernel.TenantID(tenantID)
	item.OperationID = kernel.BulkOperationID(operationID)
	item.ResourceID = resourceID
	item.Status = bulkops.ItemStatus(status)

	if errorMessage.Valid {
		item.ErrorMessage = errorMessage.String
	}
	if processedAt.Valid {
		t := processedAt.Time
		item.ProcessedAt = &t
	}

	return &item, nil
}

// nullString converts an empty string to a sql.NullString.
func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// nullTime converts a *time.Time to sql.NullTime.
func nullTime(t *time.Time) sql.NullTime {
	if t == nil {
		return sql.NullTime{}
	}
	return sql.NullTime{Time: *t, Valid: true}
}

// ensure nullTime is used (suppress unused warning)
var _ = nullTime
