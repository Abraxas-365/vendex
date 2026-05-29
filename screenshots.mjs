import { chromium } from 'playwright';
import { mkdirSync } from 'fs';

const BASE = 'http://localhost:5174';
const DIR = './screenshots';

const pages = [
  { path: '/admin',              name: '01-admin-dashboard',   title: 'Admin Dashboard' },
  { path: '/admin/products',     name: '02-admin-products',    title: 'Products Management' },
  { path: '/admin/orders',       name: '03-admin-orders',      title: 'Orders Management' },
  { path: '/admin/pages',        name: '04-admin-pages',       title: 'CMS Pages' },
  { path: '/admin/promos',       name: '05-admin-promos',      title: 'Promo Codes' },
  { path: '/admin/media',        name: '06-admin-media',       title: 'Media Gallery' },
  { path: '/admin/agent',        name: '07-admin-agent',       title: 'AI Agent Chat' },
  { path: '/admin/marketplace',  name: '08-admin-marketplace', title: 'Plugin Marketplace' },
  { path: '/',                   name: '09-store-home',        title: 'Store Homepage' },
  { path: '/products',           name: '10-store-products',    title: 'Product Catalog' },
  { path: '/cart',               name: '11-store-cart',        title: 'Shopping Cart' },
  { path: '/checkout',           name: '12-store-checkout',    title: 'Checkout' },
];

mkdirSync(DIR, { recursive: true });

const browser = await chromium.launch();
const context = await browser.newContext({
  viewport: { width: 1440, height: 900 },
  extraHTTPHeaders: { 'X-Tenant-ID': 'default' },
});

for (const pg of pages) {
  const page = await context.newPage();
  console.log(`Capturing: ${pg.title} (${pg.path})`);
  try {
    await page.goto(`${BASE}${pg.path}`, { waitUntil: 'networkidle', timeout: 15000 });
    await page.waitForTimeout(1000); // let animations settle
    await page.screenshot({ path: `${DIR}/${pg.name}.png`, fullPage: true });
    console.log(`  -> ${DIR}/${pg.name}.png`);
  } catch (err) {
    console.error(`  ERROR on ${pg.path}: ${err.message}`);
    // take screenshot anyway to capture error state
    try {
      await page.screenshot({ path: `${DIR}/${pg.name}.png`, fullPage: true });
    } catch (_) {}
  }
  await page.close();
}

await browser.close();
console.log('\nDone! Screenshots saved to ./screenshots/');
