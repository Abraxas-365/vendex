package cartrecovery

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	// ErrNotFound is returned when a recovery email cannot be found.
	ErrNotFound = errx.New("recovery email not found", errx.TypeNotFound)

	// ErrAlreadyScheduled is returned when a recovery email already exists for the cart.
	ErrAlreadyScheduled = errx.New("recovery already scheduled for this cart", errx.TypeConflict)

	// ErrInvalidStatus is returned on an invalid status transition.
	ErrInvalidStatus = errx.New("invalid recovery status transition", errx.TypeValidation)
)
