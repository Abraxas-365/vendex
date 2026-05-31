-- Template pages for storefront — editable via CMS / AI workspace agent.
-- Slug prefix _ marks these as system templates (not shown in page nav).

-- ═══════════════════════════════════════════════════════════════════════
-- PLP template — Urban Threads (fashion)
-- ═══════════════════════════════════════════════════════════════════════
INSERT INTO pages (id, tenant_id, slug, title, html, css, status, content_type, published_at)
VALUES (
  'tpl_fashion_plp', 'tnt_fashion', '_plp', 'Shop Collection',
  '
<div style="background: #111; color: #fff; padding: 4rem 3rem; border-radius: 0; margin: -2rem -1rem 2rem; position: relative; overflow: hidden;">
  <div style="position: absolute; top: 0; right: 0; width: 40%; height: 100%; background: linear-gradient(135deg, transparent 0%, rgba(190,24,93,0.08) 100%);"></div>
  <div style="position: relative; max-width: 500px;">
    <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.3em; color: rgba(255,255,255,0.4); margin: 0 0 1rem; font-weight: 500;">Collection — 2026</p>
    <h2 style="font-size: 2.5rem; font-weight: 300; margin: 0 0 1rem; letter-spacing: -0.03em; line-height: 1.15;">The <span style="font-weight: 700;">Edit</span></h2>
    <p style="font-size: 0.95rem; color: rgba(255,255,255,0.5); margin: 0; line-height: 1.7; font-weight: 300;">Pieces that define your style. Every item hand-selected for quality, fit, and timeless appeal.</p>
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
  '
<div style="margin-top: 3rem; padding-top: 2.5rem; border-top: 1px solid #f0f0f0;">
  <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 0;">
    <div style="text-align: center; padding: 2rem 1.5rem; border-right: 1px solid #f0f0f0;">
      <div style="width: 44px; height: 44px; border: 2px solid #be185d; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem;">
        <div style="width: 0; height: 0; border-left: 8px solid #be185d; border-top: 5px solid transparent; border-bottom: 5px solid transparent; margin-left: 3px;"></div>
      </div>
      <h4 style="font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #111; margin: 0 0 0.4rem;">Express Delivery</h4>
      <p style="font-size: 0.8rem; color: #999; margin: 0; line-height: 1.5;">Free on orders over $100</p>
    </div>
    <div style="text-align: center; padding: 2rem 1.5rem; border-right: 1px solid #f0f0f0;">
      <div style="width: 44px; height: 44px; border: 2px solid #be185d; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem;">
        <div style="width: 14px; height: 8px; border-left: 2px solid #be185d; border-bottom: 2px solid #be185d; transform: rotate(-45deg); margin-bottom: 3px;"></div>
      </div>
      <h4 style="font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #111; margin: 0 0 0.4rem;">Quality Guarantee</h4>
      <p style="font-size: 0.8rem; color: #999; margin: 0; line-height: 1.5;">Handpicked premium materials</p>
    </div>
    <div style="text-align: center; padding: 2rem 1.5rem; border-right: 1px solid #f0f0f0;">
      <div style="width: 44px; height: 44px; border: 2px solid #be185d; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem;">
        <div style="font-size: 1rem; color: #be185d; font-weight: 300;">30</div>
      </div>
      <h4 style="font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #111; margin: 0 0 0.4rem;">Easy Returns</h4>
      <p style="font-size: 0.8rem; color: #999; margin: 0; line-height: 1.5;">30-day no-questions-asked</p>
    </div>
    <div style="text-align: center; padding: 2rem 1.5rem;">
      <div style="width: 44px; height: 44px; border: 2px solid #be185d; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem;">
        <div style="width: 10px; height: 10px; border: 2px solid #be185d; border-radius: 50%;"></div>
      </div>
      <h4 style="font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #111; margin: 0 0 0.4rem;">Secure Checkout</h4>
      <p style="font-size: 0.8rem; color: #999; margin: 0; line-height: 1.5;">256-bit SSL encryption</p>
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
  '
<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1px; background: #e5e5e5; border-radius: 1rem; overflow: hidden;">
  <div style="background: #111; padding: 3rem; display: flex; flex-direction: column; justify-content: space-between; min-height: 280px;">
    <div>
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.3em; color: rgba(255,255,255,0.35); margin: 0 0 1rem; font-weight: 500;">Limited Edition</p>
      <h3 style="font-size: 1.75rem; font-weight: 300; color: #fff; margin: 0 0 0.75rem; line-height: 1.2;">The Midnight<br><span style="font-weight: 700;">Collection</span></h3>
      <p style="font-size: 0.85rem; color: rgba(255,255,255,0.45); margin: 0; line-height: 1.6; font-weight: 300;">Dark tones. Bold silhouettes.<br>Designed for those who own the night.</p>
    </div>
    <a href="/products" style="display: inline-block; color: #fff; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.15em; text-decoration: none; border-bottom: 1px solid rgba(255,255,255,0.3); padding-bottom: 2px; margin-top: 1.5rem; width: fit-content;">Shop Now</a>
  </div>
  <div style="background: #faf9f7; padding: 3rem; display: flex; flex-direction: column; justify-content: space-between; min-height: 280px;">
    <div>
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.3em; color: #be185d; margin: 0 0 1rem; font-weight: 600;">New Arrivals</p>
      <h3 style="font-size: 1.75rem; font-weight: 300; color: #111; margin: 0 0 0.75rem; line-height: 1.2;">Summer<br><span style="font-weight: 700;">Essentials</span></h3>
      <p style="font-size: 0.85rem; color: #999; margin: 0; line-height: 1.6; font-weight: 300;">Lightweight fabrics, effortless style.<br>Your warm-weather wardrobe starts here.</p>
    </div>
    <a href="/products" style="display: inline-block; color: #111; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.15em; text-decoration: none; border-bottom: 1px solid rgba(0,0,0,0.2); padding-bottom: 2px; margin-top: 1.5rem; width: fit-content;">Explore</a>
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
  '
<div style="background: #111; border-radius: 0.5rem; padding: 0.875rem 1.5rem; display: flex; align-items: center; justify-content: center; gap: 0.75rem;">
  <span style="background: #be185d; color: #fff; font-size: 0.6rem; font-weight: 700; padding: 0.2rem 0.5rem; border-radius: 2px; text-transform: uppercase; letter-spacing: 0.08em;">New</span>
  <p style="margin: 0; font-size: 0.8rem; color: rgba(255,255,255,0.7); font-weight: 400;">Use code <span style="color: #fff; font-weight: 600;">THREADS20</span> for 20% off your first order</p>
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
