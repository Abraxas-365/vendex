-- 037_multi_storefront.up.sql
-- Multi-storefront: allows a single tenant to run multiple storefronts.

CREATE TABLE storefronts (
    id               TEXT        NOT NULL,
    tenant_id        TEXT        NOT NULL,
    name             TEXT        NOT NULL,
    slug             TEXT        NOT NULL,
    domain           TEXT,
    description      TEXT        NOT NULL DEFAULT '',
    theme_id         TEXT        NOT NULL DEFAULT '',
    logo_url         TEXT        NOT NULL DEFAULT '',
    default_locale   TEXT        NOT NULL DEFAULT 'en',
    default_currency TEXT        NOT NULL DEFAULT 'USD',
    is_active        BOOLEAN     NOT NULL DEFAULT true,
    is_default       BOOLEAN     NOT NULL DEFAULT false,
    settings         JSONB       NOT NULL DEFAULT '{}',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    UNIQUE (tenant_id, slug),
    UNIQUE (tenant_id, domain)
);

CREATE INDEX idx_storefronts_tenant   ON storefronts(tenant_id);
CREATE INDEX idx_storefronts_domain   ON storefronts(domain) WHERE domain IS NOT NULL;

CREATE TABLE storefront_catalogs (
    id            TEXT        NOT NULL,
    tenant_id     TEXT        NOT NULL,
    storefront_id TEXT        NOT NULL REFERENCES storefronts(id) ON DELETE CASCADE,
    catalog_id    TEXT        NOT NULL,
    sort_order    INT         NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    UNIQUE (tenant_id, storefront_id, catalog_id)
);

CREATE INDEX idx_storefront_catalogs_tenant ON storefront_catalogs(tenant_id, storefront_id);
