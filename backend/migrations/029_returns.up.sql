CREATE TABLE IF NOT EXISTS return_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    order_id UUID NOT NULL,
    customer_id UUID NOT NULL,
    status TEXT NOT NULL DEFAULT 'requested', -- 'requested','approved','rejected','received','refunded','exchanged','closed'
    reason TEXT NOT NULL,
    notes TEXT,
    admin_notes TEXT,
    refund_amount_cents BIGINT DEFAULT 0,
    refund_currency TEXT DEFAULT 'USD',
    resolution TEXT, -- 'refund', 'exchange', 'store_credit'
    tracking_number TEXT,
    carrier TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_return_requests_tenant ON return_requests(tenant_id);
CREATE INDEX idx_return_requests_order ON return_requests(tenant_id, order_id);

CREATE TABLE IF NOT EXISTS return_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    return_id UUID NOT NULL REFERENCES return_requests(id) ON DELETE CASCADE,
    tenant_id TEXT NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID,
    quantity INT NOT NULL DEFAULT 1,
    reason TEXT,
    condition TEXT DEFAULT 'unopened', -- 'unopened', 'like_new', 'used', 'damaged'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_return_items_return ON return_items(return_id);
