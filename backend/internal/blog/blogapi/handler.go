package blogapi

import (
	"strconv"

	"github.com/Abraxas-365/hada-commerce/internal/blog"
	"github.com/Abraxas-365/hada-commerce/internal/blog/blogsrv"
	"github.com/Abraxas-365/hada-commerce/internal/errx"
	"github.com/Abraxas-365/hada-commerce/internal/kernel"
	"github.com/gofiber/fiber/v2"
)

// Handler exposes HTTP endpoints for the blog domain.
type Handler struct {
	svc *blogsrv.Service
}

// NewHandler creates a new blog API handler.
func NewHandler(svc *blogsrv.Service) *Handler {
	return &Handler{svc: svc}
}

// RegisterRoutes registers protected (admin) blog routes.
func (h *Handler) RegisterRoutes(router fiber.Router) {
	g := router.Group("/blog")

	// Post routes
	g.Post("/posts", h.CreatePost)
	g.Get("/posts", h.ListPosts)
	g.Get("/posts/:id", h.GetPostByID)
	g.Put("/posts/:id", h.UpdatePost)
	g.Delete("/posts/:id", h.DeletePost)
	g.Put("/posts/:id/publish", h.PublishPost)
	g.Put("/posts/:id/archive", h.ArchivePost)
	g.Post("/posts/:id/categories/:categoryId", h.AddPostCategory)
	g.Delete("/posts/:id/categories/:categoryId", h.RemovePostCategory)

	// Category routes
	g.Post("/categories", h.CreateCategory)
	g.Get("/categories", h.ListCategories)
	g.Put("/categories/:id", h.UpdateCategory)
	g.Delete("/categories/:id", h.DeleteCategory)
}

// RegisterPublicRoutes registers public blog routes (no auth required).
func (h *Handler) RegisterPublicRoutes(router fiber.Router) {
	g := router.Group("/blog")

	// Public post listing — published only
	g.Get("", h.ListPublicPosts)
	// Public post by slug
	g.Get("/:slug", h.GetPublicPostBySlug)
	// Public categories
	g.Get("/categories", h.ListPublicCategories)
}

// ============================================================================
// Post handlers
// ============================================================================

type createPostRequest struct {
	Title           string   `json:"title"`
	Slug            string   `json:"slug"`
	Excerpt         string   `json:"excerpt"`
	Content         string   `json:"content"`
	FeaturedImage   string   `json:"featured_image"`
	AuthorID        string   `json:"author_id"`
	AuthorName      string   `json:"author_name"`
	Tags            []string `json:"tags"`
	MetaTitle       string   `json:"meta_title"`
	MetaDescription string   `json:"meta_description"`
}

// CreatePost handles POST /blog/posts.
func (h *Handler) CreatePost(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createPostRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	post, err := h.svc.CreatePost(c.Context(), authCtx.TenantID, blog.CreatePostInput{
		Title:           req.Title,
		Slug:            req.Slug,
		Excerpt:         req.Excerpt,
		Content:         req.Content,
		FeaturedImage:   req.FeaturedImage,
		AuthorID:        req.AuthorID,
		AuthorName:      req.AuthorName,
		Tags:            req.Tags,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	})
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(post)
}

// GetPostByID handles GET /blog/posts/:id.
func (h *Handler) GetPostByID(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BlogPostID(c.Params("id"))

	post, err := h.svc.GetPostByID(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(post)
}

// ListPosts handles GET /blog/posts.
func (h *Handler) ListPosts(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	filter := blog.ListPostsFilter{
		Status:     c.Query("status"),
		CategoryID: c.Query("category_id"),
		Tag:        c.Query("tag"),
		Page:       page,
		PageSize:   pageSize,
	}

	result, err := h.svc.ListPosts(c.Context(), authCtx.TenantID, filter)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

type updatePostRequest struct {
	Title           *string  `json:"title"`
	Slug            *string  `json:"slug"`
	Excerpt         *string  `json:"excerpt"`
	Content         *string  `json:"content"`
	FeaturedImage   *string  `json:"featured_image"`
	AuthorID        *string  `json:"author_id"`
	AuthorName      *string  `json:"author_name"`
	Tags            []string `json:"tags"`
	MetaTitle       *string  `json:"meta_title"`
	MetaDescription *string  `json:"meta_description"`
}

// UpdatePost handles PUT /blog/posts/:id.
func (h *Handler) UpdatePost(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BlogPostID(c.Params("id"))

	var req updatePostRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	post, err := h.svc.UpdatePost(c.Context(), authCtx.TenantID, id, blog.UpdatePostInput{
		Title:           req.Title,
		Slug:            req.Slug,
		Excerpt:         req.Excerpt,
		Content:         req.Content,
		FeaturedImage:   req.FeaturedImage,
		AuthorID:        req.AuthorID,
		AuthorName:      req.AuthorName,
		Tags:            req.Tags,
		MetaTitle:       req.MetaTitle,
		MetaDescription: req.MetaDescription,
	})
	if err != nil {
		return err
	}

	return c.JSON(post)
}

// DeletePost handles DELETE /blog/posts/:id.
func (h *Handler) DeletePost(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BlogPostID(c.Params("id"))

	if err := h.svc.DeletePost(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// PublishPost handles PUT /blog/posts/:id/publish.
func (h *Handler) PublishPost(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BlogPostID(c.Params("id"))

	post, err := h.svc.PublishPost(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(post)
}

// ArchivePost handles PUT /blog/posts/:id/archive.
func (h *Handler) ArchivePost(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BlogPostID(c.Params("id"))

	post, err := h.svc.ArchivePost(c.Context(), authCtx.TenantID, id)
	if err != nil {
		return err
	}

	return c.JSON(post)
}

// AddPostCategory handles POST /blog/posts/:id/categories/:categoryId.
func (h *Handler) AddPostCategory(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	postID := kernel.BlogPostID(c.Params("id"))
	categoryID := kernel.BlogCategoryID(c.Params("categoryId"))

	if err := h.svc.AddPostCategory(c.Context(), authCtx.TenantID, postID, categoryID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RemovePostCategory handles DELETE /blog/posts/:id/categories/:categoryId.
func (h *Handler) RemovePostCategory(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	postID := kernel.BlogPostID(c.Params("id"))
	categoryID := kernel.BlogCategoryID(c.Params("categoryId"))

	if err := h.svc.RemovePostCategory(c.Context(), authCtx.TenantID, postID, categoryID); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ============================================================================
// Category handlers
// ============================================================================

type createCategoryRequest struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	ParentID    *string `json:"parent_id"`
	SortOrder   int     `json:"sort_order"`
}

// CreateCategory handles POST /blog/categories.
func (h *Handler) CreateCategory(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	var req createCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	in := blog.CreateCategoryInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
	if req.ParentID != nil {
		pid := kernel.BlogCategoryID(*req.ParentID)
		in.ParentID = &pid
	}

	cat, err := h.svc.CreateCategory(c.Context(), authCtx.TenantID, in)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).JSON(cat)
}

// ListCategories handles GET /blog/categories (admin).
func (h *Handler) ListCategories(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)

	cats, err := h.svc.ListCategories(c.Context(), authCtx.TenantID)
	if err != nil {
		return err
	}

	return c.JSON(cats)
}

type updateCategoryRequest struct {
	Name        *string `json:"name"`
	Slug        *string `json:"slug"`
	Description *string `json:"description"`
	ParentID    *string `json:"parent_id"`
	SortOrder   *int    `json:"sort_order"`
}

// UpdateCategory handles PUT /blog/categories/:id.
func (h *Handler) UpdateCategory(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BlogCategoryID(c.Params("id"))

	var req updateCategoryRequest
	if err := c.BodyParser(&req); err != nil {
		return errx.New("invalid request body", errx.TypeValidation)
	}

	in := blog.UpdateCategoryInput{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	}
	if req.ParentID != nil {
		pid := kernel.BlogCategoryID(*req.ParentID)
		in.ParentID = &pid
	}

	cat, err := h.svc.UpdateCategory(c.Context(), authCtx.TenantID, id, in)
	if err != nil {
		return err
	}

	return c.JSON(cat)
}

// DeleteCategory handles DELETE /blog/categories/:id.
func (h *Handler) DeleteCategory(c *fiber.Ctx) error {
	authCtx := c.Locals("auth").(*kernel.AuthContext)
	id := kernel.BlogCategoryID(c.Params("id"))

	if err := h.svc.DeleteCategory(c.Context(), authCtx.TenantID, id); err != nil {
		return err
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// ============================================================================
// Public handlers
// ============================================================================

// ListPublicPosts handles GET /blog — lists published posts only.
func (h *Handler) ListPublicPosts(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID is required", errx.TypeValidation)
	}

	page, _ := strconv.Atoi(c.Query("page"))
	pageSize, _ := strconv.Atoi(c.Query("page_size"))

	filter := blog.ListPostsFilter{
		Status:     string(blog.StatusPublished),
		CategoryID: c.Query("category_id"),
		Tag:        c.Query("tag"),
		Page:       page,
		PageSize:   pageSize,
	}

	result, err := h.svc.ListPosts(c.Context(), tenantID, filter)
	if err != nil {
		return err
	}

	return c.JSON(result)
}

// GetPublicPostBySlug handles GET /blog/:slug — returns a published post by slug.
func (h *Handler) GetPublicPostBySlug(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID is required", errx.TypeValidation)
	}

	slug := c.Params("slug")

	post, err := h.svc.GetPostBySlug(c.Context(), tenantID, slug)
	if err != nil {
		return err
	}

	if post.Status != blog.StatusPublished {
		return blog.ErrPostNotFound
	}

	return c.JSON(post)
}

// ListPublicCategories handles GET /blog/categories (public).
func (h *Handler) ListPublicCategories(c *fiber.Ctx) error {
	tenantID := kernel.TenantID(c.Get("X-Tenant-ID"))
	if tenantID == "" {
		return errx.New("tenant ID is required", errx.TypeValidation)
	}

	cats, err := h.svc.ListCategories(c.Context(), tenantID)
	if err != nil {
		return err
	}

	return c.JSON(cats)
}
