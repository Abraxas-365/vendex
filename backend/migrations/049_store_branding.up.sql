-- Store branding: tenant-specific storefront customization
CREATE TABLE store_branding (
    tenant_id       VARCHAR(255) PRIMARY KEY REFERENCES tenants(id) ON DELETE CASCADE,
    tagline         TEXT NOT NULL DEFAULT '',
    hero_title      TEXT NOT NULL DEFAULT '',
    hero_subtitle   TEXT NOT NULL DEFAULT '',
    accent_color    VARCHAR(20) NOT NULL DEFAULT '#6366f1',
    bg_style        VARCHAR(20) NOT NULL DEFAULT 'gradient',
    trust_badges    JSONB NOT NULL DEFAULT '[]',
    announcement    TEXT NOT NULL DEFAULT '',
    updated_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Seed branding for demo tech store
INSERT INTO store_branding (tenant_id, tagline, hero_title, hero_subtitle, accent_color, bg_style, trust_badges, announcement) VALUES
('tnt_demo', 'Premium Tech', 'Welcome to Vendex Demo', 'Discover the latest in tech — from flagship phones to studio headphones.', '#6366f1', 'gradient',
 '[{"icon":"truck","title":"Free Shipping","desc":"On orders over $50"},{"icon":"shield","title":"2-Year Warranty","desc":"Full coverage included"},{"icon":"refresh","title":"Easy Returns","desc":"30-day hassle-free returns"},{"icon":"headphones","title":"24/7 Support","desc":"Expert help anytime"}]',
 'New arrivals this season');

-- Seed branding for fashion store
INSERT INTO store_branding (tenant_id, tagline, hero_title, hero_subtitle, accent_color, bg_style, trust_badges, announcement) VALUES
('tnt_fashion', 'Curated Fashion', 'Elevate Your Style', 'Timeless pieces crafted with intention — where luxury meets everyday elegance.', '#1a1a1a', 'minimal',
 '[{"icon":"leaf","title":"Sustainable","desc":"Ethically sourced materials"},{"icon":"gem","title":"Premium Quality","desc":"Handpicked fabrics"},{"icon":"truck","title":"Free Shipping","desc":"On orders over $100"},{"icon":"refresh","title":"Easy Returns","desc":"14-day free returns"}]',
 'Summer collection now available');

-- Seed store settings for fashion tenant
INSERT INTO store_settings (tenant_id, store_name, store_email, currency, timezone) VALUES
('tnt_fashion', 'Urban Threads', 'hello@urbanthreads.co', 'USD', 'America/New_York')
ON CONFLICT (tenant_id) DO UPDATE SET store_name = EXCLUDED.store_name;

-- Update demo store name
UPDATE store_settings SET store_name = 'Vendex Tech' WHERE tenant_id = 'tnt_demo';
