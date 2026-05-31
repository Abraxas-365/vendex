package subscription

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrNotFound          = errx.NotFound("subscription not found")
	ErrAlreadyActive     = errx.Conflict("subscription already active")
	ErrNotActive         = errx.Business("subscription is not active")
	ErrAlreadyPaused     = errx.Business("subscription is already paused")
	ErrAlreadyCancelled  = errx.Business("subscription is already cancelled")
	ErrInvalidInterval   = errx.Validation("invalid billing interval")
	ErrBillingNotFound   = errx.NotFound("billing record not found")
)
