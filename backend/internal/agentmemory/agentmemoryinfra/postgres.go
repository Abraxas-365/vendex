// Package agentmemoryinfra implements the agentmemory.Repository using PostgreSQL.
package agentmemoryinfra

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	"github.com/Abraxas-365/hada-commerce/internal/agentmemory"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PostgresRepository implements agentmemory.Repository.
type PostgresRepository struct{ db *sqlx.DB }

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// dbMemory is the sqlx-scannable row for an agent_memory record.
type dbMemory struct {
	ID        string         `db:"id"`
	TenantID  string         `db:"tenant_id"`
	Category  string         `db:"category"`
	Title     string         `db:"title"`
	Content   string         `db:"content"`
	Tags      pq.StringArray `db:"tags"`
	Source    string         `db:"source"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

func fromDB(row dbMemory) agentmemory.Memory {
	tags := []string(row.Tags)
	if tags == nil {
		tags = []string{}
	}
	return agentmemory.Memory{
		ID:        kernel.AgentMemoryID(row.ID),
		TenantID:  kernel.TenantID(row.TenantID),
		Category:  row.Category,
		Title:     row.Title,
		Content:   row.Content,
		Tags:      tags,
		Source:    row.Source,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

// Create inserts a new memory record.
func (r *PostgresRepository) Create(ctx context.Context, m agentmemory.Memory) (agentmemory.Memory, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO agent_memories
			(id, tenant_id, category, title, content, tags, source, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		string(m.ID), string(m.TenantID), m.Category, m.Title, m.Content,
		pq.StringArray(m.Tags), m.Source, m.CreatedAt, m.UpdatedAt,
	)
	if err != nil {
		return agentmemory.Memory{}, errx.Wrap(err, "create agent memory", errx.TypeInternal)
	}
	return m, nil
}

// GetByID retrieves a memory by ID, scoped to a tenant.
func (r *PostgresRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentMemoryID) (agentmemory.Memory, error) {
	var row dbMemory
	err := r.db.GetContext(ctx, &row, `
		SELECT id, tenant_id, category, title, content, tags, source, created_at, updated_at
		FROM agent_memories
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err == sql.ErrNoRows {
		return agentmemory.Memory{}, agentmemory.ErrNotFound
	}
	if err != nil {
		return agentmemory.Memory{}, errx.Wrap(err, "get agent memory", errx.TypeInternal)
	}
	return fromDB(row), nil
}

// Update updates an existing memory record.
func (r *PostgresRepository) Update(ctx context.Context, m agentmemory.Memory) (agentmemory.Memory, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE agent_memories
		SET category = $1, title = $2, content = $3, tags = $4, updated_at = $5
		WHERE id = $6 AND tenant_id = $7`,
		m.Category, m.Title, m.Content,
		pq.StringArray(m.Tags), m.UpdatedAt,
		string(m.ID), string(m.TenantID),
	)
	if err != nil {
		return agentmemory.Memory{}, errx.Wrap(err, "update agent memory", errx.TypeInternal)
	}
	return m, nil
}

// Delete removes a memory record by ID, scoped to a tenant.
func (r *PostgresRepository) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentMemoryID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM agent_memories
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "delete agent memory", errx.TypeInternal)
	}
	return nil
}

// List returns paginated memories for a tenant.
func (r *PostgresRepository) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[agentmemory.Memory], error) {
	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM agent_memories WHERE tenant_id = $1`, string(tenantID),
	).Scan(&total); err != nil {
		return kernel.Paginated[agentmemory.Memory]{}, errx.Wrap(err, "count agent memories", errx.TypeInternal)
	}

	var rows []dbMemory
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, category, title, content, tags, source, created_at, updated_at
		FROM agent_memories
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[agentmemory.Memory]{}, errx.Wrap(err, "list agent memories", errx.TypeInternal)
	}

	items := make([]agentmemory.Memory, len(rows))
	for i, row := range rows {
		items[i] = fromDB(row)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// Search returns memories matching the given options, scoped to a tenant.
// It supports full-text search on title+content, category filtering,
// and tag filtering (all specified tags must match).
func (r *PostgresRepository) Search(ctx context.Context, tenantID kernel.TenantID, opts agentmemory.MemorySearchOptions, p kernel.PaginationOptions) (kernel.Paginated[agentmemory.Memory], error) {
	args := []interface{}{string(tenantID)}
	conditions := []string{"tenant_id = $1"}
	argIdx := 2

	if opts.Category != "" {
		conditions = append(conditions, "category = $"+itoa(argIdx))
		args = append(args, opts.Category)
		argIdx++
	}

	if len(opts.Tags) > 0 {
		conditions = append(conditions, "tags @> $"+itoa(argIdx))
		args = append(args, pq.StringArray(opts.Tags))
		argIdx++
	}

	if opts.Query != "" {
		conditions = append(conditions, "search_vector @@ plainto_tsquery('english', $"+itoa(argIdx)+")")
		args = append(args, opts.Query)
		argIdx++
	}

	where := strings.Join(conditions, " AND ")

	// Count total
	var total int
	countQuery := "SELECT COUNT(*) FROM agent_memories WHERE " + where
	if err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return kernel.Paginated[agentmemory.Memory]{}, errx.Wrap(err, "count agent memories search", errx.TypeInternal)
	}

	// Fetch page — rank by search relevance if a query was given
	orderBy := "created_at DESC"
	if opts.Query != "" {
		// argIdx-1 is the position of the query parameter added above
		orderBy = "ts_rank(search_vector, plainto_tsquery('english', $" + itoa(argIdx-1) + ")) DESC, created_at DESC"
	}

	listArgs := append(args, p.Limit(), p.Offset())
	selectQuery := `SELECT id, tenant_id, category, title, content, tags, source, created_at, updated_at
		FROM agent_memories
		WHERE ` + where + `
		ORDER BY ` + orderBy + `
		LIMIT $` + itoa(argIdx) + ` OFFSET $` + itoa(argIdx+1)

	var rows []dbMemory
	err := r.db.SelectContext(ctx, &rows, selectQuery, listArgs...)
	if err != nil {
		return kernel.Paginated[agentmemory.Memory]{}, errx.Wrap(err, "search agent memories", errx.TypeInternal)
	}

	items := make([]agentmemory.Memory, len(rows))
	for i, row := range rows {
		items[i] = fromDB(row)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// itoa converts an int to a string (avoids strconv import for small helper).
func itoa(i int) string {
	if i < 10 {
		return string(rune('0' + i))
	}
	// For larger values use strconv-style manual conversion
	var buf [20]byte
	pos := len(buf)
	for i >= 10 {
		pos--
		buf[pos] = byte('0' + i%10)
		i /= 10
	}
	pos--
	buf[pos] = byte('0' + i)
	return string(buf[pos:])
}

// Ensure interface compliance at compile time.
var _ agentmemory.Repository = (*PostgresRepository)(nil)
