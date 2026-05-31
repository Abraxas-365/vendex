package marketplace

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrVendorNotFound        = errx.New("vendor not found", errx.TypeNotFound)
	ErrVendorProductNotFound = errx.New("vendor product not found", errx.TypeNotFound)
	ErrVendorOrderNotFound   = errx.New("vendor order not found", errx.TypeNotFound)
	ErrConflict              = errx.New("vendor already exists", errx.TypeConflict)
	ErrInvalidInput          = errx.New("invalid input", errx.TypeValidation)

	// Preset errors
	ErrPresetNotFound     = errx.New("preset not found", errx.TypeNotFound)
	ErrPresetSlugTaken    = errx.New("preset slug already taken", errx.TypeConflict)
	ErrPresetNotInstalled = errx.New("preset not installed", errx.TypeNotFound)
)
