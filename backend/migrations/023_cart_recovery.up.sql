CREATE TABLE cart_recovery_emails (
    id           UUID        NOT NULL DEFAULT gen_random_uuid(),
    tenant_id    TEXT        NOT NULL,
    cart_id      TEXT        NOT NULL,
    customer_id  TEXT        NOT NULL,
    email        TEXT        NOT NULL,
    step         INT         NOT NULL DEFAULT 1,
    status       TEXT        NOT NULL DEFAULT 'pending',
    discount_code TEXT,
    sent_at      TIMESTAMPTZ,
    clicked_at   TIMESTAMPTZ,
    converted_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_cart_recovery_tenant_status ON cart_recovery_emails(tenant_id, status);
CREATE INDEX idx_cart_recovery_cart          ON cart_recovery_emails(tenant_id, cart_id);
