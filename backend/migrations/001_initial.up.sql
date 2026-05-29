-- =============================================================================
-- hada-commerce: Initial schema migration
-- All IDs are VARCHAR(36) UUIDs stored as strings.
-- Money fields are split into _amount (BIGINT, cents) + _currency (CHAR(3)).
-- JSONB is used for arrays and maps.
-- tenant_id is present on every table for multi-tenancy isolation.
-- =============================================================================

-- ---------------------------------------------------------------------------
-- 1. customers
-- ---------------------------------------------------------------------------
CREATE TABLE customers (
    id          VARCHAR(36)  NOT NULL,
    tenant_id   VARCHAR(36)  NOT NULL,
    email       VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    phone       VARCHAR(50)  NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT uq_customers_tenant_email UNIQUE (tenant_id, email)
);

CREATE INDEX idx_customers_tenant_id ON customers (tenant_id);
CREATE INDEX idx_customers_tenant_email ON customers (tenant_id, email);

-- ---------------------------------------------------------------------------
-- 2. customer_addresses
-- Stored as a child table (one row per address); mirrors customer.Address.
-- ---------------------------------------------------------------------------
CREATE TABLE customer_addresses (
    id          VARCHAR(36)  NOT NULL,
    customer_id VARCHAR(36)  NOT NULL,
    tenant_id   VARCHAR(36)  NOT NULL,
    street      VARCHAR(255) NOT NULL DEFAULT '',
    city        VARCHAR(100) NOT NULL DEFAULT '',
    state       VARCHAR(100) NOT NULL DEFAULT '',
    country     VARCHAR(100) NOT NULL DEFAULT '',
    postal_code VARCHAR(20)  NOT NULL DEFAULT '',
    is_default  BOOLEAN      NOT NULL DEFAULT FALSE,

    PRIMARY KEY (id),
    CONSTRAINT fk_customer_addresses_customer
        FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE CASCADE
);

CREATE INDEX idx_customer_addresses_customer_id ON customer_addresses (customer_id);
CREATE INDEX idx_customer_addresses_tenant_id   ON customer_addresses (tenant_id);

-- ---------------------------------------------------------------------------
-- 3. categories
-- Self-referencing via parent_id (nullable).
-- ---------------------------------------------------------------------------
CREATE TABLE categories (
    id          VARCHAR(36)  NOT NULL,
    tenant_id   VARCHAR(36)  NOT NULL,
    name        VARCHAR(255) NOT NULL,
    slug        VARCHAR(255) NOT NULL,
    parent_id   VARCHAR(36)  NULL,
    description TEXT         NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT uq_categories_tenant_slug UNIQUE (tenant_id, slug),
    CONSTRAINT fk_categories_parent
        FOREIGN KEY (parent_id) REFERENCES categories (id) ON DELETE SET NULL
);

CREATE INDEX idx_categories_tenant_id   ON categories (tenant_id);
CREATE INDEX idx_categories_tenant_slug ON categories (tenant_id, slug);
CREATE INDEX idx_categories_parent_id   ON categories (parent_id);

-- ---------------------------------------------------------------------------
-- 4. products
-- price → price_amount (BIGINT cents) + price_currency (CHAR(3)).
-- images and tags → JSONB arrays.
-- ---------------------------------------------------------------------------
CREATE TABLE products (
    id              VARCHAR(36)  NOT NULL,
    tenant_id       VARCHAR(36)  NOT NULL,
    name            VARCHAR(255) NOT NULL,
    description     TEXT         NOT NULL DEFAULT '',
    sku             VARCHAR(100) NOT NULL,
    price_amount    BIGINT       NOT NULL DEFAULT 0,
    price_currency  CHAR(3)      NOT NULL DEFAULT 'USD',
    images          JSONB        NOT NULL DEFAULT '[]',
    category_id     VARCHAR(36)  NULL,
    tags            JSONB        NOT NULL DEFAULT '[]',
    status          VARCHAR(20)  NOT NULL DEFAULT 'draft',
    stock           INTEGER      NOT NULL DEFAULT 0,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT uq_products_tenant_sku UNIQUE (tenant_id, sku),
    CONSTRAINT fk_products_category
        FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE SET NULL,
    CONSTRAINT chk_products_status
        CHECK (status IN ('draft', 'active', 'archived')),
    CONSTRAINT chk_products_stock
        CHECK (stock >= 0)
);

CREATE INDEX idx_products_tenant_id   ON products (tenant_id);
CREATE INDEX idx_products_tenant_sku  ON products (tenant_id, sku);
CREATE INDEX idx_products_category_id ON products (category_id);
CREATE INDEX idx_products_status      ON products (tenant_id, status);

-- ---------------------------------------------------------------------------
-- 5. collections
-- product_ids → JSONB array of ProductID strings.
-- rules       → JSONB object (map[string]any).
-- ---------------------------------------------------------------------------
CREATE TABLE collections (
    id           VARCHAR(36)  NOT NULL,
    tenant_id    VARCHAR(36)  NOT NULL,
    name         VARCHAR(255) NOT NULL,
    slug         VARCHAR(255) NOT NULL,
    description  TEXT         NOT NULL DEFAULT '',
    product_ids  JSONB        NOT NULL DEFAULT '[]',
    is_automatic BOOLEAN      NOT NULL DEFAULT FALSE,
    rules        JSONB        NOT NULL DEFAULT '{}',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT uq_collections_tenant_slug UNIQUE (tenant_id, slug)
);

CREATE INDEX idx_collections_tenant_id   ON collections (tenant_id);
CREATE INDEX idx_collections_tenant_slug ON collections (tenant_id, slug);

-- ---------------------------------------------------------------------------
-- 6. orders
-- total_amount → total_amount + total_currency.
-- shipping_address → JSONB (order.Address struct).
-- notes is added as a convenience column (not in entity, kept nullable).
-- ---------------------------------------------------------------------------
CREATE TABLE orders (
    id                  VARCHAR(36)  NOT NULL,
    tenant_id           VARCHAR(36)  NOT NULL,
    customer_id         VARCHAR(36)  NOT NULL,
    status              VARCHAR(20)  NOT NULL DEFAULT 'pending',
    total_amount        BIGINT       NOT NULL DEFAULT 0,
    total_currency      CHAR(3)      NOT NULL DEFAULT 'USD',
    shipping_address    JSONB        NOT NULL DEFAULT '{}',
    notes               TEXT         NULL,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT fk_orders_customer
        FOREIGN KEY (customer_id) REFERENCES customers (id) ON DELETE RESTRICT,
    CONSTRAINT chk_orders_status
        CHECK (status IN ('pending', 'confirmed', 'processing', 'shipped', 'delivered', 'cancelled'))
);

CREATE INDEX idx_orders_tenant_id   ON orders (tenant_id);
CREATE INDEX idx_orders_customer_id ON orders (customer_id);
CREATE INDEX idx_orders_status      ON orders (tenant_id, status);

-- ---------------------------------------------------------------------------
-- 7. order_items
-- Maps to order.OrderItem; unit_price and total are Money structs.
-- No tenant_id column — tenant is derived through orders.
-- ---------------------------------------------------------------------------
CREATE TABLE order_items (
    id                  VARCHAR(36)  NOT NULL,
    order_id            VARCHAR(36)  NOT NULL,
    product_id          VARCHAR(36)  NOT NULL,
    product_name        VARCHAR(255) NOT NULL,
    quantity            INTEGER      NOT NULL DEFAULT 1,
    unit_price_amount   BIGINT       NOT NULL DEFAULT 0,
    unit_price_currency CHAR(3)      NOT NULL DEFAULT 'USD',
    total_amount        BIGINT       NOT NULL DEFAULT 0,
    total_currency      CHAR(3)      NOT NULL DEFAULT 'USD',

    PRIMARY KEY (id),
    CONSTRAINT fk_order_items_order
        FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE,
    CONSTRAINT chk_order_items_quantity
        CHECK (quantity > 0)
);

CREATE INDEX idx_order_items_order_id   ON order_items (order_id);
CREATE INDEX idx_order_items_product_id ON order_items (product_id);

-- ---------------------------------------------------------------------------
-- 8. pages
-- meta → JSONB (storefront.PageMeta struct).
-- published_at is nullable (NULL until first publish).
-- ---------------------------------------------------------------------------
CREATE TABLE pages (
    id           VARCHAR(36)  NOT NULL,
    tenant_id    VARCHAR(36)  NOT NULL,
    slug         VARCHAR(255) NOT NULL,
    title        VARCHAR(500) NOT NULL,
    html         TEXT         NOT NULL DEFAULT '',
    css          TEXT         NOT NULL DEFAULT '',
    meta         JSONB        NOT NULL DEFAULT '{}',
    status       VARCHAR(20)  NOT NULL DEFAULT 'draft',
    version      INTEGER      NOT NULL DEFAULT 1,
    created_by   VARCHAR(255) NOT NULL DEFAULT '',
    published_at TIMESTAMPTZ  NULL,
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT uq_pages_tenant_slug UNIQUE (tenant_id, slug),
    CONSTRAINT chk_pages_status
        CHECK (status IN ('draft', 'pending_review', 'published', 'archived')),
    CONSTRAINT chk_pages_version
        CHECK (version >= 1)
);

CREATE INDEX idx_pages_tenant_id   ON pages (tenant_id);
CREATE INDEX idx_pages_tenant_slug ON pages (tenant_id, slug);
CREATE INDEX idx_pages_status      ON pages (tenant_id, status);

-- ---------------------------------------------------------------------------
-- 9. page_versions
-- Append-only history; never UPDATE or DELETE rows.
-- Maps to storefront.PageVersion.
-- ---------------------------------------------------------------------------
CREATE TABLE page_versions (
    id         VARCHAR(36)  NOT NULL,
    page_id    VARCHAR(36)  NOT NULL,
    tenant_id  VARCHAR(36)  NOT NULL,
    version    INTEGER      NOT NULL,
    html       TEXT         NOT NULL DEFAULT '',
    css        TEXT         NOT NULL DEFAULT '',
    edited_by  VARCHAR(255) NOT NULL DEFAULT '',
    comment    TEXT         NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT uq_page_versions_page_version UNIQUE (page_id, version),
    CONSTRAINT fk_page_versions_page
        FOREIGN KEY (page_id) REFERENCES pages (id) ON DELETE CASCADE
);

CREATE INDEX idx_page_versions_page_id   ON page_versions (page_id);
CREATE INDEX idx_page_versions_tenant_id ON page_versions (tenant_id);

-- ---------------------------------------------------------------------------
-- 10. promos
-- value is BIGINT (cents for fixed_amount; integer 0–100 for percentage).
-- min_order_amount and max_uses are nullable (nil = unlimited/no minimum).
-- starts_at / ends_at are nullable (nil = no bound).
-- ---------------------------------------------------------------------------
CREATE TABLE promos (
    id               VARCHAR(36)  NOT NULL,
    tenant_id        VARCHAR(36)  NOT NULL,
    code             VARCHAR(100) NOT NULL,
    type             VARCHAR(20)  NOT NULL,
    value            BIGINT       NOT NULL DEFAULT 0,
    min_order_amount BIGINT       NULL,
    max_uses         INTEGER      NULL,
    used_count       INTEGER      NOT NULL DEFAULT 0,
    starts_at        TIMESTAMPTZ  NULL,
    ends_at          TIMESTAMPTZ  NULL,
    active           BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT uq_promos_tenant_code UNIQUE (tenant_id, code),
    CONSTRAINT chk_promos_type
        CHECK (type IN ('percentage', 'fixed_amount', 'free_shipping')),
    CONSTRAINT chk_promos_used_count
        CHECK (used_count >= 0),
    CONSTRAINT chk_promos_value
        CHECK (value >= 0)
);

CREATE INDEX idx_promos_tenant_id   ON promos (tenant_id);
CREATE INDEX idx_promos_tenant_code ON promos (tenant_id, code);
CREATE INDEX idx_promos_active      ON promos (tenant_id, active);

-- ---------------------------------------------------------------------------
-- 11. media
-- size is BIGINT (bytes). url and alt are plain text.
-- ---------------------------------------------------------------------------
CREATE TABLE media (
    id           VARCHAR(36)  NOT NULL,
    tenant_id    VARCHAR(36)  NOT NULL,
    filename     VARCHAR(500) NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    size         BIGINT       NOT NULL DEFAULT 0,
    url          TEXT         NOT NULL,
    alt          TEXT         NOT NULL DEFAULT '',
    uploaded_by  VARCHAR(255) NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id)
);

CREATE INDEX idx_media_tenant_id ON media (tenant_id);
