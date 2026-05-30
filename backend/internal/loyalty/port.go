package loyalty

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines the data-access contract for the loyalty domain.
// All methods are tenant-scoped.
type Repository interface {
	// Account operations
	GetOrCreateAccount(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (LoyaltyAccount, error)
	GetAccountByID(ctx context.Context, tenantID kernel.TenantID, id kernel.LoyaltyAccountID) (LoyaltyAccount, error)
	GetAccountByCustomerID(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (LoyaltyAccount, error)
	UpdateAccount(ctx context.Context, account LoyaltyAccount) (LoyaltyAccount, error)
	ListAccounts(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[LoyaltyAccount], error)

	// Transaction operations
	CreateTransaction(ctx context.Context, tx LoyaltyTransaction) (LoyaltyTransaction, error)
	ListTransactions(ctx context.Context, tenantID kernel.TenantID, accountID kernel.LoyaltyAccountID, p kernel.PaginationOptions) (kernel.Paginated[LoyaltyTransaction], error)

	// Reward operations
	CreateReward(ctx context.Context, reward LoyaltyReward) (LoyaltyReward, error)
	GetRewardByID(ctx context.Context, tenantID kernel.TenantID, id kernel.RewardID) (LoyaltyReward, error)
	UpdateReward(ctx context.Context, reward LoyaltyReward) (LoyaltyReward, error)
	ListRewards(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[LoyaltyReward], error)
}
