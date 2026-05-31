package collection

import "github.com/Abraxas-365/hada-commerce/internal/errx"

// Domain errors for the collection bounded context.
var (
	ErrNotFound         = errx.New("collection not found", errx.TypeNotFound)
	ErrProductNotFound  = errx.New("collection product not found", errx.TypeNotFound)
	ErrDuplicateSlug    = errx.New("slug already taken", errx.TypeBusiness)
	ErrAlreadyInCollection = errx.New("product already in collection", errx.TypeBusiness)
	ErrInvalidType      = errx.New("invalid collection type: must be 'manual' or 'auto'", errx.TypeValidation)
	ErrNameRequired     = errx.New("name is required", errx.TypeValidation)
	ErrSlugRequired     = errx.New("slug is required", errx.TypeValidation)
)
