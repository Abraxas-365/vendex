CREATE TABLE IF NOT EXISTS reviews (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    product_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    rating INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
    title TEXT,
    body TEXT,
    status TEXT NOT NULL DEFAULT 'pending', -- 'pending', 'approved', 'rejected'
    verified_purchase BOOLEAN DEFAULT false,
    helpful_count INT DEFAULT 0,
    images TEXT[], -- array of image URLs
    admin_response TEXT,
    admin_responded_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_reviews_tenant ON reviews(tenant_id);
CREATE INDEX idx_reviews_product ON reviews(tenant_id, product_id);
CREATE INDEX idx_reviews_customer ON reviews(tenant_id, customer_id);
CREATE INDEX idx_reviews_status ON reviews(tenant_id, status);
