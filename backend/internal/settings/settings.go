package settings

import (
	"time"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// StoreSettings holds all configuration for a tenant's store.
type StoreSettings struct {
	TenantID       kernel.TenantID `json:"tenant_id" db:"tenant_id"`
	StoreName      string          `json:"store_name" db:"store_name"`
	StoreEmail     string          `json:"store_email" db:"store_email"`
	StorePhone     string          `json:"store_phone" db:"store_phone"`
	Currency       string          `json:"currency" db:"currency"`
	Timezone       string          `json:"timezone" db:"timezone"`
	Address        StoreAddress    `json:"address" db:"address"`
	LogoURL        string          `json:"logo_url" db:"logo_url"`
	FaviconURL     string          `json:"favicon_url" db:"favicon_url"`
	SocialLinks    SocialLinks     `json:"social_links" db:"social_links"`
	CheckoutConfig CheckoutConfig  `json:"checkout_config" db:"checkout_config"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

// StoreAddress holds the physical address of a store.
type StoreAddress struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	Country string `json:"country"`
	Zip     string `json:"zip"`
}

// SocialLinks holds URLs for a store's social media profiles.
type SocialLinks struct {
	Instagram string `json:"instagram"`
	Twitter   string `json:"twitter"`
	Facebook  string `json:"facebook"`
}

// CheckoutConfig controls checkout behaviour for a store.
type CheckoutConfig struct {
	GuestCheckout bool `json:"guest_checkout"`
	RequirePhone  bool `json:"require_phone"`
}

// DefaultSettings returns sensible defaults for a new tenant.
func DefaultSettings(tenantID kernel.TenantID) *StoreSettings {
	return &StoreSettings{
		TenantID:       tenantID,
		StoreName:      "My Store",
		Currency:       "USD",
		Timezone:       "UTC",
		CheckoutConfig: CheckoutConfig{GuestCheckout: true},
		UpdatedAt:      time.Now(),
	}
}
