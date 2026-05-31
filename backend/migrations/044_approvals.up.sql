-- 044: approval_requests table for human-in-the-loop agent action review

CREATE TABLE IF NOT EXISTS approval_requests (
    id           UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id    VARCHAR(255) NOT NULL,
    session_id   VARCHAR(255) NOT NULL DEFAULT '',
    tool_name    VARCHAR(255) NOT NULL,
    tool_input   JSONB        NOT NULL DEFAULT '{}',
    status       VARCHAR(50)  NOT NULL DEFAULT 'pending',
    reason       TEXT         NOT NULL DEFAULT '',
    requested_by VARCHAR(255) NOT NULL DEFAULT '',
    reviewed_by  VARCHAR(255) NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    reviewed_at  TIMESTAMPTZ
);

CREATE INDEX idx_approval_requests_tenant_status  ON approval_requests(tenant_id, status);
CREATE INDEX idx_approval_requests_tenant_created ON approval_requests(tenant_id, created_at DESC);
