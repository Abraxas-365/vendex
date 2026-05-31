package bloginfra

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Abraxas-365/hada-commerce/internal/blog"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// PostgresRepo implements blog.Repository using sqlx.
type PostgresRepo struct {
	db *sqlx.DB
}

// NewPostgresRepo creates a new PostgreSQL-backed blog repository.
func NewPostgresRepo(db *sqlx.DB) *PostgresRepo {
	return &PostgresRepo{db: db}
}

// ============================================================================
// Post operations
// ============================================================================

func (r *PostgresRepo) CreatePost(ctx context.Context, post *blog.BlogPost) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO blog_posts (
			id, tenant_id, title, slug, excerpt, content, featured_image,
			author_id, author_name, status, published_at, tags,
			meta_title, meta_description, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`,
		string(post.ID), string(post.TenantID),
		post.Title, post.Slug, post.Excerpt, post.Content, post.FeaturedImage,
		post.AuthorID, post.AuthorName, string(post.Status), post.PublishedAt,
		pq.StringArray(post.Tags),
		post.MetaTitle, post.MetaDescription, post.CreatedAt, post.UpdatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return blog.ErrPostSlugConflict
		}
		return errx.Wrap(err, "creating blog post", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetPostByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogPostID) (*blog.BlogPost, error) {
	post, err := r.scanPost(r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, title, slug, excerpt, content, featured_image,
		       author_id, author_name, status, published_at, tags,
		       meta_title, meta_description, created_at, updated_at
		FROM blog_posts
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	))
	if err != nil {
		return nil, err
	}

	cats, err := r.GetPostCategories(ctx, id)
	if err != nil {
		return nil, err
	}
	post.Categories = cats

	return post, nil
}

func (r *PostgresRepo) GetPostBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*blog.BlogPost, error) {
	post, err := r.scanPost(r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, title, slug, excerpt, content, featured_image,
		       author_id, author_name, status, published_at, tags,
		       meta_title, meta_description, created_at, updated_at
		FROM blog_posts
		WHERE slug = $1 AND tenant_id = $2`,
		slug, string(tenantID),
	))
	if err != nil {
		return nil, err
	}

	cats, err := r.GetPostCategories(ctx, post.ID)
	if err != nil {
		return nil, err
	}
	post.Categories = cats

	return post, nil
}

func (r *PostgresRepo) ListPosts(ctx context.Context, tenantID kernel.TenantID, filter blog.ListPostsFilter) (kernel.Paginated[blog.BlogPost], error) {
	var zero kernel.Paginated[blog.BlogPost]

	args := []any{string(tenantID)}
	where := "WHERE tenant_id = $1"
	paramIdx := 2

	if filter.Status != "" {
		where += fmt.Sprintf(" AND status = $%d", paramIdx)
		args = append(args, filter.Status)
		paramIdx++
	}

	if filter.Tag != "" {
		where += fmt.Sprintf(" AND $%d = ANY(tags)", paramIdx)
		args = append(args, filter.Tag)
		paramIdx++
	}

	if filter.CategoryID != "" {
		where += fmt.Sprintf(` AND id IN (
			SELECT post_id FROM blog_post_categories
			WHERE category_id = $%d AND tenant_id = $1
		)`, paramIdx)
		args = append(args, filter.CategoryID)
		paramIdx++
	}

	var total int
	countQ := "SELECT COUNT(*) FROM blog_posts " + where
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return zero, errx.Wrap(err, "counting blog posts", errx.TypeInternal)
	}

	pg := kernel.NewPaginationOptions(filter.Page, filter.PageSize)

	dataQ := fmt.Sprintf(`
		SELECT id, tenant_id, title, slug, excerpt, content, featured_image,
		       author_id, author_name, status, published_at, tags,
		       meta_title, meta_description, created_at, updated_at
		FROM blog_posts %s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d`, where, paramIdx, paramIdx+1)
	args = append(args, pg.Limit(), pg.Offset())

	rows, err := r.db.QueryContext(ctx, dataQ, args...)
	if err != nil {
		return zero, errx.Wrap(err, "listing blog posts", errx.TypeInternal)
	}
	defer rows.Close()

	var posts []blog.BlogPost
	for rows.Next() {
		post, err := r.scanPostRow(rows)
		if err != nil {
			return zero, err
		}
		cats, err := r.GetPostCategories(ctx, post.ID)
		if err != nil {
			return zero, err
		}
		post.Categories = cats
		posts = append(posts, *post)
	}
	if err := rows.Err(); err != nil {
		return zero, errx.Wrap(err, "iterating blog posts", errx.TypeInternal)
	}

	if posts == nil {
		posts = []blog.BlogPost{}
	}

	return kernel.NewPaginated(posts, pg.Page, pg.PageSize, total), nil
}

func (r *PostgresRepo) UpdatePost(ctx context.Context, post *blog.BlogPost) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE blog_posts SET
			title = $1, slug = $2, excerpt = $3, content = $4, featured_image = $5,
			author_id = $6, author_name = $7, status = $8, published_at = $9,
			tags = $10, meta_title = $11, meta_description = $12, updated_at = $13
		WHERE id = $14 AND tenant_id = $15`,
		post.Title, post.Slug, post.Excerpt, post.Content, post.FeaturedImage,
		post.AuthorID, post.AuthorName, string(post.Status), post.PublishedAt,
		pq.StringArray(post.Tags),
		post.MetaTitle, post.MetaDescription, post.UpdatedAt,
		string(post.ID), string(post.TenantID),
	)
	if err != nil {
		if isUniqueViolation(err) {
			return blog.ErrPostSlugConflict
		}
		return errx.Wrap(err, "updating blog post", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) DeletePost(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogPostID) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM blog_posts WHERE id = $1 AND tenant_id = $2",
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting blog post", errx.TypeInternal)
	}
	return nil
}

// ============================================================================
// Category operations
// ============================================================================

func (r *PostgresRepo) CreateCategory(ctx context.Context, cat *blog.BlogCategory) error {
	var parentID *string
	if cat.ParentID != nil {
		s := string(*cat.ParentID)
		parentID = &s
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO blog_categories (id, tenant_id, name, slug, description, parent_id, sort_order, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		string(cat.ID), string(cat.TenantID),
		cat.Name, cat.Slug, cat.Description, parentID, cat.SortOrder, cat.CreatedAt,
	)
	if err != nil {
		if isUniqueViolation(err) {
			return blog.ErrCategorySlugConflict
		}
		return errx.Wrap(err, "creating blog category", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetCategoryByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogCategoryID) (*blog.BlogCategory, error) {
	return r.scanCategory(r.db.QueryRowContext(ctx, `
		SELECT id, tenant_id, name, slug, description, parent_id, sort_order, created_at
		FROM blog_categories
		WHERE id = $1 AND tenant_id = $2`,
		string(id), string(tenantID),
	))
}

func (r *PostgresRepo) ListCategories(ctx context.Context, tenantID kernel.TenantID) ([]blog.BlogCategory, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, tenant_id, name, slug, description, parent_id, sort_order, created_at
		FROM blog_categories
		WHERE tenant_id = $1
		ORDER BY sort_order ASC, name ASC`,
		string(tenantID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "listing blog categories", errx.TypeInternal)
	}
	defer rows.Close()

	var cats []blog.BlogCategory
	for rows.Next() {
		cat, err := r.scanCategoryRow(rows)
		if err != nil {
			return nil, err
		}
		cats = append(cats, *cat)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating blog categories", errx.TypeInternal)
	}

	if cats == nil {
		cats = []blog.BlogCategory{}
	}
	return cats, nil
}

func (r *PostgresRepo) UpdateCategory(ctx context.Context, cat *blog.BlogCategory) error {
	var parentID *string
	if cat.ParentID != nil {
		s := string(*cat.ParentID)
		parentID = &s
	}

	_, err := r.db.ExecContext(ctx, `
		UPDATE blog_categories SET
			name = $1, slug = $2, description = $3, parent_id = $4, sort_order = $5
		WHERE id = $6 AND tenant_id = $7`,
		cat.Name, cat.Slug, cat.Description, parentID, cat.SortOrder,
		string(cat.ID), string(cat.TenantID),
	)
	if err != nil {
		if isUniqueViolation(err) {
			return blog.ErrCategorySlugConflict
		}
		return errx.Wrap(err, "updating blog category", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) DeleteCategory(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogCategoryID) error {
	_, err := r.db.ExecContext(ctx,
		"DELETE FROM blog_categories WHERE id = $1 AND tenant_id = $2",
		string(id), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "deleting blog category", errx.TypeInternal)
	}
	return nil
}

// ============================================================================
// Post-category association
// ============================================================================

func (r *PostgresRepo) AddPostCategory(ctx context.Context, tenantID kernel.TenantID, postID kernel.BlogPostID, categoryID kernel.BlogCategoryID) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO blog_post_categories (post_id, category_id, tenant_id)
		VALUES ($1, $2, $3)
		ON CONFLICT DO NOTHING`,
		string(postID), string(categoryID), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "adding post category", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) RemovePostCategory(ctx context.Context, tenantID kernel.TenantID, postID kernel.BlogPostID, categoryID kernel.BlogCategoryID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM blog_post_categories
		WHERE post_id = $1 AND category_id = $2 AND tenant_id = $3`,
		string(postID), string(categoryID), string(tenantID),
	)
	if err != nil {
		return errx.Wrap(err, "removing post category", errx.TypeInternal)
	}
	return nil
}

func (r *PostgresRepo) GetPostCategories(ctx context.Context, postID kernel.BlogPostID) ([]blog.BlogCategory, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.tenant_id, c.name, c.slug, c.description, c.parent_id, c.sort_order, c.created_at
		FROM blog_categories c
		JOIN blog_post_categories pc ON pc.category_id = c.id
		WHERE pc.post_id = $1
		ORDER BY c.sort_order ASC, c.name ASC`,
		string(postID),
	)
	if err != nil {
		return nil, errx.Wrap(err, "getting post categories", errx.TypeInternal)
	}
	defer rows.Close()

	var cats []blog.BlogCategory
	for rows.Next() {
		cat, err := r.scanCategoryRow(rows)
		if err != nil {
			return nil, err
		}
		cats = append(cats, *cat)
	}
	if err := rows.Err(); err != nil {
		return nil, errx.Wrap(err, "iterating post categories", errx.TypeInternal)
	}

	if cats == nil {
		cats = []blog.BlogCategory{}
	}
	return cats, nil
}

// ============================================================================
// Scan helpers
// ============================================================================

func (r *PostgresRepo) scanPost(row *sql.Row) (*blog.BlogPost, error) {
	var post blog.BlogPost
	var id, tenantID, status string
	var tags pq.StringArray

	err := row.Scan(
		&id, &tenantID, &post.Title, &post.Slug, &post.Excerpt, &post.Content, &post.FeaturedImage,
		&post.AuthorID, &post.AuthorName, &status, &post.PublishedAt, &tags,
		&post.MetaTitle, &post.MetaDescription, &post.CreatedAt, &post.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, blog.ErrPostNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning blog post", errx.TypeInternal)
	}

	post.ID = kernel.BlogPostID(id)
	post.TenantID = kernel.TenantID(tenantID)
	post.Status = blog.PostStatus(status)
	post.Tags = []string(tags)
	if post.Tags == nil {
		post.Tags = []string{}
	}

	return &post, nil
}

func (r *PostgresRepo) scanPostRow(rows *sql.Rows) (*blog.BlogPost, error) {
	var post blog.BlogPost
	var id, tenantID, status string
	var tags pq.StringArray

	err := rows.Scan(
		&id, &tenantID, &post.Title, &post.Slug, &post.Excerpt, &post.Content, &post.FeaturedImage,
		&post.AuthorID, &post.AuthorName, &status, &post.PublishedAt, &tags,
		&post.MetaTitle, &post.MetaDescription, &post.CreatedAt, &post.UpdatedAt,
	)
	if err != nil {
		return nil, errx.Wrap(err, "scanning blog post row", errx.TypeInternal)
	}

	post.ID = kernel.BlogPostID(id)
	post.TenantID = kernel.TenantID(tenantID)
	post.Status = blog.PostStatus(status)
	post.Tags = []string(tags)
	if post.Tags == nil {
		post.Tags = []string{}
	}

	return &post, nil
}

func (r *PostgresRepo) scanCategory(row *sql.Row) (*blog.BlogCategory, error) {
	var cat blog.BlogCategory
	var id, tenantID string
	var parentID *string

	err := row.Scan(
		&id, &tenantID, &cat.Name, &cat.Slug, &cat.Description,
		&parentID, &cat.SortOrder, &cat.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, blog.ErrCategoryNotFound
	}
	if err != nil {
		return nil, errx.Wrap(err, "scanning blog category", errx.TypeInternal)
	}

	cat.ID = kernel.BlogCategoryID(id)
	cat.TenantID = kernel.TenantID(tenantID)
	if parentID != nil {
		pid := kernel.BlogCategoryID(*parentID)
		cat.ParentID = &pid
	}

	return &cat, nil
}

func (r *PostgresRepo) scanCategoryRow(rows *sql.Rows) (*blog.BlogCategory, error) {
	var cat blog.BlogCategory
	var id, tenantID string
	var parentID *string

	err := rows.Scan(
		&id, &tenantID, &cat.Name, &cat.Slug, &cat.Description,
		&parentID, &cat.SortOrder, &cat.CreatedAt,
	)
	if err != nil {
		return nil, errx.Wrap(err, "scanning blog category row", errx.TypeInternal)
	}

	cat.ID = kernel.BlogCategoryID(id)
	cat.TenantID = kernel.TenantID(tenantID)
	if parentID != nil {
		pid := kernel.BlogCategoryID(*parentID)
		cat.ParentID = &pid
	}

	return &cat, nil
}

// isUniqueViolation checks if the error is a PostgreSQL unique constraint violation.
func isUniqueViolation(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}

// Ensure interface compliance.
var _ blog.Repository = (*PostgresRepo)(nil)
