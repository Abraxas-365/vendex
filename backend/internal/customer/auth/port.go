package auth

import (
	"context"

	"github.com/Abraxas-365/vendex/internal/kernel"
)

// CredentialsRepository defines persistence operations for customer credentials.
type CredentialsRepository interface {
	Create(ctx context.Context, creds *CustomerCredentials) error
	GetByEmail(ctx context.Context, tenantID kernel.TenantID, email string) (*CustomerCredentials, error)
	GetByCustomerID(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (*CustomerCredentials, error)
	UpdatePassword(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, passwordHash string) error
}
