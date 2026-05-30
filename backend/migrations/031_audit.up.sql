CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    user_email TEXT,
    action TEXT NOT NULL, -- 'create', 'update', 'delete', 'login', 'export', etc.
    resource_type TEXT NOT NULL, -- 'product', 'order', 'customer', 'settings', etc.
    resource_id TEXT,
    changes JSONB, -- {field: {old: x, new: y}}
    metadata JSONB, -- extra context (IP, user agent, etc.)
    ip_address TEXT,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_logs_tenant ON audit_logs(tenant_id);
CREATE INDEX idx_audit_logs_user ON audit_logs(tenant_id, user_id);
CREATE INDEX idx_audit_logs_resource ON audit_logs(tenant_id, resource_type, resource_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(tenant_id, action);
CREATE INDEX idx_audit_logs_created ON audit_logs(tenant_id, created_at DESC);
