-- Add slug to tenants for subdomain resolution ({slug}.vendex.ai)
ALTER TABLE tenants ADD COLUMN slug VARCHAR(63) UNIQUE;

-- Backfill existing tenants with a slug derived from company_name
UPDATE tenants SET slug = LOWER(REPLACE(REPLACE(company_name, ' ', '-'), '''', '')) WHERE slug IS NULL;

-- Make slug NOT NULL after backfill
ALTER TABLE tenants ALTER COLUMN slug SET NOT NULL;

-- Add index on storefronts.domain for fast custom domain lookups (already covered by WHERE clause scans but explicit index helps)
CREATE INDEX IF NOT EXISTS idx_storefronts_domain ON storefronts (domain) WHERE domain IS NOT NULL;

-- Add a tenant_domains table for multiple custom domains per tenant
CREATE TABLE tenant_domains (
    id          VARCHAR(255) PRIMARY KEY DEFAULT gen_random_uuid()::text,
    tenant_id   VARCHAR(255) NOT NULL REFERENCES tenants(id) ON DELETE CASCADE,
    domain      VARCHAR(255) NOT NULL UNIQUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    verified_at TIMESTAMP,
    created_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tenant_domains_tenant ON tenant_domains (tenant_id);
CREATE UNIQUE INDEX idx_tenant_domains_domain ON tenant_domains (domain);
