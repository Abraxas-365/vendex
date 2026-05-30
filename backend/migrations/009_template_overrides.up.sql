-- Per-tenant template overrides for storefront block rendering.
CREATE TABLE template_overrides (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID         NOT NULL REFERENCES tenants(id),
    block_type VARCHAR(100) NOT NULL,
    template   TEXT         NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (tenant_id, block_type)
);

CREATE INDEX idx_template_overrides_tenant ON template_overrides(tenant_id);
