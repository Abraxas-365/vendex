package catalog

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

var (
	ErrCategoryNotFound      = errx.New("CATEGORY_NOT_FOUND", "category not found", http.StatusNotFound)
	ErrCategoryDuplicateSlug = errx.New("CATEGORY_DUPLICATE_SLUG", "category with this slug already exists", http.StatusConflict)
	ErrCollectionNotFound    = errx.New("COLLECTION_NOT_FOUND", "collection not found", http.StatusNotFound)
	ErrCollectionDupSlug     = errx.New("COLLECTION_DUPLICATE_SLUG", "collection with this slug already exists", http.StatusConflict)
)
