-- Add has_variants flag to products
ALTER TABLE products ADD COLUMN has_variants BOOLEAN NOT NULL DEFAULT FALSE;

-- Product options (e.g. "Size", "Color") with predefined values stored as JSONB
CREATE TABLE product_options (
    id          VARCHAR(36)  NOT NULL,
    product_id  VARCHAR(36)  NOT NULL,
    tenant_id   VARCHAR(36)  NOT NULL,
    name        VARCHAR(255) NOT NULL,
    position    INTEGER      NOT NULL DEFAULT 0,
    values      JSONB        NOT NULL DEFAULT '[]',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_product_options_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

-- Product variants (specific option combinations with own price/SKU/stock)
CREATE TABLE product_variants (
    id              VARCHAR(36)  NOT NULL,
    product_id      VARCHAR(36)  NOT NULL,
    tenant_id       VARCHAR(36)  NOT NULL,
    sku             VARCHAR(255) NOT NULL,
    price_amount    BIGINT       NOT NULL DEFAULT 0,
    price_currency  CHAR(3)      NOT NULL DEFAULT 'USD',
    stock           INTEGER      NOT NULL DEFAULT 0,
    options         JSONB        NOT NULL DEFAULT '{}',
    active          BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_product_variants_product FOREIGN KEY (product_id) REFERENCES products(id) ON DELETE CASCADE
);

CREATE INDEX idx_product_options_product ON product_options(product_id);
CREATE INDEX idx_product_options_tenant  ON product_options(tenant_id);
CREATE INDEX idx_product_variants_product ON product_variants(product_id);
CREATE INDEX idx_product_variants_tenant  ON product_variants(tenant_id);
CREATE INDEX idx_product_variants_sku     ON product_variants(tenant_id, sku);
