package shard

import (
	"sync"
	"tasks-service-demo/internal/models"
)

// ShardUnit is a lightweight, optimized storage unit for shard-based stores
// Removes unnecessary overhead from MemoryStore when used within sharded architecture
type ShardUnit struct {
	tasks map[int]*models.Task
	mu    sync.RWMutex
}

// NewShardUnit creates a new shard unit with pre-allocated capacity
func NewShardUnit(capacity int) *ShardUnit {
	if capacity <= 0 {
		capacity = 64 // Default capacity for better memory layout
	}
	
	return &ShardUnit{
		tasks: make(map[int]*models.Task, capacity),
	}
}

// Set stores a task with given ID (ID generation handled by parent ShardStore)
func (s *ShardUnit) Set(id int, task *models.Task) {
	s.mu.Lock()
	s.tasks[id] = task
	s.mu.Unlock()
}

// Get retrieves a task by ID
func (s *ShardUnit) Get(id int) (*models.Task, bool) {
	s.mu.RLock()
	task, exists := s.tasks[id]
	s.mu.RUnlock()
	return task, exists
}

// Update modifies an existing task
func (s *ShardUnit) Update(id int, task *models.Task) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.tasks[id]; !exists {
		return false
	}
	
	s.tasks[id] = task
	return true
}

// Delete removes a task by ID
func (s *ShardUnit) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if _, exists := s.tasks[id]; !exists {
		return false
	}
	
	delete(s.tasks, id)
	return true
}

// GetAll returns all tasks in this shard unit (for bulk operations)
func (s *ShardUnit) GetAll() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	tasks := make([]*models.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks
}

// Count returns the number of tasks in this shard unit
func (s *ShardUnit) Count() int {
	s.mu.RLock()
	count := len(s.tasks)
	s.mu.RUnlock()
	return count
}

// GetTasksUnsafe returns tasks map without locking (for use when parent already holds lock)
func (s *ShardUnit) GetTasksUnsafe() map[int]*models.Task {
	return s.tasks
}