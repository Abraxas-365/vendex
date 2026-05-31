// Package agenttriggerinfra provides Postgres implementations of the agenttrigger repositories.
package agenttriggerinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/agenttrigger/agenttrigger"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// ─── Trigger Repository ────────────────────────────────────────────────────────

// PostgresTriggerRepo implements agenttrigger.TriggerRepository.
type PostgresTriggerRepo struct{ db *sqlx.DB }

// NewPostgresTriggerRepo creates a new PostgresTriggerRepo.
func NewPostgresTriggerRepo(db *sqlx.DB) *PostgresTriggerRepo {
	return &PostgresTriggerRepo{db: db}
}

// dbTrigger is the sqlx-scannable row for an agent_trigger.
type dbTrigger struct {
	ID          string       `db:"id"`
	TenantID    string       `db:"tenant_id"`
	Name        string       `db:"name"`
	EventType   string       `db:"event_type"`
	Prompt      string       `db:"prompt"`
	PresetID    string       `db:"preset_id"`
	Enabled     bool         `db:"enabled"`
	Cooldown    int          `db:"cooldown"`
	LastFiredAt sql.NullTime `db:"last_fired_at"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
}

func fromDBTrigger(row dbTrigger) agenttrigger.Trigger {
	var lastFiredAt *time.Time
	if row.LastFiredAt.Valid {
		t := row.LastFiredAt.Time
		lastFiredAt = &t
	}
	return agenttrigger.Trigger{
		ID:          kernel.AgentTriggerID(row.ID),
		TenantID:    kernel.TenantID(row.TenantID),
		Name:        row.Name,
		EventType:   row.EventType,
		Prompt:      row.Prompt,
		PresetID:    row.PresetID,
		Enabled:     row.Enabled,
		Cooldown:    row.Cooldown,
		LastFiredAt: lastFiredAt,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func (r *PostgresTriggerRepo) Create(ctx context.Context, t agenttrigger.Trigger) (agenttrigger.Trigger, error) {
	var lastFiredAt interface{}
	if t.LastFiredAt != nil {
		lastFiredAt = *t.LastFiredAt
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO agent_triggers
			(id, tenant_id, name, event_type, prompt, preset_id, enabled, cooldown, last_fired_at, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`,
		string(t.ID), string(t.TenantID), t.Name, t.EventType, t.Prompt,
		t.PresetID, t.Enabled, t.Cooldown, lastFiredAt, t.CreatedAt, t.UpdatedAt,
	)
	if err != nil {
		return agenttrigger.Trigger{}, errx.Wrap(err, "create agent trigger", errx.TypeInternal)
	}
	return t, nil
}

func (r *PostgresTriggerRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentTriggerID) (agenttrigger.Trigger, error) {
	var row dbTrigger
	err := r.db.GetContext(ctx, &row, `
		SELECT id, tenant_id, name, event_type, prompt, preset_id, enabled, cooldown,
		       last_fired_at, created_at, updated_at
		FROM agent_triggers
		WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err == sql.ErrNoRows {
		return agenttrigger.Trigger{}, agenttrigger.ErrNotFound
	}
	if err != nil {
		return agenttrigger.Trigger{}, errx.Wrap(err, "get agent trigger", errx.TypeInternal)
	}
	return fromDBTrigger(row), nil
}

func (r *PostgresTriggerRepo) Update(ctx context.Context, t agenttrigger.Trigger) (agenttrigger.Trigger, error) {
	var lastFiredAt interface{}
	if t.LastFiredAt != nil {
		lastFiredAt = *t.LastFiredAt
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE agent_triggers
		SET name=$2, event_type=$3, prompt=$4, preset_id=$5, enabled=$6,
		    cooldown=$7, last_fired_at=$8, updated_at=$9
		WHERE id=$1 AND tenant_id=$10`,
		string(t.ID), t.Name, t.EventType, t.Prompt, t.PresetID,
		t.Enabled, t.Cooldown, lastFiredAt, t.UpdatedAt, string(t.TenantID),
	)
	if err != nil {
		return agenttrigger.Trigger{}, errx.Wrap(err, "update agent trigger", errx.TypeInternal)
	}
	return t, nil
}

func (r *PostgresTriggerRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.AgentTriggerID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM agent_triggers WHERE id=$1 AND tenant_id=$2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "delete agent trigger", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresTriggerRepo) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[agenttrigger.Trigger], error) {
	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM agent_triggers WHERE tenant_id=$1`, string(tenantID),
	).Scan(&total); err != nil {
		return kernel.Paginated[agenttrigger.Trigger]{}, errx.Wrap(err, "count agent triggers", errx.TypeInternal)
	}

	var rows []dbTrigger
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, name, event_type, prompt, preset_id, enabled, cooldown,
		       last_fired_at, created_at, updated_at
		FROM agent_triggers
		WHERE tenant_id=$1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[agenttrigger.Trigger]{}, errx.Wrap(err, "list agent triggers", errx.TypeInternal)
	}

	items := make([]agenttrigger.Trigger, len(rows))
	for i, row := range rows {
		items[i] = fromDBTrigger(row)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// ListByEventType returns all enabled triggers for an event type across all tenants.
func (r *PostgresTriggerRepo) ListByEventType(ctx context.Context, eventType string) ([]agenttrigger.Trigger, error) {
	var rows []dbTrigger
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, name, event_type, prompt, preset_id, enabled, cooldown,
		       last_fired_at, created_at, updated_at
		FROM agent_triggers
		WHERE event_type=$1 AND enabled=true
		ORDER BY created_at ASC`,
		eventType,
	)
	if err != nil {
		return nil, errx.Wrap(err, "list triggers by event type", errx.TypeInternal)
	}

	items := make([]agenttrigger.Trigger, len(rows))
	for i, row := range rows {
		items[i] = fromDBTrigger(row)
	}
	return items, nil
}

// UpdateLastFired updates only the last_fired_at timestamp for cooldown tracking.
func (r *PostgresTriggerRepo) UpdateLastFired(ctx context.Context, id kernel.AgentTriggerID, firedAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE agent_triggers SET last_fired_at=$2, updated_at=$3 WHERE id=$1`,
		string(id), firedAt, firedAt,
	)
	if err != nil {
		return errx.Wrap(err, "update trigger last_fired_at", errx.TypeInternal)
	}
	return nil
}

// Ensure interface compliance.
var _ agenttrigger.TriggerRepository = (*PostgresTriggerRepo)(nil)

// ─── Trigger Log Repository ────────────────────────────────────────────────────

// PostgresTriggerLogRepo implements agenttrigger.TriggerLogRepository.
type PostgresTriggerLogRepo struct{ db *sqlx.DB }

// NewPostgresTriggerLogRepo creates a new PostgresTriggerLogRepo.
func NewPostgresTriggerLogRepo(db *sqlx.DB) *PostgresTriggerLogRepo {
	return &PostgresTriggerLogRepo{db: db}
}

// dbTriggerLog is the sqlx-scannable row for an agent_trigger_log.
type dbTriggerLog struct {
	ID            string    `db:"id"`
	TriggerID     string    `db:"trigger_id"`
	TenantID      string    `db:"tenant_id"`
	EventType     string    `db:"event_type"`
	EventPayload  []byte    `db:"event_payload"`
	AgentResponse string    `db:"agent_response"`
	Status        string    `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
}

func fromDBTriggerLog(row dbTriggerLog) agenttrigger.TriggerLog {
	payload := json.RawMessage(row.EventPayload)
	if len(payload) == 0 {
		payload = json.RawMessage("{}")
	}
	return agenttrigger.TriggerLog{
		ID:            kernel.TriggerLogID(row.ID),
		TriggerID:     kernel.AgentTriggerID(row.TriggerID),
		TenantID:      kernel.TenantID(row.TenantID),
		EventType:     row.EventType,
		EventPayload:  payload,
		AgentResponse: row.AgentResponse,
		Status:        row.Status,
		CreatedAt:     row.CreatedAt,
	}
}

func (r *PostgresTriggerLogRepo) Create(ctx context.Context, log agenttrigger.TriggerLog) (agenttrigger.TriggerLog, error) {
	payload := log.EventPayload
	if len(payload) == 0 {
		payload = json.RawMessage("{}")
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO agent_trigger_logs
			(id, trigger_id, tenant_id, event_type, event_payload, agent_response, status, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		string(log.ID), string(log.TriggerID), string(log.TenantID),
		log.EventType, []byte(payload), log.AgentResponse, log.Status, log.CreatedAt,
	)
	if err != nil {
		return agenttrigger.TriggerLog{}, errx.Wrap(err, "create trigger log", errx.TypeInternal)
	}
	return log, nil
}

func (r *PostgresTriggerLogRepo) ListByTrigger(ctx context.Context, tenantID kernel.TenantID, triggerID kernel.AgentTriggerID, p kernel.PaginationOptions) (kernel.Paginated[agenttrigger.TriggerLog], error) {
	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM agent_trigger_logs WHERE trigger_id=$1 AND tenant_id=$2`,
		string(triggerID), string(tenantID),
	).Scan(&total); err != nil {
		return kernel.Paginated[agenttrigger.TriggerLog]{}, errx.Wrap(err, "count trigger logs", errx.TypeInternal)
	}

	var rows []dbTriggerLog
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, trigger_id, tenant_id, event_type, event_payload, agent_response, status, created_at
		FROM agent_trigger_logs
		WHERE trigger_id=$1 AND tenant_id=$2
		ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
		string(triggerID), string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[agenttrigger.TriggerLog]{}, errx.Wrap(err, "list trigger logs", errx.TypeInternal)
	}

	items := make([]agenttrigger.TriggerLog, len(rows))
	for i, row := range rows {
		items[i] = fromDBTriggerLog(row)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// Ensure interface compliance.
var _ agenttrigger.TriggerLogRepository = (*PostgresTriggerLogRepo)(nil)
