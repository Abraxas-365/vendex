-- 022: Gift Cards domain

CREATE TABLE gift_cards (
    id                      TEXT        NOT NULL,
    tenant_id               TEXT        NOT NULL,
    code                    TEXT        NOT NULL,
    initial_amount_cents    BIGINT      NOT NULL,
    initial_amount_currency TEXT        NOT NULL DEFAULT 'USD',
    balance_cents           BIGINT      NOT NULL,
    balance_currency        TEXT        NOT NULL DEFAULT 'USD',
    expires_at              TIMESTAMPTZ,
    active                  BOOLEAN     NOT NULL DEFAULT true,
    created_by              TEXT        NOT NULL DEFAULT '',
    created_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    UNIQUE (tenant_id, code)
);

CREATE INDEX idx_gift_cards_tenant ON gift_cards(tenant_id);

CREATE TABLE gift_card_transactions (
    id               TEXT        NOT NULL,
    gift_card_id     TEXT        NOT NULL REFERENCES gift_cards(id) ON DELETE CASCADE,
    tenant_id        TEXT        NOT NULL,
    type             TEXT        NOT NULL, -- 'credit' or 'debit'
    amount_cents     BIGINT      NOT NULL,
    amount_currency  TEXT        NOT NULL DEFAULT 'USD',
    order_id         TEXT,
    note             TEXT        NOT NULL DEFAULT '',
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_gc_transactions_card   ON gift_card_transactions(gift_card_id);
CREATE INDEX idx_gc_transactions_tenant ON gift_card_transactions(tenant_id);
