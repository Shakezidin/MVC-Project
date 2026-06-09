package response

import (
	"errors"
	"net/http"
)

// AppError is a domain-aware error with HTTP mapping.
type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Details    map[string]string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// Error codes used across the application.
const (
	ErrCodeInternal          = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest        = "BAD_REQUEST"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeValidation        = "VALIDATION_ERROR"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeTooManyRequests   = "TOO_MANY_REQUESTS"
	ErrCodeRequestTimeout    = "REQUEST_TIMEOUT"
)

func NewInternalError(err error) *AppError {
	return &AppError{
		Code:       ErrCodeInternal,
		Message:    "something went wrong",
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeBadRequest,
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
	}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeUnauthorized,
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeForbidden,
		Message:    message,
		HTTPStatus: http.StatusForbidden,
	}
}

func NewNotFoundError(message string) *AppError {
	return &AppError{
		Code:       ErrCodeNotFound,
		Message:    message,
		HTTPStatus: http.StatusNotFound,
	}
}

func NewValidationError(details map[string]string) *AppError {
	return &AppError{
		Code:       ErrCodeValidation,
		Message:    "validation failed",
		HTTPStatus: http.StatusUnprocessableEntity,
		Details:    details,
	}
}

// AsAppError attempts to cast an error to AppError.
func AsAppError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return NewInternalError(err)
}
