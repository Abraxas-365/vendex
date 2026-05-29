package errx

import (
	"fmt"
	"net/http"
)

// Type categorises domain errors for HTTP mapping.
type Type string

const (
	TypeValidation    Type = "validation"    // 400
	TypeAuthorization Type = "authorization" // 401
	TypeNotFound      Type = "not_found"     // 404
	TypeConflict      Type = "conflict"      // 409
	TypeBusiness      Type = "business"      // 422
	TypeInternal      Type = "internal"      // 500
)

// httpStatus maps an error Type to an HTTP status code.
func httpStatus(t Type) int {
	switch t {
	case TypeValidation:
		return http.StatusBadRequest
	case TypeAuthorization:
		return http.StatusUnauthorized
	case TypeNotFound:
		return http.StatusNotFound
	case TypeConflict:
		return http.StatusConflict
	case TypeBusiness:
		return http.StatusUnprocessableEntity
	case TypeInternal:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

// Error is a structured domain error with a human message, type, and optional detail.
type Error struct {
	message string
	errType Type
	detail  string
}

// New creates a new domain error with the given message and type.
func New(message string, errType Type) *Error {
	return &Error{message: message, errType: errType}
}

// Wrap wraps a Go error into an *Error with a descriptive message.
// Use this to convert infrastructure errors (DB, network) into domain errors.
func Wrap(err error, description string, errType Type) *Error {
	return &Error{
		message: description,
		errType: errType,
		detail:  err.Error(),
	}
}

// Error implements the error interface.
func (e *Error) Error() string {
	if e.detail != "" {
		return fmt.Sprintf("%s: %s", e.message, e.detail)
	}
	return e.message
}

// Message returns the human-readable error message.
func (e *Error) Message() string { return e.message }

// Type returns the error type.
func (e *Error) Type() Type { return e.errType }

// HTTPStatus returns the HTTP status code for this error.
func (e *Error) HTTPStatus() int { return httpStatus(e.errType) }

// --- package-level helpers for use in middleware / handlers ---

// HTTPStatus returns the HTTP status code for an error.
// If err is not an *Error, returns 500.
func HTTPStatus(err error) int {
	if e, ok := err.(*Error); ok {
		return e.HTTPStatus()
	}
	return http.StatusInternalServerError
}

// Message returns the human-readable message for an error.
// If err is not an *Error, returns "internal error".
func Message(err error) string {
	if e, ok := err.(*Error); ok {
		return e.Message()
	}
	return "internal error"
}

// Is checks whether err is an *Error of the given type.
func Is(err error, errType Type) bool {
	if e, ok := err.(*Error); ok {
		return e.errType == errType
	}
	return false
}

// IsNotFound returns true if err is a TypeNotFound error.
func IsNotFound(err error) bool { return Is(err, TypeNotFound) }

// IsConflict returns true if err is a TypeConflict error.
func IsConflict(err error) bool { return Is(err, TypeConflict) }

// IsValidation returns true if err is a TypeValidation error.
func IsValidation(err error) bool { return Is(err, TypeValidation) }
