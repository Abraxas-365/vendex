package abtestinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/Abraxas-365/vendex/internal/abtest"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements abtest.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed A/B test repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// compile-time check
var _ abtest.Repository = (*PostgresRepo)(nil)

// ─── Experiment ──────────────────────────────────────────────────────────────

func (r *PostgresRepo) CreateExperiment(ctx context.Context, e *abtest.Experiment) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO experiments
			(id, tenant_id, name, description, type, status, traffic_percent, started_at, ended_at, winner_variant_id, created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		string(e.ID), string(e.TenantID), e.Name, e.Description,
		string(e.Type), string(e.Status), e.TrafficPercent,
		e.StartedAt, e.EndedAt, nullableVariantID(e.WinnerVariantID),
		e.CreatedAt, e.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting experiment", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetExperimentByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID) (*abtest.Experiment, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, description, type, status, traffic_percent,
		       started_at, ended_at, winner_variant_id, created_at, updated_at
		FROM experiments
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	e, err := scanExperiment(row)
	if err == sql.ErrNoRows {
		return nil, abtest.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning experiment", errx.TypeInternal)
	}
	return e, nil
}

func (r *PostgresRepo) UpdateExperiment(ctx context.Context, e *abtest.Experiment) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE experiments SET
			name             = $3,
			description      = $4,
			type             = $5,
			status           = $6,
			traffic_percent  = $7,
			started_at       = $8,
			ended_at         = $9,
			winner_variant_id = $10,
			updated_at       = $11
		WHERE id = $1 AND tenant_id = $2`,
		string(e.ID), string(e.TenantID),
		e.Name, e.Description, string(e.Type), string(e.Status),
		e.TrafficPercent, e.StartedAt, e.EndedAt,
		nullableVariantID(e.WinnerVariantID), e.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "updating experiment", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) DeleteExperiment(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM experiments WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting experiment", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) ListExperiments(ctx context.Context, tenantID kernel.TenantID, status string, pg kernel.PaginationOptions) (kernel.Paginated[abtest.Experiment], error) {
	var total int
	if status != "" {
		if err := r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM experiments WHERE tenant_id = $1 AND status = $2`,
			string(tenantID), status,
		).Scan(&total); err != nil {
			return kernel.Paginated[abtest.Experiment]{}, errx.Wrap(err, "counting experiments", errx.TypeInternal)
		}
	} else {
		if err := r.db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM experiments WHERE tenant_id = $1`,
			string(tenantID),
		).Scan(&total); err != nil {
			return kernel.Paginated[abtest.Experiment]{}, errx.Wrap(err, "counting experiments", errx.TypeInternal)
		}
	}

	var rows *sql.Rows
	var err error
	if status != "" {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, tenant_id, name, description, type, status, traffic_percent,
			       started_at, ended_at, winner_variant_id, created_at, updated_at
			FROM experiments
			WHERE tenant_id = $1 AND status = $2
			ORDER BY created_at DESC
			LIMIT $3 OFFSET $4`,
			string(tenantID), status, pg.Limit(), pg.Offset(),
		)
	} else {
		rows, err = r.db.QueryContext(ctx, `
			SELECT id, tenant_id, name, description, type, status, traffic_percent,
			       started_at, ended_at, winner_variant_id, created_at, updated_at
			FROM experiments
			WHERE tenant_id = $1
			ORDER BY created_at DESC
			LIMIT $2 OFFSET $3`,
			string(tenantID), pg.Limit(), pg.Offset(),
		)
	}
	if err != nil {
		return kernel.Paginated[abtest.Experiment]{}, errx.Wrap(err, "listing experiments", errx.TypeInternal)
	}
	defer rows.Close()

	items := make([]abtest.Experiment, 0)
	for rows.Next() {
		e, err := scanExperimentRow(rows)
		if err != nil {
			return kernel.Paginated[abtest.Experiment]{}, errx.Wrap(err, "scanning experiment row", errx.TypeInternal)
		}
		items = append(items, *e)
	}

	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

// ─── Variant ─────────────────────────────────────────────────────────────────

func (r *PostgresRepo) CreateVariant(ctx context.Context, v *abtest.ExperimentVariant) error {
	configJSON, err := json.Marshal(v.Config)
	if err != nil {
		return errx.Wrap(err, "marshaling config", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO experiment_variants
			(id, tenant_id, experiment_id, name, description, weight, is_control, config, visitors, conversions, revenue_cents, created_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		string(v.ID), string(v.TenantID), string(v.ExperimentID),
		v.Name, v.Description, v.Weight, v.IsControl,
		string(configJSON), v.Visitors, v.Conversions, v.RevenueCents,
		v.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting variant", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetVariantByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentVariantID) (*abtest.ExperimentVariant, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, experiment_id, name, description, weight, is_control, config,
		       visitors, conversions, revenue_cents, created_at
		FROM experiment_variants
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	v, err := scanVariant(row)
	if err == sql.ErrNoRows {
		return nil, abtest.ErrVariantNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning variant", errx.TypeInternal)
	}
	return v, nil
}

func (r *PostgresRepo) ListVariants(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID) ([]abtest.ExperimentVariant, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, experiment_id, name, description, weight, is_control, config,
		       visitors, conversions, revenue_cents, created_at
		FROM experiment_variants
		WHERE experiment_id = $1 AND tenant_id = $2
		ORDER BY created_at ASC`,
		string(experimentID), string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "listing variants", errx.TypeInternal)
	}
	defer rows.Close()

	var items []abtest.ExperimentVariant
	for rows.Next() {
		v, err := scanVariantRow(rows)
		if err != nil {
			return nil, errx.Wrap(err, "scanning variant row", errx.TypeInternal)
		}
		items = append(items, *v)
	}
	return items, nil
}

func (r *PostgresRepo) DeleteVariant(ctx context.Context, tenantID kernel.TenantID, id kernel.ExperimentVariantID) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM experiment_variants WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting variant", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) IncrementVariantVisitors(ctx context.Context, variantID kernel.ExperimentVariantID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE experiment_variants SET visitors = visitors + 1 WHERE id = $1`,
		string(variantID),
	)
	if err != nil {
		return errx.Wrap(err, "incrementing visitors", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) IncrementVariantConversions(ctx context.Context, variantID kernel.ExperimentVariantID, revenueCents int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE experiment_variants SET conversions = conversions + 1, revenue_cents = revenue_cents + $2 WHERE id = $1`,
		string(variantID), revenueCents,
	)
	if err != nil {
		return errx.Wrap(err, "incrementing conversions", errx.TypeInternal)
	}
	return nil
}

// ─── Assignment ──────────────────────────────────────────────────────────────

func (r *PostgresRepo) CreateAssignment(ctx context.Context, a *abtest.ExperimentAssignment) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO experiment_assignments
			(id, tenant_id, experiment_id, variant_id, visitor_id, converted, revenue_cents, assigned_at, converted_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (tenant_id, experiment_id, visitor_id) DO NOTHING`,
		string(a.ID), string(a.TenantID), string(a.ExperimentID),
		string(a.VariantID), a.VisitorID, a.Converted,
		a.RevenueCents, a.AssignedAt, a.ConvertedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting assignment", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetAssignment(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID, visitorID string) (*abtest.ExperimentAssignment, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, experiment_id, variant_id, visitor_id, converted, revenue_cents, assigned_at, converted_at
		FROM experiment_assignments
		WHERE tenant_id = $1 AND experiment_id = $2 AND visitor_id = $3`,
		string(tenantID), string(experimentID), visitorID,
	)
	a, err := scanAssignment(row)
	if err == sql.ErrNoRows {
		return nil, abtest.ErrAssignmentNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning assignment", errx.TypeInternal)
	}
	return a, nil
}

func (r *PostgresRepo) RecordConversion(ctx context.Context, tenantID kernel.TenantID, experimentID kernel.ExperimentID, visitorID string, revenueCents int64) error {
	now := time.Now()
	_, err := r.db.ExecContext(ctx, `
		UPDATE experiment_assignments
		SET converted = true, revenue_cents = $4, converted_at = $5
		WHERE tenant_id = $1 AND experiment_id = $2 AND visitor_id = $3 AND converted = false`,
		string(tenantID), string(experimentID), visitorID, revenueCents, now,
	)
	if err != nil {
		return errx.Wrap(err, "recording conversion", errx.TypeInternal)
	}
	return nil
}

// ─── Scan helpers ────────────────────────────────────────────────────────────

type experimentScanner interface {
	Scan(dest ...interface{}) error
}

func scanExperiment(s experimentScanner) (*abtest.Experiment, error) {
	var e abtest.Experiment
	var id, tenantID, expType, status string
	var winnerID sql.NullString

	err := s.Scan(
		&id, &tenantID, &e.Name, &e.Description,
		&expType, &status, &e.TrafficPercent,
		&e.StartedAt, &e.EndedAt, &winnerID,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	e.ID = kernel.ExperimentID(id)
	e.TenantID = kernel.TenantID(tenantID)
	e.Type = abtest.ExperimentType(expType)
	e.Status = abtest.ExperimentStatus(status)
	if winnerID.Valid {
		vid := kernel.ExperimentVariantID(winnerID.String)
		e.WinnerVariantID = &vid
	}

	return &e, nil
}

func scanExperimentRow(rows *sql.Rows) (*abtest.Experiment, error) {
	var e abtest.Experiment
	var id, tenantID, expType, status string
	var winnerID sql.NullString

	err := rows.Scan(
		&id, &tenantID, &e.Name, &e.Description,
		&expType, &status, &e.TrafficPercent,
		&e.StartedAt, &e.EndedAt, &winnerID,
		&e.CreatedAt, &e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	e.ID = kernel.ExperimentID(id)
	e.TenantID = kernel.TenantID(tenantID)
	e.Type = abtest.ExperimentType(expType)
	e.Status = abtest.ExperimentStatus(status)
	if winnerID.Valid {
		vid := kernel.ExperimentVariantID(winnerID.String)
		e.WinnerVariantID = &vid
	}

	return &e, nil
}

type variantScanner interface {
	Scan(dest ...interface{}) error
}

func scanVariant(s variantScanner) (*abtest.ExperimentVariant, error) {
	var v abtest.ExperimentVariant
	var id, tenantID, experimentID string
	var configJSON string

	err := s.Scan(
		&id, &tenantID, &experimentID,
		&v.Name, &v.Description, &v.Weight, &v.IsControl,
		&configJSON, &v.Visitors, &v.Conversions, &v.RevenueCents,
		&v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	v.ID = kernel.ExperimentVariantID(id)
	v.TenantID = kernel.TenantID(tenantID)
	v.ExperimentID = kernel.ExperimentID(experimentID)

	if configJSON != "" {
		if err := json.Unmarshal([]byte(configJSON), &v.Config); err != nil {
			v.Config = map[string]interface{}{}
		}
	} else {
		v.Config = map[string]interface{}{}
	}

	return &v, nil
}

func scanVariantRow(rows *sql.Rows) (*abtest.ExperimentVariant, error) {
	var v abtest.ExperimentVariant
	var id, tenantID, experimentID string
	var configJSON string

	err := rows.Scan(
		&id, &tenantID, &experimentID,
		&v.Name, &v.Description, &v.Weight, &v.IsControl,
		&configJSON, &v.Visitors, &v.Conversions, &v.RevenueCents,
		&v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	v.ID = kernel.ExperimentVariantID(id)
	v.TenantID = kernel.TenantID(tenantID)
	v.ExperimentID = kernel.ExperimentID(experimentID)

	if configJSON != "" {
		if err := json.Unmarshal([]byte(configJSON), &v.Config); err != nil {
			v.Config = map[string]interface{}{}
		}
	} else {
		v.Config = map[string]interface{}{}
	}

	return &v, nil
}

type assignmentScanner interface {
	Scan(dest ...interface{}) error
}

func scanAssignment(s assignmentScanner) (*abtest.ExperimentAssignment, error) {
	var a abtest.ExperimentAssignment
	var id, tenantID, experimentID, variantID string

	err := s.Scan(
		&id, &tenantID, &experimentID, &variantID,
		&a.VisitorID, &a.Converted, &a.RevenueCents,
		&a.AssignedAt, &a.ConvertedAt,
	)
	if err != nil {
		return nil, err
	}

	a.ID = kernel.ExperimentAssignmentID(id)
	a.TenantID = kernel.TenantID(tenantID)
	a.ExperimentID = kernel.ExperimentID(experimentID)
	a.VariantID = kernel.ExperimentVariantID(variantID)

	return &a, nil
}

// nullableVariantID converts an optional ExperimentVariantID to a sql-compatible value.
func nullableVariantID(id *kernel.ExperimentVariantID) interface{} {
	if id == nil {
		return nil
	}
	return string(*id)
}
