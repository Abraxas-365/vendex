package searchinfra

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/Abraxas-365/vendex/internal/search"
	"github.com/jmoiron/sqlx"
)

// PostgresRepo implements search.Repository using PostgreSQL full-text search.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed search repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// Search performs full-text product search with optional filters.
func (r *PostgresRepo) Search(ctx context.Context, tenantID kernel.TenantID, q search.SearchQuery) (*search.SearchResult, error) {
	args := []any{string(tenantID)}
	paramIdx := 2 // $1 is tenant_id

	var whereClauses []string
	whereClauses = append(whereClauses, "tenant_id = $1")

	// Full-text search on the pre-computed search_vector column.
	hasQuery := strings.TrimSpace(q.Query) != ""
	if hasQuery {
		whereClauses = append(whereClauses, fmt.Sprintf("search_vector @@ plainto_tsquery('english', $%d)", paramIdx))
		args = append(args, q.Query)
		paramIdx++
	}

	// Category filter.
	if q.CategoryID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("category_id = $%d", paramIdx))
		args = append(args, q.CategoryID)
		paramIdx++
	}

	// Tags filter (product must have ALL requested tags).
	if len(q.Tags) > 0 {
		tagsJSON, err := json.Marshal(q.Tags)
		if err != nil {
			return nil, errx.Wrap(err, "marshaling tags filter", errx.TypeInternal)
		}
		whereClauses = append(whereClauses, fmt.Sprintf("tags::jsonb @> $%d::jsonb", paramIdx))
		args = append(args, string(tagsJSON))
		paramIdx++
	}

	// Status filter.
	if q.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", paramIdx))
		args = append(args, q.Status)
		paramIdx++
	}

	// Price range filters.
	if q.MinPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("price_amount >= $%d", paramIdx))
		args = append(args, *q.MinPrice)
		paramIdx++
	}
	if q.MaxPrice != nil {
		whereClauses = append(whereClauses, fmt.Sprintf("price_amount <= $%d", paramIdx))
		args = append(args, *q.MaxPrice)
		paramIdx++
	}

	whereClause := strings.Join(whereClauses, " AND ")

	// Count query.
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM products WHERE %s", whereClause)
	var total int
	if err := r.db.QueryRowContext(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, errx.Wrap(err, "counting search results", errx.TypeInternal)
	}

	// Build ORDER BY clause.
	orderBy := buildOrderBy(q.SortBy, hasQuery)

	// Select query with rank when doing full-text search.
	var rankExpr string
	if hasQuery {
		rankExpr = fmt.Sprintf(", ts_rank(search_vector, plainto_tsquery('english', $%d)) AS rank", paramIdx)
		args = append(args, q.Query)
		paramIdx++
	} else {
		rankExpr = ", 0.0 AS rank"
	}

	dataSQL := fmt.Sprintf(`
		SELECT id, tenant_id, name, description, price_amount, price_currency,
		       sku, images, category_id, tags, status, stock, created_at%s
		FROM products
		WHERE %s
		ORDER BY %s
		LIMIT $%d OFFSET $%d`,
		rankExpr, whereClause, orderBy, paramIdx, paramIdx+1,
	)
	args = append(args, q.PageSize, (q.Page-1)*q.PageSize)

	rows, err := r.db.QueryContext(ctx, dataSQL, args...)
	if err != nil {
		return nil, errx.Wrap(err, "querying search results", errx.TypeInternal)
	}
	defer rows.Close()

	var hits []search.ProductHit
	for rows.Next() {
		var hit search.ProductHit
		var imagesJSON, tagsJSON, tenantIDStr string
		if err := rows.Scan(
			&hit.ID, &tenantIDStr, &hit.Name, &hit.Description,
			&hit.Price.Amount, &hit.Price.Currency,
			&hit.SKU, &imagesJSON, &hit.CategoryID, &tagsJSON,
			&hit.Status, &hit.Stock, &hit.CreatedAt, &hit.Rank,
		); err != nil {
			return nil, errx.Wrap(err, "scanning search result row", errx.TypeInternal)
		}
		_ = tenantIDStr // tenant_id read but not exposed in ProductHit
		_ = json.Unmarshal([]byte(imagesJSON), &hit.Images)
		_ = json.Unmarshal([]byte(tagsJSON), &hit.Tags)
		if hit.Images == nil {
			hit.Images = []string{}
		}
		if hit.Tags == nil {
			hit.Tags = []string{}
		}
		hits = append(hits, hit)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating search results", errx.TypeInternal)
	}

	if hits == nil {
		hits = []search.ProductHit{}
	}

	totalPages := total / q.PageSize
	if total%q.PageSize > 0 {
		totalPages++
	}

	return &search.SearchResult{
		Products:   hits,
		Total:      total,
		Page:       q.Page,
		PageSize:   q.PageSize,
		TotalPages: totalPages,
		Query:      q.Query,
	}, nil
}

// Suggest returns autocomplete suggestions matching the given prefix.
func (r *PostgresRepo) Suggest(ctx context.Context, tenantID kernel.TenantID, prefix string, limit int) ([]search.SearchSuggestion, error) {
	const q = `
		SELECT name, COUNT(*) AS cnt
		FROM products
		WHERE tenant_id = $1
		  AND status = 'active'
		  AND name ILIKE $2
		GROUP BY name
		ORDER BY cnt DESC, name ASC
		LIMIT $3`

	rows, err := r.db.QueryContext(ctx, q, string(tenantID), prefix+"%", limit)
	if err != nil {
		return nil, errx.Wrap(err, "querying suggestions", errx.TypeInternal)
	}
	defer rows.Close()

	var suggestions []search.SearchSuggestion
	for rows.Next() {
		var s search.SearchSuggestion
		if err := rows.Scan(&s.Term, &s.Count); err != nil {
			return nil, errx.Wrap(err, "scanning suggestion row", errx.TypeInternal)
		}
		suggestions = append(suggestions, s)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating suggestions", errx.TypeInternal)
	}

	if suggestions == nil {
		suggestions = []search.SearchSuggestion{}
	}
	return suggestions, nil
}

// buildOrderBy returns the SQL ORDER BY clause based on the sort option.
func buildOrderBy(sortBy string, hasQuery bool) string {
	switch sortBy {
	case "price_asc":
		return "price_amount ASC"
	case "price_desc":
		return "price_amount DESC"
	case "name":
		return "name ASC"
	case "created_at":
		return "created_at DESC"
	case "relevance":
		if hasQuery {
			return "rank DESC, created_at DESC"
		}
		return "created_at DESC"
	default:
		if hasQuery {
			return "rank DESC, created_at DESC"
		}
		return "created_at DESC"
	}
}

// Ensure interface compliance.
var _ search.Repository = (*PostgresRepo)(nil)
