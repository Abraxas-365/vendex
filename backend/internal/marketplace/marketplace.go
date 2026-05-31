package marketplace

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// VendorStatus represents the approval state of a vendor.
type VendorStatus string

const (
	VendorStatusPending  VendorStatus = "pending"
	VendorStatusApproved VendorStatus = "approved"
	VendorStatusRejected VendorStatus = "rejected"
	VendorStatusSuspended VendorStatus = "suspended"
)

// Vendor represents a marketplace seller/partner.
type Vendor struct {
	ID          kernel.VendorID  `json:"id" db:"id"`
	TenantID    kernel.TenantID  `json:"tenant_id" db:"tenant_id"`
	Name        string           `json:"name" db:"name"`
	Description string           `json:"description" db:"description"`
	Email       string           `json:"email" db:"email"`
	Phone       string           `json:"phone" db:"phone"`
	Status      VendorStatus     `json:"status" db:"status"`
	Commission  float64          `json:"commission" db:"commission"` // percentage, e.g. 10.0 = 10%
	CreatedAt   time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at" db:"updated_at"`
}

// VendorProduct links a product to a vendor with vendor-specific pricing.
type VendorProduct struct {
	ID        string          `json:"id" db:"id"`
	TenantID  kernel.TenantID `json:"tenant_id" db:"tenant_id"`
	VendorID  kernel.VendorID `json:"vendor_id" db:"vendor_id"`
	ProductID kernel.ProductID `json:"product_id" db:"product_id"`
	Price     kernel.Money    `json:"price" db:"-"`
	Stock     int             `json:"stock" db:"stock"`
	CreatedAt time.Time       `json:"created_at" db:"created_at"`
}

// VendorOrder tracks the vendor's share of an order.
type VendorOrder struct {
	ID          string          `json:"id" db:"id"`
	TenantID    kernel.TenantID `json:"tenant_id" db:"tenant_id"`
	VendorID    kernel.VendorID `json:"vendor_id" db:"vendor_id"`
	OrderID     kernel.OrderID  `json:"order_id" db:"order_id"`
	Amount      kernel.Money    `json:"amount" db:"-"`
	Commission  kernel.Money    `json:"commission_amount" db:"-"`
	Status      string          `json:"status" db:"status"`
	CreatedAt   time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at" db:"updated_at"`
}

// CreateVendorRequest holds input for creating a new vendor.
type CreateVendorRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Commission  float64 `json:"commission"`
}

// UpdateVendorRequest holds input for updating a vendor.
type UpdateVendorRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Email       *string  `json:"email"`
	Phone       *string  `json:"phone"`
	Status      *string  `json:"status"`
	Commission  *float64 `json:"commission"`
}
