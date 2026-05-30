-- Navigation menus for storefront header and footer layout.
CREATE TABLE navigation_menus (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id  UUID         NOT NULL REFERENCES tenants(id),
    location   VARCHAR(20)  NOT NULL CHECK (location IN ('header', 'footer')),
    label      VARCHAR(255) NOT NULL,
    url        VARCHAR(500) NOT NULL,
    position   INTEGER      NOT NULL DEFAULT 0,
    parent_id  UUID         REFERENCES navigation_menus(id),
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_navigation_menus_tenant ON navigation_menus(tenant_id);
CREATE INDEX idx_navigation_menus_location ON navigation_menus(tenant_id, location);
