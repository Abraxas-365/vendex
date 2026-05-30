DROP INDEX IF EXISTS idx_product_variants_sku;
DROP INDEX IF EXISTS idx_product_variants_tenant;
DROP INDEX IF EXISTS idx_product_variants_product;
DROP INDEX IF EXISTS idx_product_options_tenant;
DROP INDEX IF EXISTS idx_product_options_product;

DROP TABLE IF EXISTS product_variants;
DROP TABLE IF EXISTS product_options;

ALTER TABLE products DROP COLUMN IF EXISTS has_variants;
