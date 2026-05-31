package agentsessioninfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/agentsession"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// ─── Session Repository ───────────────────────────────────────────────────────

// PostgresSessionRepo implements agentsession.SessionRepository.
type PostgresSessionRepo struct{ db *sqlx.DB }

// NewPostgresSessionRepo creates a new PostgresSessionRepo.
func NewPostgresSessionRepo(db *sqlx.DB) *PostgresSessionRepo {
	return &PostgresSessionRepo{db: db}
}

// dbSession is the sqlx-scannable row for an agent_session.
type dbSession struct {
	ID          string         `db:"id"`
	TenantID    string         `db:"tenant_id"`
	PresetID    string         `db:"preset_id"`
	ContainerID string         `db:"container_id"`
	Status      string         `db:"status"`
	FrontendURL string         `db:"frontend_url"`
	Metadata    []byte         `db:"metadata"`
	CreatedAt   time.Time      `db:"created_at"`
	UpdatedAt   time.Time      `db:"updated_at"`
	StoppedAt   sql.NullTime   `db:"stopped_at"`
}

func fromDBSession(row dbSession) agentsession.Session {
	meta := json.RawMessage(row.Metadata)
	if len(meta) == 0 {
		meta = json.RawMessage("{}")
	}
	var stoppedAt *time.Time
	if row.StoppedAt.Valid {
		t := row.StoppedAt.Time
		stoppedAt = &t
	}
	return agentsession.Session{
		ID:          kernel.AgentSessionID(row.ID),
		TenantID:    kernel.TenantID(row.TenantID),
		PresetID:    kernel.PresetID(row.PresetID),
		ContainerID: row.ContainerID,
		Status:      agentsession.SessionStatus(row.Status),
		FrontendURL: row.FrontendURL,
		Metadata:    meta,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
		StoppedAt:   stoppedAt,
	}
}

func (r *PostgresSessionRepo) Create(ctx context.Context, s agentsession.Session) (agentsession.Session, error) {
	meta := s.Metadata
	if len(meta) == 0 {
		meta = json.RawMessage("{}")
	}

	var stoppedAt interface{}
	if s.StoppedAt != nil {
		stoppedAt = *s.StoppedAt
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO agent_sessions
			(id, tenant_id, preset_id, container_id, status, frontend_url, metadata,
			 created_at, updated_at, stopped_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		string(s.ID), string(s.TenantID), string(s.PresetID), s.ContainerID,
		string(s.Status), s.FrontendURL, []byte(meta),
		s.CreatedAt, s.UpdatedAt, stoppedAt,
	)
	if err != nil {
		return agentsession.Session{}, errx.Wrap(err, "create agent session", errx.TypeInternal)
	}
	return s, nil
}

func (r *PostgresSessionRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentSessionID) (agentsession.Session, error) {
	var row dbSession
	err := r.db.GetContext(ctx, &row, `
		SELECT id, tenant_id, preset_id, container_id, status, frontend_url, metadata,
		       created_at, updated_at, stopped_at
		FROM agent_sessions
		WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err == sql.ErrNoRows {
		return agentsession.Session{}, agentsession.ErrSessionNotFound
	}
	if err != nil {
		return agentsession.Session{}, errx.Wrap(err, "get agent session", errx.TypeInternal)
	}
	return fromDBSession(row), nil
}

func (r *PostgresSessionRepo) Update(ctx context.Context, s agentsession.Session) (agentsession.Session, error) {
	meta := s.Metadata
	if len(meta) == 0 {
		meta = json.RawMessage("{}")
	}

	var stoppedAt interface{}
	if s.StoppedAt != nil {
		stoppedAt = *s.StoppedAt
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE agent_sessions
		SET container_id=$2, status=$3, frontend_url=$4, metadata=$5,
		    updated_at=$6, stopped_at=$7
		WHERE id=$1 AND tenant_id=$8`,
		string(s.ID), s.ContainerID, string(s.Status), s.FrontendURL,
		[]byte(meta), s.UpdatedAt, stoppedAt, string(s.TenantID),
	)
	if err != nil {
		return agentsession.Session{}, errx.Wrap(err, "update agent session", errx.TypeInternal)
	}
	return s, nil
}

func (r *PostgresSessionRepo) ListByTenant(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[agentsession.Session], error) {
	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM agent_sessions WHERE tenant_id=$1`, string(tenantID),
	).Scan(&total); err != nil {
		return kernel.Paginated[agentsession.Session]{}, errx.Wrap(err, "count agent sessions", errx.TypeInternal)
	}

	var rows []dbSession
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, preset_id, container_id, status, frontend_url, metadata,
		       created_at, updated_at, stopped_at
		FROM agent_sessions
		WHERE tenant_id=$1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[agentsession.Session]{}, errx.Wrap(err, "list agent sessions", errx.TypeInternal)
	}

	items := make([]agentsession.Session, len(rows))
	for i, row := range rows {
		items[i] = fromDBSession(row)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

func (r *PostgresSessionRepo) ListActive(ctx context.Context, tenantID kernel.TenantID) ([]agentsession.Session, error) {
	var rows []dbSession
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, preset_id, container_id, status, frontend_url, metadata,
		       created_at, updated_at, stopped_at
		FROM agent_sessions
		WHERE tenant_id=$1 AND status IN ('creating', 'running')
		ORDER BY created_at DESC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "list active agent sessions", errx.TypeInternal)
	}

	items := make([]agentsession.Session, len(rows))
	for i, row := range rows {
		items[i] = fromDBSession(row)
	}
	return items, nil
}

// Ensure interface compliance.
var _ agentsession.SessionRepository = (*PostgresSessionRepo)(nil)

// ─── Chat Repository ──────────────────────────────────────────────────────────

// PostgresChatRepo implements agentsession.ChatRepository.
type PostgresChatRepo struct{ db *sqlx.DB }

// NewPostgresChatRepo creates a new PostgresChatRepo.
func NewPostgresChatRepo(db *sqlx.DB) *PostgresChatRepo {
	return &PostgresChatRepo{db: db}
}

// dbChatMessage is the sqlx-scannable row for an agent_chat_message.
type dbChatMessage struct {
	ID        string    `db:"id"`
	SessionID string    `db:"session_id"`
	Role      string    `db:"role"`
	Content   string    `db:"content"`
	ToolName  string    `db:"tool_name"`
	CreatedAt time.Time `db:"created_at"`
}

func fromDBChatMessage(row dbChatMessage) agentsession.ChatMessage {
	return agentsession.ChatMessage{
		ID:        row.ID,
		SessionID: kernel.AgentSessionID(row.SessionID),
		Role:      row.Role,
		Content:   row.Content,
		ToolName:  row.ToolName,
		CreatedAt: row.CreatedAt,
	}
}

func (r *PostgresChatRepo) SaveMessage(ctx context.Context, msg agentsession.ChatMessage) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO agent_chat_messages (id, session_id, role, content, tool_name, created_at)
		VALUES ($1,$2,$3,$4,$5,$6)`,
		msg.ID, string(msg.SessionID), msg.Role, msg.Content, msg.ToolName, msg.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "save chat message", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresChatRepo) ListMessages(ctx context.Context, sessionID kernel.AgentSessionID, p kernel.PaginationOptions) (kernel.Paginated[agentsession.ChatMessage], error) {
	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM agent_chat_messages WHERE session_id=$1`, string(sessionID),
	).Scan(&total); err != nil {
		return kernel.Paginated[agentsession.ChatMessage]{}, errx.Wrap(err, "count chat messages", errx.TypeInternal)
	}

	var rows []dbChatMessage
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, session_id, role, content, tool_name, created_at
		FROM agent_chat_messages
		WHERE session_id=$1
		ORDER BY created_at ASC LIMIT $2 OFFSET $3`,
		string(sessionID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[agentsession.ChatMessage]{}, errx.Wrap(err, "list chat messages", errx.TypeInternal)
	}

	items := make([]agentsession.ChatMessage, len(rows))
	for i, row := range rows {
		items[i] = fromDBChatMessage(row)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// Ensure interface compliance.
var _ agentsession.ChatRepository = (*PostgresChatRepo)(nil)
