package bigcache

import (
	"testing"
	"tasks-service-demo/internal/models"
)

func TestBigCacheStore_Create(t *testing.T) {
	store := NewBigCacheStore()
	defer store.Close()

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

func TestBigCacheStore_GetByID(t *testing.T) {
	store := NewBigCacheStore()
	defer store.Close()

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

func TestBigCacheStore_GetByID_NotFound(t *testing.T) {
	store := NewBigCacheStore()
	defer store.Close()

	_, err := store.GetByID(999)
	if err == nil {
		t.Fatal("Expected error for non-existent task")
	}
}

func TestBigCacheStore_Update(t *testing.T) {
	store := NewBigCacheStore()
	defer store.Close()

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

func TestBigCacheStore_Update_NotFound(t *testing.T) {
	store := NewBigCacheStore()
	defer store.Close()

	updatedTask := &models.Task{
		Name:   "Updated Task",
		Status: 1,
	}
	err := store.Update(999, updatedTask)
	if err == nil {
		t.Fatal("Expected error for non-existent task")
	}
}

func TestBigCacheStore_Delete(t *testing.T) {
	store := NewBigCacheStore()
	defer store.Close()

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

func TestBigCacheStore_Delete_NotFound(t *testing.T) {
	store := NewBigCacheStore()
	defer store.Close()

	err := store.Delete(999)
	if err == nil {
		t.Fatal("Expected error for non-existent task")
	}
}

func TestBigCacheStore_Stats(t *testing.T) {
	store := NewBigCacheStore()
	defer store.Close()

	// Create some tasks
	for i := 0; i < 10; i++ {
		task := &models.Task{
			Name:   "Test Task",
			Status: 0,
		}
		store.Create(task)
	}

	stats := store.Stats()
	if stats.Hits+stats.Misses == 0 {
		t.Log("Stats might not be enabled or no operations performed yet")
	}

	length := store.Len()
	if length != 10 {
		t.Errorf("Expected 10 entries, got %d", length)
	}
}