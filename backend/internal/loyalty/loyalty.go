package loyalty

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Tier constants represent loyalty membership levels.
const (
	TierBronze   = "bronze"
	TierSilver   = "silver"
	TierGold     = "gold"
	TierPlatinum = "platinum"
)

// TransactionType constants describe the direction or reason for a points change.
const (
	TransactionTypeEarn   = "earn"
	TransactionTypeRedeem = "redeem"
	TransactionTypeExpire = "expire"
	TransactionTypeAdjust = "adjust"
)

// RewardType constants describe what a reward delivers.
const (
	RewardTypeDiscount     = "discount"
	RewardTypeFreeShipping = "free_shipping"
	RewardTypeGiftCard     = "gift_card"
)

// LoyaltyAccount holds a customer's accumulated points, tier status, and lifetime earnings.
type LoyaltyAccount struct {
	ID             kernel.LoyaltyAccountID `json:"id"              db:"id"`
	TenantID       kernel.TenantID         `json:"tenant_id"       db:"tenant_id"`
	CustomerID     kernel.CustomerID       `json:"customer_id"     db:"customer_id"`
	PointsBalance  int                     `json:"points_balance"  db:"points_balance"`
	Tier           string                  `json:"tier"            db:"tier"`
	LifetimePoints int                     `json:"lifetime_points" db:"lifetime_points"`
	CreatedAt      time.Time               `json:"created_at"      db:"created_at"`
	UpdatedAt      time.Time               `json:"updated_at"      db:"updated_at"`
}

// LoyaltyTransaction records a single points change against an account.
type LoyaltyTransaction struct {
	ID        kernel.LoyaltyTransactionID `json:"id"         db:"id"`
	TenantID  kernel.TenantID             `json:"tenant_id"  db:"tenant_id"`
	AccountID kernel.LoyaltyAccountID     `json:"account_id" db:"account_id"`
	Type      string                      `json:"type"       db:"type"`
	Points    int                         `json:"points"     db:"points"`
	Reference string                      `json:"reference"  db:"reference"`
	Note      string                      `json:"note"       db:"note"`
	CreatedAt time.Time                   `json:"created_at" db:"created_at"`
}

// LoyaltyReward describes a reward that can be redeemed using points.
type LoyaltyReward struct {
	ID          kernel.RewardID `json:"id"           db:"id"`
	TenantID    kernel.TenantID `json:"tenant_id"    db:"tenant_id"`
	Name        string          `json:"name"         db:"name"`
	Description string          `json:"description"  db:"description"`
	PointsCost  int             `json:"points_cost"  db:"points_cost"`
	RewardType  string          `json:"reward_type"  db:"reward_type"`
	ValueCents  int             `json:"value_cents"  db:"value_cents"`
	Active      bool            `json:"active"       db:"active"`
	CreatedAt   time.Time       `json:"created_at"   db:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"   db:"updated_at"`
}

// -----------------------------------------------------------------------------
// Input DTOs
// -----------------------------------------------------------------------------

// EarnPointsInput carries data needed to earn points for an account.
type EarnPointsInput struct {
	CustomerID kernel.CustomerID
	Points     int
	Reference  string
	Note       string
}

// RedeemPointsInput carries data needed to redeem points against a reward.
type RedeemPointsInput struct {
	CustomerID kernel.CustomerID
	RewardID   kernel.RewardID
	Note       string
}

// AdjustPointsInput carries data for a manual admin adjustment.
type AdjustPointsInput struct {
	Points    int
	Reference string
	Note      string
}

// CreateRewardInput holds data for creating a new reward.
type CreateRewardInput struct {
	Name        string
	Description string
	PointsCost  int
	RewardType  string
	ValueCents  int
}

// UpdateRewardInput holds optional fields that may be changed on a reward.
type UpdateRewardInput struct {
	Name        *string
	Description *string
	PointsCost  *int
	RewardType  *string
	ValueCents  *int
	Active      *bool
}
