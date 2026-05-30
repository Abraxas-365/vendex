-- Subscriptions: recurring billing contracts

CREATE TABLE subscriptions (
    id              TEXT        NOT NULL,
    tenant_id       TEXT        NOT NULL,
    customer_id     TEXT        NOT NULL,
    product_id      TEXT        NOT NULL,
    variant_id      TEXT,
    price_amount    BIGINT      NOT NULL,
    price_currency  TEXT        NOT NULL DEFAULT 'USD',
    interval        TEXT        NOT NULL,
    status          TEXT        NOT NULL DEFAULT 'active',
    next_billing_date TIMESTAMPTZ NOT NULL,
    last_billed_at  TIMESTAMPTZ,
    cancelled_at    TIMESTAMPTZ,
    paused_at       TIMESTAMPTZ,
    trial_ends_at   TIMESTAMPTZ,
    metadata        JSONB       NOT NULL DEFAULT '{}',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_subscriptions_tenant    ON subscriptions(tenant_id);
CREATE INDEX idx_subscriptions_customer  ON subscriptions(tenant_id, customer_id);
CREATE INDEX idx_subscriptions_billing   ON subscriptions(tenant_id, status, next_billing_date);

-- Billing records: one row per billing attempt

CREATE TABLE billing_records (
    id               TEXT        NOT NULL,
    subscription_id  TEXT        NOT NULL REFERENCES subscriptions(id) ON DELETE CASCADE,
    tenant_id        TEXT        NOT NULL,
    amount_cents     BIGINT      NOT NULL,
    amount_currency  TEXT        NOT NULL DEFAULT 'USD',
    status           TEXT        NOT NULL DEFAULT 'pending',
    order_id         TEXT,
    failure_reason   TEXT,
    billed_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_billing_records_sub    ON billing_records(subscription_id);
CREATE INDEX idx_billing_records_tenant ON billing_records(tenant_id);
