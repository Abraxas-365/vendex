-- =============================================================================
-- hada-commerce: Marketplace & Plugin Registry migration
-- IDs are VARCHAR(36) UUIDs stored as strings (consistent with 001_initial).
-- JSONB is used for arrays and maps.
-- plugin_installations is scoped by tenant_id for multi-tenancy isolation.
-- =============================================================================

-- ---------------------------------------------------------------------------
-- 1. plugins — the global plugin catalogue (one row per plugin)
-- ---------------------------------------------------------------------------
CREATE TABLE plugins (
    id           VARCHAR(36)  NOT NULL,
    name         TEXT         NOT NULL UNIQUE,
    display_name TEXT         NOT NULL,
    description  TEXT         NOT NULL DEFAULT '',
    author       TEXT         NOT NULL DEFAULT '',
    icon         TEXT         NOT NULL DEFAULT '',
    category     TEXT         NOT NULL DEFAULT 'community'
                              CHECK (category IN ('official', 'community', 'custom')),
    tags         JSONB        NOT NULL DEFAULT '[]',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

CREATE INDEX idx_plugins_category ON plugins (category);
CREATE INDEX idx_plugins_name     ON plugins (name);

-- ---------------------------------------------------------------------------
-- 2. plugin_versions — one row per published version of a plugin
-- ---------------------------------------------------------------------------
CREATE TABLE plugin_versions (
    id               VARCHAR(36)  NOT NULL,
    plugin_id        VARCHAR(36)  NOT NULL REFERENCES plugins (id) ON DELETE CASCADE,
    version          TEXT         NOT NULL,
    changelog        TEXT         NOT NULL DEFAULT '',
    permissions      JSONB        NOT NULL DEFAULT '[]',
    manifest_json    TEXT         NOT NULL DEFAULT '{}',
    frontend_url     TEXT         NOT NULL DEFAULT '',
    backend_entry    TEXT         NOT NULL DEFAULT '',
    min_platform_ver TEXT         NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT uq_plugin_versions_plugin_version UNIQUE (plugin_id, version)
);

CREATE INDEX idx_plugin_versions_plugin ON plugin_versions (plugin_id);

-- ---------------------------------------------------------------------------
-- 3. plugin_installations — per-tenant plugin activation record
-- ---------------------------------------------------------------------------
CREATE TABLE plugin_installations (
    id           VARCHAR(36)  NOT NULL,
    tenant_id    VARCHAR(36)  NOT NULL,
    plugin_id    VARCHAR(36)  NOT NULL REFERENCES plugins (id) ON DELETE CASCADE,
    version_id   VARCHAR(36)  NOT NULL REFERENCES plugin_versions (id),
    status       TEXT         NOT NULL DEFAULT 'active'
                              CHECK (status IN ('active', 'inactive', 'failed')),
    settings     JSONB        NOT NULL DEFAULT '{}',
    installed_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT uq_installations_tenant_plugin UNIQUE (tenant_id, plugin_id)
);

CREATE INDEX idx_installations_tenant        ON plugin_installations (tenant_id);
CREATE INDEX idx_installations_tenant_status ON plugin_installations (tenant_id, status);
