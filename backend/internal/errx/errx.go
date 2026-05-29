// Package errx provides typed domain errors for hada-commerce.
// Errors carry a human-readable message and a type tag that maps to an HTTP status code.
package errx

import (
	"errors"
	"fmt"
	"net/http"
)

// Type categorises an error for HTTP status mapping.
type Type string

const (
	TypeNotFound     Type = "not_found"     // 404
	TypeConflict     Type = "conflict"      // 409
	TypeValidation   Type = "validation"    // 400
	TypeAuthorization Type = "authorization" // 401
	TypeBusiness     Type = "business"      // 422
	TypeInternal     Type = "internal"      // 500
)

// HTTPStatus maps an error type to its HTTP status code.
func (t Type) HTTPStatus() int {
	switch t {
	case TypeNotFound:
		return http.StatusNotFound
	case TypeConflict:
		return http.StatusConflict
	case TypeValidation:
		return http.StatusBadRequest
	case TypeAuthorization:
		return http.StatusUnauthorized
	case TypeBusiness:
		return http.StatusUnprocessableEntity
	case TypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Error is a domain error with a message and a type.
type Error struct {
	message string
	errType Type
	wrapped error
}

// New creates a new domain error with the given message and type.
func New(message string, t Type) *Error {
	return &Error{message: message, errType: t}
}

// Wrap creates a new domain error wrapping an underlying error.
func Wrap(err error, message string, t Type) *Error {
	return &Error{message: message, errType: t, wrapped: err}
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.wrapped != nil {
		return fmt.Sprintf("%s: %v", e.message, e.wrapped)
	}
	return e.message
}

// Unwrap implements errors.Unwrap for chain traversal.
func (e *Error) Unwrap() error {
	return e.wrapped
}

// Is enables errors.Is matching — two *Error values are equal if they point to the same instance.
func (e *Error) Is(target error) bool {
	var t *Error
	if errors.As(target, &t) {
		return e == t
	}
	return false
}

// Message returns the human-readable error message (without wrapped error detail).
func (e *Error) Message() string {
	return e.message
}

// Type returns the error type.
func (e *Error) Type() Type {
	return e.errType
}

// HTTPStatus returns the HTTP status code for this error.
func (e *Error) HTTPStatus() int {
	return e.errType.HTTPStatus()
}

// HTTPStatus extracts the HTTP status from any error.
// Falls back to 500 for non-domain errors.
func HTTPStatus(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.HTTPStatus()
	}
	return http.StatusInternalServerError
}

// Code returns a machine-readable code string derived from the error type.
func Code(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return string(e.errType)
	}
	return string(TypeInternal)
}

// Message extracts the human-readable message from any error.
func Message(err error) string {
	var e *Error
	if errors.As(err, &e) {
		return e.message
	}
	return err.Error()
}
