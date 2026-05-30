package renderer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"strings"
)

// blockTemplates maps block type names to their html/template definitions.
// All templates receive a map[string]interface{} as their data (parsed settings).
var blockTemplates = map[string]*template.Template{}

func init() {
	for name, src := range rawBlockTemplates {
		t, err := template.New(name).Funcs(templateFuncs).Parse(src)
		if err != nil {
			panic(fmt.Sprintf("renderer: failed to parse block template %q: %v", name, err))
		}
		blockTemplates[name] = t
	}
}

// templateFuncs provides helpers available to all block templates.
var templateFuncs = template.FuncMap{
	"safeHTML": func(s string) template.HTML {
		return template.HTML(s) // #nosec G203 — caller-controlled trusted HTML
	},
	"orDefault": func(val, def interface{}) interface{} {
		if val == nil {
			return def
		}
		if s, ok := val.(string); ok && s == "" {
			return def
		}
		return val
	},
	"strDefault": func(val interface{}, def string) string {
		if val == nil {
			return def
		}
		if s, ok := val.(string); ok && s != "" {
			return s
		}
		return def
	},
	"boolVal": func(val interface{}) bool {
		if val == nil {
			return false
		}
		if b, ok := val.(bool); ok {
			return b
		}
		return false
	},
}

// renderBlock renders a single block given its type name and Settings JSON.
// Returns an empty string if the block type is unknown (graceful degradation).
func renderBlock(blockType string, settings json.RawMessage) (string, error) {
	t, ok := blockTemplates[blockType]
	if !ok {
		// Unknown block type — render a placeholder comment
		return fmt.Sprintf("<!-- unknown block type: %s -->", template.HTMLEscapeString(blockType)), nil
	}

	// Parse settings into a map
	data := make(map[string]interface{})
	if len(settings) > 0 {
		if err := json.Unmarshal(settings, &data); err != nil {
			return "", fmt.Errorf("renderer: failed to parse settings for block %q: %w", blockType, err)
		}
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("renderer: failed to execute template for block %q: %w", blockType, err)
	}

	return buf.String(), nil
}

// rawBlockTemplates contains the Go html/template source for each built-in block type.
var rawBlockTemplates = map[string]string{

	// ── hero ──────────────────────────────────────────────────────────────────
	"hero": `<section class="section hero-section" style="
  {{- if strDefault .background_color "" }}background-color: {{strDefault .background_color ""}};{{end}}
  {{- if strDefault .background_image "" }}background-image: url('{{strDefault .background_image ""}}'); background-size: cover; background-position: center;{{end}}
  {{- if strDefault .text_color "" }}color: {{strDefault .text_color ""}};{{end}}
  text-align: {{strDefault .alignment "center"}};
">
  <div class="container">
    {{- if strDefault .heading "" }}
    <h1 class="hero-heading" style="font-size: 3rem; margin-bottom: 1rem; {{if strDefault .text_color ""}}color: {{strDefault .text_color ""}};{{end}}">
      {{strDefault .heading ""}}
    </h1>
    {{- end}}
    {{- if strDefault .subheading "" }}
    <p class="hero-subheading" style="font-size: 1.25rem; margin-bottom: 2rem; opacity: 0.85;">
      {{strDefault .subheading ""}}
    </p>
    {{- end}}
    {{- if strDefault .button_text "" }}
    <a href="{{strDefault .button_url "#"}}" class="btn btn-primary" style="font-size: 1.1rem; padding: 1rem 2.5rem; border-radius: var(--border-radius-full); box-shadow: var(--shadow-md);">
      {{strDefault .button_text "Learn More"}}
    </a>
    {{- end}}
  </div>
</section>`,

	// ── rich_text ─────────────────────────────────────────────────────────────
	"rich_text": `<section class="section rich-text-section">
  <div class="container">
    <div class="rich-text-content" style="max-width: 800px; margin: 0 auto;">
      {{safeHTML (strDefault .content "")}}
    </div>
  </div>
</section>`,

	// ── product_grid ──────────────────────────────────────────────────────────
	"product_grid": `<section class="section product-grid-section">
  <div class="container">
    {{- if strDefault .title "" }}
    <h2 style="margin-bottom: 2rem; text-align: center;">{{strDefault .title ""}}</h2>
    {{- end}}
    <div class="product-grid" style="
      display: grid;
      grid-template-columns: repeat({{strDefault .columns "4"}}, 1fr);
      gap: calc(var(--spacing-unit) * 6);
    ">
      <div class="product-grid-placeholder" style="
        grid-column: 1 / -1;
        padding: 3rem;
        background: var(--color-surface);
        border: 2px dashed var(--color-border);
        border-radius: var(--border-radius-lg);
        text-align: center;
        color: var(--color-text-muted);
      ">
        <p style="font-size: 1.1rem;">Products loading&hellip;</p>
        {{- if strDefault .collection_id "" }}
        <p style="margin-top: 0.5rem; font-size: 0.875rem;">Collection: {{strDefault .collection_id ""}}</p>
        {{- end}}
      </div>
    </div>
  </div>
</section>`,

	// ── featured_collection ───────────────────────────────────────────────────
	"featured_collection": `<section class="section featured-collection-section">
  <div class="container">
    {{- if strDefault .title "" }}
    <h2 style="margin-bottom: 2rem; text-align: center;">{{strDefault .title ""}}</h2>
    {{- end}}
    <div class="featured-collection-placeholder" style="
      padding: 3rem;
      background: var(--color-surface);
      border: 2px dashed var(--color-border);
      border-radius: var(--border-radius-lg);
      text-align: center;
      color: var(--color-text-muted);
    ">
      <p style="font-size: 1.1rem;">Featured Collection loading&hellip;</p>
      {{- if strDefault .collection_id "" }}
      <p style="margin-top: 0.5rem; font-size: 0.875rem;">Collection: {{strDefault .collection_id ""}}</p>
      {{- end}}
    </div>
  </div>
</section>`,

	// ── image_block ───────────────────────────────────────────────────────────
	"image_block": `{{- if strDefault .src "" }}
<section class="section image-block-section">
  <div class="container" style="text-align: center;">
    {{- if strDefault .link_url "" }}
    <a href="{{strDefault .link_url ""}}" style="display: inline-block;">
    {{- end}}
    <img
      src="{{strDefault .src ""}}"
      alt="{{strDefault .alt ""}}"
      {{- if strDefault .width "" }} width="{{strDefault .width ""}}"{{end}}
      {{- if strDefault .height "" }} height="{{strDefault .height ""}}"{{end}}
      style="border-radius: var(--border-radius-md); box-shadow: var(--shadow-md); max-width: 100%;"
    />
    {{- if strDefault .link_url "" }}
    </a>
    {{- end}}
    {{- if strDefault .caption "" }}
    <p style="margin-top: 1rem; color: var(--color-text-muted); font-size: 0.9rem; font-style: italic;">
      {{strDefault .caption ""}}
    </p>
    {{- end}}
  </div>
</section>
{{- end}}`,

	// ── video_block ───────────────────────────────────────────────────────────
	"video_block": `{{- if strDefault .url "" }}
<section class="section video-block-section">
  <div class="container">
    {{- if strDefault .title "" }}
    <h2 style="margin-bottom: 1.5rem; text-align: center;">{{strDefault .title ""}}</h2>
    {{- end}}
    <div class="video-wrapper" style="position: relative; padding-bottom: 56.25%; height: 0; overflow: hidden; border-radius: var(--border-radius-lg); box-shadow: var(--shadow-lg);">
      {{- if eq (strDefault .type "video") "youtube" }}
      <iframe
        src="{{strDefault .url ""}}"
        style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; border: 0;"
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture"
        allowfullscreen
        title="{{strDefault .title "Video"}}"
      ></iframe>
      {{- else}}
      <video
        src="{{strDefault .url ""}}"
        style="position: absolute; top: 0; left: 0; width: 100%; height: 100%; object-fit: cover;"
        controls
        {{- if boolVal .autoplay}} autoplay muted playsinline{{end}}
        title="{{strDefault .title "Video"}}"
      ></video>
      {{- end}}
    </div>
  </div>
</section>
{{- end}}`,

	// ── banner ────────────────────────────────────────────────────────────────
	"banner": `<div class="banner-section" style="
  background-color: {{strDefault .background_color "var(--color-primary)"}};
  color: {{strDefault .text_color "#ffffff"}};
  padding: 0.75rem var(--spacing-unit);
  text-align: center;
  position: relative;
">
  {{- if strDefault .link_url "" }}
  <a href="{{strDefault .link_url ""}}" style="color: inherit; font-weight: 600;">
    {{strDefault .text ""}}
  </a>
  {{- else}}
  <span style="font-weight: 600;">{{strDefault .text ""}}</span>
  {{- end}}
  {{- if boolVal .dismissible }}
  <button
    onclick="this.parentElement.style.display='none'"
    style="
      position: absolute;
      right: 1rem;
      top: 50%;
      transform: translateY(-50%);
      background: transparent;
      border: none;
      color: inherit;
      font-size: 1.25rem;
      cursor: pointer;
      padding: 0;
      line-height: 1;
    "
    aria-label="Dismiss"
  >&times;</button>
  {{- end}}
</div>`,

	// ── newsletter_signup ─────────────────────────────────────────────────────
	"newsletter_signup": `<section class="section newsletter-section" style="background: var(--color-surface);">
  <div class="container" style="max-width: 600px; text-align: center;">
    {{- if strDefault .heading "" }}
    <h2 style="margin-bottom: 1rem;">{{strDefault .heading ""}}</h2>
    {{- end}}
    {{- if strDefault .subheading "" }}
    <p style="margin-bottom: 2rem; color: var(--color-text-muted);">{{strDefault .subheading ""}}</p>
    {{- end}}
    <form class="newsletter-form" onsubmit="return false;" style="display: flex; gap: 0.5rem; flex-wrap: wrap; justify-content: center;">
      <input
        type="email"
        placeholder="{{strDefault .placeholder "Enter your email"}}"
        style="
          flex: 1;
          min-width: 240px;
          padding: 0.75rem 1rem;
          border: 1px solid var(--color-border);
          border-radius: var(--border-radius-md);
          font-size: var(--font-base-size);
          background: var(--color-background);
          color: var(--color-text);
          outline: none;
        "
      />
      <button type="submit" class="btn btn-primary">
        {{strDefault .button_text "Subscribe"}}
      </button>
    </form>
  </div>
</section>`,

	// ── testimonials ──────────────────────────────────────────────────────────
	"testimonials": `<section class="section testimonials-section">
  <div class="container">
    <div class="testimonials-grid" style="
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: calc(var(--spacing-unit) * 6);
    ">
      {{/* items is expected to be []interface{} of map[string]interface{} */}}
      <!-- Testimonials are loaded dynamically -->
    </div>
  </div>
</section>`,

	// ── cta ───────────────────────────────────────────────────────────────────
	"cta": `<section class="section cta-section" style="
  background-color: {{strDefault .background_color "var(--color-primary)"}};
  color: {{strDefault .text_color "#ffffff"}};
  text-align: center;
">
  <div class="container">
    {{- if strDefault .heading "" }}
    <h2 style="font-size: 2.5rem; margin-bottom: 1rem; {{if strDefault .text_color ""}}color: {{strDefault .text_color ""}};{{end}}">
      {{strDefault .heading ""}}
    </h2>
    {{- end}}
    {{- if strDefault .subheading "" }}
    <p style="font-size: 1.2rem; margin-bottom: 2rem; opacity: 0.9;">
      {{strDefault .subheading ""}}
    </p>
    {{- end}}
    {{- if strDefault .button_text "" }}
    <a
      href="{{strDefault .button_url "#"}}"
      class="btn"
      style="
        background: #ffffff;
        color: {{strDefault .background_color "var(--color-primary)"}};
        font-size: 1.1rem;
        padding: 1rem 2.5rem;
        border-radius: var(--border-radius-full);
        box-shadow: var(--shadow-md);
        font-weight: 700;
      "
    >
      {{strDefault .button_text "Get Started"}}
    </a>
    {{- end}}
  </div>
</section>`,
}

// renderTestimonialsWithItems handles the testimonials block specially because its
// settings contain a nested array that requires Go-side iteration.
func renderTestimonialsBlock(settings json.RawMessage) (string, error) {
	data := make(map[string]interface{})
	if len(settings) > 0 {
		if err := json.Unmarshal(settings, &data); err != nil {
			return "", fmt.Errorf("renderer: failed to parse testimonials settings: %w", err)
		}
	}

	var sb strings.Builder
	sb.WriteString(`<section class="section testimonials-section">
  <div class="container">
    <div class="testimonials-grid" style="
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: calc(var(--spacing-unit) * 6);
    ">
`)

	items, _ := data["items"].([]interface{})
	for _, raw := range items {
		item, ok := raw.(map[string]interface{})
		if !ok {
			continue
		}
		name := stringVal(item["name"])
		role := stringVal(item["role"])
		content := stringVal(item["content"])
		avatar := stringVal(item["avatar"])

		sb.WriteString(`      <div class="testimonial-card" style="
        background: var(--color-surface);
        border: 1px solid var(--color-border);
        border-radius: var(--border-radius-lg);
        padding: 2rem;
        box-shadow: var(--shadow-sm);
      ">
`)
		sb.WriteString(fmt.Sprintf(`        <p style="color: var(--color-text); margin-bottom: 1.5rem; font-style: italic; line-height: 1.7;">%s</p>
`, template.HTMLEscapeString(content)))
		sb.WriteString(`        <div style="display: flex; align-items: center; gap: 0.75rem;">
`)
		if avatar != "" {
			sb.WriteString(fmt.Sprintf(`          <img src="%s" alt="%s" style="width: 48px; height: 48px; border-radius: 50%%; object-fit: cover;" />
`,
				template.HTMLEscapeString(avatar),
				template.HTMLEscapeString(name),
			))
		}
		sb.WriteString(`          <div>
`)
		sb.WriteString(fmt.Sprintf(`            <p style="font-weight: 600; color: var(--color-text);">%s</p>
`, template.HTMLEscapeString(name)))
		if role != "" {
			sb.WriteString(fmt.Sprintf(`            <p style="font-size: 0.875rem; color: var(--color-text-muted);">%s</p>
`, template.HTMLEscapeString(role)))
		}
		sb.WriteString(`          </div>
        </div>
      </div>
`)
	}

	sb.WriteString(`    </div>
  </div>
</section>`)

	return sb.String(), nil
}

// stringVal safely converts an interface{} to a string, returning "" if nil or not a string.
func stringVal(v interface{}) string {
	if v == nil {
		return ""
	}
	s, _ := v.(string)
	return s
}
