package renderer_test

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/storefront"
	"github.com/Abraxas-365/vendex/internal/storefront/renderer"
	"github.com/Abraxas-365/vendex/internal/theme"
)

// ── mock ThemeGetter ──────────────────────────────────────────────────────────

type mockThemeGetter struct {
	theme *theme.Theme
	err   error
}

func (m *mockThemeGetter) GetActiveTheme(_ context.Context, _ kernel.TenantID) (*theme.Theme, error) {
	return m.theme, m.err
}

// themeWithTokens constructs a mock ThemeGetter that returns the given tokens.
func themeWithTokens(tokens theme.ThemeTokens) renderer.ThemeGetter {
	return &mockThemeGetter{
		theme: &theme.Theme{
			ID:       kernel.ThemeID("theme-1"),
			Name:     "Test Theme",
			IsActive: true,
			Tokens:   tokens,
		},
	}
}

// errorThemeGetter returns an error — renderer should fall back to defaults.
func errorThemeGetter() renderer.ThemeGetter {
	return &mockThemeGetter{err: errors.New("theme not found")}
}

// ── page builders ─────────────────────────────────────────────────────────────

func simplePage(tenantID kernel.TenantID, contentType storefront.ContentType) *storefront.Page {
	return &storefront.Page{
		ID:          kernel.PageID("page-1"),
		TenantID:    tenantID,
		Slug:        "home",
		Title:       "Home Page",
		ContentType: contentType,
		HTML:        "<p>Hello, world!</p>",
		Status:      storefront.PageStatusPublished,
	}
}

func pageWithSections(tenantID kernel.TenantID, sections []storefront.Section) *storefront.Page {
	p := simplePage(tenantID, storefront.ContentTypeBlocks)
	p.Sections = sections
	return p
}

func richTextSection(id, content string) storefront.Section {
	settings, _ := json.Marshal(map[string]string{"content": content})
	return storefront.Section{
		ID:       id,
		Type:     "rich_text",
		Settings: settings,
	}
}

func heroSection(id, heading string) storefront.Section {
	settings, _ := json.Marshal(map[string]string{"heading": heading, "button_text": "Shop Now"})
	return storefront.Section{
		ID:       id,
		Type:     "hero",
		Settings: settings,
	}
}

// ── tests ─────────────────────────────────────────────────────────────────────

func TestRenderPage_NilPage_ReturnsError(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))
	_, err := r.RenderPage(context.Background(), nil)
	if err == nil {
		t.Fatal("expected an error for nil page, got nil")
	}
}

func TestRenderPage_HTMLContentType_IncludesRawHTML(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))
	page := simplePage("tenant-1", storefront.ContentTypeHTML)
	page.HTML = "<h1>Custom HTML</h1>"

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}
	if !strings.Contains(html, "<h1>Custom HTML</h1>") {
		t.Errorf("expected raw HTML to appear in output")
	}
}

func TestRenderPage_ContainsHTMLDoctype(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))
	page := simplePage("tenant-1", storefront.ContentTypeHTML)

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}
	if !strings.HasPrefix(strings.TrimSpace(html), "<!DOCTYPE html>") {
		t.Error("expected output to start with <!DOCTYPE html>")
	}
}

func TestRenderPage_ContainsTitle(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))
	page := simplePage("tenant-1", storefront.ContentTypeHTML)
	page.Title = "My Storefront"

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}
	if !strings.Contains(html, "<title>My Storefront</title>") {
		t.Errorf("expected <title>My Storefront</title> in output")
	}
}

func TestRenderPage_ContainsCSSVariables(t *testing.T) {
	tokens := theme.DefaultTokens()
	tokens.Colors.Primary = "#abcdef"

	r := renderer.New(themeWithTokens(tokens))
	page := simplePage("tenant-1", storefront.ContentTypeHTML)

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}

	// The :root block with custom properties must appear inside a <style> tag.
	if !strings.Contains(html, ":root {") {
		t.Error("expected CSS :root block in output")
	}
	if !strings.Contains(html, "--color-primary: #abcdef;") {
		t.Error("expected --color-primary custom property with custom value")
	}
}

func TestRenderPage_FallsBackToDefaultTokens_WhenThemeErrors(t *testing.T) {
	r := renderer.New(errorThemeGetter())
	page := simplePage("tenant-1", storefront.ContentTypeHTML)

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage should succeed even when theme fetch fails: %v", err)
	}

	// Default primary is #6366f1.
	if !strings.Contains(html, "--color-primary: #6366f1;") {
		t.Error("expected default --color-primary in output when theme fetch fails")
	}
}

func TestRenderPage_RichTextSection(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))

	content := "<p>Welcome to our store!</p>"
	page := pageWithSections("tenant-1", []storefront.Section{
		richTextSection("s1", content),
	})

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}

	if !strings.Contains(html, content) {
		t.Errorf("expected rich_text content %q to appear in output", content)
	}
	if !strings.Contains(html, "rich-text-section") {
		t.Error("expected rich-text-section CSS class in output")
	}
}

func TestRenderPage_HeroSection(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))

	heading := "Big Sale Today!"
	page := pageWithSections("tenant-1", []storefront.Section{
		heroSection("s1", heading),
	})

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}

	if !strings.Contains(html, heading) {
		t.Errorf("expected hero heading %q in output", heading)
	}
	if !strings.Contains(html, "hero-section") {
		t.Error("expected hero-section CSS class in output")
	}
	if !strings.Contains(html, "Shop Now") {
		t.Error("expected button_text 'Shop Now' in output")
	}
}

func TestRenderPage_MultipleSections(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))

	page := pageWithSections("tenant-1", []storefront.Section{
		richTextSection("s1", "<p>Section One</p>"),
		richTextSection("s2", "<p>Section Two</p>"),
		heroSection("s3", "Buy Everything"),
	})

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}

	if !strings.Contains(html, "<p>Section One</p>") {
		t.Error("expected Section One content")
	}
	if !strings.Contains(html, "<p>Section Two</p>") {
		t.Error("expected Section Two content")
	}
	if !strings.Contains(html, "Buy Everything") {
		t.Error("expected hero heading")
	}
}

func TestRenderPage_UnknownBlockType_GracefullyDegrades(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))

	settings, _ := json.Marshal(map[string]string{"foo": "bar"})
	page := pageWithSections("tenant-1", []storefront.Section{
		{ID: "s1", Type: "totally_unknown_block_xyz", Settings: settings},
	})

	// Unknown block type should not error — it renders a comment placeholder.
	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage should not error on unknown block type: %v", err)
	}
	if !strings.Contains(html, "unknown block type") {
		t.Errorf("expected placeholder comment for unknown block type, got output without it")
	}
}

func TestRenderPage_ContainsHeaderAndFooter(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))
	page := simplePage("tenant-1", storefront.ContentTypeHTML)
	page.Title = "My Store"

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}

	if !strings.Contains(html, "site-header") {
		t.Error("expected site-header in output")
	}
	if !strings.Contains(html, "site-footer") {
		t.Error("expected site-footer in output")
	}
	if !strings.Contains(html, "<main>") {
		t.Error("expected <main> in output")
	}
}

func TestRenderPage_BlocksWithNestedBlocks(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))

	richSettings, _ := json.Marshal(map[string]string{"content": "<b>Nested block content</b>"})
	section := storefront.Section{
		ID:   "s1",
		Type: "container",
		Blocks: []storefront.Block{
			{ID: "b1", Type: "rich_text", Settings: richSettings},
		},
	}

	page := pageWithSections("tenant-1", []storefront.Section{section})

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage with nested blocks: %v", err)
	}

	if !strings.Contains(html, "<b>Nested block content</b>") {
		t.Error("expected nested block content in output")
	}
	// section-wrapper div should be rendered with data attributes
	if !strings.Contains(html, "section-wrapper") {
		t.Error("expected section-wrapper div for section with blocks")
	}
}

func TestRenderPage_PageCSS_IncludedInOutput(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))
	page := simplePage("tenant-1", storefront.ContentTypeHTML)
	page.CSS = ".custom { color: hotpink; }"

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}

	if !strings.Contains(html, ".custom { color: hotpink; }") {
		t.Error("expected per-page CSS to appear in output")
	}
}

func TestRenderPage_NilConfig_DoesNotPanic(t *testing.T) {
	// NewWithConfig with all-nil optional deps must not panic.
	r := renderer.NewWithConfig(themeWithTokens(theme.DefaultTokens()), renderer.Config{})
	page := simplePage("tenant-1", storefront.ContentTypeHTML)

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("NewWithConfig with nil deps should succeed: %v", err)
	}
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("expected valid HTML output")
	}
}

func TestRenderPage_StoreName_DefaultsToPageTitle(t *testing.T) {
	r := renderer.New(themeWithTokens(theme.DefaultTokens()))
	page := simplePage("tenant-1", storefront.ContentTypeHTML)
	page.Title = "Fancy Store"

	html, err := r.RenderPage(context.Background(), page)
	if err != nil {
		t.Fatalf("RenderPage: %v", err)
	}

	// When no settings getter is configured, StoreName falls back to page.Title.
	if !strings.Contains(html, "Fancy Store") {
		t.Error("expected store name (page title) in header")
	}
}
