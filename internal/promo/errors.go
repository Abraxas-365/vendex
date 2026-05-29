package promo

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

var (
	ErrPromoNotFound    = errx.New("PROMO_NOT_FOUND", "promo not found", http.StatusNotFound)
	ErrCodeAlreadyExists = errx.New("PROMO_CODE_EXISTS", "a promo with this code already exists", http.StatusConflict)
	ErrPromoExpired     = errx.New("PROMO_EXPIRED", "this promo code has expired", http.StatusUnprocessableEntity)
	ErrPromoNotStarted  = errx.New("PROMO_NOT_STARTED", "this promo code is not yet active", http.StatusUnprocessableEntity)
	ErrPromoInactive    = errx.New("PROMO_INACTIVE", "this promo code is inactive", http.StatusUnprocessableEntity)
	ErrPromoMaxUses     = errx.New("PROMO_MAX_USES", "this promo code has reached its usage limit", http.StatusUnprocessableEntity)
	ErrPromoMinOrder    = errx.New("PROMO_MIN_ORDER", "order total does not meet the minimum required for this promo", http.StatusUnprocessableEntity)
)
