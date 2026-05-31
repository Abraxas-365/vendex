package approval

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	// ErrNotFound is returned when an approval request does not exist for the tenant.
	ErrNotFound = errx.New("approval request not found", errx.TypeNotFound)

	// ErrAlreadyReviewed is returned when trying to review an already-reviewed request.
	ErrAlreadyReviewed = errx.New("approval request already reviewed", errx.TypeBusiness)
)
