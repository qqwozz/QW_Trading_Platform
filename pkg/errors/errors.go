// Package errors provides application-specific error types that carry HTTP
// status codes for consistent API error responses.
package errors

import (
	"fmt"
	"net/http"
)

// AppError represents an application error with an HTTP status code, a user-facing
// message, and an optional wrapped underlying error.
type AppError struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

// Error returns a formatted error string. If an underlying error is present,
// it is appended to the message.
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap returns the underlying error, enabling use with errors.Is and errors.As.
func (e *AppError) Unwrap() error {
	return e.Err
}

// NotFound creates a 404 Not Found error.
func NotFound(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

// NotFoundErr creates a 404 Not Found error wrapping an underlying error.
func NotFoundErr(message string, err error) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message, Err: err}
}

// BadRequest creates a 400 Bad Request error.
func BadRequest(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

// BadRequestErr creates a 400 Bad Request error wrapping an underlying error.
func BadRequestErr(message string, err error) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message, Err: err}
}

// Unauthorized creates a 401 Unauthorized error.
func Unauthorized(message string) *AppError {
	return &AppError{Code: http.StatusUnauthorized, Message: message}
}

// Forbidden creates a 403 Forbidden error.
func Forbidden(message string) *AppError {
	return &AppError{Code: http.StatusForbidden, Message: message}
}

// Conflict creates a 409 Conflict error.
func Conflict(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}

// Internal creates a 500 Internal Server Error.
func Internal(message string) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message}
}

// InternalErr creates a 500 Internal Server Error wrapping an underlying error.
func InternalErr(message string, err error) *AppError {
	return &AppError{Code: http.StatusInternalServerError, Message: message, Err: err}
}

// FromError converts a generic error into an AppError. If the error is already
// an *AppError, it is returned as-is; otherwise it becomes a 500 Internal error.
func FromError(err error) *AppError {
	if appErr, ok := err.(*AppError); ok {
		return appErr
	}
	return Internal(err.Error())
}
