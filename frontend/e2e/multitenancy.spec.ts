import { test, expect } from '@playwright/test'

const FASHION = 'http://fashion.localhost:5173'
const TECH = 'http://vendex-demo-store.localhost:5173'

test.describe('Multi-tenancy storefront', () => {
  // ── Home page branding ────────────────────────────────────────────────────

  test('fashion store: home shows Urban Threads branding with rose accent', async ({ page }) => {
    await page.goto(`${FASHION}/`, { waitUntil: 'networkidle' })

    // Navbar shows store name
    await expect(page.locator('nav')).toContainText('Urban Threads')

    // Hero
    await expect(page.locator('h1')).toContainText('Elevate Your Style', { timeout: 10000 })

    // Rose accent (#be185d) should appear somewhere
    const accentEl = page.locator('[style*="#be185d"]').first()
    await expect(accentEl).toBeVisible()

    await page.screenshot({ path: 'e2e/screenshots/fashion-home.png', fullPage: true })
  })

  test('tech store: home shows Vendex Tech branding with purple accent', async ({ page }) => {
    await page.goto(`${TECH}/`, { waitUntil: 'networkidle' })

    await expect(page.locator('nav')).toContainText('Vendex Tech')
    await expect(page.locator('h1')).toContainText('Welcome to Vendex Demo', { timeout: 10000 })

    // Purple accent (#6366f1)
    const accentEl = page.locator('[style*="#6366f1"]').first()
    await expect(accentEl).toBeVisible()

    await page.screenshot({ path: 'e2e/screenshots/tech-home.png', fullPage: true })
  })

  // ── Product list ──────────────────────────────────────────────────────────

  test('fashion store: products page shows fashion items', async ({ page }) => {
    await page.goto(`${FASHION}/products`, { waitUntil: 'networkidle' })
    const body = await page.textContent('body')

    expect(body).toContain('Silk Midi Dress')
    expect(body).not.toContain('iPhone')

    await page.screenshot({ path: 'e2e/screenshots/fashion-products.png', fullPage: true })
  })

  test('tech store: products page shows tech items', async ({ page }) => {
    await page.goto(`${TECH}/products`, { waitUntil: 'networkidle' })
    const body = await page.textContent('body')

    expect(body).toContain('iPhone')
    expect(body).not.toContain('Silk Midi Dress')

    await page.screenshot({ path: 'e2e/screenshots/tech-products.png', fullPage: true })
  })

  // ── CMS pages ─────────────────────────────────────────────────────────────

  test('fashion store: About page loads with Urban Threads content', async ({ page }) => {
    await page.goto(`${FASHION}/pages/about`, { waitUntil: 'networkidle' })
    const body = await page.textContent('body')

    expect(body).toContain('Urban Threads')
    expect(body).toContain('sustainable')

    await page.screenshot({ path: 'e2e/screenshots/fashion-about.png', fullPage: true })
  })

  test('tech store: About page loads with Vendex Tech content', async ({ page }) => {
    await page.goto(`${TECH}/pages/about`, { waitUntil: 'networkidle' })
    const body = await page.textContent('body')

    expect(body).toContain('Vendex Tech')
    expect(body).toContain('e-commerce')

    await page.screenshot({ path: 'e2e/screenshots/tech-about.png', fullPage: true })
  })

  test('fashion store: FAQ page loads', async ({ page }) => {
    await page.goto(`${FASHION}/pages/faq`, { waitUntil: 'networkidle' })
    const body = await page.textContent('body')

    expect(body).toContain('Frequently Asked Questions')
    expect(body).toContain('cashmere')

    await page.screenshot({ path: 'e2e/screenshots/fashion-faq.png', fullPage: true })
  })

  test('fashion store: Shipping & Returns page loads', async ({ page }) => {
    await page.goto(`${FASHION}/pages/shipping-returns`, { waitUntil: 'networkidle' })
    const body = await page.textContent('body')

    expect(body).toContain('Shipping & Returns')
    expect(body).toContain('14-day')

    await page.screenshot({ path: 'e2e/screenshots/fashion-shipping.png', fullPage: true })
  })

  // ── Footer ────────────────────────────────────────────────────────────────

  test('fashion store: footer shows store name and CMS links', async ({ page }) => {
    await page.goto(`${FASHION}/`, { waitUntil: 'networkidle' })
    const footer = page.locator('footer')

    await expect(footer).toContainText('Urban Threads')
    await expect(footer).toContainText('About')
    await expect(footer).toContainText('Contact')
    await expect(footer).toContainText('Privacy Policy')
    await expect(footer).toContainText('Shipping & Returns')

    // Footer links should navigate to CMS pages
    await footer.getByText('About').click()
    await page.waitForURL('**/pages/about')
    await expect(page.locator('body')).toContainText('Urban Threads')
  })

  test('tech store: footer shows Vendex Tech branding', async ({ page }) => {
    await page.goto(`${TECH}/`, { waitUntil: 'networkidle' })
    const footer = page.locator('footer')

    await expect(footer).toContainText('Vendex Tech')
    await expect(footer).toContainText('hello@vendex.ai')
  })
})
