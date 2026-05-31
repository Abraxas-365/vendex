-- Template pages for storefront — editable via CMS / AI workspace agent.
-- Slug prefix _ marks these as system templates (not shown in page nav).

-- ═══════════════════════════════════════════════════════════════════════
-- PLP template — Urban Threads (fashion)
-- ═══════════════════════════════════════════════════════════════════════
INSERT INTO pages (id, tenant_id, slug, title, html, css, status, content_type, published_at)
VALUES (
  'tpl_fashion_plp', 'tnt_fashion', '_plp', 'Shop Collection',
  '<div style="background: linear-gradient(135deg, #be185d 0%, #9d174d 100%); color: #fff; padding: 3rem 2rem; border-radius: 1rem; margin-bottom: 2rem; text-align: center;">
  <h2 style="font-size: 2rem; font-weight: 700; margin: 0 0 0.5rem;">Our Collection</h2>
  <p style="font-size: 1.1rem; opacity: 0.9; max-width: 600px; margin: 0 auto;">Curated fashion pieces designed for the modern wardrobe. From everyday essentials to statement looks.</p>
</div>',
  '', 'published', 'html', NOW()
)
ON CONFLICT (tenant_id, slug) DO UPDATE SET
  title = EXCLUDED.title,
  html  = EXCLUDED.html,
  status = 'published',
  published_at = NOW();

-- ═══════════════════════════════════════════════════════════════════════
-- PLP template — Vendex Tech
-- ═══════════════════════════════════════════════════════════════════════
INSERT INTO pages (id, tenant_id, slug, title, html, css, status, content_type, published_at)
VALUES (
  'tpl_demo_plp', 'tnt_demo', '_plp', 'Shop Tech',
  '<div style="background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%); color: #fff; padding: 3rem 2rem; border-radius: 1rem; margin-bottom: 2rem; text-align: center;">
  <h2 style="font-size: 2rem; font-weight: 700; margin: 0 0 0.5rem;">Tech Essentials</h2>
  <p style="font-size: 1.1rem; opacity: 0.9; max-width: 600px; margin: 0 auto;">Premium electronics and accessories. Discover the latest in phones, laptops, audio, and smart devices.</p>
</div>',
  '', 'published', 'html', NOW()
)
ON CONFLICT (tenant_id, slug) DO UPDATE SET
  title = EXCLUDED.title,
  html  = EXCLUDED.html,
  status = 'published',
  published_at = NOW();

-- ═══════════════════════════════════════════════════════════════════════
-- PDP template — Urban Threads (fashion)
-- ═══════════════════════════════════════════════════════════════════════
INSERT INTO pages (id, tenant_id, slug, title, html, css, status, content_type, published_at)
VALUES (
  'tpl_fashion_pdp', 'tnt_fashion', '_pdp', 'Product Details',
  '<div style="border-top: 1px solid #f3f4f6; margin-top: 2rem; padding-top: 2rem;">
  <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1.5rem;">
    <div style="text-align: center; padding: 1.5rem;">
      <div style="font-size: 1.5rem; margin-bottom: 0.5rem;">🚚</div>
      <h4 style="font-weight: 600; margin: 0 0 0.25rem; color: #111827;">Free Shipping</h4>
      <p style="font-size: 0.875rem; color: #6b7280; margin: 0;">On orders over $100</p>
    </div>
    <div style="text-align: center; padding: 1.5rem;">
      <div style="font-size: 1.5rem; margin-bottom: 0.5rem;">↩️</div>
      <h4 style="font-weight: 600; margin: 0 0 0.25rem; color: #111827;">Easy Returns</h4>
      <p style="font-size: 0.875rem; color: #6b7280; margin: 0;">30-day return policy</p>
    </div>
    <div style="text-align: center; padding: 1.5rem;">
      <div style="font-size: 1.5rem; margin-bottom: 0.5rem;">💎</div>
      <h4 style="font-weight: 600; margin: 0 0 0.25rem; color: #111827;">Premium Quality</h4>
      <p style="font-size: 0.875rem; color: #6b7280; margin: 0;">Handpicked materials</p>
    </div>
  </div>
</div>',
  '', 'published', 'html', NOW()
)
ON CONFLICT (tenant_id, slug) DO UPDATE SET
  title = EXCLUDED.title,
  html  = EXCLUDED.html,
  status = 'published',
  published_at = NOW();

-- ═══════════════════════════════════════════════════════════════════════
-- PDP template — Vendex Tech
-- ═══════════════════════════════════════════════════════════════════════
INSERT INTO pages (id, tenant_id, slug, title, html, css, status, content_type, published_at)
VALUES (
  'tpl_demo_pdp', 'tnt_demo', '_pdp', 'Product Details',
  '<div style="border-top: 1px solid #f3f4f6; margin-top: 2rem; padding-top: 2rem;">
  <div style="display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1.5rem;">
    <div style="text-align: center; padding: 1.5rem;">
      <div style="font-size: 1.5rem; margin-bottom: 0.5rem;">🚀</div>
      <h4 style="font-weight: 600; margin: 0 0 0.25rem; color: #111827;">Fast Shipping</h4>
      <p style="font-size: 0.875rem; color: #6b7280; margin: 0;">2-day delivery available</p>
    </div>
    <div style="text-align: center; padding: 1.5rem;">
      <div style="font-size: 1.5rem; margin-bottom: 0.5rem;">🛡️</div>
      <h4 style="font-weight: 600; margin: 0 0 0.25rem; color: #111827;">1-Year Warranty</h4>
      <p style="font-size: 0.875rem; color: #6b7280; margin: 0;">Full manufacturer coverage</p>
    </div>
    <div style="text-align: center; padding: 1.5rem;">
      <div style="font-size: 1.5rem; margin-bottom: 0.5rem;">🔧</div>
      <h4 style="font-weight: 600; margin: 0 0 0.25rem; color: #111827;">Tech Support</h4>
      <p style="font-size: 0.875rem; color: #6b7280; margin: 0;">Expert help when you need it</p>
    </div>
  </div>
</div>',
  '', 'published', 'html', NOW()
)
ON CONFLICT (tenant_id, slug) DO UPDATE SET
  title = EXCLUDED.title,
  html  = EXCLUDED.html,
  status = 'published',
  published_at = NOW();

-- ═══════════════════════════════════════════════════════════════════════
-- Home template — Urban Threads (fashion)
-- ═══════════════════════════════════════════════════════════════════════
INSERT INTO pages (id, tenant_id, slug, title, html, css, status, content_type, published_at)
VALUES (
  'tpl_fashion_home', 'tnt_fashion', '_home', 'Home Content',
  '<div style="background: #fdf2f8; border-radius: 1rem; padding: 2.5rem; text-align: center;">
  <h3 style="font-size: 1.5rem; font-weight: 700; color: #831843; margin: 0 0 0.75rem;">New Season, New Style</h3>
  <p style="color: #9d174d; font-size: 1rem; max-width: 500px; margin: 0 auto;">Explore our latest arrivals — fresh looks for every occasion, crafted with care and designed to last.</p>
</div>',
  '', 'published', 'html', NOW()
)
ON CONFLICT (tenant_id, slug) DO UPDATE SET
  title = EXCLUDED.title,
  html  = EXCLUDED.html,
  status = 'published',
  published_at = NOW();

-- ═══════════════════════════════════════════════════════════════════════
-- Home template — Vendex Tech
-- ═══════════════════════════════════════════════════════════════════════
INSERT INTO pages (id, tenant_id, slug, title, html, css, status, content_type, published_at)
VALUES (
  'tpl_demo_home', 'tnt_demo', '_home', 'Home Content',
  '<div style="background: #eef2ff; border-radius: 1rem; padding: 2.5rem; text-align: center;">
  <h3 style="font-size: 1.5rem; font-weight: 700; color: #3730a3; margin: 0 0 0.75rem;">Deals of the Week</h3>
  <p style="color: #4f46e5; font-size: 1rem; max-width: 500px; margin: 0 auto;">Save up to 30% on selected electronics. Limited time offer on our best-selling tech essentials.</p>
</div>',
  '', 'published', 'html', NOW()
)
ON CONFLICT (tenant_id, slug) DO UPDATE SET
  title = EXCLUDED.title,
  html  = EXCLUDED.html,
  status = 'published',
  published_at = NOW();

-- ═══════════════════════════════════════════════════════════════════════
-- Cart template — Urban Threads (fashion)
-- ═══════════════════════════════════════════════════════════════════════
INSERT INTO pages (id, tenant_id, slug, title, html, css, status, content_type, published_at)
VALUES (
  'tpl_fashion_cart', 'tnt_fashion', '_cart', 'Cart Promo',
  '<div style="background: linear-gradient(90deg, #fdf2f8, #fce7f3); border-radius: 0.75rem; padding: 1rem 1.5rem; display: flex; align-items: center; justify-content: center; gap: 0.5rem;">
  <span style="font-size: 1.1rem;">✨</span>
  <p style="margin: 0; font-size: 0.875rem; color: #831843; font-weight: 500;">Use code <strong>STYLE15</strong> for 15% off your first order!</p>
</div>',
  '', 'published', 'html', NOW()
)
ON CONFLICT (tenant_id, slug) DO UPDATE SET
  title = EXCLUDED.title,
  html  = EXCLUDED.html,
  status = 'published',
  published_at = NOW();

-- ═══════════════════════════════════════════════════════════════════════
-- Cart template — Vendex Tech
-- ═══════════════════════════════════════════════════════════════════════
INSERT INTO pages (id, tenant_id, slug, title, html, css, status, content_type, published_at)
VALUES (
  'tpl_demo_cart', 'tnt_demo', '_cart', 'Cart Promo',
  '<div style="background: linear-gradient(90deg, #eef2ff, #e0e7ff); border-radius: 0.75rem; padding: 1rem 1.5rem; display: flex; align-items: center; justify-content: center; gap: 0.5rem;">
  <span style="font-size: 1.1rem;">🔥</span>
  <p style="margin: 0; font-size: 0.875rem; color: #3730a3; font-weight: 500;">Use code <strong>TECH10</strong> for 10% off on accessories!</p>
</div>',
  '', 'published', 'html', NOW()
)
ON CONFLICT (tenant_id, slug) DO UPDATE SET
  title = EXCLUDED.title,
  html  = EXCLUDED.html,
  status = 'published',
  published_at = NOW();
