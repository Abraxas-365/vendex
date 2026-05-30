package storefrontinfra

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/Abraxas-365/hada-commerce/internal/storefront/renderer"
)

// NavMenuPostgresRepo implements renderer.NavMenuRepository using PostgreSQL.
type NavMenuPostgresRepo struct {
	db *sqlx.DB
}

// NewNavMenuPostgresRepo creates a new PostgreSQL-backed nav menu repository.
func NewNavMenuPostgresRepo(db *sqlx.DB) *NavMenuPostgresRepo {
	return &NavMenuPostgresRepo{db: db}
}

// ListByLocation returns navigation menu items for a tenant at the given location,
// ordered by position ascending.
func (r *NavMenuPostgresRepo) ListByLocation(ctx context.Context, tenantID kernel.TenantID, location renderer.NavLocation) ([]renderer.NavMenuItem, error) {
	const q = `
		SELECT id, label, url, position, COALESCE(parent_id::text, '') AS parent_id
		FROM navigation_menus
		WHERE tenant_id = $1
		  AND location  = $2
		ORDER BY position ASC, created_at ASC
	`

	rows, err := r.db.QueryContext(ctx, q, string(tenantID), string(location))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []renderer.NavMenuItem
	for rows.Next() {
		var item renderer.NavMenuItem
		if err := rows.Scan(&item.ID, &item.Label, &item.URL, &item.Position, &item.ParentID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}
