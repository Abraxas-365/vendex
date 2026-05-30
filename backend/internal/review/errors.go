package review

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrNotFound          = errx.New("review not found", errx.TypeNotFound)
	ErrInvalidRating     = errx.New("rating must be between 1 and 5", errx.TypeValidation)
	ErrInvalidStatus     = errx.New("invalid review status", errx.TypeValidation)
	ErrAlreadyModerated  = errx.New("review has already been moderated", errx.TypeBusiness)
	ErrProductIDRequired = errx.New("product_id is required", errx.TypeValidation)
	ErrCustomerIDRequired = errx.New("customer_id is required", errx.TypeValidation)
)
