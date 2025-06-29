package services

import (
	"tasks-service-demo/internal/models"
	"tasks-service-demo/internal/storage"
	"tasks-service-demo/internal/storage/naive"
	"testing"
)

func setupTestService() *TaskService {
	storage.ResetStore()
	storage.InitStore(naive.NewMemoryStore())
	return NewTaskService()
}

func TestTaskService_GetAllTasks(t *testing.T) {
	service := setupTestService()

	// Test empty service
	tasks := service.GetAllTasks()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}

	// Add a task through service
	req := &models.CreateTaskRequest{Name: "Test Task", Status: 0}
	service.CreateTask(req)

	tasks = service.GetAllTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestTaskService_GetTaskByID(t *testing.T) {
	service := setupTestService()

	// Create a task
	req := &models.CreateTaskRequest{Name: "Test Task", Status: 0}
	task, _ := service.CreateTask(req)

	// Test getting existing task
	retrieved, err := service.GetTaskByID(task.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved.Name != "Test Task" {
		t.Errorf("Expected name 'Test Task', got '%s'", retrieved.Name)
	}

	// Test getting non-existent task
	_, err = service.GetTaskByID(999)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestTaskService_CreateTask(t *testing.T) {
	service := setupTestService()

	// Test valid request
	req := &models.CreateTaskRequest{
		Name:   "New Task",
		Status: 0,
	}

	task, err := service.CreateTask(req)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if task.Name != "New Task" {
		t.Errorf("Expected name 'New Task', got '%s'", task.Name)
	}

	if task.ID == 0 {
		t.Error("Expected task ID to be set")
	}

	// Test invalid request (empty name)
	invalidReq := &models.CreateTaskRequest{
		Name:   "",
		Status: 0,
	}

	_, err = service.CreateTask(invalidReq)
	if err == nil {
		t.Error("Expected validation error for empty name")
	}

	// Test invalid status
	invalidReq2 := &models.CreateTaskRequest{
		Name:   "Test",
		Status: 2, // Invalid status
	}

	_, err = service.CreateTask(invalidReq2)
	if err == nil {
		t.Error("Expected validation error for invalid status")
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	service := setupTestService()

	// Create a task first
	createReq := &models.CreateTaskRequest{
		Name:   "Original Task",
		Status: 0,
	}
	task, _ := service.CreateTask(createReq)

	// Test valid update
	updateReq := &models.UpdateTaskRequest{
		Name:   "Updated Task",
		Status: 1,
	}

	updatedTask, err := service.UpdateTask(task.ID, updateReq)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if updatedTask.Name != "Updated Task" {
		t.Errorf("Expected name 'Updated Task', got '%s'", updatedTask.Name)
	}

	if updatedTask.Status != 1 {
		t.Errorf("Expected status 1, got %d", updatedTask.Status)
	}

	// Test updating non-existent task
	_, err = service.UpdateTask(999, updateReq)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}

	// Test invalid update request
	invalidReq := &models.UpdateTaskRequest{
		Name:   "", // Empty name
		Status: 0,
	}

	_, err = service.UpdateTask(task.ID, invalidReq)
	if err == nil {
		t.Error("Expected validation error for empty name")
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	service := setupTestService()

	// Create a task first
	createReq := &models.CreateTaskRequest{
		Name:   "Task to Delete",
		Status: 0,
	}
	task, _ := service.CreateTask(createReq)

	// Delete the task
	err := service.DeleteTask(task.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = service.GetTaskByID(task.ID)
	if err == nil {
		t.Error("Expected error for deleted task")
	}

	// Test deleting non-existent task
	err = service.DeleteTask(999)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestTaskService_ValidationIntegration(t *testing.T) {
	service := setupTestService()

	// Test that service properly validates before calling store
	tests := []struct {
		name        string
		req         *models.CreateTaskRequest
		expectError bool
	}{
		{"Valid task", &models.CreateTaskRequest{Name: "Valid", Status: 0}, false},
		{"Empty name", &models.CreateTaskRequest{Name: "", Status: 0}, true},
		{"Invalid status", &models.CreateTaskRequest{Name: "Test", Status: 2}, true},
		{"Valid completed", &models.CreateTaskRequest{Name: "Done", Status: 1}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.CreateTask(tt.req)

			if tt.expectError && err == nil {
				t.Error("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
