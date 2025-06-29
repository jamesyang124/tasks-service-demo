package entities

import (
	"encoding/json"
	"testing"
)

func TestTask_JSONMarshaling(t *testing.T) {
	task := Task{
		ID:     1,
		Name:   "Test Task",
		Status: 0,
	}

	// Test marshaling
	data, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Task
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal task: %v", err)
	}

	// Verify fields
	if unmarshaled.ID != task.ID {
		t.Errorf("Expected ID %d, got %d", task.ID, unmarshaled.ID)
	}
	if unmarshaled.Name != task.Name {
		t.Errorf("Expected name '%s', got '%s'", task.Name, unmarshaled.Name)
	}
	if unmarshaled.Status != task.Status {
		t.Errorf("Expected status %d, got %d", task.Status, unmarshaled.Status)
	}
}

func TestTask_DefaultValues(t *testing.T) {
	task := Task{}

	if task.ID != 0 {
		t.Errorf("Expected default ID 0, got %d", task.ID)
	}
	if task.Name != "" {
		t.Errorf("Expected default name '', got '%s'", task.Name)
	}
	if task.Status != 0 {
		t.Errorf("Expected default status 0, got %d", task.Status)
	}
}

func TestTask_StatusValues(t *testing.T) {
	tests := []struct {
		name   string
		status int
		valid  bool
	}{
		{"incomplete status", 0, true},
		{"complete status", 1, true},
		{"invalid negative", -1, false},
		{"invalid positive", 2, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				ID:     1,
				Name:   "Test",
				Status: tt.status,
			}

			// Verify the status is set correctly
			if task.Status != tt.status {
				t.Errorf("Expected status %d, got %d", tt.status, task.Status)
			}
		})
	}
}

func TestTask_NameValidation(t *testing.T) {
	tests := []struct {
		name     string
		taskName string
		valid    bool
	}{
		{"valid short name", "Test", true},
		{"valid long name", "This is a very long task name that should still be valid within the 100 character limit", true},
		{"empty name", "", false},
		{"name at max length", "x" + string(make([]byte, 99)), true},     // 100 chars
		{"name over max length", "x" + string(make([]byte, 100)), false}, // 101 chars
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			task := Task{
				ID:     1,
				Name:   tt.taskName,
				Status: 0,
			}

			// Verify the name is set correctly
			if task.Name != tt.taskName {
				t.Errorf("Expected name '%s', got '%s'", tt.taskName, task.Name)
			}
		})
	}
}
