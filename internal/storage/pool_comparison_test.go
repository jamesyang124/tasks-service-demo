package storage

import (
	"testing"
	"tasks-service-demo/internal/models"
)

const BenchDatasetSize = 10000 // Smaller dataset for focused comparison

// Benchmark ChannelStore with sync.Pool
func BenchmarkChannelStore_WithPool_Create(b *testing.B) {
	store := NewChannelStore(4)
	defer store.Shutdown()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &models.Task{
				Name:   "Benchmark Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}

func BenchmarkChannelStore_WithPool_Read(b *testing.B) {
	store := NewChannelStore(4)
	defer store.Shutdown()
	
	// Pre-populate
	for i := 1; i <= BenchDatasetSize; i++ {
		task := &models.Task{
			Name:   "Benchmark Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % BenchDatasetSize) + 1
			store.GetByID(targetID)
			i++
		}
	})
}

func BenchmarkChannelStore_WithPool_Update(b *testing.B) {
	store := NewChannelStore(4)
	defer store.Shutdown()
	
	// Pre-populate
	for i := 1; i <= BenchDatasetSize; i++ {
		task := &models.Task{
			Name:   "Benchmark Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % BenchDatasetSize) + 1
			updatedTask := &models.Task{
				Name:   "Updated Task",
				Status: 1,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

// Benchmark ChannelStoreNoPool without sync.Pool
func BenchmarkChannelStore_NoPool_Create(b *testing.B) {
	store := NewChannelStoreNoPool(4)
	defer store.Shutdown()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &models.Task{
				Name:   "Benchmark Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}

func BenchmarkChannelStore_NoPool_Read(b *testing.B) {
	store := NewChannelStoreNoPool(4)
	defer store.Shutdown()
	
	// Pre-populate
	for i := 1; i <= BenchDatasetSize; i++ {
		task := &models.Task{
			Name:   "Benchmark Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % BenchDatasetSize) + 1
			store.GetByID(targetID)
			i++
		}
	})
}

func BenchmarkChannelStore_NoPool_Update(b *testing.B) {
	store := NewChannelStoreNoPool(4)
	defer store.Shutdown()
	
	// Pre-populate
	for i := 1; i <= BenchDatasetSize; i++ {
		task := &models.Task{
			Name:   "Benchmark Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % BenchDatasetSize) + 1
			updatedTask := &models.Task{
				Name:   "Updated Task",
				Status: 1,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}