package audit

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	// ErrNotFound is returned when an audit entry cannot be located.
	ErrNotFound = errx.New("audit entry not found", errx.TypeNotFound)

	// ErrInvalidInput is returned for missing or malformed input.
	ErrInvalidInput = errx.New("invalid audit input", errx.TypeValidation)
)
