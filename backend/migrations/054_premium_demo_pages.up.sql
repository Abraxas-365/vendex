-- 054: Upgrade tnt_demo (Vendex Tech) CMS pages + templates to premium content
-- ============================================================================

-- ── About ───────────────────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.3em; color: #6366f1; font-weight: 600; margin: 0 0 2rem; text-align: center;">Our Mission</p>

  <p style="font-size: 1.25rem; color: #111; line-height: 1.8; margin: 0 0 2rem; font-weight: 300; text-align: center;">Technology should empower, not overwhelm. <strong style="font-weight: 600;">We make cutting-edge tech accessible to everyone.</strong></p>

  <div style="width: 40px; height: 1px; background: #ddd; margin: 0 auto 2rem;"></div>

  <p style="font-size: 0.95rem; color: #666; line-height: 1.9; margin: 0 0 1.5rem;">Founded in 2024, Vendex Tech started with a clear vision: bridge the gap between premium technology and the people who use it. We curate the best devices — smartphones, laptops, audio gear, and smart home essentials — so you never have to compromise on quality or value.</p>

  <p style="font-size: 0.95rem; color: #666; line-height: 1.9; margin: 0 0 3rem;">Every product in our catalog is tested by our in-house team. We negotiate directly with manufacturers to bring you competitive pricing, and back every purchase with our industry-leading 2-year warranty.</p>

  <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 0; border-top: 1px solid #eee; border-bottom: 1px solid #eee;">
    <div style="text-align: center; padding: 2rem 1rem; border-right: 1px solid #eee;">
      <p style="font-size: 2rem; font-weight: 200; color: #111; margin: 0; letter-spacing: -0.02em;">200<span style="font-size: 1rem; color: #6366f1;">+</span></p>
      <p style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.15em; color: #999; margin: 0.5rem 0 0;">Products Tested</p>
    </div>
    <div style="text-align: center; padding: 2rem 1rem; border-right: 1px solid #eee;">
      <p style="font-size: 2rem; font-weight: 200; color: #111; margin: 0; letter-spacing: -0.02em;">24<span style="font-size: 1rem; color: #6366f1;">h</span></p>
      <p style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.15em; color: #999; margin: 0.5rem 0 0;">Avg Ship Time</p>
    </div>
    <div style="text-align: center; padding: 2rem 1rem;">
      <p style="font-size: 2rem; font-weight: 200; color: #111; margin: 0; letter-spacing: -0.02em;">2<span style="font-size: 1rem; color: #6366f1;">yr</span></p>
      <p style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.15em; color: #999; margin: 0.5rem 0 0;">Full Warranty</p>
    </div>
  </div>

  <div style="margin-top: 3rem; background: #111; border-radius: 0.5rem; padding: 2.5rem; text-align: center;">
    <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.25em; color: rgba(255,255,255,0.35); margin: 0 0 0.75rem;">Why Vendex</p>
    <p style="font-size: 1rem; color: rgba(255,255,255,0.8); margin: 0; line-height: 1.7; font-weight: 300;">Expert curation. Honest pricing. Real support.<br>Every device we sell, we''d use ourselves.</p>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_demo' AND slug = 'about';

-- ── FAQ ─────────────────────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 1.1rem; color: #666; text-align: center; margin: 0 0 3rem; line-height: 1.7; font-weight: 300;">Quick answers to common questions.</p>

  <div style="display: flex; flex-direction: column; gap: 1px; background: #eee; border-radius: 0.5rem; overflow: hidden;">
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">What warranty do your products carry?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">Every product includes a 2-year Vendex warranty on top of the manufacturer''s coverage. If anything goes wrong, we handle the replacement — no runaround.</p>
    </div>
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">How fast is shipping?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">Standard: 3–5 days. Express: 1–2 days. Next-day available in select cities. Orders over $75 ship free.</p>
    </div>
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">Can I return a product I''ve opened?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">Yes — 30-day returns on all items, opened or not. We''ll send a prepaid label. Refunds process in 3–5 business days.</p>
    </div>
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">Do you price-match?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">Absolutely. Find a lower price at an authorized retailer within 14 days of purchase and we''ll refund the difference, plus an extra 5%.</p>
    </div>
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">Do you ship internationally?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">We currently ship to the US and Canada. International expansion is coming soon — join our mailing list for updates.</p>
    </div>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_demo' AND slug = 'faq';

-- ── Contact ─────────────────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 1.1rem; color: #666; text-align: center; margin: 0 0 3rem; line-height: 1.7; font-weight: 300;">Real humans. Real answers. Average response time under 30 minutes.</p>

  <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 0; border: 1px solid #eee; border-radius: 0.5rem; overflow: hidden; margin-bottom: 3rem;">
    <div style="padding: 2rem 1.5rem; text-align: center; border-right: 1px solid #eee;">
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.2em; color: #999; margin: 0 0 0.5rem;">Email</p>
      <p style="font-size: 0.85rem; color: #111; font-weight: 500; margin: 0;">hello@vendex.ai</p>
    </div>
    <div style="padding: 2rem 1.5rem; text-align: center; border-right: 1px solid #eee;">
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.2em; color: #999; margin: 0 0 0.5rem;">Live Chat</p>
      <p style="font-size: 0.85rem; color: #111; font-weight: 500; margin: 0;">24/7 via dashboard</p>
    </div>
    <div style="padding: 2rem 1.5rem; text-align: center;">
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.2em; color: #999; margin: 0 0 0.5rem;">Phone</p>
      <p style="font-size: 0.85rem; color: #111; font-weight: 500; margin: 0;">+1 (555) 012-3456</p>
    </div>
  </div>

  <div style="border: 1px solid #eee; border-radius: 0.5rem; padding: 2rem;">
    <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.2em; color: #999; margin: 0 0 0.75rem;">Headquarters</p>
    <p style="font-size: 0.95rem; color: #111; font-weight: 500; margin: 0 0 0.25rem;">123 Commerce Avenue, Suite 400</p>
    <p style="font-size: 0.85rem; color: #666; margin: 0;">New York, NY 10001 — Mon–Fri, 9am–6pm ET</p>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_demo' AND slug = 'contact';

-- ── Privacy ─────────────────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 0.75rem; color: #999; text-align: center; margin: 0 0 3rem;">Last updated May 2026</p>

  <div style="display: flex; flex-direction: column; gap: 2.5rem;">
    <div>
      <h3 style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.2em; color: #6366f1; font-weight: 600; margin: 0 0 0.75rem;">Information We Collect</h3>
      <p style="font-size: 0.9rem; color: #555; margin: 0; line-height: 1.8;">We collect only what''s needed to process your orders: name, email, shipping address, and payment details. Browsing analytics help us improve product recommendations — always anonymized.</p>
    </div>
    <div>
      <h3 style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.2em; color: #6366f1; font-weight: 600; margin: 0 0 0.75rem;">How We Use Your Data</h3>
      <p style="font-size: 0.9rem; color: #555; margin: 0; line-height: 1.8;">Order fulfillment, customer support, and service improvement. We never sell personal information to third parties. Marketing emails are strictly opt-in.</p>
    </div>
    <div>
      <h3 style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.2em; color: #6366f1; font-weight: 600; margin: 0 0 0.75rem;">Data Protection</h3>
      <p style="font-size: 0.9rem; color: #555; margin: 0; line-height: 1.8;">TLS 1.3 encryption on all connections. Payment processing through PCI-DSS Level 1 certified partners. Infrastructure hosted in SOC 2 compliant data centers.</p>
    </div>
    <div>
      <h3 style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.2em; color: #6366f1; font-weight: 600; margin: 0 0 0.75rem;">Your Rights</h3>
      <p style="font-size: 0.9rem; color: #555; margin: 0; line-height: 1.8;">Request access, correction, or deletion of your personal data at any time. Contact privacy@vendex.ai — we respond within 48 hours.</p>
    </div>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_demo' AND slug = 'privacy';

-- ── Shipping & Returns ──────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 1.1rem; color: #666; text-align: center; margin: 0 0 3rem; line-height: 1.7; font-weight: 300;">Fast. Free. No-hassle returns.</p>

  <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1px; background: #eee; border-radius: 0.5rem; overflow: hidden; margin-bottom: 3rem;">
    <div style="background: #fff; padding: 2.5rem 2rem;">
      <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.25em; color: #6366f1; font-weight: 600; margin: 0 0 1.5rem;">Shipping</p>
      <div style="display: flex; flex-direction: column; gap: 1rem; font-size: 0.85rem;">
        <div style="display: flex; justify-content: space-between;">
          <span style="color: #666;">Standard — 3–5 days</span>
          <span style="color: #111; font-weight: 600;">$6</span>
        </div>
        <div style="display: flex; justify-content: space-between;">
          <span style="color: #666;">Express — 1–2 days</span>
          <span style="color: #111; font-weight: 600;">$12</span>
        </div>
        <div style="display: flex; justify-content: space-between;">
          <span style="color: #666;">Next Day</span>
          <span style="color: #111; font-weight: 600;">$20</span>
        </div>
        <div style="display: flex; justify-content: space-between; padding-top: 0.75rem; border-top: 1px solid #f0f0f0;">
          <span style="color: #111; font-weight: 500;">Orders over $75</span>
          <span style="color: #6366f1; font-weight: 700;">Free</span>
        </div>
      </div>
    </div>
    <div style="background: #fff; padding: 2.5rem 2rem;">
      <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.25em; color: #6366f1; font-weight: 600; margin: 0 0 1.5rem;">Returns</p>
      <div style="display: flex; flex-direction: column; gap: 0.75rem; font-size: 0.85rem; color: #666;">
        <p style="margin: 0;">30-day return window</p>
        <p style="margin: 0;">Opened products accepted</p>
        <p style="margin: 0;">Free prepaid return label</p>
        <p style="margin: 0;">Refund in 3–5 business days</p>
      </div>
    </div>
  </div>

  <div style="background: #111; border-radius: 0.5rem; padding: 1.25rem 2rem; text-align: center;">
    <p style="margin: 0; font-size: 0.8rem; color: rgba(255,255,255,0.6);">Need help? <span style="color: #fff; font-weight: 500;">hello@vendex.ai</span></p>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_demo' AND slug = 'shipping-returns';

-- ── Upgrade template pages ──────────────────────────────────────────────────

-- PLP template — premium version
UPDATE pages SET html = '
<div style="background: #111; color: #fff; padding: 4rem 3rem; border-radius: 0; margin: -2rem -1rem 2rem; position: relative; overflow: hidden;">
  <div style="position: absolute; top: 0; right: 0; width: 40%; height: 100%; background: linear-gradient(135deg, transparent 0%, rgba(99,102,241,0.08) 100%);"></div>
  <div style="position: relative; max-width: 500px;">
    <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.3em; color: rgba(255,255,255,0.4); margin: 0 0 1rem; font-weight: 500;">Catalog — 2026</p>
    <h2 style="font-size: 2.5rem; font-weight: 300; margin: 0 0 1rem; letter-spacing: -0.03em; line-height: 1.15;">Tech <span style="font-weight: 700;">Essentials</span></h2>
    <p style="font-size: 0.95rem; color: rgba(255,255,255,0.5); margin: 0; line-height: 1.7; font-weight: 300;">Every device tested, reviewed, and backed by our 2-year warranty. Premium tech, honest pricing.</p>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_demo' AND slug = '_plp';

-- PDP template — premium version
UPDATE pages SET html = '
<div style="margin-top: 3rem; padding-top: 2.5rem; border-top: 1px solid #f0f0f0;">
  <div style="display: grid; grid-template-columns: repeat(4, 1fr); gap: 0;">
    <div style="text-align: center; padding: 2rem 1.5rem; border-right: 1px solid #f0f0f0;">
      <div style="width: 44px; height: 44px; border: 2px solid #6366f1; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem;">
        <div style="width: 0; height: 0; border-left: 8px solid #6366f1; border-top: 5px solid transparent; border-bottom: 5px solid transparent; margin-left: 3px;"></div>
      </div>
      <h4 style="font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #111; margin: 0 0 0.4rem;">Fast Shipping</h4>
      <p style="font-size: 0.8rem; color: #999; margin: 0; line-height: 1.5;">Free on orders over $75</p>
    </div>
    <div style="text-align: center; padding: 2rem 1.5rem; border-right: 1px solid #f0f0f0;">
      <div style="width: 44px; height: 44px; border: 2px solid #6366f1; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem;">
        <div style="width: 14px; height: 8px; border-left: 2px solid #6366f1; border-bottom: 2px solid #6366f1; transform: rotate(-45deg); margin-bottom: 3px;"></div>
      </div>
      <h4 style="font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #111; margin: 0 0 0.4rem;">2-Year Warranty</h4>
      <p style="font-size: 0.8rem; color: #999; margin: 0; line-height: 1.5;">Full coverage, no fine print</p>
    </div>
    <div style="text-align: center; padding: 2rem 1.5rem; border-right: 1px solid #f0f0f0;">
      <div style="width: 44px; height: 44px; border: 2px solid #6366f1; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem;">
        <div style="font-size: 1rem; color: #6366f1; font-weight: 300;">30</div>
      </div>
      <h4 style="font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #111; margin: 0 0 0.4rem;">Easy Returns</h4>
      <p style="font-size: 0.8rem; color: #999; margin: 0; line-height: 1.5;">Even on opened products</p>
    </div>
    <div style="text-align: center; padding: 2rem 1.5rem;">
      <div style="width: 44px; height: 44px; border: 2px solid #6366f1; border-radius: 50%; display: flex; align-items: center; justify-content: center; margin: 0 auto 1rem;">
        <div style="width: 10px; height: 10px; border: 2px solid #6366f1; border-radius: 50%;"></div>
      </div>
      <h4 style="font-size: 0.8rem; font-weight: 700; text-transform: uppercase; letter-spacing: 0.1em; color: #111; margin: 0 0 0.4rem;">Price Match</h4>
      <p style="font-size: 0.8rem; color: #999; margin: 0; line-height: 1.5;">Best price guaranteed + 5%</p>
    </div>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_demo' AND slug = '_pdp';

-- Home template — premium version
UPDATE pages SET html = '
<div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1px; background: #e5e5e5; border-radius: 1rem; overflow: hidden;">
  <div style="background: #111; padding: 3rem; display: flex; flex-direction: column; justify-content: space-between; min-height: 280px;">
    <div>
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.3em; color: rgba(255,255,255,0.35); margin: 0 0 1rem; font-weight: 500;">Featured</p>
      <h3 style="font-size: 1.75rem; font-weight: 300; color: #fff; margin: 0 0 0.75rem; line-height: 1.2;">Next-Gen<br><span style="font-weight: 700;">Audio</span></h3>
      <p style="font-size: 0.85rem; color: rgba(255,255,255,0.45); margin: 0; line-height: 1.6; font-weight: 300;">Studio-grade headphones and speakers.<br>Hear every detail, feel every beat.</p>
    </div>
    <a href="/products" style="display: inline-block; color: #fff; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.15em; text-decoration: none; border-bottom: 1px solid rgba(255,255,255,0.3); padding-bottom: 2px; margin-top: 1.5rem; width: fit-content;">Shop Now</a>
  </div>
  <div style="background: #f8f8ff; padding: 3rem; display: flex; flex-direction: column; justify-content: space-between; min-height: 280px;">
    <div>
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.3em; color: #6366f1; margin: 0 0 1rem; font-weight: 600;">This Week</p>
      <h3 style="font-size: 1.75rem; font-weight: 300; color: #111; margin: 0 0 0.75rem; line-height: 1.2;">Smart Home<br><span style="font-weight: 700;">Starter Kits</span></h3>
      <p style="font-size: 0.85rem; color: #999; margin: 0; line-height: 1.6; font-weight: 300;">Everything you need to get started.<br>Up to 30% off curated bundles.</p>
    </div>
    <a href="/products" style="display: inline-block; color: #111; font-size: 0.75rem; text-transform: uppercase; letter-spacing: 0.15em; text-decoration: none; border-bottom: 1px solid rgba(0,0,0,0.2); padding-bottom: 2px; margin-top: 1.5rem; width: fit-content;">Explore</a>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_demo' AND slug = '_home';

-- Cart template — premium version
UPDATE pages SET html = '
<div style="background: #111; border-radius: 0.5rem; padding: 0.875rem 1.5rem; display: flex; align-items: center; justify-content: center; gap: 0.75rem;">
  <span style="background: #6366f1; color: #fff; font-size: 0.6rem; font-weight: 700; padding: 0.2rem 0.5rem; border-radius: 2px; text-transform: uppercase; letter-spacing: 0.08em;">Deal</span>
  <p style="margin: 0; font-size: 0.8rem; color: rgba(255,255,255,0.7); font-weight: 400;">Free shipping on orders over $75 — use code <span style="color: #fff; font-weight: 600;">VENDEX10</span> for 10% off accessories</p>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_demo' AND slug = '_cart';
