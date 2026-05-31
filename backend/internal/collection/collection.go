package collection

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// CollectionType distinguishes manual (curated) from automatic (rule-based) collections.
type CollectionType string

const (
	// CollectionManual is a hand-curated list of products.
	CollectionManual CollectionType = "manual"
	// CollectionAuto is a rule-based collection whose members are computed dynamically.
	CollectionAuto CollectionType = "auto"
)

// CollectionRule is a single filter rule used by automatic collections.
type CollectionRule struct {
	// Field is the product attribute to filter on: "price", "tag", "category", "vendor", "type".
	Field string `json:"field"`
	// Operator is the comparison: "eq", "neq", "gt", "lt", "contains", "starts_with".
	Operator string `json:"operator"`
	// Value is the scalar value to compare against.
	Value string `json:"value"`
}

// Collection represents a curated or rule-based group of products.
type Collection struct {
	ID              kernel.CollectionID `json:"id"               db:"id"`
	TenantID        kernel.TenantID     `json:"tenant_id"        db:"tenant_id"`
	Name            string              `json:"name"             db:"name"`
	Slug            string              `json:"slug"             db:"slug"`
	Description     string              `json:"description"      db:"description"`
	ImageURL        string              `json:"image_url"        db:"image_url"`
	Type            CollectionType      `json:"type"             db:"type"`
	Rules           []CollectionRule    `json:"rules"            db:"rules"`
	IsActive        bool                `json:"is_active"        db:"is_active"`
	SortOrder       int                 `json:"sort_order"       db:"sort_order"`
	MetaTitle       string              `json:"meta_title"       db:"meta_title"`
	MetaDescription string              `json:"meta_description" db:"meta_description"`
	PublishedAt     *time.Time          `json:"published_at"     db:"published_at"`
	ProductCount    int                 `json:"product_count"    db:"-"`
	CreatedAt       time.Time           `json:"created_at"       db:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"       db:"updated_at"`
}

// IsManual returns true when this collection's products are managed manually.
func (c *Collection) IsManual() bool { return c.Type == CollectionManual }

// IsAuto returns true when this collection's members are computed from rules.
func (c *Collection) IsAuto() bool { return c.Type == CollectionAuto }

// CollectionProduct records a product's membership in a collection.
type CollectionProduct struct {
	ID           kernel.CollectionProductID `json:"id"            db:"id"`
	TenantID     kernel.TenantID            `json:"tenant_id"     db:"tenant_id"`
	CollectionID kernel.CollectionID        `json:"collection_id" db:"collection_id"`
	ProductID    string                     `json:"product_id"    db:"product_id"`
	SortOrder    int                        `json:"sort_order"    db:"sort_order"`
	AddedAt      time.Time                  `json:"added_at"      db:"added_at"`
}

// --- Request / response DTOs ---

// CreateInput carries the fields required to create a collection.
type CreateInput struct {
	Name            string
	Slug            string
	Description     string
	ImageURL        string
	Type            CollectionType
	Rules           []CollectionRule
	IsActive        bool
	SortOrder       int
	MetaTitle       string
	MetaDescription string
	PublishedAt     *time.Time
}

// UpdateInput holds optional fields for a collection update.
// Only non-nil pointer fields are applied.
type UpdateInput struct {
	Name            *string
	Slug            *string
	Description     *string
	ImageURL        *string
	Type            *CollectionType
	Rules           []CollectionRule
	IsActive        *bool
	SortOrder       *int
	MetaTitle       *string
	MetaDescription *string
	PublishedAt     *time.Time
	ClearPublishedAt bool
}
