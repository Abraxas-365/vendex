package socialauth

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// Repository defines the persistence interface for the social auth domain.
type Repository interface {
	// Create persists a new social account link.
	Create(ctx context.Context, sa SocialAccount) (SocialAccount, error)

	// GetByID retrieves a social account by its ID, scoped to the tenant.
	GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.SocialAccountID) (SocialAccount, error)

	// GetByProvider retrieves a social account by provider + provider_user_id, scoped to the tenant.
	GetByProvider(ctx context.Context, tenantID kernel.TenantID, provider, providerUserID string) (SocialAccount, error)

	// ListByCustomer returns all social accounts linked to a given customer, scoped to the tenant.
	ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) ([]SocialAccount, error)

	// List returns a paginated list of all social accounts for the tenant.
	List(ctx context.Context, tenantID kernel.TenantID, pg kernel.PaginationOptions) (kernel.Paginated[SocialAccount], error)

	// Delete removes a social account by ID, scoped to the tenant.
	Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.SocialAccountID) error
}
