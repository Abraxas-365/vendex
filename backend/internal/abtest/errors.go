package abtest

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrNotFound              = errx.New("experiment not found", errx.TypeNotFound)
	ErrVariantNotFound       = errx.New("experiment variant not found", errx.TypeNotFound)
	ErrAssignmentNotFound    = errx.New("experiment assignment not found", errx.TypeNotFound)
	ErrAlreadyRunning        = errx.New("experiment is already running", errx.TypeBusiness)
	ErrNotRunning            = errx.New("experiment is not running", errx.TypeBusiness)
	ErrAlreadyCompleted      = errx.New("experiment is already completed", errx.TypeBusiness)
	ErrInsufficientVariants  = errx.New("experiment requires at least 2 variants to start", errx.TypeValidation)
	ErrCannotModifyRunning   = errx.New("cannot modify a running experiment", errx.TypeBusiness)
	ErrInvalidTrafficPercent = errx.New("traffic percent must be between 1 and 100", errx.TypeValidation)
	ErrAlreadyAssigned       = errx.New("visitor already assigned to this experiment", errx.TypeConflict)
	ErrInvalidWeight         = errx.New("variant weight must be greater than 0", errx.TypeValidation)
)
