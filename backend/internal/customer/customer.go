package customer

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Address represents a physical address.
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	IsDefault  bool   `json:"is_default"`
}

// Customer represents a commerce customer entity.
type Customer struct {
	ID        kernel.CustomerID `json:"id" db:"id"`
	TenantID  kernel.TenantID   `json:"tenant_id" db:"tenant_id"`
	Email     kernel.Email      `json:"email" db:"email"`
	Name      string            `json:"name" db:"name"`
	Phone     string            `json:"phone" db:"phone"`
	Addresses []Address         `json:"addresses"`
	CreatedAt time.Time         `json:"created_at" db:"created_at"`
	UpdatedAt time.Time         `json:"updated_at" db:"updated_at"`
}

// DefaultAddress returns the customer's default address, or nil if none.
func (c *Customer) DefaultAddress() *Address {
	for i := range c.Addresses {
		if c.Addresses[i].IsDefault {
			return &c.Addresses[i]
		}
	}
	return nil
}

// AddAddress adds an address. If isDefault, clears default from other addresses.
func (c *Customer) AddAddress(addr Address) {
	if addr.IsDefault {
		for i := range c.Addresses {
			c.Addresses[i].IsDefault = false
		}
	}
	c.Addresses = append(c.Addresses, addr)
	c.UpdatedAt = time.Now()
}

// SetDefaultAddress sets the address at the given index as default.
func (c *Customer) SetDefaultAddress(idx int) bool {
	if idx < 0 || idx >= len(c.Addresses) {
		return false
	}
	for i := range c.Addresses {
		c.Addresses[i].IsDefault = (i == idx)
	}
	c.UpdatedAt = time.Now()
	return true
}
