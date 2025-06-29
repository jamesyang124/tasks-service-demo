package errors

// Package errors provides structured error types and helpers for the application.

// AppError represents a structured application error with error code
type AppError struct {
	Code    int    `json:"code"`    // Error code for API responses
	Message string `json:"message"` // Human-readable error message
	Type    string `json:"type"`    // Error type for categorization
	Cause   error  `json:"-"`       // Original error, not serialized
}

// Error implements the error interface for AppError.
func (e *AppError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

// WithCause adds the underlying cause to the error and returns a new AppError.
func (e *AppError) WithCause(cause error) *AppError {
	return &AppError{
		Code:    e.Code,
		Message: e.Message,
		Type:    e.Type,
		Cause:   cause,
	}
}

// NewValidationError creates a new AppError of type VALIDATION_ERROR.
func NewValidationError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Type:    "VALIDATION_ERROR",
	}
}

// Pre-defined app errors for common scenarios.
var (
	// ErrTaskNotFound is returned when a requested task is not found
	ErrTaskNotFound = &AppError{
		Code:    ErrCodeTaskNotFound,
		Message: "Task not found",
		Type:    "NOT_FOUND",
	}
	// ErrTaskInvalidInput is returned when task input validation fails
	ErrTaskInvalidInput = &AppError{
		Code:    ErrCodeTaskInvalidInput,
		Message: "Invalid input provided",
		Type:    "VALIDATION_ERROR",
	}
	// ErrTaskCannotBeNil is returned when a nil task is provided
	ErrTaskCannotBeNil = &AppError{
		Code:    ErrCodeTaskInvalidInput,
		Message: "task cannot be nil",
		Type:    "VALIDATION_ERROR",
	}
	// ErrInternalError is returned for internal server errors
	ErrInternalError = &AppError{
		Code:    ErrCodeInternalError,
		Message: "Internal server error",
		Type:    "INTERNAL_ERROR",
	}
	// ErrStorageError is returned when storage operations fail
	ErrStorageError = &AppError{
		Code:    ErrCodeStorageError,
		Message: "storage operation error",
		Type:    "STORAGE_ERROR",
	}
)
