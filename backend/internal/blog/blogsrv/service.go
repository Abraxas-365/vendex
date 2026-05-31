package blogsrv

import (
	"context"
	"time"

	"github.com/Abraxas-365/vendex/internal/blog"
	"github.com/Abraxas-365/vendex/internal/errx"
	"github.com/Abraxas-365/vendex/internal/eventbus"
	"github.com/Abraxas-365/vendex/internal/kernel"
	"github.com/google/uuid"
)

// Service handles blog business logic.
type Service struct {
	repo blog.Repository
	bus  eventbus.Bus
}

// New creates a new blog service.
func New(repo blog.Repository, bus eventbus.Bus) *Service {
	return &Service{repo: repo, bus: bus}
}

// ============================================================================
// Post operations
// ============================================================================

// CreatePost creates a new blog post in draft status.
func (s *Service) CreatePost(ctx context.Context, tenantID kernel.TenantID, in blog.CreatePostInput) (*blog.BlogPost, error) {
	if in.Title == "" {
		return nil, blog.ErrTitleRequired
	}
	if in.Content == "" {
		return nil, blog.ErrContentRequired
	}
	if in.Slug == "" {
		return nil, blog.ErrSlugRequired
	}

	tags := in.Tags
	if tags == nil {
		tags = []string{}
	}

	now := time.Now().UTC()
	post := &blog.BlogPost{
		ID:              kernel.BlogPostID(uuid.NewString()),
		TenantID:        tenantID,
		Title:           in.Title,
		Slug:            in.Slug,
		Excerpt:         in.Excerpt,
		Content:         in.Content,
		FeaturedImage:   in.FeaturedImage,
		AuthorID:        in.AuthorID,
		AuthorName:      in.AuthorName,
		Status:          blog.StatusDraft,
		Tags:            tags,
		MetaTitle:       in.MetaTitle,
		MetaDescription: in.MetaDescription,
		Categories:      []blog.BlogCategory{},
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.CreatePost(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// GetPostByID retrieves a blog post by ID, scoped to tenant.
func (s *Service) GetPostByID(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogPostID) (*blog.BlogPost, error) {
	return s.repo.GetPostByID(ctx, tenantID, id)
}

// GetPostBySlug retrieves a blog post by slug, scoped to tenant.
func (s *Service) GetPostBySlug(ctx context.Context, tenantID kernel.TenantID, slug string) (*blog.BlogPost, error) {
	return s.repo.GetPostBySlug(ctx, tenantID, slug)
}

// ListPosts returns a paginated list of blog posts matching the given filter.
func (s *Service) ListPosts(ctx context.Context, tenantID kernel.TenantID, filter blog.ListPostsFilter) (kernel.Paginated[blog.BlogPost], error) {
	return s.repo.ListPosts(ctx, tenantID, filter)
}

// UpdatePost applies partial updates to an existing blog post.
func (s *Service) UpdatePost(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogPostID, in blog.UpdatePostInput) (*blog.BlogPost, error) {
	post, err := s.repo.GetPostByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if in.Title != nil {
		if *in.Title == "" {
			return nil, blog.ErrTitleRequired
		}
		post.Title = *in.Title
	}
	if in.Slug != nil {
		if *in.Slug == "" {
			return nil, blog.ErrSlugRequired
		}
		post.Slug = *in.Slug
	}
	if in.Excerpt != nil {
		post.Excerpt = *in.Excerpt
	}
	if in.Content != nil {
		if *in.Content == "" {
			return nil, blog.ErrContentRequired
		}
		post.Content = *in.Content
	}
	if in.FeaturedImage != nil {
		post.FeaturedImage = *in.FeaturedImage
	}
	if in.AuthorID != nil {
		post.AuthorID = *in.AuthorID
	}
	if in.AuthorName != nil {
		post.AuthorName = *in.AuthorName
	}
	if in.Tags != nil {
		post.Tags = in.Tags
	}
	if in.MetaTitle != nil {
		post.MetaTitle = *in.MetaTitle
	}
	if in.MetaDescription != nil {
		post.MetaDescription = *in.MetaDescription
	}
	post.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

// DeletePost removes a blog post.
func (s *Service) DeletePost(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogPostID) error {
	_, err := s.repo.GetPostByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return s.repo.DeletePost(ctx, tenantID, id)
}

// PublishPost transitions a post to published status.
func (s *Service) PublishPost(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogPostID) (*blog.BlogPost, error) {
	post, err := s.repo.GetPostByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if post.Status == blog.StatusPublished {
		return nil, blog.ErrAlreadyPublished
	}

	now := time.Now().UTC()
	post.Status = blog.StatusPublished
	post.PublishedAt = &now
	post.UpdatedAt = now

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, errx.Wrap(err, "publishing post", errx.TypeInternal)
	}

	if evt, err := eventbus.NewEvent(eventbus.BlogPostPublished, tenantID, eventbus.BlogPostPayload{
		PostID: string(post.ID),
		Title:  post.Title,
		Slug:   post.Slug,
	}); err == nil {
		_ = s.bus.Publish(ctx, evt)
	}

	return post, nil
}

// ArchivePost transitions a post to archived status.
func (s *Service) ArchivePost(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogPostID) (*blog.BlogPost, error) {
	post, err := s.repo.GetPostByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if post.Status == blog.StatusArchived {
		return nil, blog.ErrAlreadyArchived
	}

	post.Status = blog.StatusArchived
	post.UpdatedAt = time.Now().UTC()

	if err := s.repo.UpdatePost(ctx, post); err != nil {
		return nil, errx.Wrap(err, "archiving post", errx.TypeInternal)
	}

	return post, nil
}

// ============================================================================
// Category operations
// ============================================================================

// CreateCategory creates a new blog category.
func (s *Service) CreateCategory(ctx context.Context, tenantID kernel.TenantID, in blog.CreateCategoryInput) (*blog.BlogCategory, error) {
	if in.Name == "" {
		return nil, errx.New("name is required", errx.TypeValidation)
	}
	if in.Slug == "" {
		return nil, errx.New("slug is required", errx.TypeValidation)
	}

	cat := &blog.BlogCategory{
		ID:          kernel.BlogCategoryID(uuid.NewString()),
		TenantID:    tenantID,
		Name:        in.Name,
		Slug:        in.Slug,
		Description: in.Description,
		ParentID:    in.ParentID,
		SortOrder:   in.SortOrder,
		CreatedAt:   time.Now().UTC(),
	}

	if err := s.repo.CreateCategory(ctx, cat); err != nil {
		return nil, err
	}

	return cat, nil
}

// ListCategories returns all categories for a tenant.
func (s *Service) ListCategories(ctx context.Context, tenantID kernel.TenantID) ([]blog.BlogCategory, error) {
	return s.repo.ListCategories(ctx, tenantID)
}

// UpdateCategory applies partial updates to an existing category.
func (s *Service) UpdateCategory(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogCategoryID, in blog.UpdateCategoryInput) (*blog.BlogCategory, error) {
	cat, err := s.repo.GetCategoryByID(ctx, tenantID, id)
	if err != nil {
		return nil, err
	}

	if in.Name != nil {
		cat.Name = *in.Name
	}
	if in.Slug != nil {
		cat.Slug = *in.Slug
	}
	if in.Description != nil {
		cat.Description = *in.Description
	}
	if in.ParentID != nil {
		cat.ParentID = in.ParentID
	}
	if in.SortOrder != nil {
		cat.SortOrder = *in.SortOrder
	}

	if err := s.repo.UpdateCategory(ctx, cat); err != nil {
		return nil, err
	}

	return cat, nil
}

// DeleteCategory removes a blog category.
func (s *Service) DeleteCategory(ctx context.Context, tenantID kernel.TenantID, id kernel.BlogCategoryID) error {
	_, err := s.repo.GetCategoryByID(ctx, tenantID, id)
	if err != nil {
		return err
	}
	return s.repo.DeleteCategory(ctx, tenantID, id)
}

// AddPostCategory associates a category with a post.
func (s *Service) AddPostCategory(ctx context.Context, tenantID kernel.TenantID, postID kernel.BlogPostID, categoryID kernel.BlogCategoryID) error {
	return s.repo.AddPostCategory(ctx, tenantID, postID, categoryID)
}

// RemovePostCategory removes a category association from a post.
func (s *Service) RemovePostCategory(ctx context.Context, tenantID kernel.TenantID, postID kernel.BlogPostID, categoryID kernel.BlogCategoryID) error {
	return s.repo.RemovePostCategory(ctx, tenantID, postID, categoryID)
}
