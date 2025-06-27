package models

import (
	"testing"
)

func TestValidateStruct_CreateTaskRequest_Success(t *testing.T) {
	tests := []struct {
		name string
		req  CreateTaskRequest
	}{
		{"valid task with status 0", CreateTaskRequest{Name: "Test Task", Status: 0}},
		{"valid task with status 1", CreateTaskRequest{Name: "Completed Task", Status: 1}},
		{"long name within limit", CreateTaskRequest{Name: "This is a longer task name that should still be valid", Status: 0}},
		{"single character name", CreateTaskRequest{Name: "T", Status: 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(&tt.req)
			if err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
		})
	}
}

func TestValidateStruct_CreateTaskRequest_ValidationErrors(t *testing.T) {
	tests := []struct {
		name          string
		req           CreateTaskRequest
		expectedField string
		expectedMsg   string
	}{
		{"empty name", CreateTaskRequest{Name: "", Status: 0}, "name", "name is required"},
		{"invalid status 2", CreateTaskRequest{Name: "Test", Status: 2}, "status", "status must be 0 (incomplete) or 1 (complete)"},
		{"invalid status -1", CreateTaskRequest{Name: "Test", Status: -1}, "status", "status must be 0 (incomplete) or 1 (complete)"},
		{"name too long", CreateTaskRequest{Name: "This is a very long task name that exceeds the maximum allowed length of 100 characters for testing validation system", Status: 0}, "name", "name must be at most 100 characters long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(&tt.req)
			if err == nil {
				t.Error("Expected validation error")
				return
			}

			validationErr, ok := err.(*ValidationError)
			if !ok {
				t.Errorf("Expected ValidationError, got %T", err)
				return
			}

			if validationErr.Field != tt.expectedField {
				t.Errorf("Expected field '%s', got '%s'", tt.expectedField, validationErr.Field)
			}

			if validationErr.Message != tt.expectedMsg {
				t.Errorf("Expected message '%s', got '%s'", tt.expectedMsg, validationErr.Message)
			}
		})
	}
}

func TestValidateStruct_UpdateTaskRequest(t *testing.T) {
	validReq := UpdateTaskRequest{Name: "Updated Task", Status: 1}
	err := ValidateStruct(&validReq)
	if err != nil {
		t.Errorf("Expected no validation error for valid request, got: %v", err)
	}

	invalidReq := UpdateTaskRequest{Name: "", Status: 0}
	err = ValidateStruct(&invalidReq)
	if err == nil {
		t.Error("Expected validation error for empty name")
	}
}

func TestValidateStruct_Task(t *testing.T) {
	validTask := Task{ID: 1, Name: "Test Task", Status: 0}
	err := ValidateStruct(&validTask)
	if err != nil {
		t.Errorf("Expected no validation error for valid task, got: %v", err)
	}

	invalidTask := Task{ID: 1, Name: "", Status: 0}
	err = ValidateStruct(&invalidTask)
	if err == nil {
		t.Error("Expected validation error for empty name")
	}
}

func TestValidationError_Error(t *testing.T) {
	err := &ValidationError{
		Field:   "name",
		Message: "name is required",
	}

	expected := "name is required"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}
