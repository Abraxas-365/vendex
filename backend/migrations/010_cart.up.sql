CREATE TABLE carts (
    id          VARCHAR(36)  NOT NULL,
    tenant_id   VARCHAR(36)  NOT NULL,
    customer_id VARCHAR(36)  NULL,
    session_id  VARCHAR(255) NULL,
    currency    CHAR(3)      NOT NULL DEFAULT 'USD',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ  NOT NULL DEFAULT (NOW() + INTERVAL '30 days'),
    PRIMARY KEY (id)
);

CREATE TABLE cart_items (
    id                  VARCHAR(36)  NOT NULL,
    cart_id             VARCHAR(36)  NOT NULL,
    tenant_id           VARCHAR(36)  NOT NULL,
    product_id          VARCHAR(36)  NOT NULL,
    variant_id          VARCHAR(36)  NULL,
    quantity            INTEGER      NOT NULL DEFAULT 1,
    unit_price_amount   BIGINT       NOT NULL DEFAULT 0,
    unit_price_currency CHAR(3)      NOT NULL DEFAULT 'USD',
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_cart_items_cart FOREIGN KEY (cart_id) REFERENCES carts(id) ON DELETE CASCADE,
    CONSTRAINT chk_cart_items_quantity CHECK (quantity > 0)
);

CREATE INDEX idx_carts_tenant ON carts(tenant_id);
CREATE INDEX idx_carts_session ON carts(tenant_id, session_id);
CREATE INDEX idx_carts_customer ON carts(tenant_id, customer_id);
CREATE INDEX idx_cart_items_cart ON cart_items(cart_id);
