package recommendation

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrRuleNotFound  = errx.NotFound("recommendation rule not found")
	ErrInvalidInput  = errx.Validation("invalid recommendation input")
	ErrProductIDReq  = errx.Validation("product_id is required")
	ErrVisitorIDReq  = errx.Validation("visitor_id is required")
	ErrRuleNameReq   = errx.Validation("rule name is required")
	ErrRuleTypeReq   = errx.Validation("rule type is required")
)
