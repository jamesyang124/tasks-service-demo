package errors

import (
	"testing"
)

func TestErrorCodes(t *testing.T) {
	tests := []struct {
		name     string
		code     int
		category string
		minRange int
		maxRange int
	}{
		{"TaskNotFound", ErrCodeTaskNotFound, "task", 1000, 1999},
		{"TaskInvalidInput", ErrCodeTaskInvalidInput, "task", 1000, 1999},
		{"TaskNameRequired", ErrCodeTaskNameRequired, "task", 1000, 1999},
		{"TaskNameTooLong", ErrCodeTaskNameTooLong, "task", 1000, 1999},
		{"TaskInvalidStatus", ErrCodeTaskInvalidStatus, "task", 1000, 1999},
		{"InvalidJSON", ErrCodeInvalidJSON, "request", 2000, 2999},
		{"InvalidID", ErrCodeInvalidID, "request", 2000, 2999},
		{"MissingFields", ErrCodeMissingFields, "request", 2000, 2999},
		{"InternalError", ErrCodeInternalError, "system", 5000, 5999},
		{"StorageError", ErrCodeStorageError, "system", 5000, 5999},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.code < tt.minRange || tt.code > tt.maxRange {
				t.Errorf("Error code %d is outside expected range %d-%d for %s category",
					tt.code, tt.minRange, tt.maxRange, tt.category)
			}
		})
	}
}

func TestErrorCodeUniqueness(t *testing.T) {
	codes := []int{
		ErrCodeTaskNotFound,
		ErrCodeTaskInvalidInput,
		ErrCodeTaskNameRequired,
		ErrCodeTaskNameTooLong,
		ErrCodeTaskInvalidStatus,
		ErrCodeInvalidJSON,
		ErrCodeInvalidID,
		ErrCodeMissingFields,
		ErrCodeInternalError,
		ErrCodeStorageError,
	}

	seen := make(map[int]bool)
	for _, code := range codes {
		if seen[code] {
			t.Errorf("Duplicate error code found: %d", code)
		}
		seen[code] = true
	}
}

func TestErrorResponse_Creation(t *testing.T) {
	code := ErrCodeTaskNotFound

	response := ErrorResponse{
		Code:    code,
		Message: "Additional context",
	}

	if response.Code != code {
		t.Errorf("Expected code %d, got %d", code, response.Code)
	}
	if response.Message != "Additional context" {
		t.Errorf("Expected message 'Additional context', got '%s'", response.Message)
	}
}

func TestAppError_Error(t *testing.T) {
	appErr := &AppError{
		Code:    ErrCodeTaskInvalidInput,
		Message: "test error",
		Type:    "VALIDATION_ERROR",
	}

	if appErr.Error() != "test error" {
		t.Errorf("Expected error message 'test error', got '%s'", appErr.Error())
	}

	if appErr.Code != ErrCodeTaskInvalidInput {
		t.Errorf("Expected code %d, got %d", ErrCodeTaskInvalidInput, appErr.Code)
	}
}
