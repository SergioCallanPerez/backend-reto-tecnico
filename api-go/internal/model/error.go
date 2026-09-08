package model

import "net/http"

type ErrorCode string

const (
	ErrCodeValidation   ErrorCode = "VALIDATION_ERROR"
	ErrCodeUpstream     ErrorCode = "UPSTREAM_ERROR"
	ErrCodeInternal     ErrorCode = "INTERNAL_ERROR"
	ErrCodeUnauthorized ErrorCode = "UNAUTHORIZED"
)

type AppError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Detail  string    `json:"detail,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) HTTPStatus() int {
	switch e.Code {
	case ErrCodeValidation:
		return http.StatusBadRequest
	case ErrCodeUpstream:
		return http.StatusBadGateway
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}

func NewValidationError(message, detail string) *AppError {
	return &AppError{Code: ErrCodeValidation, Message: message, Detail: detail}
}

func NewUpstreamError(message, detail string) *AppError {
	return &AppError{Code: ErrCodeUpstream, Message: message, Detail: detail}
}

func NewInternalError(message string) *AppError {
	return &AppError{Code: ErrCodeInternal, Message: message}
}

func NewUnauthorizedError(message string) *AppError {
	return &AppError{Code: ErrCodeUnauthorized, Message: message}
}
