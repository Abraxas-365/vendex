package renderer

import (
	"fmt"
	"strings"

	"github.com/Abraxas-365/vendex/internal/theme"
)

// TokensToCSS converts a ThemeTokens struct into a CSS :root block containing
// all design tokens as CSS custom properties.
func TokensToCSS(tokens theme.ThemeTokens) string {
	var sb strings.Builder

	sb.WriteString(":root {\n")

	// Color tokens
	sb.WriteString(fmt.Sprintf("  --color-primary: %s;\n", tokens.Colors.Primary))
	sb.WriteString(fmt.Sprintf("  --color-secondary: %s;\n", tokens.Colors.Secondary))
	sb.WriteString(fmt.Sprintf("  --color-background: %s;\n", tokens.Colors.Background))
	sb.WriteString(fmt.Sprintf("  --color-surface: %s;\n", tokens.Colors.Surface))
	sb.WriteString(fmt.Sprintf("  --color-text: %s;\n", tokens.Colors.Text))
	sb.WriteString(fmt.Sprintf("  --color-text-muted: %s;\n", tokens.Colors.TextMuted))
	sb.WriteString(fmt.Sprintf("  --color-border: %s;\n", tokens.Colors.Border))
	sb.WriteString(fmt.Sprintf("  --color-error: %s;\n", tokens.Colors.Error))
	sb.WriteString(fmt.Sprintf("  --color-success: %s;\n", tokens.Colors.Success))
	sb.WriteString(fmt.Sprintf("  --color-warning: %s;\n", tokens.Colors.Warning))
	sb.WriteString(fmt.Sprintf("  --color-info: %s;\n", tokens.Colors.Info))

	// Typography tokens
	sb.WriteString(fmt.Sprintf("  --font-heading: %s;\n", tokens.Typography.FontHeading))
	sb.WriteString(fmt.Sprintf("  --font-body: %s;\n", tokens.Typography.FontBody))
	sb.WriteString(fmt.Sprintf("  --font-base-size: %s;\n", tokens.Typography.BaseSize))
	sb.WriteString(fmt.Sprintf("  --font-scale-ratio: %g;\n", tokens.Typography.ScaleRatio))

	// Spacing tokens
	sb.WriteString(fmt.Sprintf("  --spacing-unit: %s;\n", tokens.Spacing.Unit))
	sb.WriteString(fmt.Sprintf("  --spacing-section-padding: %s;\n", tokens.Spacing.SectionPadding))

	// Border tokens
	sb.WriteString(fmt.Sprintf("  --border-radius-sm: %s;\n", tokens.Borders.RadiusSm))
	sb.WriteString(fmt.Sprintf("  --border-radius-md: %s;\n", tokens.Borders.RadiusMd))
	sb.WriteString(fmt.Sprintf("  --border-radius-lg: %s;\n", tokens.Borders.RadiusLg))
	sb.WriteString(fmt.Sprintf("  --border-radius-full: %s;\n", tokens.Borders.RadiusFull))

	// Shadow tokens
	sb.WriteString(fmt.Sprintf("  --shadow-sm: %s;\n", tokens.Shadows.Sm))
	sb.WriteString(fmt.Sprintf("  --shadow-md: %s;\n", tokens.Shadows.Md))
	sb.WriteString(fmt.Sprintf("  --shadow-lg: %s;\n", tokens.Shadows.Lg))

	sb.WriteString("}\n")

	return sb.String()
}

// baseCSS returns the base CSS reset and utility styles that use the CSS custom properties.
const baseCSS = `
* { margin: 0; padding: 0; box-sizing: border-box; }

body {
  font-family: var(--font-body);
  font-size: var(--font-base-size);
  color: var(--color-text);
  background: var(--color-background);
  line-height: 1.6;
}

h1, h2, h3, h4, h5, h6 {
  font-family: var(--font-heading);
  line-height: 1.2;
  color: var(--color-text);
}

a {
  color: var(--color-primary);
  text-decoration: none;
}

a:hover {
  text-decoration: underline;
}

.section {
  padding: var(--spacing-section-padding) var(--spacing-unit);
  width: 100%;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 calc(var(--spacing-unit) * 4);
}

img {
  max-width: 100%;
  height: auto;
  display: block;
}

button, .btn {
  cursor: pointer;
  display: inline-block;
  padding: 0.75rem 1.5rem;
  border: none;
  border-radius: var(--border-radius-md);
  font-family: var(--font-body);
  font-size: var(--font-base-size);
  font-weight: 600;
  transition: opacity 0.2s ease, transform 0.1s ease;
}

button:hover, .btn:hover {
  opacity: 0.9;
  transform: translateY(-1px);
}

.btn-primary {
  background: var(--color-primary);
  color: #ffffff;
  box-shadow: var(--shadow-sm);
}

.btn-secondary {
  background: var(--color-secondary);
  color: #ffffff;
  box-shadow: var(--shadow-sm);
}
`
