import { defineConfig } from '@playwright/test'

export default defineConfig({
  testDir: './e2e',
  timeout: 60_000,
  use: {
    baseURL: 'http://localhost:5173',
    viewport: { width: 1440, height: 900 },
    screenshot: 'on', // Capture screenshots on every test (pass or fail)
    trace: 'on-first-retry', // Full trace on retry for debugging
    video: 'retain-on-failure', // Record video but only keep on failure
  },
  // HTML report — run `npx playwright show-report` to open
  reporter: [
    ['html', { open: 'never', outputFolder: 'playwright-report' }],
    ['list'], // Console output during run
  ],
  webServer: {
    command: 'npm run dev',
    port: 5173,
    reuseExistingServer: true,
    timeout: 30_000,
  },
  projects: [
    {
      name: 'screenshots',
      use: { browserName: 'chromium' },
      testMatch: 'screenshots.spec.ts',
    },
    {
      name: 'multitenancy',
      use: { browserName: 'chromium' },
      testMatch: 'multitenancy.spec.ts',
      timeout: 60_000,
    },
    {
      name: 'agent',
      use: { browserName: 'chromium' },
      testMatch: 'agent-edit.spec.ts',
      timeout: 300_000, // 5 min — LLM inference + Docker provisioning
    },
    {
      name: 'admin-features',
      use: { browserName: 'chromium' },
      testMatch: 'admin-features.spec.ts',
      timeout: 60_000,
    },
    {
      name: 'ux-review',
      use: { browserName: 'chromium' },
      testMatch: 'ux-review.spec.ts',
      timeout: 90_000,
    },
  ],
})
