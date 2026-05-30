# Database Migration Guide

## How Migrations Work

Migrations are stored in sequential numbered files in `backend/migrations/`:
- Numbering: `001_`, `002_`, `003_`, etc.
- Each migration has `.up.sql` (apply) and `.down.sql` (rollback)
- Applied sequentially in order
- Once applied, never re-run (append-only)

## Current Migrations

### 001_initial.up.sql (13KB)
**Core commerce tables**:
- `tenants` - Multi-tenant organization
- `customers` - Commerce customers
- `products` - Product catalog
- `categories` - Product categories
- `product_variants` - Product variants
- `collections` - Product groupings
- And other core tables

**Key schema**: 
```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    phone VARCHAR(20),
    addresses JSONB,
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    UNIQUE(tenant_id, email),
    FOREIGN KEY(tenant_id) REFERENCES tenants(id)
);
```

### 004_iam.up.sql (12KB)
**Identity & Access Management tables**:
- `users` - IAM users (different from customers)
- `roles` - Authorization roles
- `permissions` - Fine-grained permissions
- `user_roles` - User→Role assignments
- `api_keys` - Programmatic access tokens
- `sessions` - User session tracking
- `refresh_tokens` - Token rotation storage
- `password_reset_tokens` - Password reset flows
- `oauth_accounts` - OAuth provider linking

**Key schemas**:
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255),
    created_at TIMESTAMP,
    updated_at TIMESTAMP,
    FOREIGN KEY(tenant_id) REFERENCES tenants(id)
);

CREATE TABLE api_keys (
    id VARCHAR(36) PRIMARY KEY,
    user_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    key VARCHAR(255) NOT NULL UNIQUE,
    scopes TEXT,
    expires_at TIMESTAMP,
    created_at TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(tenant_id) REFERENCES tenants(id)
);

CREATE TABLE refresh_tokens (
    id VARCHAR(36) PRIMARY KEY,
    user_id UUID NOT NULL,
    tenant_id UUID NOT NULL,
    token TEXT NOT NULL,
    expires_at TIMESTAMP,
    is_revoked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id),
    FOREIGN KEY(tenant_id) REFERENCES tenants(id)
);
```

### 010_cart.up.sql
**Shopping cart**:
- `carts` - Cart entity
- `cart_items` - Items in cart with quantity & price

### 011_shipping.up.sql
**Shipping management**:
- `shipping_zones` - Geographic shipping regions
- `shipping_rates` - Cost by zone and weight

### 012_tax.up.sql
**Tax calculation**:
- `tax_rates` - Tax percentage by region/category

### 013_payment.up.sql
**Payment processing**:
- `payments` - Payment records
- `payment_transactions` - Transaction history

## Creating a New Migration

### Step 1: Create Migration File
```bash
# In backend/migrations/
# Pick next number (e.g., 014 if 013 is latest)
touch 014_new_feature.up.sql
touch 014_new_feature.down.sql
```

### Step 2: Write UP Migration
```sql
-- backend/migrations/014_new_feature.up.sql

-- Create table
CREATE TABLE some_table (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    FOREIGN KEY(tenant_id) REFERENCES tenants(id)
);

-- Create indexes
CREATE INDEX idx_some_table_tenant ON some_table(tenant_id);
CREATE UNIQUE INDEX idx_some_table_unique ON some_table(tenant_id, name);

-- Add constraints
ALTER TABLE some_table ADD CONSTRAINT check_name_length CHECK (length(name) > 0);
```

### Step 3: Write DOWN Migration (Rollback)
```sql
-- backend/migrations/014_new_feature.down.sql

DROP TABLE IF EXISTS some_table;
```

### Step 4: Test Locally
```bash
# Ensure database is running
export DATABASE_URL="postgres://user:pass@localhost:5432/hada_test"

# Apply migration
go run ./cmd migrate up

# Verify schema in psql
psql $DATABASE_URL -c "\dt"  # List tables
psql $DATABASE_URL -c "\d some_table"  # Describe table

# Rollback to test down migration
go run ./cmd migrate down

# Re-apply
go run ./cmd migrate up
```

## Migration Patterns

### Adding Column to Existing Table
```sql
-- UP
ALTER TABLE orders ADD COLUMN tracking_number VARCHAR(255);
ALTER TABLE orders ADD CONSTRAINT fk_shipment_id FOREIGN KEY(shipment_id) REFERENCES shipments(id);

-- DOWN
ALTER TABLE orders DROP COLUMN tracking_number;
```

### Renaming Table
```sql
-- UP
ALTER TABLE old_name RENAME TO new_name;

-- DOWN
ALTER TABLE new_name RENAME TO old_name;
```

### Changing Column Type
```sql
-- UP
ALTER TABLE products ALTER COLUMN sku TYPE VARCHAR(100);

-- DOWN
ALTER TABLE products ALTER COLUMN sku TYPE VARCHAR(50);
```

### Adding Index
```sql
-- UP
CREATE INDEX idx_customers_email ON customers(email);
CREATE UNIQUE INDEX idx_api_keys_tenant_key ON api_keys(tenant_id, key);

-- DOWN
DROP INDEX IF EXISTS idx_customers_email;
DROP INDEX IF EXISTS idx_api_keys_tenant_key;
```

### Seeding Data
```sql
-- UP (avoid unless really needed)
INSERT INTO categories (id, tenant_id, name) VALUES 
    (gen_random_uuid(), '..tenant_id..', 'Electronics');

-- DOWN
DELETE FROM categories WHERE name = 'Electronics';
```

## Column Conventions

### IDs
```sql
id UUID PRIMARY KEY DEFAULT gen_random_uuid()        -- Primary key
tenant_id UUID NOT NULL                              -- Tenant scoping
user_id UUID NOT NULL (FOREIGN KEY users.id)        -- Reference
```

### Timestamps
```sql
created_at TIMESTAMP DEFAULT NOW()
updated_at TIMESTAMP DEFAULT NOW()
expires_at TIMESTAMP                                  -- Optional expiry
```

### JSON/Array Data
```sql
addresses JSONB                    -- Nested objects/arrays
scopes TEXT                        -- Comma-separated or JSON
metadata JSONB DEFAULT '{}'::jsonb
```

### Soft Deletes (Optional, usually avoid)
```sql
deleted_at TIMESTAMP NULL          -- If needed, don't use soft delete without reason
```

## Constraints Best Practices

### NOT NULL for Required Fields
```sql
CREATE TABLE customers (
    id UUID PRIMARY KEY,
    tenant_id UUID NOT NULL,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    ...
);
```

### UNIQUE Constraints for Key Fields
```sql
CREATE UNIQUE INDEX idx_customer_email_per_tenant ON customers(tenant_id, email);
```

### FOREIGN KEYS for Relationships
```sql
ALTER TABLE orders ADD CONSTRAINT fk_order_customer 
    FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE;
```

### CHECK Constraints for Validation
```sql
ALTER TABLE products ADD CONSTRAINT check_price_positive CHECK(price > 0);
```

## Common Issues & Solutions

### Migration Already Applied
**Problem**: Running a migration twice causes duplicate table error
**Solution**: Migrations are append-only; never modify `.up.sql` file once committed
**Fix**: Create new migration to fix it (e.g., 015_fix_previous_issue.sql)

### Foreign Key Constraint Fails
**Problem**: Can't drop table because other tables reference it
**Solution**: Drop dependent tables first, or add `ON DELETE CASCADE`
```sql
ALTER TABLE orders DROP CONSTRAINT IF EXISTS fk_order_customer;
ALTER TABLE orders ADD CONSTRAINT fk_order_customer 
    FOREIGN KEY(customer_id) REFERENCES customers(id) ON DELETE CASCADE;
```

### JSONB Index for Performance
**Problem**: Querying JSONB columns is slow
**Solution**: Create GIN index
```sql
CREATE INDEX idx_customers_addresses ON customers USING GIN(addresses);
```

### Large Table Migrations (Production)
**Pattern**: Add column with default, then backfill, then remove default
```sql
-- Step 1: Add with default
ALTER TABLE large_table ADD COLUMN new_col VARCHAR(255) DEFAULT 'temporary';

-- Step 2 (separate migration): Backfill during low-traffic time
UPDATE large_table SET new_col = computed_value WHERE new_col = 'temporary';

-- Step 3 (separate migration): Remove default and set NOT NULL
ALTER TABLE large_table ALTER COLUMN new_col DROP DEFAULT;
ALTER TABLE large_table ALTER COLUMN new_col SET NOT NULL;
```

## Deployment Checklist

Before deploying to production:

- [ ] Test UP migration on local database
- [ ] Test DOWN migration rollback
- [ ] Verify no breaking changes to existing queries
- [ ] Ensure FOREIGN KEY constraints are correct
- [ ] Check for long-running queries on production-sized data
- [ ] Plan for data backfill if adding NOT NULL columns
- [ ] Communicate schema changes to API team
- [ ] Version control: commit both `.up.sql` and `.down.sql`

## PostgreSQL Useful Commands

```bash
# Connect to database
psql postgres://user:pass@host:port/dbname

# List tables
\dt

# Describe table
\d table_name

# Show indexes
\di table_name

# Show constraints
\d table_name (includes constraints)

# Show foreign keys
SELECT * FROM information_schema.table_constraints 
WHERE table_name = 'table_name';

# Show column info
SELECT * FROM information_schema.columns 
WHERE table_name = 'table_name';
```

