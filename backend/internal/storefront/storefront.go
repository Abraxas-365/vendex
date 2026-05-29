package storefront

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PageStatus represents the lifecycle state of a storefront page.
type PageStatus string

const (
	PageStatusDraft         PageStatus = "draft"
	PageStatusPendingReview PageStatus = "pending_review"
	PageStatusPublished     PageStatus = "published"
	PageStatusArchived      PageStatus = "archived"
)

// PageMeta holds SEO and Open Graph metadata for a page.
type PageMeta struct {
	Description string   `json:"description"`
	OGTitle     string   `json:"og_title"`
	OGImage     string   `json:"og_image"`
	Keywords    []string `json:"keywords"`
}

// Page is the core CMS entity representing a storefront page.
// Business rules:
//   - Pages created by an agent always start as pending_review.
//   - Admin-created pages start as draft.
//   - Publishing requires an explicit Publish() call — never auto-publish.
//   - Every edit produces a new PageVersion snapshot.
//   - Only published pages are served publicly.
type Page struct {
	ID          kernel.PageID  `json:"id"`
	TenantID    kernel.TenantID `json:"tenant_id"`
	Slug        string          `json:"slug"`
	Title       string          `json:"title"`
	HTML        string          `json:"html"`
	CSS         string          `json:"css"`
	Meta        PageMeta        `json:"meta"`
	Status      PageStatus      `json:"status"`
	Version     int             `json:"version"`
	CreatedBy   string          `json:"created_by"`
	PublishedAt *time.Time      `json:"published_at,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// IsPublished returns true when the page is live.
func (p *Page) IsPublished() bool {
	return p.Status == PageStatusPublished
}

// CanBePublished returns true when the page is in a state that allows publishing.
func (p *Page) CanBePublished() bool {
	return p.Status == PageStatusDraft || p.Status == PageStatusPendingReview
}

// CanBeEdited returns true when the page is not archived.
func (p *Page) CanBeEdited() bool {
	return p.Status != PageStatusArchived
}

// PageVersion is a full snapshot of a page's HTML/CSS at a given version number.
// The history is append-only — versions are never deleted.
type PageVersion struct {
	ID        kernel.PageVersionID `json:"id"`
	PageID    kernel.PageID        `json:"page_id"`
	TenantID  kernel.TenantID      `json:"tenant_id"`
	Version   int                  `json:"version"`
	HTML      string               `json:"html"`
	CSS       string               `json:"css"`
	EditedBy  string               `json:"edited_by"`
	Comment   string               `json:"comment"`
	CreatedAt time.Time            `json:"created_at"`
}

// TemplateTag represents a dynamic tag embedded in page content, e.g. {{products "featured" limit=8}}.
// The storefront renderer resolves these tags at serve-time by delegating to the appropriate domain.
type TemplateTag struct {
	// Type identifies which resolver handles this tag (e.g. "products", "promo_banner", "category").
	Type string `json:"type"`
	// Args holds the parsed arguments from the tag syntax.
	Args map[string]any `json:"args"`
}
