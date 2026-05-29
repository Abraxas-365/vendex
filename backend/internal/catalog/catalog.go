package catalog

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Category represents a product category in a tenant's catalog.
type Category struct {
	ID          kernel.CategoryID
	TenantID    kernel.TenantID
	Name        string
	Slug        string
	ParentID    *kernel.CategoryID
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// IsRoot returns true if this category has no parent.
func (c *Category) IsRoot() bool {
	return c.ParentID == nil
}

// Collection represents a curated or automatic grouping of products.
type Collection struct {
	ID          kernel.CollectionID
	TenantID    kernel.TenantID
	Name        string
	Slug        string
	Description string
	ProductIDs  []kernel.ProductID
	IsAutomatic bool
	Rules       map[string]any
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
