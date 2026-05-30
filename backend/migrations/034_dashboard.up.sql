-- Dashboard domain: reads from existing tables (orders, customers, products)
-- No new tables needed, but adding useful indexes for reporting queries.

CREATE INDEX IF NOT EXISTS idx_orders_tenant_created ON orders(tenant_id, created_at);
CREATE INDEX IF NOT EXISTS idx_orders_tenant_status  ON orders(tenant_id, status);
