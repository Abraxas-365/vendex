-- Add SEO fields to products table
ALTER TABLE products ADD COLUMN meta_title VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE products ADD COLUMN meta_description TEXT NOT NULL DEFAULT '';
ALTER TABLE products ADD COLUMN slug VARCHAR(255) NOT NULL DEFAULT '';
ALTER TABLE products ADD COLUMN canonical_url TEXT NOT NULL DEFAULT '';

-- Backfill slugs from existing product names
UPDATE products SET slug = LOWER(REGEXP_REPLACE(REGEXP_REPLACE(TRIM(name), '[^a-zA-Z0-9]+', '-', 'g'), '^-+|-+$', '', 'g'))
WHERE slug = '';

-- Unique index: one slug per tenant
CREATE UNIQUE INDEX idx_products_tenant_slug ON products(tenant_id, slug);
