// Package renderer converts block-based storefront pages into full HTML5 documents
// using the active theme's design tokens as CSS custom properties.
package renderer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/settings"
	"github.com/Abraxas-365/vendex/internal/storefront"
	"github.com/Abraxas-365/vendex/internal/theme"
)

// ThemeGetter resolves the active theme for a tenant.
type ThemeGetter interface {
	GetActiveTheme(ctx context.Context, tenantID kernel.TenantID) (*theme.Theme, error)
}

// Renderer renders storefront pages as full HTML5 documents.
type Renderer struct {
	themeGetter      ThemeGetter
	productLister    ProductLister
	collectionGetter CollectionGetter
	settingsGetter   SettingsGetter
	navRepo          NavMenuRepository
	overrideRepo     TemplateOverrideRepository
}

// Config holds optional dependencies for the Renderer.
// All fields are optional — the renderer degrades gracefully when nil.
type Config struct {
	ProductLister    ProductLister
	CollectionGetter CollectionGetter
	SettingsGetter   SettingsGetter
	NavRepo          NavMenuRepository
	OverrideRepo     TemplateOverrideRepository
}

// New creates a new Renderer backed by the given ThemeGetter (no data resolvers).
func New(themeGetter ThemeGetter) *Renderer {
	return &Renderer{themeGetter: themeGetter}
}

// NewWithConfig creates a new Renderer with full data-resolver support.
func NewWithConfig(themeGetter ThemeGetter, cfg Config) *Renderer {
	return &Renderer{
		themeGetter:      themeGetter,
		productLister:    cfg.ProductLister,
		collectionGetter: cfg.CollectionGetter,
		settingsGetter:   cfg.SettingsGetter,
		navRepo:          cfg.NavRepo,
		overrideRepo:     cfg.OverrideRepo,
	}
}

// RenderPage produces a complete HTML5 document for the given page.
//
// For ContentType=="html": wraps Page.HTML + Page.CSS in a themed layout.
// For ContentType=="blocks": iterates sections, renders each block template, wraps in layout.
//
// The tenant is resolved from page.TenantID. This method satisfies the PageRenderer
// interface used by storefrontapi.Handler.
func (r *Renderer) RenderPage(ctx context.Context, page *storefront.Page) (string, error) {
	if page == nil {
		return "", errx.New("page is nil", errx.TypeValidation)
	}

	tenantID := page.TenantID

	// 1. Resolve active theme (falls back to defaults on failure)
	tokens := theme.DefaultTokens()
	if t, err := r.themeGetter.GetActiveTheme(ctx, tenantID); err == nil && t != nil {
		tokens = t.Tokens
	}

	// 2. Generate CSS custom properties from theme tokens
	rootCSS := TokensToCSS(tokens)

	// 3. Fetch store settings for branding (logo, store name, social links)
	var storeSettings *settings.StoreSettings
	if r.settingsGetter != nil {
		if ss, err := r.settingsGetter.Get(ctx, tenantID); err == nil {
			storeSettings = ss
		}
		// Silently ignore settings errors — layout falls back to defaults
	}

	// 4. Fetch navigation menus
	var headerNav, footerNav []NavMenuItem
	if r.navRepo != nil {
		if items, err := r.navRepo.ListByLocation(ctx, tenantID, NavLocationHeader); err == nil {
			headerNav = items
		}
		if items, err := r.navRepo.ListByLocation(ctx, tenantID, NavLocationFooter); err == nil {
			footerNav = items
		}
	}

	// 5. Render the page body according to content type
	var bodyHTML string
	var err error

	switch page.ContentType {
	case storefront.ContentTypeBlocks:
		bodyHTML, err = r.renderSections(ctx, tenantID, page.Sections)
		if err != nil {
			return "", fmt.Errorf("renderer: failed to render sections: %w", err)
		}
	default:
		// ContentTypeHTML — use raw HTML field
		bodyHTML = page.HTML
	}

	// 6. Build store name and branding for layout
	storeName := page.Title
	var logoURL string
	var socialLinks settings.SocialLinks
	if storeSettings != nil {
		if storeSettings.StoreName != "" {
			storeName = storeSettings.StoreName
		}
		logoURL = storeSettings.LogoURL
		socialLinks = storeSettings.SocialLinks
	}

	// 7. Wrap in HTML5 layout with header + footer
	return buildLayout(layoutData{
		Title:       page.Title,
		Description: page.Meta.Description,
		OGTitle:     page.Meta.OGTitle,
		OGImage:     page.Meta.OGImage,
		RootCSS:     rootCSS,
		BaseCSS:     baseCSS,
		PageCSS:     page.CSS,
		BodyHTML:    bodyHTML,
		StoreName:   storeName,
		LogoURL:     logoURL,
		SocialLinks: socialLinks,
		HeaderNav:   headerNav,
		FooterNav:   footerNav,
		Year:        time.Now().Year(),
	})
}

// renderSections iterates each section and renders its blocks.
func (r *Renderer) renderSections(ctx context.Context, tenantID kernel.TenantID, sections []storefront.Section) (string, error) {
	var sb strings.Builder

	for _, section := range sections {
		if len(section.Blocks) > 0 {
			sb.WriteString(fmt.Sprintf(`<div class="section-wrapper" data-section-id="%s" data-section-type="%s">`,
				template.HTMLEscapeString(section.ID),
				template.HTMLEscapeString(section.Type),
			))
			for _, block := range section.Blocks {
				blockHTML, err := r.renderBlockCtx(ctx, tenantID, block.Type, block.Settings)
				if err != nil {
					return "", err
				}
				sb.WriteString(blockHTML)
			}
			sb.WriteString("</div>")
		} else {
			// Section itself is the renderable unit
			sectionHTML, err := r.renderBlockCtx(ctx, tenantID, section.Type, section.Settings)
			if err != nil {
				return "", err
			}
			sb.WriteString(sectionHTML)
		}
	}

	return sb.String(), nil
}

// renderBlockCtx renders a single block with context for live data fetching and
// per-tenant template overrides.
func (r *Renderer) renderBlockCtx(ctx context.Context, tenantID kernel.TenantID, blockType string, rawSettings json.RawMessage) (string, error) {
	// Check for tenant-specific template override first.
	if r.overrideRepo != nil {
		override, err := r.overrideRepo.GetByBlockType(ctx, tenantID, blockType)
		if err == nil && override != nil {
			return renderBlockWithOverride(blockType, rawSettings, override.Template)
		}
		// Ignore override lookup errors — fall through to built-in template
	}

	// Dispatch to data-enriched renderers for dynamic block types.
	switch blockType {
	case "product_grid":
		return r.renderProductGridBlock(ctx, tenantID, rawSettings)
	case "featured_collection":
		return r.renderFeaturedCollectionBlock(ctx, tenantID, rawSettings)
	case "testimonials":
		return renderTestimonialsBlock(rawSettings)
	default:
		return renderBlock(blockType, rawSettings)
	}
}

// ── Layout ────────────────────────────────────────────────────────────────────

// layoutData holds all fields needed to render the full HTML5 document.
type layoutData struct {
	Title       string
	Description string
	OGTitle     string
	OGImage     string
	RootCSS     string // :root { ... } block from theme tokens
	BaseCSS     string // base reset + utility styles
	PageCSS     string // per-page custom CSS from page.CSS
	BodyHTML    string // rendered section/block HTML
	// Branding + navigation
	StoreName   string
	LogoURL     string
	SocialLinks settings.SocialLinks
	HeaderNav   []NavMenuItem
	FooterNav   []NavMenuItem
	Year        int
}

// pageLayoutTmpl is the outer HTML5 document shell with sticky header and footer.
var pageLayoutTmpl = template.Must(template.New("layout").Funcs(template.FuncMap{
	"safeCSS":  func(s string) template.CSS { return template.CSS(s) },
	"safeHTML": func(s string) template.HTML { return template.HTML(s) }, // #nosec G203
}).Parse(`<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>{{.Title}}</title>
  {{- if .Description}}
  <meta name="description" content="{{.Description}}">
  {{- end}}
  {{- if .OGTitle}}
  <meta property="og:title" content="{{.OGTitle}}">
  {{- end}}
  {{- if .OGImage}}
  <meta property="og:image" content="{{.OGImage}}">
  {{- end}}
  <style>
    {{safeCSS .RootCSS}}
    {{safeCSS .BaseCSS}}
  </style>
  {{- if .PageCSS}}
  <style>
    {{safeCSS .PageCSS}}
  </style>
  {{- end}}
</head>
<body>

  <!-- ── Header ─────────────────────────────────────────────────────────────── -->
  <header class="site-header" style="
    background: var(--color-background);
    border-bottom: 1px solid var(--color-border);
    padding: 0 calc(var(--spacing-unit) * 4);
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 64px;
    position: sticky;
    top: 0;
    z-index: 100;
    box-shadow: var(--shadow-sm);
  ">
    <a class="site-logo" href="/" style="
      display: flex;
      align-items: center;
      gap: 0.75rem;
      text-decoration: none;
      color: var(--color-text);
      font-weight: 700;
      font-size: 1.25rem;
      font-family: var(--font-heading);
    ">
      {{- if .LogoURL}}
      <img src="{{.LogoURL}}" alt="{{.StoreName}} logo" style="height: 40px; width: auto; object-fit: contain;" />
      {{- end}}
      <span>{{.StoreName}}</span>
    </a>
    {{- if .HeaderNav}}
    <nav class="site-nav" aria-label="Main navigation" style="
      display: flex;
      align-items: center;
      gap: 0.25rem;
    ">
      {{- range .HeaderNav}}
      <a href="{{.URL}}" style="
        padding: 0.5rem 1rem;
        border-radius: var(--border-radius-md);
        color: var(--color-text);
        text-decoration: none;
        font-size: 0.9rem;
        font-weight: 500;
      ">{{.Label}}</a>
      {{- end}}
    </nav>
    {{- end}}
  </header>

  <!-- ── Page content ───────────────────────────────────────────────────────── -->
  <main>
    {{safeHTML .BodyHTML}}
  </main>

  <!-- ── Footer ─────────────────────────────────────────────────────────────── -->
  <footer class="site-footer" style="
    background: var(--color-surface);
    border-top: 1px solid var(--color-border);
    padding: calc(var(--spacing-unit) * 8) calc(var(--spacing-unit) * 4);
    margin-top: calc(var(--spacing-unit) * 8);
  ">
    <div class="container" style="
      display: flex;
      flex-direction: column;
      align-items: center;
      gap: 1.5rem;
      text-align: center;
    ">
      {{- if .FooterNav}}
      <nav aria-label="Footer navigation" style="display: flex; flex-wrap: wrap; justify-content: center; gap: 0.25rem;">
        {{- range .FooterNav}}
        <a href="{{.URL}}" style="
          padding: 0.25rem 0.75rem;
          color: var(--color-text-muted);
          text-decoration: none;
          font-size: 0.875rem;
        ">{{.Label}}</a>
        {{- end}}
      </nav>
      {{- end}}
      <div class="footer-social" style="display: flex; gap: 1rem; align-items: center;">
        {{- if .SocialLinks.Instagram}}
        <a href="{{.SocialLinks.Instagram}}" target="_blank" rel="noopener noreferrer"
           style="color: var(--color-text-muted); text-decoration: none; font-size: 0.875rem;">Instagram</a>
        {{- end}}
        {{- if .SocialLinks.Twitter}}
        <a href="{{.SocialLinks.Twitter}}" target="_blank" rel="noopener noreferrer"
           style="color: var(--color-text-muted); text-decoration: none; font-size: 0.875rem;">Twitter</a>
        {{- end}}
        {{- if .SocialLinks.Facebook}}
        <a href="{{.SocialLinks.Facebook}}" target="_blank" rel="noopener noreferrer"
           style="color: var(--color-text-muted); text-decoration: none; font-size: 0.875rem;">Facebook</a>
        {{- end}}
      </div>
      <p style="color: var(--color-text-muted); font-size: 0.8rem; margin: 0;">
        &copy; {{.Year}} {{.StoreName}}. All rights reserved.
      </p>
    </div>
  </footer>

</body>
</html>
`))

// buildLayout executes the page layout template and returns the final HTML string.
func buildLayout(data layoutData) (string, error) {
	var buf bytes.Buffer
	if err := pageLayoutTmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("renderer: failed to execute layout template: %w", err)
	}
	return buf.String(), nil
}
