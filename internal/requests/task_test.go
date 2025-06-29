package requests

import (
	"testing"
)

func TestCreateTaskRequest_Validation(t *testing.T) {
	tests := []struct {
		name          string
		request       CreateTaskRequest
		expectError   bool
		expectedField string
		expectedMsg   string
	}{
		{
			name:        "valid request",
			request:     CreateTaskRequest{Name: "Test Task", Status: 0},
			expectError: false,
		},
		{
			name:        "valid with status 1",
			request:     CreateTaskRequest{Name: "Completed Task", Status: 1},
			expectError: false,
		},
		{
			name:          "empty name",
			request:       CreateTaskRequest{Name: "", Status: 0},
			expectError:   true,
			expectedField: "name",
			expectedMsg:   "name is required",
		},
		{
			name:          "invalid status 2",
			request:       CreateTaskRequest{Name: "Test", Status: 2},
			expectError:   true,
			expectedField: "status",
			expectedMsg:   "status must be 0 (incomplete) or 1 (complete)",
		},
		{
			name:          "invalid status -1",
			request:       CreateTaskRequest{Name: "Test", Status: -1},
			expectError:   true,
			expectedField: "status",
			expectedMsg:   "status must be 0 (incomplete) or 1 (complete)",
		},
		{
			name:          "name too long",
			request:       CreateTaskRequest{Name: string(make([]byte, 101)), Status: 0},
			expectError:   true,
			expectedField: "name",
			expectedMsg:   "name must be at most 100 characters long",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.expectError {
				if err == nil {
					t.Error("Expected validation error")
					return
				}

				appErr := err

				if appErr.Message != tt.expectedMsg {
					t.Errorf("Expected message '%s', got '%s'", tt.expectedMsg, appErr.Message)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, got: %v", err)
				}
			}
		})
	}
}

func TestUpdateTaskRequest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		request     UpdateTaskRequest
		expectError bool
	}{
		{
			name:        "valid update request",
			request:     UpdateTaskRequest{Name: "Updated Task", Status: 1},
			expectError: false,
		},
		{
			name:        "empty name",
			request:     UpdateTaskRequest{Name: "", Status: 0},
			expectError: true,
		},
		{
			name:        "invalid status",
			request:     UpdateTaskRequest{Name: "Valid Name", Status: 3},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.expectError && err == nil {
				t.Error("Expected validation error")
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no validation error, got: %v", err)
			}
		})
	}
}

func TestValidatableInterface(t *testing.T) {
	// Test that both request types implement Validatable interface
	var createReq Validatable = CreateTaskRequest{Name: "Test", Status: 0}
	var updateReq Validatable = UpdateTaskRequest{Name: "Test", Status: 1}

	err := createReq.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid CreateTaskRequest, got: %v", err)
	}

	err = updateReq.Validate()
	if err != nil {
		t.Errorf("Expected no error for valid UpdateTaskRequest, got: %v", err)
	}
}
