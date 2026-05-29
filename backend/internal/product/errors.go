package product

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/kernel/errx"
)

var (
	ErrNotFound     = errx.New("PRODUCT_NOT_FOUND", "product not found", http.StatusNotFound)
	ErrDuplicateSKU = errx.New("PRODUCT_DUPLICATE_SKU", "product with this SKU already exists", http.StatusConflict)
	ErrOutOfStock   = errx.New("PRODUCT_OUT_OF_STOCK", "product is out of stock", http.StatusUnprocessableEntity)
	ErrInvalidPrice = errx.New("PRODUCT_INVALID_PRICE", "product price must be positive", http.StatusBadRequest)
)
