package customer

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

var (
	ErrNotFound       = errx.New("CUSTOMER_NOT_FOUND", "customer not found", http.StatusNotFound)
	ErrDuplicateEmail = errx.New("CUSTOMER_DUPLICATE_EMAIL", "customer with this email already exists", http.StatusConflict)
	ErrInvalidEmail   = errx.New("CUSTOMER_INVALID_EMAIL", "invalid email address", http.StatusBadRequest)
)
