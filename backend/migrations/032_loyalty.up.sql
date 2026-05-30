CREATE TABLE loyalty_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    customer_id UUID NOT NULL,
    points_balance INT NOT NULL DEFAULT 0,
    tier TEXT NOT NULL DEFAULT 'bronze',
    lifetime_points INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, customer_id)
);
CREATE INDEX idx_loyalty_accounts_tenant ON loyalty_accounts(tenant_id);

CREATE TABLE loyalty_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    account_id UUID NOT NULL REFERENCES loyalty_accounts(id),
    type TEXT NOT NULL, -- earn, redeem, expire, adjust
    points INT NOT NULL,
    reference TEXT, -- order_id, reward_id, etc.
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_loyalty_transactions_tenant ON loyalty_transactions(tenant_id);
CREATE INDEX idx_loyalty_transactions_account ON loyalty_transactions(account_id);

CREATE TABLE loyalty_rewards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT,
    points_cost INT NOT NULL,
    reward_type TEXT NOT NULL, -- discount, free_shipping, gift_card
    value_cents INT NOT NULL DEFAULT 0,
    active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_loyalty_rewards_tenant ON loyalty_rewards(tenant_id);
