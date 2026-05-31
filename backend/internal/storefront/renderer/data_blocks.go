// Package renderer — data_blocks.go renders product_grid and featured_collection
// blocks by fetching real data from the product and catalog services.
package renderer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"

	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/product"
)

// ── product_grid ──────────────────────────────────────────────────────────────

// renderProductGridBlock renders a product grid with real product data.
// If no ProductLister is configured it falls back to the static placeholder.
func (r *Renderer) renderProductGridBlock(ctx context.Context, tenantID kernel.TenantID, settings json.RawMessage) (string, error) {
	if r.productLister == nil {
		return renderBlock("product_grid", settings)
	}

	// Parse settings
	cfg := struct {
		Title      string `json:"title"`
		CollID     string `json:"collection_id"`
		CategoryID string `json:"category_id"`
		Columns    string `json:"columns"`
		Limit      int    `json:"limit"`
	}{Columns: "4", Limit: 8}
	if len(settings) > 0 {
		_ = json.Unmarshal(settings, &cfg)
	}
	if cfg.Limit <= 0 || cfg.Limit > 100 {
		cfg.Limit = 8
	}

	pg := kernel.NewPaginationOptions(1, cfg.Limit)

	var products []product.Product
	if cfg.CategoryID != "" {
		res, err := r.productLister.ListByCategory(ctx, tenantID, kernel.CategoryID(cfg.CategoryID), pg)
		if err == nil {
			products = res.Items
		}
	} else {
		res, err := r.productLister.List(ctx, tenantID, pg)
		if err == nil {
			products = res.Items
		}
	}

	return buildProductGrid(cfg.Title, cfg.Columns, products), nil
}

// buildProductGrid produces the product grid HTML from a list of products.
func buildProductGrid(title, columns string, products []product.Product) string {
	if columns == "" {
		columns = "4"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf(`<section class="section product-grid-section">
  <div class="container">
`))
	if title != "" {
		sb.WriteString(fmt.Sprintf(`    <h2 style="margin-bottom: 2rem; text-align: center;">%s</h2>
`, template.HTMLEscapeString(title)))
	}
	sb.WriteString(fmt.Sprintf(`    <div class="product-grid" style="
      display: grid;
      grid-template-columns: repeat(%s, 1fr);
      gap: calc(var(--spacing-unit) * 6);
    ">
`, template.HTMLEscapeString(columns)))

	if len(products) == 0 {
		sb.WriteString(`      <div style="grid-column: 1 / -1; padding: 3rem; text-align: center; color: var(--color-text-muted);">
        <p>No products available.</p>
      </div>
`)
	} else {
		for _, p := range products {
			sb.WriteString(buildProductCard(p))
		}
	}

	sb.WriteString(`    </div>
  </div>
</section>`)
	return sb.String()
}

// ── featured_collection ───────────────────────────────────────────────────────

// renderFeaturedCollectionBlock renders a featured collection with real collection + product data.
// Falls back to static placeholder if no CollectionGetter is configured.
func (r *Renderer) renderFeaturedCollectionBlock(ctx context.Context, tenantID kernel.TenantID, settings json.RawMessage) (string, error) {
	if r.collectionGetter == nil {
		return renderBlock("featured_collection", settings)
	}

	cfg := struct {
		Title string `json:"title"`
		CollID string `json:"collection_id"`
		Columns string `json:"columns"`
		Limit  int    `json:"limit"`
	}{Columns: "4", Limit: 8}
	if len(settings) > 0 {
		_ = json.Unmarshal(settings, &cfg)
	}
	if cfg.Limit <= 0 || cfg.Limit > 100 {
		cfg.Limit = 8
	}

	if cfg.CollID == "" {
		// No collection configured — render static placeholder
		return renderBlock("featured_collection", settings)
	}

	coll, err := r.collectionGetter.GetCollectionByID(ctx, tenantID, kernel.CollectionID(cfg.CollID))
	if err != nil || coll == nil {
		// Collection not found — fall back gracefully
		return renderBlock("featured_collection", settings)
	}

	// Use collection title if none provided in settings
	title := cfg.Title
	if title == "" {
		title = coll.Name
	}

	// Fetch products for the collection
	var products []product.Product
	if r.productLister != nil && len(coll.ProductIDs) > 0 {
		limit := cfg.Limit
		if limit > len(coll.ProductIDs) {
			limit = len(coll.ProductIDs)
		}
		pg := kernel.NewPaginationOptions(1, limit)
		res, err := r.productLister.List(ctx, tenantID, pg)
		if err == nil {
			// Filter to collection products
			collSet := make(map[kernel.ProductID]bool, len(coll.ProductIDs))
			for _, pid := range coll.ProductIDs {
				collSet[pid] = true
			}
			for _, p := range res.Items {
				if collSet[p.ID] {
					products = append(products, p)
				}
				if len(products) >= cfg.Limit {
					break
				}
			}
		}
	}

	columns := cfg.Columns
	if columns == "" {
		columns = "4"
	}

	var sb strings.Builder
	sb.WriteString(`<section class="section featured-collection-section">
  <div class="container">
`)
	if title != "" {
		sb.WriteString(fmt.Sprintf(`    <h2 style="margin-bottom: 2rem; text-align: center;">%s</h2>
`, template.HTMLEscapeString(title)))
	}
	if coll.Description != "" {
		sb.WriteString(fmt.Sprintf(`    <p style="text-align: center; color: var(--color-text-muted); margin-bottom: 2rem;">%s</p>
`, template.HTMLEscapeString(coll.Description)))
	}
	sb.WriteString(fmt.Sprintf(`    <div class="product-grid" style="
      display: grid;
      grid-template-columns: repeat(%s, 1fr);
      gap: calc(var(--spacing-unit) * 6);
    ">
`, template.HTMLEscapeString(columns)))

	if len(products) == 0 {
		sb.WriteString(`      <div style="grid-column: 1 / -1; padding: 3rem; text-align: center; color: var(--color-text-muted);">
        <p>No products in this collection yet.</p>
      </div>
`)
	} else {
		for _, p := range products {
			sb.WriteString(buildProductCard(p))
		}
	}

	sb.WriteString(`    </div>
  </div>
</section>`)
	return sb.String(), nil
}

// ── Shared product card ───────────────────────────────────────────────────────

// buildProductCard produces an HTML card for a single product.
func buildProductCard(p product.Product) string {
	imageHTML := `<div style="
        background: var(--color-surface);
        height: 200px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--color-text-muted);
        font-size: 0.8rem;
      ">No image</div>`
	if len(p.Images) > 0 && p.Images[0] != "" {
		imageHTML = fmt.Sprintf(`<img
        src="%s"
        alt="%s"
        style="width: 100%%; height: 200px; object-fit: cover; display: block;"
        loading="lazy"
      />`,
			template.HTMLEscapeString(p.Images[0]),
			template.HTMLEscapeString(p.Name),
		)
	}

	price := fmt.Sprintf("%s %.2f", p.Price.Currency, float64(p.Price.Amount)/100.0)
	productURL := fmt.Sprintf("/products/%s", string(p.ID))

	return fmt.Sprintf(`      <div class="product-card" style="
        background: var(--color-background);
        border: 1px solid var(--color-border);
        border-radius: var(--border-radius-lg);
        overflow: hidden;
        box-shadow: var(--shadow-sm);
        transition: box-shadow 0.2s ease;
      ">
        <a href="%s" style="display: block; text-decoration: none; color: inherit;">
          %s
          <div style="padding: 1rem;">
            <h3 style="margin: 0 0 0.5rem; font-size: 1rem; font-weight: 600; color: var(--color-text);">%s</h3>
            <p style="margin: 0; font-size: 0.95rem; font-weight: 700; color: var(--color-primary);">%s</p>
          </div>
        </a>
      </div>
`,
		template.HTMLEscapeString(productURL),
		imageHTML,
		template.HTMLEscapeString(p.Name),
		template.HTMLEscapeString(price),
	)
}

// ── Template overrides ────────────────────────────────────────────────────────

// renderBlockWithOverride renders a block using a tenant-supplied Go template string.
// The template receives the parsed settings map[string]interface{} as data.
func renderBlockWithOverride(blockType string, settings json.RawMessage, tmplSrc string) (string, error) {
	t, err := template.New(blockType + "_override").Funcs(templateFuncs).Parse(tmplSrc)
	if err != nil {
		// Invalid override template — fall through to built-in
		return renderBlock(blockType, settings)
	}

	data := make(map[string]interface{})
	if len(settings) > 0 {
		if err := json.Unmarshal(settings, &data); err != nil {
			return "", fmt.Errorf("renderer: failed to parse settings for override block %q: %w", blockType, err)
		}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		// Execution failure — fall back to built-in
		return renderBlock(blockType, settings)
	}
	return buf.String(), nil
}
