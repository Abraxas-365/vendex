CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    body TEXT,
    type TEXT NOT NULL DEFAULT 'info', -- info, warning, success, error
    resource_type TEXT, -- order, product, customer, return, etc.
    resource_id TEXT,
    read BOOLEAN NOT NULL DEFAULT false,
    read_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_notifications_tenant_user ON notifications(tenant_id, user_id);
CREATE INDEX idx_notifications_unread ON notifications(tenant_id, user_id, read) WHERE read = false;
CREATE INDEX idx_notifications_created ON notifications(tenant_id, created_at DESC);
