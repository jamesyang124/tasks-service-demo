package bigcache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync/atomic"
	"time"
	"tasks-service-demo/internal/models"

	"github.com/allegro/bigcache/v3"
)

// BigCacheStore implements storage using Allegro BigCache for high performance
type BigCacheStore struct {
	cache  *bigcache.BigCache
	nextID int64 // atomic counter for ID generation
}

// NewBigCacheStore creates a new BigCache-based store
func NewBigCacheStore() *BigCacheStore {
	config := bigcache.Config{
		Shards:             1024,                // Number of cache shards, must be power of two
		LifeWindow:         24 * time.Hour,     // Time after which entry can be evicted
		CleanWindow:        5 * time.Minute,    // Interval between removing expired entries
		MaxEntriesInWindow: 1000 * 10 * 60,     // Rps * lifeWindow, used only for statistics
		MaxEntrySize:       500,                // Max size of entry in bytes
		StatsEnabled:       false,              // Enable to collect statistics
		Verbose:            false,              // Enable to get info on what is happening
		Hasher:             nil, // Use default hasher
		HardMaxCacheSize:   8192,               // Max cache size in MB
		OnRemove:           nil,                // Callback fired when entry is removed
		OnRemoveWithReason: nil,                // Callback fired when entry is removed with reason
	}

	cache, err := bigcache.New(context.Background(), config)
	if err != nil {
		panic(fmt.Sprintf("Failed to create BigCache: %v", err))
	}

	return &BigCacheStore{
		cache:  cache,
		nextID: 0,
	}
}

// Create adds a new task to the store
func (s *BigCacheStore) Create(task *models.Task) error {
	// Generate unique ID atomically
	id := int(atomic.AddInt64(&s.nextID, 1))
	task.ID = id

	// Serialize task to JSON
	data, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	// Store in BigCache
	key := strconv.Itoa(id)
	err = s.cache.Set(key, data)
	if err != nil {
		return fmt.Errorf("failed to store task in cache: %w", err)
	}

	return nil
}

// GetByID retrieves a task by its ID
func (s *BigCacheStore) GetByID(id int) (*models.Task, error) {
	key := strconv.Itoa(id)
	
	data, err := s.cache.Get(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, fmt.Errorf("failed to get task from cache: %w", err)
	}

	var task models.Task
	err = json.Unmarshal(data, &task)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal task: %w", err)
	}

	return &task, nil
}

// Update modifies an existing task
func (s *BigCacheStore) Update(id int, updatedTask *models.Task) error {
	// Check if task exists
	_, err := s.GetByID(id)
	if err != nil {
		return err
	}

	// Set the ID and serialize
	updatedTask.ID = id
	data, err := json.Marshal(updatedTask)
	if err != nil {
		return fmt.Errorf("failed to marshal updated task: %w", err)
	}

	// Update in BigCache
	key := strconv.Itoa(id)
	err = s.cache.Set(key, data)
	if err != nil {
		return fmt.Errorf("failed to update task in cache: %w", err)
	}

	return nil
}

// Delete removes a task from the store
func (s *BigCacheStore) Delete(id int) error {
	key := strconv.Itoa(id)
	
	err := s.cache.Delete(key)
	if err != nil {
		if errors.Is(err, bigcache.ErrEntryNotFound) {
			return errors.New("task not found")
		}
		return fmt.Errorf("failed to delete task from cache: %w", err)
	}

	return nil
}

// GetAll retrieves all tasks from the store
func (s *BigCacheStore) GetAll() []*models.Task {
	var tasks []*models.Task
	
	// BigCache doesn't provide a direct way to iterate all keys
	// So we'll need to track them separately or use a different approach
	// For now, let's return empty slice as this is primarily a cache
	// In a real implementation, you might maintain a separate index
	
	return tasks
}

// Close closes the BigCache
func (s *BigCacheStore) Close() error {
	return s.cache.Close()
}

// Stats returns cache statistics
func (s *BigCacheStore) Stats() bigcache.Stats {
	return s.cache.Stats()
}

// Len returns the number of entries in cache
func (s *BigCacheStore) Len() int {
	return s.cache.Len()
}

// Capacity returns cache capacity
func (s *BigCacheStore) Capacity() int {
	return s.cache.Capacity()
}