package bulkops

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	// ErrNotFound is returned when a bulk operation does not exist.
	ErrNotFound = errx.NotFound("bulk operation not found")

	// ErrInvalidInput is returned for malformed requests.
	ErrInvalidInput = errx.Validation("invalid bulk operation input")

	// ErrNoResourceIDs is returned when no resource IDs are provided.
	ErrNoResourceIDs = errx.Validation("no resource IDs provided")

	// ErrAlreadyCompleted is returned when an operation has already completed.
	ErrAlreadyCompleted = errx.Business("operation already completed")

	// ErrAlreadyProcessing is returned when an operation is already running.
	ErrAlreadyProcessing = errx.Business("operation is already processing")

	// ErrCannotCancel is returned when the operation is in a terminal state.
	ErrCannotCancel = errx.Business("operation cannot be cancelled in its current state")
)
