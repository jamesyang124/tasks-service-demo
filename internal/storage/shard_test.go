package storage

import (
	"fmt"
	"testing"
	"tasks-service-demo/internal/models"
)

func TestNewShardStore(t *testing.T) {
	// Test default shards (now CPU cores * 2, min 4, max 16)
	store := NewShardStore(0)
	if store.numShards < 4 || store.numShards > 16 {
		t.Errorf("Expected 4-16 default shards based on CPU cores, got %d", store.numShards)
	}
	
	// Test custom shards
	store = NewShardStore(8)
	if store.numShards != 8 {
		t.Errorf("Expected 8 shards, got %d", store.numShards)
	}
	
	if len(store.shards) != 8 {
		t.Errorf("Expected 8 shard instances, got %d", len(store.shards))
	}
}

func TestShardStore_Create(t *testing.T) {
	store := NewShardStore(4)
	
	task := &models.Task{
		Name:   "Shard Test Task",
		Status: 0,
	}
	
	err := store.Create(task)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	
	if task.ID == 0 {
		t.Error("Expected task ID to be set")
	}
	
	// Verify task is stored in correct shard
	retrieved, err := store.GetByID(task.ID)
	if err != nil {
		t.Errorf("Expected to retrieve task, got error: %v", err)
	}
	
	if retrieved.Name != "Shard Test Task" {
		t.Errorf("Expected name 'Shard Test Task', got '%s'", retrieved.Name)
	}
}

func TestShardStore_GetByID(t *testing.T) {
	store := NewShardStore(4)
	
	// Create multiple tasks
	tasks := []*models.Task{
		{Name: "Task 1", Status: 0},
		{Name: "Task 2", Status: 1},
		{Name: "Task 3", Status: 0},
	}
	
	for _, task := range tasks {
		store.Create(task)
	}
	
	// Test retrieving existing tasks
	for _, task := range tasks {
		retrieved, err := store.GetByID(task.ID)
		if err != nil {
			t.Errorf("Expected to find task %d, got error: %v", task.ID, err)
		}
		if retrieved.Name != task.Name {
			t.Errorf("Expected name '%s', got '%s'", task.Name, retrieved.Name)
		}
	}
	
	// Test non-existent task
	_, err := store.GetByID(999)
	if err == nil {
		t.Error("Expected error for non-existent task")
	}
}

func TestShardStore_GetAll(t *testing.T) {
	store := NewShardStore(4)
	
	// Test empty store
	tasks := store.GetAll()
	if len(tasks) != 0 {
		t.Errorf("Expected 0 tasks, got %d", len(tasks))
	}
	
	// Create tasks across different shards
	numTasks := 10
	for i := 0; i < numTasks; i++ {
		task := &models.Task{
			Name:   fmt.Sprintf("Task %d", i),
			Status: i % 2,
		}
		store.Create(task)
	}
	
	// Verify all tasks are returned
	allTasks := store.GetAll()
	if len(allTasks) != numTasks {
		t.Errorf("Expected %d tasks, got %d", numTasks, len(allTasks))
	}
}

func TestShardStore_Update(t *testing.T) {
	store := NewShardStore(4)
	
	// Create a task
	task := &models.Task{Name: "Original", Status: 0}
	store.Create(task)
	
	// Update the task
	updatedTask := &models.Task{Name: "Updated", Status: 1}
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

func TestShardStore_Delete(t *testing.T) {
	store := NewShardStore(4)
	
	// Create a task
	task := &models.Task{Name: "To Delete", Status: 0}
	store.Create(task)
	
	// Verify task exists
	_, err := store.GetByID(task.ID)
	if err != nil {
		t.Errorf("Expected task to exist before deletion")
	}
	
	// Delete the task
	err = store.Delete(task.ID)
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

func TestShardStore_ShardDistribution(t *testing.T) {
	store := NewShardStore(4)
	
	// Create many tasks to test distribution
	numTasks := 100
	for i := 0; i < numTasks; i++ {
		task := &models.Task{
			Name:   fmt.Sprintf("Task %d", i),
			Status: i % 2,
		}
		store.Create(task)
	}
	
	// Check shard statistics
	stats := store.GetShardStats()
	
	if stats["numShards"] != 4 {
		t.Errorf("Expected 4 shards, got %v", stats["numShards"])
	}
	
	if stats["totalTasks"] != numTasks {
		t.Errorf("Expected %d total tasks, got %v", numTasks, stats["totalTasks"])
	}
	
	// Verify tasks are distributed across shards (not all in one shard)
	tasksPerShard := stats["tasksPerShard"].([]int)
	nonEmptyShards := 0
	for _, count := range tasksPerShard {
		if count > 0 {
			nonEmptyShards++
		}
	}
	
	if nonEmptyShards < 2 {
		t.Errorf("Expected tasks to be distributed across multiple shards, got %d non-empty shards", nonEmptyShards)
	}
}

func TestShardStore_ConsistentHashing(t *testing.T) {
	store := NewShardStore(4)
	
	// Test that same ID always maps to same shard
	testID := 42
	shard1 := store.getShardByID(testID)
	shard2 := store.getShardByID(testID)
	
	if shard1 != shard2 {
		t.Errorf("Expected consistent shard mapping, got %d and %d", shard1, shard2)
	}
	
	// Test that different IDs can map to different shards
	differentShards := make(map[int]bool)
	for i := 1; i <= 20; i++ {
		shard := store.getShardByID(i)
		differentShards[shard] = true
	}
	
	if len(differentShards) < 2 {
		t.Error("Expected different IDs to map to different shards")
	}
}

func TestShardStore_ConcurrentAccess(t *testing.T) {
	store := NewShardStore(4)
	
	// Test concurrent creates
	done := make(chan bool, 10)
	
	for i := 0; i < 10; i++ {
		go func(id int) {
			task := &models.Task{
				Name:   fmt.Sprintf("Concurrent Task %d", id),
				Status: 0,
			}
			err := store.Create(task)
			if err != nil {
				t.Errorf("Concurrent create failed: %v", err)
			}
			done <- true
		}(i)
	}
	
	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
	
	// Verify all tasks were created
	tasks := store.GetAll()
	if len(tasks) != 10 {
		t.Errorf("Expected 10 tasks after concurrent creates, got %d", len(tasks))
	}
}

func TestShardStore_GetShard(t *testing.T) {
	store := NewShardStore(4)
	
	// Test valid shard index
	shard := store.GetShard(0)
	if shard == nil {
		t.Error("Expected valid shard, got nil")
	}
	
	// Test invalid shard indices
	if store.GetShard(-1) != nil {
		t.Error("Expected nil for negative index")
	}
	
	if store.GetShard(4) != nil {
		t.Error("Expected nil for out-of-bounds index")
	}
}