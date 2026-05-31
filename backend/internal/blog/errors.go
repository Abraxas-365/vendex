package blog

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrPostNotFound         = errx.New("blog post not found", errx.TypeNotFound)
	ErrCategoryNotFound     = errx.New("blog category not found", errx.TypeNotFound)
	ErrPostSlugConflict     = errx.New("a post with this slug already exists", errx.TypeConflict)
	ErrCategorySlugConflict = errx.New("a category with this slug already exists", errx.TypeConflict)
	ErrAlreadyPublished     = errx.New("post is already published", errx.TypeBusiness)
	ErrAlreadyArchived      = errx.New("post is already archived", errx.TypeBusiness)
	ErrTitleRequired        = errx.New("title is required", errx.TypeValidation)
	ErrContentRequired      = errx.New("content is required", errx.TypeValidation)
	ErrSlugRequired         = errx.New("slug is required", errx.TypeValidation)
)
