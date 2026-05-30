CREATE TABLE customer_credentials (
    id              VARCHAR(36)  NOT NULL,
    customer_id     VARCHAR(36)  NOT NULL,
    tenant_id       VARCHAR(36)  NOT NULL,
    email           VARCHAR(255) NOT NULL,
    password_hash   TEXT         NOT NULL,
    created_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id),
    CONSTRAINT fk_customer_credentials_customer FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE CASCADE,
    CONSTRAINT uq_customer_credentials_email_tenant UNIQUE (tenant_id, email)
);

CREATE INDEX idx_customer_credentials_tenant ON customer_credentials(tenant_id);
CREATE INDEX idx_customer_credentials_customer ON customer_credentials(customer_id);
