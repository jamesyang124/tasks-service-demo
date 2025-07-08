package xsync

import (
	"sync"
	"testing"
	"tasks-service-demo/internal/entities"
	apperrors "tasks-service-demo/internal/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestXSyncStore_Create(t *testing.T) {
	store := NewXSyncStore()
	
	task := &entities.Task{
		Name:   "Test Task",
		Status: 0,
	}
	
	err := store.Create(task)
	assert.Nil(t, err)
	assert.Equal(t, 1, task.ID)
	
	// Test second task gets incremented ID
	task2 := &entities.Task{
		Name:   "Test Task 2",
		Status: 1,
	}
	
	err = store.Create(task2)
	assert.Nil(t, err)
	assert.Equal(t, 2, task2.ID)
}

func TestXSyncStore_GetByID(t *testing.T) {
	store := NewXSyncStore()
	
	// Create a task
	task := &entities.Task{
		Name:   "Test Task",
		Status: 0,
	}
	store.Create(task)
	
	// Test successful retrieval
	retrieved, err := store.GetByID(task.ID)
	assert.Nil(t, err)
	assert.Equal(t, task.Name, retrieved.Name)
	assert.Equal(t, task.Status, retrieved.Status)
	
	// Test non-existent task
	_, err = store.GetByID(999)
	assert.NotNil(t, err)
	assert.Equal(t, apperrors.ErrTaskNotFound, err)
}

func TestXSyncStore_GetAll(t *testing.T) {
	store := NewXSyncStore()
	
	// Test empty store
	tasks := store.GetAll()
	assert.Empty(t, tasks)
	
	// Create multiple tasks
	task1 := &entities.Task{Name: "Task 1", Status: 0}
	task2 := &entities.Task{Name: "Task 2", Status: 1}
	
	store.Create(task1)
	store.Create(task2)
	
	tasks = store.GetAll()
	assert.Len(t, tasks, 2)
	
	// Verify tasks are in the result (order might vary)
	taskNames := []string{tasks[0].Name, tasks[1].Name}
	assert.Contains(t, taskNames, "Task 1")
	assert.Contains(t, taskNames, "Task 2")
}

func TestXSyncStore_Update(t *testing.T) {
	store := NewXSyncStore()
	
	// Create a task
	task := &entities.Task{
		Name:   "Original Task",
		Status: 0,
	}
	store.Create(task)
	
	// Update the task
	updatedTask := &entities.Task{
		Name:   "Updated Task",
		Status: 1,
	}
	
	err := store.Update(task.ID, updatedTask)
	assert.Nil(t, err)
	assert.Equal(t, task.ID, updatedTask.ID)
	
	// Verify update
	retrieved, err := store.GetByID(task.ID)
	assert.Nil(t, err)
	assert.Equal(t, "Updated Task", retrieved.Name)
	assert.Equal(t, 1, retrieved.Status)
	
	// Test updating non-existent task
	err = store.Update(999, updatedTask)
	assert.NotNil(t, err)
	assert.Equal(t, apperrors.ErrTaskNotFound, err)
}

func TestXSyncStore_Delete(t *testing.T) {
	store := NewXSyncStore()
	
	// Create a task
	task := &entities.Task{
		Name:   "Task to Delete",
		Status: 0,
	}
	store.Create(task)
	
	// Delete the task
	err := store.Delete(task.ID)
	assert.Nil(t, err)
	
	// Verify deletion
	_, err = store.GetByID(task.ID)
	assert.NotNil(t, err)
	assert.Equal(t, apperrors.ErrTaskNotFound, err)
	
	// Test deleting non-existent task
	err = store.Delete(999)
	assert.NotNil(t, err)
	assert.Equal(t, apperrors.ErrTaskNotFound, err)
}

func TestXSyncStore_ConcurrentOperations(t *testing.T) {
	store := NewXSyncStore()
	
	const numGoroutines = 100
	const numOperations = 10
	
	var wg sync.WaitGroup
	
	// Test concurrent creates
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				task := &entities.Task{
					Name:   "Concurrent Task",
					Status: 0,
				}
				err := store.Create(task)
				require.Nil(t, err)
			}
		}(i)
	}
	
	wg.Wait()
	
	// Verify all tasks were created
	tasks := store.GetAll()
	assert.Len(t, tasks, numGoroutines*numOperations)
	
	// Test concurrent reads
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < numOperations; j++ {
				allTasks := store.GetAll()
				assert.NotEmpty(t, allTasks)
			}
		}()
	}
	
	wg.Wait()
}

func TestXSyncStore_AtomicIDGeneration(t *testing.T) {
	store := NewXSyncStore()
	
	const numGoroutines = 50
	var wg sync.WaitGroup
	ids := make([]int, 0, numGoroutines)
	var mu sync.Mutex
	
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			task := &entities.Task{
				Name:   "Test Task",
				Status: 0,
			}
			store.Create(task)
			
			mu.Lock()
			ids = append(ids, task.ID)
			mu.Unlock()
		}()
	}
	
	wg.Wait()
	
	// Verify all IDs are unique
	idSet := make(map[int]bool)
	for _, id := range ids {
		assert.False(t, idSet[id], "Duplicate ID found: %d", id)
		idSet[id] = true
	}
	
	assert.Len(t, idSet, numGoroutines)
}