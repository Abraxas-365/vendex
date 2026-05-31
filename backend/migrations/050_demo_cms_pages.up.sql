-- 050: Seed polished CMS pages for both demo tenants + update fashion branding
-- ============================================================================

-- ── Update fashion branding: black → rose (luxurious fashion feel) ───────────
UPDATE store_branding
SET accent_color = '#be185d',
    updated_at   = NOW()
WHERE tenant_id = 'tnt_fashion';

-- ── Update existing tnt_demo pages with richer content ──────────────────────

UPDATE pages SET html = '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">About Vendex Tech</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">Where technology meets everyday life.</p>

  <p style="line-height:1.8;margin-bottom:1.5rem;">
    Founded in 2024, Vendex Tech started with a simple belief: everyone deserves access to premium technology
    without the premium headache. We curate the best phones, laptops, audio gear, and smart-home devices from
    brands you trust — and back every purchase with our industry-leading 2-year warranty.
  </p>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Our Mission</h2>
  <p style="line-height:1.8;margin-bottom:1.5rem;">
    To make cutting-edge technology accessible, understandable, and delightful for every customer.
    Whether you''re upgrading your first smartphone or building a home studio, our team of tech enthusiasts
    is here to guide you.
  </p>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Why Choose Us</h2>
  <ul style="line-height:2;padding-left:1.25rem;list-style:disc;margin-bottom:1.5rem;">
    <li>Hand-tested products — every item passes our quality review</li>
    <li>2-year warranty on all electronics</li>
    <li>Free shipping on orders over $50</li>
    <li>30-day hassle-free returns</li>
    <li>Expert support from real tech enthusiasts, not scripts</li>
  </ul>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Our Team</h2>
  <p style="line-height:1.8;">
    We''re a small, passionate team of engineers, designers, and product nerds based in New York.
    We test every product we sell — if we wouldn''t buy it ourselves, it doesn''t make the catalog.
  </p>
</div>
', meta = ''{"description":"Learn about Vendex Tech — premium technology, curated for you."}'', updated_at = NOW()
WHERE id = 'pg_about' AND tenant_id = 'tnt_demo';

UPDATE pages SET html = '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">Contact Us</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">We''d love to hear from you.</p>

  <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:2rem;margin-bottom:2.5rem;">
    <div style="background:#f9fafb;border-radius:12px;padding:1.5rem;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Email</h3>
      <p style="color:#6b7280;">hello@vendex.ai</p>
    </div>
    <div style="background:#f9fafb;border-radius:12px;padding:1.5rem;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Phone</h3>
      <p style="color:#6b7280;">+1 (555) 010-0100</p>
      <p style="color:#9ca3af;font-size:0.875rem;">Mon–Fri, 9am–6pm ET</p>
    </div>
    <div style="background:#f9fafb;border-radius:12px;padding:1.5rem;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Office</h3>
      <p style="color:#6b7280;">123 Commerce Ave<br/>New York, NY 10001</p>
    </div>
  </div>

  <h2 style="font-size:1.25rem;font-weight:600;margin-bottom:1rem;">Support Hours</h2>
  <p style="line-height:1.8;color:#6b7280;">
    Our support team is available Monday through Friday, 9:00 AM to 6:00 PM Eastern Time.
    For urgent issues outside business hours, email us and we''ll respond first thing next morning.
  </p>
</div>
', meta = ''{"description":"Get in touch with the Vendex Tech team."}'', updated_at = NOW()
WHERE id = 'pg_contact' AND tenant_id = 'tnt_demo';

UPDATE pages SET html = '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">Frequently Asked Questions</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">Quick answers to common questions.</p>

  <div style="border-top:1px solid #e5e7eb;">
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">How long does shipping take?</h3>
      <p style="color:#6b7280;line-height:1.7;">Standard shipping takes 5–7 business days. Express shipping (2–3 business days) is available at checkout for an additional fee. All orders over $50 ship free.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">What is your return policy?</h3>
      <p style="color:#6b7280;line-height:1.7;">We offer 30-day hassle-free returns on all products. Items must be in original condition with packaging. We''ll provide a prepaid return label.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Do you offer warranties?</h3>
      <p style="color:#6b7280;line-height:1.7;">Yes! Every electronic product comes with our 2-year warranty covering manufacturing defects and hardware failures. This is in addition to the manufacturer''s warranty.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Can I track my order?</h3>
      <p style="color:#6b7280;line-height:1.7;">Absolutely. Once your order ships, you''ll receive a tracking number via email. You can also track your order from your account dashboard.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Do you ship internationally?</h3>
      <p style="color:#6b7280;line-height:1.7;">Currently we ship to the US and Canada. International shipping is coming soon — join our newsletter to be the first to know.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">How can I contact support?</h3>
      <p style="color:#6b7280;line-height:1.7;">Email us at hello@vendex.ai or call +1 (555) 010-0100 during business hours. We typically respond to emails within 4 hours.</p>
    </div>
  </div>
</div>
', meta = ''{"description":"Frequently asked questions about Vendex Tech orders, shipping, and returns."}'', updated_at = NOW()
WHERE id = 'pg_faq' AND tenant_id = 'tnt_demo';

UPDATE pages SET html = '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">Privacy Policy</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">Last updated: January 2026</p>

  <div style="line-height:1.8;color:#374151;">
    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Information We Collect</h2>
    <p style="margin-bottom:1rem;">We collect information you provide directly: name, email, shipping address, and payment details when you place an order. We also collect browsing data (pages visited, device type) to improve your experience.</p>

    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">How We Use Your Information</h2>
    <ul style="padding-left:1.25rem;list-style:disc;margin-bottom:1rem;">
      <li>Process and fulfill your orders</li>
      <li>Send order confirmations and shipping updates</li>
      <li>Improve our website and product offerings</li>
      <li>Respond to your questions and support requests</li>
    </ul>

    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Data Protection</h2>
    <p style="margin-bottom:1rem;">We use industry-standard encryption (TLS 1.3) to protect your data in transit. Payment information is processed by our PCI-compliant payment partner and never stored on our servers.</p>

    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Your Rights</h2>
    <p style="margin-bottom:1rem;">You may request access to, correction of, or deletion of your personal data at any time by contacting us at hello@vendex.ai.</p>

    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Cookies</h2>
    <p style="margin-bottom:1rem;">We use essential cookies to keep your cart and session active. Analytics cookies help us understand how visitors use our site. You can disable non-essential cookies in your browser settings.</p>

    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Contact</h2>
    <p>Questions about this policy? Email us at hello@vendex.ai.</p>
  </div>
</div>
', meta = ''{"description":"Vendex Tech privacy policy — how we collect, use, and protect your data."}'', updated_at = NOW()
WHERE id = 'pg_privacy' AND tenant_id = 'tnt_demo';

-- ── New shipping & returns page for tnt_demo ────────────────────────────────

INSERT INTO pages (id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, content_type, sections)
VALUES ('pg_shipping', 'tnt_demo', 'shipping-returns', 'Shipping & Returns', '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">Shipping & Returns</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">Everything you need to know about getting your order.</p>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Shipping</h2>
  <div style="background:#f9fafb;border-radius:12px;padding:1.5rem;margin-bottom:1.5rem;">
    <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:1.5rem;">
      <div>
        <p style="font-weight:600;">Standard</p>
        <p style="color:#6b7280;font-size:0.875rem;">5–7 business days</p>
        <p style="color:#6b7280;font-size:0.875rem;">Free on orders over $50</p>
      </div>
      <div>
        <p style="font-weight:600;">Express</p>
        <p style="color:#6b7280;font-size:0.875rem;">2–3 business days</p>
        <p style="color:#6b7280;font-size:0.875rem;">$9.99 flat rate</p>
      </div>
      <div>
        <p style="font-weight:600;">Next Day</p>
        <p style="color:#6b7280;font-size:0.875rem;">1 business day</p>
        <p style="color:#6b7280;font-size:0.875rem;">$19.99 flat rate</p>
      </div>
    </div>
  </div>

  <p style="line-height:1.8;color:#374151;margin-bottom:1.5rem;">All orders are processed within 1 business day. You''ll receive a tracking number via email once your package ships. We currently ship to the US and Canada.</p>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Returns</h2>
  <ul style="line-height:2;padding-left:1.25rem;list-style:disc;color:#374151;margin-bottom:1.5rem;">
    <li>30-day return window from delivery date</li>
    <li>Items must be in original condition with packaging</li>
    <li>Free prepaid return labels provided</li>
    <li>Refund processed within 5 business days of receiving the return</li>
    <li>Original shipping costs are non-refundable</li>
  </ul>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Exchanges</h2>
  <p style="line-height:1.8;color:#374151;">Need a different size or color? Contact us within 30 days and we''ll arrange an exchange at no extra cost. If the new item costs more, we''ll charge the difference.</p>
</div>
', '', '{"description":"Vendex Tech shipping rates, delivery times, and return policy."}', 'published', 1, 'system', NOW(), 'html', '[]')
ON CONFLICT (tenant_id, slug) DO NOTHING;

-- ── New tnt_fashion CMS pages ───────────────────────────────────────────────

INSERT INTO pages (id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, content_type, sections)
VALUES ('pg_f_about', 'tnt_fashion', 'about', 'About Urban Threads', '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">About Urban Threads</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">Curated fashion, crafted with intention.</p>

  <p style="line-height:1.8;margin-bottom:1.5rem;">
    Urban Threads was born from a love for timeless style and sustainable craftsmanship. We believe that great fashion
    doesn''t have to come at the expense of the planet — or your wallet. Every piece in our collection is carefully
    selected for its quality, design, and ethical production.
  </p>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Our Philosophy</h2>
  <p style="line-height:1.8;margin-bottom:1.5rem;">
    Less, but better. We''d rather you own five pieces you love than fifty you forget about. Our buyers travel the
    world seeking artisans and small-batch producers who share our commitment to quality and sustainability.
  </p>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Why Urban Threads</h2>
  <ul style="line-height:2;padding-left:1.25rem;list-style:disc;margin-bottom:1.5rem;">
    <li>Ethically sourced materials from certified suppliers</li>
    <li>Handpicked fabrics — cashmere, silk, organic cotton, linen</li>
    <li>Small-batch production to reduce waste</li>
    <li>Free shipping on orders over $100</li>
    <li>14-day free returns — no questions asked</li>
  </ul>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Our Promise</h2>
  <p style="line-height:1.8;">
    Every garment tells a story. From the farms where our cotton grows to the ateliers where our dresses are sewn,
    we maintain full transparency about our supply chain. When you wear Urban Threads, you wear something meaningful.
  </p>
</div>
', '', '{"description":"About Urban Threads — ethically sourced, timeless fashion."}', 'published', 1, 'system', NOW(), 'html', '[]')
ON CONFLICT (tenant_id, slug) DO NOTHING;

INSERT INTO pages (id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, content_type, sections)
VALUES ('pg_f_contact', 'tnt_fashion', 'contact', 'Contact Us', '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">Get In Touch</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">We''re here to help with styling, orders, and everything in between.</p>

  <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(200px,1fr));gap:2rem;margin-bottom:2.5rem;">
    <div style="background:#fdf2f8;border-radius:12px;padding:1.5rem;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Email</h3>
      <p style="color:#6b7280;">hello@urbanthreads.co</p>
    </div>
    <div style="background:#fdf2f8;border-radius:12px;padding:1.5rem;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Phone</h3>
      <p style="color:#6b7280;">+1 (555) 234-5678</p>
      <p style="color:#9ca3af;font-size:0.875rem;">Mon–Sat, 10am–7pm ET</p>
    </div>
    <div style="background:#fdf2f8;border-radius:12px;padding:1.5rem;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Showroom</h3>
      <p style="color:#6b7280;">456 Fashion District<br/>Brooklyn, NY 11201</p>
    </div>
  </div>

  <h2 style="font-size:1.25rem;font-weight:600;margin-bottom:1rem;">Styling Appointments</h2>
  <p style="line-height:1.8;color:#6b7280;">
    Book a complimentary styling session at our Brooklyn showroom. Our stylists will help you find pieces
    that complement your wardrobe and personal style. Email us to schedule.
  </p>
</div>
', '', '{"description":"Contact Urban Threads — styling appointments, orders, and support."}', 'published', 1, 'system', NOW(), 'html', '[]')
ON CONFLICT (tenant_id, slug) DO NOTHING;

INSERT INTO pages (id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, content_type, sections)
VALUES ('pg_f_faq', 'tnt_fashion', 'faq', 'FAQ', '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">Frequently Asked Questions</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">Everything you need to know about shopping with us.</p>

  <div style="border-top:1px solid #e5e7eb;">
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">How do I find my size?</h3>
      <p style="color:#6b7280;line-height:1.7;">Each product page includes a detailed size guide with measurements. If you''re between sizes, we recommend sizing up for a relaxed fit or sizing down for a tailored look.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">What fabrics do you use?</h3>
      <p style="color:#6b7280;line-height:1.7;">We specialize in natural, premium fabrics: Mongolian cashmere, mulberry silk, organic cotton, Belgian linen, and Italian merino wool. Each product description details the exact composition.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">How should I care for my garments?</h3>
      <p style="color:#6b7280;line-height:1.7;">Care instructions are included with every piece. Generally, we recommend cold hand-wash or dry cleaning for silk and cashmere. Cotton and linen can be machine-washed on a gentle cycle.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Can I return or exchange an item?</h3>
      <p style="color:#6b7280;line-height:1.7;">Yes! We offer 14-day free returns. Items must be unworn, unwashed, and in original packaging with tags attached. Exchanges are free — just contact us.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Do you offer gift wrapping?</h3>
      <p style="color:#6b7280;line-height:1.7;">Yes, complimentary gift wrapping is available at checkout. Each order arrives in our signature recycled tissue paper with a hand-written note if you choose.</p>
    </div>
    <div style="padding:1.25rem 0;border-bottom:1px solid #e5e7eb;">
      <h3 style="font-weight:600;margin-bottom:0.5rem;">Are your products sustainable?</h3>
      <p style="color:#6b7280;line-height:1.7;">Sustainability is at our core. We source from certified ethical suppliers, use biodegradable packaging, and produce in small batches to minimize waste. We''re also a member of the Sustainable Apparel Coalition.</p>
    </div>
  </div>
</div>
', '', '{"description":"Frequently asked questions about Urban Threads — sizing, fabrics, care, and returns."}', 'published', 1, 'system', NOW(), 'html', '[]')
ON CONFLICT (tenant_id, slug) DO NOTHING;

INSERT INTO pages (id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, content_type, sections)
VALUES ('pg_f_privacy', 'tnt_fashion', 'privacy', 'Privacy Policy', '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">Privacy Policy</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">Last updated: January 2026</p>

  <div style="line-height:1.8;color:#374151;">
    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Information We Collect</h2>
    <p style="margin-bottom:1rem;">When you shop with Urban Threads, we collect your name, email, shipping address, and payment information to process your orders. We also gather browsing data to personalize your shopping experience.</p>

    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">How We Use Your Data</h2>
    <ul style="padding-left:1.25rem;list-style:disc;margin-bottom:1rem;">
      <li>Process and fulfill your orders</li>
      <li>Send order updates and shipping notifications</li>
      <li>Personalize product recommendations</li>
      <li>Improve our website and collections</li>
    </ul>

    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Data Security</h2>
    <p style="margin-bottom:1rem;">Your data is encrypted with TLS 1.3 in transit. Payment processing is handled by our PCI-compliant partner — we never store your card details.</p>

    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Your Rights</h2>
    <p style="margin-bottom:1rem;">You can request access to, correction of, or deletion of your personal data at any time. Contact us at hello@urbanthreads.co.</p>

    <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Contact</h2>
    <p>Questions? Reach out to hello@urbanthreads.co.</p>
  </div>
</div>
', '', '{"description":"Urban Threads privacy policy — how we handle your data."}', 'published', 1, 'system', NOW(), 'html', '[]')
ON CONFLICT (tenant_id, slug) DO NOTHING;

INSERT INTO pages (id, tenant_id, slug, title, html, css, meta, status, version, created_by, published_at, content_type, sections)
VALUES ('pg_f_shipping', 'tnt_fashion', 'shipping-returns', 'Shipping & Returns', '
<div style="max-width:720px;margin:0 auto;">
  <h1 style="font-size:2rem;font-weight:700;margin-bottom:0.5rem;">Shipping & Returns</h1>
  <p style="color:#6b7280;margin-bottom:2rem;">Delivered with care, returned without hassle.</p>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Shipping</h2>
  <div style="background:#fdf2f8;border-radius:12px;padding:1.5rem;margin-bottom:1.5rem;">
    <div style="display:grid;grid-template-columns:repeat(auto-fit,minmax(180px,1fr));gap:1.5rem;">
      <div>
        <p style="font-weight:600;">Standard</p>
        <p style="color:#6b7280;font-size:0.875rem;">5–7 business days</p>
        <p style="color:#6b7280;font-size:0.875rem;">Free on orders over $100</p>
      </div>
      <div>
        <p style="font-weight:600;">Express</p>
        <p style="color:#6b7280;font-size:0.875rem;">2–3 business days</p>
        <p style="color:#6b7280;font-size:0.875rem;">$12.00 flat rate</p>
      </div>
      <div>
        <p style="font-weight:600;">Next Day</p>
        <p style="color:#6b7280;font-size:0.875rem;">1 business day</p>
        <p style="color:#6b7280;font-size:0.875rem;">$24.00 flat rate</p>
      </div>
    </div>
  </div>

  <p style="line-height:1.8;color:#374151;margin-bottom:1.5rem;">Every order is carefully folded in recycled tissue paper and packaged in our signature boxes. Gift wrapping is complimentary.</p>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Returns</h2>
  <ul style="line-height:2;padding-left:1.25rem;list-style:disc;color:#374151;margin-bottom:1.5rem;">
    <li>14-day return window from delivery date</li>
    <li>Items must be unworn, unwashed, with tags attached</li>
    <li>Free prepaid return labels included in every order</li>
    <li>Refund processed within 3–5 business days</li>
  </ul>

  <h2 style="font-size:1.25rem;font-weight:600;margin:2rem 0 1rem;">Exchanges</h2>
  <p style="line-height:1.8;color:#374151;">Need a different size? We offer free exchanges within 14 days. Contact us and we''ll arrange everything — no extra shipping costs.</p>
</div>
', '', '{"description":"Urban Threads shipping rates, delivery packaging, and return policy."}', 'published', 1, 'system', NOW(), 'html', '[]')
ON CONFLICT (tenant_id, slug) DO NOTHING;
