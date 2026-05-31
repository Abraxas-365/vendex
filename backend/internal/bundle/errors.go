package bundle

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	// ErrNotFound is returned when a bundle cannot be located.
	ErrNotFound = errx.New("bundle not found", errx.TypeNotFound)

	// ErrSlugTaken is returned when a slug already exists for the tenant.
	ErrSlugTaken = errx.New("bundle with this slug already exists", errx.TypeConflict)

	// ErrNoItems is returned when an operation requires at least one item but none exist.
	ErrNoItems = errx.New("bundle must have at least one item", errx.TypeBusiness)

	// ErrItemNotFound is returned when a bundle item cannot be located.
	ErrItemNotFound = errx.New("bundle item not found", errx.TypeNotFound)

	// ErrInvalidDiscount is returned when the discount value is out of range.
	ErrInvalidDiscount = errx.New("invalid discount value", errx.TypeValidation)
)
