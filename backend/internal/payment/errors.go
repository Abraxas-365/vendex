package payment

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrPaymentNotFound      = errx.New("payment not found", errx.TypeNotFound)
	ErrRefundNotFound       = errx.New("refund not found", errx.TypeNotFound)
	ErrAlreadyPaid          = errx.New("payment has already been completed", errx.TypeConflict)
	ErrInvalidAmount        = errx.New("payment amount must be greater than zero", errx.TypeValidation)
	ErrRefundExceedsPayment = errx.New("refund amount exceeds original payment amount", errx.TypeBusiness)
	ErrProviderError        = errx.New("payment provider returned an error", errx.TypeExternal)
)
