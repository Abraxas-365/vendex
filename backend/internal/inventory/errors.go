package inventory

import "github.com/Abraxas-365/vendex/internal/errx"

var (
	ErrWarehouseNotFound    = errx.New("warehouse not found", errx.TypeNotFound)
	ErrStockLevelNotFound   = errx.New("stock level not found", errx.TypeNotFound)
	ErrMovementNotFound     = errx.New("stock movement not found", errx.TypeNotFound)
	ErrWarehouseNameConflict = errx.New("warehouse with this name already exists", errx.TypeConflict)
	ErrInvalidQuantity      = errx.New("quantity must not result in negative stock", errx.TypeValidation)
	ErrInvalidMovementType  = errx.New("invalid movement type", errx.TypeValidation)
	ErrWarehouseInactive    = errx.New("warehouse is inactive", errx.TypeBusiness)
)
