-- Multi-currency exchange rate management
CREATE TABLE currency_rates (
    id              TEXT            NOT NULL,
    tenant_id       TEXT            NOT NULL,
    base_currency   TEXT            NOT NULL,
    target_currency TEXT            NOT NULL,
    rate            NUMERIC(20, 10) NOT NULL,
    auto_update     BOOLEAN         NOT NULL DEFAULT false,
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE (tenant_id, base_currency, target_currency)
);

CREATE INDEX idx_currency_rates_tenant ON currency_rates(tenant_id);
