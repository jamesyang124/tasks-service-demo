package shard

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"tasks-service-demo/internal/models"

	"github.com/bytedance/gopkg/util/gopool"
)

// ShardStoreGopool uses ByteDance gopool for per-core worker optimization
type ShardStoreGopool struct {
	shards    []*ShardUnit
	numShards int
	nextID    int64 // atomic counter for lock-free ID generation
	shardMask int   // bitmask for power-of-2 optimization
	
	// Per-core worker pools using ByteDance gopool
	pools     []gopool.Pool // One pool per CPU core
	numCores  int
	coreMask  int // bitmask for core selection
}

// NewShardStoreGopool creates a new shard store with ByteDance gopool per-core workers
func NewShardStoreGopool(numShards int) *ShardStoreGopool {
	numCores := runtime.NumCPU()
	
	if numShards <= 0 {
		// Default to CPU cores × 2, minimum 4, maximum 64
		numShards = numCores * 2
		if numShards < 4 {
			numShards = 4
		}
		if numShards > 64 {
			numShards = 64
		}
	}
	
	// Round up to next power of 2 for bitwise optimization
	numShards = nextPowerOfTwo(numShards)
	shardMask := numShards - 1
	
	// Round cores to power of 2 for bitwise optimization
	numCores = nextPowerOfTwo(numCores)
	coreMask := numCores - 1

	// Pre-allocate shards with expected capacity
	shards := make([]*ShardUnit, numShards)
	for i := 0; i < numShards; i++ {
		shards[i] = NewShardUnit(64) // Pre-allocate map capacity to reduce rehashing
	}

	// Create per-core worker pools
	pools := make([]gopool.Pool, numCores)
	for i := 0; i < numCores; i++ {
		// Create pool with size optimized for M4 Pro (2 workers per core)
		pools[i] = gopool.NewPool("shard-core-"+string(rune(i+'0')), 2, gopool.NewConfig())
	}
	
	return &ShardStoreGopool{
		shards:    shards,
		numShards: numShards,
		nextID:    0,
		shardMask: shardMask,
		pools:     pools,
		numCores:  numCores,
		coreMask:  coreMask,
	}
}

// getCoreIndex returns the core index for a given shard using consistent hashing
func (s *ShardStoreGopool) getCoreIndex(shardIndex int) int {
	return shardIndex & s.coreMask
}

// getShardByID returns the shard index for a given ID using bitwise AND
func (s *ShardStoreGopool) getShardByID(id int) int {
	return id & s.shardMask
}

// generateID generates a globally unique ID across all shards using atomic operations
func (s *ShardStoreGopool) generateID() int {
	return int(atomic.AddInt64(&s.nextID, 1))
}

// Create stores a task in the appropriate shard
func (s *ShardStoreGopool) Create(task *models.Task) error {
	task.ID = s.generateID()
	shardIndex := s.getShardByID(task.ID)
	shard := s.shards[shardIndex]

	// Use ShardUnit API for better encapsulation
	shard.Set(task.ID, task)

	return nil
}

// GetByID retrieves a task by ID from the appropriate shard
func (s *ShardStoreGopool) GetByID(id int) (*models.Task, error) {
	shardIndex := s.getShardByID(id)
	shard := s.shards[shardIndex]

	// Use ShardUnit API for better encapsulation
	task, exists := shard.Get(id)
	if !exists {
		return nil, errors.New("task not found")
	}
	return task, nil
}

// GetAll retrieves all tasks from all shards using per-core gopool workers
func (s *ShardStoreGopool) GetAll() []*models.Task {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var allTasks []*models.Task

	// Process shards using per-core pools for optimal CPU utilization
	for i := 0; i < len(s.shards); i++ {
		wg.Add(1)
		
		// Capture variables for closure
		shardIndex := i
		shard := s.shards[i]
		
		// Select core pool using consistent hashing
		coreIndex := s.getCoreIndex(shardIndex)
		pool := s.pools[coreIndex]
		
		// Submit work to the core-specific pool
		pool.Go(func() {
			defer wg.Done()
			
			// Use ShardUnit API for better encapsulation
			tasks := shard.GetAll()
			
			// Collect results with minimal contention
			mu.Lock()
			allTasks = append(allTasks, tasks...)
			mu.Unlock()
		})
	}

	wg.Wait()
	return allTasks
}

// Update modifies a task in the appropriate shard
func (s *ShardStoreGopool) Update(id int, updatedTask *models.Task) error {
	shardIndex := s.getShardByID(id)
	shard := s.shards[shardIndex]

	// Use ShardUnit API for better encapsulation
	updatedTask.ID = id
	if !shard.Update(id, updatedTask) {
		return errors.New("task not found")
	}
	return nil
}

// Delete removes a task from the appropriate shard
func (s *ShardStoreGopool) Delete(id int) error {
	shardIndex := s.getShardByID(id)
	shard := s.shards[shardIndex]

	// Use ShardUnit API for better encapsulation
	if !shard.Delete(id) {
		return errors.New("task not found")
	}
	return nil
}

// Close gracefully shuts down all worker pools  
func (s *ShardStoreGopool) Close() error {
	// ByteDance gopool handles cleanup automatically
	// No explicit close needed for gopool.Pool
	return nil
}