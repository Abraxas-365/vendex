package blog

import (
	"time"

	"github.com/Abraxas-365/hada-commerce/internal/kernel"
)

// PostStatus represents the lifecycle of a blog post.
type PostStatus string

const (
	StatusDraft     PostStatus = "draft"
	StatusPublished PostStatus = "published"
	StatusArchived  PostStatus = "archived"
)

// BlogPost is the aggregate root for a blog post.
type BlogPost struct {
	ID              kernel.BlogPostID     `json:"id"`
	TenantID        kernel.TenantID       `json:"tenant_id"`
	Title           string                `json:"title"`
	Slug            string                `json:"slug"`
	Excerpt         string                `json:"excerpt"`
	Content         string                `json:"content"`
	FeaturedImage   string                `json:"featured_image"`
	AuthorID        string                `json:"author_id"`
	AuthorName      string                `json:"author_name"`
	Status          PostStatus            `json:"status"`
	PublishedAt     *time.Time            `json:"published_at,omitempty"`
	Tags            []string              `json:"tags"`
	Categories      []BlogCategory        `json:"categories,omitempty"`
	MetaTitle       string                `json:"meta_title"`
	MetaDescription string                `json:"meta_description"`
	CreatedAt       time.Time             `json:"created_at"`
	UpdatedAt       time.Time             `json:"updated_at"`
}

// BlogCategory is a category for organizing blog posts.
type BlogCategory struct {
	ID          kernel.BlogCategoryID  `json:"id"`
	TenantID    kernel.TenantID        `json:"tenant_id"`
	Name        string                 `json:"name"`
	Slug        string                 `json:"slug"`
	Description string                 `json:"description"`
	ParentID    *kernel.BlogCategoryID `json:"parent_id,omitempty"`
	SortOrder   int                    `json:"sort_order"`
	CreatedAt   time.Time              `json:"created_at"`
}

// CreatePostInput holds all data needed to create a blog post.
type CreatePostInput struct {
	Title           string
	Slug            string
	Excerpt         string
	Content         string
	FeaturedImage   string
	AuthorID        string
	AuthorName      string
	Tags            []string
	MetaTitle       string
	MetaDescription string
}

// UpdatePostInput holds the fields that can be updated on a blog post.
type UpdatePostInput struct {
	Title           *string
	Slug            *string
	Excerpt         *string
	Content         *string
	FeaturedImage   *string
	AuthorID        *string
	AuthorName      *string
	Tags            []string
	MetaTitle       *string
	MetaDescription *string
}

// ListPostsFilter holds filter criteria for listing blog posts.
type ListPostsFilter struct {
	Status     string
	CategoryID string
	Tag        string
	Page       int
	PageSize   int
}

// CreateCategoryInput holds all data needed to create a blog category.
type CreateCategoryInput struct {
	Name        string
	Slug        string
	Description string
	ParentID    *kernel.BlogCategoryID
	SortOrder   int
}

// UpdateCategoryInput holds the fields that can be updated on a blog category.
type UpdateCategoryInput struct {
	Name        *string
	Slug        *string
	Description *string
	ParentID    *kernel.BlogCategoryID
	SortOrder   *int
}
