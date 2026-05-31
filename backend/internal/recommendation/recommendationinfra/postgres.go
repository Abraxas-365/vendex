package recommendationinfra

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/recommendation"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements recommendation.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepository creates a new PostgreSQL-backed recommendation repository.
func NewPostgresRepository(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Ensure interface compliance at compile time.
var _ recommendation.Repository = (*PostgresRepo)(nil)

// ---------------------------------------------------------------------------
// Tracking
// ---------------------------------------------------------------------------

// TrackView persists a product view event.
func (r *PostgresRepo) TrackView(ctx context.Context, view recommendation.ProductView) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO product_views (id, tenant_id, visitor_id, customer_id, product_id, source, viewed_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		string(view.ID),
		string(view.TenantID),
		view.VisitorID,
		nullString(view.CustomerID),
		view.ProductID,
		nullString(view.Source),
		view.ViewedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting product view", errx.TypeInternal)
	}
	return nil
}

// TrackInteraction persists a product interaction event.
func (r *PostgresRepo) TrackInteraction(ctx context.Context, interaction recommendation.ProductInteraction) error {
	metaJSON, err := json.Marshal(interaction.Metadata)
	if err != nil {
		return errx.Wrap(err, "marshalling interaction metadata", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO product_interactions (id, tenant_id, visitor_id, customer_id, product_id, interaction_type, metadata, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(interaction.ID),
		string(interaction.TenantID),
		interaction.VisitorID,
		nullString(interaction.CustomerID),
		interaction.ProductID,
		interaction.InteractionType,
		metaJSON,
		interaction.CreatedAt,
	)
	if err != nil {
		return errx.Wrap(err, "inserting product interaction", errx.TypeInternal)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Recommendations
// ---------------------------------------------------------------------------

// GetFrequentlyBoughtTogether returns products co-purchased with the given product.
// It joins order_items to find products that appear in the same orders.
func (r *PostgresRepo) GetFrequentlyBoughtTogether(ctx context.Context, tenantID kernel.TenantID, productID string, limit int) ([]recommendation.RecommendedProduct, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT oi2.product_id, COUNT(*) AS score
		FROM order_items oi1
		JOIN order_items oi2
		  ON oi1.order_id = oi2.order_id
		 AND oi1.tenant_id = oi2.tenant_id
		WHERE oi1.tenant_id = $1
		  AND oi1.product_id = $2
		  AND oi2.product_id != $2
		GROUP BY oi2.product_id
		ORDER BY score DESC
		LIMIT $3`,
		string(tenantID), productID, limit,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying frequently bought together", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRecommendedProducts(rows, "frequently_bought_together")
}

// GetTrending returns products with the highest interaction count within the given duration.
func (r *PostgresRepo) GetTrending(ctx context.Context, tenantID kernel.TenantID, limit int, since time.Duration) ([]recommendation.RecommendedProduct, error) {
	// Convert duration to a PostgreSQL interval string.
	interval := fmt.Sprintf("%d seconds", int(since.Seconds()))

	rows, err := r.db.QueryContext(ctx, `
		SELECT product_id, COUNT(*) AS score
		FROM product_interactions
		WHERE tenant_id = $1
		  AND created_at > NOW() - $2::interval
		GROUP BY product_id
		ORDER BY score DESC
		LIMIT $3`,
		string(tenantID), interval, limit,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying trending products", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRecommendedProducts(rows, "trending")
}

// GetRecentlyViewed returns the most recently viewed products for a visitor.
func (r *PostgresRepo) GetRecentlyViewed(ctx context.Context, tenantID kernel.TenantID, visitorID string, limit int) ([]recommendation.RecommendedProduct, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT ON (product_id) product_id, 1.0 AS score
		FROM product_views
		WHERE tenant_id = $1
		  AND visitor_id = $2
		ORDER BY product_id, viewed_at DESC
		LIMIT $3`,
		string(tenantID), visitorID, limit,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying recently viewed", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRecommendedProducts(rows, "recently_viewed")
}

// GetPersonalized returns recommended products based on the visitor's interaction history.
// It finds products that share a visitor with the current visitor's purchased/carted products.
func (r *PostgresRepo) GetPersonalized(ctx context.Context, tenantID kernel.TenantID, visitorID string, limit int) ([]recommendation.RecommendedProduct, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT i2.product_id, COUNT(*) AS score
		FROM product_interactions i1
		JOIN product_interactions i2
		  ON i1.product_id != i2.product_id
		 AND i1.tenant_id  = i2.tenant_id
		 AND i1.visitor_id != i2.visitor_id
		WHERE i1.tenant_id  = $1
		  AND i1.visitor_id = $2
		  AND i2.product_id NOT IN (
		      SELECT product_id
		      FROM product_interactions
		      WHERE tenant_id = $1 AND visitor_id = $2
		  )
		GROUP BY i2.product_id
		ORDER BY score DESC
		LIMIT $3`,
		string(tenantID), visitorID, limit,
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying personalized recommendations", errx.TypeInternal)
	}
	defer rows.Close()

	return scanRecommendedProducts(rows, "personalized")
}

// ---------------------------------------------------------------------------
// Rules
// ---------------------------------------------------------------------------

// CreateRule persists a new recommendation rule.
func (r *PostgresRepo) CreateRule(ctx context.Context, rule recommendation.RecommendationRule) (recommendation.RecommendationRule, error) {
	configJSON, err := json.Marshal(rule.Config)
	if err != nil {
		return recommendation.RecommendationRule{}, errx.Wrap(err, "marshalling rule config", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		INSERT INTO recommendation_rules (id, tenant_id, name, type, config, is_active, priority, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		string(rule.ID),
		string(rule.TenantID),
		rule.Name,
		rule.Type,
		configJSON,
		rule.IsActive,
		rule.Priority,
		rule.CreatedAt,
		rule.UpdatedAt,
	)
	if err != nil {
		return recommendation.RecommendationRule{}, errx.Wrap(err, "inserting recommendation rule", errx.TypeInternal)
	}
	return rule, nil
}

// GetRuleByID returns a recommendation rule scoped to a tenant.
func (r *PostgresRepo) GetRuleByID(ctx context.Context, tenantID kernel.TenantID, id kernel.RecommendationRuleID) (recommendation.RecommendationRule, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, type, config, is_active, priority, created_at, updated_at
		FROM recommendation_rules
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)

	rule, err := scanRule(row.Scan)
	if err == sql.ErrNoRows {
		return recommendation.RecommendationRule{}, recommendation.ErrRuleNotFound
	}
	if err != nil {
		return recommendation.RecommendationRule{}, errx.Wrap(err, "scanning recommendation rule", errx.TypeInternal)
	}
	return rule, nil
}

// ListRules returns all recommendation rules for a tenant.
func (r *PostgresRepo) ListRules(ctx context.Context, tenantID kernel.TenantID) ([]recommendation.RecommendationRule, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, type, config, is_active, priority, created_at, updated_at
		FROM recommendation_rules
		WHERE tenant_id = $1
		ORDER BY priority DESC, created_at DESC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "querying recommendation rules", errx.TypeInternal)
	}
	defer rows.Close()

	var rules []recommendation.RecommendationRule
	for rows.Next() {
		rule, err := scanRule(rows.Scan)
		if err != nil {
			return nil, errx.Wrap(err, "scanning recommendation rule row", errx.TypeInternal)
		}
		rules = append(rules, rule)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating recommendation rules", errx.TypeInternal)
	}
	if rules == nil {
		rules = []recommendation.RecommendationRule{}
	}
	return rules, nil
}

// UpdateRule replaces a recommendation rule's mutable fields.
func (r *PostgresRepo) UpdateRule(ctx context.Context, rule recommendation.RecommendationRule) (recommendation.RecommendationRule, error) {
	configJSON, err := json.Marshal(rule.Config)
	if err != nil {
		return recommendation.RecommendationRule{}, errx.Wrap(err, "marshalling rule config", errx.TypeInternal)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE recommendation_rules
		SET name = $1, type = $2, config = $3, is_active = $4, priority = $5, updated_at = $6
		WHERE id = $7 AND tenant_id = $8`,
		rule.Name,
		rule.Type,
		configJSON,
		rule.IsActive,
		rule.Priority,
		rule.UpdatedAt,
		string(rule.ID),
		string(rule.TenantID),
	)
	if err != nil {
		return recommendation.RecommendationRule{}, errx.Wrap(err, "updating recommendation rule", errx.TypeInternal)
	}
	return rule, nil
}

// DeleteRule removes a recommendation rule.
func (r *PostgresRepo) DeleteRule(ctx context.Context, tenantID kernel.TenantID, id kernel.RecommendationRuleID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM recommendation_rules WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting recommendation rule", errx.TypeInternal)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Scan helpers
// ---------------------------------------------------------------------------

type scanFunc func(dest ...any) error

func scanRule(scan scanFunc) (recommendation.RecommendationRule, error) {
	var rule recommendation.RecommendationRule
	var id, tenantID string
	var configJSON []byte

	err := scan(
		&id, &tenantID,
		&rule.Name, &rule.Type,
		&configJSON,
		&rule.IsActive, &rule.Priority,
		&rule.CreatedAt, &rule.UpdatedAt,
	)
	if err != nil {
		return recommendation.RecommendationRule{}, err
	}

	rule.ID = kernel.NewRecommendationRuleID(id)
	rule.TenantID = kernel.TenantID(tenantID)

	if len(configJSON) > 0 {
		if err := json.Unmarshal(configJSON, &rule.Config); err != nil {
			return recommendation.RecommendationRule{}, errx.Wrap(err, "unmarshalling rule config", errx.TypeInternal)
		}
	}

	return rule, nil
}

// scanRecommendedProducts reads (product_id, score) rows into a slice.
func scanRecommendedProducts(rows *sql.Rows, reason string) ([]recommendation.RecommendedProduct, error) {
	var items []recommendation.RecommendedProduct
	for rows.Next() {
		var p recommendation.RecommendedProduct
		if err := rows.Scan(&p.ProductID, &p.Score); err != nil {
			return nil, errx.Wrap(err, "scanning recommended product", errx.TypeInternal)
		}
		p.Reason = reason
		items = append(items, p)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating recommended products", errx.TypeInternal)
	}
	if items == nil {
		items = []recommendation.RecommendedProduct{}
	}
	return items, nil
}

// nullString converts an empty string to sql.NullString.
func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}
