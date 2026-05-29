import { chromium } from 'playwright';
import { mkdirSync, readFileSync, writeFileSync } from 'fs';

const BASE = 'http://localhost:5174';
const DIR = './screenshots';

const pages = [
  // Admin pages
  { path: '/admin',              name: '01-admin-dashboard',   title: 'Admin Dashboard' },
  { path: '/admin/products',     name: '02-admin-products',    title: 'Products Management' },
  { path: '/admin/catalog',      name: '03-admin-catalog',     title: 'Catalog (Categories & Collections)' },
  { path: '/admin/orders',       name: '04-admin-orders',      title: 'Orders Management' },
  { path: '/admin/customers',    name: '05-admin-customers',   title: 'Customers' },
  { path: '/admin/pages',        name: '06-admin-pages',       title: 'CMS Pages' },
  { path: '/admin/promos',       name: '07-admin-promos',      title: 'Promo Codes' },
  { path: '/admin/media',        name: '08-admin-media',       title: 'Media Gallery' },
  { path: '/admin/agent',        name: '09-admin-agent',       title: 'AI Agent Chat' },
  { path: '/admin/marketplace',  name: '10-admin-marketplace', title: 'Plugin Marketplace' },
  { path: '/admin/settings',     name: '11-admin-settings',    title: 'Settings' },
  // Storefront pages
  { path: '/',                   name: '12-store-home',        title: 'Store Homepage' },
  { path: '/products',           name: '13-store-products',    title: 'Product Catalog' },
  { path: '/cart',               name: '14-store-cart',        title: 'Shopping Cart' },
  { path: '/checkout',           name: '15-store-checkout',    title: 'Checkout' },
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
    await page.waitForTimeout(1500); // let animations + charts settle
    await page.screenshot({ path: `${DIR}/${pg.name}.png`, fullPage: true });
    console.log(`  -> ${DIR}/${pg.name}.png`);
  } catch (err) {
    console.error(`  ERROR on ${pg.path}: ${err.message}`);
    try {
      await page.screenshot({ path: `${DIR}/${pg.name}.png`, fullPage: true });
    } catch (_) {}
  }
  await page.close();
}

await browser.close();
console.log('\nGenerating HTML report...');

// Generate HTML report with embedded images
let html = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Hada Commerce — Screenshot Report</title>
  <style>
    * { margin: 0; padding: 0; box-sizing: border-box; }
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #f8fafc; color: #1e293b; }
    header { background: linear-gradient(135deg, #4f46e5, #7c3aed); color: white; padding: 3rem 2rem; text-align: center; }
    header h1 { font-size: 2rem; font-weight: 700; }
    header p { margin-top: 0.5rem; opacity: 0.85; font-size: 1.1rem; }
    .toc { max-width: 900px; margin: 2rem auto; padding: 1.5rem 2rem; background: white; border-radius: 12px; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
    .toc h2 { font-size: 1.1rem; margin-bottom: 1rem; color: #4f46e5; }
    .toc ul { list-style: none; display: grid; grid-template-columns: 1fr 1fr; gap: 0.4rem; }
    .toc a { text-decoration: none; color: #4f46e5; font-size: 0.95rem; padding: 0.3rem 0; display: block; }
    .toc a:hover { text-decoration: underline; }
    .screenshots { max-width: 1200px; margin: 0 auto; padding: 2rem; }
    .screenshot { margin-bottom: 3rem; background: white; border-radius: 12px; overflow: hidden; box-shadow: 0 1px 3px rgba(0,0,0,0.1); }
    .screenshot h2 { padding: 1rem 1.5rem; font-size: 1.2rem; border-bottom: 1px solid #e2e8f0; color: #1e293b; }
    .screenshot h2 span { color: #94a3b8; font-weight: 400; font-size: 0.9rem; margin-left: 0.5rem; }
    .screenshot img { width: 100%; display: block; }
    footer { text-align: center; padding: 2rem; color: #94a3b8; font-size: 0.85rem; }
  </style>
</head>
<body>
  <header>
    <h1>Hada Commerce — Screenshot Report</h1>
    <p>Generated ${new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric', hour: '2-digit', minute: '2-digit' })}</p>
  </header>
  <div class="toc">
    <h2>Table of Contents</h2>
    <ul>
`;

for (const pg of pages) {
  html += `      <li><a href="#${pg.name}">${pg.title}</a></li>\n`;
}

html += `    </ul>
  </div>
  <div class="screenshots">
`;

for (const pg of pages) {
  let imgSrc = '';
  try {
    const imgBuf = readFileSync(`${DIR}/${pg.name}.png`);
    imgSrc = `data:image/png;base64,${imgBuf.toString('base64')}`;
  } catch {
    imgSrc = '';
  }
  html += `    <div class="screenshot" id="${pg.name}">
      <h2>${pg.title}<span>${pg.path}</span></h2>
      ${imgSrc ? `<img src="${imgSrc}" alt="${pg.title}" loading="lazy" />` : '<p style="padding:2rem;color:#ef4444;">Screenshot not found</p>'}
    </div>
`;
}

html += `  </div>
  <footer>Hada Commerce &mdash; ${pages.length} pages captured</footer>
</body>
</html>`;

writeFileSync('report.html', html);
console.log(`\nDone! ${pages.length} screenshots captured.`);
console.log('Report saved to ./report.html');
