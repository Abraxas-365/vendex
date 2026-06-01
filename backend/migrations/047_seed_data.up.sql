-- ============================================================
-- Vendex Seed Data — realistic demo store for development/testing
-- ============================================================

-- Tenant
INSERT INTO tenants (id, company_name, status, subscription_plan, max_users, current_users, trial_expires_at)
VALUES ('tnt_demo', 'Vendex Demo Store', 'ACTIVE', 'PROFESSIONAL', 10, 2, NOW() + INTERVAL '90 days')
ON CONFLICT (id) DO NOTHING;

-- Store Settings
INSERT INTO store_settings (tenant_id, store_name, store_email, store_phone, currency, timezone, address, logo_url, social_links)
VALUES ('tnt_demo', 'Vendex Demo', 'hello@vendex.ai', '+1-555-0100', 'USD', 'America/New_York',
  '{"street":"123 Commerce Ave","city":"New York","state":"NY","zip":"10001","country":"US"}',
  '', '{"twitter":"https://twitter.com/vendex","instagram":"https://instagram.com/vendex"}')
ON CONFLICT (tenant_id) DO NOTHING;

-- ============================================================
-- Categories (6 top-level, 8 sub-categories)
-- ============================================================
INSERT INTO categories (id, tenant_id, name, slug, parent_id, description) VALUES
  ('cat_electronics', 'tnt_demo', 'Electronics', 'electronics', NULL, 'Gadgets, devices, and accessories'),
  ('cat_clothing',    'tnt_demo', 'Clothing',    'clothing',    NULL, 'Apparel for all occasions'),
  ('cat_home',        'tnt_demo', 'Home & Kitchen', 'home-kitchen', NULL, 'Everything for your home'),
  ('cat_sports',      'tnt_demo', 'Sports & Outdoors', 'sports-outdoors', NULL, 'Gear for active lifestyles'),
  ('cat_books',       'tnt_demo', 'Books',        'books',       NULL, 'Fiction, non-fiction, and more'),
  ('cat_beauty',      'tnt_demo', 'Beauty & Health', 'beauty-health', NULL, 'Personal care and wellness')
ON CONFLICT (id) DO NOTHING;

INSERT INTO categories (id, tenant_id, name, slug, parent_id, description) VALUES
  ('cat_phones',      'tnt_demo', 'Phones',       'phones',       'cat_electronics', 'Smartphones and accessories'),
  ('cat_laptops',     'tnt_demo', 'Laptops',      'laptops',      'cat_electronics', 'Notebooks and workstations'),
  ('cat_audio',       'tnt_demo', 'Audio',        'audio',        'cat_electronics', 'Headphones, speakers, and earbuds'),
  ('cat_mens',        'tnt_demo', 'Men''s',       'mens',         'cat_clothing',    'Men''s apparel'),
  ('cat_womens',      'tnt_demo', 'Women''s',     'womens',       'cat_clothing',    'Women''s apparel'),
  ('cat_furniture',   'tnt_demo', 'Furniture',    'furniture',    'cat_home',        'Tables, chairs, and storage'),
  ('cat_fitness',     'tnt_demo', 'Fitness',      'fitness',      'cat_sports',      'Gym and workout equipment'),
  ('cat_skincare',    'tnt_demo', 'Skincare',     'skincare',     'cat_beauty',      'Cleansers, moisturizers, serums')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Products (20 products across categories)
-- ============================================================
INSERT INTO products (id, tenant_id, name, description, sku, price_amount, price_currency, category_id, tags, status, stock, has_variants, slug, meta_title, meta_description) VALUES
  -- Electronics
  ('prd_iphone',     'tnt_demo', 'iPhone 16 Pro',           'The latest iPhone with A18 Pro chip, 48MP camera system, and titanium design.', 'IPHONE-16-PRO', 109900, 'USD', 'cat_phones',   '["apple","smartphone","5g"]', 'active', 45, true,  'iphone-16-pro',           'iPhone 16 Pro - Buy Now',         'Get the iPhone 16 Pro with titanium design.'),
  ('prd_macbook',    'tnt_demo', 'MacBook Pro 14"',         'M4 Pro chip, 18GB RAM, 512GB SSD. Built for professionals.', 'MBP-14-M4',     199900, 'USD', 'cat_laptops',  '["apple","laptop","m4"]',    'active', 22, true,  'macbook-pro-14',          'MacBook Pro 14" M4',              'MacBook Pro with M4 Pro chip.'),
  ('prd_airpods',    'tnt_demo', 'AirPods Pro 3',           'Active noise cancellation, adaptive audio, and USB-C charging.', 'AIRPODS-PRO-3', 24900,  'USD', 'cat_audio',    '["apple","earbuds","anc"]',  'active', 120, false, 'airpods-pro-3',           'AirPods Pro 3',                   'AirPods Pro with adaptive audio.'),
  ('prd_galaxy',     'tnt_demo', 'Samsung Galaxy S25 Ultra', 'Snapdragon 8 Elite, 200MP camera, S Pen included, titanium frame.', 'GALAXY-S25-U',  119900, 'USD', 'cat_phones',   '["samsung","android","5g"]', 'active', 38, true,  'samsung-galaxy-s25-ultra','Galaxy S25 Ultra',                'Samsung Galaxy S25 Ultra.'),
  ('prd_pixel',      'tnt_demo', 'Google Pixel 9',          'Tensor G4 chip, best-in-class camera AI, 7 years of updates.', 'PIXEL-9',       89900,  'USD', 'cat_phones',   '["google","android","ai"]',  'active', 55, false, 'google-pixel-9',          'Google Pixel 9',                  'Pixel 9 with Tensor G4.'),
  ('prd_sony_wh',    'tnt_demo', 'Sony WH-1000XM6',        'Industry-leading noise cancellation, 40hr battery, multipoint.', 'SONY-WH1000XM6',34900,  'USD', 'cat_audio',    '["sony","headphones","anc"]','active', 70, false, 'sony-wh-1000xm6',        'Sony WH-1000XM6 Headphones',      'Sony WH-1000XM6 wireless headphones.'),
  -- Clothing
  ('prd_tshirt',     'tnt_demo', 'Essential Cotton T-Shirt', 'Premium 100% organic cotton, pre-shrunk, comfortable fit.', 'ESS-TEE-001',   2900,   'USD', 'cat_mens',     '["basics","cotton","tee"]',  'active', 200, true,  'essential-cotton-tshirt', 'Essential Cotton T-Shirt',        'Premium organic cotton tee.'),
  ('prd_jeans',      'tnt_demo', 'Slim Fit Denim Jeans',     'Stretch denim, tapered leg, mid-rise. Dark indigo wash.', 'SLIM-JEAN-001', 7900,   'USD', 'cat_mens',     '["denim","jeans","slim"]',   'active', 85, true,  'slim-fit-denim-jeans',    'Slim Fit Denim Jeans',            'Stretch denim slim fit jeans.'),
  ('prd_dress',      'tnt_demo', 'Floral Midi Dress',        'Lightweight floral print dress with adjustable waist tie.', 'FLR-DRESS-001', 6500,   'USD', 'cat_womens',   '["dress","floral","midi"]',  'active', 60, true,  'floral-midi-dress',       'Floral Midi Dress',               'Beautiful floral midi dress.'),
  ('prd_hoodie',     'tnt_demo', 'Zip-Up Fleece Hoodie',     'Heavyweight fleece, full zip, kangaroo pockets.', 'FLEECE-HOOD-01',4900,   'USD', 'cat_mens',     '["hoodie","fleece","warm"]',  'active', 110, true,  'zip-up-fleece-hoodie',    'Zip-Up Fleece Hoodie',            'Heavyweight fleece hoodie.'),
  -- Home & Kitchen
  ('prd_coffeemaker','tnt_demo', 'Breville Barista Express', 'Semi-auto espresso machine with integrated grinder.', 'BREVILLE-BE',   69900,  'USD', 'cat_home',     '["coffee","espresso","breville"]', 'active', 18, false, 'breville-barista-express','Breville Barista Express',        'Semi-auto espresso machine.'),
  ('prd_blender',    'tnt_demo', 'Vitamix E310',             'Professional-grade blender, variable speed, 48oz container.', 'VITAMIX-E310',  34900,  'USD', 'cat_home',     '["blender","vitamix","kitchen"]', 'active', 30, false, 'vitamix-e310',            'Vitamix E310 Blender',            'Professional-grade Vitamix blender.'),
  ('prd_desk',       'tnt_demo', 'Standing Desk Pro',        'Electric sit-stand desk, 60x30", memory presets, cable tray.', 'DESK-STAND-01', 59900,  'USD', 'cat_furniture','["desk","standing","ergonomic"]', 'active', 25, false, 'standing-desk-pro',       'Standing Desk Pro',               'Electric sit-stand desk.'),
  -- Sports
  ('prd_yogamat',    'tnt_demo', 'Premium Yoga Mat',         'Non-slip, 6mm thick, eco-friendly TPE material.', 'YOGA-MAT-001',  3900,   'USD', 'cat_fitness',  '["yoga","mat","fitness"]',   'active', 150, false, 'premium-yoga-mat',        'Premium Yoga Mat',                'Non-slip eco-friendly yoga mat.'),
  ('prd_dumbbell',   'tnt_demo', 'Adjustable Dumbbell Set',  'From 5 to 52.5 lbs per dumbbell, replaces 15 sets.', 'ADJ-DBELL-001', 34900,  'USD', 'cat_fitness',  '["dumbbell","weights","gym"]','active', 40, false, 'adjustable-dumbbell-set', 'Adjustable Dumbbell Set',         'Adjustable dumbbells 5-52.5 lbs.'),
  -- Books
  ('prd_book_ai',    'tnt_demo', 'AI Engineering',           'Practical guide to building AI-powered applications. Chip Huyen.', 'BOOK-AI-ENG',   4499,   'USD', 'cat_books',    '["ai","engineering","tech"]','active', 80, false, 'ai-engineering-book',     'AI Engineering by Chip Huyen',    'Guide to building AI apps.'),
  ('prd_book_design','tnt_demo', 'Refactoring UI',           'Learn UI design from the creators of Tailwind CSS.', 'BOOK-REF-UI',   7900,   'USD', 'cat_books',    '["design","ui","tailwind"]','active', 65, false, 'refactoring-ui',          'Refactoring UI',                  'Learn UI design.'),
  -- Beauty
  ('prd_serum',      'tnt_demo', 'Vitamin C Brightening Serum','20% Vitamin C, hyaluronic acid, ferulic acid. 1oz.', 'SERUM-VITC-01', 2800,   'USD', 'cat_skincare', '["serum","vitaminc","brightening"]', 'active', 95, false, 'vitamin-c-brightening-serum','Vitamin C Brightening Serum',  'Vitamin C serum with HA.'),
  ('prd_sunscreen',  'tnt_demo', 'Daily Mineral Sunscreen SPF 50','Zinc oxide, lightweight, no white cast. 2oz.', 'SUN-SPF50-01',  1900,   'USD', 'cat_skincare', '["sunscreen","spf50","mineral"]', 'active', 130, false, 'daily-mineral-sunscreen-spf50','Daily Mineral Sunscreen SPF 50','Lightweight mineral sunscreen.'),
  -- Draft product
  ('prd_draft',      'tnt_demo', 'Mystery Box (Coming Soon)', 'Surprise box of curated products. Details TBA.', 'MYSTERY-BOX',   4900,   'USD', NULL,          '["mystery","surprise"]',     'draft', 0, false, 'mystery-box-coming-soon',  'Mystery Box',                     'Surprise curated box.')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Product Variants (for products with has_variants=true)
-- ============================================================
-- iPhone variants
INSERT INTO product_variants (id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active) VALUES
  ('var_ip_128_nt', 'prd_iphone', 'tnt_demo', 'IPHONE-16-PRO-128-NT', 109900, 'USD', 15, '{"storage":"128GB","color":"Natural Titanium"}', true),
  ('var_ip_256_nt', 'prd_iphone', 'tnt_demo', 'IPHONE-16-PRO-256-NT', 119900, 'USD', 12, '{"storage":"256GB","color":"Natural Titanium"}', true),
  ('var_ip_256_bk', 'prd_iphone', 'tnt_demo', 'IPHONE-16-PRO-256-BK', 119900, 'USD', 10, '{"storage":"256GB","color":"Black Titanium"}', true),
  ('var_ip_512_wt', 'prd_iphone', 'tnt_demo', 'IPHONE-16-PRO-512-WT', 139900, 'USD', 8,  '{"storage":"512GB","color":"White Titanium"}', true)
ON CONFLICT (id) DO NOTHING;

-- MacBook variants
INSERT INTO product_variants (id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active) VALUES
  ('var_mb_18_512',  'prd_macbook', 'tnt_demo', 'MBP-14-M4-18-512',  199900, 'USD', 8,  '{"ram":"18GB","storage":"512GB","color":"Space Black"}', true),
  ('var_mb_36_1tb',  'prd_macbook', 'tnt_demo', 'MBP-14-M4-36-1TB',  249900, 'USD', 6,  '{"ram":"36GB","storage":"1TB","color":"Space Black"}', true),
  ('var_mb_18_512s', 'prd_macbook', 'tnt_demo', 'MBP-14-M4-18-512S', 199900, 'USD', 8,  '{"ram":"18GB","storage":"512GB","color":"Silver"}', true)
ON CONFLICT (id) DO NOTHING;

-- Galaxy variants
INSERT INTO product_variants (id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active) VALUES
  ('var_gs_256_bk',  'prd_galaxy', 'tnt_demo', 'GALAXY-S25U-256-BK', 119900, 'USD', 12, '{"storage":"256GB","color":"Titanium Black"}', true),
  ('var_gs_512_gy',  'prd_galaxy', 'tnt_demo', 'GALAXY-S25U-512-GY', 139900, 'USD', 10, '{"storage":"512GB","color":"Titanium Gray"}', true),
  ('var_gs_1tb_bl',  'prd_galaxy', 'tnt_demo', 'GALAXY-S25U-1TB-BL', 159900, 'USD', 6,  '{"storage":"1TB","color":"Titanium Blue"}', true)
ON CONFLICT (id) DO NOTHING;

-- T-shirt variants
INSERT INTO product_variants (id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active) VALUES
  ('var_tee_s_wh',  'prd_tshirt', 'tnt_demo', 'ESS-TEE-S-WH',  2900, 'USD', 40, '{"size":"S","color":"White"}', true),
  ('var_tee_m_wh',  'prd_tshirt', 'tnt_demo', 'ESS-TEE-M-WH',  2900, 'USD', 50, '{"size":"M","color":"White"}', true),
  ('var_tee_l_bk',  'prd_tshirt', 'tnt_demo', 'ESS-TEE-L-BK',  2900, 'USD', 45, '{"size":"L","color":"Black"}', true),
  ('var_tee_xl_bk', 'prd_tshirt', 'tnt_demo', 'ESS-TEE-XL-BK', 2900, 'USD', 35, '{"size":"XL","color":"Black"}', true),
  ('var_tee_m_nv',  'prd_tshirt', 'tnt_demo', 'ESS-TEE-M-NV',  2900, 'USD', 30, '{"size":"M","color":"Navy"}', true)
ON CONFLICT (id) DO NOTHING;

-- Jeans variants
INSERT INTO product_variants (id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active) VALUES
  ('var_jeans_30',  'prd_jeans', 'tnt_demo', 'SLIM-JEAN-30', 7900, 'USD', 20, '{"waist":"30","length":"32"}', true),
  ('var_jeans_32',  'prd_jeans', 'tnt_demo', 'SLIM-JEAN-32', 7900, 'USD', 25, '{"waist":"32","length":"32"}', true),
  ('var_jeans_34',  'prd_jeans', 'tnt_demo', 'SLIM-JEAN-34', 7900, 'USD', 20, '{"waist":"34","length":"32"}', true),
  ('var_jeans_36',  'prd_jeans', 'tnt_demo', 'SLIM-JEAN-36', 7900, 'USD', 20, '{"waist":"36","length":"34"}', true)
ON CONFLICT (id) DO NOTHING;

-- Dress variants
INSERT INTO product_variants (id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active) VALUES
  ('var_dress_s',  'prd_dress', 'tnt_demo', 'FLR-DRESS-S',  6500, 'USD', 15, '{"size":"S"}', true),
  ('var_dress_m',  'prd_dress', 'tnt_demo', 'FLR-DRESS-M',  6500, 'USD', 20, '{"size":"M"}', true),
  ('var_dress_l',  'prd_dress', 'tnt_demo', 'FLR-DRESS-L',  6500, 'USD', 15, '{"size":"L"}', true)
ON CONFLICT (id) DO NOTHING;

-- Hoodie variants
INSERT INTO product_variants (id, product_id, tenant_id, sku, price_amount, price_currency, stock, options, active) VALUES
  ('var_hood_m_gy',  'prd_hoodie', 'tnt_demo', 'FLEECE-HOOD-M-GY',  4900, 'USD', 30, '{"size":"M","color":"Charcoal"}', true),
  ('var_hood_l_gy',  'prd_hoodie', 'tnt_demo', 'FLEECE-HOOD-L-GY',  4900, 'USD', 30, '{"size":"L","color":"Charcoal"}', true),
  ('var_hood_l_nv',  'prd_hoodie', 'tnt_demo', 'FLEECE-HOOD-L-NV',  4900, 'USD', 25, '{"size":"L","color":"Navy"}', true),
  ('var_hood_xl_bk', 'prd_hoodie', 'tnt_demo', 'FLEECE-HOOD-XL-BK', 4900, 'USD', 25, '{"size":"XL","color":"Black"}', true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Customers (12)
-- ============================================================
INSERT INTO customers (id, tenant_id, email, name, phone) VALUES
  ('cst_alice',   'tnt_demo', 'alice.johnson@gmail.com',   'Alice Johnson',   '+1-555-0101'),
  ('cst_bob',     'tnt_demo', 'bob.smith@outlook.com',     'Bob Smith',       '+1-555-0102'),
  ('cst_carol',   'tnt_demo', 'carol.davis@yahoo.com',     'Carol Davis',     '+1-555-0103'),
  ('cst_david',   'tnt_demo', 'david.wilson@gmail.com',    'David Wilson',    '+1-555-0104'),
  ('cst_emma',    'tnt_demo', 'emma.brown@icloud.com',     'Emma Brown',      '+1-555-0105'),
  ('cst_frank',   'tnt_demo', 'frank.lee@gmail.com',       'Frank Lee',       '+1-555-0106'),
  ('cst_grace',   'tnt_demo', 'grace.martinez@outlook.com','Grace Martinez',  '+1-555-0107'),
  ('cst_henry',   'tnt_demo', 'henry.taylor@gmail.com',    'Henry Taylor',    '+1-555-0108'),
  ('cst_iris',    'tnt_demo', 'iris.chen@gmail.com',        'Iris Chen',       '+1-555-0109'),
  ('cst_james',   'tnt_demo', 'james.anderson@yahoo.com',  'James Anderson',  '+1-555-0110'),
  ('cst_kate',    'tnt_demo', 'kate.thomas@icloud.com',    'Kate Thomas',     '+1-555-0111'),
  ('cst_leo',     'tnt_demo', 'leo.garcia@gmail.com',      'Leo Garcia',      '+1-555-0112')
ON CONFLICT (id) DO NOTHING;

-- Customer addresses
INSERT INTO customer_addresses (id, customer_id, tenant_id, street, city, state, postal_code, country, is_default) VALUES
  ('addr_alice', 'cst_alice', 'tnt_demo', '456 Oak St',     'Brooklyn',      'NY', '11201', 'US', true),
  ('addr_bob',   'cst_bob',   'tnt_demo', '789 Pine Ave',   'San Francisco', 'CA', '94102', 'US', true),
  ('addr_carol', 'cst_carol', 'tnt_demo', '321 Elm Blvd',   'Austin',        'TX', '78701', 'US', true),
  ('addr_david', 'cst_david', 'tnt_demo', '555 Maple Dr',   'Seattle',       'WA', '98101', 'US', true),
  ('addr_emma',  'cst_emma',  'tnt_demo', '100 Cherry Ln',  'Portland',      'OR', '97201', 'US', true),
  ('addr_frank', 'cst_frank', 'tnt_demo', '200 Birch Rd',   'Chicago',       'IL', '60601', 'US', true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Orders (15 orders with varied statuses and dates)
-- ============================================================
INSERT INTO orders (id, tenant_id, customer_id, status, total_amount, total_currency, subtotal_amount, shipping_amount, tax_amount, discount_amount, shipping_address, payment_status, payment_method, shipping_method, created_at) VALUES
  ('ord_001', 'tnt_demo', 'cst_alice', 'delivered',  115895, 'USD', 109900, 0, 8995, 3000, '{"street":"456 Oak St","city":"Brooklyn","state":"NY","zip":"11201","country":"US"}', 'paid', 'card', 'Standard', NOW() - INTERVAL '28 days'),
  ('ord_002', 'tnt_demo', 'cst_bob',   'delivered',  211893, 'USD', 199900, 0, 16393, 4400, '{"street":"789 Pine Ave","city":"San Francisco","state":"CA","zip":"94102","country":"US"}', 'paid', 'card', 'Express', NOW() - INTERVAL '25 days'),
  ('ord_003', 'tnt_demo', 'cst_carol', 'delivered',  7258,  'USD',  6500, 0, 758,  0,    '{"street":"321 Elm Blvd","city":"Austin","state":"TX","zip":"78701","country":"US"}', 'paid', 'card', 'Standard', NOW() - INTERVAL '22 days'),
  ('ord_004', 'tnt_demo', 'cst_david', 'delivered',  26382, 'USD',  24900, 0, 2482, 1000, '{"street":"555 Maple Dr","city":"Seattle","state":"WA","zip":"98101","country":"US"}', 'paid', 'paypal', 'Standard', NOW() - INTERVAL '20 days'),
  ('ord_005', 'tnt_demo', 'cst_emma',  'shipped',    38890, 'USD',  34900, 999, 2991, 0,  '{"street":"100 Cherry Ln","city":"Portland","state":"OR","zip":"97201","country":"US"}', 'paid', 'card', 'Express', NOW() - INTERVAL '14 days'),
  ('ord_006', 'tnt_demo', 'cst_frank', 'shipped',    75492, 'USD',  69900, 0, 5592, 0,   '{"street":"200 Birch Rd","city":"Chicago","state":"IL","zip":"60601","country":"US"}', 'paid', 'card', 'Standard', NOW() - INTERVAL '12 days'),
  ('ord_007', 'tnt_demo', 'cst_alice', 'delivered',  9412,  'USD',  8700, 0, 712,  0,    '{"street":"456 Oak St","city":"Brooklyn","state":"NY","zip":"11201","country":"US"}', 'paid', 'card', 'Standard', NOW() - INTERVAL '18 days'),
  ('ord_008', 'tnt_demo', 'cst_grace', 'processing', 127790,'USD', 119900, 0, 7890, 0,   '{"street":"77 Sunset Blvd","city":"Miami","state":"FL","zip":"33101","country":"US"}', 'paid', 'card', 'Express', NOW() - INTERVAL '5 days'),
  ('ord_009', 'tnt_demo', 'cst_henry', 'confirmed',  4745,  'USD',  4499, 0, 246,  0,    '{"street":"88 River Rd","city":"Denver","state":"CO","zip":"80201","country":"US"}', 'paid', 'paypal', 'Standard', NOW() - INTERVAL '3 days'),
  ('ord_010', 'tnt_demo', 'cst_iris',  'pending',    36390, 'USD',  34900, 0, 1490, 0,   '{"street":"99 Garden Way","city":"Boston","state":"MA","zip":"02101","country":"US"}', 'pending', NULL, NULL, NOW() - INTERVAL '1 day'),
  ('ord_011', 'tnt_demo', 'cst_james', 'confirmed',  11580, 'USD',  10800, 0, 780,  0,   '{"street":"12 Lake St","city":"Minneapolis","state":"MN","zip":"55401","country":"US"}', 'paid', 'card', 'Standard', NOW() - INTERVAL '2 days'),
  ('ord_012', 'tnt_demo', 'cst_kate',  'pending',    2132,  'USD',  1900, 0, 232,  0,    '{"street":"45 Hill Ave","city":"Nashville","state":"TN","zip":"37201","country":"US"}', 'pending', NULL, NULL, NOW()),
  ('ord_013', 'tnt_demo', 'cst_bob',   'delivered',  37690, 'USD',  34900, 999, 1791, 0, '{"street":"789 Pine Ave","city":"San Francisco","state":"CA","zip":"94102","country":"US"}', 'paid', 'card', 'Express', NOW() - INTERVAL '30 days'),
  ('ord_014', 'tnt_demo', 'cst_leo',   'cancelled',  89900, 'USD',  89900, 0, 0,    0,   '{"street":"33 Park Row","city":"Phoenix","state":"AZ","zip":"85001","country":"US"}', 'refunded', 'card', NULL, NOW() - INTERVAL '10 days'),
  ('ord_015', 'tnt_demo', 'cst_carol', 'delivered',  5474,  'USD',  4900, 0, 574,  0,    '{"street":"321 Elm Blvd","city":"Austin","state":"TX","zip":"78701","country":"US"}', 'paid', 'card', 'Standard', NOW() - INTERVAL '15 days')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Order Items
-- ============================================================
INSERT INTO order_items (id, order_id, product_id, product_name, quantity, unit_price_amount, unit_price_currency, total_amount, total_currency) VALUES
  ('oi_001', 'ord_001', 'prd_iphone',     'iPhone 16 Pro',            1, 109900, 'USD', 109900, 'USD'),
  ('oi_002', 'ord_002', 'prd_macbook',    'MacBook Pro 14"',          1, 199900, 'USD', 199900, 'USD'),
  ('oi_003', 'ord_003', 'prd_dress',      'Floral Midi Dress',        1, 6500,   'USD', 6500,   'USD'),
  ('oi_004', 'ord_004', 'prd_airpods',    'AirPods Pro 3',            1, 24900,  'USD', 24900,  'USD'),
  ('oi_005', 'ord_005', 'prd_sony_wh',    'Sony WH-1000XM6',         1, 34900,  'USD', 34900,  'USD'),
  ('oi_006', 'ord_006', 'prd_coffeemaker','Breville Barista Express', 1, 69900,  'USD', 69900,  'USD'),
  ('oi_007a','ord_007', 'prd_tshirt',     'Essential Cotton T-Shirt', 2, 2900,   'USD', 5800,   'USD'),
  ('oi_007b','ord_007', 'prd_serum',      'Vitamin C Brightening Serum',1,2800,   'USD', 2800,   'USD'),
  ('oi_008', 'ord_008', 'prd_galaxy',     'Samsung Galaxy S25 Ultra', 1, 119900, 'USD', 119900, 'USD'),
  ('oi_009', 'ord_009', 'prd_book_ai',    'AI Engineering',           1, 4499,   'USD', 4499,   'USD'),
  ('oi_010', 'ord_010', 'prd_dumbbell',   'Adjustable Dumbbell Set',  1, 34900,  'USD', 34900,  'USD'),
  ('oi_011a','ord_011', 'prd_tshirt',     'Essential Cotton T-Shirt', 2, 2900,   'USD', 5800,   'USD'),
  ('oi_011b','ord_011', 'prd_hoodie',     'Zip-Up Fleece Hoodie',     1, 4900,   'USD', 4900,   'USD'),
  ('oi_012', 'ord_012', 'prd_sunscreen',  'Daily Mineral Sunscreen',  1, 1900,   'USD', 1900,   'USD'),
  ('oi_013', 'ord_013', 'prd_blender',    'Vitamix E310',             1, 34900,  'USD', 34900,  'USD'),
  ('oi_014', 'ord_014', 'prd_pixel',      'Google Pixel 9',           1, 89900,  'USD', 89900,  'USD'),
  ('oi_015', 'ord_015', 'prd_hoodie',     'Zip-Up Fleece Hoodie',     1, 4900,   'USD', 4900,   'USD')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Collections (4)
-- ============================================================
INSERT INTO collections (id, tenant_id, name, slug, description, type, is_active, sort_order, published_at) VALUES
  ('col_bestsellers', 'tnt_demo', 'Best Sellers',    'best-sellers',    'Our most popular products',         'manual', true, 1, NOW()),
  ('col_new',         'tnt_demo', 'New Arrivals',    'new-arrivals',    'Just landed — the latest products', 'manual', true, 2, NOW()),
  ('col_tech',        'tnt_demo', 'Tech Essentials', 'tech-essentials', 'Must-have gadgets and accessories', 'manual', true, 3, NOW()),
  ('col_summer',      'tnt_demo', 'Summer Sale',     'summer-sale',     'Hot deals for the season',          'manual', true, 4, NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO collection_products (id, tenant_id, collection_id, product_id, sort_order) VALUES
  ('cp_01', 'tnt_demo', 'col_bestsellers', 'prd_iphone',     1),
  ('cp_02', 'tnt_demo', 'col_bestsellers', 'prd_airpods',    2),
  ('cp_03', 'tnt_demo', 'col_bestsellers', 'prd_tshirt',     3),
  ('cp_04', 'tnt_demo', 'col_bestsellers', 'prd_coffeemaker',4),
  ('cp_05', 'tnt_demo', 'col_new',         'prd_galaxy',     1),
  ('cp_06', 'tnt_demo', 'col_new',         'prd_sony_wh',    2),
  ('cp_07', 'tnt_demo', 'col_new',         'prd_book_ai',    3),
  ('cp_08', 'tnt_demo', 'col_tech',        'prd_iphone',     1),
  ('cp_09', 'tnt_demo', 'col_tech',        'prd_macbook',    2),
  ('cp_10', 'tnt_demo', 'col_tech',        'prd_galaxy',     3),
  ('cp_11', 'tnt_demo', 'col_tech',        'prd_airpods',    4),
  ('cp_12', 'tnt_demo', 'col_tech',        'prd_pixel',      5),
  ('cp_13', 'tnt_demo', 'col_summer',      'prd_tshirt',     1),
  ('cp_14', 'tnt_demo', 'col_summer',      'prd_dress',      2),
  ('cp_15', 'tnt_demo', 'col_summer',      'prd_sunscreen',  3),
  ('cp_16', 'tnt_demo', 'col_summer',      'prd_yogamat',    4)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Promos (5)
-- ============================================================
INSERT INTO promos (id, tenant_id, code, type, value, min_order_amount, max_uses, used_count, starts_at, ends_at, active, stackable) VALUES
  ('prm_welcome',  'tnt_demo', 'WELCOME10',  'percentage',   1000, 2000,  NULL, 24, NOW() - INTERVAL '60 days', NOW() + INTERVAL '30 days', true,  false),
  ('prm_summer',   'tnt_demo', 'SUMMER25',   'percentage',   2500, 5000,  100,  12, NOW() - INTERVAL '15 days', NOW() + INTERVAL '45 days', true,  false),
  ('prm_flat20',   'tnt_demo', 'SAVE20',     'fixed_amount', 2000, 10000, 50,   8,  NOW() - INTERVAL '30 days', NOW() + INTERVAL '60 days', true,  false),
  ('prm_freeship', 'tnt_demo', 'FREESHIP',   'free_shipping', 0,   5000,  NULL, 35, NOW() - INTERVAL '90 days', NULL,                       true,  true),
  ('prm_expired',  'tnt_demo', 'FLASH50',    'percentage',   5000, 0,     20,   20, NOW() - INTERVAL '60 days', NOW() - INTERVAL '30 days', false, false)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Pages (5)
-- ============================================================
INSERT INTO pages (id, tenant_id, slug, title, html, status, version, created_by, published_at, content_type) VALUES
  ('pg_about',    'tnt_demo', 'about',    'About Us',        '<h1>About Vendex Demo</h1><p>We''re a modern e-commerce platform powered by AI. Our mission is to help merchants build and grow their online stores effortlessly.</p><p>Founded in 2024, we believe the future of commerce is intelligent, personalized, and accessible to everyone.</p>', 'published', 1, 'admin', NOW() - INTERVAL '30 days', 'html'),
  ('pg_contact',  'tnt_demo', 'contact',  'Contact Us',      '<h1>Contact Us</h1><p>Email: hello@vendex.ai</p><p>Phone: +1-555-0100</p><p>Address: 123 Commerce Ave, New York, NY 10001</p>', 'published', 1, 'admin', NOW() - INTERVAL '30 days', 'html'),
  ('pg_faq',      'tnt_demo', 'faq',      'FAQ',             '<h1>Frequently Asked Questions</h1><h3>How long does shipping take?</h3><p>Standard shipping takes 5-7 business days. Express shipping is 2-3 business days.</p><h3>What is your return policy?</h3><p>We offer 30-day returns on all unused items in original packaging.</p><h3>Do you ship internationally?</h3><p>Yes! We ship to over 50 countries.</p>', 'published', 2, 'admin', NOW() - INTERVAL '20 days', 'html'),
  ('pg_privacy',  'tnt_demo', 'privacy',  'Privacy Policy',  '<h1>Privacy Policy</h1><p>Your privacy is important to us. We collect only the data needed to process your orders and improve your experience.</p>', 'published', 1, 'admin', NOW() - INTERVAL '30 days', 'html'),
  ('pg_landing',  'tnt_demo', 'summer-sale-2025', 'Summer Sale 2025', '<h1>Summer Sale!</h1><p>Up to 25% off on selected items. Use code SUMMER25 at checkout.</p>', 'draft', 1, 'admin', NULL, 'html')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Shipping Zones & Rates
-- ============================================================
INSERT INTO shipping_zones (id, tenant_id, name, countries, states) VALUES
  ('sz_domestic', 'tnt_demo', 'Domestic US',   '["US"]', '[]'),
  ('sz_canada',   'tnt_demo', 'Canada',        '["CA"]', '[]'),
  ('sz_intl',     'tnt_demo', 'International', '["GB","DE","FR","JP","AU"]', '[]')
ON CONFLICT (id) DO NOTHING;

INSERT INTO shipping_rates (id, zone_id, tenant_id, name, type, price_amount, price_currency, est_days_min, est_days_max, min_order_amount, active) VALUES
  ('sr_dom_std',  'sz_domestic', 'tnt_demo', 'Standard Shipping', 'flat',  999,  'USD', 5, 7, NULL,  true),
  ('sr_dom_exp',  'sz_domestic', 'tnt_demo', 'Express Shipping',  'flat',  1999, 'USD', 2, 3, NULL,  true),
  ('sr_dom_free', 'sz_domestic', 'tnt_demo', 'Free Shipping',     'free',  0,    'USD', 5, 7, 7500,  true),
  ('sr_ca_std',   'sz_canada',   'tnt_demo', 'Standard to Canada','flat',  1499, 'USD', 7, 14, NULL, true),
  ('sr_intl_std', 'sz_intl',     'tnt_demo', 'International Standard','flat',2499,'USD', 10, 21,NULL, true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Tax Rates
-- ============================================================
INSERT INTO tax_rates (id, tenant_id, name, rate, country, state, priority, active) VALUES
  ('tax_ny',  'tnt_demo', 'New York Sales Tax',     0.08875, 'US', 'NY', 1, true),
  ('tax_ca',  'tnt_demo', 'California Sales Tax',   0.07250, 'US', 'CA', 1, true),
  ('tax_tx',  'tnt_demo', 'Texas Sales Tax',        0.06250, 'US', 'TX', 1, true),
  ('tax_wa',  'tnt_demo', 'Washington Sales Tax',   0.06500, 'US', 'WA', 1, true),
  ('tax_il',  'tnt_demo', 'Illinois Sales Tax',     0.06250, 'US', 'IL', 1, true),
  ('tax_gst', 'tnt_demo', 'Canada GST',             0.05000, 'CA', '',   1, true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Blog Posts (4)
-- ============================================================
INSERT INTO blog_categories (id, tenant_id, name, slug, description, sort_order) VALUES
  ('a1000000-0000-0000-0000-000000000001', 'tnt_demo', 'News',    'news',    'Company news and announcements', 1),
  ('a1000000-0000-0000-0000-000000000002', 'tnt_demo', 'Guides',  'guides',  'How-to guides and tutorials',     2),
  ('a1000000-0000-0000-0000-000000000003', 'tnt_demo', 'Reviews', 'reviews', 'Product reviews and comparisons', 3)
ON CONFLICT (id) DO NOTHING;

INSERT INTO blog_posts (id, tenant_id, title, slug, excerpt, content, author_id, author_name, status, published_at, tags) VALUES
  (gen_random_uuid(), 'tnt_demo', 'Welcome to Vendex Demo Store',
   'welcome-to-vendex', 'Meet your new AI-powered commerce platform.',
   '## Welcome!\n\nWe''re excited to launch the Vendex demo store. This platform showcases what''s possible when you combine modern e-commerce with AI-powered tools.\n\n### What makes Vendex different?\n\n- **92+ AI tools** for managing every aspect of your store\n- **Docker-based agent workspaces** for custom preset workflows\n- **Smart analytics** with conversion funnels and revenue insights\n\nStay tuned for more updates!',
   'admin', 'Admin', 'published', NOW() - INTERVAL '25 days', ARRAY['announcement','launch']),

  (gen_random_uuid(), 'tnt_demo', 'Top 5 Tech Gifts for 2025',
   'top-5-tech-gifts-2025', 'Our curated list of the best tech gifts this year.',
   '## Top 5 Tech Gifts\n\n1. **iPhone 16 Pro** — The ultimate smartphone with titanium design\n2. **MacBook Pro 14"** — M4 Pro power for creators\n3. **AirPods Pro 3** — Best-in-class noise cancellation\n4. **Sony WH-1000XM6** — Premium wireless headphones\n5. **Google Pixel 9** — AI-first photography\n\nAll available now in our [Tech Essentials](/collections/tech-essentials) collection.',
   'admin', 'Admin', 'published', NOW() - INTERVAL '18 days', ARRAY['tech','gifts','guide']),

  (gen_random_uuid(), 'tnt_demo', 'How to Choose the Right Espresso Machine',
   'how-to-choose-espresso-machine', 'Everything you need to know before buying.',
   '## Buying Guide: Espresso Machines\n\nChoosing an espresso machine can be overwhelming. Here''s what to consider:\n\n### Budget\n- **Under $300**: Manual/basic semi-auto\n- **$300-$700**: Semi-auto with grinder (like our Breville Barista Express)\n- **$700+**: Super-automatic\n\n### Key Features\n- Built-in grinder vs. separate\n- Boiler type (single vs. dual)\n- Pressure gauge\n- Steam wand quality\n\nOur recommendation? The **Breville Barista Express** hits the sweet spot of price, quality, and features.',
   'admin', 'Admin', 'published', NOW() - INTERVAL '10 days', ARRAY['coffee','guide','kitchen']),

  (gen_random_uuid(), 'tnt_demo', 'Summer Sale Preview: What to Expect',
   'summer-sale-preview', 'Sneak peek at our upcoming summer deals.',
   '## Summer Sale is Coming!\n\nGet ready for our biggest sale of the season. Here''s what to expect:\n\n- **Up to 25% off** select clothing and accessories\n- **Free shipping** on orders over $75\n- **New arrivals** in our summer collection\n\nUse code **SUMMER25** to save. Sale starts soon!',
   'admin', 'Admin', 'draft', NULL, ARRAY['sale','summer','preview'])
ON CONFLICT DO NOTHING;

-- ============================================================
-- Gift Cards (3)
-- ============================================================
INSERT INTO gift_cards (id, tenant_id, code, initial_amount_cents, initial_amount_currency, balance_cents, balance_currency, expires_at, active, created_by) VALUES
  ('gc_001', 'tnt_demo', 'GIFT-ABCD-1234', 5000,  'USD', 5000,  'USD', NOW() + INTERVAL '1 year', true,  'admin'),
  ('gc_002', 'tnt_demo', 'GIFT-EFGH-5678', 10000, 'USD', 7500,  'USD', NOW() + INTERVAL '1 year', true,  'admin'),
  ('gc_003', 'tnt_demo', 'GIFT-IJKL-9012', 2500,  'USD', 0,     'USD', NOW() - INTERVAL '30 days', false, 'admin')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Warehouse & Stock Levels
-- ============================================================
INSERT INTO warehouses (id, tenant_id, name, address, is_default, active) VALUES
  ('a0000000-0000-0000-0000-000000000001', 'tnt_demo', 'Main Warehouse',    '100 Logistics Pkwy, Newark, NJ 07102', true,  true),
  ('a0000000-0000-0000-0000-000000000002', 'tnt_demo', 'West Coast Fulfillment', '500 Harbor Blvd, Long Beach, CA 90802', false, true)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Subscriptions (2)
-- ============================================================
INSERT INTO subscriptions (id, tenant_id, customer_id, product_id, price_amount, price_currency, interval, status, next_billing_date, last_billed_at) VALUES
  ('sub_001', 'tnt_demo', 'cst_emma',  'prd_serum',     2800, 'USD', 'monthly', 'active', NOW() + INTERVAL '15 days', NOW() - INTERVAL '15 days'),
  ('sub_002', 'tnt_demo', 'cst_grace', 'prd_sunscreen', 1900, 'USD', 'monthly', 'active', NOW() + INTERVAL '22 days', NOW() - INTERVAL '8 days')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Agent Presets (single workspace preset)
-- ============================================================
INSERT INTO presets (id, tenant_id, name, slug, description, version, image, system_prompt, status, visibility, icon, tags, install_count) VALUES
  ('preset_webdev', 'tnt_demo', 'Vendex Workspace', 'workspace', 'Full workspace with HTML/CSS editor, live preview, and Chromium screenshots.', '1.0.0', 'vendex-preset:latest', 'You are a web designer assistant. Help the merchant build beautiful store pages.', 'active', 'public', 'palette', '["design","web","html"]', 0)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Customer Groups (2)
-- ============================================================
INSERT INTO customer_groups (id, tenant_id, name, description, created_at, updated_at) VALUES
  ('grp_vip',      'tnt_demo', 'VIP Customers',      'High-value repeat customers',  NOW() - INTERVAL '60 days', NOW()),
  ('grp_wholesale', 'tnt_demo', 'Wholesale Buyers',   'B2B wholesale customers',      NOW() - INTERVAL '45 days', NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO customer_group_memberships (id, group_id, customer_id, tenant_id, assigned_at) VALUES
  ('cgm_01', 'grp_vip', 'cst_alice', 'tnt_demo', NOW() - INTERVAL '50 days'),
  ('cgm_02', 'grp_vip', 'cst_bob',   'tnt_demo', NOW() - INTERVAL '40 days'),
  ('cgm_03', 'grp_vip', 'cst_emma',  'tnt_demo', NOW() - INTERVAL '30 days')
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Currency Rates (3)
-- ============================================================
INSERT INTO currency_rates (id, tenant_id, base_currency, target_currency, rate) VALUES
  ('cr_eur', 'tnt_demo', 'USD', 'EUR', 0.9200000000),
  ('cr_gbp', 'tnt_demo', 'USD', 'GBP', 0.7900000000),
  ('cr_cad', 'tnt_demo', 'USD', 'CAD', 1.3600000000)
ON CONFLICT (id) DO NOTHING;

-- ============================================================
-- Done!
-- ============================================================
