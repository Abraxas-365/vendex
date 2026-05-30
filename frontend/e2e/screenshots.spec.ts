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
  { name: 'store-products', path: '/products', title: 'Product List' },
  { name: 'store-cart', path: '/cart', title: 'Shopping Cart' },
  { name: 'store-checkout', path: '/checkout', title: 'Checkout' },

  // Auth
  { name: 'auth-login', path: '/login', title: 'Login' },

  // Admin pages
  { name: 'admin-dashboard', path: '/admin', title: 'Admin Dashboard' },
  { name: 'admin-products', path: '/admin/products', title: 'Admin Products' },
  { name: 'admin-orders', path: '/admin/orders', title: 'Admin Orders' },
  { name: 'admin-customers', path: '/admin/customers', title: 'Admin Customers' },
  { name: 'admin-catalog', path: '/admin/catalog', title: 'Admin Catalog' },
  { name: 'admin-pages', path: '/admin/pages', title: 'Admin Pages' },
  { name: 'admin-promos', path: '/admin/promos', title: 'Admin Promos' },
  { name: 'admin-media', path: '/admin/media', title: 'Admin Media' },
  { name: 'admin-marketplace', path: '/admin/marketplace', title: 'Admin Marketplace' },
  { name: 'admin-settings', path: '/admin/settings', title: 'Admin Settings' },
  { name: 'admin-agent', path: '/admin/agent', title: 'AI Agent Chat' },
]

test.beforeAll(() => {
  if (!fs.existsSync(SCREENSHOT_DIR)) {
    fs.mkdirSync(SCREENSHOT_DIR, { recursive: true })
  }
})

// Mock auth so admin pages render without redirect
test.beforeEach(async ({ page }) => {
  // Set fake auth tokens in localStorage so protected routes render
  await page.addInitScript(() => {
    localStorage.setItem('hada_access_token', 'mock-token-for-screenshots')
    localStorage.setItem('hada_refresh_token', 'mock-refresh-token')
    localStorage.setItem(
      'auth_user',
      JSON.stringify({
        id: 'usr_mock',
        tenant_id: 'tnt_mock',
        email: 'demo@hada.commerce',
        name: 'Demo User',
        status: 'active',
        scopes: ['admin'],
      }),
    )
    localStorage.setItem(
      'auth_tenant',
      JSON.stringify({
        id: 'tnt_mock',
        name: 'Demo Store',
        slug: 'demo-store',
        plan: 'pro',
        is_active: true,
      }),
    )
  })

  // Intercept all API calls so pages render even without a backend
  await page.route('**/api/v1/**', (route) => {
    const url = route.request().url()
    const json = (body: unknown) =>
      route.fulfill({ status: 200, contentType: 'application/json', body: JSON.stringify(body) })

    if (url.includes('/auth/me')) {
      return json({
        user: { id: 'usr_mock', tenant_id: 'tnt_mock', email: 'demo@hada.commerce', name: 'Demo User', status: 'active', scopes: ['admin'] },
        tenant: { id: 'tnt_mock', name: 'Demo Store', slug: 'demo-store', plan: 'pro', is_active: true },
      })
    }
    if (url.includes('/analytics/dashboard')) {
      return json({ total_products: 42, total_orders: 156, total_customers: 89, total_revenue: 1284500, currency: 'USD', pending_orders: 12, active_promos: 3, pending_pages: 2 })
    }
    if (url.includes('/analytics/revenue')) {
      const points = Array.from({ length: 30 }, (_, i) => {
        const d = new Date(); d.setDate(d.getDate() - (29 - i))
        return { date: d.toISOString().slice(0, 10), amount: Math.floor(Math.random() * 80000) + 20000 }
      })
      return json(points)
    }
    if (url.includes('/analytics/top-products')) {
      return json([
        { product_id: 'p1', name: 'Vitamin C Serum', total_sold: 245, revenue: 489500 },
        { product_id: 'p2', name: 'Hydrating Cream', total_sold: 198, revenue: 395000 },
        { product_id: 'p3', name: 'Retinol Night Oil', total_sold: 156, revenue: 467400 },
        { product_id: 'p4', name: 'SPF 50 Sunscreen', total_sold: 134, revenue: 267200 },
        { product_id: 'p5', name: 'Cleansing Balm', total_sold: 112, revenue: 223600 },
      ])
    }
    if (url.includes('/analytics/order-status')) {
      return json([
        { status: 'pending', count: 12 },
        { status: 'processing', count: 8 },
        { status: 'shipped', count: 24 },
        { status: 'delivered', count: 98 },
        { status: 'cancelled', count: 5 },
      ])
    }
    if (url.includes('/analytics/recent-orders')) {
      return json([
        { id: 'o1', customer_name: 'Sarah Chen', total: 12900, currency: 'USD', status: 'pending', created_at: new Date().toISOString() },
        { id: 'o2', customer_name: 'James Wilson', total: 8500, currency: 'USD', status: 'shipped', created_at: new Date().toISOString() },
        { id: 'o3', customer_name: 'Maria Garcia', total: 24700, currency: 'USD', status: 'delivered', created_at: new Date().toISOString() },
        { id: 'o4', customer_name: 'Alex Kim', total: 5600, currency: 'USD', status: 'processing', created_at: new Date().toISOString() },
        { id: 'o5', customer_name: 'Emma Brown', total: 18300, currency: 'USD', status: 'pending', created_at: new Date().toISOString() },
      ])
    }
    if (url.includes('/plugins/installed')) {
      return json({ items: [], total: 0, page: 1, page_size: 20, total_pages: 0 })
    }
    if (url.match(/\/plugins\/?(\?.*)?$/) && !url.includes('/settings') && !url.includes('/installed')) {
      return json({ items: [], total: 0, page: 1, page_size: 20, total_pages: 0 })
    }
    // Settings endpoint (match /settings but not /plugins/.../settings)
    if (url.match(/\/settings\/?(\?.*)?$/) && !url.includes('/plugins/')) {
      return json({
        tenant_id: 'tnt_mock', store_name: 'Hada Store', store_email: 'hello@hada.store',
        store_phone: '+1 555-0123', currency: 'USD', timezone: 'America/New_York',
        address: { street: '123 Commerce St', city: 'San Francisco', state: 'CA', zip: '94102', country: 'US' },
        logo_url: '', favicon_url: '',
        social_links: { instagram: '', twitter: '', facebook: '', tiktok: '' },
        checkout_config: { require_phone: false, require_shipping: true, allow_notes: true },
        updated_at: new Date().toISOString(),
      })
    }
    // Default: return empty paginated results for list endpoints
    return json({ items: [], total: 0, page: 1, page_size: 20, total_pages: 0 })
  })
})

for (const p of pages) {
  test(`screenshot: ${p.title}`, async ({ page }) => {
    await page.goto(p.path, { waitUntil: 'networkidle', timeout: 10000 }).catch(() => {
      // Page may error if no backend — that's fine, screenshot what we have
    })
    // Wait a moment for any client-side rendering
    await page.waitForTimeout(1000)
    await page.screenshot({
      path: path.join(SCREENSHOT_DIR, `${p.name}.png`),
      fullPage: true,
    })
  })
}

// After all screenshots, generate HTML gallery
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
  .join('\n')}
  </div>
</body>
</html>`

  fs.writeFileSync(path.join(SCREENSHOT_DIR, 'index.html'), html)
})
