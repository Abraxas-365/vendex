package returns

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrNotFound         = errx.New("return request not found", errx.TypeNotFound)
	ErrInvalidStatus    = errx.New("invalid return status transition", errx.TypeBusiness)
	ErrInvalidInput     = errx.New("invalid return request input", errx.TypeValidation)
	ErrNoItems          = errx.New("return request must contain at least one item", errx.TypeValidation)
	ErrAlreadyClosed    = errx.New("return request is already closed", errx.TypeConflict)
	ErrAlreadyRejected  = errx.New("return request is already rejected", errx.TypeConflict)
)
