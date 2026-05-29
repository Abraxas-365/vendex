package catalog

import (
	"github.com/Abraxas-365/hada-commerce/internal/errx"
)

var (
	// Category errors
	ErrCategoryNotFound    = errx.New("category not found", errx.TypeNotFound)
	ErrCategoryDuplicateSlug = errx.New("a category with this slug already exists", errx.TypeConflict)
	ErrCategoryInvalidInput  = errx.New("invalid category input", errx.TypeValidation)

	// Collection errors
	ErrCollectionNotFound = errx.New("collection not found", errx.TypeNotFound)
	ErrCollectionDupSlug  = errx.New("a collection with this slug already exists", errx.TypeConflict)
	ErrCollectionInvalid  = errx.New("invalid collection input", errx.TypeValidation)
)
