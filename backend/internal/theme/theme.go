package theme

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Theme represents a design token set for a tenant's storefront.
type Theme struct {
	ID        kernel.ThemeID  `json:"id"`
	TenantID  kernel.TenantID `json:"tenant_id"`
	Name      string          `json:"name"`
	IsActive  bool            `json:"is_active"`
	Tokens    ThemeTokens     `json:"tokens"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// ThemeTokens holds all design tokens for a theme.
type ThemeTokens struct {
	Colors     ColorTokens     `json:"colors"`
	Typography TypographyTokens `json:"typography"`
	Spacing    SpacingTokens   `json:"spacing"`
	Borders    BorderTokens    `json:"borders"`
	Shadows    ShadowTokens    `json:"shadows"`
}

// ColorTokens holds color design tokens.
type ColorTokens struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Background string `json:"background"`
	Surface    string `json:"surface"`
	Text       string `json:"text"`
	TextMuted  string `json:"text_muted"`
	Border     string `json:"border"`
	Error      string `json:"error"`
	Success    string `json:"success"`
	Warning    string `json:"warning"`
	Info       string `json:"info"`
}

// TypographyTokens holds typography design tokens.
type TypographyTokens struct {
	FontHeading string  `json:"font_heading"`
	FontBody    string  `json:"font_body"`
	BaseSize    string  `json:"base_size"`
	ScaleRatio  float64 `json:"scale_ratio"`
}

// SpacingTokens holds spacing design tokens.
type SpacingTokens struct {
	Unit           string `json:"unit"`
	SectionPadding string `json:"section_padding"`
}

// BorderTokens holds border design tokens.
type BorderTokens struct {
	RadiusSm   string `json:"radius_sm"`
	RadiusMd   string `json:"radius_md"`
	RadiusLg   string `json:"radius_lg"`
	RadiusFull string `json:"radius_full"`
}

// ShadowTokens holds shadow design tokens.
type ShadowTokens struct {
	Sm string `json:"sm"`
	Md string `json:"md"`
	Lg string `json:"lg"`
}

// DefaultTokens returns a sensible set of design token defaults.
func DefaultTokens() ThemeTokens {
	return ThemeTokens{
		Colors: ColorTokens{
			Primary:    "#6366f1", // indigo-500
			Secondary:  "#8b5cf6", // violet-500
			Background: "#ffffff",
			Surface:    "#f9fafb",
			Text:       "#111827",
			TextMuted:  "#6b7280",
			Border:     "#e5e7eb",
			Error:      "#ef4444",
			Success:    "#22c55e",
			Warning:    "#f59e0b",
			Info:       "#3b82f6",
		},
		Typography: TypographyTokens{
			FontHeading: "Inter, sans-serif",
			FontBody:    "Inter, sans-serif",
			BaseSize:    "16px",
			ScaleRatio:  1.25,
		},
		Spacing: SpacingTokens{
			Unit:           "4px",
			SectionPadding: "64px",
		},
		Borders: BorderTokens{
			RadiusSm:   "4px",
			RadiusMd:   "8px",
			RadiusLg:   "16px",
			RadiusFull: "9999px",
		},
		Shadows: ShadowTokens{
			Sm: "0 1px 2px 0 rgb(0 0 0 / 0.05)",
			Md: "0 4px 6px -1px rgb(0 0 0 / 0.1), 0 2px 4px -2px rgb(0 0 0 / 0.1)",
			Lg: "0 10px 15px -3px rgb(0 0 0 / 0.1), 0 4px 6px -4px rgb(0 0 0 / 0.1)",
		},
	}
}
