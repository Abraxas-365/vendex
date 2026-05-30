package currency

import "github.com/Abraxas-365/hada-commerce/internal/errx"

// Domain errors for the currency bounded context.
var (
	ErrRateNotFound         = errx.New("exchange rate not found", errx.TypeNotFound)
	ErrUnsupportedCurrency  = errx.New("unsupported currency", errx.TypeValidation)
	ErrSameCurrency         = errx.New("source and target currency are the same", errx.TypeValidation)
	ErrInvalidRate          = errx.New("exchange rate must be positive", errx.TypeValidation)
	ErrDuplicateRate        = errx.New("exchange rate already exists for this currency pair", errx.TypeConflict)
)
