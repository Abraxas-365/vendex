package product

import (
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
	ID          kernel.ProductID
	TenantID    kernel.TenantID
	Name        string
	Description string
	Price       kernel.Money
	SKU         string
	Images      []string
	CategoryID  kernel.CategoryID
	Tags        []string
	Status      Status
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
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
