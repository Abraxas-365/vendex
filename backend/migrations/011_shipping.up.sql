CREATE TABLE shipping_zones (
    id          VARCHAR(36)  NOT NULL,
    tenant_id   VARCHAR(36)  NOT NULL,
    name        VARCHAR(255) NOT NULL,
    countries   JSONB        NOT NULL DEFAULT '[]',
    states      JSONB        NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE shipping_rates (
    id                  VARCHAR(36)   NOT NULL,
    zone_id             VARCHAR(36)   NOT NULL,
    tenant_id           VARCHAR(36)   NOT NULL,
    name                VARCHAR(255)  NOT NULL,
    type                VARCHAR(50)   NOT NULL,
    price_amount        BIGINT        NOT NULL DEFAULT 0,
    price_currency      CHAR(3)       NOT NULL DEFAULT 'USD',
    min_weight          DECIMAL(10,2) NULL,
    max_weight          DECIMAL(10,2) NULL,
    min_order_amount    BIGINT        NULL,
    max_order_amount    BIGINT        NULL,
    est_days_min        INTEGER       NULL,
    est_days_max        INTEGER       NULL,
    active              BOOLEAN       NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ   NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_shipping_rates_zone FOREIGN KEY (zone_id) REFERENCES shipping_zones(id) ON DELETE CASCADE,
    CONSTRAINT chk_shipping_rate_type CHECK (type IN ('flat', 'weight_based', 'price_based', 'free'))
);

CREATE INDEX idx_shipping_zones_tenant ON shipping_zones(tenant_id);
CREATE INDEX idx_shipping_rates_zone ON shipping_rates(zone_id);
CREATE INDEX idx_shipping_rates_tenant ON shipping_rates(tenant_id);
