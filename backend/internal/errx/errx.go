package errx

import "fmt"

// ErrorType categorizes domain errors for consistent HTTP status mapping.
type ErrorType string

const (
	TypeNotFound     ErrorType = "NOT_FOUND"     // 404
	TypeConflict     ErrorType = "CONFLICT"       // 409
	TypeValidation   ErrorType = "VALIDATION"     // 400
	TypeAuthorization ErrorType = "AUTHORIZATION" // 401
	TypeBusiness     ErrorType = "BUSINESS"       // 422
	TypeInternal     ErrorType = "INTERNAL"       // 500
)

func (t ErrorType) HTTPStatus() int {
	switch t {
	case TypeNotFound:
		return 404
	case TypeConflict:
		return 409
	case TypeValidation:
		return 400
	case TypeAuthorization:
		return 401
	case TypeBusiness:
		return 422
	case TypeInternal:
		return 500
	default:
		return 500
	}
}

// Error is a structured domain error with message, type, and optional detail.
type Error struct {
	message string
	errType ErrorType
	detail  string
	wrapped error
}

// New creates a new domain error with the given message and type.
func New(message string, errType ErrorType) *Error {
	return &Error{message: message, errType: errType}
}

// Wrap wraps an underlying error with a domain error type and description.
func Wrap(err error, description string, errType ErrorType) *Error {
	return &Error{
		message: description,
		errType: errType,
		wrapped: err,
		detail:  err.Error(),
	}
}

func (e *Error) Error() string {
	if e.detail != "" {
		return fmt.Sprintf("%s: %s", e.message, e.detail)
	}
	return e.message
}

func (e *Error) Message() string   { return e.message }
func (e *Error) Type() ErrorType   { return e.errType }
func (e *Error) HTTPStatus() int   { return e.errType.HTTPStatus() }
func (e *Error) Unwrap() error     { return e.wrapped }

// Is checks if err is an *Error with the same message (same sentinel).
func Is(err error, target *Error) bool {
	if err == nil || target == nil {
		return false
	}
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.message == target.message && e.errType == target.errType
}

// HTTPStatus returns the HTTP status code for an error.
// Returns 500 if err is not an *Error.
func HTTPStatus(err error) int {
	if e, ok := err.(*Error); ok {
		return e.HTTPStatus()
	}
	return 500
}

// Message returns the human-readable message for an error.
// Returns "internal error" if err is not an *Error.
func Message(err error) string {
	if e, ok := err.(*Error); ok {
		return e.message
	}
	return "internal error"
}

// Type returns the ErrorType for an error.
// Returns TypeInternal if err is not an *Error.
func Type(err error) ErrorType {
	if e, ok := err.(*Error); ok {
		return e.errType
	}
	return TypeInternal
}
