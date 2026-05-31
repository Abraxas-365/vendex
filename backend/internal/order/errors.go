package order

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrNotFound          = errx.New("order not found", errx.TypeNotFound)
	ErrEmptyOrder        = errx.New("order must contain at least one item", errx.TypeValidation)
	ErrInvalidTransition = errx.New("invalid order status transition", errx.TypeBusiness)
	ErrAlreadyCancelled  = errx.New("order is already cancelled", errx.TypeConflict)
)
