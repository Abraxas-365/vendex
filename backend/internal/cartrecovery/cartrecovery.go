package cartrecovery

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Recovery status constants.
const (
	StatusPending   = "pending"
	StatusSent      = "sent"
	StatusClicked   = "clicked"
	StatusConverted = "converted"
)

// RecoveryStatus is the lifecycle state of a recovery email.
type RecoveryStatus string

// RecoveryEmail tracks a single recovery email sent (or to be sent) for an abandoned cart.
type RecoveryEmail struct {
	ID           kernel.RecoveryID  `json:"id"            db:"id"`
	TenantID     kernel.TenantID    `json:"tenant_id"     db:"tenant_id"`
	CartID       kernel.CartID      `json:"cart_id"       db:"cart_id"`
	CustomerID   kernel.CustomerID  `json:"customer_id"   db:"customer_id"`
	Email        string             `json:"email"         db:"email"`
	Step         int                `json:"step"          db:"step"`
	Status       RecoveryStatus     `json:"status"        db:"status"`
	DiscountCode *string            `json:"discount_code" db:"discount_code"`
	SentAt       *time.Time         `json:"sent_at"       db:"sent_at"`
	ClickedAt    *time.Time         `json:"clicked_at"    db:"clicked_at"`
	ConvertedAt  *time.Time         `json:"converted_at"  db:"converted_at"`
	CreatedAt    time.Time          `json:"created_at"    db:"created_at"`
}

// RecoveryStats holds aggregate counts for a tenant's recovery emails.
type RecoveryStats struct {
	Total          int     `json:"total"`
	Pending        int     `json:"pending"`
	Sent           int     `json:"sent"`
	Clicked        int     `json:"clicked"`
	Converted      int     `json:"converted"`
	ConversionRate float64 `json:"conversion_rate"`
}
