package reviewinfra

import (
	"context"
	"database/sql"
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/review"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// PostgresRepo implements review.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed review repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance.
var _ review.Repository = (*PostgresRepo)(nil)

// Create inserts a new review row.
func (r *PostgresRepo) Create(ctx context.Context, rv review.Review) (review.Review, error) {
	const q = `
		INSERT INTO reviews
			(id, tenant_id, product_id, customer_id, rating, title, body, status,
			 verified_purchase, helpful_count, images, admin_response, admin_responded_at,
			 created_at, updated_at)
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
		RETURNING id, tenant_id, product_id, customer_id, rating,
		          COALESCE(title, ''), COALESCE(body, ''), status,
		          verified_purchase, helpful_count,
		          COALESCE(images, '{}'), admin_response, admin_responded_at,
		          created_at, updated_at`

	row := r.db.QueryRowContext(ctx, q,
		string(rv.ID),
		string(rv.TenantID),
		string(rv.ProductID),
		string(rv.CustomerID),
		rv.Rating,
		nullableString(rv.Title),
		nullableString(rv.Body),
		string(rv.Status),
		rv.VerifiedPurchase,
		rv.HelpfulCount,
		pq.Array(rv.Images),
		rv.AdminResponse,
		rv.AdminRespondedAt,
		rv.CreatedAt,
		rv.UpdatedAt,
	)
	return scanReview(row)
}

// GetByID retrieves a review by ID scoped to tenant.
func (r *PostgresRepo) GetByID(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID) (review.Review, error) {
	const q = `
		SELECT id, tenant_id, product_id, customer_id, rating,
		       COALESCE(title, ''), COALESCE(body, ''), status,
		       verified_purchase, helpful_count,
		       COALESCE(images, '{}'), admin_response, admin_responded_at,
		       created_at, updated_at
		FROM reviews
		WHERE id = $1 AND tenant_id = $2`

	row := r.db.QueryRowContext(ctx, q, string(id), string(tenantID))
	rv, err := scanReview(row)
	if err == sql.ErrNoRows {
		return review.Review{}, review.ErrNotFound
	}
	return rv, err
}

// ListByProduct returns paginated reviews for a product, optionally filtered by status.
func (r *PostgresRepo) ListByProduct(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID, status string, pg kernel.PaginationOptions) (kernel.Paginated[review.Review], error) {
	args := []any{string(tenantID), string(productID)}
	whereStatus := ""
	if status != "" {
		args = append(args, status)
		whereStatus = " AND status = $3"
	}

	var total int
	countQ := "SELECT COUNT(*) FROM reviews WHERE tenant_id = $1 AND product_id = $2" + whereStatus
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return kernel.Paginated[review.Review]{}, errx.Wrap(err, "counting product reviews", errx.TypeInternal)
	}

	limitIdx := len(args) + 1
	offsetIdx := limitIdx + 1
	dataQ := `
		SELECT id, tenant_id, product_id, customer_id, rating,
		       COALESCE(title, ''), COALESCE(body, ''), status,
		       verified_purchase, helpful_count,
		       COALESCE(images, '{}'), admin_response, admin_responded_at,
		       created_at, updated_at
		FROM reviews
		WHERE tenant_id = $1 AND product_id = $2` + whereStatus + `
		ORDER BY created_at DESC
		LIMIT $` + itoa(limitIdx) + ` OFFSET $` + itoa(offsetIdx)

	args = append(args, pg.Limit(), pg.Offset())
	rows, err := r.db.QueryContext(ctx, dataQ, args...)
	if err != nil {
		return kernel.Paginated[review.Review]{}, errx.Wrap(err, "querying product reviews", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRows(rows, pg, total)
}

// ListByCustomer returns paginated reviews submitted by a customer.
func (r *PostgresRepo) ListByCustomer(ctx context.Context, tenantID kernel.TenantID, customerID kernel.CustomerID, pg kernel.PaginationOptions) (kernel.Paginated[review.Review], error) {
	const countQ = `SELECT COUNT(*) FROM reviews WHERE tenant_id = $1 AND customer_id = $2`
	var total int
	if err := r.db.QueryRowContext(ctx, countQ, string(tenantID), string(customerID)).Scan(&total); err != nil {
		return kernel.Paginated[review.Review]{}, errx.Wrap(err, "counting customer reviews", errx.TypeInternal)
	}

	const dataQ = `
		SELECT id, tenant_id, product_id, customer_id, rating,
		       COALESCE(title, ''), COALESCE(body, ''), status,
		       verified_purchase, helpful_count,
		       COALESCE(images, '{}'), admin_response, admin_responded_at,
		       created_at, updated_at
		FROM reviews
		WHERE tenant_id = $1 AND customer_id = $2
		ORDER BY created_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.QueryContext(ctx, dataQ, string(tenantID), string(customerID), pg.Limit(), pg.Offset())
	if err != nil {
		return kernel.Paginated[review.Review]{}, errx.Wrap(err, "querying customer reviews", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRows(rows, pg, total)
}

// List returns all reviews for a tenant, optionally filtered by status.
func (r *PostgresRepo) List(ctx context.Context, tenantID kernel.TenantID, status string, pg kernel.PaginationOptions) (kernel.Paginated[review.Review], error) {
	args := []any{string(tenantID)}
	whereStatus := ""
	if status != "" {
		args = append(args, status)
		whereStatus = " AND status = $2"
	}

	var total int
	countQ := "SELECT COUNT(*) FROM reviews WHERE tenant_id = $1" + whereStatus
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return kernel.Paginated[review.Review]{}, errx.Wrap(err, "counting reviews", errx.TypeInternal)
	}

	limitIdx := len(args) + 1
	offsetIdx := limitIdx + 1
	dataQ := `
		SELECT id, tenant_id, product_id, customer_id, rating,
		       COALESCE(title, ''), COALESCE(body, ''), status,
		       verified_purchase, helpful_count,
		       COALESCE(images, '{}'), admin_response, admin_responded_at,
		       created_at, updated_at
		FROM reviews
		WHERE tenant_id = $1` + whereStatus + `
		ORDER BY created_at DESC
		LIMIT $` + itoa(limitIdx) + ` OFFSET $` + itoa(offsetIdx)

	args = append(args, pg.Limit(), pg.Offset())
	rows, err := r.db.QueryContext(ctx, dataQ, args...)
	if err != nil {
		return kernel.Paginated[review.Review]{}, errx.Wrap(err, "querying reviews", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRows(rows, pg, total)
}

// UpdateStatus changes the moderation status of a review and returns the updated row.
func (r *PostgresRepo) UpdateStatus(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID, status review.ReviewStatus) (review.Review, error) {
	const q = `
		UPDATE reviews SET status = $1, updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
		RETURNING id, tenant_id, product_id, customer_id, rating,
		          COALESCE(title, ''), COALESCE(body, ''), status,
		          verified_purchase, helpful_count,
		          COALESCE(images, '{}'), admin_response, admin_responded_at,
		          created_at, updated_at`

	row := r.db.QueryRowContext(ctx, q, string(status), string(id), string(tenantID))
	rv, err := scanReview(row)
	if err == sql.ErrNoRows {
		return review.Review{}, review.ErrNotFound
	}
	return rv, err
}

// GetStats returns aggregated rating statistics for a product (approved reviews only).
func (r *PostgresRepo) GetStats(ctx context.Context, tenantID kernel.TenantID, productID kernel.ProductID) (review.ReviewStats, error) {
	const q = `
		SELECT rating, COUNT(*) AS cnt
		FROM reviews
		WHERE tenant_id = $1 AND product_id = $2 AND status = 'approved'
		GROUP BY rating`

	rows, err := r.db.QueryContext(ctx, q, string(tenantID), string(productID))
	if err != nil {
		return review.ReviewStats{}, errx.Wrap(err, "getting review stats", errx.TypeInternal)
	}
	defer rows.Close()

	dist := make(map[int]int)
	total := 0
	var weightedSum int64

	for rows.Next() {
		var rating, cnt int
		if err := rows.Scan(&rating, &cnt); err != nil {
			return review.ReviewStats{}, errx.Wrap(err, "scanning stats row", errx.TypeInternal)
		}
		dist[rating] = cnt
		total += cnt
		weightedSum += int64(rating * cnt)
	}
	if err := rows.Err(); err != nil {
		return review.ReviewStats{}, errx.Wrap(err, "iterating stats rows", errx.TypeInternal)
	}

	var avg float64
	if total > 0 {
		avg = float64(weightedSum) / float64(total)
	}

	return review.ReviewStats{
		AverageRating: avg,
		TotalReviews:  total,
		Distribution:  dist,
	}, nil
}

// IncrementHelpful atomically increments helpful_count.
func (r *PostgresRepo) IncrementHelpful(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID) error {
	const q = `UPDATE reviews SET helpful_count = helpful_count + 1, updated_at = NOW() WHERE id = $1 AND tenant_id = $2`
	_, err := r.db.ExecContext(ctx, q, string(id), string(tenantID))
	if err != nil {
		return errx.Wrap(err, "incrementing helpful count", errx.TypeInternal)
	}
	return nil
}

// SetAdminResponse stores an admin response and timestamps it.
func (r *PostgresRepo) SetAdminResponse(ctx context.Context, tenantID kernel.TenantID, id kernel.ReviewID, response string) (review.Review, error) {
	const q = `
		UPDATE reviews
		SET admin_response = $1, admin_responded_at = NOW(), updated_at = NOW()
		WHERE id = $2 AND tenant_id = $3
		RETURNING id, tenant_id, product_id, customer_id, rating,
		          COALESCE(title, ''), COALESCE(body, ''), status,
		          verified_purchase, helpful_count,
		          COALESCE(images, '{}'), admin_response, admin_responded_at,
		          created_at, updated_at`

	row := r.db.QueryRowContext(ctx, q, response, string(id), string(tenantID))
	rv, err := scanReview(row)
	if err == sql.ErrNoRows {
		return review.Review{}, review.ErrNotFound
	}
	return rv, err
}

// ─── helpers ──────────────────────────────────────────────────────────────────

// rowScanner unifies *sql.Row and *sql.Rows for scanning.
type rowScanner interface {
	Scan(dest ...any) error
}

func scanReview(row rowScanner) (review.Review, error) {
	var rv review.Review
	var id, tenantID, productID, customerID, status string
	var images pq.StringArray
	var adminResponse sql.NullString
	var adminRespondedAt sql.NullTime
	var createdAt, updatedAt time.Time

	err := row.Scan(
		&id, &tenantID, &productID, &customerID,
		&rv.Rating, &rv.Title, &rv.Body, &status,
		&rv.VerifiedPurchase, &rv.HelpfulCount,
		&images, &adminResponse, &adminRespondedAt,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return review.Review{}, err
	}

	rv.ID = kernel.ReviewID(id)
	rv.TenantID = kernel.TenantID(tenantID)
	rv.ProductID = kernel.ProductID(productID)
	rv.CustomerID = kernel.CustomerID(customerID)
	rv.Status = review.ReviewStatus(status)
	rv.Images = []string(images)
	rv.CreatedAt = createdAt
	rv.UpdatedAt = updatedAt

	if adminResponse.Valid {
		s := adminResponse.String
		rv.AdminResponse = &s
	}
	if adminRespondedAt.Valid {
		t := adminRespondedAt.Time
		rv.AdminRespondedAt = &t
	}

	if rv.Images == nil {
		rv.Images = []string{}
	}
	return rv, nil
}

func scanRows(rows *sql.Rows, pg kernel.PaginationOptions, total int) (kernel.Paginated[review.Review], error) {
	var items []review.Review
	for rows.Next() {
		var rv review.Review
		var id, tenantID, productID, customerID, status string
		var images pq.StringArray
		var adminResponse sql.NullString
		var adminRespondedAt sql.NullTime
		var createdAt, updatedAt time.Time

		if err := rows.Scan(
			&id, &tenantID, &productID, &customerID,
			&rv.Rating, &rv.Title, &rv.Body, &status,
			&rv.VerifiedPurchase, &rv.HelpfulCount,
			&images, &adminResponse, &adminRespondedAt,
			&createdAt, &updatedAt,
		); err != nil {
			return kernel.Paginated[review.Review]{}, errx.Wrap(err, "scanning review row", errx.TypeInternal)
		}

		rv.ID = kernel.ReviewID(id)
		rv.TenantID = kernel.TenantID(tenantID)
		rv.ProductID = kernel.ProductID(productID)
		rv.CustomerID = kernel.CustomerID(customerID)
		rv.Status = review.ReviewStatus(status)
		rv.Images = []string(images)
		rv.CreatedAt = createdAt
		rv.UpdatedAt = updatedAt

		if adminResponse.Valid {
			s := adminResponse.String
			rv.AdminResponse = &s
		}
		if adminRespondedAt.Valid {
			t := adminRespondedAt.Time
			rv.AdminRespondedAt = &t
		}

		if rv.Images == nil {
			rv.Images = []string{}
		}
		items = append(items, rv)
	}
	if err := rows.Err(); err != nil {
		return kernel.Paginated[review.Review]{}, errx.Wrap(err, "iterating review rows", errx.TypeInternal)
	}
	if items == nil {
		items = []review.Review{}
	}
	return kernel.NewPaginated(items, pg.Page, pg.PageSize, total), nil
}

// nullableString converts empty string to nil for SQL NULL.
func nullableString(s string) any {
	if s == "" {
		return nil
	}
	return s
}

// itoa converts an int to string without strconv import.
func itoa(n int) string {
	// Simple positive integer to string
	if n == 0 {
		return "0"
	}
	buf := make([]byte, 0, 10)
	for n > 0 {
		buf = append([]byte{byte('0' + n%10)}, buf...)
		n /= 10
	}
	return string(buf)
}
