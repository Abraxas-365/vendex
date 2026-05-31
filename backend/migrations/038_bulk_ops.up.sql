-- 038_bulk_ops.up.sql
-- Bulk operations domain: tracks admin-initiated batch operations on products/orders.

CREATE TABLE bulk_operations (
    id              TEXT        NOT NULL,
    tenant_id       TEXT        NOT NULL,
    type            TEXT        NOT NULL,           -- 'price_update', 'status_change', 'tag_add', 'tag_remove', 'delete'
    resource_type   TEXT        NOT NULL,           -- 'product', 'order'
    status          TEXT        NOT NULL DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed', 'cancelled'
    total_items     INT         NOT NULL DEFAULT 0,
    processed_items INT         NOT NULL DEFAULT 0,
    failed_items    INT         NOT NULL DEFAULT 0,
    parameters      JSONB       NOT NULL DEFAULT '{}',
    errors          JSONB                DEFAULT '[]',
    created_by      TEXT        NOT NULL,
    started_at      TIMESTAMPTZ,
    completed_at    TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_bulk_ops_tenant   ON bulk_operations(tenant_id);
CREATE INDEX idx_bulk_ops_status   ON bulk_operations(tenant_id, status);

CREATE TABLE bulk_operation_items (
    id           TEXT        NOT NULL,
    tenant_id    TEXT        NOT NULL,
    operation_id TEXT        NOT NULL REFERENCES bulk_operations(id) ON DELETE CASCADE,
    resource_id  TEXT        NOT NULL,
    status       TEXT        NOT NULL DEFAULT 'pending', -- 'pending', 'success', 'failed'
    error_message TEXT,
    processed_at TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE INDEX idx_bulk_op_items_op ON bulk_operation_items(operation_id);
CREATE INDEX idx_bulk_op_items_tenant ON bulk_operation_items(tenant_id);
