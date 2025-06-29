package shard

import (
	"runtime"
	"sync/atomic"
	"tasks-service-demo/internal/entities"
	apperrors "tasks-service-demo/internal/errors"
)

// ShardStore distributes tasks across multiple shard units using optimized sharding
type ShardStore struct {
	shards    []*ShardUnit // Array of shard units for distributed storage
	numShards int          // Total number of shards
	nextID    int64        // Atomic counter for lock-free ID generation
	shardMask int          // Bitmask for power-of-2 optimization

	// Performance optimizations - Pre-allocated worker slice
	workers    []chan workerJob // One dedicated worker channel per shard
	workerDone []chan struct{}  // Synchronization channels for graceful shutdown
}

// workerJob represents a shard processing job for worker pool
type workerJob struct {
	shardIndex int
	shard      *ShardUnit
	results    chan<- shardResult
}

// shardResult represents the result from processing a shard
type shardResult struct {
	tasks []*entities.Task
	index int
}

// isPowerOfTwo checks if a number is a power of 2
func isPowerOfTwo(n int) bool {
	return n > 0 && (n&(n-1)) == 0
}

// nextPowerOfTwo returns the next power of 2 >= n
func nextPowerOfTwo(n int) int {
	if n <= 1 {
		return 1
	}
	// Handle case where n is already power of 2
	if isPowerOfTwo(n) {
		return n
	}

	// Find next power of 2
	power := 1
	for power < n {
		power <<= 1
	}
	return power
}

// NewShardStore creates a new shard store with specified number of shards
// Optimized for power-of-2 shard counts for better CPU cache performance
func NewShardStore(numShards int) *ShardStore {
	if numShards <= 0 {
		// Default to CPU cores × 2, minimum 4, maximum 64
		numShards = runtime.NumCPU() * 2
		if numShards < 4 {
			numShards = 4
		}
		if numShards > 64 {
			numShards = 64
		}
	}

	// Round up to next power of 2 for bitwise optimization
	numShards = nextPowerOfTwo(numShards)
	shardMask := numShards - 1 // For bitwise AND operation

	// Pre-allocate shards with expected capacity for better memory layout
	shards := make([]*ShardUnit, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = NewShardUnit(64) // Pre-allocate map capacity to reduce rehashing
	}

	// Initialize pre-allocated worker channels (one per shard for optimal locality)
	workers := make([]chan workerJob, numShards)
	workerDone := make([]chan struct{}, numShards)

	for i := 0; i < numShards; i++ {
		workers[i] = make(chan workerJob, 2) // Small buffer for better throughput
		workerDone[i] = make(chan struct{})
	}

	store := &ShardStore{
		shards:     shards,
		numShards:  numShards,
		nextID:     0, // Start from 0 for atomic operations
		shardMask:  shardMask,
		workers:    workers,
		workerDone: workerDone,
	}

	// Start dedicated worker goroutines (one per shard for optimal locality)
	for i := 0; i < numShards; i++ {
		go store.dedicatedWorker(i)
	}

	return store
}

// getShardByID returns the shard index for a given ID using bitwise AND
// For power-of-2 shard counts, bitwise AND is faster than modulo
func (s *ShardStore) getShardByID(id int) int {
	return id & s.shardMask
}

// dedicatedWorker processes jobs for a specific shard (optimal CPU cache locality)
func (s *ShardStore) dedicatedWorker(workerID int) {
	workerChan := s.workers[workerID]

	for {
		select {
		case job := <-workerChan:
			// Process job for this dedicated worker's shard using ShardUnit
			tasks := job.shard.GetAll()

			// Send result
			job.results <- shardResult{tasks: tasks, index: job.shardIndex}

		case <-s.workerDone[workerID]:
			// Graceful shutdown signal
			return
		}
	}
}

// generateID generates a globally unique ID across all shards using atomic operations
func (s *ShardStore) generateID() int {
	return int(atomic.AddInt64(&s.nextID, 1))
}

// Create stores a task in the appropriate shard
func (s *ShardStore) Create(task *entities.Task) *apperrors.AppError {
	if task == nil {
		return apperrors.ErrTaskCannotBeNil
	}

	// Generate global ID
	task.ID = s.generateID()

	// Determine shard based on ID
	shardIndex := s.getShardByID(task.ID)

	// Access shard directly (no global mutex needed - array is immutable)
	shard := s.shards[shardIndex]

	// Store in the shard using ShardUnit API
	shard.Set(task.ID, task)

	return nil
}

// GetByID retrieves a task by ID from the appropriate shard
func (s *ShardStore) GetByID(id int) (*entities.Task, *apperrors.AppError) {
	shardIndex := s.getShardByID(id)

	// Access shard directly (no global mutex needed)
	shard := s.shards[shardIndex]

	// Use ShardUnit API for better encapsulation
	task, exists := shard.Get(id)
	if !exists {
		return nil, apperrors.ErrTaskNotFound
	}
	return task, nil
}

// GetAll retrieves all tasks from all shards using dedicated workers
func (s *ShardStore) GetAll() []*entities.Task {
	// Create result channel for this operation
	results := make(chan shardResult, s.numShards)

	// Submit jobs directly to dedicated workers (optimal shard-worker affinity)
	for i, shard := range s.shards {
		s.workers[i] <- workerJob{
			shardIndex: i,
			shard:      shard,
			results:    results,
		}
	}

	// Collect results from all shards
	var allTasks []*entities.Task
	for i := 0; i < s.numShards; i++ {
		result := <-results
		allTasks = append(allTasks, result.tasks...)
	}

	return allTasks
}

// Update modifies a task in the appropriate shard
func (s *ShardStore) Update(id int, updatedTask *entities.Task) *apperrors.AppError {
	shardIndex := s.getShardByID(id)

	// Access shard directly
	shard := s.shards[shardIndex]

	// Use ShardUnit API for better encapsulation
	updatedTask.ID = id
	if !shard.Update(id, updatedTask) {
		return apperrors.ErrTaskNotFound
	}
	return nil
}

// Delete removes a task from the appropriate shard
func (s *ShardStore) Delete(id int) *apperrors.AppError {
	shardIndex := s.getShardByID(id)

	// Access shard directly
	shard := s.shards[shardIndex]

	// Use ShardUnit API for better encapsulation
	if !shard.Delete(id) {
		return apperrors.ErrTaskNotFound
	}
	return nil
}

// Close gracefully shuts down all worker goroutines
func (s *ShardStore) Close() *apperrors.AppError {
	// Signal all workers to stop
	for i := 0; i < s.numShards; i++ {
		close(s.workerDone[i])
	}

	// Close worker channels to prevent new jobs
	for i := 0; i < s.numShards; i++ {
		close(s.workers[i])
	}

	return nil
}
