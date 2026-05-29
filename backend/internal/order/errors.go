package order

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

var (
	ErrNotFound          = errx.New("ORDER_NOT_FOUND", "order not found", http.StatusNotFound)
	ErrEmptyOrder        = errx.New("ORDER_EMPTY", "order must contain at least one item", http.StatusBadRequest)
	ErrInvalidTransition = errx.New("ORDER_INVALID_TRANSITION", "invalid order status transition", http.StatusUnprocessableEntity)
	ErrAlreadyCancelled  = errx.New("ORDER_ALREADY_CANCELLED", "order is already cancelled", http.StatusConflict)
)
