// Package renderer converts block-based storefront pages into full HTML5 documents
// using the active theme's design tokens as CSS custom properties.
package renderer

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/storefront"
	"github.com/Abraxas-365/hada-commerce/internal/theme"
)

// ThemeGetter resolves the active theme for a tenant.
type ThemeGetter interface {
	GetActiveTheme(ctx context.Context, tenantID kernel.TenantID) (*theme.Theme, error)
}

// Renderer renders storefront pages as full HTML5 documents.
type Renderer struct {
	themeGetter ThemeGetter
}

// New creates a new Renderer backed by the given ThemeGetter.
func New(themeGetter ThemeGetter) *Renderer {
	return &Renderer{themeGetter: themeGetter}
}

// RenderPage produces a complete HTML5 document for the given page.
//
// For ContentType=="html": wraps Page.HTML + Page.CSS in a themed layout.
// For ContentType=="blocks": iterates sections, renders each block template, wraps in layout.
func (r *Renderer) RenderPage(ctx context.Context, page *storefront.Page) (string, error) {
	if page == nil {
		return "", errx.New("page is nil", errx.TypeValidation)
	}

	// 1. Resolve active theme (falls back to defaults on failure)
	tokens := theme.DefaultTokens()
	if t, err := r.themeGetter.GetActiveTheme(ctx, page.TenantID); err == nil && t != nil {
		tokens = t.Tokens
	}

	// 2. Generate CSS custom properties from theme tokens
	rootCSS := TokensToCSS(tokens)

	// 3. Render the page body according to content type
	var bodyHTML string
	var err error

	switch page.ContentType {
	case storefront.ContentTypeBlocks:
		bodyHTML, err = renderSections(page.Sections)
		if err != nil {
			return "", fmt.Errorf("renderer: failed to render sections: %w", err)
		}
	default:
		// ContentTypeHTML — use raw HTML field
		bodyHTML = page.HTML
	}

	// 4. Wrap in HTML5 layout
	return buildLayout(layoutData{
		Title:       page.Title,
		Description: page.Meta.Description,
		OGTitle:     page.Meta.OGTitle,
		OGImage:     page.Meta.OGImage,
		RootCSS:     rootCSS,
		BaseCSS:     baseCSS,
		PageCSS:     page.CSS,
		BodyHTML:    bodyHTML,
	})
}

// renderSections iterates each section and renders its blocks.
func renderSections(sections []storefront.Section) (string, error) {
	var sb strings.Builder

	for _, section := range sections {
		// Render blocks nested inside the section (if any)
		if len(section.Blocks) > 0 {
			sb.WriteString(fmt.Sprintf(`<div class="section-wrapper" data-section-id="%s" data-section-type="%s">`,
				template.HTMLEscapeString(section.ID),
				template.HTMLEscapeString(section.Type),
			))
			for _, block := range section.Blocks {
				var blockHTML string
				var err error

				if block.Type == "testimonials" {
					blockHTML, err = renderTestimonialsBlock(block.Settings)
				} else {
					blockHTML, err = renderBlock(block.Type, block.Settings)
				}
				if err != nil {
					return "", err
				}
				sb.WriteString(blockHTML)
			}
			sb.WriteString("</div>")
		} else {
			// Section itself is the renderable unit
			var sectionHTML string
			var err error

			if section.Type == "testimonials" {
				sectionHTML, err = renderTestimonialsBlock(section.Settings)
			} else {
				sectionHTML, err = renderBlock(section.Type, section.Settings)
			}
			if err != nil {
				return "", err
			}
			sb.WriteString(sectionHTML)
		}
	}

	return sb.String(), nil
}

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
}

// pageLayoutTmpl is the outer HTML5 document shell.
var pageLayoutTmpl = template.Must(template.New("layout").Funcs(template.FuncMap{
	"safeCSS": func(s string) template.CSS { return template.CSS(s) },
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
  <main>
    {{safeHTML .BodyHTML}}
  </main>
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
