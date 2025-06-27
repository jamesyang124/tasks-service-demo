package storage

import (
	"testing"
	"tasks-service-demo/internal/models"
)

func BenchmarkMemoryStore_Create(b *testing.B) {
	store := NewMemoryStore()
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		task := &models.Task{
			Name:   "Benchmark Task",
			Status: 0,
		}
		store.Create(task)
	}
}

func BenchmarkMemoryStore_GetByID(b *testing.B) {
	store := NewMemoryStore()
	
	// Setup: Create 1000 tasks
	for i := 0; i < 1000; i++ {
		task := &models.Task{
			Name:   "Setup Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetByID((i % 1000) + 1)
	}
}

func BenchmarkMemoryStore_GetAll(b *testing.B) {
	store := NewMemoryStore()
	
	// Setup: Create 1000 tasks
	for i := 0; i < 1000; i++ {
		task := &models.Task{
			Name:   "Setup Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetAll()
	}
}

func BenchmarkMemoryStore_Update(b *testing.B) {
	store := NewMemoryStore()
	
	// Setup: Create 1000 tasks
	for i := 0; i < 1000; i++ {
		task := &models.Task{
			Name:   "Setup Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		updatedTask := &models.Task{
			Name:   "Updated Task",
			Status: 1,
		}
		store.Update((i%1000)+1, updatedTask)
	}
}

func BenchmarkMemoryStore_Delete(b *testing.B) {
	b.StopTimer()
	
	for i := 0; i < b.N; i++ {
		store := NewMemoryStore()
		
		// Setup: Create 1000 tasks
		for j := 0; j < 1000; j++ {
			task := &models.Task{
				Name:   "Setup Task",
				Status: 0,
			}
			store.Create(task)
		}
		
		b.StartTimer()
		store.Delete((i % 1000) + 1)
		b.StopTimer()
	}
}

func BenchmarkMemoryStore_ConcurrentCreate(b *testing.B) {
	store := NewMemoryStore()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &models.Task{
				Name:   "Concurrent Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}

func BenchmarkMemoryStore_ConcurrentRead(b *testing.B) {
	store := NewMemoryStore()
	
	// Setup: Create 1000 tasks
	for i := 0; i < 1000; i++ {
		task := &models.Task{
			Name:   "Setup Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			store.GetByID((i % 1000) + 1)
			i++
		}
	})
}