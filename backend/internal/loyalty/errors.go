package loyalty

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrAccountNotFound    = errx.NotFound("loyalty account not found")
	ErrInsufficientPoints = errx.Business("insufficient points balance")
	ErrRewardNotFound     = errx.NotFound("loyalty reward not found")
	ErrRewardInactive     = errx.Business("loyalty reward is not active")
	ErrInvalidPoints      = errx.Validation("points must be greater than zero")
)
