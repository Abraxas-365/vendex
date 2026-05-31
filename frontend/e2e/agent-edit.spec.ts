/**
 * ============================================================================
 * AI Agent E2E Test — Page Editing via Natural Language
 * ============================================================================
 *
 * This test validates the full AI agent pipeline end-to-end:
 *
 *   1. Authenticates as a tenant admin using OTP (passwordless login)
 *   2. Sends a natural language instruction to the AI agent
 *   3. The agent uses its tools (list_pages → update_page) to modify CMS content
 *   4. Navigates to the public storefront and verifies the edit is live
 *
 * Prerequisites:
 *   - Backend running on localhost:8080
 *   - PostgreSQL running on localhost:5433 (docker-compose)
 *   - Frontend dev server running on localhost:5173
 *   - Docker daemon running (for agent workspace provisioning)
 *   - Seeded admin user: admin@urbanthreads.co (tnt_fashion)
 *   - Seeded pages including "about" page for tnt_fashion
 *
 * Run:
 *   npx playwright test --project=agent
 *
 * HTML Report:
 *   npx playwright show-report
 *
 * ============================================================================
 */

import { test, expect, Page } from '@playwright/test'
import { Client } from 'pg'

// ---------------------------------------------------------------------------
// Configuration
// ---------------------------------------------------------------------------

const API = 'http://localhost:8080'
const FASHION_STORE = 'http://fashion.localhost:5173'
const TENANT_ID = 'tnt_fashion'
const ADMIN_EMAIL = 'admin@urbanthreads.co'

/** Unique marker so we can verify the agent's edit — timestamped to avoid stale cache hits */
const EDIT_MARKER = `Agent E2E Test ${Date.now()}`

/** Database connection config (matches docker-compose.yml) */
const DB_CONFIG = {
  host: 'localhost',
  port: 5433,
  user: 'hada',
  password: 'hada',
  database: 'hada',
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

/**
 * Authenticate as the fashion tenant admin via OTP passwordless flow.
 *
 * Flow:
 *   1. POST /auth/passwordless/login/initiate → triggers OTP email (console in dev)
 *   2. Query PostgreSQL directly to grab the OTP code
 *   3. POST /auth/passwordless/login/verify → returns JWT access_token
 *
 * @returns JWT access token string
 */
async function loginAsAdmin(): Promise<string> {
  console.log('[auth] Initiating OTP login for', ADMIN_EMAIL)

  // Step 1: Request OTP
  const initRes = await fetch(`${API}/auth/passwordless/login/initiate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: ADMIN_EMAIL, tenant_id: TENANT_ID }),
  })
  expect(initRes.ok, `Login initiate failed: HTTP ${initRes.status} — ${await initRes.text()}`).toBeTruthy()
  console.log('[auth] OTP requested successfully')

  // Step 2: Retrieve OTP code from the database
  const db = new Client(DB_CONFIG)
  await db.connect()
  const { rows } = await db.query(
    `SELECT code FROM otps
     WHERE contact = $1 AND verified_at IS NULL
     ORDER BY created_at DESC LIMIT 1`,
    [ADMIN_EMAIL],
  )
  await db.end()

  expect(rows.length, 'No OTP found in database — is the backend running?').toBeGreaterThan(0)
  const otpCode = rows[0].code
  console.log('[auth] OTP code retrieved from DB:', otpCode)

  // Step 3: Verify OTP and get JWT tokens
  const verifyRes = await fetch(`${API}/auth/passwordless/login/verify`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: ADMIN_EMAIL, code: otpCode, tenant_id: TENANT_ID }),
  })
  expect(verifyRes.ok, `Login verify failed: HTTP ${verifyRes.status} — ${await verifyRes.text()}`).toBeTruthy()

  const tokens = await verifyRes.json()
  expect(tokens.access_token, 'No access_token in login response').toBeTruthy()
  console.log('[auth] Login successful — JWT obtained')

  return tokens.access_token
}

/**
 * Send a natural language message to the AI agent and collect the full response.
 *
 * The agent chat endpoint streams SSE events:
 *   - text_delta: incremental text tokens
 *   - tool_start: agent is calling a tool
 *   - tool_end: tool returned a result
 *   - turn_end: agent finished responding
 *   - error: something went wrong
 *
 * This function collects all events and returns structured results.
 *
 * @param token - JWT access token
 * @param message - Natural language instruction for the agent
 * @returns Object with fullText (agent response), toolCalls (tools used), and raw events
 */
async function sendAgentMessage(
  token: string,
  message: string,
): Promise<{ fullText: string; toolCalls: string[]; events: any[] }> {
  console.log('[agent] Sending message:', message.slice(0, 100) + '...')

  const res = await fetch(`${API}/api/v1/agent/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({
      message,
      session_id: `e2e-test-${Date.now()}`,
    }),
  })
  expect(res.ok, `Agent chat failed: HTTP ${res.status} — ${await res.text()}`).toBeTruthy()

  // Parse the SSE stream
  const reader = res.body!.getReader()
  const decoder = new TextDecoder()
  let fullText = ''
  const toolCalls: string[] = []
  const events: any[] = []
  let done = false

  while (!done) {
    const { value, done: streamDone } = await reader.read()
    if (streamDone) break

    const chunk = decoder.decode(value, { stream: true })
    const lines = chunk.split('\n')

    for (const line of lines) {
      if (!line.startsWith('data: ')) continue
      const data = line.slice(6)
      try {
        const evt = JSON.parse(data)
        events.push(evt)

        switch (evt.kind) {
          case 'text_delta':
            if (evt.text) fullText += evt.text
            break
          case 'tool_start':
            console.log(`[agent]   → Tool call: ${evt.tool_name}`)
            toolCalls.push(evt.tool_name)
            break
          case 'tool_end':
            console.log(`[agent]   ← Tool result: ${evt.tool_name} (${evt.result?.length ?? 0} chars)`)
            break
          case 'turn_end':
            console.log('[agent] Turn complete')
            done = true
            break
          case 'error':
            console.error('[agent] ERROR:', evt.error)
            throw new Error(`Agent error: ${evt.error}`)
            break
        }
      } catch (e) {
        if (e instanceof SyntaxError) continue // partial JSON chunk
        throw e
      }
    }
  }

  console.log(`[agent] Response: ${fullText.slice(0, 200)}...`)
  console.log(`[agent] Tools used: ${toolCalls.join(', ') || 'none'}`)

  return { fullText, toolCalls, events }
}

/**
 * Take a named screenshot and log it.
 */
async function screenshot(page: Page, name: string, description: string) {
  const path = `e2e/screenshots/${name}.png`
  await page.screenshot({ path, fullPage: true })
  console.log(`[screenshot] ${description} → ${path}`)
}

// ---------------------------------------------------------------------------
// Test Suite
// ---------------------------------------------------------------------------

test.describe('AI Agent E2E — Page Editing', () => {
  // Agent needs time: LLM inference + tool execution + potential Docker provisioning
  test.setTimeout(300_000) // 5 minutes

  let token: string

  test.beforeAll(async () => {
    token = await loginAsAdmin()
  })

  test('agent edits a CMS page and the change appears on the storefront', async ({ page }) => {
    // =========================================================================
    // STEP 1: Verify the storefront page BEFORE the agent edit
    // =========================================================================
    console.log('\n━━━ STEP 1: Capture storefront state BEFORE agent edit ━━━')

    await page.goto(`${FASHION_STORE}/pages/about`, { waitUntil: 'networkidle', timeout: 30_000 })
    await screenshot(page, '01-before-agent-edit', 'About page before AI agent edits it')

    // Confirm our marker is NOT already there
    const bodyBefore = await page.locator('body').textContent()
    expect(bodyBefore).not.toContain(EDIT_MARKER)
    console.log('[verify] Confirmed: marker NOT present before edit')

    // =========================================================================
    // STEP 2: Send instruction to AI agent
    // =========================================================================
    console.log('\n━━━ STEP 2: Instruct AI agent to edit the About page ━━━')

    const prompt = [
      `I need you to edit the About page for our store.`,
      `First, use list_pages to find the page with slug "about".`,
      `Then use update_page to modify its HTML content.`,
      `Add this exact HTML at the very beginning of the content:`,
      `<h2 class="agent-marker">${EDIT_MARKER}</h2>`,
      `Keep all existing content below it — do not delete anything.`,
    ].join(' ')

    const { fullText, toolCalls } = await sendAgentMessage(token, prompt)

    // Verify the agent actually used the expected tools
    expect(toolCalls, 'Agent should have called list_pages').toContain('list_pages')
    expect(toolCalls, 'Agent should have called update_page').toContain('update_page')
    console.log('[verify] Agent used the correct tools: list_pages → update_page')

    // Verify agent response indicates success
    expect(
      fullText.toLowerCase(),
      'Agent response should indicate success',
    ).toMatch(/updated|success|done|modified|changed|page/i)
    console.log('[verify] Agent confirmed the page was updated')

    // =========================================================================
    // STEP 3: Verify the edit is live on the storefront
    // =========================================================================
    console.log('\n━━━ STEP 3: Verify edit is live on storefront ━━━')

    // Navigate to the About page again (fresh load)
    await page.goto(`${FASHION_STORE}/pages/about`, { waitUntil: 'networkidle', timeout: 30_000 })
    await screenshot(page, '02-after-agent-edit', 'About page AFTER AI agent edited it')

    // Assert the marker text is now present
    await expect(
      page.locator('body'),
      `Expected the agent's edit marker to be visible on the storefront`,
    ).toContainText(EDIT_MARKER, { timeout: 30_000 })
    console.log('[verify] SUCCESS — Agent edit is live on the storefront!')

    // Check the marker is in an h2 tag specifically
    const markerEl = page.locator(`h2:has-text("${EDIT_MARKER}")`)
    await expect(markerEl, 'Marker should be in an <h2> element').toBeVisible({ timeout: 10_000 })
    await screenshot(page, '03-marker-highlighted', 'Close-up showing the agent-added h2 marker')

    // =========================================================================
    // STEP 4: Navigate other pages to ensure no side effects
    // =========================================================================
    console.log('\n━━━ STEP 4: Verify no side effects on other pages ━━━')

    // Check Home page still loads correctly
    await page.goto(`${FASHION_STORE}/`, { waitUntil: 'networkidle', timeout: 30_000 })
    await screenshot(page, '04-home-after-edit', 'Home page — should be unaffected')
    await expect(page.locator('body')).not.toContainText(EDIT_MARKER)
    console.log('[verify] Home page unaffected ✓')

    // Check Products page still loads correctly
    await page.goto(`${FASHION_STORE}/products`, { waitUntil: 'networkidle', timeout: 30_000 })
    await screenshot(page, '05-products-after-edit', 'Products page — should be unaffected')
    await expect(page.locator('body')).not.toContainText(EDIT_MARKER)
    console.log('[verify] Products page unaffected ✓')

    console.log('\n━━━ ALL CHECKS PASSED ━━━')
  })
})
