package socialauthsrv

import (
	"context"
	"errors"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/socialauth"
	"github.com/google/uuid"
)

// Service implements business logic for the social auth domain.
type Service struct {
	repo socialauth.Repository
}

// New creates a new social auth service.
func New(repo socialauth.Repository) *Service {
	return &Service{repo: repo}
}

// LinkAccount links a social provider account to a customer.
// Returns ErrAlreadyLinked if the provider account is already linked.
// Returns ErrProviderNotSupported if the provider is unknown.
func (s *Service) LinkAccount(ctx context.Context, tenantID kernel.TenantID, input socialauth.LinkInput) (socialauth.SocialAccount, error) {
	if !socialauth.SupportedProviders[input.Provider] {
		return socialauth.SocialAccount{}, socialauth.ErrProviderNotSupported()
	}

	// Check if this provider account is already linked (to any customer in this tenant).
	existing, err := s.repo.GetByProvider(ctx, tenantID, input.Provider, input.ProviderUserID)
	if err == nil {
		// Already exists — return the existing record wrapped as conflict.
		_ = existing
		return socialauth.SocialAccount{}, socialauth.ErrAlreadyLinked()
	}
	// Only proceed if error was "not found"; any other error is unexpected.
	if !isNotFound(err) {
		return socialauth.SocialAccount{}, err
	}

	now := time.Now().UTC()
	sa := socialauth.SocialAccount{
		ID:             kernel.SocialAccountID(uuid.NewString()),
		TenantID:       tenantID,
		CustomerID:     input.CustomerID,
		Provider:       input.Provider,
		ProviderUserID: input.ProviderUserID,
		Email:          input.Email,
		Name:           input.Name,
		AvatarURL:      input.AvatarURL,
		AccessToken:    input.AccessToken,
		RefreshToken:   input.RefreshToken,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	return s.repo.Create(ctx, sa)
}

// UnlinkAccount removes a linked social account by ID.
func (s *Service) UnlinkAccount(ctx context.Context, tenantID kernel.TenantID, id kernel.SocialAccountID) error {
	return s.repo.Delete(ctx, tenantID, id)
}

// FindByProvider looks up a social account by provider and provider user ID.
func (s *Service) FindByProvider(ctx context.Context, tenantID kernel.TenantID, provider, providerUserID string) (socialauth.SocialAccount, error) {
	if !socialauth.SupportedProviders[provider] {
		return socialauth.SocialAccount{}, socialauth.ErrProviderNotSupported()
	}
	return s.repo.GetByProvider(ctx, tenantID, provider, providerUserID)
}

// ListByCustomer returns all linked social accounts for a customer.
func (s *Service) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) ([]socialauth.SocialAccount, error) {
	return s.repo.ListByCustomer(ctx, tenantID, customerID)
}

// GetProviderConfig returns the list of supported OAuth provider configurations.
func (s *Service) GetProviderConfig() []socialauth.ProviderConfig {
	return socialauth.SupportedProviderConfigs()
}

// HandleCallback handles an incoming OAuth callback code.
// In this implementation we do NOT perform an actual token exchange with the
// provider — that is delegated to the frontend or an edge function.
// This method validates that the provider is supported and returns
// ErrInvalidOAuthCode when the code is empty (basic guard).
func (s *Service) HandleCallback(_ context.Context, input socialauth.OAuthCallbackInput) error {
	if !socialauth.SupportedProviders[input.Provider] {
		return socialauth.ErrProviderNotSupported()
	}
	if input.Code == "" {
		return socialauth.ErrInvalidOAuthCode()
	}
	return nil
}

// List returns a paginated list of all social accounts for the tenant (admin).
func (s *Service) List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[socialauth.SocialAccount], error) {
	return s.repo.List(ctx, tenantID, pg)
}

// isNotFound returns true when err is a social_auth NOT_FOUND errx error.
func isNotFound(err error) bool {
	if err == nil {
		return false
	}
	var e *errx.Error
	if errors.As(err, &e) {
		return e.Code == socialauth.CodeNotFound.Code
	}
	return false
}
