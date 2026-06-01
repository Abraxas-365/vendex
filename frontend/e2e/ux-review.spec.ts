import { test } from '@playwright/test'
import * as fs from 'fs'
import * as path from 'path'
import { fileURLToPath } from 'url'

const __filename = fileURLToPath(import.meta.url)
const __dirname = path.dirname(__filename)
const SCREENSHOT_DIR = path.join(__dirname, '..', 'screenshots', 'ux-review')

// Ensure output directory exists
fs.mkdirSync(SCREENSHOT_DIR, { recursive: true })

const TENANTS = [
  { id: 'tnt_fashion', name: 'fashion', host: 'fashion.localhost' },
  { id: 'tnt_demo', name: 'tech', host: 'demo.localhost' },
]

const STORE_PAGES = [
  { name: 'home', path: '/', waitFor: 2000 },
  { name: 'products', path: '/products', waitFor: 2000 },
  { name: 'cart', path: '/cart', waitFor: 1000 },
  { name: 'about', path: '/pages/about', waitFor: 1000 },
  { name: 'contact', path: '/pages/contact', waitFor: 1000 },
  { name: 'faq', path: '/pages/faq', waitFor: 1000 },
  { name: 'shipping-returns', path: '/pages/shipping-returns', waitFor: 1000 },
]

for (const tenant of TENANTS) {
  test.describe(`${tenant.name} storefront`, () => {
    for (const page of STORE_PAGES) {
      test(`${tenant.name} - ${page.name}`, async ({ browser }) => {
        const context = await browser.newContext({
          extraHTTPHeaders: { 'X-Tenant-ID': tenant.id },
          viewport: { width: 1440, height: 900 },
        })
        const p = await context.newPage()
        await p.goto(`http://localhost:5173${page.path}`)
        await p.waitForTimeout(page.waitFor)
        await p.screenshot({
          path: path.join(SCREENSHOT_DIR, `${tenant.name}-${page.name}.png`),
          fullPage: true,
        })
        await context.close()
      })
    }

    // Product detail - grab first product link
    test(`${tenant.name} - product-detail`, async ({ browser }) => {
      const context = await browser.newContext({
        extraHTTPHeaders: { 'X-Tenant-ID': tenant.id },
        viewport: { width: 1440, height: 900 },
      })
      const p = await context.newPage()
      await p.goto('http://localhost:5173/products')
      await p.waitForTimeout(2000)
      // Click first product card
      const firstProduct = p.locator('a[href*="/products/"]').first()
      if (await firstProduct.isVisible()) {
        await firstProduct.click()
        await p.waitForTimeout(2000)
        await p.screenshot({
          path: path.join(SCREENSHOT_DIR, `${tenant.name}-product-detail.png`),
          fullPage: true,
        })
      }
      await context.close()
    })

    // Mobile viewport
    test(`${tenant.name} - mobile-home`, async ({ browser }) => {
      const context = await browser.newContext({
        extraHTTPHeaders: { 'X-Tenant-ID': tenant.id },
        viewport: { width: 390, height: 844 },
      })
      const p = await context.newPage()
      await p.goto('http://localhost:5173/')
      await p.waitForTimeout(2000)
      await p.screenshot({
        path: path.join(SCREENSHOT_DIR, `${tenant.name}-mobile-home.png`),
        fullPage: true,
      })
      await context.close()
    })

    test(`${tenant.name} - mobile-products`, async ({ browser }) => {
      const context = await browser.newContext({
        extraHTTPHeaders: { 'X-Tenant-ID': tenant.id },
        viewport: { width: 390, height: 844 },
      })
      const p = await context.newPage()
      await p.goto('http://localhost:5173/products')
      await p.waitForTimeout(2000)
      await p.screenshot({
        path: path.join(SCREENSHOT_DIR, `${tenant.name}-mobile-products.png`),
        fullPage: true,
      })
      await context.close()
    })
  })
}

// Admin pages (needs auth)
test.describe('admin pages', () => {
  const ADMIN_PAGES = [
    { name: 'dashboard', path: '/admin' },
    { name: 'products', path: '/admin/products' },
    { name: 'orders', path: '/admin/orders' },
    { name: 'customers', path: '/admin/customers' },
    { name: 'settings', path: '/admin/settings' },
    { name: 'inventory', path: '/admin/inventory' },
    { name: 'pages', path: '/admin/pages' },
    { name: 'theme', path: '/admin/theme' },
  ]

  for (const page of ADMIN_PAGES) {
    test(`admin - ${page.name}`, async ({ browser }) => {
      const context = await browser.newContext({
        extraHTTPHeaders: { 'X-Tenant-ID': 'tnt_fashion' },
        viewport: { width: 1440, height: 900 },
        storageState: undefined,
      })
      const p = await context.newPage()
      
      // Login via OTP
      await p.goto('http://localhost:5173/login')
      await p.waitForTimeout(1000)
      
      // Fill email and submit
      const emailInput = p.locator('input[type="email"]')
      if (await emailInput.isVisible()) {
        await emailInput.fill('admin@hada.test')
        await p.locator('button[type="submit"]').click()
        await p.waitForTimeout(1500)
        
        // Get OTP from DB
        const { execSync } = await import('child_process')
        const otp = execSync(
          `psql "postgres://hada:hada@localhost:5433/hada" -t -A -c "SELECT code FROM otp_codes WHERE email='admin@hada.test' ORDER BY created_at DESC LIMIT 1"`
        ).toString().trim()
        
        if (otp) {
          const otpInput = p.locator('input[name="code"], input[placeholder*="code"], input[type="text"]').first()
          if (await otpInput.isVisible()) {
            await otpInput.fill(otp)
            await p.locator('button[type="submit"]').click()
            await p.waitForTimeout(2000)
          }
        }
      }
      
      // Navigate to admin page
      await p.goto(`http://localhost:5173${page.path}`)
      await p.waitForTimeout(2000)
      await p.screenshot({
        path: path.join(SCREENSHOT_DIR, `admin-${page.name}.png`),
        fullPage: true,
      })
      await context.close()
    })
  }
})
