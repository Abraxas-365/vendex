-- Agent memory knowledge base
-- Stores persistent, searchable knowledge entries per tenant.
-- Used by the AI agent to retain brand guidelines, product taxonomy,
-- tone of voice, and past decisions across chat sessions.

CREATE TABLE IF NOT EXISTS agent_memories (
    id          TEXT        NOT NULL,
    tenant_id   TEXT        NOT NULL,
    category    TEXT        NOT NULL DEFAULT 'general',
    title       TEXT        NOT NULL,
    content     TEXT        NOT NULL,
    tags        TEXT[]      NOT NULL DEFAULT '{}',
    source      TEXT        NOT NULL DEFAULT 'human',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    search_vector TSVECTOR GENERATED ALWAYS AS (
        setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
        setweight(to_tsvector('english', coalesce(content, '')), 'B')
    ) STORED,
    PRIMARY KEY (id)
);

CREATE INDEX idx_agent_memories_tenant    ON agent_memories(tenant_id);
CREATE INDEX idx_agent_memories_category  ON agent_memories(tenant_id, category);
CREATE INDEX idx_agent_memories_tags      ON agent_memories USING GIN(tags);
CREATE INDEX idx_agent_memories_search    ON agent_memories USING GIN(search_vector);
