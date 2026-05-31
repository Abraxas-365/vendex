-- Template pages for PLP and PDP — editable via CMS / AI workspace agent.
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
