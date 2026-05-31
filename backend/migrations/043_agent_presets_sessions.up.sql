-- Agent Presets
CREATE TABLE IF NOT EXISTS presets (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    name TEXT NOT NULL,
    slug TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    version TEXT NOT NULL DEFAULT '1.0.0',
    image TEXT NOT NULL,
    frontend_port INTEGER NOT NULL DEFAULT 8080,
    system_prompt TEXT NOT NULL DEFAULT '',
    tools_manifest JSONB NOT NULL DEFAULT '[]',
    status TEXT NOT NULL DEFAULT 'draft',
    visibility TEXT NOT NULL DEFAULT 'private',
    icon TEXT NOT NULL DEFAULT '',
    tags JSONB NOT NULL DEFAULT '[]',
    install_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_presets_slug ON presets(slug);
CREATE INDEX IF NOT EXISTS idx_presets_tenant ON presets(tenant_id);
CREATE INDEX IF NOT EXISTS idx_presets_status_visibility ON presets(status, visibility);

-- Preset Installations
CREATE TABLE IF NOT EXISTS preset_installs (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    preset_id TEXT NOT NULL REFERENCES presets(id) ON DELETE CASCADE,
    installed_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    config JSONB NOT NULL DEFAULT '{}',
    UNIQUE(tenant_id, preset_id)
);

CREATE INDEX IF NOT EXISTS idx_preset_installs_tenant ON preset_installs(tenant_id);

-- Agent Sessions
CREATE TABLE IF NOT EXISTS agent_sessions (
    id TEXT PRIMARY KEY,
    tenant_id TEXT NOT NULL,
    preset_id TEXT NOT NULL REFERENCES presets(id),
    container_id TEXT NOT NULL DEFAULT '',
    status TEXT NOT NULL DEFAULT 'creating',
    frontend_url TEXT NOT NULL DEFAULT '',
    metadata JSONB NOT NULL DEFAULT '{}',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    stopped_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_agent_sessions_tenant ON agent_sessions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_agent_sessions_tenant_status ON agent_sessions(tenant_id, status);

-- Chat Messages
CREATE TABLE IF NOT EXISTS agent_chat_messages (
    id TEXT PRIMARY KEY,
    session_id TEXT NOT NULL REFERENCES agent_sessions(id) ON DELETE CASCADE,
    role TEXT NOT NULL,
    content TEXT NOT NULL,
    tool_name TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_agent_chat_session ON agent_chat_messages(session_id, created_at);
