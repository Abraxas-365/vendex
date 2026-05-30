CREATE TABLE IF NOT EXISTS warehouses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    name TEXT NOT NULL,
    address TEXT,
    is_default BOOLEAN DEFAULT false,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_warehouses_tenant ON warehouses(tenant_id);

CREATE TABLE IF NOT EXISTS stock_levels (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    quantity INT NOT NULL DEFAULT 0,
    reserved INT NOT NULL DEFAULT 0,
    low_stock_threshold INT NOT NULL DEFAULT 5,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, product_id, variant_id, warehouse_id)
);
CREATE INDEX idx_stock_levels_tenant ON stock_levels(tenant_id);
CREATE INDEX idx_stock_levels_product ON stock_levels(tenant_id, product_id);

CREATE TABLE IF NOT EXISTS stock_movements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    product_id UUID NOT NULL,
    variant_id UUID,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    type TEXT NOT NULL, -- 'received', 'sold', 'returned', 'adjusted', 'transferred'
    quantity INT NOT NULL,
    reference TEXT,
    note TEXT,
    created_by TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_stock_movements_tenant ON stock_movements(tenant_id);
