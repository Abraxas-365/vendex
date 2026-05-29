package catalog

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Category represents a product category in a tenant's catalog.
type Category struct {
	ID          kernel.CategoryID  `json:"id" db:"id"`
	TenantID    kernel.TenantID    `json:"tenant_id" db:"tenant_id"`
	Name        string             `json:"name" db:"name"`
	Slug        string             `json:"slug" db:"slug"`
	ParentID    *kernel.CategoryID `json:"parent_id" db:"parent_id"`
	Description string             `json:"description" db:"description"`
	CreatedAt   time.Time          `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" db:"updated_at"`
}

// IsRoot returns true if this category has no parent.
func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

// Collection represents a curated or automatic grouping of products.
type Collection struct {
	ID          kernel.CollectionID `json:"id" db:"id"`
	TenantID    kernel.TenantID     `json:"tenant_id" db:"tenant_id"`
	Name        string              `json:"name" db:"name"`
	Slug        string              `json:"slug" db:"slug"`
	Description string              `json:"description" db:"description"`
	ProductIDs  []kernel.ProductID  `json:"product_ids" db:"product_ids"`
	IsAutomatic bool                `json:"is_automatic" db:"is_automatic"`
	Rules       map[string]any      `json:"rules" db:"rules"`
	CreatedAt   time.Time           `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at" db:"updated_at"`
}

// ContainsProduct checks if a product is in this collection.
func (c *Collection) ContainsProduct(id kernel.ProductID) bool {
	for _, pid := range c.ProductIDs {
		if pid == id {
			return true
		}
	}
	return false
}

// AddProduct adds a product to the collection if not already present.
func (c *Collection) AddProduct(id kernel.ProductID) bool {
	if c.ContainsProduct(id) {
		return false
	}
	c.ProductIDs = append(c.ProductIDs, id)
	c.UpdatedAt = time.Now()
	return true
}

// RemoveProduct removes a product from the collection.
func (c *Collection) RemoveProduct(id kernel.ProductID) bool {
	for i, pid := range c.ProductIDs {
		if pid == id {
			c.ProductIDs = append(c.ProductIDs[:i], c.ProductIDs[i+1:]...)
			c.UpdatedAt = time.Now()
			return true
		}
	}
	return false
}
