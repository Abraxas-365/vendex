-- 053: Upgrade tnt_fashion regular CMS pages to premium editorial content
-- ============================================================================

-- ── About ───────────────────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.3em; color: #be185d; font-weight: 600; margin: 0 0 2rem; text-align: center;">Our Story</p>

  <p style="font-size: 1.25rem; color: #111; line-height: 1.8; margin: 0 0 2rem; font-weight: 300; text-align: center;">Urban Threads was born from a simple belief: <strong style="font-weight: 600;">exceptional style shouldn''t require compromise.</strong></p>

  <div style="width: 40px; height: 1px; background: #ddd; margin: 0 auto 2rem;"></div>

  <p style="font-size: 0.95rem; color: #666; line-height: 1.9; margin: 0 0 1.5rem;">We curate pieces that blend contemporary design with timeless craftsmanship. From the silk midi dresses perfect for evening events to the cashmere knits that elevate your everyday — each item is chosen with intent.</p>

  <p style="font-size: 0.95rem; color: #666; line-height: 1.9; margin: 0 0 3rem;">We partner with independent ateliers and designers who share our commitment to quality materials, ethical production, and designs that transcend seasons.</p>

  <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 0; border-top: 1px solid #eee; border-bottom: 1px solid #eee;">
    <div style="text-align: center; padding: 2rem 1rem; border-right: 1px solid #eee;">
      <p style="font-size: 2rem; font-weight: 200; color: #111; margin: 0; letter-spacing: -0.02em;">50<span style="font-size: 1rem; color: #be185d;">+</span></p>
      <p style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.15em; color: #999; margin: 0.5rem 0 0;">Designer Partners</p>
    </div>
    <div style="text-align: center; padding: 2rem 1rem; border-right: 1px solid #eee;">
      <p style="font-size: 2rem; font-weight: 200; color: #111; margin: 0; letter-spacing: -0.02em;">12<span style="font-size: 1rem; color: #be185d;">k</span></p>
      <p style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.15em; color: #999; margin: 0.5rem 0 0;">Happy Customers</p>
    </div>
    <div style="text-align: center; padding: 2rem 1rem;">
      <p style="font-size: 2rem; font-weight: 200; color: #111; margin: 0; letter-spacing: -0.02em;">98<span style="font-size: 1rem; color: #be185d;">%</span></p>
      <p style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.15em; color: #999; margin: 0.5rem 0 0;">Satisfaction Rate</p>
    </div>
  </div>

  <div style="margin-top: 3rem; background: #111; border-radius: 0.5rem; padding: 2.5rem; text-align: center;">
    <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.25em; color: rgba(255,255,255,0.35); margin: 0 0 0.75rem;">Our Commitment</p>
    <p style="font-size: 1rem; color: rgba(255,255,255,0.8); margin: 0; line-height: 1.7; font-weight: 300;">Sustainable sourcing. Ethical production. Timeless design.<br>Every purchase supports independent artisans worldwide.</p>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_fashion' AND slug = 'about';

-- ── FAQ ─────────────────────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 1.1rem; color: #666; text-align: center; margin: 0 0 3rem; line-height: 1.7; font-weight: 300;">Everything you need to know about shopping with us.</p>

  <div style="display: flex; flex-direction: column; gap: 1px; background: #eee; border-radius: 0.5rem; overflow: hidden;">
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">How do I find my size?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">Each product page includes a detailed size guide. Between sizes? Size up for relaxed, down for tailored. Our stylists are available via chat to help.</p>
    </div>
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">What is your return policy?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">30 days, unworn with original tags. Initiate from your account — we''ll send a prepaid label. Refunds in 3–5 business days.</p>
    </div>
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">How long does shipping take?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">Standard: 3–5 days. Express: 1–2 days ($12, or free over $200). All orders include tracking.</p>
    </div>
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">Are your products sustainably made?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">Yes. Certified ethical manufacturers, eco-friendly packaging, carbon offset. Look for the Conscious Choice badge.</p>
    </div>
    <div style="background: #fff; padding: 1.75rem 2rem;">
      <h4 style="font-weight: 600; color: #111; margin: 0 0 0.5rem; font-size: 0.9rem;">Do you offer gift cards?</h4>
      <p style="color: #888; font-size: 0.85rem; margin: 0; line-height: 1.7;">Digital gift cards in $50, $100, $250, and $500. Never expire, usable on any item.</p>
    </div>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_fashion' AND slug = 'faq';

-- ── Contact ─────────────────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 1.1rem; color: #666; text-align: center; margin: 0 0 3rem; line-height: 1.7; font-weight: 300;">We''d love to hear from you. Our team responds within 2 hours during business days.</p>

  <div style="display: grid; grid-template-columns: repeat(3, 1fr); gap: 0; border: 1px solid #eee; border-radius: 0.5rem; overflow: hidden; margin-bottom: 3rem;">
    <div style="padding: 2rem 1.5rem; text-align: center; border-right: 1px solid #eee;">
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.2em; color: #999; margin: 0 0 0.5rem;">Email</p>
      <p style="font-size: 0.85rem; color: #111; font-weight: 500; margin: 0;">hello@urbanthreads.co</p>
    </div>
    <div style="padding: 2rem 1.5rem; text-align: center; border-right: 1px solid #eee;">
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.2em; color: #999; margin: 0 0 0.5rem;">Live Chat</p>
      <p style="font-size: 0.85rem; color: #111; font-weight: 500; margin: 0;">Mon–Fri, 9–6 EST</p>
    </div>
    <div style="padding: 2rem 1.5rem; text-align: center;">
      <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.2em; color: #999; margin: 0 0 0.5rem;">Phone</p>
      <p style="font-size: 0.85rem; color: #111; font-weight: 500; margin: 0;">+1 (555) 234-5678</p>
    </div>
  </div>

  <div style="border: 1px solid #eee; border-radius: 0.5rem; padding: 2rem;">
    <p style="font-size: 0.6rem; text-transform: uppercase; letter-spacing: 0.2em; color: #999; margin: 0 0 0.75rem;">Showroom</p>
    <p style="font-size: 0.95rem; color: #111; font-weight: 500; margin: 0 0 0.25rem;">123 Fashion Avenue, SoHo</p>
    <p style="font-size: 0.85rem; color: #666; margin: 0;">New York, NY 10012 — Tue–Sat, 10am–7pm</p>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_fashion' AND slug = 'contact';

-- ── Privacy ─────────────────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 0.75rem; color: #999; text-align: center; margin: 0 0 3rem;">Last updated May 2026</p>

  <div style="display: flex; flex-direction: column; gap: 2.5rem;">
    <div>
      <h3 style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.2em; color: #be185d; font-weight: 600; margin: 0 0 0.75rem;">Information We Collect</h3>
      <p style="font-size: 0.9rem; color: #555; margin: 0; line-height: 1.8;">We collect information you provide directly — name, email, shipping address, and payment details when you make a purchase. We also collect browsing data to improve your experience.</p>
    </div>
    <div>
      <h3 style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.2em; color: #be185d; font-weight: 600; margin: 0 0 0.75rem;">How We Use Your Data</h3>
      <p style="font-size: 0.9rem; color: #555; margin: 0; line-height: 1.8;">Exclusively to process orders, provide support, and improve our services. We never sell your personal information. Marketing emails are opt-in only.</p>
    </div>
    <div>
      <h3 style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.2em; color: #be185d; font-weight: 600; margin: 0 0 0.75rem;">Data Protection</h3>
      <p style="font-size: 0.9rem; color: #555; margin: 0; line-height: 1.8;">All transactions encrypted with TLS 1.3. Payment data processed by PCI-DSS Level 1 certified providers. Regular security audits.</p>
    </div>
    <div>
      <h3 style="font-size: 0.7rem; text-transform: uppercase; letter-spacing: 0.2em; color: #be185d; font-weight: 600; margin: 0 0 0.75rem;">Your Rights</h3>
      <p style="font-size: 0.9rem; color: #555; margin: 0; line-height: 1.8;">Access, modify, or delete your personal data at any time. Contact privacy@urbanthreads.co with any questions.</p>
    </div>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_fashion' AND slug = 'privacy';

-- ── Shipping & Returns ──────────────────────────────────────────────────────
UPDATE pages SET html = '
<div style="max-width: 640px; margin: 0 auto;">
  <p style="font-size: 1.1rem; color: #666; text-align: center; margin: 0 0 3rem; line-height: 1.7; font-weight: 300;">Simple. Transparent. No surprises.</p>

  <div style="display: grid; grid-template-columns: 1fr 1fr; gap: 1px; background: #eee; border-radius: 0.5rem; overflow: hidden; margin-bottom: 3rem;">
    <div style="background: #fff; padding: 2.5rem 2rem;">
      <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.25em; color: #be185d; font-weight: 600; margin: 0 0 1.5rem;">Shipping</p>
      <div style="display: flex; flex-direction: column; gap: 1rem; font-size: 0.85rem;">
        <div style="display: flex; justify-content: space-between;">
          <span style="color: #666;">Standard — 3–5 days</span>
          <span style="color: #111; font-weight: 600;">$8</span>
        </div>
        <div style="display: flex; justify-content: space-between;">
          <span style="color: #666;">Express — 1–2 days</span>
          <span style="color: #111; font-weight: 600;">$15</span>
        </div>
        <div style="display: flex; justify-content: space-between; padding-top: 0.75rem; border-top: 1px solid #f0f0f0;">
          <span style="color: #111; font-weight: 500;">Orders over $100</span>
          <span style="color: #be185d; font-weight: 700;">Free</span>
        </div>
      </div>
    </div>
    <div style="background: #fff; padding: 2.5rem 2rem;">
      <p style="font-size: 0.65rem; text-transform: uppercase; letter-spacing: 0.25em; color: #be185d; font-weight: 600; margin: 0 0 1.5rem;">Returns</p>
      <div style="display: flex; flex-direction: column; gap: 0.75rem; font-size: 0.85rem; color: #666;">
        <p style="margin: 0;">30-day return window</p>
        <p style="margin: 0;">Free prepaid return label</p>
        <p style="margin: 0;">Refund in 3–5 business days</p>
        <p style="margin: 0;">Exchange for size or color</p>
      </div>
    </div>
  </div>

  <div style="background: #111; border-radius: 0.5rem; padding: 1.25rem 2rem; text-align: center;">
    <p style="margin: 0; font-size: 0.8rem; color: rgba(255,255,255,0.6);">Questions? <span style="color: #fff; font-weight: 500;">hello@urbanthreads.co</span></p>
  </div>
</div>',
updated_at = NOW()
WHERE tenant_id = 'tnt_fashion' AND slug = 'shipping-returns';
