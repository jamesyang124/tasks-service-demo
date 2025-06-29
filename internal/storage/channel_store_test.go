package storage

import (
	"testing"
	"tasks-service-demo/internal/models"
)

func TestChannelStore_Create(t *testing.T) {
	store := NewChannelStore(4)
	defer store.Shutdown()

	task := &models.Task{
		Name:   "Test Task",
		Status: 0,
	}

	err := store.Create(task)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if task.ID == 0 {
		t.Fatalf("Expected task ID to be set, got %d", task.ID)
	}
}

func TestChannelStore_GetByID(t *testing.T) {
	store := NewChannelStore(4)
	defer store.Shutdown()

	// Create a task first
	task := &models.Task{
		Name:   "Test Task",
		Status: 0,
	}
	err := store.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Retrieve the task
	retrieved, err := store.GetByID(task.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.Name != task.Name {
		t.Errorf("Expected name %s, got %s", task.Name, retrieved.Name)
	}

	if retrieved.Status != task.Status {
		t.Errorf("Expected status %d, got %d", task.Status, retrieved.Status)
	}
}

func TestChannelStore_GetByID_NotFound(t *testing.T) {
	store := NewChannelStore(4)
	defer store.Shutdown()

	_, err := store.GetByID(999)
	if err == nil {
		t.Fatal("Expected error for non-existent task")
	}
}

func TestChannelStore_Update(t *testing.T) {
	store := NewChannelStore(4)
	defer store.Shutdown()

	// Create a task first
	task := &models.Task{
		Name:   "Original Task",
		Status: 0,
	}
	err := store.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Update the task
	updatedTask := &models.Task{
		Name:   "Updated Task",
		Status: 1,
	}
	err = store.Update(task.ID, updatedTask)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify the update
	retrieved, err := store.GetByID(task.ID)
	if err != nil {
		t.Fatalf("Failed to retrieve updated task: %v", err)
	}

	if retrieved.Name != updatedTask.Name {
		t.Errorf("Expected name %s, got %s", updatedTask.Name, retrieved.Name)
	}

	if retrieved.Status != updatedTask.Status {
		t.Errorf("Expected status %d, got %d", updatedTask.Status, retrieved.Status)
	}
}

func TestChannelStore_Update_NotFound(t *testing.T) {
	store := NewChannelStore(4)
	defer store.Shutdown()

	updatedTask := &models.Task{
		Name:   "Updated Task",
		Status: 1,
	}
	err := store.Update(999, updatedTask)
	if err == nil {
		t.Fatal("Expected error for non-existent task")
	}
}

func TestChannelStore_Delete(t *testing.T) {
	store := NewChannelStore(4)
	defer store.Shutdown()

	// Create a task first
	task := &models.Task{
		Name:   "Test Task",
		Status: 0,
	}
	err := store.Create(task)
	if err != nil {
		t.Fatalf("Failed to create task: %v", err)
	}

	// Delete the task
	err = store.Delete(task.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = store.GetByID(task.ID)
	if err == nil {
		t.Fatal("Expected error when retrieving deleted task")
	}
}

func TestChannelStore_Delete_NotFound(t *testing.T) {
	store := NewChannelStore(4)
	defer store.Shutdown()

	err := store.Delete(999)
	if err == nil {
		t.Fatal("Expected error for non-existent task")
	}
}

func TestChannelStore_GetAll(t *testing.T) {
	store := NewChannelStore(4)
	defer store.Shutdown()

	// Create multiple tasks
	task1 := &models.Task{Name: "Task 1", Status: 0}
	task2 := &models.Task{Name: "Task 2", Status: 1}
	task3 := &models.Task{Name: "Task 3", Status: 0}

	store.Create(task1)
	store.Create(task2)
	store.Create(task3)

	// Get all tasks
	tasks := store.GetAll()
	if len(tasks) != 3 {
		t.Fatalf("Expected 3 tasks, got %d", len(tasks))
	}

	// Verify task names exist (order might vary due to sharding)
	taskNames := make(map[string]bool)
	for _, task := range tasks {
		taskNames[task.Name] = true
	}

	if !taskNames["Task 1"] || !taskNames["Task 2"] || !taskNames["Task 3"] {
		t.Error("Not all tasks found in GetAll result")
	}
}

func TestChannelStore_ConcurrentOperations(t *testing.T) {
	store := NewChannelStore(8)
	defer store.Shutdown()

	const numGoroutines = 10
	const tasksPerGoroutine = 100

	// Create tasks concurrently
	done := make(chan bool, numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(workerID int) {
			for j := 0; j < tasksPerGoroutine; j++ {
				task := &models.Task{
					Name:   "Concurrent Task",
					Status: workerID % 2,
				}
				err := store.Create(task)
				if err != nil {
					t.Errorf("Failed to create task: %v", err)
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	// Verify total count
	tasks := store.GetAll()
	expectedCount := numGoroutines * tasksPerGoroutine
	if len(tasks) != expectedCount {
		t.Errorf("Expected %d tasks, got %d", expectedCount, len(tasks))
	}
}

func TestChannelStore_IDGeneration(t *testing.T) {
	store := NewChannelStore(4)
	defer store.Shutdown()

	// Create multiple tasks and verify unique IDs
	tasks := make([]*models.Task, 10)
	ids := make(map[int]bool)

	for i := 0; i < 10; i++ {
		task := &models.Task{
			Name:   "ID Test Task",
			Status: 0,
		}
		err := store.Create(task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		if ids[task.ID] {
			t.Fatalf("Duplicate ID generated: %d", task.ID)
		}
		ids[task.ID] = true
		tasks[i] = task
	}

	// Verify all IDs are unique and sequential
	for i := 1; i <= 10; i++ {
		if !ids[i] {
			t.Errorf("Missing ID: %d", i)
		}
	}
}