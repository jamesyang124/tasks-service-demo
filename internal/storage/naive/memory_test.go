package naive

import (
	"tasks-service-demo/internal/entities"
	"testing"
)

func TestMemoryStore_Create(t *testing.T) {
	store := NewMemoryStore()

	task := &entities.Task{
		Name:   "Test Task",
		Status: 0,
	}

	err := store.Create(task)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if task.ID == 0 {
		t.Error("Expected task ID to be set")
	}

	if task.ID != 1 {
		t.Errorf("Expected first task ID to be 1, got %d", task.ID)
	}
}

func TestMemoryStore_GetByID(t *testing.T) {
	store := NewMemoryStore()

	// Create a task first
	task := &entities.Task{Name: "Test Task", Status: 0}
	store.Create(task)

	// Test getting existing task
	retrieved, err := store.GetByID(task.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved.Name != "Test Task" {
		t.Errorf("Expected name 'Test Task', got '%s'", retrieved.Name)
	}

	// Test getting non-existent task
	_, err = store.GetByID(999)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestMemoryStore_GetAll(t *testing.T) {
	store := NewMemoryStore()

	// Test empty store
	tasks := store.GetAll()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}

	// Add tasks
	task1 := &entities.Task{Name: "Task 1", Status: 0}
	task2 := &entities.Task{Name: "Task 2", Status: 1}
	store.Create(task1)
	store.Create(task2)

	tasks = store.GetAll()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(tasks))
	}
}

func TestMemoryStore_Update(t *testing.T) {
	store := NewMemoryStore()

	// Create a task first
	task := &entities.Task{Name: "Original", Status: 0}
	store.Create(task)

	// Update the task
	updatedTask := &entities.Task{Name: "Updated", Status: 1}
	err := store.Update(task.ID, updatedTask)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify update
	retrieved, _ := store.GetByID(task.ID)
	if retrieved.Name != "Updated" {
		t.Errorf("Expected name 'Updated', got '%s'", retrieved.Name)
	}
	if retrieved.Status != 1 {
		t.Errorf("Expected status 1, got %d", retrieved.Status)
	}

	// Test updating non-existent task
	err = store.Update(999, updatedTask)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestMemoryStore_Delete(t *testing.T) {
	store := NewMemoryStore()

	// Create a task first
	task := &entities.Task{Name: "To Delete", Status: 0}
	store.Create(task)

	// Delete the task
	err := store.Delete(task.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Verify deletion
	_, err = store.GetByID(task.ID)
	if err == nil {
		t.Error("Expected error for deleted task")
	}

	// Test deleting non-existent task
	err = store.Delete(999)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestMemoryStore_ConcurrentAccess(t *testing.T) {
	store := NewMemoryStore()

	// Test concurrent creates
	done := make(chan bool, 2)

	go func() {
		task := &entities.Task{Name: "Task 1", Status: 0}
		store.Create(task)
		done <- true
	}()

	go func() {
		task := &entities.Task{Name: "Task 2", Status: 1}
		store.Create(task)
		done <- true
	}()

	<-done
	<-done

	tasks := store.GetAll()
	if len(tasks) != 2 {
		t.Errorf("Expected 2 tasks after concurrent creates, got %d", len(tasks))
	}
}
