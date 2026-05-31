package multistore

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Storefront represents a single storefront configuration for a tenant.
// A single tenant can run multiple storefronts (e.g. wholesale, retail, B2B, regional).
type Storefront struct {
	ID              kernel.StorefrontEntryID `json:"id"`
	TenantID        kernel.TenantID          `json:"tenant_id"`
	Name            string                   `json:"name"`
	Slug            string                   `json:"slug"`
	Domain          *string                  `json:"domain,omitempty"`
	Description     string                   `json:"description,omitempty"`
	ThemeID         string                   `json:"theme_id,omitempty"`
	LogoURL         string                   `json:"logo_url,omitempty"`
	DefaultLocale   string                   `json:"default_locale"`
	DefaultCurrency string                   `json:"default_currency"`
	IsActive        bool                     `json:"is_active"`
	IsDefault       bool                     `json:"is_default"`
	Settings        map[string]interface{}   `json:"settings"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
}

// StorefrontCatalog links a storefront to a catalog (product grouping).
type StorefrontCatalog struct {
	ID           kernel.StorefrontCatalogID `json:"id"`
	TenantID     kernel.TenantID            `json:"tenant_id"`
	StorefrontID kernel.StorefrontEntryID   `json:"storefront_id"`
	CatalogID    string                     `json:"catalog_id"`
	SortOrder    int                        `json:"sort_order"`
	CreatedAt    time.Time                  `json:"created_at"`
}

// CreateInput holds the data required to create a new storefront.
type CreateInput struct {
	Name            string
	Slug            string
	Domain          *string
	Description     string
	ThemeID         string
	LogoURL         string
	DefaultLocale   string
	DefaultCurrency string
	IsActive        bool
	Settings        map[string]interface{}
}

// UpdateInput holds the fields that can be updated on a storefront.
// Nil pointer fields are treated as "no change".
type UpdateInput struct {
	Name            *string
	Domain          *string
	Description     *string
	ThemeID         *string
	LogoURL         *string
	DefaultLocale   *string
	DefaultCurrency *string
	IsActive        *bool
	Settings        map[string]interface{}
}
