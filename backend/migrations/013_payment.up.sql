-- Migration 013: Payment and Refund tables

CREATE TABLE payments (
    id                  VARCHAR(36)  NOT NULL,
    tenant_id           VARCHAR(36)  NOT NULL,
    order_id            VARCHAR(36)  NOT NULL,
    amount_amount       BIGINT       NOT NULL DEFAULT 0,
    amount_currency     CHAR(3)      NOT NULL DEFAULT 'USD',
    status              VARCHAR(50)  NOT NULL DEFAULT 'pending',
    provider            VARCHAR(50)  NOT NULL,
    provider_payment_id VARCHAR(255) NULL,
    provider_data       JSONB        NULL DEFAULT '{}',
    method              VARCHAR(50)  NULL,
    error_message       TEXT         NULL,
    paid_at             TIMESTAMPTZ  NULL,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_payments_order FOREIGN KEY (order_id) REFERENCES orders(id),
    CONSTRAINT chk_payment_status CHECK (status IN ('pending', 'processing', 'completed', 'failed', 'refunded'))
);

CREATE TABLE refunds (
    id                  VARCHAR(36)  NOT NULL,
    tenant_id           VARCHAR(36)  NOT NULL,
    payment_id          VARCHAR(36)  NOT NULL,
    order_id            VARCHAR(36)  NOT NULL,
    amount_amount       BIGINT       NOT NULL DEFAULT 0,
    amount_currency     CHAR(3)      NOT NULL DEFAULT 'USD',
    reason              TEXT         NULL,
    status              VARCHAR(50)  NOT NULL DEFAULT 'pending',
    provider_refund_id  VARCHAR(255) NULL,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_refunds_payment FOREIGN KEY (payment_id) REFERENCES payments(id),
    CONSTRAINT fk_refunds_order FOREIGN KEY (order_id) REFERENCES orders(id),
    CONSTRAINT chk_refund_status CHECK (status IN ('pending', 'completed', 'failed'))
);

CREATE INDEX idx_payments_tenant ON payments(tenant_id);
CREATE INDEX idx_payments_order ON payments(order_id);
CREATE INDEX idx_refunds_payment ON refunds(payment_id);
CREATE INDEX idx_refunds_order ON refunds(order_id);
CREATE INDEX idx_refunds_tenant ON refunds(tenant_id);
