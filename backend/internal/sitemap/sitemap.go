package sitemap

import (
	"fmt"
	"strings"
	"time"
)

// URLEntry represents a single URL entry in the sitemap.
type URLEntry struct {
	Loc        string
	LastMod    string // RFC3339
	ChangeFreq string // "daily", "weekly", "monthly"
	Priority   string // "0.5", "0.8", "1.0"
}

// Sitemap holds a collection of URL entries to be serialised as XML.
type Sitemap struct {
	URLs []URLEntry
}

// ToXML renders the sitemap as an XML string conforming to the
// sitemaps.org/schemas/sitemap/0.9 namespace.
func (s *Sitemap) ToXML() string {
	var b strings.Builder
	b.WriteString(`<?xml version="1.0" encoding="UTF-8"?>`)
	b.WriteString("\n")
	b.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">`)
	b.WriteString("\n")

	for _, u := range s.URLs {
		b.WriteString("  <url>\n")
		b.WriteString(fmt.Sprintf("    <loc>%s</loc>\n", escapeXML(u.Loc)))
		if u.LastMod != "" {
			b.WriteString(fmt.Sprintf("    <lastmod>%s</lastmod>\n", u.LastMod))
		}
		if u.ChangeFreq != "" {
			b.WriteString(fmt.Sprintf("    <changefreq>%s</changefreq>\n", u.ChangeFreq))
		}
		if u.Priority != "" {
			b.WriteString(fmt.Sprintf("    <priority>%s</priority>\n", u.Priority))
		}
		b.WriteString("  </url>\n")
	}

	b.WriteString("</urlset>\n")
	return b.String()
}

// FormatTime formats a time.Time value as RFC3339 for use in lastmod fields.
func FormatTime(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// escapeXML escapes the five predefined XML entities in a string.
func escapeXML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&apos;")
	return s
}
