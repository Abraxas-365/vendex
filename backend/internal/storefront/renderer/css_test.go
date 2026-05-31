package renderer_test

import (
	"strings"
	"testing"

	"github.com/Abraxas-365/vendex/internal/storefront/renderer"
	"github.com/Abraxas-365/vendex/internal/theme"
)

// ── helpers ───────────────────────────────────────────────────────────────────

func assertContains(t *testing.T, haystack, needle string) {
	t.Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("expected CSS to contain %q\nGot:\n%s", needle, haystack)
	}
}

// ── TokensToCSS tests ─────────────────────────────────────────────────────────

func TestTokensToCSS_WrapsInRootBlock(t *testing.T) {
	css := renderer.TokensToCSS(theme.DefaultTokens())
	if !strings.HasPrefix(css, ":root {") {
		t.Errorf("expected CSS to start with \":root {\", got: %q", css[:min(50, len(css))])
	}
	if !strings.HasSuffix(strings.TrimSpace(css), "}") {
		t.Errorf("expected CSS to end with \"}\", got: %q", css[max(0, len(css)-50):])
	}
}

func TestTokensToCSS_DefaultTokens_ColorProperties(t *testing.T) {
	tokens := theme.DefaultTokens()
	css := renderer.TokensToCSS(tokens)

	checks := []string{
		"--color-primary: " + tokens.Colors.Primary + ";",
		"--color-secondary: " + tokens.Colors.Secondary + ";",
		"--color-background: " + tokens.Colors.Background + ";",
		"--color-surface: " + tokens.Colors.Surface + ";",
		"--color-text: " + tokens.Colors.Text + ";",
		"--color-text-muted: " + tokens.Colors.TextMuted + ";",
		"--color-border: " + tokens.Colors.Border + ";",
		"--color-error: " + tokens.Colors.Error + ";",
		"--color-success: " + tokens.Colors.Success + ";",
		"--color-warning: " + tokens.Colors.Warning + ";",
		"--color-info: " + tokens.Colors.Info + ";",
	}

	for _, want := range checks {
		assertContains(t, css, want)
	}
}

func TestTokensToCSS_DefaultTokens_TypographyProperties(t *testing.T) {
	tokens := theme.DefaultTokens()
	css := renderer.TokensToCSS(tokens)

	assertContains(t, css, "--font-heading: "+tokens.Typography.FontHeading+";")
	assertContains(t, css, "--font-body: "+tokens.Typography.FontBody+";")
	assertContains(t, css, "--font-base-size: "+tokens.Typography.BaseSize+";")
	// Scale ratio is formatted with %g — just check it contains the key.
	assertContains(t, css, "--font-scale-ratio:")
}

func TestTokensToCSS_DefaultTokens_SpacingProperties(t *testing.T) {
	tokens := theme.DefaultTokens()
	css := renderer.TokensToCSS(tokens)

	assertContains(t, css, "--spacing-unit: "+tokens.Spacing.Unit+";")
	assertContains(t, css, "--spacing-section-padding: "+tokens.Spacing.SectionPadding+";")
}

func TestTokensToCSS_DefaultTokens_BorderProperties(t *testing.T) {
	tokens := theme.DefaultTokens()
	css := renderer.TokensToCSS(tokens)

	assertContains(t, css, "--border-radius-sm: "+tokens.Borders.RadiusSm+";")
	assertContains(t, css, "--border-radius-md: "+tokens.Borders.RadiusMd+";")
	assertContains(t, css, "--border-radius-lg: "+tokens.Borders.RadiusLg+";")
	assertContains(t, css, "--border-radius-full: "+tokens.Borders.RadiusFull+";")
}

func TestTokensToCSS_DefaultTokens_ShadowProperties(t *testing.T) {
	tokens := theme.DefaultTokens()
	css := renderer.TokensToCSS(tokens)

	assertContains(t, css, "--shadow-sm: "+tokens.Shadows.Sm+";")
	assertContains(t, css, "--shadow-md: "+tokens.Shadows.Md+";")
	assertContains(t, css, "--shadow-lg: "+tokens.Shadows.Lg+";")
}

func TestTokensToCSS_CustomTokens(t *testing.T) {
	tokens := theme.ThemeTokens{
		Colors: theme.ColorTokens{
			Primary:    "#ff0000",
			Secondary:  "#00ff00",
			Background: "#000000",
			Surface:    "#111111",
			Text:       "#ffffff",
			TextMuted:  "#aaaaaa",
			Border:     "#333333",
			Error:      "#ff4444",
			Success:    "#44ff44",
			Warning:    "#ffaa00",
			Info:       "#4444ff",
		},
		Typography: theme.TypographyTokens{
			FontHeading: "Georgia, serif",
			FontBody:    "Helvetica, sans-serif",
			BaseSize:    "18px",
			ScaleRatio:  1.333,
		},
		Spacing: theme.SpacingTokens{
			Unit:           "8px",
			SectionPadding: "80px",
		},
		Borders: theme.BorderTokens{
			RadiusSm:   "2px",
			RadiusMd:   "6px",
			RadiusLg:   "12px",
			RadiusFull: "100%",
		},
		Shadows: theme.ShadowTokens{
			Sm: "0 1px 3px rgba(0,0,0,0.1)",
			Md: "0 4px 8px rgba(0,0,0,0.2)",
			Lg: "0 8px 16px rgba(0,0,0,0.3)",
		},
	}

	css := renderer.TokensToCSS(tokens)

	assertContains(t, css, "--color-primary: #ff0000;")
	assertContains(t, css, "--color-secondary: #00ff00;")
	assertContains(t, css, "--font-heading: Georgia, serif;")
	assertContains(t, css, "--font-body: Helvetica, sans-serif;")
	assertContains(t, css, "--font-base-size: 18px;")
	assertContains(t, css, "--spacing-unit: 8px;")
	assertContains(t, css, "--spacing-section-padding: 80px;")
	assertContains(t, css, "--border-radius-sm: 2px;")
	assertContains(t, css, "--border-radius-full: 100%;")
	assertContains(t, css, "--shadow-sm: 0 1px 3px rgba(0,0,0,0.1);")
	assertContains(t, css, "--shadow-lg: 0 8px 16px rgba(0,0,0,0.3);")
}

func TestTokensToCSS_EmptyTokens(t *testing.T) {
	// Empty tokens — all values are zero-value strings. Should not panic
	// and should still emit all CSS custom property names.
	css := renderer.TokensToCSS(theme.ThemeTokens{})

	requiredProps := []string{
		"--color-primary:",
		"--color-secondary:",
		"--color-background:",
		"--color-surface:",
		"--color-text:",
		"--color-text-muted:",
		"--color-border:",
		"--color-error:",
		"--color-success:",
		"--color-warning:",
		"--color-info:",
		"--font-heading:",
		"--font-body:",
		"--font-base-size:",
		"--font-scale-ratio:",
		"--spacing-unit:",
		"--spacing-section-padding:",
		"--border-radius-sm:",
		"--border-radius-md:",
		"--border-radius-lg:",
		"--border-radius-full:",
		"--shadow-sm:",
		"--shadow-md:",
		"--shadow-lg:",
	}

	for _, prop := range requiredProps {
		assertContains(t, css, prop)
	}
}

func TestTokensToCSS_TableDriven_Colors(t *testing.T) {
	tests := []struct {
		name     string
		color    string
		property string
		field    func(*theme.ColorTokens) *string
	}{
		{"primary", "#abc123", "--color-primary:", func(c *theme.ColorTokens) *string { return &c.Primary }},
		{"secondary", "#def456", "--color-secondary:", func(c *theme.ColorTokens) *string { return &c.Secondary }},
		{"error", "#ff0000", "--color-error:", func(c *theme.ColorTokens) *string { return &c.Error }},
		{"success", "#00ff00", "--color-success:", func(c *theme.ColorTokens) *string { return &c.Success }},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tokens := theme.DefaultTokens()
			*tc.field(&tokens.Colors) = tc.color

			css := renderer.TokensToCSS(tokens)
			want := tc.property + " " + tc.color + ";"
			assertContains(t, css, want)
		})
	}
}

// min/max helpers for Go < 1.21 (project uses go 1.25 but let's be safe).
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
