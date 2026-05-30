CREATE TABLE customer_groups (
    id          VARCHAR(36)  NOT NULL,
    tenant_id   VARCHAR(36)  NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT         NULL DEFAULT '',
    rules       JSONB        NOT NULL DEFAULT '{}',
    auto_assign BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id)
);

CREATE TABLE customer_group_memberships (
    id          VARCHAR(36) NOT NULL,
    group_id    VARCHAR(36) NOT NULL,
    customer_id VARCHAR(36) NOT NULL,
    tenant_id   VARCHAR(36) NOT NULL,
    assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_memberships_group FOREIGN KEY (group_id) REFERENCES customer_groups(id) ON DELETE CASCADE,
    CONSTRAINT fk_memberships_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    CONSTRAINT uq_membership UNIQUE (group_id, customer_id)
);

CREATE INDEX idx_customer_groups_tenant ON customer_groups(tenant_id);
CREATE INDEX idx_memberships_group ON customer_group_memberships(group_id);
CREATE INDEX idx_memberships_customer ON customer_group_memberships(customer_id);
CREATE INDEX idx_memberships_tenant ON customer_group_memberships(tenant_id);
