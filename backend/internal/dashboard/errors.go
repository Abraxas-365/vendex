package dashboard

import "github.com/Abraxas-365/hada-commerce/internal/errx"

// ErrInvalidDateRange is returned when the From date is after the To date.
var ErrInvalidDateRange = errx.New("invalid date range: from must be before to", errx.TypeValidation)
