package webhook

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrNotFound         = errx.NotFound("webhook not found")
	ErrDeliveryNotFound = errx.NotFound("webhook delivery not found")
	ErrInvalidURL       = errx.Validation("invalid webhook URL")
	ErrNoEvents         = errx.Validation("at least one event type is required")
	ErrAlreadyActive    = errx.Conflict("webhook is already active")
	ErrAlreadyInactive  = errx.Conflict("webhook is already inactive")
)
