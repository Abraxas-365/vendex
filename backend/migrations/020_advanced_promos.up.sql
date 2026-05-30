-- Advanced promo features: buy-X-get-Y, product/category targeting, stackable promos

ALTER TABLE promos ADD COLUMN target_product_ids  JSONB   NULL DEFAULT '[]';
ALTER TABLE promos ADD COLUMN target_category_ids JSONB   NULL DEFAULT '[]';
ALTER TABLE promos ADD COLUMN customer_group_id   VARCHAR(36) NULL;
ALTER TABLE promos ADD COLUMN stackable           BOOLEAN NOT NULL DEFAULT FALSE;
ALTER TABLE promos ADD COLUMN buy_quantity        INTEGER NULL;
ALTER TABLE promos ADD COLUMN get_quantity        INTEGER NULL;
ALTER TABLE promos ADD COLUMN get_product_id      VARCHAR(36) NULL;
ALTER TABLE promos ADD COLUMN get_discount        BIGINT  NULL;

-- Extend the type constraint to include the new buy_x_get_y type
ALTER TABLE promos DROP CONSTRAINT IF EXISTS chk_promos_type;
ALTER TABLE promos ADD CONSTRAINT chk_promos_type
    CHECK (type IN ('percentage', 'fixed_amount', 'free_shipping', 'buy_x_get_y'));
