package errors

// Package errors provides structured error types and helpers for the application.

// ErrorResponse represents a standardized API error response
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message,omitempty"`
}

// ToResponse creates an error response from AppError
func ToResponse(appErr *AppError) ErrorResponse {
	return ErrorResponse{
		Code:    appErr.Code,
		Message: appErr.Message,
	}
}

var (
	// ErrInternalErrorResponse is a pre-defined error response for internal server errors.
	ErrInternalErrorResponse = &AppError{
		Code:    ErrCodeInternalError,
		Message: ErrInternalError.Message,
	}
)
