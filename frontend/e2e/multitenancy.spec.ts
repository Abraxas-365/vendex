import { test, expect } from '@playwright/test'

test.describe('Multi-tenancy storefront', () => {
  test('fashion.localhost shows Urban Threads branding', async ({ page }) => {
    // Intercept the /store/info call to see what the backend actually returns
    const infoResponses: any[] = []
    page.on('response', async (resp) => {
      if (resp.url().includes('/store/info')) {
        const body = await resp.json().catch(() => null)
        infoResponses.push({ url: resp.url(), status: resp.status(), body })
      }
    })

    await page.goto('http://fashion.localhost:5173/', { waitUntil: 'networkidle' })

    console.log('=== fashion.localhost /store/info responses ===')
    console.log(JSON.stringify(infoResponses, null, 2))

    // Check the hero title is "Elevate Your Style" (fashion store)
    const heroTitle = page.locator('h1')
    await expect(heroTitle).toContainText('Elevate Your Style', { timeout: 10000 })

    await page.screenshot({ path: 'e2e/screenshots/fashion-store.png', fullPage: true })
  })

  test('vendex-demo-store.localhost shows Vendex Tech branding', async ({ page }) => {
    const infoResponses: any[] = []
    page.on('response', async (resp) => {
      if (resp.url().includes('/store/info')) {
        const body = await resp.json().catch(() => null)
        infoResponses.push({ url: resp.url(), status: resp.status(), body })
      }
    })

    await page.goto('http://vendex-demo-store.localhost:5173/', { waitUntil: 'networkidle' })

    console.log('=== vendex-demo-store.localhost /store/info responses ===')
    console.log(JSON.stringify(infoResponses, null, 2))

    // Check the hero title is "Welcome to Vendex Demo" (tech store)
    const heroTitle = page.locator('h1')
    await expect(heroTitle).toContainText('Welcome to Vendex Demo', { timeout: 10000 })

    await page.screenshot({ path: 'e2e/screenshots/tech-store.png', fullPage: true })
  })

  test('both stores show different products', async ({ page }) => {
    // Fashion store
    await page.goto('http://fashion.localhost:5173/products', { waitUntil: 'networkidle' })
    const fashionProducts = await page.locator('[data-testid="product-card"], .product-card, a[href*="/products/"]').count()
    const fashionText = await page.textContent('body')

    await page.screenshot({ path: 'e2e/screenshots/fashion-products.png', fullPage: true })

    // Tech store
    await page.goto('http://vendex-demo-store.localhost:5173/products', { waitUntil: 'networkidle' })
    const techProducts = await page.locator('[data-testid="product-card"], .product-card, a[href*="/products/"]').count()
    const techText = await page.textContent('body')

    await page.screenshot({ path: 'e2e/screenshots/tech-products.png', fullPage: true })

    console.log(`Fashion products count: ${fashionProducts}`)
    console.log(`Tech products count: ${techProducts}`)
    console.log(`Fashion has "Silk Midi Dress": ${fashionText?.includes('Silk Midi Dress')}`)
    console.log(`Tech has "iPhone": ${techText?.includes('iPhone')}`)

    // They should be different
    expect(fashionText).toContain('Silk Midi Dress')
    expect(techText).toContain('iPhone')
  })
})
