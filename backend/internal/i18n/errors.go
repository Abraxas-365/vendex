package i18n

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	// ErrNotFound is returned when a translation does not exist.
	ErrNotFound = errx.New("translation not found", errx.TypeNotFound)

	// ErrInvalidLocale is returned when an unsupported locale is provided.
	ErrInvalidLocale = errx.New("unsupported locale", errx.TypeValidation)

	// ErrInvalidEntityType is returned when an invalid entity type is provided.
	ErrInvalidEntityType = errx.New("invalid entity type", errx.TypeValidation)

	// ErrInvalidField is returned when an invalid translatable field name is provided.
	ErrInvalidField = errx.New("invalid field name", errx.TypeValidation)
)
