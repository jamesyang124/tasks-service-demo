package storage

import (
	"testing"
	"tasks-service-demo/internal/models"
)

func BenchmarkShardStore_Create(b *testing.B) {
	store := NewShardStore(8)
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		task := &models.Task{
			Name:   "Benchmark Task",
			Status: 0,
		}
		store.Create(task)
	}
}

func BenchmarkShardStore_GetByID(b *testing.B) {
	store := NewShardStore(8)
	
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

func BenchmarkShardStore_GetAll(b *testing.B) {
	store := NewShardStore(8)
	
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

func BenchmarkShardStore_Update(b *testing.B) {
	store := NewShardStore(8)
	
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

func BenchmarkShardStore_Delete(b *testing.B) {
	b.StopTimer()
	
	for i := 0; i < b.N; i++ {
		store := NewShardStore(8)
		
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

func BenchmarkShardStore_ConcurrentCreate(b *testing.B) {
	store := NewShardStore(8)
	
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

func BenchmarkShardStore_ConcurrentRead(b *testing.B) {
	store := NewShardStore(8)
	
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

// Shard-specific benchmarks
func BenchmarkShardStore_4Shards(b *testing.B) {
	store := NewShardStore(4)
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &models.Task{
				Name:   "4-Shard Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}

func BenchmarkShardStore_8Shards(b *testing.B) {
	store := NewShardStore(8)
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &models.Task{
				Name:   "8-Shard Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}

func BenchmarkShardStore_16Shards(b *testing.B) {
	store := NewShardStore(16)
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &models.Task{
				Name:   "16-Shard Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}

func BenchmarkShardStore_32Shards(b *testing.B) {
	store := NewShardStore(32)
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &models.Task{
				Name:   "32-Shard Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}