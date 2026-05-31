-- Rollback migration 037
DROP TABLE IF EXISTS collection_products;

DROP INDEX IF EXISTS idx_collections_active;

ALTER TABLE collections
    DROP COLUMN IF EXISTS published_at,
    DROP COLUMN IF EXISTS meta_description,
    DROP COLUMN IF EXISTS meta_title,
    DROP COLUMN IF EXISTS sort_order,
    DROP COLUMN IF EXISTS is_active,
    DROP COLUMN IF EXISTS type,
    DROP COLUMN IF EXISTS image_url;
