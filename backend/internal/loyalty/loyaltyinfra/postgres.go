package loyaltyinfra

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/loyalty"
)

// PostgresRepository implements loyalty.Repository using PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// -----------------------------------------------------------------------------
// Account operations
// -----------------------------------------------------------------------------

// GetOrCreateAccount retrieves the account for (tenantID, customerID), creating it if absent.
func (r *PostgresRepository) GetOrCreateAccount(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (loyalty.LoyaltyAccount, error) {
	// Try fetch first.
	account, err := r.GetAccountByCustomerID(ctx, tenantID, customerID)
	if err == nil {
		return account, nil
	}
	if !errx.IsNotFound(err) {
		return loyalty.LoyaltyAccount{}, err
	}

	// Create a new account.
	now := time.Now().UTC()
	_, err = r.db.ExecContext(ctx, `
		INSERT INTO loyalty_accounts
			(tenant_id, customer_id, points_balance, tier, lifetime_points, created_at, updated_at)
		VALUES ($1, $2, 0, 'bronze', 0, $3, $3)
		ON CONFLICT (tenant_id, customer_id) DO NOTHING`,
		string(tenantID), string(customerID), now,
	)
	if err != nil {
		return loyalty.LoyaltyAccount{}, errx.Wrap(err, "create loyalty account", errx.TypeInternal)
	}

	return r.GetAccountByCustomerID(ctx, tenantID, customerID)
}

// GetAccountByID retrieves a loyalty account by its primary key.
func (r *PostgresRepository) GetAccountByID(ctx context.Context, tenantID kernel.TenantID, id kernel.LoyaltyAccountID) (loyalty.LoyaltyAccount, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, customer_id, points_balance, tier, lifetime_points, created_at, updated_at
		FROM loyalty_accounts
		WHERE tenant_id = $1 AND id = $2`,
		string(tenantID), string(id),
	)
	return scanAccount(row)
}

// GetAccountByCustomerID retrieves a loyalty account for a given customer.
func (r *PostgresRepository) GetAccountByCustomerID(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID) (loyalty.LoyaltyAccount, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, customer_id, points_balance, tier, lifetime_points, created_at, updated_at
		FROM loyalty_accounts
		WHERE tenant_id = $1 AND customer_id = $2`,
		string(tenantID), string(customerID),
	)
	return scanAccount(row)
}

// UpdateAccount persists changes to an existing loyalty account.
func (r *PostgresRepository) UpdateAccount(ctx context.Context, account loyalty.LoyaltyAccount) (loyalty.LoyaltyAccount, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE loyalty_accounts
		SET points_balance = $3, tier = $4, lifetime_points = $5, updated_at = $6
		WHERE tenant_id = $1 AND id = $2`,
		string(account.TenantID), string(account.ID),
		account.PointsBalance, account.Tier, account.LifetimePoints, account.UpdatedAt,
	)
	if err != nil {
		return loyalty.LoyaltyAccount{}, errx.Wrap(err, "update loyalty account", errx.TypeInternal)
	}
	return account, nil
}

// ListAccounts returns a paginated list of loyalty accounts for a tenant.
func (r *PostgresRepository) ListAccounts(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[loyalty.LoyaltyAccount], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM loyalty_accounts WHERE tenant_id = $1`, string(tenantID)).Scan(&total); err != nil {
		return kernel.Paginated[loyalty.LoyaltyAccount]{}, errx.Wrap(err, "count loyalty accounts", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, customer_id, points_balance, tier, lifetime_points, created_at, updated_at
		FROM loyalty_accounts
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[loyalty.LoyaltyAccount]{}, errx.Wrap(err, "list loyalty accounts", errx.TypeInternal)
	}
	defer rows.Close()

	var items []loyalty.LoyaltyAccount
	for rows.Next() {
		acc, err := scanAccountRow(rows)
		if err != nil {
			return kernel.Paginated[loyalty.LoyaltyAccount]{}, err
		}
		items = append(items, acc)
	}
	if err := rows.Err(); err != nil {
		return kernel.Paginated[loyalty.LoyaltyAccount]{}, errx.Wrap(err, "iterate loyalty accounts", errx.TypeInternal)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// -----------------------------------------------------------------------------
// Transaction operations
// -----------------------------------------------------------------------------

// CreateTransaction inserts a loyalty transaction record.
func (r *PostgresRepository) CreateTransaction(ctx context.Context, tx loyalty.LoyaltyTransaction) (loyalty.LoyaltyTransaction, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO loyalty_transactions
			(id, tenant_id, account_id, type, points, reference, note, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(tx.ID), string(tx.TenantID), string(tx.AccountID),
		tx.Type, tx.Points, tx.Reference, tx.Note, tx.CreatedAt,
	)
	if err != nil {
		return loyalty.LoyaltyTransaction{}, errx.Wrap(err, "create loyalty transaction", errx.TypeInternal)
	}
	return tx, nil
}

// ListTransactions returns a paginated list of transactions for an account.
func (r *PostgresRepository) ListTransactions(ctx context.Context, tenantID kernel.TenantID, accountID kernel.LoyaltyAccountID, p kernel.PaginationOptions) (kernel.Paginated[loyalty.LoyaltyTransaction], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM loyalty_transactions WHERE tenant_id = $1 AND account_id = $2`,
		string(tenantID), string(accountID),
	).Scan(&total); err != nil {
		return kernel.Paginated[loyalty.LoyaltyTransaction]{}, errx.Wrap(err, "count loyalty transactions", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, account_id, type, points, reference, note, created_at
		FROM loyalty_transactions
		WHERE tenant_id = $1 AND account_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`,
		string(tenantID), string(accountID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[loyalty.LoyaltyTransaction]{}, errx.Wrap(err, "list loyalty transactions", errx.TypeInternal)
	}
	defer rows.Close()

	var items []loyalty.LoyaltyTransaction
	for rows.Next() {
		tx, err := scanTransactionRow(rows)
		if err != nil {
			return kernel.Paginated[loyalty.LoyaltyTransaction]{}, err
		}
		items = append(items, tx)
	}
	if err := rows.Err(); err != nil {
		return kernel.Paginated[loyalty.LoyaltyTransaction]{}, errx.Wrap(err, "iterate loyalty transactions", errx.TypeInternal)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// -----------------------------------------------------------------------------
// Reward operations
// -----------------------------------------------------------------------------

// CreateReward inserts a new loyalty reward row.
func (r *PostgresRepository) CreateReward(ctx context.Context, reward loyalty.LoyaltyReward) (loyalty.LoyaltyReward, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO loyalty_rewards
			(id, tenant_id, name, description, points_cost, reward_type, value_cents, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		string(reward.ID), string(reward.TenantID),
		reward.Name, reward.Description, reward.PointsCost,
		reward.RewardType, reward.ValueCents, reward.Active,
		reward.CreatedAt, reward.UpdatedAt,
	)
	if err != nil {
		return loyalty.LoyaltyReward{}, errx.Wrap(err, "create loyalty reward", errx.TypeInternal)
	}
	return reward, nil
}

// GetRewardByID retrieves a reward by its primary key.
func (r *PostgresRepository) GetRewardByID(ctx context.Context, tenantID kernel.TenantID, id kernel.RewardID) (loyalty.LoyaltyReward, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, description, points_cost, reward_type, value_cents, active, created_at, updated_at
		FROM loyalty_rewards
		WHERE tenant_id = $1 AND id = $2`,
		string(tenantID), string(id),
	)
	return scanReward(row)
}

// UpdateReward persists changes to an existing loyalty reward.
func (r *PostgresRepository) UpdateReward(ctx context.Context, reward loyalty.LoyaltyReward) (loyalty.LoyaltyReward, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE loyalty_rewards
		SET name = $3, description = $4, points_cost = $5, reward_type = $6,
		    value_cents = $7, active = $8, updated_at = $9
		WHERE tenant_id = $1 AND id = $2`,
		string(reward.TenantID), string(reward.ID),
		reward.Name, reward.Description, reward.PointsCost,
		reward.RewardType, reward.ValueCents, reward.Active, reward.UpdatedAt,
	)
	if err != nil {
		return loyalty.LoyaltyReward{}, errx.Wrap(err, "update loyalty reward", errx.TypeInternal)
	}
	return reward, nil
}

// ListRewards returns a paginated list of rewards for a tenant.
func (r *PostgresRepository) ListRewards(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[loyalty.LoyaltyReward], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM loyalty_rewards WHERE tenant_id = $1`, string(tenantID)).Scan(&total); err != nil {
		return kernel.Paginated[loyalty.LoyaltyReward]{}, errx.Wrap(err, "count loyalty rewards", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, description, points_cost, reward_type, value_cents, active, created_at, updated_at
		FROM loyalty_rewards
		WHERE tenant_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[loyalty.LoyaltyReward]{}, errx.Wrap(err, "list loyalty rewards", errx.TypeInternal)
	}
	defer rows.Close()

	var items []loyalty.LoyaltyReward
	for rows.Next() {
		reward, err := scanRewardRow(rows)
		if err != nil {
			return kernel.Paginated[loyalty.LoyaltyReward]{}, err
		}
		items = append(items, reward)
	}
	if err := rows.Err(); err != nil {
		return kernel.Paginated[loyalty.LoyaltyReward]{}, errx.Wrap(err, "iterate loyalty rewards", errx.TypeInternal)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// -----------------------------------------------------------------------------
// Scan helpers
// -----------------------------------------------------------------------------

func scanAccount(row *sql.Row) (loyalty.LoyaltyAccount, error) {
	var acc loyalty.LoyaltyAccount
	var idStr, tenantStr, customerStr string
	err := row.Scan(
		&idStr, &tenantStr, &customerStr,
		&acc.PointsBalance, &acc.Tier, &acc.LifetimePoints,
		&acc.CreatedAt, &acc.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return loyalty.LoyaltyAccount{}, loyalty.ErrAccountNotFound
	}
	if err != nil {
		return loyalty.LoyaltyAccount{}, errx.Wrap(err, "scan loyalty account", errx.TypeInternal)
	}
	acc.ID = kernel.LoyaltyAccountID(idStr)
	acc.TenantID = kernel.TenantID(tenantStr)
	acc.CustomerID = kernel.CustomerID(customerStr)
	return acc, nil
}

func scanAccountRow(rows *sql.Rows) (loyalty.LoyaltyAccount, error) {
	var acc loyalty.LoyaltyAccount
	var idStr, tenantStr, customerStr string
	err := rows.Scan(
		&idStr, &tenantStr, &customerStr,
		&acc.PointsBalance, &acc.Tier, &acc.LifetimePoints,
		&acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		return loyalty.LoyaltyAccount{}, errx.Wrap(err, "scan loyalty account row", errx.TypeInternal)
	}
	acc.ID = kernel.LoyaltyAccountID(idStr)
	acc.TenantID = kernel.TenantID(tenantStr)
	acc.CustomerID = kernel.CustomerID(customerStr)
	return acc, nil
}

func scanTransactionRow(rows *sql.Rows) (loyalty.LoyaltyTransaction, error) {
	var tx loyalty.LoyaltyTransaction
	var idStr, tenantStr, accountStr string
	var reference, note sql.NullString
	err := rows.Scan(
		&idStr, &tenantStr, &accountStr,
		&tx.Type, &tx.Points, &reference, &note,
		&tx.CreatedAt,
	)
	if err != nil {
		return loyalty.LoyaltyTransaction{}, errx.Wrap(err, "scan loyalty transaction row", errx.TypeInternal)
	}
	tx.ID = kernel.LoyaltyTransactionID(idStr)
	tx.TenantID = kernel.TenantID(tenantStr)
	tx.AccountID = kernel.LoyaltyAccountID(accountStr)
	if reference.Valid {
		tx.Reference = reference.String
	}
	if note.Valid {
		tx.Note = note.String
	}
	return tx, nil
}

func scanReward(row *sql.Row) (loyalty.LoyaltyReward, error) {
	var reward loyalty.LoyaltyReward
	var idStr, tenantStr string
	var description sql.NullString
	err := row.Scan(
		&idStr, &tenantStr,
		&reward.Name, &description, &reward.PointsCost,
		&reward.RewardType, &reward.ValueCents, &reward.Active,
		&reward.CreatedAt, &reward.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return loyalty.LoyaltyReward{}, loyalty.ErrRewardNotFound
	}
	if err != nil {
		return loyalty.LoyaltyReward{}, errx.Wrap(err, "scan loyalty reward", errx.TypeInternal)
	}
	reward.ID = kernel.RewardID(idStr)
	reward.TenantID = kernel.TenantID(tenantStr)
	if description.Valid {
		reward.Description = description.String
	}
	return reward, nil
}

func scanRewardRow(rows *sql.Rows) (loyalty.LoyaltyReward, error) {
	var reward loyalty.LoyaltyReward
	var idStr, tenantStr string
	var description sql.NullString
	err := rows.Scan(
		&idStr, &tenantStr,
		&reward.Name, &description, &reward.PointsCost,
		&reward.RewardType, &reward.ValueCents, &reward.Active,
		&reward.CreatedAt, &reward.UpdatedAt,
	)
	if err != nil {
		return loyalty.LoyaltyReward{}, errx.Wrap(err, "scan loyalty reward row", errx.TypeInternal)
	}
	reward.ID = kernel.RewardID(idStr)
	reward.TenantID = kernel.TenantID(tenantStr)
	if description.Valid {
		reward.Description = description.String
	}
	return reward, nil
}

// Ensure interface compliance at compile time.
var _ loyalty.Repository = (*PostgresRepository)(nil)
