package agentmemory

import "github.com/Abraxas-365/hada-commerce/internal/errx"

var (
	ErrNotFound        = errx.New("memory not found", errx.TypeNotFound)
	ErrInvalidCategory = errx.New("invalid memory category", errx.TypeValidation)
	ErrInvalidInput    = errx.New("invalid memory input", errx.TypeValidation)
)
