-- 042: Recommendations domain

-- Product views table — tracks when a visitor views a product page.
CREATE TABLE product_views (
    id          TEXT        NOT NULL,
    tenant_id   TEXT        NOT NULL,
    visitor_id  TEXT        NOT NULL,
    customer_id TEXT,
    product_id  TEXT        NOT NULL,
    source      TEXT,
    viewed_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_product_views_tenant  ON product_views(tenant_id);
CREATE INDEX idx_product_views_visitor ON product_views(tenant_id, visitor_id);
CREATE INDEX idx_product_views_product ON product_views(tenant_id, product_id);
CREATE INDEX idx_product_views_recent  ON product_views(tenant_id, viewed_at DESC);

-- Product interactions table — tracks richer behavioural signals.
CREATE TABLE product_interactions (
    id               TEXT        NOT NULL,
    tenant_id        TEXT        NOT NULL,
    visitor_id       TEXT        NOT NULL,
    customer_id      TEXT,
    product_id       TEXT        NOT NULL,
    interaction_type TEXT        NOT NULL,
    metadata         JSONB       NOT NULL DEFAULT '{}',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_product_interactions_tenant  ON product_interactions(tenant_id);
CREATE INDEX idx_product_interactions_product ON product_interactions(tenant_id, product_id, interaction_type);
CREATE INDEX idx_product_interactions_visitor ON product_interactions(tenant_id, visitor_id);

-- Recommendation rules table — admin-configurable rules for the recommendation engine.
CREATE TABLE recommendation_rules (
    id         TEXT        NOT NULL,
    tenant_id  TEXT        NOT NULL,
    name       TEXT        NOT NULL,
    type       TEXT        NOT NULL,
    config     JSONB       NOT NULL DEFAULT '{}',
    is_active  BOOLEAN     NOT NULL DEFAULT TRUE,
    priority   INT         NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_recommendation_rules_tenant ON recommendation_rules(tenant_id);
