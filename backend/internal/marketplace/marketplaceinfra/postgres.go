package marketplaceinfra

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/marketplace"
)

// ─── Vendor Repository ────────────────────────────────────────────────────────

// PostgresVendorRepo implements marketplace.VendorRepository.
type PostgresVendorRepo struct{ db *sqlx.DB }

// NewPostgresVendorRepo creates a new PostgresVendorRepo.
func NewPostgresVendorRepo(db *sqlx.DB) *PostgresVendorRepo {
	return &PostgresVendorRepo{db: db}
}

// dbVendor is the sqlx-scannable row for a vendor.
type dbVendor struct {
	ID          string    `db:"id"`
	TenantID    string    `db:"tenant_id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Email       string    `db:"email"`
	Phone       string    `db:"phone"`
	Status      string    `db:"status"`
	Commission  float64   `db:"commission"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

func fromDBVendor(row dbVendor) marketplace.Vendor {
	return marketplace.Vendor{
		ID:          kernel.VendorID(row.ID),
		TenantID:    kernel.TenantID(row.TenantID),
		Name:        row.Name,
		Description: row.Description,
		Email:       row.Email,
		Phone:       row.Phone,
		Status:      marketplace.VendorStatus(row.Status),
		Commission:  row.Commission,
		CreatedAt:   row.CreatedAt,
		UpdatedAt:   row.UpdatedAt,
	}
}

func (r *PostgresVendorRepo) Create(ctx context.Context, v marketplace.Vendor) (marketplace.Vendor, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO marketplace_vendors
			(id, tenant_id, name, description, email, phone, status, commission, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		string(v.ID), string(v.TenantID), v.Name, v.Description,
		v.Email, v.Phone, string(v.Status), v.Commission, v.CreatedAt, v.UpdatedAt,
	)
	if err != nil {
		return marketplace.Vendor{}, errx.Wrap(err, "create vendor", errx.TypeInternal)
	}
	return v, nil
}

func (r *PostgresVendorRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.VendorID) (marketplace.Vendor, error) {
	var row dbVendor
	err := r.db.GetContext(ctx, &row, `
		SELECT id, tenant_id, name, description, email, phone, status, commission, created_at, updated_at
		FROM marketplace_vendors WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	if err == sql.ErrNoRows {
		return marketplace.Vendor{}, marketplace.ErrVendorNotFound
	}
	if err != nil {
		return marketplace.Vendor{}, errx.Wrap(err, "get vendor", errx.TypeInternal)
	}
	return fromDBVendor(row), nil
}

func (r *PostgresVendorRepo) Update(ctx context.Context, v marketplace.Vendor) (marketplace.Vendor, error) {
	_, err := r.db.ExecContext(ctx, `
		UPDATE marketplace_vendors
		SET name=$3, description=$4, email=$5, phone=$6, status=$7, commission=$8, updated_at=$9
		WHERE tenant_id=$1 AND id=$2`,
		string(v.TenantID), string(v.ID),
		v.Name, v.Description, v.Email, v.Phone, string(v.Status), v.Commission, v.UpdatedAt,
	)
	if err != nil {
		return marketplace.Vendor{}, errx.Wrap(err, "update vendor", errx.TypeInternal)
	}
	return v, nil
}

func (r *PostgresVendorRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id kernel.VendorID) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM marketplace_vendors WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), string(id),
	)
	if err != nil {
		return errx.Wrap(err, "delete vendor", errx.TypeInternal)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return marketplace.ErrVendorNotFound
	}
	return nil
}

func (r *PostgresVendorRepo) List(ctx context.Context, tenantID kernel.TenantID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.Vendor], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM marketplace_vendors WHERE tenant_id=$1`, string(tenantID)).Scan(&total); err != nil {
		return kernel.Paginated[marketplace.Vendor]{}, errx.Wrap(err, "count vendors", errx.TypeInternal)
	}

	var rows []dbVendor
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, name, description, email, phone, status, commission, created_at, updated_at
		FROM marketplace_vendors WHERE tenant_id=$1
		ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		string(tenantID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[marketplace.Vendor]{}, errx.Wrap(err, "list vendors", errx.TypeInternal)
	}

	vendors := make([]marketplace.Vendor, len(rows))
	for i, row := range rows {
		vendors[i] = fromDBVendor(row)
	}
	return kernel.NewPaginated(vendors, p.Page, p.PageSize, total), nil
}

// Ensure interface compliance.
var _ marketplace.VendorRepository = (*PostgresVendorRepo)(nil)

// ─── VendorProduct Repository ──────────────────────────────────────────────────

// PostgresVendorProductRepo implements marketplace.VendorProductRepository.
type PostgresVendorProductRepo struct{ db *sqlx.DB }

// NewPostgresVendorProductRepo creates a new PostgresVendorProductRepo.
func NewPostgresVendorProductRepo(db *sqlx.DB) *PostgresVendorProductRepo {
	return &PostgresVendorProductRepo{db: db}
}

// dbVendorProduct is the sqlx-scannable row for a vendor-product link.
type dbVendorProduct struct {
	ID          string    `db:"id"`
	TenantID    string    `db:"tenant_id"`
	VendorID    string    `db:"vendor_id"`
	ProductID   string    `db:"product_id"`
	PriceCents  int64     `db:"price_cents"`
	Currency    string    `db:"currency"`
	Stock       int       `db:"stock"`
	CreatedAt   time.Time `db:"created_at"`
}

func fromDBVendorProduct(row dbVendorProduct) marketplace.VendorProduct {
	return marketplace.VendorProduct{
		ID:        row.ID,
		TenantID:  kernel.TenantID(row.TenantID),
		VendorID:  kernel.VendorID(row.VendorID),
		ProductID: kernel.ProductID(row.ProductID),
		Price:     kernel.Money{Amount: row.PriceCents, Currency: row.Currency},
		Stock:     row.Stock,
		CreatedAt: row.CreatedAt,
	}
}

func (r *PostgresVendorProductRepo) Create(ctx context.Context, vp marketplace.VendorProduct) (marketplace.VendorProduct, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO marketplace_vendor_products
			(id, tenant_id, vendor_id, product_id, price_cents, currency, stock, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		vp.ID, string(vp.TenantID), string(vp.VendorID), string(vp.ProductID),
		vp.Price.Amount, vp.Price.Currency, vp.Stock, vp.CreatedAt,
	)
	if err != nil {
		return marketplace.VendorProduct{}, errx.Wrap(err, "create vendor product", errx.TypeInternal)
	}
	return vp, nil
}

func (r *PostgresVendorProductRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id string) (marketplace.VendorProduct, error) {
	var row dbVendorProduct
	err := r.db.GetContext(ctx, &row, `
		SELECT id, tenant_id, vendor_id, product_id, price_cents, currency, stock, created_at
		FROM marketplace_vendor_products WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), id,
	)
	if err == sql.ErrNoRows {
		return marketplace.VendorProduct{}, marketplace.ErrVendorProductNotFound
	}
	if err != nil {
		return marketplace.VendorProduct{}, errx.Wrap(err, "get vendor product", errx.TypeInternal)
	}
	return fromDBVendorProduct(row), nil
}

func (r *PostgresVendorProductRepo) Delete(ctx context.Context, tenantID kernel.TenantID, id string) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM marketplace_vendor_products WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), id,
	)
	if err != nil {
		return errx.Wrap(err, "delete vendor product", errx.TypeInternal)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return marketplace.ErrVendorProductNotFound
	}
	return nil
}

func (r *PostgresVendorProductRepo) ListByVendor(ctx context.Context, tenantID kernel.TenantID, vendorID kernel.VendorID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.VendorProduct], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM marketplace_vendor_products WHERE tenant_id=$1 AND vendor_id=$2`, string(tenantID), string(vendorID)).Scan(&total); err != nil {
		return kernel.Paginated[marketplace.VendorProduct]{}, errx.Wrap(err, "count vendor products", errx.TypeInternal)
	}

	var rows []dbVendorProduct
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, vendor_id, product_id, price_cents, currency, stock, created_at
		FROM marketplace_vendor_products WHERE tenant_id=$1 AND vendor_id=$2
		ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
		string(tenantID), string(vendorID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[marketplace.VendorProduct]{}, errx.Wrap(err, "list vendor products", errx.TypeInternal)
	}

	vps := make([]marketplace.VendorProduct, len(rows))
	for i, row := range rows {
		vps[i] = fromDBVendorProduct(row)
	}
	return kernel.NewPaginated(vps, p.Page, p.PageSize, total), nil
}

// Ensure interface compliance.
var _ marketplace.VendorProductRepository = (*PostgresVendorProductRepo)(nil)

// ─── VendorOrder Repository ────────────────────────────────────────────────────

// PostgresVendorOrderRepo implements marketplace.VendorOrderRepository.
type PostgresVendorOrderRepo struct{ db *sqlx.DB }

// NewPostgresVendorOrderRepo creates a new PostgresVendorOrderRepo.
func NewPostgresVendorOrderRepo(db *sqlx.DB) *PostgresVendorOrderRepo {
	return &PostgresVendorOrderRepo{db: db}
}

// dbVendorOrder is the sqlx-scannable row for a vendor order.
type dbVendorOrder struct {
	ID              string    `db:"id"`
	TenantID        string    `db:"tenant_id"`
	VendorID        string    `db:"vendor_id"`
	OrderID         string    `db:"order_id"`
	AmountCents     int64     `db:"amount_cents"`
	Currency        string    `db:"currency"`
	CommissionCents int64     `db:"commission_cents"`
	Status          string    `db:"status"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

func fromDBVendorOrder(row dbVendorOrder) marketplace.VendorOrder {
	return marketplace.VendorOrder{
		ID:         row.ID,
		TenantID:   kernel.TenantID(row.TenantID),
		VendorID:   kernel.VendorID(row.VendorID),
		OrderID:    kernel.OrderID(row.OrderID),
		Amount:     kernel.Money{Amount: row.AmountCents, Currency: row.Currency},
		Commission: kernel.Money{Amount: row.CommissionCents, Currency: row.Currency},
		Status:     row.Status,
		CreatedAt:  row.CreatedAt,
		UpdatedAt:  row.UpdatedAt,
	}
}

func (r *PostgresVendorOrderRepo) Create(ctx context.Context, vo marketplace.VendorOrder) (marketplace.VendorOrder, error) {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO marketplace_vendor_orders
			(id, tenant_id, vendor_id, order_id, amount_cents, currency, commission_cents, status, created_at, updated_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`,
		vo.ID, string(vo.TenantID), string(vo.VendorID), string(vo.OrderID),
		vo.Amount.Amount, vo.Amount.Currency, vo.Commission.Amount,
		vo.Status, vo.CreatedAt, vo.UpdatedAt,
	)
	if err != nil {
		return marketplace.VendorOrder{}, errx.Wrap(err, "create vendor order", errx.TypeInternal)
	}
	return vo, nil
}

func (r *PostgresVendorOrderRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id string) (marketplace.VendorOrder, error) {
	var row dbVendorOrder
	err := r.db.GetContext(ctx, &row, `
		SELECT id, tenant_id, vendor_id, order_id, amount_cents, currency, commission_cents, status, created_at, updated_at
		FROM marketplace_vendor_orders WHERE tenant_id=$1 AND id=$2`,
		string(tenantID), id,
	)
	if err == sql.ErrNoRows {
		return marketplace.VendorOrder{}, marketplace.ErrVendorOrderNotFound
	}
	if err != nil {
		return marketplace.VendorOrder{}, errx.Wrap(err, "get vendor order", errx.TypeInternal)
	}
	return fromDBVendorOrder(row), nil
}

func (r *PostgresVendorOrderRepo) ListByVendor(ctx context.Context, tenantID kernel.TenantID, vendorID kernel.VendorID, p kernel.PaginationOptions) (kernel.Paginated[marketplace.VendorOrder], error) {
	var total int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM marketplace_vendor_orders WHERE tenant_id=$1 AND vendor_id=$2`, string(tenantID), string(vendorID)).Scan(&total); err != nil {
		return kernel.Paginated[marketplace.VendorOrder]{}, errx.Wrap(err, "count vendor orders", errx.TypeInternal)
	}

	var rows []dbVendorOrder
	err := r.db.SelectContext(ctx, &rows, `
		SELECT id, tenant_id, vendor_id, order_id, amount_cents, currency, commission_cents, status, created_at, updated_at
		FROM marketplace_vendor_orders WHERE tenant_id=$1 AND vendor_id=$2
		ORDER BY created_at DESC LIMIT $3 OFFSET $4`,
		string(tenantID), string(vendorID), p.Limit(), p.Offset(),
	)
	if err != nil {
		return kernel.Paginated[marketplace.VendorOrder]{}, errx.Wrap(err, "list vendor orders", errx.TypeInternal)
	}

	vos := make([]marketplace.VendorOrder, len(rows))
	for i, row := range rows {
		vos[i] = fromDBVendorOrder(row)
	}
	return kernel.NewPaginated(vos, p.Page, p.PageSize, total), nil
}

// Ensure interface compliance.
var _ marketplace.VendorOrderRepository = (*PostgresVendorOrderRepo)(nil)
