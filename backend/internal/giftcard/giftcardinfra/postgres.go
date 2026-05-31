package giftcardinfra

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/giftcard"
	"github.com/Abraxas-365/vendex/internal/kernel"
)

// PostgresRepository implements giftcard.Repository using PostgreSQL.
type PostgresRepository struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgresRepository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create inserts a new gift card row.
func (r *PostgresRepository) Create(ctx context.Context, gc *giftcard.GiftCard) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO gift_cards
			(id, tenant_id, code,
			 initial_amount_cents, initial_amount_currency,
			 balance_cents, balance_currency,
			 expires_at, active, created_by, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		string(gc.ID), string(gc.TenantID), gc.Code,
		gc.InitialAmount.Amount, gc.InitialAmount.Currency,
		gc.Balance.Amount, gc.Balance.Currency,
		gc.ExpiresAt, gc.Active, gc.CreatedBy, gc.CreatedAt, gc.UpdatedAt,
	)
	if err != nil {
		if isDuplicateCode(err) {
			return giftcard.ErrDuplicateCode
		}
		return errx.Wrap(err, "insert gift card", errx.TypeInternal)
	}
	return nil
}

// GetByID retrieves a gift card by primary key, scoped to the tenant.
func (r *PostgresRepository) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.GiftCardID) (*giftcard.GiftCard, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, code,
		       initial_amount_cents, initial_amount_currency,
		       balance_cents, balance_currency,
		       expires_at, active, created_by, created_at, updated_at
		FROM gift_cards WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	return scanGiftCard(row)
}

// GetByCode retrieves a gift card by its code (case-insensitive), scoped to the tenant.
func (r *PostgresRepository) GetByCode(ctx context.Context, tenantID kernel.TenantID, code string) (*giftcard.GiftCard, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, code,
		       initial_amount_cents, initial_amount_currency,
		       balance_cents, balance_currency,
		       expires_at, active, created_by, created_at, updated_at
		FROM gift_cards WHERE tenant_id=$1 AND UPPER(code)=UPPER($2)`,
		string(tenantID), code,
	)
	return scanGiftCard(row)
}

// List returns gift cards for a tenant with pagination.
func (r *PostgresRepository) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[giftcard.GiftCard], error) {
	var zero kernel.Paginated[giftcard.GiftCard]

	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM gift_cards WHERE tenant_id=$1`, string(tenantID)).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "count gift cards", errx.TypeInternal)
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, code,
		       initial_amount_cents, initial_amount_currency,
		       balance_cents, balance_currency,
		       expires_at, active, created_by, created_at, updated_at
		FROM gift_cards WHERE tenant_id=$1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return zero, errx.Wrap(err, "list gift cards", errx.TypeInternal)
	}
	defer rows.Close()

	var items []giftcard.GiftCard
	for rows.Next() {
		gc, err := scanGiftCardRow(rows)
		if err != nil {
			return zero, err
		}
		items = append(items, *gc)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterate gift cards", errx.TypeInternal)
	}
	return kernel.NewPaginated(items, p.Page, p.PageSize, total), nil
}

// Update persists mutations to an existing gift card.
func (r *PostgresRepository) Update(ctx context.Context, gc *giftcard.GiftCard) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE gift_cards
		SET code=$3,
		    initial_amount_cents=$4, initial_amount_currency=$5,
		    balance_cents=$6, balance_currency=$7,
		    expires_at=$8, active=$9, created_by=$10, updated_at=$11
		WHERE tenant_id=$1 AND id=$2`,
		string(gc.TenantID), string(gc.ID),
		gc.Code,
		gc.InitialAmount.Amount, gc.InitialAmount.Currency,
		gc.Balance.Amount, gc.Balance.Currency,
		gc.ExpiresAt, gc.Active, gc.CreatedBy, gc.UpdatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "update gift card", errx.TypeInternal)
	}
	return nil
}

// Delete removes a gift card by ID.
func (r *PostgresRepository) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.GiftCardID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM gift_cards WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	if err != nil {
		return errx.Wrap(err, "delete gift card", errx.TypeInternal)
	}
	return nil
}

// CreateTransaction inserts a new gift card transaction row.
func (r *PostgresRepository) CreateTransaction(ctx context.Context, tx *giftcard.GiftCardTransaction) error {
	var orderID sql.NullString
	if tx.OrderID != nil {
		orderID = sql.NullString{String: string(*tx.OrderID), Valid: true}
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO gift_card_transactions
			(id, gift_card_id, tenant_id, type,
			 amount_cents, amount_currency,
			 order_id, note, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		string(tx.ID), string(tx.GiftCardID), string(tx.TenantID),
		tx.Type,
		tx.Amount.Amount, tx.Amount.Currency,
		orderID, tx.Note, tx.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "insert gift card transaction", errx.TypeInternal)
	}
	return nil
}

// ListTransactions returns all transactions for a given gift card, ordered by creation time.
func (r *PostgresRepository) ListTransactions(ctx context.Context, tenantID kernel.TenantID, giftCardID kernel.GiftCardID) ([]giftcard.GiftCardTransaction, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, gift_card_id, tenant_id, type,
		       amount_cents, amount_currency,
		       order_id, note, created_at
		FROM gift_card_transactions
		WHERE tenant_id=$1 AND gift_card_id=$2
		ORDER BY created_at DESC`,
		string(tenantID), string(giftCardID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "list gift card transactions", errx.TypeInternal)
	}
	defer rows.Close()

	var txs []giftcard.GiftCardTransaction
	for rows.Next() {
		tx, err := scanTransactionRow(rows)
		if err != nil {
			return nil, err
		}
		txs = append(txs, *tx)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterate gift card transactions", errx.TypeInternal)
	}
	return txs, nil
}

// ---------------------------------------------------------------------------
// Scan helpers
// ---------------------------------------------------------------------------

func scanGiftCard(row *sql.Row) (*giftcard.GiftCard, error) {
	var gc giftcard.GiftCard
	var idStr, tenantStr string
	var expiresAt sql.NullTime
	var updatedAt time.Time

	err := row.Scan(
		&idStr, &tenantStr, &gc.Code,
		&gc.InitialAmount.Amount, &gc.InitialAmount.Currency,
		&gc.Balance.Amount, &gc.Balance.Currency,
		&expiresAt, &gc.Active, &gc.CreatedBy, &gc.CreatedAt, &updatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, giftcard.ErrNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scan gift card", errx.TypeInternal)
	}

	gc.ID = kernel.GiftCardID(idStr)
	gc.TenantID = kernel.TenantID(tenantStr)
	gc.UpdatedAt = updatedAt
	if expiresAt.Valid {
		t := expiresAt.Time
		gc.ExpiresAt = &t
	}
	return &gc, nil
}

func scanGiftCardRow(rows *sql.Rows) (*giftcard.GiftCard, error) {
	var gc giftcard.GiftCard
	var idStr, tenantStr string
	var expiresAt sql.NullTime
	var updatedAt time.Time

	err := rows.Scan(
		&idStr, &tenantStr, &gc.Code,
		&gc.InitialAmount.Amount, &gc.InitialAmount.Currency,
		&gc.Balance.Amount, &gc.Balance.Currency,
		&expiresAt, &gc.Active, &gc.CreatedBy, &gc.CreatedAt, &updatedAt,
	)
	if err != nil {
		return nil, errx.Wrap(err, "scan gift card row", errx.TypeInternal)
	}

	gc.ID = kernel.GiftCardID(idStr)
	gc.TenantID = kernel.TenantID(tenantStr)
	gc.UpdatedAt = updatedAt
	if expiresAt.Valid {
		t := expiresAt.Time
		gc.ExpiresAt = &t
	}
	return &gc, nil
}

func scanTransactionRow(rows *sql.Rows) (*giftcard.GiftCardTransaction, error) {
	var tx giftcard.GiftCardTransaction
	var idStr, cardIDStr, tenantStr string
	var orderID sql.NullString

	err := rows.Scan(
		&idStr, &cardIDStr, &tenantStr, &tx.Type,
		&tx.Amount.Amount, &tx.Amount.Currency,
		&orderID, &tx.Note, &tx.CreatedAt,
	)
	if err != nil {
		return nil, errx.Wrap(err, "scan gift card transaction row", errx.TypeInternal)
	}

	tx.ID = kernel.GiftCardTransactionID(idStr)
	tx.GiftCardID = kernel.GiftCardID(cardIDStr)
	tx.TenantID = kernel.TenantID(tenantStr)
	if orderID.Valid {
		oid := kernel.OrderID(orderID.String)
		tx.OrderID = &oid
	}
	return &tx, nil
}

// isDuplicateCode detects unique constraint violations on the code column.
func isDuplicateCode(err error) bool {
	return strings.Contains(err.Error(), "unique") || strings.Contains(err.Error(), "duplicate")
}

// Ensure interface compliance.
var _ giftcard.Repository = (*PostgresRepository)(nil)
