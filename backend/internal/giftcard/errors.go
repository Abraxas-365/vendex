package giftcard

import (
	"github.com/Abraxas-365/vendex/internal/errx"
)

var (
	ErrNotFound            = errx.NotFound("gift card not found")
	ErrInsufficientBalance = errx.Business("insufficient gift card balance")
	ErrExpired             = errx.Business("gift card has expired")
	ErrInactive            = errx.Business("gift card is not active")
	ErrInvalidCode         = errx.Validation("invalid gift card code")
	ErrDuplicateCode       = errx.Conflict("gift card code already exists")
	ErrTransactionNotFound = errx.NotFound("gift card transaction not found")
)
