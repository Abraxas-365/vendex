package errx

import "fmt"

// Error is a structured domain error with code, message, and HTTP status.
type Error struct {
	code    string
	message string
	status  int
	detail  string
}

func New(code, message string, status int) *Error {
	return &Error{code: code, message: message, status: status}
}

func (e *Error) Error() string {
	if e.detail != "" {
		return fmt.Sprintf("%s: %s", e.message, e.detail)
	}
	return e.message
}

func (e *Error) Code() string    { return e.code }
func (e *Error) Message() string { return e.message }
func (e *Error) Status() int     { return e.status }

func Wrap(err *Error, detail string) *Error {
	return &Error{
		code:    err.code,
		message: err.message,
		status:  err.status,
		detail:  detail,
	}
}

func Is(err error, target *Error) bool {
	if err == nil || target == nil {
		return false
	}
	e, ok := err.(*Error)
	if !ok {
		return false
	}
	return e.code == target.code
}

func HTTPStatus(err error) int {
	if e, ok := err.(*Error); ok {
		return e.status
	}
	return 500
}

func Code(err error) string {
	if e, ok := err.(*Error); ok {
		return e.code
	}
	return "INTERNAL_ERROR"
}

func Message(err error) string {
	if e, ok := err.(*Error); ok {
		return e.message
	}
	return "internal error"
}
