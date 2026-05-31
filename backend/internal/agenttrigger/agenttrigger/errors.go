package agenttrigger

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrNotFound         = errx.NotFound("agent trigger not found")
	ErrInvalidEventType = errx.Validation("unsupported event type")
	ErrCooldownActive   = errx.Business("trigger is in cooldown period")
	ErrInvalidInput     = errx.Validation("invalid trigger input")
)
