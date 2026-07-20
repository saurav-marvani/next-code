package vcs

import "fmt"

// Error represents a VCS-related error
type Error struct {
	Code    string
	Message string
	Err     error
}

// NewError creates a new VCS error
func NewError(code, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Wrap wraps an error with VCS context
func (e *Error) Wrap(err error) *Error {
	e.Err = err
	return e
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Common VCS errors
var (
	ErrAuthentication   = NewError("auth_failed", "Authentication failed")
	ErrNotFound         = NewError("not_found", "Resource not found")
	ErrRateLimited      = NewError("rate_limited", "Rate limit exceeded")
	ErrUnauthorized     = NewError("unauthorized", "Unauthorized access")
	ErrInvalidInput     = NewError("invalid_input", "Invalid input")
	ErrServerError      = NewError("server_error", "Server error")
	ErrNotImplemented   = NewError("not_implemented", "Not implemented")
)
