package product

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrNotFound      = errx.New("product not found", errx.TypeNotFound)
	ErrDuplicateSKU  = errx.New("product with this SKU already exists", errx.TypeConflict)
	ErrDuplicateSlug = errx.New("product with this slug already exists", errx.TypeConflict)
	ErrOutOfStock    = errx.New("product is out of stock", errx.TypeBusiness)
	ErrInvalidPrice  = errx.New("product price must be positive", errx.TypeValidation)

	// Variant & option errors.
	ErrOptionNotFound      = errx.New("product option not found", errx.TypeNotFound)
	ErrVariantNotFound     = errx.New("product variant not found", errx.TypeNotFound)
	ErrDuplicateVariantSKU = errx.New("variant with this SKU already exists", errx.TypeConflict)
)
