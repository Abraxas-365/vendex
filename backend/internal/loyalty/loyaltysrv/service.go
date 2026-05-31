package loyaltysrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/loyalty"
	"github.com/google/uuid"
)

// Service implements the loyalty domain's business logic.
type Service struct {
	repo loyalty.Repository
	bus  eventbus.Bus
}

// New creates a new loyalty Service.
func New(repo loyalty.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// -----------------------------------------------------------------------------
// Account operations
// -----------------------------------------------------------------------------

// GetOrCreateAccount fetches the loyalty account for a customer, creating one if needed.
func (s *Service) GetOrCreateAccount(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (loyalty.LoyaltyAccount, error) {
	return s.repo.GetOrCreateAccount(ctx, tenantID, customerID)
}

// GetAccount fetches a loyalty account by its ID.
func (s *Service) GetAccount(ctx context.Context, tenantID kernel.TenantID, id kernel.LoyaltyAccountID) (loyalty.LoyaltyAccount, error) {
	return s.repo.GetAccountByID(ctx, tenantID, id)
}

// ListAccounts returns a paginated list of loyalty accounts for a tenant.
func (s *Service) ListAccounts(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[loyalty.LoyaltyAccount], error) {
	return s.repo.ListAccounts(ctx, tenantID, p)
}

// EarnPoints credits points to a customer's loyalty account and records the transaction.
func (s *Service) EarnPoints(ctx context.Context, tenantID kernel.TenantID, input loyalty.EarnPointsInput) (loyalty.LoyaltyAccount, error) {
	if input.Points <= 0 {
		return loyalty.LoyaltyAccount{}, loyalty.ErrInvalidPoints
	}

	account, err := s.repo.GetOrCreateAccount(ctx, tenantID, input.CustomerID)
	if err != nil {
		return loyalty.LoyaltyAccount{}, err
	}

	account.PointsBalance += input.Points
	account.LifetimePoints += input.Points
	account.Tier = computeTier(account.LifetimePoints)
	account.UpdatedAt = time.Now().UTC()

	account, err = s.repo.UpdateAccount(ctx, account)
	if err != nil {
		return loyalty.LoyaltyAccount{}, err
	}

	tx := loyalty.LoyaltyTransaction{
		ID:        kernel.NewLoyaltyTransactionID(uuid.NewString()),
		TenantID:  tenantID,
		AccountID: account.ID,
		Type:      loyalty.TransactionTypeEarn,
		Points:    input.Points,
		Reference: input.Reference,
		Note:      input.Note,
		CreatedAt: time.Now().UTC(),
	}
	if _, err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return loyalty.LoyaltyAccount{}, err
	}

	if s.bus != nil {
		event, _ := eventbus.NewEvent(eventbus.LoyaltyPointsEarned, tenantID, map[string]any{
			"account_id":  string(account.ID),
			"customer_id": string(input.CustomerID),
			"points":      input.Points,
		})
		_ = s.bus.Publish(ctx, event)
	}

	return account, nil
}

// RedeemPoints deducts points from a customer's account in exchange for a reward.
func (s *Service) RedeemPoints(ctx context.Context, tenantID kernel.TenantID, input loyalty.RedeemPointsInput) (loyalty.LoyaltyAccount, error) {
	reward, err := s.repo.GetRewardByID(ctx, tenantID, input.RewardID)
	if err != nil {
		return loyalty.LoyaltyAccount{}, err
	}
	if !reward.Active {
		return loyalty.LoyaltyAccount{}, loyalty.ErrRewardInactive
	}

	account, err := s.repo.GetOrCreateAccount(ctx, tenantID, input.CustomerID)
	if err != nil {
		return loyalty.LoyaltyAccount{}, err
	}

	if account.PointsBalance < reward.PointsCost {
		return loyalty.LoyaltyAccount{}, loyalty.ErrInsufficientPoints
	}

	account.PointsBalance -= reward.PointsCost
	account.UpdatedAt = time.Now().UTC()

	account, err = s.repo.UpdateAccount(ctx, account)
	if err != nil {
		return loyalty.LoyaltyAccount{}, err
	}

	tx := loyalty.LoyaltyTransaction{
		ID:        kernel.NewLoyaltyTransactionID(uuid.NewString()),
		TenantID:  tenantID,
		AccountID: account.ID,
		Type:      loyalty.TransactionTypeRedeem,
		Points:    -reward.PointsCost,
		Reference: string(input.RewardID),
		Note:      input.Note,
		CreatedAt: time.Now().UTC(),
	}
	if _, err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return loyalty.LoyaltyAccount{}, err
	}

	if s.bus != nil {
		event, _ := eventbus.NewEvent(eventbus.LoyaltyPointsRedeemed, tenantID, map[string]any{
			"account_id":  string(account.ID),
			"customer_id": string(input.CustomerID),
			"reward_id":   string(input.RewardID),
			"points":      reward.PointsCost,
		})
		_ = s.bus.Publish(ctx, event)
	}

	return account, nil
}

// AdjustPoints performs a manual admin adjustment (positive or negative).
func (s *Service) AdjustPoints(ctx context.Context, tenantID kernel.TenantID, accountID kernel.LoyaltyAccountID, input loyalty.AdjustPointsInput) (loyalty.LoyaltyAccount, error) {
	account, err := s.repo.GetAccountByID(ctx, tenantID, accountID)
	if err != nil {
		return loyalty.LoyaltyAccount{}, err
	}

	account.PointsBalance += input.Points
	if account.PointsBalance < 0 {
		account.PointsBalance = 0
	}
	// Only update lifetime points when adding
	if input.Points > 0 {
		account.LifetimePoints += input.Points
		account.Tier = computeTier(account.LifetimePoints)
	}
	account.UpdatedAt = time.Now().UTC()

	account, err = s.repo.UpdateAccount(ctx, account)
	if err != nil {
		return loyalty.LoyaltyAccount{}, err
	}

	tx := loyalty.LoyaltyTransaction{
		ID:        kernel.NewLoyaltyTransactionID(uuid.NewString()),
		TenantID:  tenantID,
		AccountID: account.ID,
		Type:      loyalty.TransactionTypeAdjust,
		Points:    input.Points,
		Reference: input.Reference,
		Note:      input.Note,
		CreatedAt: time.Now().UTC(),
	}
	if _, err := s.repo.CreateTransaction(ctx, tx); err != nil {
		return loyalty.LoyaltyAccount{}, err
	}

	return account, nil
}

// ListTransactions returns a paginated list of transactions for an account.
func (s *Service) ListTransactions(ctx context.Context, tenantID kernel.TenantID, accountID kernel.LoyaltyAccountID, p kernel.PaginationOptions) (kernel.Paginated[loyalty.LoyaltyTransaction], error) {
	return s.repo.ListTransactions(ctx, tenantID, accountID, p)
}

// -----------------------------------------------------------------------------
// Reward operations
// -----------------------------------------------------------------------------

// CreateReward creates a new redeemable reward for a tenant.
func (s *Service) CreateReward(ctx context.Context, tenantID kernel.TenantID, input loyalty.CreateRewardInput) (loyalty.LoyaltyReward, error) {
	now := time.Now().UTC()
	reward := loyalty.LoyaltyReward{
		ID:          kernel.NewRewardID(uuid.NewString()),
		TenantID:    tenantID,
		Name:        input.Name,
		Description: input.Description,
		PointsCost:  input.PointsCost,
		RewardType:  input.RewardType,
		ValueCents:  input.ValueCents,
		Active:      true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	reward, err := s.repo.CreateReward(ctx, reward)
	if err != nil {
		return loyalty.LoyaltyReward{}, err
	}

	if s.bus != nil {
		event, _ := eventbus.NewEvent(eventbus.LoyaltyRewardCreated, tenantID, map[string]any{
			"reward_id": string(reward.ID),
			"name":      reward.Name,
		})
		_ = s.bus.Publish(ctx, event)
	}

	return reward, nil
}

// GetRewardByID fetches a single reward.
func (s *Service) GetRewardByID(ctx context.Context, tenantID kernel.TenantID, id kernel.RewardID) (loyalty.LoyaltyReward, error) {
	return s.repo.GetRewardByID(ctx, tenantID, id)
}

// ListRewards returns a paginated list of rewards for a tenant.
func (s *Service) ListRewards(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[loyalty.LoyaltyReward], error) {
	return s.repo.ListRewards(ctx, tenantID, p)
}

// UpdateReward applies partial updates to a reward.
func (s *Service) UpdateReward(ctx context.Context, tenantID kernel.TenantID, id kernel.RewardID, input loyalty.UpdateRewardInput) (loyalty.LoyaltyReward, error) {
	reward, err := s.repo.GetRewardByID(ctx, tenantID, id)
	if err != nil {
		return loyalty.LoyaltyReward{}, err
	}

	if input.Name != nil {
		reward.Name = *input.Name
	}
	if input.Description != nil {
		reward.Description = *input.Description
	}
	if input.PointsCost != nil {
		reward.PointsCost = *input.PointsCost
	}
	if input.RewardType != nil {
		reward.RewardType = *input.RewardType
	}
	if input.ValueCents != nil {
		reward.ValueCents = *input.ValueCents
	}
	if input.Active != nil {
		reward.Active = *input.Active
	}
	reward.UpdatedAt = time.Now().UTC()

	return s.repo.UpdateReward(ctx, reward)
}

// -----------------------------------------------------------------------------
// Tier computation
// -----------------------------------------------------------------------------

// computeTier maps lifetime points to a tier name.
func computeTier(lifetimePoints int) string {
	switch {
	case lifetimePoints >= 10000:
		return loyalty.TierPlatinum
	case lifetimePoints >= 5000:
		return loyalty.TierGold
	case lifetimePoints >= 1000:
		return loyalty.TierSilver
	default:
		return loyalty.TierBronze
	}
}
