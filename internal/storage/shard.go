package storage

import (
	"runtime"
	"sync"
	"tasks-service-demo/internal/models"
)

// ShardStore distributes tasks across multiple memory stores using modulo hashing
type ShardStore struct {
	shards    []*MemoryStore
	numShards int
	mu        sync.RWMutex
	nextID    int
	idMutex   sync.Mutex
}

// NewShardStore creates a new shard store with specified number of shards
// For in-memory cache, optimal shard count is CPU cores × 2 for Go applications
// Hot shard issues are minimal since all shards are in local memory
func NewShardStore(numShards int) *ShardStore {
	if numShards <= 0 {
		// Default to CPU cores × 2, minimum 4, maximum 16
		numShards = runtime.NumCPU() * 2
		if numShards < 4 {
			numShards = 4
		}
		if numShards > 16 {
			numShards = 16
		}
	}

	shards := make([]*MemoryStore, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = NewMemoryStore()
	}

	return &ShardStore{
		shards:    shards,
		numShards: numShards,
		nextID:    1,
	}
}

// getShardByID returns the shard index for a given ID using simple modulo
// For sequential integer IDs, direct modulo is more efficient than hashing
func (s *ShardStore) getShardByID(id int) int {
	return id % s.numShards
}

// generateID generates a globally unique ID across all shards
func (s *ShardStore) generateID() int {
	s.idMutex.Lock()
	defer s.idMutex.Unlock()
	id := s.nextID
	s.nextID++
	return id
}

// Create stores a task in the appropriate shard
func (s *ShardStore) Create(task *models.Task) error {
	// Generate global ID
	task.ID = s.generateID()

	// Determine shard based on ID
	shardIndex := s.getShardByID(task.ID)

	s.mu.RLock()
	shard := s.shards[shardIndex]
	s.mu.RUnlock()

	// Store directly in the shard with our global ID
	shard.mu.Lock()
	shard.tasks[task.ID] = task
	shard.mu.Unlock()

	return nil
}

// GetByID retrieves a task by ID from the appropriate shard
func (s *ShardStore) GetByID(id int) (*models.Task, error) {
	shardIndex := s.getShardByID(id)

	s.mu.RLock()
	shard := s.shards[shardIndex]
	s.mu.RUnlock()

	// Use the shard's thread-safe GetByID method
	return shard.GetByID(id)
}

// GetAll retrieves all tasks from all shards
func (s *ShardStore) GetAll() []*models.Task {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var allTasks []*models.Task

	// Collect tasks from all shards
	for _, shard := range s.shards {
		shardTasks := shard.GetAll()
		allTasks = append(allTasks, shardTasks...)
	}

	return allTasks
}

// Update modifies a task in the appropriate shard
func (s *ShardStore) Update(id int, updatedTask *models.Task) error {
	shardIndex := s.getShardByID(id)

	s.mu.RLock()
	shard := s.shards[shardIndex]
	s.mu.RUnlock()

	// Update using the shard's thread-safe update method
	updatedTask.ID = id
	return shard.Update(id, updatedTask)
}

// Delete removes a task from the appropriate shard
func (s *ShardStore) Delete(id int) error {
	shardIndex := s.getShardByID(id)

	s.mu.RLock()
	shard := s.shards[shardIndex]
	s.mu.RUnlock()

	// Use the shard's thread-safe delete method
	return shard.Delete(id)
}

// GetShardStats returns statistics about shard distribution
func (s *ShardStore) GetShardStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := make(map[string]interface{})
	stats["numShards"] = s.numShards
	stats["totalTasks"] = len(s.GetAll())

	shardCounts := make([]int, s.numShards)
	for i, shard := range s.shards {
		shardCounts[i] = len(shard.GetAll())
	}
	stats["tasksPerShard"] = shardCounts

	return stats
}

// GetShard returns a specific shard (useful for testing/debugging)
func (s *ShardStore) GetShard(index int) *MemoryStore {
	if index < 0 || index >= s.numShards {
		return nil
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.shards[index]
}
