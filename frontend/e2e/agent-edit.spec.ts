import { test, expect } from '@playwright/test'
import { Client } from 'pg'

const API = 'http://localhost:8080'
const FASHION_STORE = 'http://fashion.localhost:5173'
const TENANT_ID = 'tnt_fashion'
const ADMIN_EMAIL = 'admin@urbanthreads.co'

// Unique marker to verify the agent actually edited the page
const EDIT_MARKER = `Agent E2E Test ${Date.now()}`

/**
 * Helper: authenticate as the fashion tenant admin via OTP.
 * Returns a JWT access token.
 */
async function loginAsAdmin(): Promise<string> {
  // 1. Initiate login — triggers OTP generation
  const initRes = await fetch(`${API}/auth/passwordless/login/initiate`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: ADMIN_EMAIL, tenant_id: TENANT_ID }),
  })
  expect(initRes.ok, `initiate login failed: ${initRes.status}`).toBeTruthy()

  // 2. Grab the OTP code directly from the database
  const db = new Client({
    host: 'localhost',
    port: 5433,
    user: 'hada',
    password: 'hada',
    database: 'hada',
  })
  await db.connect()
  const { rows } = await db.query(
    `SELECT code FROM otps WHERE contact = $1 AND verified_at IS NULL ORDER BY created_at DESC LIMIT 1`,
    [ADMIN_EMAIL],
  )
  await db.end()
  expect(rows.length).toBeGreaterThan(0)
  const otpCode = rows[0].code

  // 3. Verify OTP — returns JWT tokens
  const verifyRes = await fetch(`${API}/auth/passwordless/login/verify`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ email: ADMIN_EMAIL, code: otpCode, tenant_id: TENANT_ID }),
  })
  expect(verifyRes.ok, `verify login failed: ${verifyRes.status}`).toBeTruthy()
  const tokens = await verifyRes.json()
  expect(tokens.access_token).toBeTruthy()
  return tokens.access_token
}

/**
 * Helper: send a message to the agent chat and collect the full SSE response.
 * Waits for the turn_end event.
 */
async function sendAgentMessage(token: string, message: string): Promise<string> {
  const res = await fetch(`${API}/api/v1/agent/chat`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${token}`,
    },
    body: JSON.stringify({ message, session_id: `e2e-test-${Date.now()}` }),
  })
  expect(res.ok, `agent chat failed: ${res.status}`).toBeTruthy()

  // Read SSE stream
  const reader = res.body!.getReader()
  const decoder = new TextDecoder()
  let fullText = ''
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
        if (evt.kind === 'text_delta' && evt.text) {
          fullText += evt.text
        }
        if (evt.kind === 'turn_end') {
          done = true
          break
        }
        if (evt.kind === 'error') {
          throw new Error(`Agent error: ${evt.error}`)
        }
      } catch (e) {
        if (e instanceof SyntaxError) continue // partial JSON
        throw e
      }
    }
  }

  return fullText
}

test.describe('AI Agent E2E — page editing', () => {
  test.setTimeout(120_000) // Agent calls can take time (LLM + Docker)

  let token: string

  test.beforeAll(async () => {
    token = await loginAsAdmin()
  })

  test('agent edits a CMS page and changes appear on the storefront', async ({ page }) => {
    // 1. Ask the agent to update the About page with our unique marker
    const prompt = `First use list_pages to find the page with slug "about", then use update_page to update its HTML. Add this exact text at the very beginning of the html: <h2>${EDIT_MARKER}</h2> — followed by the existing content. Do not remove the existing content, just prepend that h2 tag.`

    const response = await sendAgentMessage(token, prompt)
    // Verify agent acknowledged the edit (it should mention success or the tool)
    expect(response.toLowerCase()).toMatch(/updated|success|done|page/i)

    // 2. Navigate to the storefront About page and verify the marker is present
    await page.goto(`${FASHION_STORE}/pages/about`, { waitUntil: 'networkidle' })
    await expect(page.locator('body')).toContainText(EDIT_MARKER, { timeout: 15_000 })

    // 3. Screenshot for evidence
    await page.screenshot({ path: 'e2e/screenshots/agent-edit-verified.png', fullPage: true })
  })
})
