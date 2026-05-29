package customer

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrNotFound       = errx.New("customer not found", errx.TypeNotFound)
	ErrDuplicateEmail = errx.New("customer with this email already exists", errx.TypeConflict)
	ErrInvalidEmail   = errx.New("invalid email address", errx.TypeValidation)
)
