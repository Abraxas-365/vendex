package notification

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrNotFound     = errx.NotFound("notification not found")
	ErrInvalidInput = errx.Validation("invalid notification input")
)
