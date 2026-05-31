CREATE TABLE experiments (
    id              TEXT        NOT NULL,
    tenant_id       TEXT        NOT NULL,
    name            TEXT        NOT NULL,
    description     TEXT        NOT NULL DEFAULT '',
    type            TEXT        NOT NULL DEFAULT 'page',
    status          TEXT        NOT NULL DEFAULT 'draft',
    traffic_percent INT         NOT NULL DEFAULT 100,
    started_at      TIMESTAMPTZ,
    ended_at        TIMESTAMPTZ,
    winner_variant_id TEXT,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_experiments_tenant ON experiments(tenant_id);
CREATE INDEX idx_experiments_status ON experiments(tenant_id, status);

CREATE TABLE experiment_variants (
    id            TEXT        NOT NULL,
    tenant_id     TEXT        NOT NULL,
    experiment_id TEXT        NOT NULL REFERENCES experiments(id) ON DELETE CASCADE,
    name          TEXT        NOT NULL,
    description   TEXT        NOT NULL DEFAULT '',
    weight        INT         NOT NULL DEFAULT 50,
    is_control    BOOLEAN     NOT NULL DEFAULT false,
    config        JSONB       NOT NULL DEFAULT '{}',
    visitors      INT         NOT NULL DEFAULT 0,
    conversions   INT         NOT NULL DEFAULT 0,
    revenue_cents BIGINT      NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_exp_variants_experiment ON experiment_variants(experiment_id);
CREATE INDEX idx_exp_variants_tenant ON experiment_variants(tenant_id);

CREATE TABLE experiment_assignments (
    id            TEXT        NOT NULL,
    tenant_id     TEXT        NOT NULL,
    experiment_id TEXT        NOT NULL REFERENCES experiments(id) ON DELETE CASCADE,
    variant_id    TEXT        NOT NULL REFERENCES experiment_variants(id) ON DELETE CASCADE,
    visitor_id    TEXT        NOT NULL,
    converted     BOOLEAN     NOT NULL DEFAULT false,
    revenue_cents BIGINT      NOT NULL DEFAULT 0,
    assigned_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    converted_at  TIMESTAMPTZ,
    PRIMARY KEY (id),
    UNIQUE (tenant_id, experiment_id, visitor_id)
);

CREATE INDEX idx_exp_assignments_visitor ON experiment_assignments(tenant_id, visitor_id);
CREATE INDEX idx_exp_assignments_variant ON experiment_assignments(variant_id);
