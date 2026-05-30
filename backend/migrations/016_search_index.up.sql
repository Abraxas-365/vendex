-- Add tsvector column for fast full-text search on products.
ALTER TABLE products ADD COLUMN IF NOT EXISTS search_vector TSVECTOR;

-- Populate existing rows.
UPDATE products
SET search_vector = to_tsvector('english',
    coalesce(name, '') || ' ' ||
    coalesce(description, '') || ' ' ||
    coalesce(sku, '')
);

-- Create GIN index for efficient full-text queries.
CREATE INDEX IF NOT EXISTS idx_products_search ON products USING GIN (search_vector);

-- Function to keep search_vector up to date on INSERT/UPDATE.
CREATE OR REPLACE FUNCTION products_search_vector_update() RETURNS trigger AS $$
BEGIN
    NEW.search_vector := to_tsvector('english',
        coalesce(NEW.name, '') || ' ' ||
        coalesce(NEW.description, '') || ' ' ||
        coalesce(NEW.sku, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger that fires before every INSERT or UPDATE on products.
DROP TRIGGER IF EXISTS trg_products_search_vector ON products;
CREATE TRIGGER trg_products_search_vector
    BEFORE INSERT OR UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION products_search_vector_update();
