package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrNotFound      = &AppError{Code: "NOT_FOUND", Message: "Resource not found"}
	ErrUnauthorized  = &AppError{Code: "UNAUTHORIZED", Message: "Unauthorized"}
	ErrForbidden     = &AppError{Code: "FORBIDDEN", Message: "Forbidden"}
	ErrBadRequest    = &AppError{Code: "BAD_REQUEST", Message: "Bad request"}
	ErrInternalError = &AppError{Code: "INTERNAL_ERROR", Message: "Internal server error"}
)

func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
	}
}

func NewBadRequestError(message string) *AppError {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: message,
	}
}

func NewForbiddenError(message string) *AppError {
	return &AppError{
		Code:    "FORBIDDEN",
		Message: message,
	}
}



