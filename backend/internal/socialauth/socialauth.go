package socialauth

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Provider constants for supported OAuth providers.
const (
	ProviderGoogle   = "google"
	ProviderFacebook = "facebook"
)

// SupportedProviders lists all supported OAuth providers.
var SupportedProviders = map[string]bool{
	ProviderGoogle:   true,
	ProviderFacebook: true,
}

// SocialAccount represents a linked OAuth social account for a customer.
type SocialAccount struct {
	ID             kernel.SocialAccountID `json:"id"`
	TenantID       kernel.TenantID        `json:"tenant_id"`
	CustomerID     kernel.CustomerID      `json:"customer_id"`
	Provider       string                 `json:"provider"`
	ProviderUserID string                 `json:"provider_user_id"`
	Email          string                 `json:"email,omitempty"`
	Name           string                 `json:"name,omitempty"`
	AvatarURL      string                 `json:"avatar_url,omitempty"`
	AccessToken    string                 `json:"access_token,omitempty"`
	RefreshToken   string                 `json:"refresh_token,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at"`
}

// LinkInput holds the data needed to link a social account to a customer.
type LinkInput struct {
	CustomerID     kernel.CustomerID `json:"customer_id"`
	Provider       string            `json:"provider"`
	ProviderUserID string            `json:"provider_user_id"`
	Email          string            `json:"email"`
	Name           string            `json:"name"`
	AvatarURL      string            `json:"avatar_url"`
	AccessToken    string            `json:"access_token"`
	RefreshToken   string            `json:"refresh_token"`
}

// OAuthCallbackInput holds the data from an OAuth callback flow.
type OAuthCallbackInput struct {
	Provider    string `json:"provider"`
	Code        string `json:"code"`
	RedirectURI string `json:"redirect_uri"`
}

// ProviderConfig holds configuration for an OAuth provider.
type ProviderConfig struct {
	Provider string `json:"provider"`
	Name     string `json:"name"`
	Enabled  bool   `json:"enabled"`
}

// SupportedProviderConfigs returns the list of supported provider configs.
func SupportedProviderConfigs() []ProviderConfig {
	return []ProviderConfig{
		{Provider: ProviderGoogle, Name: "Google", Enabled: true},
		{Provider: ProviderFacebook, Name: "Facebook", Enabled: true},
	}
}
