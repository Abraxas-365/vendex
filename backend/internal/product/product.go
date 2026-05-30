package product

import (
	"strings"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Status represents the lifecycle state of a product.
type Status string

const (
	StatusDraft    Status = "draft"
	StatusActive   Status = "active"
	StatusArchived Status = "archived"
)

// Product is the core entity for sellable items.
type Product struct {
	ID          kernel.ProductID  `json:"id" db:"id"`
	TenantID    kernel.TenantID   `json:"tenant_id" db:"tenant_id"`
	Name        string            `json:"name" db:"name"`
	Description string            `json:"description" db:"description"`
	Price       kernel.Money      `json:"price"`
	SKU         string            `json:"sku" db:"sku"`
	Images      []string          `json:"images"`
	CategoryID  kernel.CategoryID `json:"category_id" db:"category_id"`
	Tags        []string          `json:"tags"`
	Status      Status            `json:"status" db:"status"`
	Stock       int               `json:"stock" db:"stock"`
	HasVariants     bool              `json:"has_variants" db:"has_variants"`
	MetaTitle       string            `json:"meta_title,omitempty" db:"meta_title"`
	MetaDescription string            `json:"meta_description,omitempty" db:"meta_description"`
	Slug            string            `json:"slug" db:"slug"`
	CanonicalURL    string            `json:"canonical_url,omitempty" db:"canonical_url"`
	Options         []ProductOption   `json:"options,omitempty"`
	Variants        []ProductVariant  `json:"variants,omitempty"`
	CreatedAt       time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at" db:"updated_at"`
}

// GenerateSlug creates a URL-safe slug from a product name.
// It lowercases, replaces non-alphanumeric characters with hyphens,
// collapses consecutive hyphens, and trims leading/trailing hyphens.
func GenerateSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	var b strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteByte('-')
		}
	}
	result := b.String()
	for strings.Contains(result, "--") {
		result = strings.ReplaceAll(result, "--", "-")
	}
	return strings.Trim(result, "-")
}

// Activate transitions the product to active status.
func (p *Product) Activate() {
	p.Status = StatusActive
	p.UpdatedAt = time.Now()
}

// Archive transitions the product to archived status.
func (p *Product) Archive() {
	p.Status = StatusArchived
	p.UpdatedAt = time.Now()
}

// IsAvailable returns true if the product is active and in stock.
func (p *Product) IsAvailable() bool {
	return p.Status == StatusActive && p.Stock > 0
}

// DeductStock reduces stock by qty. Returns false if insufficient stock.
func (p *Product) DeductStock(qty int) bool {
	if p.Stock < qty {
		return false
	}
	p.Stock -= qty
	p.UpdatedAt = time.Now()
	return true
}

// AddStock increases stock by qty.
func (p *Product) AddStock(qty int) {
	p.Stock += qty
	p.UpdatedAt = time.Now()
}
