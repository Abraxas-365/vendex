package socialauth

import (
	"net/http"

	"github.com/Abraxas-365/hada-commerce/internal/errx"
)

var ErrRegistry = errx.NewRegistry("SOCIAL_AUTH")

var (
	CodeNotFound             = ErrRegistry.Register("NOT_FOUND", errx.TypeNotFound, http.StatusNotFound, "social account not found")
	CodeAlreadyLinked        = ErrRegistry.Register("ALREADY_LINKED", errx.TypeConflict, http.StatusConflict, "social account already linked to a customer")
	CodeProviderNotSupported = ErrRegistry.Register("PROVIDER_NOT_SUPPORTED", errx.TypeValidation, http.StatusBadRequest, "oauth provider is not supported")
	CodeInvalidOAuthCode     = ErrRegistry.Register("INVALID_OAUTH_CODE", errx.TypeValidation, http.StatusBadRequest, "invalid or expired oauth code")
)

func ErrNotFound() *errx.Error             { return ErrRegistry.New(CodeNotFound) }
func ErrAlreadyLinked() *errx.Error        { return ErrRegistry.New(CodeAlreadyLinked) }
func ErrProviderNotSupported() *errx.Error { return ErrRegistry.New(CodeProviderNotSupported) }
func ErrInvalidOAuthCode() *errx.Error     { return ErrRegistry.New(CodeInvalidOAuthCode) }
