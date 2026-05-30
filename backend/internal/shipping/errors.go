package shipping

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrZoneNotFound     = errx.New("shipping zone not found", errx.TypeNotFound)
	ErrRateNotFound     = errx.New("shipping rate not found", errx.TypeNotFound)
	ErrNoRatesAvailable = errx.New("no shipping rates available for this address", errx.TypeBusiness)
)
