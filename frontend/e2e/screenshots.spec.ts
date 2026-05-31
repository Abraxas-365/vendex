import { test } from '@playwright/test'
import * as fs from 'fs'
import * as path from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const SCREENSHOT_DIR = path.join(__dirname, '..', 'screenshots')

// All pages to screenshot
const pages = [
  // Store pages
  { name: 'store-home', path: '/', title: 'Store Home' },
  { name: 'store-products', path: '/products', title: 'Product Catalog' },
  { name: 'store-cart', path: '/cart', title: 'Shopping Cart' },
  { name: 'store-checkout', path: '/checkout', title: 'Checkout' },

  // Auth
  { name: 'auth-login', path: '/login', title: 'Login' },

  // Admin pages
  { name: 'admin-dashboard', path: '/admin', title: 'Dashboard' },
  { name: 'admin-products', path: '/admin/products', title: 'Products' },
  { name: 'admin-orders', path: '/admin/orders', title: 'Orders' },
  { name: 'admin-customers', path: '/admin/customers', title: 'Customers' },
  { name: 'admin-catalog', path: '/admin/catalog', title: 'Catalog' },
  { name: 'admin-collections', path: '/admin/collections', title: 'Collections' },
  { name: 'admin-pages', path: '/admin/pages', title: 'Pages' },
  { name: 'admin-promos', path: '/admin/promos', title: 'Promotions' },
  { name: 'admin-media', path: '/admin/media', title: 'Media Library' },
  { name: 'admin-marketplace', path: '/admin/marketplace', title: 'Marketplace' },
  { name: 'admin-theme', path: '/admin/theme', title: 'Theme Editor' },
  { name: 'admin-settings', path: '/admin/settings', title: 'Settings' },
  { name: 'admin-agent', path: '/admin/agent', title: 'AI Agent' },

  // Fulfillment & Logistics
  { name: 'admin-shipping', path: '/admin/shipping', title: 'Shipping' },
  { name: 'admin-tax', path: '/admin/tax', title: 'Tax' },
  { name: 'admin-inventory', path: '/admin/inventory', title: 'Inventory' },
  { name: 'admin-returns', path: '/admin/returns', title: 'Returns' },

  // Payments & Finance
  { name: 'admin-payments', path: '/admin/payments', title: 'Payments' },
  { name: 'admin-gift-cards', path: '/admin/gift-cards', title: 'Gift Cards' },
  { name: 'admin-currency', path: '/admin/currency-rates', title: 'Currency Rates' },

  // Customers & Loyalty
  { name: 'admin-customer-groups', path: '/admin/customer-groups', title: 'Customer Groups' },
  { name: 'admin-loyalty', path: '/admin/loyalty', title: 'Loyalty' },
  { name: 'admin-subscriptions', path: '/admin/subscriptions', title: 'Subscriptions' },

  // Marketing & Growth
  { name: 'admin-cart-recovery', path: '/admin/cart-recovery', title: 'Cart Recovery' },
  { name: 'admin-reviews', path: '/admin/reviews', title: 'Reviews' },
  { name: 'admin-bundles', path: '/admin/bundles', title: 'Bundles' },
  { name: 'admin-ab-testing', path: '/admin/ab-testing', title: 'A/B Testing' },
  { name: 'admin-recommendations', path: '/admin/recommendations', title: 'Recommendations' },

  // Content & CMS
  { name: 'admin-blog', path: '/admin/blog', title: 'Blog' },
  { name: 'admin-translations', path: '/admin/translations', title: 'Translations' },

  // Analytics & Reporting
  { name: 'admin-reporting', path: '/admin/reporting', title: 'Reporting' },

  // Operations & Tech
  { name: 'admin-import-export', path: '/admin/import-export', title: 'Import / Export' },
  { name: 'admin-webhooks', path: '/admin/webhooks', title: 'Webhooks' },
  { name: 'admin-audit', path: '/admin/audit-logs', title: 'Audit Logs' },
  { name: 'admin-bulk-operations', path: '/admin/bulk-operations', title: 'Bulk Operations' },

  // Multi-channel & Social
  { name: 'admin-multistores', path: '/admin/storefronts', title: 'Storefronts' },
  { name: 'admin-social-accounts', path: '/admin/social-accounts', title: 'Social Accounts' },
  { name: 'admin-notifications', path: '/admin/notifications', title: 'Notifications' },
]

// ── Mock Data ────────────────────────────────────────────────────────────────

const now = new Date().toISOString()
const daysAgo = (n: number) => { const d = new Date(); d.setDate(d.getDate() - n); return d.toISOString() }

const mockProducts = [
  { id: 'p1', tenant_id: 'tnt_mock', name: 'Vitamin C Brightening Serum', description: 'A powerful antioxidant serum with 20% vitamin C to brighten and even skin tone.', sku: 'SKU-VC-001', price: { amount: 4500, currency: 'USD' }, images: ['https://placehold.co/400x400/f97316/white?text=Vitamin+C'], category_id: 'cat1', tags: ['skincare', 'serum', 'bestseller'], status: 'active', stock: 142, created_at: daysAgo(60), updated_at: daysAgo(2) },
  { id: 'p2', tenant_id: 'tnt_mock', name: 'Hydrating Rose Cream', description: 'Deep hydration cream with rose extract and hyaluronic acid.', sku: 'SKU-HR-002', price: { amount: 3200, currency: 'USD' }, images: ['https://placehold.co/400x400/ec4899/white?text=Rose+Cream'], category_id: 'cat1', tags: ['skincare', 'moisturizer'], status: 'active', stock: 89, created_at: daysAgo(45), updated_at: daysAgo(5) },
  { id: 'p3', tenant_id: 'tnt_mock', name: 'Retinol Night Oil', description: 'Advanced retinol formula for overnight skin renewal and anti-aging.', sku: 'SKU-RN-003', price: { amount: 5800, currency: 'USD' }, images: ['https://placehold.co/400x400/8b5cf6/white?text=Retinol'], category_id: 'cat1', tags: ['skincare', 'anti-aging'], status: 'active', stock: 67, created_at: daysAgo(30), updated_at: daysAgo(1) },
  { id: 'p4', tenant_id: 'tnt_mock', name: 'SPF 50 Daily Sunscreen', description: 'Lightweight mineral sunscreen with broad spectrum SPF 50 protection.', sku: 'SKU-SS-004', price: { amount: 2800, currency: 'USD' }, images: ['https://placehold.co/400x400/eab308/white?text=SPF+50'], category_id: 'cat2', tags: ['sun-protection', 'daily'], status: 'active', stock: 215, created_at: daysAgo(90), updated_at: daysAgo(3) },
  { id: 'p5', tenant_id: 'tnt_mock', name: 'Gentle Cleansing Balm', description: 'Melting balm cleanser that removes makeup and impurities without stripping.', sku: 'SKU-CB-005', price: { amount: 2400, currency: 'USD' }, images: ['https://placehold.co/400x400/10b981/white?text=Cleanser'], category_id: 'cat2', tags: ['cleanser', 'gentle'], status: 'active', stock: 178, created_at: daysAgo(75), updated_at: daysAgo(7) },
  { id: 'p6', tenant_id: 'tnt_mock', name: 'Niacinamide Pore Serum', description: '10% niacinamide serum to minimize pores and control oil production.', sku: 'SKU-NP-006', price: { amount: 3800, currency: 'USD' }, images: ['https://placehold.co/400x400/06b6d4/white?text=Niacinamide'], category_id: 'cat1', tags: ['skincare', 'serum', 'pore-care'], status: 'active', stock: 94, created_at: daysAgo(20), updated_at: daysAgo(1) },
  { id: 'p7', tenant_id: 'tnt_mock', name: 'Exfoliating AHA Toner', description: 'Glycolic acid toner for gentle chemical exfoliation and radiance.', sku: 'SKU-AT-007', price: { amount: 2900, currency: 'USD' }, images: ['https://placehold.co/400x400/f43f5e/white?text=AHA+Toner'], category_id: 'cat3', tags: ['toner', 'exfoliant'], status: 'draft', stock: 0, created_at: daysAgo(5), updated_at: daysAgo(1) },
]

const mockOrders = [
  { id: 'o1', tenant_id: 'tnt_mock', customer_id: 'c1', items: [{ id: 'oi1', product_id: 'p1', product_name: 'Vitamin C Brightening Serum', quantity: 2, unit_price: { amount: 4500, currency: 'USD' }, total: { amount: 9000, currency: 'USD' } }, { id: 'oi2', product_id: 'p4', product_name: 'SPF 50 Daily Sunscreen', quantity: 1, unit_price: { amount: 2800, currency: 'USD' }, total: { amount: 2800, currency: 'USD' } }], status: 'pending', total_amount: { amount: 11800, currency: 'USD' }, shipping_address: { street: '123 Main St', city: 'San Francisco', state: 'CA', country: 'US', postal_code: '94102', is_default: true }, created_at: daysAgo(0), updated_at: daysAgo(0) },
  { id: 'o2', tenant_id: 'tnt_mock', customer_id: 'c2', items: [{ id: 'oi3', product_id: 'p2', product_name: 'Hydrating Rose Cream', quantity: 1, unit_price: { amount: 3200, currency: 'USD' }, total: { amount: 3200, currency: 'USD' } }], status: 'shipped', total_amount: { amount: 3200, currency: 'USD' }, shipping_address: { street: '456 Oak Ave', city: 'Los Angeles', state: 'CA', country: 'US', postal_code: '90001', is_default: true }, created_at: daysAgo(3), updated_at: daysAgo(1) },
  { id: 'o3', tenant_id: 'tnt_mock', customer_id: 'c3', items: [{ id: 'oi4', product_id: 'p3', product_name: 'Retinol Night Oil', quantity: 1, unit_price: { amount: 5800, currency: 'USD' }, total: { amount: 5800, currency: 'USD' } }, { id: 'oi5', product_id: 'p5', product_name: 'Gentle Cleansing Balm', quantity: 2, unit_price: { amount: 2400, currency: 'USD' }, total: { amount: 4800, currency: 'USD' } }], status: 'delivered', total_amount: { amount: 10600, currency: 'USD' }, shipping_address: { street: '789 Pine Rd', city: 'New York', state: 'NY', country: 'US', postal_code: '10001', is_default: true }, created_at: daysAgo(7), updated_at: daysAgo(2) },
  { id: 'o4', tenant_id: 'tnt_mock', customer_id: 'c4', items: [{ id: 'oi6', product_id: 'p6', product_name: 'Niacinamide Pore Serum', quantity: 3, unit_price: { amount: 3800, currency: 'USD' }, total: { amount: 11400, currency: 'USD' } }], status: 'processing', total_amount: { amount: 11400, currency: 'USD' }, shipping_address: { street: '321 Elm St', city: 'Chicago', state: 'IL', country: 'US', postal_code: '60601', is_default: true }, created_at: daysAgo(1), updated_at: daysAgo(0) },
  { id: 'o5', tenant_id: 'tnt_mock', customer_id: 'c5', items: [{ id: 'oi7', product_id: 'p1', product_name: 'Vitamin C Brightening Serum', quantity: 1, unit_price: { amount: 4500, currency: 'USD' }, total: { amount: 4500, currency: 'USD' } }], status: 'confirmed', total_amount: { amount: 4500, currency: 'USD' }, shipping_address: { street: '654 Maple Dr', city: 'Austin', state: 'TX', country: 'US', postal_code: '73301', is_default: true }, created_at: daysAgo(2), updated_at: daysAgo(1) },
]

const mockCustomers = [
  { id: 'c1', tenant_id: 'tnt_mock', email: 'sarah.chen@example.com', name: 'Sarah Chen', phone: '+1 415-555-0101', addresses: [{ street: '123 Main St', city: 'San Francisco', state: 'CA', country: 'US', postal_code: '94102', is_default: true }], created_at: daysAgo(120), updated_at: daysAgo(0) },
  { id: 'c2', tenant_id: 'tnt_mock', email: 'james.wilson@example.com', name: 'James Wilson', phone: '+1 310-555-0202', addresses: [{ street: '456 Oak Ave', city: 'Los Angeles', state: 'CA', country: 'US', postal_code: '90001', is_default: true }], created_at: daysAgo(90), updated_at: daysAgo(1) },
  { id: 'c3', tenant_id: 'tnt_mock', email: 'maria.garcia@example.com', name: 'Maria Garcia', phone: '+1 212-555-0303', addresses: [{ street: '789 Pine Rd', city: 'New York', state: 'NY', country: 'US', postal_code: '10001', is_default: true }], created_at: daysAgo(60), updated_at: daysAgo(2) },
  { id: 'c4', tenant_id: 'tnt_mock', email: 'alex.kim@example.com', name: 'Alex Kim', phone: '+1 312-555-0404', addresses: [{ street: '321 Elm St', city: 'Chicago', state: 'IL', country: 'US', postal_code: '60601', is_default: true }], created_at: daysAgo(30), updated_at: daysAgo(0) },
  { id: 'c5', tenant_id: 'tnt_mock', email: 'emma.brown@example.com', name: 'Emma Brown', phone: '+1 512-555-0505', addresses: [{ street: '654 Maple Dr', city: 'Austin', state: 'TX', country: 'US', postal_code: '73301', is_default: true }], created_at: daysAgo(15), updated_at: daysAgo(1) },
]

const mockCategories = [
  { id: 'cat1', tenant_id: 'tnt_mock', name: 'Serums & Treatments', slug: 'serums-treatments', parent_id: null, description: 'Targeted treatments and serums for specific skin concerns.', created_at: daysAgo(120), updated_at: daysAgo(10) },
  { id: 'cat2', tenant_id: 'tnt_mock', name: 'Cleansers & Sun Care', slug: 'cleansers-sun-care', parent_id: null, description: 'Daily cleansing and sun protection products.', created_at: daysAgo(120), updated_at: daysAgo(10) },
  { id: 'cat3', tenant_id: 'tnt_mock', name: 'Toners & Exfoliants', slug: 'toners-exfoliants', parent_id: null, description: 'Prep and polish your skin with our toners and exfoliants.', created_at: daysAgo(60), updated_at: daysAgo(5) },
]

const mockCollections = [
  { id: 'col1', tenant_id: 'tnt_mock', name: 'Best Sellers', slug: 'best-sellers', description: 'Our most popular products loved by thousands.', product_ids: ['p1', 'p2', 'p3'], is_automatic: false, rules: {}, created_at: daysAgo(90), updated_at: daysAgo(1) },
  { id: 'col2', tenant_id: 'tnt_mock', name: 'New Arrivals', slug: 'new-arrivals', description: 'Just launched — discover our latest formulations.', product_ids: ['p6', 'p7'], is_automatic: false, rules: {}, created_at: daysAgo(30), updated_at: daysAgo(5) },
]

const mockPages = [
  { id: 'pg1', tenant_id: 'tnt_mock', slug: 'home', title: 'Home', html: '', css: '', meta: { description: 'Welcome to Hada Store', og_title: 'Hada Store', og_image: '', keywords: ['skincare', 'beauty'] }, content_type: 'blocks', sections: [
    { id: 's1', type: 'hero', settings: { title: 'Radiant Skin Starts Here', subtitle: 'Discover our curated collection of clean beauty essentials', button_text: 'Shop Now', button_url: '/products', background_color: '#4f46e5' }, blocks: [] },
    { id: 's2', type: 'featured_collection', settings: { collection_id: 'col1', title: 'Best Sellers', columns: 3 }, blocks: [] },
    { id: 's3', type: 'rich_text', settings: { content: '<h2>Our Promise</h2><p>Every product is dermatologist-tested, cruelty-free, and made with sustainable ingredients.</p>' }, blocks: [] },
  ], status: 'published', version: 3, created_by: 'usr_mock', published_at: daysAgo(7), created_at: daysAgo(60), updated_at: daysAgo(1) },
  { id: 'pg2', tenant_id: 'tnt_mock', slug: 'about', title: 'About Us', html: '<h1>About Hada</h1><p>We believe great skincare should be accessible to everyone.</p>', css: '', meta: { description: 'About Hada Store', og_title: 'About Us', og_image: '', keywords: ['about'] }, content_type: 'html', sections: [], status: 'published', version: 1, created_by: 'usr_mock', published_at: daysAgo(30), created_at: daysAgo(60), updated_at: daysAgo(30) },
  { id: 'pg3', tenant_id: 'tnt_mock', slug: 'summer-sale', title: 'Summer Sale 2026', html: '', css: '', meta: { description: 'Summer sale — up to 40% off', og_title: 'Summer Sale', og_image: '', keywords: ['sale', 'summer'] }, content_type: 'blocks', sections: [
    { id: 's4', type: 'banner', settings: { text: 'Summer Sale — Up to 40% Off!', background_color: '#dc2626', text_color: '#ffffff' }, blocks: [] },
    { id: 's5', type: 'product_grid', settings: { title: 'Sale Items', product_ids: ['p1', 'p2', 'p4', 'p5'], columns: 4 }, blocks: [] },
  ], status: 'draft', version: 1, created_by: 'usr_mock', created_at: daysAgo(3), updated_at: daysAgo(1) },
]

const mockBlockTypes = [
  { id: 'bt-hero', name: 'hero', display_name: 'Hero Banner', category: 'content', icon: 'layout', schema: { title: { type: 'text', label: 'Title' }, subtitle: { type: 'text', label: 'Subtitle' }, button_text: { type: 'text', label: 'Button Text' }, button_url: { type: 'url', label: 'Button URL' }, background_color: { type: 'color', label: 'Background' } }, default_settings: { title: 'Welcome', subtitle: '', button_text: 'Learn More', button_url: '/', background_color: '#4f46e5' }, created_at: daysAgo(90), updated_at: daysAgo(90) },
  { id: 'bt-rich-text', name: 'rich_text', display_name: 'Rich Text', category: 'content', icon: 'type', schema: { content: { type: 'richtext', label: 'Content' } }, default_settings: { content: '<p>Enter your text here...</p>' }, created_at: daysAgo(90), updated_at: daysAgo(90) },
  { id: 'bt-product-grid', name: 'product_grid', display_name: 'Product Grid', category: 'commerce', icon: 'grid', schema: { title: { type: 'text', label: 'Title' }, columns: { type: 'number', label: 'Columns', min: 2, max: 6 }, product_ids: { type: 'product_list', label: 'Products' } }, default_settings: { title: 'Products', columns: 3, product_ids: [] }, created_at: daysAgo(90), updated_at: daysAgo(90) },
  { id: 'bt-featured-collection', name: 'featured_collection', display_name: 'Featured Collection', category: 'commerce', icon: 'star', schema: { collection_id: { type: 'collection', label: 'Collection' }, title: { type: 'text', label: 'Title' }, columns: { type: 'number', label: 'Columns' } }, default_settings: { collection_id: '', title: 'Featured', columns: 3 }, created_at: daysAgo(90), updated_at: daysAgo(90) },
  { id: 'bt-banner', name: 'banner', display_name: 'Banner', category: 'content', icon: 'alert-circle', schema: { text: { type: 'text', label: 'Banner Text' }, background_color: { type: 'color', label: 'Background' }, text_color: { type: 'color', label: 'Text Color' } }, default_settings: { text: 'Announcement', background_color: '#1e40af', text_color: '#ffffff' }, created_at: daysAgo(90), updated_at: daysAgo(90) },
  { id: 'bt-image', name: 'image_block', display_name: 'Image', category: 'media', icon: 'image', schema: { src: { type: 'image', label: 'Image' }, alt: { type: 'text', label: 'Alt Text' }, caption: { type: 'text', label: 'Caption' } }, default_settings: { src: '', alt: '', caption: '' }, created_at: daysAgo(90), updated_at: daysAgo(90) },
  { id: 'bt-video', name: 'video_block', display_name: 'Video', category: 'media', icon: 'play-circle', schema: { url: { type: 'url', label: 'Video URL' }, autoplay: { type: 'boolean', label: 'Autoplay' } }, default_settings: { url: '', autoplay: false }, created_at: daysAgo(90), updated_at: daysAgo(90) },
  { id: 'bt-newsletter', name: 'newsletter_signup', display_name: 'Newsletter Signup', category: 'content', icon: 'mail', schema: { heading: { type: 'text', label: 'Heading' }, placeholder: { type: 'text', label: 'Placeholder' }, button_text: { type: 'text', label: 'Button Text' } }, default_settings: { heading: 'Stay in the loop', placeholder: 'Enter your email', button_text: 'Subscribe' }, created_at: daysAgo(90), updated_at: daysAgo(90) },
  { id: 'bt-testimonials', name: 'testimonials', display_name: 'Testimonials', category: 'content', icon: 'message-circle', schema: { title: { type: 'text', label: 'Title' }, items: { type: 'array', label: 'Testimonials' } }, default_settings: { title: 'What our customers say', items: [] }, created_at: daysAgo(90), updated_at: daysAgo(90) },
  { id: 'bt-cta', name: 'cta', display_name: 'Call to Action', category: 'content', icon: 'zap', schema: { title: { type: 'text', label: 'Title' }, description: { type: 'text', label: 'Description' }, button_text: { type: 'text', label: 'Button Text' }, button_url: { type: 'url', label: 'URL' } }, default_settings: { title: 'Ready to start?', description: '', button_text: 'Get Started', button_url: '/' }, created_at: daysAgo(90), updated_at: daysAgo(90) },
]

const mockThemes = [
  {
    id: 'th1', tenant_id: 'tnt_mock', name: 'Default Light', is_active: true,
    tokens: {
      colors: { primary: '#4f46e5', secondary: '#7c3aed', background: '#ffffff', surface: '#f8fafc', text: '#0f172a', text_muted: '#64748b', border: '#e2e8f0', error: '#ef4444', success: '#22c55e', warning: '#f59e0b', info: '#3b82f6' },
      typography: { font_heading: 'Inter', font_body: 'Inter', base_size: '16px', scale_ratio: 1.25 },
      spacing: { unit: '0.25rem', section_padding: '4rem' },
      borders: { radius_sm: '0.25rem', radius_md: '0.5rem', radius_lg: '1rem', radius_full: '9999px' },
      shadows: { sm: '0 1px 2px rgba(0,0,0,0.05)', md: '0 4px 6px rgba(0,0,0,0.07)', lg: '0 10px 15px rgba(0,0,0,0.1)' },
    },
    created_at: daysAgo(60), updated_at: daysAgo(1),
  },
  {
    id: 'th2', tenant_id: 'tnt_mock', name: 'Dark Elegance', is_active: false,
    tokens: {
      colors: { primary: '#a78bfa', secondary: '#f472b6', background: '#0f172a', surface: '#1e293b', text: '#f1f5f9', text_muted: '#94a3b8', border: '#334155', error: '#f87171', success: '#4ade80', warning: '#fbbf24', info: '#60a5fa' },
      typography: { font_heading: 'Playfair Display', font_body: 'Inter', base_size: '16px', scale_ratio: 1.333 },
      spacing: { unit: '0.25rem', section_padding: '5rem' },
      borders: { radius_sm: '0.125rem', radius_md: '0.375rem', radius_lg: '0.75rem', radius_full: '9999px' },
      shadows: { sm: '0 1px 3px rgba(0,0,0,0.3)', md: '0 4px 8px rgba(0,0,0,0.4)', lg: '0 12px 24px rgba(0,0,0,0.5)' },
    },
    created_at: daysAgo(30), updated_at: daysAgo(5),
  },
  {
    id: 'th3', tenant_id: 'tnt_mock', name: 'Warm Minimal', is_active: false,
    tokens: {
      colors: { primary: '#d97706', secondary: '#b45309', background: '#fffbeb', surface: '#fef3c7', text: '#1c1917', text_muted: '#78716c', border: '#e7e5e4', error: '#dc2626', success: '#16a34a', warning: '#ea580c', info: '#0284c7' },
      typography: { font_heading: 'DM Serif Display', font_body: 'DM Sans', base_size: '17px', scale_ratio: 1.2 },
      spacing: { unit: '0.25rem', section_padding: '3.5rem' },
      borders: { radius_sm: '0.5rem', radius_md: '0.75rem', radius_lg: '1.5rem', radius_full: '9999px' },
      shadows: { sm: '0 1px 2px rgba(0,0,0,0.04)', md: '0 2px 4px rgba(0,0,0,0.06)', lg: '0 8px 16px rgba(0,0,0,0.08)' },
    },
    created_at: daysAgo(15), updated_at: daysAgo(3),
  },
]

const mockPromos = [
  { id: 'pr1', tenant_id: 'tnt_mock', code: 'SUMMER25', type: 'percentage', value: 25, min_order_amount: 5000, max_uses: 500, used_count: 142, starts_at: daysAgo(10), ends_at: new Date(Date.now() + 20 * 86400000).toISOString(), active: true, created_at: daysAgo(10) },
  { id: 'pr2', tenant_id: 'tnt_mock', code: 'FREESHIP', type: 'free_shipping', value: 0, min_order_amount: 3000, max_uses: 1000, used_count: 387, starts_at: daysAgo(30), ends_at: new Date(Date.now() + 60 * 86400000).toISOString(), active: true, created_at: daysAgo(30) },
  { id: 'pr3', tenant_id: 'tnt_mock', code: 'WELCOME10', type: 'fixed_amount', value: 1000, min_order_amount: 2000, max_uses: undefined, used_count: 56, starts_at: daysAgo(90), active: true, created_at: daysAgo(90) },
]

const mockMedia = [
  { id: 'm1', tenant_id: 'tnt_mock', filename: 'hero-banner.jpg', content_type: 'image/jpeg', size: 245000, url: 'https://placehold.co/1200x400/4f46e5/white?text=Hero+Banner', created_at: daysAgo(30), updated_at: daysAgo(30) },
  { id: 'm2', tenant_id: 'tnt_mock', filename: 'product-lifestyle.jpg', content_type: 'image/jpeg', size: 189000, url: 'https://placehold.co/800x600/ec4899/white?text=Lifestyle', created_at: daysAgo(25), updated_at: daysAgo(25) },
  { id: 'm3', tenant_id: 'tnt_mock', filename: 'store-logo.png', content_type: 'image/png', size: 12000, url: 'https://placehold.co/200x200/0f172a/white?text=HADA', created_at: daysAgo(60), updated_at: daysAgo(60) },
  { id: 'm4', tenant_id: 'tnt_mock', filename: 'about-team.jpg', content_type: 'image/jpeg', size: 320000, url: 'https://placehold.co/1000x500/10b981/white?text=Team+Photo', created_at: daysAgo(20), updated_at: daysAgo(20) },
]

const mockPlugins = [
  { id: 'plg1', name: 'reviews-pro', display_name: 'Reviews Pro', description: 'Advanced product reviews with photos, ratings, and Q&A.', author: 'Hada Labs', icon: 'star', category: 'marketing', tags: ['reviews', 'social-proof', 'ugc'], created_at: daysAgo(180), updated_at: daysAgo(10) },
  { id: 'plg2', name: 'email-campaigns', display_name: 'Email Campaigns', description: 'Automated email marketing with abandoned cart recovery and newsletters.', author: 'Hada Labs', icon: 'mail', category: 'marketing', tags: ['email', 'automation', 'marketing'], created_at: daysAgo(150), updated_at: daysAgo(5) },
  { id: 'plg3', name: 'analytics-plus', display_name: 'Analytics Plus', description: 'Enhanced analytics with cohort analysis, funnel tracking, and custom reports.', author: 'DataFlow Inc', icon: 'bar-chart', category: 'analytics', tags: ['analytics', 'reporting', 'data'], created_at: daysAgo(120), updated_at: daysAgo(15) },
]

const mockInstallations = [
  { id: 'inst1', tenant_id: 'tnt_mock', plugin_id: 'plg1', version_id: 'v1', status: 'active', settings: { display_mode: 'inline', require_purchase: true }, installed_at: daysAgo(30), updated_at: daysAgo(5) },
  { id: 'inst2', tenant_id: 'tnt_mock', plugin_id: 'plg2', version_id: 'v2', status: 'active', settings: { sender_email: 'hello@hada.store' }, installed_at: daysAgo(20), updated_at: daysAgo(2) },
]

const mockSettings = {
  tenant_id: 'tnt_mock', store_name: 'Hada Beauty', store_email: 'hello@hada.store',
  store_phone: '+1 555-0123', currency: 'USD', timezone: 'America/New_York',
  address: { street: '123 Commerce St', city: 'San Francisco', state: 'CA', zip: '94102', country: 'US' },
  logo_url: '', favicon_url: '',
  social_links: { instagram: 'https://instagram.com/hadabeauty', twitter: 'https://x.com/hadabeauty', facebook: '', tiktok: 'https://tiktok.com/@hadabeauty' },
  checkout_config: { require_phone: false, require_shipping: true, allow_notes: true },
  updated_at: now,
}

// ── Test Setup ───────────────────────────────────────────────────────────────

test.beforeAll(() => {
  if (!fs.existsSync(SCREENSHOT_DIR)) {
    fs.mkdirSync(SCREENSHOT_DIR, { recursive: true })
  }
})

test.beforeEach(async ({ page }) => {
  // Set fake auth tokens
  await page.addInitScript(() => {
    localStorage.setItem('hada_access_token', 'mock-token-for-screenshots')
    localStorage.setItem('hada_refresh_token', 'mock-refresh-token')
    localStorage.setItem(
      'auth_user',
      JSON.stringify({
        id: 'usr_mock', tenant_id: 'tnt_mock', email: 'demo@hada.commerce',
        name: 'Demo User', status: 'active', scopes: ['admin'],
      }),
    )
    localStorage.setItem(
      'auth_tenant',
      JSON.stringify({
        id: 'tnt_mock', name: 'Hada Beauty', slug: 'hada-beauty',
        plan: 'pro', is_active: true,
      }),
    )
  })

  // Intercept all API calls with rich mock data
  await page.route('**/api/v1/**', (route) => {
    const url = route.request().url()
    const method = route.request().method()
    const json = (body: unknown) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(body) })
    const paginated = (items: unknown[]) =>
      json({ items, total: items.length, page: 1, page_size: 20, total_pages: 1 })

    // Auth
    if (url.includes('/auth/me')) {
      return json({
        user: { id: 'usr_mock', tenant_id: 'tnt_mock', email: 'demo@hada.commerce', name: 'Demo User', status: 'active', scopes: ['admin'] },
        tenant: { id: 'tnt_mock', name: 'Hada Beauty', slug: 'hada-beauty', plan: 'pro', is_active: true },
      })
    }

    // Analytics
    if (url.includes('/analytics/dashboard')) {
      return json({ total_products: 42, total_orders: 156, total_customers: 89, total_revenue: 1284500, currency: 'USD', pending_orders: 12, active_promos: 3, pending_pages: 2 })
    }
    if (url.includes('/analytics/revenue')) {
      return json(Array.from({ length: 30 }, (_, i) => {
        const d = new Date(); d.setDate(d.getDate() - (29 - i))
        return { date: d.toISOString().slice(0, 10), amount: Math.floor(20000 + Math.sin(i * 0.5) * 15000 + i * 2000) }
      }))
    }
    if (url.includes('/analytics/top-products')) {
      return json([
        { product_id: 'p1', name: 'Vitamin C Brightening Serum', total_sold: 245, revenue: 1102500 },
        { product_id: 'p3', name: 'Retinol Night Oil', total_sold: 156, revenue: 904800 },
        { product_id: 'p2', name: 'Hydrating Rose Cream', total_sold: 198, revenue: 633600 },
        { product_id: 'p4', name: 'SPF 50 Daily Sunscreen', total_sold: 134, revenue: 375200 },
        { product_id: 'p5', name: 'Gentle Cleansing Balm', total_sold: 112, revenue: 268800 },
      ])
    }
    if (url.includes('/analytics/order-status')) {
      return json([
        { status: 'pending', count: 12 }, { status: 'processing', count: 8 },
        { status: 'shipped', count: 24 }, { status: 'delivered', count: 98 }, { status: 'cancelled', count: 5 },
      ])
    }
    if (url.includes('/analytics/recent-orders')) {
      return json(mockOrders.slice(0, 5).map(o => ({
        id: o.id, customer_name: mockCustomers.find(c => c.id === o.customer_id)?.name ?? 'Unknown',
        total: o.total_amount.amount, currency: 'USD', status: o.status, created_at: o.created_at,
      })))
    }

    // Themes
    if (url.includes('/themes/active')) {
      return json(mockThemes[0])
    }
    if (url.match(/\/themes\/?(\?.*)?$/) && !url.includes('/active')) {
      if (method === 'GET') return json(mockThemes)
      return json(mockThemes[0])
    }
    if (url.match(/\/themes\/th\d/)) {
      const id = url.match(/\/themes\/(th\d)/)?.[1]
      return json(mockThemes.find(t => t.id === id) ?? mockThemes[0])
    }

    // Block types
    if (url.match(/\/block-types\/?(\?.*)?$/) || url.match(/\/storefront\/block-types/)) {
      return json(mockBlockTypes)
    }

    // Storefront pages
    if (url.includes('/storefront/pages/by-slug/')) {
      return json(mockPages[0])
    }
    if (url.match(/\/storefront\/pages\/(pg\d)/)) {
      const id = url.match(/\/storefront\/pages\/(pg\d)/)?.[1]
      return json(mockPages.find(p => p.id === id) ?? mockPages[0])
    }
    if (url.match(/\/storefront\/pages\/?(\?.*)?$/)) {
      return paginated(mockPages)
    }
    // Also match /pages/ for legacy routes
    if (url.match(/\/pages\/(pg\d)/)) {
      const id = url.match(/\/pages\/(pg\d)/)?.[1]
      return json(mockPages.find(p => p.id === id) ?? mockPages[0])
    }
    if (url.match(/\/pages\/?(\?.*)?$/) && !url.includes('/storefront/')) {
      return paginated(mockPages)
    }

    // Products
    if (url.match(/\/products\/(p\d)/)) {
      const id = url.match(/\/products\/(p\d)/)?.[1]
      return json(mockProducts.find(p => p.id === id) ?? mockProducts[0])
    }
    if (url.match(/\/products\/?(\?.*)?$/)) {
      return paginated(mockProducts)
    }

    // Orders
    if (url.match(/\/orders\/(o\d)/)) {
      const id = url.match(/\/orders\/(o\d)/)?.[1]
      return json(mockOrders.find(o => o.id === id) ?? mockOrders[0])
    }
    if (url.match(/\/orders\/?(\?.*)?$/)) {
      return paginated(mockOrders)
    }

    // Customers
    if (url.match(/\/customers\/(c\d)/)) {
      const id = url.match(/\/customers\/(c\d)/)?.[1]
      return json(mockCustomers.find(c => c.id === id) ?? mockCustomers[0])
    }
    if (url.match(/\/customers\/?(\?.*)?$/)) {
      return paginated(mockCustomers)
    }

    // Catalog
    if (url.match(/\/categories\/?(\?.*)?$/)) {
      return paginated(mockCategories)
    }
    if (url.match(/\/collections\/?(\?.*)?$/)) {
      return paginated(mockCollections)
    }

    // Promos
    if (url.match(/\/promos\/?(\?.*)?$/)) {
      return paginated(mockPromos)
    }

    // Media
    if (url.match(/\/media\/?(\?.*)?$/)) {
      return paginated(mockMedia)
    }

    // Plugins
    if (url.includes('/plugins/installed')) {
      return paginated(mockInstallations)
    }
    if (url.includes('/plugins/js-manifest')) {
      return json({ scripts: [] })
    }
    if (url.match(/\/plugins\/(plg\d)/)) {
      const id = url.match(/\/plugins\/(plg\d)/)?.[1]
      return json(mockPlugins.find(p => p.id === id) ?? mockPlugins[0])
    }
    if (url.match(/\/plugins\/?(\?.*)?$/) && !url.includes('/installed') && !url.includes('/js-manifest')) {
      return paginated(mockPlugins)
    }

    // Marketplace (legacy)
    if (url.includes('/marketplace/installed')) {
      return paginated(mockInstallations)
    }
    if (url.match(/\/marketplace\/plugins\/([\w-]+)/)) {
      return json({ ...mockPlugins[0], latest_version: { id: 'v1', plugin_id: 'plg1', version: '1.2.0', changelog: 'Bug fixes', frontend_url: 'https://cdn.hada.store/plugins/reviews-pro/widget.js', backend_entry: 'https://plugins.hada.store/reviews-pro/webhook', config_schema: {}, created_at: daysAgo(10) } })
    }

    // Settings
    if (url.match(/\/settings\/?(\?.*)?$/) && !url.includes('/plugins/')) {
      return json(mockSettings)
    }

    // Shipping
    if (url.match(/\/shipping\/zones\/[\w-]+\/rates\/?(\?.*)?$/)) {
      return paginated(mockShippingRates)
    }
    if (url.match(/\/shipping\/zones\/?(\?.*)?$/)) {
      return paginated(mockShippingZones)
    }
    if (url.match(/\/shipping\/rates\/?(\?.*)?$/)) {
      return paginated(mockShippingRates)
    }

    // Tax
    if (url.match(/\/tax\/rates\/?(\?.*)?$/)) {
      return paginated(mockTaxRates)
    }

    // Payments
    if (url.match(/\/payments\/order\//)) {
      return paginated(mockPayments)
    }
    if (url.match(/\/payments\/[\w-]+\/refunds\/?(\?.*)?$/)) {
      return paginated([])
    }
    if (url.match(/\/payments\/?(\?.*)?$/)) {
      return paginated(mockPayments)
    }

    // Customer Groups
    if (url.match(/\/customer-groups\/[\w-]+\/members\/?(\?.*)?$/)) {
      return paginated(mockCustomers.slice(0, 2))
    }
    if (url.match(/\/customer-groups\/?(\?.*)?$/)) {
      return paginated(mockCustomerGroups)
    }

    // Gift Cards
    if (url.match(/\/gift-cards\/[\w-]+\/transactions\/?(\?.*)?$/)) {
      return paginated([])
    }
    if (url.match(/\/gift-cards\/?(\?.*)?$/)) {
      return paginated(mockGiftCards)
    }

    // Cart Recovery
    if (url.match(/\/cart-recovery\/stats\/?(\?.*)?$/)) {
      return json({ total_abandoned: 47, recovered: 12, recovery_rate: 25.5, revenue_recovered: 189400 })
    }
    if (url.match(/\/cart-recovery\/?(\?.*)?$/)) {
      return paginated(mockCartRecovery)
    }

    // Currency Rates
    if (url.match(/\/currencies\/?(\?.*)?$/)) {
      return json(['USD', 'EUR', 'GBP', 'CAD', 'AUD', 'JPY', 'CNY', 'MXN'])
    }
    if (url.match(/\/currency\/convert\/?(\?.*)?$/)) {
      return json({ from: 'USD', to: 'EUR', amount: 100, converted: 92.5, rate: 0.925 })
    }
    if (url.match(/\/currency-rates\/?(\?.*)?$/)) {
      return paginated(mockCurrencyRates)
    }

    // Translations
    if (url.match(/\/i18n\/supported-locales\/?(\?.*)?$/)) {
      return json(['en', 'es', 'fr', 'de', 'ja', 'zh'])
    }
    if (url.match(/\/i18n\/[\w-]+\/[\w-]+\/locales\/?(\?.*)?$/)) {
      return json(['en', 'es', 'fr'])
    }
    if (url.match(/\/i18n\//)) {
      return json({ locale: 'es', fields: { name: 'Serum de Vitamina C', description: 'Un poderoso suero antioxidante.' } })
    }

    // Subscriptions
    if (url.match(/\/subscriptions\/[\w-]+\/billing\/?(\?.*)?$/)) {
      return paginated([])
    }
    if (url.match(/\/subscriptions\/due\/?(\?.*)?$/)) {
      return paginated(mockSubscriptions.slice(0, 1))
    }
    if (url.match(/\/subscriptions\/?(\?.*)?$/)) {
      return paginated(mockSubscriptions)
    }

    // Inventory
    if (url.match(/\/inventory\/stock\/low\/?(\?.*)?$/)) {
      return paginated(mockLowStockAlerts)
    }
    if (url.match(/\/inventory\/stock\/[\w-]+\/?(\?.*)?$/)) {
      return paginated([])
    }
    if (url.match(/\/inventory\/warehouses\/?(\?.*)?$/)) {
      return paginated(mockWarehouses)
    }

    // Reviews
    if (url.match(/\/reviews\/?(\?.*)?$/)) {
      return paginated(mockReviews)
    }

    // Returns
    if (url.match(/\/returns\/[\w-]+\/?(\?.*)?$/) && !url.match(/\/returns\/?(\?.*)?$/)) {
      return json(mockReturns[0])
    }
    if (url.match(/\/returns\/?(\?.*)?$/)) {
      return paginated(mockReturns)
    }

    // Webhooks
    if (url.match(/\/webhooks\/[\w-]+\/deliveries\/?(\?.*)?$/)) {
      return paginated(mockWebhookDeliveries)
    }
    if (url.match(/\/webhooks\/?(\?.*)?$/)) {
      return paginated(mockWebhooks)
    }

    // Audit Logs
    if (url.match(/\/audit\/stats\/?(\?.*)?$/)) {
      return json({ total_events: 1243, unique_users: 5, top_actions: ['update', 'create', 'delete'] })
    }
    if (url.match(/\/audit\/[\w-]+\/?(\?.*)?$/) && !url.match(/\/audit\/stats/)) {
      return json(mockAuditLogs[0])
    }
    if (url.match(/\/audit\/?(\?.*)?$/)) {
      return paginated(mockAuditLogs)
    }

    // Loyalty
    if (url.match(/\/loyalty\/accounts\/[\w-]+\/transactions\/?(\?.*)?$/)) {
      return paginated([])
    }
    if (url.match(/\/loyalty\/accounts\/?(\?.*)?$/)) {
      return paginated(mockLoyaltyAccounts)
    }
    if (url.match(/\/loyalty\/rewards\/?(\?.*)?$/)) {
      return paginated(mockLoyaltyRewards)
    }

    // Bundles
    if (url.match(/\/bundles\/[\w-]+\/price\/?(\?.*)?$/)) {
      return json({ bundle_id: 'bun1', total: { amount: 7500, currency: 'USD' }, savings: { amount: 1500, currency: 'USD' } })
    }
    if (url.match(/\/bundles\/?(\?.*)?$/)) {
      return paginated(mockBundles)
    }

    // Reporting / Dashboard
    if (url.match(/\/dashboard\/sales\/?(\?.*)?$/)) {
      return json({ total_revenue: 1284500, order_count: 156, average_order_value: 8234, refund_total: 45000 })
    }
    if (url.match(/\/dashboard\/top-products\/?(\?.*)?$/)) {
      return json([
        { product_id: 'p1', name: 'Vitamin C Brightening Serum', total_sold: 245, revenue: 1102500 },
        { product_id: 'p3', name: 'Retinol Night Oil', total_sold: 156, revenue: 904800 },
        { product_id: 'p2', name: 'Hydrating Rose Cream', total_sold: 198, revenue: 633600 },
        { product_id: 'p4', name: 'SPF 50 Daily Sunscreen', total_sold: 134, revenue: 375200 },
        { product_id: 'p5', name: 'Gentle Cleansing Balm', total_sold: 112, revenue: 268800 },
      ])
    }
    if (url.match(/\/dashboard\/revenue\/?(\?.*)?$/)) {
      return json(Array.from({ length: 30 }, (_, i) => {
        const d = new Date(); d.setDate(d.getDate() - (29 - i))
        return { date: d.toISOString().slice(0, 10), amount: Math.floor(20000 + Math.sin(i * 0.5) * 15000 + i * 2000) }
      }))
    }
    if (url.match(/\/dashboard\/customers\/?(\?.*)?$/)) {
      return json({ new_customers: 23, returning_customers: 66, churn_rate: 3.2, lifetime_value: 18500 })
    }
    if (url.match(/\/dashboard\/funnel\/?(\?.*)?$/)) {
      return json([
        { step: 'Visited', count: 4200 }, { step: 'Added to Cart', count: 1890 },
        { step: 'Checkout Started', count: 720 }, { step: 'Purchased', count: 540 },
      ])
    }

    // Social Accounts
    if (url.match(/\/social-accounts\/?(\?.*)?$/)) {
      return paginated(mockSocialAccounts)
    }

    // Notifications
    if (url.match(/\/notifications\/unread-count\/?(\?.*)?$/)) {
      return json({ count: 3 })
    }
    if (url.match(/\/notifications\/?(\?.*)?$/)) {
      return paginated(mockNotifications)
    }

    // Storefronts (Multistores)
    if (url.match(/\/storefronts\/?(\?.*)?$/)) {
      return paginated(mockStorefronts)
    }

    // Blog
    if (url.match(/\/blog\/categories\/?(\?.*)?$/)) {
      return paginated(mockBlogCategories)
    }
    if (url.match(/\/blog\/?(\?.*)?$/)) {
      return paginated(mockBlogPosts)
    }

    // A/B Testing
    if (url.match(/\/experiments\/[\w-]+\/results\/?(\?.*)?$/)) {
      return json({ experiment_id: 'exp1', variants: [{ name: 'control', conversions: 45, visitors: 500 }, { name: 'variant_a', conversions: 62, visitors: 498 }] })
    }
    if (url.match(/\/experiments\/?(\?.*)?$/)) {
      return paginated(mockExperiments)
    }

    // Recommendations
    if (url.match(/\/recommendations\/?(\?.*)?$/)) {
      return paginated(mockRecommendations)
    }

    // Bulk Operations
    if (url.match(/\/bulk-operations\/[\w-]+\/items\/?(\?.*)?$/)) {
      return paginated([])
    }
    if (url.match(/\/bulk-operations\/?(\?.*)?$/)) {
      return paginated(mockBulkOperations)
    }

    // Default fallback
    return json({ items: [], total: 0, page: 1, page_size: 20, total_pages: 0 })
  })
})

// ── Screenshot Tests ─────────────────────────────────────────────────────────

for (const p of pages) {
  test(`screenshot: ${p.title}`, async ({ page }) => {
    await page.goto(p.path, { waitUntil: 'networkidle', timeout: 15000 }).catch(() => {})
    await page.waitForTimeout(1500)
    await page.screenshot({
      path: path.join(SCREENSHOT_DIR, `${p.name}.png`),
      fullPage: true,
    })
  })
}

// ── Generate HTML Gallery ────────────────────────────────────────────────────

test.afterAll(() => {
  const files = fs.readdirSync(SCREENSHOT_DIR).filter((f) => f.endsWith('.png')).sort()

  const html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Hada Commerce - UI Screenshots</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; background: #0a0a0a; color: #e5e5e5; padding: 2rem; }
    h1 { text-align: center; font-size: 2rem; margin-bottom: 0.5rem; }
    .subtitle { text-align: center; color: #888; margin-bottom: 2rem; }
    .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(600px, 1fr)); gap: 2rem; max-width: 1800px; margin: 0 auto; }
    .card { background: #141414; border: 1px solid #222; border-radius: 12px; overflow: hidden; transition: transform 0.2s; }
    .card:hover { transform: translateY(-4px); border-color: #444; }
    .card img { width: 100%; height: auto; display: block; border-bottom: 1px solid #222; }
    .card .label { padding: 1rem; font-size: 0.95rem; font-weight: 500; }
    .card .label .path { color: #666; font-size: 0.8rem; margin-top: 0.25rem; }
  </style>
</head>
<body>
  <h1>Hada Commerce</h1>
  <p class="subtitle">UI Screenshots &mdash; ${new Date().toLocaleDateString()}</p>
  <div class="grid">
${files
  .map((f) => {
    const name = f.replace('.png', '').replace(/-/g, ' ')
    const matchedPage = pages.find((p) => p.name === f.replace('.png', ''))
    const pagePath = matchedPage?.path ?? ''
    return `    <div class="card">
      <img src="${f}" alt="${name}" loading="lazy" />
      <div class="label">${name}<div class="path">${pagePath}</div></div>
    </div>`
  })
  .join('\\n')}
  </div>
</body>
</html>`

  fs.writeFileSync(path.join(SCREENSHOT_DIR, 'index.html'), html)
})
