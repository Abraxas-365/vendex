-- Migration 037: Extend collections table and add collection_products
-- Adds richer metadata columns to the existing collections table
-- and creates the collection_products join table.

-- Add new columns to existing collections table (safe, backward-compatible)
ALTER TABLE collections
    ADD COLUMN IF NOT EXISTS image_url        TEXT,
    ADD COLUMN IF NOT EXISTS type             TEXT        NOT NULL DEFAULT 'manual',
    ADD COLUMN IF NOT EXISTS is_active        BOOLEAN     NOT NULL DEFAULT TRUE,
    ADD COLUMN IF NOT EXISTS sort_order       INT         NOT NULL DEFAULT 0,
    ADD COLUMN IF NOT EXISTS meta_title       TEXT,
    ADD COLUMN IF NOT EXISTS meta_description TEXT,
    ADD COLUMN IF NOT EXISTS published_at     TIMESTAMPTZ;

-- Rename is_automatic → type where is_automatic=true becomes 'auto'
-- (data migration — existing rows)
UPDATE collections SET type = 'auto' WHERE is_automatic = TRUE;

CREATE INDEX IF NOT EXISTS idx_collections_active ON collections(tenant_id, is_active);

-- collection_products: explicit product membership with sort order
CREATE TABLE IF NOT EXISTS collection_products (
    id            TEXT        NOT NULL,
    tenant_id     TEXT        NOT NULL,
    collection_id TEXT        NOT NULL REFERENCES collections(id) ON DELETE CASCADE,
    product_id    TEXT        NOT NULL,
    sort_order    INT         NOT NULL DEFAULT 0,
    added_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    UNIQUE (tenant_id, collection_id, product_id)
);

CREATE INDEX IF NOT EXISTS idx_collection_products_collection ON collection_products(collection_id);
CREATE INDEX IF NOT EXISTS idx_collection_products_product    ON collection_products(product_id);
CREATE INDEX IF NOT EXISTS idx_collection_products_tenant     ON collection_products(tenant_id);
