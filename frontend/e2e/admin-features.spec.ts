/**
 * ============================================================================
 * Admin Features E2E Test Suite
 * ============================================================================
 *
 * Comprehensive end-to-end tests for all major admin features:
 *   - Auth & Login (OTP passwordless)
 *   - Dashboard
 *   - Products
 *   - Orders
 *   - Pages (CMS)
 *   - Collections & Categories
 *   - Settings
 *   - Marketplace (plugins)
 *   - Preset Marketplace
 *
 * Prerequisites:
 *   - Backend running on localhost:8080
 *   - PostgreSQL running on localhost:5433 (docker-compose)
 *   - Frontend dev server running on localhost:5173
 *   - Seeded admin user: admin@vendex.ai (tnt_demo)
 *
 * Run:
 *   npx playwright test --project=admin-features
 *
 * ============================================================================
 */

import { test, expect, type Page, type BrowserContext } from '@playwright/test'
import { Client } from 'pg'

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const ADMIN_BASE = 'http://localhost:5173'
const API_BASE = 'http://localhost:8080'
const TENANT_ID = 'tnt_demo'
const ADMIN_EMAIL = 'admin@vendex.ai'

/** Database connection config (matches docker-compose.yml) */
const DB_CONFIG = {
  host: 'localhost',
  port: 5433,
  user: 'hada',
  password: 'hada',
  database: 'hada',
}

// ---------------------------------------------------------------------------
// Auth helper
// ---------------------------------------------------------------------------

/**
 * Retrieve the latest unused OTP code from the database for a given email.
 */
async function getOTPCode(email: string): Promise<string> {
  const client = new Client(DB_CONFIG)
  await client.connect()
  const result = await client.query(
    `SELECT code FROM otps
     WHERE contact = $1 AND verified_at IS NULL
     ORDER BY created_at DESC LIMIT 1`,
    [email],
  )
  await client.end()

  if (result.rows.length === 0) {
    throw new Error(`No OTP found in DB for ${email} — is the backend running?`)
  }
  return result.rows[0].code
}

/**
 * Authenticate via the OTP API directly (without going through the UI).
 * Returns { access_token, tenant_id } to inject into localStorage.
 */
async function getAuthTokens(): Promise<{ accessToken: string; tenantId: string }> {
  console.log('[auth] Initiating OTP login for', ADMIN_EMAIL)

  // Step 1: Request OTP
  const initRes = await fetch(`${API_BASE}/auth/passwordless/login/initiate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: ADMIN_EMAIL, tenant_id: TENANT_ID }),
  })

  if (!initRes.ok) {
    const text = await initRes.text()
    throw new Error(`Login initiate failed: HTTP ${initRes.status} — ${text}`)
  }
  console.log('[auth] OTP requested successfully')

  // Small delay to let the DB write the OTP
  await new Promise((r) => setTimeout(r, 500))

  // Step 2: Get OTP code from DB
  const otpCode = await getOTPCode(ADMIN_EMAIL)
  console.log('[auth] OTP code retrieved from DB:', otpCode)

  // Step 3: Verify OTP and get JWT tokens
  const verifyRes = await fetch(`${API_BASE}/auth/passwordless/login/verify`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: ADMIN_EMAIL, code: otpCode, tenant_id: TENANT_ID }),
  })
  const verifyText = await verifyRes.text()

  if (!verifyRes.ok) {
    throw new Error(`Login verify failed: HTTP ${verifyRes.status} — ${verifyText}`)
  }

  const tokens = JSON.parse(verifyText) as {
    access_token: string
    refresh_token: string
    tenant: { id: string }
  }

  if (!tokens.access_token) {
    throw new Error('No access_token in login response')
  }
  console.log('[auth] Login successful — JWT obtained')

  return {
    accessToken: tokens.access_token,
    tenantId: tokens.tenant?.id ?? TENANT_ID,
  }
}

/**
 * Inject auth tokens into the browser's localStorage so the React app
 * considers itself authenticated without going through the login UI.
 */
async function injectAuthState(
  context: BrowserContext,
  accessToken: string,
  tenantId: string,
): Promise<void> {
  // Open a blank page to set localStorage on the right origin
  const page = await context.newPage()
  await page.goto(`${ADMIN_BASE}/login`, { waitUntil: 'domcontentloaded' })
  await page.evaluate(
    ({ at, rt, tid }) => {
      localStorage.setItem('hada_access_token', at)
      localStorage.setItem('hada_refresh_token', rt)
      localStorage.setItem('hada_tenant_id', tid)
    },
    { at: accessToken, rt: accessToken, tid: tenantId },
  )
  await page.close()
  console.log('[auth] Auth state injected into browser localStorage')
}

/**
 * Take a named screenshot and log it.
 */
async function screenshot(page: Page, name: string, description: string): Promise<void> {
  const path = `e2e/screenshots/admin-features-${name}.png`
  await page.screenshot({ path, fullPage: true })
  console.log(`[screenshot] ${description} → ${path}`)
}

/**
 * Navigate to a path and wait for network to settle. Returns the page URL after navigation.
 */
async function navigate(page: Page, path: string): Promise<void> {
  await page.goto(`${ADMIN_BASE}${path}`, { waitUntil: 'networkidle', timeout: 30_000 })
}

/**
 * Verify no error page (404/500) is showing.
 */
async function assertNoErrorPage(page: Page, routeLabel: string): Promise<void> {
  // Check for common error indicators
  const bodyText = await page.locator('body').textContent()
  const hasNotFound =
    bodyText?.includes('404') && bodyText?.includes('not found')
  const hasServerError = bodyText?.includes('500') && bodyText?.includes('error')

  if (hasNotFound) {
    throw new Error(`${routeLabel}: Got 404 Not Found page`)
  }
  if (hasServerError) {
    throw new Error(`${routeLabel}: Got 500 Server Error page`)
  }
}

// ---------------------------------------------------------------------------
// Test Suite
// ---------------------------------------------------------------------------

test.describe.serial('Admin Features E2E', () => {
  test.setTimeout(60_000)

  let accessToken: string
  let tenantId: string

  // Authenticate once before all tests
  test.beforeAll(async () => {
    const tokens = await getAuthTokens()
    accessToken = tokens.accessToken
    tenantId = tokens.tenantId
  })

  // ---------------------------------------------------------------------------
  // 1. Auth & Login
  // ---------------------------------------------------------------------------

  test('1. Login via OTP and verify redirect to admin dashboard', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      // Inject auth tokens
      await injectAuthState(context, accessToken, tenantId)

      // Navigate to admin
      await navigate(page, '/admin')
      await screenshot(page, '01-dashboard', 'Admin dashboard after OTP login')

      // Should be on /admin (not /login)
      expect(page.url()).toContain('/admin')
      expect(page.url()).not.toContain('/login')

      // Admin layout should be visible — check for the sidebar
      await expect(page.locator('nav, [class*="sidebar"], [class*="w-60"]').first()).toBeVisible({
        timeout: 15_000,
      })
      console.log('[verify] Auth successful — admin layout visible')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 2. Dashboard
  // ---------------------------------------------------------------------------

  test('2. Dashboard loads with stats widgets', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)
      await navigate(page, '/admin')
      await screenshot(page, '02-dashboard', 'Admin dashboard')
      await assertNoErrorPage(page, '/admin')

      // The Dashboard page renders an h1 with "Dashboard"
      await expect(page.locator('h1').filter({ hasText: 'Dashboard' })).toBeVisible({
        timeout: 15_000,
      })

      // Stat cards: "Total Orders" and "Total Revenue" are rendered
      await expect(page.locator('body')).toContainText('Total Orders', { timeout: 15_000 })
      await expect(page.locator('body')).toContainText('Total Revenue', { timeout: 15_000 })
      console.log('[verify] Dashboard stats are visible ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 3. Products
  // ---------------------------------------------------------------------------

  test('3. Products list loads', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)

      // Watch for the API call to succeed
      const apiResponse = page.waitForResponse(
        (resp) =>
          resp.url().includes('/api/v1/products') &&
          resp.request().method() === 'GET',
        { timeout: 20_000 },
      )

      await navigate(page, '/admin/products')
      await screenshot(page, '03-products', 'Products list page')
      await assertNoErrorPage(page, '/admin/products')

      // Products page renders h1 "Products"
      await expect(page.locator('h1').filter({ hasText: 'Products' })).toBeVisible({
        timeout: 15_000,
      })

      // Verify the API call returned OK
      const resp = await apiResponse
      expect(resp.status(), `Products API returned ${resp.status()}`).toBeLessThan(400)
      console.log(`[verify] Products API returned ${resp.status()} ✓`)
      console.log('[verify] Products page loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 4. Orders
  // ---------------------------------------------------------------------------

  test('4. Orders list loads', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)

      const apiResponse = page.waitForResponse(
        (resp) =>
          resp.url().includes('/api/v1/orders') &&
          resp.request().method() === 'GET',
        { timeout: 20_000 },
      )

      await navigate(page, '/admin/orders')
      await screenshot(page, '04-orders', 'Orders list page')
      await assertNoErrorPage(page, '/admin/orders')

      // Orders page renders h1 "Orders"
      await expect(page.locator('h1').filter({ hasText: 'Orders' })).toBeVisible({
        timeout: 15_000,
      })

      const resp = await apiResponse
      expect(resp.status(), `Orders API returned ${resp.status()}`).toBeLessThan(400)
      console.log(`[verify] Orders API returned ${resp.status()} ✓`)
      console.log('[verify] Orders page loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 5. Pages (CMS)
  // ---------------------------------------------------------------------------

  test('5. Pages (CMS) list loads', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)

      const apiResponse = page.waitForResponse(
        (resp) =>
          resp.url().includes('/api/v1/pages') &&
          resp.request().method() === 'GET',
        { timeout: 20_000 },
      )

      await navigate(page, '/admin/pages')
      await screenshot(page, '05-pages', 'CMS Pages list')
      await assertNoErrorPage(page, '/admin/pages')

      // Pages page renders h1 "Pages"
      await expect(page.locator('h1').filter({ hasText: 'Pages' })).toBeVisible({
        timeout: 15_000,
      })

      const resp = await apiResponse
      expect(resp.status(), `Pages API returned ${resp.status()}`).toBeLessThan(400)
      console.log(`[verify] Pages API returned ${resp.status()} ✓`)
      console.log('[verify] Pages list loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 6. Collections
  // ---------------------------------------------------------------------------

  test('6. Collections list loads', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)

      const apiResponse = page.waitForResponse(
        (resp) =>
          resp.url().includes('/api/v1/collections') &&
          resp.request().method() === 'GET',
        { timeout: 20_000 },
      )

      await navigate(page, '/admin/collections')
      await screenshot(page, '06-collections', 'Collections list')
      await assertNoErrorPage(page, '/admin/collections')

      // Collections page renders h1 "Collections"
      await expect(page.locator('h1').filter({ hasText: 'Collections' })).toBeVisible({
        timeout: 15_000,
      })

      const resp = await apiResponse
      expect(resp.status(), `Collections API returned ${resp.status()}`).toBeLessThan(400)
      console.log(`[verify] Collections API returned ${resp.status()} ✓`)
      console.log('[verify] Collections page loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 7. Catalog (Categories)
  // ---------------------------------------------------------------------------

  test('7. Catalog (categories) page loads', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)
      await navigate(page, '/admin/catalog')
      await screenshot(page, '07-catalog', 'Catalog/Categories page')
      await assertNoErrorPage(page, '/admin/catalog')

      // Catalog page should be visible — check sidebar nav item is active
      // and the page content renders (no crash)
      const body = page.locator('body')
      await expect(body).toBeVisible({ timeout: 10_000 })

      // URL should still be /admin/catalog (no redirect to /login)
      expect(page.url()).toContain('/admin/catalog')
      console.log('[verify] Catalog page loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 8. Customers
  // ---------------------------------------------------------------------------

  test('8. Customers list loads', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)

      const apiResponse = page.waitForResponse(
        (resp) =>
          resp.url().includes('/api/v1/customers') &&
          resp.request().method() === 'GET',
        { timeout: 20_000 },
      )

      await navigate(page, '/admin/customers')
      await screenshot(page, '08-customers', 'Customers list')
      await assertNoErrorPage(page, '/admin/customers')

      await expect(page.locator('h1').filter({ hasText: 'Customers' })).toBeVisible({
        timeout: 15_000,
      })

      const resp = await apiResponse
      expect(resp.status(), `Customers API returned ${resp.status()}`).toBeLessThan(400)
      console.log(`[verify] Customers API returned ${resp.status()} ✓`)
      console.log('[verify] Customers page loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 9. Settings
  // ---------------------------------------------------------------------------

  test('9. Settings page loads', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)
      await navigate(page, '/admin/settings')
      await screenshot(page, '09-settings', 'Settings page')
      await assertNoErrorPage(page, '/admin/settings')

      // Settings page should be at the right URL
      expect(page.url()).toContain('/admin/settings')
      // Body should be visible with no crash
      await expect(page.locator('body')).toBeVisible({ timeout: 10_000 })
      console.log('[verify] Settings page loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 10. Marketplace (plugins) — the original 404 bug
  // ---------------------------------------------------------------------------

  test('10. Marketplace loads without 404 error', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)

      // Monitor network for the plugins API call
      const apiResponse = page.waitForResponse(
        (resp) =>
          resp.url().includes('/api/v1/plugins') &&
          resp.request().method() === 'GET',
        { timeout: 20_000 },
      )

      await navigate(page, '/admin/marketplace')
      await screenshot(page, '10-marketplace', 'Marketplace page')
      await assertNoErrorPage(page, '/admin/marketplace')

      // Marketplace page renders h1 "Marketplace"
      await expect(page.locator('h1').filter({ hasText: 'Marketplace' })).toBeVisible({
        timeout: 15_000,
      })

      // Verify API responded (not 404)
      const resp = await apiResponse
      expect(resp.status(), `Marketplace plugins API returned ${resp.status()} — expected non-404`).not.toBe(404)
      console.log(`[verify] Marketplace API returned ${resp.status()} ✓`)
      console.log('[verify] Marketplace loaded without 404 ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 11. Preset Marketplace
  // ---------------------------------------------------------------------------

  test('11. Preset Marketplace loads and displays presets', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)

      // Monitor network for the presets API call
      const apiResponse = page.waitForResponse(
        (resp) =>
          resp.url().includes('/api/v1/presets') &&
          resp.request().method() === 'GET',
        { timeout: 20_000 },
      )

      await navigate(page, '/admin/presets')
      await screenshot(page, '11-presets', 'Preset Marketplace page')
      await assertNoErrorPage(page, '/admin/presets')

      // Presets page renders h1 "Preset Marketplace"
      await expect(page.locator('h1').filter({ hasText: 'Preset Marketplace' })).toBeVisible({
        timeout: 15_000,
      })

      const resp = await apiResponse
      expect(resp.status(), `Presets API returned ${resp.status()} — expected non-404`).not.toBe(404)
      console.log(`[verify] Presets API returned ${resp.status()} ✓`)
      console.log('[verify] Preset Marketplace loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 12. Reviews
  // ---------------------------------------------------------------------------

  test('12. Reviews page loads', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)
      await navigate(page, '/admin/reviews')
      await screenshot(page, '12-reviews', 'Reviews page')
      await assertNoErrorPage(page, '/admin/reviews')

      await expect(page.locator('h1').filter({ hasText: 'Reviews' })).toBeVisible({
        timeout: 15_000,
      })
      console.log('[verify] Reviews page loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 13. Inventory
  // ---------------------------------------------------------------------------

  test('13. Inventory page loads', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      await injectAuthState(context, accessToken, tenantId)
      await navigate(page, '/admin/inventory')
      await screenshot(page, '13-inventory', 'Inventory page')
      await assertNoErrorPage(page, '/admin/inventory')

      expect(page.url()).toContain('/admin/inventory')
      await expect(page.locator('body')).toBeVisible({ timeout: 10_000 })
      console.log('[verify] Inventory page loaded ✓')
    } finally {
      await context.close()
    }
  })

  // ---------------------------------------------------------------------------
  // 14. Full login UI flow (smoke test)
  // ---------------------------------------------------------------------------

  test('14. OTP login UI flow smoke test', async ({ browser }) => {
    const context = await browser.newContext()
    const page = await context.newPage()

    try {
      // Get a fresh OTP
      const initRes = await fetch(`${API_BASE}/auth/passwordless/login/initiate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: ADMIN_EMAIL, tenant_id: TENANT_ID }),
      })
      expect(initRes.ok, `OTP initiate failed: ${initRes.status}`).toBeTruthy()
      await new Promise((r) => setTimeout(r, 500))
      const freshOtp = await getOTPCode(ADMIN_EMAIL)

      // Navigate to login page
      await page.goto(`${ADMIN_BASE}/login`, { waitUntil: 'networkidle', timeout: 20_000 })
      await screenshot(page, '14a-login-page', 'Login page')

      // Should show the email input
      const emailInput = page.locator('input[type="email"], input[placeholder*="email" i]').first()
      await expect(emailInput).toBeVisible({ timeout: 10_000 })

      // Fill in email and submit
      await emailInput.fill(ADMIN_EMAIL)
      const submitBtn = page.locator('button[type="submit"]').first()
      await submitBtn.click()
      await screenshot(page, '14b-after-email', 'After email submission')

      // After tenant selection (auto-selected if single tenant), should show code input
      // Wait for either tenant selector or OTP input
      const codeInput = page.locator('input[type="text"][maxlength="6"], input[inputmode="numeric"]').first()
      await expect(codeInput).toBeVisible({ timeout: 15_000 })
      await screenshot(page, '14c-otp-step', 'OTP code entry step')

      // Enter the OTP code
      await codeInput.fill(freshOtp)
      const verifyBtn = page.locator('button[type="submit"]').first()
      await verifyBtn.click()

      // Should redirect to admin after successful login
      await page.waitForURL('**/admin**', { timeout: 15_000 })
      await screenshot(page, '14d-after-login', 'Admin dashboard after UI login')

      expect(page.url()).toContain('/admin')
      console.log('[verify] OTP login UI flow works end-to-end ✓')
    } finally {
      await context.close()
    }
  })
})
