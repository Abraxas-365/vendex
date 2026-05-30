-- Seed built-in block types with rich schemas and default settings.
-- Uses ON CONFLICT (name) DO NOTHING for idempotency.
-- Note: 'marketing' is not a valid category per the CHECK constraint;
--       newsletter and cta are categorised as 'content' instead.

INSERT INTO block_types (id, name, display_name, category, icon, schema, default_settings, created_at, updated_at)
VALUES

-- 1. hero
(
  '00000000-0000-0000-0000-000000000001',
  'hero',
  'Hero Banner',
  'layout',
  'Star',
  '{
    "heading":          {"type": "string",  "label": "Heading",            "default": "Welcome to our store"},
    "subheading":       {"type": "string",  "label": "Subheading",         "default": ""},
    "background_image": {"type": "string",  "label": "Background Image",   "format": "url"},
    "cta_text":         {"type": "string",  "label": "CTA Text",           "default": ""},
    "cta_link":         {"type": "string",  "label": "CTA Link",           "format": "url"},
    "alignment":        {"type": "enum",    "label": "Alignment",          "options": ["left", "center", "right"], "default": "center"},
    "overlay_opacity":  {"type": "number",  "label": "Overlay Opacity",    "min": 0, "max": 100, "default": 40}
  }',
  '{
    "heading":         "Welcome to our store",
    "alignment":       "center",
    "overlay_opacity": 40
  }',
  NOW(),
  NOW()
),

-- 2. rich-text
(
  '00000000-0000-0000-0000-000000000002',
  'rich-text',
  'Rich Text',
  'content',
  'FileText',
  '{
    "content":   {"type": "string", "label": "Content",   "format": "html", "default": ""},
    "max_width": {"type": "enum",   "label": "Max Width", "options": ["sm", "md", "lg", "full"], "default": "md"}
  }',
  '{
    "content":   "",
    "max_width": "md"
  }',
  NOW(),
  NOW()
),

-- 3. product-grid
(
  '00000000-0000-0000-0000-000000000003',
  'product-grid',
  'Product Grid',
  'commerce',
  'Grid',
  '{
    "title":         {"type": "string",  "label": "Title",            "default": ""},
    "collection_id": {"type": "string",  "label": "Collection ID"},
    "columns":       {"type": "number",  "label": "Columns",          "min": 2, "max": 6,  "default": 4},
    "max_items":     {"type": "number",  "label": "Max Items",        "min": 1, "max": 100, "default": 8},
    "show_price":    {"type": "boolean", "label": "Show Price",       "default": true},
    "show_rating":   {"type": "boolean", "label": "Show Rating",      "default": true}
  }',
  '{
    "columns":     4,
    "max_items":   8,
    "show_price":  true,
    "show_rating": true
  }',
  NOW(),
  NOW()
),

-- 4. featured-collection
(
  '00000000-0000-0000-0000-000000000004',
  'featured-collection',
  'Featured Collection',
  'commerce',
  'ShoppingBag',
  '{
    "title":         {"type": "string", "label": "Title",          "default": ""},
    "collection_id": {"type": "string", "label": "Collection ID"},
    "layout":        {"type": "enum",   "label": "Layout",         "options": ["grid", "carousel"], "default": "grid"},
    "items_to_show": {"type": "number", "label": "Items to Show",  "min": 1, "max": 24, "default": 4}
  }',
  '{
    "layout":        "grid",
    "items_to_show": 4
  }',
  NOW(),
  NOW()
),

-- 5. image
(
  '00000000-0000-0000-0000-000000000005',
  'image',
  'Image',
  'media',
  'Image',
  '{
    "src":     {"type": "string", "label": "Image URL", "format": "url"},
    "alt":     {"type": "string", "label": "Alt Text",  "default": ""},
    "link":    {"type": "string", "label": "Link URL",  "format": "url"},
    "width":   {"type": "enum",   "label": "Width",     "options": ["sm", "md", "lg", "full"], "default": "full"},
    "caption": {"type": "string", "label": "Caption",   "default": ""}
  }',
  '{
    "width": "full"
  }',
  NOW(),
  NOW()
),

-- 6. video
(
  '00000000-0000-0000-0000-000000000006',
  'video',
  'Video',
  'media',
  'Play',
  '{
    "url":      {"type": "string",  "label": "Video URL",   "format": "url"},
    "autoplay": {"type": "boolean", "label": "Autoplay",    "default": false},
    "loop":     {"type": "boolean", "label": "Loop",        "default": false},
    "muted":    {"type": "boolean", "label": "Muted",       "default": true},
    "poster":   {"type": "string",  "label": "Poster Image","format": "url"}
  }',
  '{
    "autoplay": false,
    "loop":     false,
    "muted":    true
  }',
  NOW(),
  NOW()
),

-- 7. banner
(
  '00000000-0000-0000-0000-000000000007',
  'banner',
  'Banner',
  'layout',
  'Flag',
  '{
    "text":             {"type": "string",  "label": "Text",             "default": ""},
    "background_color": {"type": "string",  "label": "Background Color", "format": "color", "default": "#4f46e5"},
    "text_color":       {"type": "string",  "label": "Text Color",       "format": "color", "default": "#ffffff"},
    "link":             {"type": "string",  "label": "Link URL",         "format": "url"},
    "dismissible":      {"type": "boolean", "label": "Dismissible",      "default": false}
  }',
  '{
    "background_color": "#4f46e5",
    "text_color":       "#ffffff",
    "dismissible":      false
  }',
  NOW(),
  NOW()
),

-- 8. newsletter
(
  '00000000-0000-0000-0000-000000000008',
  'newsletter',
  'Newsletter Signup',
  'content',
  'Mail',
  '{
    "heading":         {"type": "string", "label": "Heading",         "default": "Subscribe to our newsletter"},
    "description":     {"type": "string", "label": "Description",     "default": ""},
    "placeholder":     {"type": "string", "label": "Placeholder",     "default": "Enter your email"},
    "button_text":     {"type": "string", "label": "Button Text",     "default": "Subscribe"},
    "success_message": {"type": "string", "label": "Success Message", "default": "Thanks for subscribing!"}
  }',
  '{
    "heading":         "Subscribe to our newsletter",
    "button_text":     "Subscribe",
    "placeholder":     "Enter your email",
    "success_message": "Thanks for subscribing!"
  }',
  NOW(),
  NOW()
),

-- 9. testimonials
(
  '00000000-0000-0000-0000-000000000009',
  'testimonials',
  'Testimonials',
  'content',
  'MessageCircle',
  '{
    "heading": {"type": "string", "label": "Heading", "default": "What our customers say"},
    "items":   {
      "type":  "array",
      "label": "Testimonials",
      "items": {
        "type": "object",
        "properties": {
          "name":       {"type": "string", "label": "Name"},
          "role":       {"type": "string", "label": "Role"},
          "quote":      {"type": "string", "label": "Quote"},
          "avatar_url": {"type": "string", "label": "Avatar URL", "format": "url"}
        }
      }
    },
    "layout":  {"type": "enum",   "label": "Layout",  "options": ["grid", "carousel"], "default": "grid"},
    "columns": {"type": "number", "label": "Columns", "min": 1, "max": 4, "default": 3}
  }',
  '{
    "heading": "What our customers say",
    "layout":  "grid",
    "columns": 3,
    "items":   []
  }',
  NOW(),
  NOW()
),

-- 10. cta (call to action)
(
  '00000000-0000-0000-0000-000000000010',
  'cta',
  'Call to Action',
  'content',
  'Megaphone',
  '{
    "heading":                {"type": "string", "label": "Heading",                "default": "Ready to get started?"},
    "description":            {"type": "string", "label": "Description",            "default": ""},
    "primary_button_text":    {"type": "string", "label": "Primary Button Text",    "default": "Shop Now"},
    "primary_button_link":    {"type": "string", "label": "Primary Button Link",    "format": "url"},
    "secondary_button_text":  {"type": "string", "label": "Secondary Button Text",  "default": ""},
    "secondary_button_link":  {"type": "string", "label": "Secondary Button Link",  "format": "url"},
    "background_color":       {"type": "string", "label": "Background Color",       "format": "color", "default": "#4f46e5"}
  }',
  '{
    "heading":             "Ready to get started?",
    "primary_button_text": "Shop Now",
    "background_color":    "#4f46e5"
  }',
  NOW(),
  NOW()
)

ON CONFLICT (name) DO NOTHING;
