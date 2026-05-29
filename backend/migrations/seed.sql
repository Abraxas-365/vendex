-- =============================================================================
-- hada-commerce seed data (matches actual DB schema)
-- =============================================================================

-- Clean up any partial inserts from before
TRUNCATE products, customers, customer_addresses, orders, order_items, pages, page_versions, promos CASCADE;

-- ---------------------------------------------------------------------------
-- Products (no slug column)
-- ---------------------------------------------------------------------------
INSERT INTO products (id, tenant_id, name, description, sku, price_amount, price_currency, stock, status, category_id, tags, images, created_at, updated_at) VALUES
  ('prod-001', 'default', 'Classic Leather Sandals',
   'Handcrafted leather sandals with cushioned sole. Perfect for summer walks.',
   'SND-001', 4999, 'USD', 150, 'active', 'cat-001',
   '["sandals", "leather", "summer", "bestseller"]',
   '["https://images.unsplash.com/photo-1603487742131-4160ec999306?w=600"]',
   NOW(), NOW()),

  ('prod-002', 'default', 'Organic Cotton T-Shirt',
   '100% organic cotton crew neck tee. Available in multiple colors.',
   'TSH-001', 2999, 'USD', 300, 'active', 'cat-002',
   '["tshirt", "cotton", "organic", "basics"]',
   '["https://images.unsplash.com/photo-1521572163474-6864f9cf17ab?w=600"]',
   NOW(), NOW()),

  ('prod-003', 'default', 'Canvas Tote Bag',
   'Durable canvas tote with reinforced handles. Eco-friendly and stylish.',
   'BAG-001', 3499, 'USD', 200, 'active', 'cat-003',
   '["bag", "tote", "canvas", "eco"]',
   '["https://images.unsplash.com/photo-1544816155-12df9643f363?w=600"]',
   NOW(), NOW()),

  ('prod-004', 'default', 'Linen Summer Dress',
   'Light and breezy linen dress. Ideal for warm weather occasions.',
   'DRS-001', 7999, 'USD', 75, 'active', 'cat-004',
   '["dress", "linen", "summer", "women"]',
   '["https://images.unsplash.com/photo-1595777457583-95e059d581b8?w=600"]',
   NOW(), NOW()),

  ('prod-005', 'default', 'Wool Beanie Hat',
   'Warm merino wool beanie. Unisex design fits all head sizes.',
   'HAT-001', 1999, 'USD', 500, 'active', 'cat-003',
   '["hat", "beanie", "wool", "winter"]',
   '["https://images.unsplash.com/photo-1576871337632-b9aef4c17ab9?w=600"]',
   NOW(), NOW()),

  ('prod-006', 'default', 'Running Sneakers Pro',
   'Lightweight running shoes with responsive cushioning.',
   'SHO-001', 12999, 'USD', 5, 'active', 'cat-001',
   '["shoes", "running", "sport", "new"]',
   '["https://images.unsplash.com/photo-1542291026-7eec264c27ff?w=600"]',
   NOW(), NOW()),

  ('prod-007', 'default', 'Draft Product (Not Published)',
   'This product is still being prepared.',
   'DRF-001', 999, 'USD', 0, 'draft', 'cat-002',
   '["draft"]', '[]', NOW(), NOW());

-- ---------------------------------------------------------------------------
-- Customers (name is single column, not first_name/last_name)
-- ---------------------------------------------------------------------------
INSERT INTO customers (id, tenant_id, email, name, phone, created_at, updated_at) VALUES
  ('cust-001', 'default', 'alice@example.com', 'Alice Johnson', '+1-555-0101', NOW(), NOW()),
  ('cust-002', 'default', 'bob@example.com',   'Bob Smith',     '+1-555-0102', NOW(), NOW()),
  ('cust-003', 'default', 'carol@example.com', 'Carol Davis',   '+1-555-0103', NOW(), NOW());

-- customer_addresses (no label column, uses postal_code not zip)
INSERT INTO customer_addresses (id, customer_id, tenant_id, street, city, state, country, postal_code, is_default) VALUES
  ('addr-001', 'cust-001', 'default', '123 Main St',   'Portland', 'OR', 'US', '97201', true),
  ('addr-002', 'cust-002', 'default', '456 Oak Ave',   'Seattle',  'WA', 'US', '98101', true),
  ('addr-003', 'cust-003', 'default', '789 Pine Blvd', 'Austin',   'TX', 'US', '73301', true);

-- ---------------------------------------------------------------------------
-- Orders (shipping_address is JSONB, not separate columns)
-- ---------------------------------------------------------------------------
INSERT INTO orders (id, tenant_id, customer_id, status, total_amount, total_currency, shipping_address, notes, created_at, updated_at) VALUES
  ('ord-001', 'default', 'cust-001', 'delivered',  8498, 'USD',
   '{"street":"123 Main St","city":"Portland","state":"OR","country":"US","postal_code":"97201"}',
   'Leave at front door', NOW() - INTERVAL '10 days', NOW() - INTERVAL '3 days'),
  ('ord-002', 'default', 'cust-002', 'shipped',    2999, 'USD',
   '{"street":"456 Oak Ave","city":"Seattle","state":"WA","country":"US","postal_code":"98101"}',
   NULL, NOW() - INTERVAL '5 days', NOW() - INTERVAL '1 day'),
  ('ord-003', 'default', 'cust-001', 'confirmed',  7999, 'USD',
   '{"street":"123 Main St","city":"Portland","state":"OR","country":"US","postal_code":"97201"}',
   NULL, NOW() - INTERVAL '2 days', NOW()),
  ('ord-004', 'default', 'cust-003', 'pending',   16498, 'USD',
   '{"street":"789 Pine Blvd","city":"Austin","state":"TX","country":"US","postal_code":"73301"}',
   'Gift wrap please', NOW() - INTERVAL '1 day', NOW()),
  ('ord-005', 'default', 'cust-002', 'cancelled',  1999, 'USD',
   '{"street":"456 Oak Ave","city":"Seattle","state":"WA","country":"US","postal_code":"98101"}',
   'Changed my mind', NOW() - INTERVAL '15 days', NOW() - INTERVAL '14 days');

-- order_items (has total_amount and total_currency columns)
INSERT INTO order_items (id, order_id, product_id, product_name, quantity, unit_price_amount, unit_price_currency, total_amount, total_currency) VALUES
  ('oi-001', 'ord-001', 'prod-001', 'Classic Leather Sandals', 1, 4999, 'USD', 4999, 'USD'),
  ('oi-002', 'ord-001', 'prod-003', 'Canvas Tote Bag',         1, 3499, 'USD', 3499, 'USD'),
  ('oi-003', 'ord-002', 'prod-002', 'Organic Cotton T-Shirt',  1, 2999, 'USD', 2999, 'USD'),
  ('oi-004', 'ord-003', 'prod-004', 'Linen Summer Dress',      1, 7999, 'USD', 7999, 'USD'),
  ('oi-005', 'ord-004', 'prod-006', 'Running Sneakers Pro',    1, 12999, 'USD', 12999, 'USD'),
  ('oi-006', 'ord-004', 'prod-003', 'Canvas Tote Bag',         1, 3499, 'USD', 3499, 'USD'),
  ('oi-007', 'ord-005', 'prod-005', 'Wool Beanie Hat',         1, 1999, 'USD', 1999, 'USD');

-- ---------------------------------------------------------------------------
-- CMS Pages (meta is JSONB, not separate columns)
-- ---------------------------------------------------------------------------
INSERT INTO pages (id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, created_at, updated_at) VALUES
  ('page-001', 'default', 'summer-sale', 'Summer Sale 2025',
'<section class="hero">
  <h1>Summer Sale</h1>
  <p>Up to 30% off on all summer essentials</p>
  <a href="/products" class="cta-btn">Shop Now</a>
</section>
<section class="features">
  <div class="feature-card">
    <h3>Free Shipping</h3>
    <p>On orders over $50</p>
  </div>
  <div class="feature-card">
    <h3>Easy Returns</h3>
    <p>30-day return policy</p>
  </div>
  <div class="feature-card">
    <h3>Secure Payment</h3>
    <p>256-bit SSL encryption</p>
  </div>
</section>',
'.hero { background: linear-gradient(135deg, #f8b400, #ff6f61); padding: 4rem 2rem; text-align: center; color: white; border-radius: 12px; margin-bottom: 2rem; }
.hero h1 { font-size: 2.5rem; margin-bottom: 0.5rem; }
.hero p { font-size: 1.2rem; opacity: 0.9; margin-bottom: 1.5rem; }
.cta-btn { display: inline-block; padding: 0.75rem 2rem; background: white; color: #ff6f61; border-radius: 8px; text-decoration: none; font-weight: bold; }
.features { display: grid; grid-template-columns: repeat(3, 1fr); gap: 1.5rem; }
.feature-card { background: #f9fafb; padding: 2rem; border-radius: 8px; text-align: center; }
.feature-card h3 { color: #1f2937; margin-bottom: 0.5rem; }',
'{"description":"Summer sale - up to 30% off","og_title":"Summer Sale 2025","keywords":"summer, sale, discount"}',
'published', 1, 'agent:store-manager', NOW(), NOW() - INTERVAL '5 days', NOW()),

  ('page-002', 'default', 'about-us', 'About Hada Store',
'<section class="about">
  <h1>About Hada Store</h1>
  <p>We are a curated marketplace focused on sustainable, handcrafted goods from independent artisans around the world.</p>
  <div class="stats">
    <div class="stat"><span class="num">500+</span><span class="label">Products</span></div>
    <div class="stat"><span class="num">50+</span><span class="label">Artisans</span></div>
    <div class="stat"><span class="num">10k+</span><span class="label">Happy Customers</span></div>
  </div>
</section>',
'.about { max-width: 700px; margin: 0 auto; text-align: center; padding: 2rem; }
.about h1 { font-size: 2rem; margin-bottom: 1rem; color: #1f2937; }
.about p { color: #6b7280; line-height: 1.8; margin-bottom: 2rem; }
.stats { display: flex; justify-content: center; gap: 3rem; }
.stat { text-align: center; }
.stat .num { display: block; font-size: 2rem; font-weight: bold; color: #4f46e5; }
.stat .label { color: #9ca3af; font-size: 0.875rem; }',
'{"description":"About Hada Store - sustainable goods","og_title":"About Hada Store","keywords":"about, store, sustainable"}',
'pending_review', 1, 'agent:content-writer', NULL, NOW() - INTERVAL '2 days', NOW()),

  ('page-003', 'default', 'faq', 'Frequently Asked Questions',
'<section class="faq">
  <h1>FAQ</h1>
  <div class="qa"><h3>What is your return policy?</h3><p>30-day hassle-free returns on all items.</p></div>
  <div class="qa"><h3>Do you ship internationally?</h3><p>Yes, we ship to over 50 countries worldwide.</p></div>
  <div class="qa"><h3>How long does shipping take?</h3><p>Standard shipping takes 5-7 business days. Express is 2-3 days.</p></div>
</section>',
'.faq { max-width: 700px; margin: 0 auto; padding: 2rem; }
.faq h1 { font-size: 2rem; margin-bottom: 1.5rem; color: #1f2937; text-align: center; }
.qa { margin-bottom: 1.5rem; padding: 1.5rem; background: #f9fafb; border-radius: 8px; }
.qa h3 { color: #1f2937; margin-bottom: 0.5rem; }
.qa p { color: #6b7280; line-height: 1.6; }',
'{"description":"FAQ - Hada Store","og_title":"FAQ","keywords":"faq, questions, help"}',
'draft', 1, 'admin:luis@hada.com', NULL, NOW() - INTERVAL '1 day', NOW());

-- Page version for published page
INSERT INTO page_versions (id, page_id, tenant_id, version, html, css, edited_by, comment, created_at) VALUES
  ('pv-001', 'page-001', 'default', 1,
   '<section class="hero"><h1>Summer Sale</h1><p>Up to 30% off on all summer essentials</p><a href="/products" class="cta-btn">Shop Now</a></section>',
   '.hero { background: linear-gradient(135deg, #f8b400, #ff6f61); padding: 4rem 2rem; text-align: center; color: white; border-radius: 12px; }',
   'agent:store-manager', 'Initial version', NOW() - INTERVAL '5 days');

-- ---------------------------------------------------------------------------
-- Promos (no description column, no min_order_currency)
-- ---------------------------------------------------------------------------
INSERT INTO promos (id, tenant_id, code, type, value, min_order_amount, max_uses, used_count, starts_at, ends_at, active, created_at) VALUES
  ('promo-001', 'default', 'SUMMER30', 'percentage', 30, 2000, 1000, 47, NOW() - INTERVAL '30 days', NOW() + INTERVAL '30 days', true, NOW() - INTERVAL '30 days'),
  ('promo-002', 'default', 'WELCOME10', 'fixed_amount', 1000, 5000, 500, 123, NOW() - INTERVAL '90 days', NOW() + INTERVAL '90 days', true, NOW() - INTERVAL '90 days'),
  ('promo-003', 'default', 'FREESHIP', 'free_shipping', 0, 0, 200, 200, NOW() - INTERVAL '60 days', NOW() - INTERVAL '1 day', false, NOW() - INTERVAL '60 days');
