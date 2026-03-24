package domain

import "errors"

// Sentinel errors used across the service layer.
// Handlers map these to HTTP status codes via HTTPError().
var (
	ErrNotFound     = errors.New("not found")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrConflict     = errors.New("conflict")
	ErrValidation   = errors.New("validation error")
	ErrBadRequest   = errors.New("bad request")
	ErrInternal     = errors.New("internal server error")
)

// AppError wraps a sentinel error with a human-readable message and optional
// field-level validation details.
type AppError struct {
	Err      error
	Message  string
	Details  map[string]string // field → reason, used for validation errors
}

func (e *AppError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

// NewNotFound creates a not-found AppError.
func NewNotFound(msg string) *AppError {
	return &AppError{Err: ErrNotFound, Message: msg}
}

// NewConflict creates a conflict AppError.
func NewConflict(msg string) *AppError {
	return &AppError{Err: ErrConflict, Message: msg}
}

// NewValidation creates a validation AppError with field-level detail.
func NewValidation(details map[string]string) *AppError {
	return &AppError{Err: ErrValidation, Message: "validation error", Details: details}
}

// NewForbidden creates a forbidden AppError.
func NewForbidden(msg string) *AppError {
	return &AppError{Err: ErrForbidden, Message: msg}
}

// NewUnauthorized creates an unauthorized AppError.
func NewUnauthorized(msg string) *AppError {
	return &AppError{Err: ErrUnauthorized, Message: msg}
}

// NewBadRequest creates a bad-request AppError.
func NewBadRequest(msg string) *AppError {
	return &AppError{Err: ErrBadRequest, Message: msg}
}
