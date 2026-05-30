-- Block types registry (global, not per-tenant)
CREATE TABLE IF NOT EXISTS block_types (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    display_name VARCHAR(255) NOT NULL,
    category VARCHAR(50) NOT NULL CHECK (category IN ('content', 'commerce', 'media', 'layout')),
    schema JSONB NOT NULL DEFAULT '{}',
    default_settings JSONB NOT NULL DEFAULT '{}',
    icon VARCHAR(50) NOT NULL DEFAULT 'box',
    plugin_id VARCHAR(36) REFERENCES plugins(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Add sections column to existing pages table
ALTER TABLE pages ADD COLUMN IF NOT EXISTS content_type VARCHAR(20) NOT NULL DEFAULT 'html' CHECK (content_type IN ('html', 'blocks'));
ALTER TABLE pages ADD COLUMN IF NOT EXISTS sections JSONB NOT NULL DEFAULT '[]';

-- Seed built-in block types
INSERT INTO block_types (id, name, display_name, category, schema, default_settings, icon) VALUES
('bt-hero', 'hero', 'Hero Banner', 'content',
 '{"properties":{"heading":{"type":"string"},"subheading":{"type":"string"},"background_image":{"type":"string"},"background_color":{"type":"string"},"cta_text":{"type":"string"},"cta_link":{"type":"string"},"alignment":{"type":"string","enum":["left","center","right"]}}}',
 '{"heading":"Welcome","subheading":"","background_color":"#4F46E5","alignment":"center"}',
 'image'),

('bt-rich-text', 'rich-text', 'Rich Text', 'content',
 '{"properties":{"content":{"type":"string"}}}',
 '{"content":""}',
 'type'),

('bt-product-grid', 'product-grid', 'Product Grid', 'commerce',
 '{"properties":{"title":{"type":"string"},"collection_id":{"type":"string"},"columns":{"type":"integer","minimum":1,"maximum":6},"limit":{"type":"integer","minimum":1,"maximum":24},"show_price":{"type":"boolean"}}}',
 '{"title":"Featured Products","columns":4,"limit":8,"show_price":true}',
 'grid'),

('bt-image-banner', 'image-banner', 'Image Banner', 'media',
 '{"properties":{"image_url":{"type":"string"},"alt_text":{"type":"string"},"link":{"type":"string"},"height":{"type":"string"}}}',
 '{"height":"400px"}',
 'image'),

('bt-newsletter', 'newsletter', 'Newsletter Signup', 'content',
 '{"properties":{"heading":{"type":"string"},"description":{"type":"string"},"button_text":{"type":"string"},"placeholder":{"type":"string"}}}',
 '{"heading":"Subscribe to our newsletter","button_text":"Subscribe","placeholder":"Enter your email"}',
 'mail'),

('bt-category-grid', 'category-grid', 'Category Grid', 'commerce',
 '{"properties":{"title":{"type":"string"},"columns":{"type":"integer","minimum":1,"maximum":6},"show_description":{"type":"boolean"}}}',
 '{"title":"Shop by Category","columns":3,"show_description":false}',
 'grid'),

('bt-cta', 'cta', 'Call to Action', 'content',
 '{"properties":{"heading":{"type":"string"},"description":{"type":"string"},"button_text":{"type":"string"},"button_link":{"type":"string"},"background_color":{"type":"string"}}}',
 '{"heading":"Ready to get started?","button_text":"Get Started","background_color":"#111827"}',
 'megaphone'),

('bt-spacer', 'spacer', 'Spacer', 'layout',
 '{"properties":{"height":{"type":"string"}}}',
 '{"height":"48px"}',
 'minus'),

('bt-divider', 'divider', 'Divider', 'layout',
 '{"properties":{"color":{"type":"string"},"thickness":{"type":"string"},"width":{"type":"string"}}}',
 '{"color":"#E5E7EB","thickness":"1px","width":"100%"}',
 'minus'),

('bt-testimonials', 'testimonials', 'Testimonials', 'content',
 '{"properties":{"title":{"type":"string"},"items":{"type":"array","items":{"type":"object","properties":{"quote":{"type":"string"},"author":{"type":"string"},"role":{"type":"string"},"avatar":{"type":"string"}}}}}}',
 '{"title":"What our customers say","items":[]}',
 'message-circle')

ON CONFLICT (id) DO NOTHING;
