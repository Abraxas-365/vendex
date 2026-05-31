-- Agent trigger rules: maps store events to automated agent prompts
CREATE TABLE IF NOT EXISTS agent_triggers (
    id          VARCHAR(255) PRIMARY KEY,
    tenant_id   VARCHAR(255) NOT NULL,
    name        VARCHAR(255) NOT NULL,
    event_type  VARCHAR(255) NOT NULL,
    prompt      TEXT         NOT NULL,
    preset_id   VARCHAR(255) NOT NULL DEFAULT '',
    enabled     BOOLEAN      NOT NULL DEFAULT true,
    cooldown    INT          NOT NULL DEFAULT 300,
    last_fired_at TIMESTAMPTZ,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_triggers_tenant ON agent_triggers(tenant_id);
CREATE INDEX IF NOT EXISTS idx_agent_triggers_event  ON agent_triggers(event_type, enabled);

-- Execution log: records each time a trigger fires
CREATE TABLE IF NOT EXISTS agent_trigger_logs (
    id             VARCHAR(255) PRIMARY KEY,
    trigger_id     VARCHAR(255) NOT NULL REFERENCES agent_triggers(id) ON DELETE CASCADE,
    tenant_id      VARCHAR(255) NOT NULL,
    event_type     VARCHAR(255) NOT NULL,
    event_payload  JSONB        NOT NULL DEFAULT '{}',
    agent_response TEXT         NOT NULL DEFAULT '',
    status         VARCHAR(50)  NOT NULL DEFAULT 'success',
    created_at     TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_trigger_logs_trigger ON agent_trigger_logs(trigger_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_agent_trigger_logs_tenant  ON agent_trigger_logs(tenant_id);
