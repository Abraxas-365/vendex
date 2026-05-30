-- Rollback advanced promo features

-- Restore original type constraint (remove buy_x_get_y)
ALTER TABLE promos DROP CONSTRAINT IF EXISTS chk_promos_type;
ALTER TABLE promos ADD CONSTRAINT chk_promos_type
    CHECK (type IN ('percentage', 'fixed_amount', 'free_shipping'));

ALTER TABLE promos DROP COLUMN IF EXISTS get_discount;
ALTER TABLE promos DROP COLUMN IF EXISTS get_product_id;
ALTER TABLE promos DROP COLUMN IF EXISTS get_quantity;
ALTER TABLE promos DROP COLUMN IF EXISTS buy_quantity;
ALTER TABLE promos DROP COLUMN IF EXISTS stackable;
ALTER TABLE promos DROP COLUMN IF EXISTS customer_group_id;
ALTER TABLE promos DROP COLUMN IF EXISTS target_category_ids;
ALTER TABLE promos DROP COLUMN IF EXISTS target_product_ids;
