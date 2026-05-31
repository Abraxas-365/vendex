package blog

import (
	"context"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// Repository defines persistence operations for the blog domain.
type Repository interface {
	// Post operations
	CreatePost(ctx context.Context, post *BlogPost) error
	GetPostByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogPostID) (*BlogPost, error)
	GetPostBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*BlogPost, error)
	ListPosts(ctx context.Context, tenantID kernel.TenantID, filter ListPostsFilter) (kernel.Paginated[BlogPost], error)
	UpdatePost(ctx context.Context, post *BlogPost) error
	DeletePost(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogPostID) error

	// Category operations
	CreateCategory(ctx context.Context, cat *BlogCategory) error
	GetCategoryByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogCategoryID) (*BlogCategory, error)
	ListCategories(ctx context.Context, tenantID kernel.TenantID) ([]BlogCategory, error)
	UpdateCategory(ctx context.Context, cat *BlogCategory) error
	DeleteCategory(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogCategoryID) error

	// Post-category association
	AddPostCategory(ctx context.Context, tenantID kernel.TenantID, postID kernel.BlogPostID, categoryID kernel.BlogCategoryID) error
	RemovePostCategory(ctx context.Context, tenantID kernel.TenantID, postID kernel.BlogPostID, categoryID kernel.BlogCategoryID) error
	GetPostCategories(ctx context.Context, postID kernel.BlogPostID) ([]BlogCategory, error)
}
