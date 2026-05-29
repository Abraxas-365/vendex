package promo

import (
	"github.com/Abraxas-365/hada-commerce/internal/errx"
)

var (
	ErrPromoNotFound    = errx.New("promo not found", errx.TypeNotFound)
	ErrCodeAlreadyExists = errx.New("a promo with this code already exists", errx.TypeConflict)
	ErrPromoExpired     = errx.New("this promo code has expired", errx.TypeBusiness)
	ErrPromoNotStarted  = errx.New("this promo code is not yet active", errx.TypeBusiness)
	ErrPromoInactive    = errx.New("this promo code is inactive", errx.TypeBusiness)
	ErrPromoMaxUses     = errx.New("this promo code has reached its usage limit", errx.TypeBusiness)
	ErrPromoMinOrder    = errx.New("order total does not meet the minimum required for this promo", errx.TypeBusiness)
)
