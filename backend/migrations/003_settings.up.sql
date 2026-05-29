CREATE TABLE IF NOT EXISTS store_settings (
    tenant_id VARCHAR(255) PRIMARY KEY,
    store_name VARCHAR(255) NOT NULL DEFAULT '',
    store_email VARCHAR(255) NOT NULL DEFAULT '',
    store_phone VARCHAR(100) NOT NULL DEFAULT '',
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    timezone VARCHAR(50) NOT NULL DEFAULT 'UTC',
    address JSONB NOT NULL DEFAULT '{}',
    logo_url TEXT NOT NULL DEFAULT '',
    favicon_url TEXT NOT NULL DEFAULT '',
    social_links JSONB NOT NULL DEFAULT '{}',
    checkout_config JSONB NOT NULL DEFAULT '{"guest_checkout": true, "require_phone": false}',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
