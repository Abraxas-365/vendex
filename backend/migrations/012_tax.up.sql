-- Tax rates table: stores tax rate configurations by jurisdiction.
CREATE TABLE tax_rates (
    id                  VARCHAR(36)   NOT NULL,
    tenant_id           VARCHAR(36)   NOT NULL,
    name                VARCHAR(255)  NOT NULL,
    rate                DECIMAL(8,6)  NOT NULL,
    country             VARCHAR(2)    NOT NULL,
    state               VARCHAR(100)  NULL DEFAULT '',
    city                VARCHAR(255)  NULL DEFAULT '',
    zip_code            VARCHAR(20)   NULL DEFAULT '',
    priority            INTEGER       NOT NULL DEFAULT 0,
    compound            BOOLEAN       NOT NULL DEFAULT FALSE,
    includes_shipping   BOOLEAN       NOT NULL DEFAULT FALSE,
    active              BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT chk_tax_rate_positive CHECK (rate >= 0 AND rate <= 1)
);

CREATE INDEX idx_tax_rates_tenant ON tax_rates(tenant_id);
CREATE INDEX idx_tax_rates_location ON tax_rates(tenant_id, country, state);
