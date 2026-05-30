-- Revert SEO fields from products table
DROP INDEX IF EXISTS idx_products_tenant_slug;
ALTER TABLE products DROP COLUMN IF EXISTS canonical_url;
ALTER TABLE products DROP COLUMN IF EXISTS slug;
ALTER TABLE products DROP COLUMN IF EXISTS meta_description;
ALTER TABLE products DROP COLUMN IF EXISTS meta_title;
