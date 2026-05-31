package cart

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Cart is the core entity for a shopping cart.
type Cart struct {
	ID         kernel.CartID     `json:"id" db:"id"`
	TenantID   kernel.TenantID   `json:"tenant_id" db:"tenant_id"`
	CustomerID kernel.CustomerID `json:"customer_id,omitempty" db:"customer_id"`
	SessionID  string            `json:"session_id,omitempty" db:"session_id"`
	Items      []CartItem        `json:"items"`
	Currency   string            `json:"currency" db:"currency"`
	CreatedAt  time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at" db:"updated_at"`
	ExpiresAt  time.Time         `json:"expires_at" db:"expires_at"`
}

// CartItem represents a line item in the cart.
type CartItem struct {
	ID        kernel.CartItemID `json:"id" db:"id"`
	CartID    kernel.CartID     `json:"cart_id" db:"cart_id"`
	TenantID  kernel.TenantID   `json:"tenant_id" db:"tenant_id"`
	ProductID kernel.ProductID  `json:"product_id" db:"product_id"`
	VariantID string            `json:"variant_id,omitempty" db:"variant_id"` // nullable, for future variant support
	Quantity  int               `json:"quantity" db:"quantity"`
	UnitPrice kernel.Money      `json:"unit_price"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
}

// --- Domain errors ---

var (
	ErrNotFound     = errx.New("cart not found", errx.TypeNotFound)
	ErrItemNotFound = errx.New("cart item not found", errx.TypeNotFound)
	ErrEmptyCart    = errx.New("cart is empty", errx.TypeBusiness)
	ErrInvalidQty   = errx.New("quantity must be greater than 0", errx.TypeValidation)
)

// --- Domain methods ---

// Subtotal returns the sum of all item quantities * unit prices.
func (c *Cart) Subtotal() kernel.Money {
	total := kernel.Money{Amount: 0, Currency: c.Currency}
	for _, item := range c.Items {
		total = total.Add(item.UnitPrice.Multiply(item.Quantity))
	}
	return total
}

// ItemCount returns the total number of units across all items.
func (c *Cart) ItemCount() int {
	count := 0
	for _, item := range c.Items {
		count += item.Quantity
	}
	return count
}

// AddItem adds an item to the cart, merging quantities if the same product_id already exists.
func (c *Cart) AddItem(item CartItem) {
	for i, existing := range c.Items {
		if existing.ProductID == item.ProductID && existing.VariantID == item.VariantID {
			c.Items[i].Quantity += item.Quantity
			c.Items[i].UpdatedAt = time.Now()
			return
		}
	}
	c.Items = append(c.Items, item)
}

// RemoveItem removes the item with the given ID from the cart.
func (c *Cart) RemoveItem(itemID kernel.CartItemID) {
	filtered := c.Items[:0]
	for _, item := range c.Items {
		if item.ID != itemID {
			filtered = append(filtered, item)
		}
	}
	c.Items = filtered
}

// UpdateItemQuantity sets the quantity of a specific item. Qty must be > 0.
func (c *Cart) UpdateItemQuantity(itemID kernel.CartItemID, qty int) error {
	if qty <= 0 {
		return ErrInvalidQty
	}
	for i, item := range c.Items {
		if item.ID == itemID {
			c.Items[i].Quantity = qty
			c.Items[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrItemNotFound
}

// Clear removes all items from the cart.
func (c *Cart) Clear() {
	c.Items = []CartItem{}
}
