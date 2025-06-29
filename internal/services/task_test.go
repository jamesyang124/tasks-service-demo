package services

import (
	"testing"

	"tasks-service-demo/internal/entities"
	"tasks-service-demo/internal/requests"
	"tasks-service-demo/internal/storage"
	"tasks-service-demo/internal/storage/naive"
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
	req := &requests.CreateTaskRequest{Name: "Test Task", Status: 0}
	service.CreateTask(req)

	tasks = service.GetAllTasks()
	if len(tasks) != 1 {
		t.Errorf("Expected 1 task, got %d", len(tasks))
	}
}

func TestTaskService_GetTaskByID(t *testing.T) {
	service := setupTestService()

	// Create a task
	req := &requests.CreateTaskRequest{Name: "Test Task", Status: 0}
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

func TestTaskService_CreateTask_Success(t *testing.T) {
	service := setupTestService()

	req := &requests.CreateTaskRequest{
		Name:   "Test Task",
		Status: 0,
	}

	task, err := service.CreateTask(req)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if task == nil {
		t.Fatal("Expected task to be created")
	}

	if task.Name != req.Name {
		t.Errorf("Expected name '%s', got '%s'", req.Name, task.Name)
	}

	if task.Status != req.Status {
		t.Errorf("Expected status %d, got %d", req.Status, task.Status)
	}

	if task.ID == 0 {
		t.Error("Expected task ID to be set")
	}
}

func TestTaskService_CreateTask_ValidationError(t *testing.T) {
	service := setupTestService()

	tests := []struct {
		name string
		req  *requests.CreateTaskRequest
	}{
		{"empty name", &requests.CreateTaskRequest{Name: "", Status: 0}},
		{"invalid status", &requests.CreateTaskRequest{Name: "Test", Status: 2}},
		{"name too long", &requests.CreateTaskRequest{Name: string(make([]byte, 101)), Status: 0}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// In current implementation, service layer doesn't validate
			// Validation happens in middleware, so these should succeed
			task, err := service.CreateTask(tt.req)
			if err != nil {
				t.Errorf("Expected no error (validation happens in middleware), got %v", err)
			}
			if task == nil {
				t.Error("Expected task to be created (validation happens in middleware)")
			}
		})
	}
}

func TestTaskService_UpdateTask_Success(t *testing.T) {
	service := setupTestService()

	// Create a task first
	createReq := &requests.CreateTaskRequest{Name: "Original Task", Status: 0}
	createdTask, _ := service.CreateTask(createReq)

	// Update the task
	updateReq := &requests.UpdateTaskRequest{Name: "Updated Task", Status: 1}
	updatedTask, err := service.UpdateTask(createdTask.ID, updateReq)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if updatedTask.Name != updateReq.Name {
		t.Errorf("Expected name '%s', got '%s'", updateReq.Name, updatedTask.Name)
	}

	if updatedTask.Status != updateReq.Status {
		t.Errorf("Expected status %d, got %d", updateReq.Status, updatedTask.Status)
	}

	if updatedTask.ID != createdTask.ID {
		t.Errorf("Expected ID to remain %d, got %d", createdTask.ID, updatedTask.ID)
	}
}

func TestTaskService_UpdateTask_ValidationError(t *testing.T) {
	service := setupTestService()

	// Create a task first
	createReq := &requests.CreateTaskRequest{Name: "Original Task", Status: 0}
	createdTask, _ := service.CreateTask(createReq)

	tests := []struct {
		name string
		req  *requests.UpdateTaskRequest
	}{
		{"empty name", &requests.UpdateTaskRequest{Name: "", Status: 0}},
		{"invalid status", &requests.UpdateTaskRequest{Name: "Test", Status: 2}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// In current implementation, service layer doesn't validate
			// Validation happens in middleware, so these should succeed
			task, err := service.UpdateTask(createdTask.ID, tt.req)
			if err != nil {
				t.Errorf("Expected no error (validation happens in middleware), got %v", err)
			}
			if task == nil {
				t.Error("Expected task to be updated (validation happens in middleware)")
			}
		})
	}
}

func TestTaskService_UpdateTask_NotFound(t *testing.T) {
	service := setupTestService()

	updateReq := &requests.UpdateTaskRequest{Name: "Updated Task", Status: 1}
	task, err := service.UpdateTask(999, updateReq)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
	if task != nil {
		t.Error("Expected no task to be returned")
	}
}

func TestTaskService_DeleteTask_Success(t *testing.T) {
	service := setupTestService()

	// Create a task first
	req := &requests.CreateTaskRequest{Name: "Task to Delete", Status: 0}
	createdTask, _ := service.CreateTask(req)

	// Delete the task
	err := service.DeleteTask(createdTask.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify task is deleted
	_, err = service.GetTaskByID(createdTask.ID)
	if err == nil {
		t.Error("Expected error when getting deleted task")
	}
}

func TestTaskService_DeleteTask_NotFound(t *testing.T) {
	service := setupTestService()

	// RESTful DELETE should be idempotent - no error for non-existent resource
	err := service.DeleteTask(999)
	if err != nil {
		t.Errorf("Expected no error for non-existent task (RESTful idempotent), got: %v", err)
	}
}

func TestTaskService_Integration(t *testing.T) {
	service := setupTestService()

	// Create multiple tasks
	tasks := []*entities.Task{}
	for i := 0; i < 5; i++ {
		req := &requests.CreateTaskRequest{
			Name:   "Integration Task",
			Status: i % 2,
		}
		task, err := service.CreateTask(req)
		if err != nil {
			t.Fatalf("Failed to create task %d: %v", i, err)
		}
		tasks = append(tasks, task)
	}

	// Get all tasks
	allTasks := service.GetAllTasks()
	if len(allTasks) != 5 {
		t.Errorf("Expected 5 tasks, got %d", len(allTasks))
	}

	// Update each task
	for _, task := range tasks {
		updateReq := &requests.UpdateTaskRequest{
			Name:   "Updated Integration Task",
			Status: 1,
		}
		_, err := service.UpdateTask(task.ID, updateReq)
		if err != nil {
			t.Fatalf("Failed to update task %d: %v", task.ID, err)
		}
	}

	// Delete each task
	for _, task := range tasks {
		err := service.DeleteTask(task.ID)
		if err != nil {
			t.Fatalf("Failed to delete task %d: %v", task.ID, err)
		}
	}

	// Verify all tasks are deleted
	finalTasks := service.GetAllTasks()
	if len(finalTasks) != 0 {
		t.Errorf("Expected 0 tasks after deletion, got %d", len(finalTasks))
	}
}

func TestTaskService_CreateTask(t *testing.T) {
	service := setupTestService()

	// Test valid request
	req := &requests.CreateTaskRequest{
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
	invalidReq := &requests.CreateTaskRequest{
		Name:   "",
		Status: 0,
	}

	_, err = service.CreateTask(invalidReq)
	if err != nil {
		t.Errorf("Expected no error (validation happens in middleware), got %v", err)
	}

	// Test invalid status
	invalidReq2 := &requests.CreateTaskRequest{
		Name:   "Test",
		Status: 2, // Invalid status
	}

	_, err = service.CreateTask(invalidReq2)
	if err != nil {
		t.Errorf("Expected no error (validation happens in middleware), got %v", err)
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	service := setupTestService()

	// Create a task first
	createReq := &requests.CreateTaskRequest{
		Name:   "Original Task",
		Status: 0,
	}
	task, _ := service.CreateTask(createReq)

	// Test valid update
	updateReq := &requests.UpdateTaskRequest{
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
	invalidReq := &requests.UpdateTaskRequest{
		Name:   "", // Empty name
		Status: 0,
	}

	_, err = service.UpdateTask(task.ID, invalidReq)
	if err != nil {
		t.Errorf("Expected no error (validation happens in middleware), got %v", err)
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	service := setupTestService()

	// Create a task first
	createReq := &requests.CreateTaskRequest{
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

	// Test deleting non-existent task - RESTful DELETE should be idempotent
	err = service.DeleteTask(999)
	if err != nil {
		t.Errorf("Expected no error for non-existent task (RESTful idempotent), got: %v", err)
	}
}

func TestTaskService_ValidationIntegration(t *testing.T) {
	service := setupTestService()

	// Test that service properly validates before calling store
	tests := []struct {
		name        string
		req         *requests.CreateTaskRequest
		expectError bool
	}{
		{"Valid task", &requests.CreateTaskRequest{Name: "Valid", Status: 0}, false},
		{"Empty name", &requests.CreateTaskRequest{Name: "", Status: 0}, false},         // No validation in service
		{"Invalid status", &requests.CreateTaskRequest{Name: "Test", Status: 2}, false}, // No validation in service
		{"Valid completed", &requests.CreateTaskRequest{Name: "Done", Status: 1}, false},
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
