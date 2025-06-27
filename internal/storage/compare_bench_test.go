package storage

import (
	"testing"
	"tasks-service-demo/internal/models"
)

// Direct comparison benchmarks between MemoryStore and ShardStore

func BenchmarkMemoryVsShard_Create(b *testing.B) {
	b.Run("MemoryStore", func(b *testing.B) {
		store := NewMemoryStore()
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			task := &models.Task{
				Name:   "Memory Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
	
	b.Run("ShardStore_8", func(b *testing.B) {
		store := NewShardStore(8)
		b.ResetTimer()
		
		for i := 0; i < b.N; i++ {
			task := &models.Task{
				Name:   "Shard Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}

func BenchmarkMemoryVsShard_ConcurrentCreate(b *testing.B) {
	b.Run("MemoryStore", func(b *testing.B) {
		store := NewMemoryStore()
		b.ResetTimer()
		
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				task := &models.Task{
					Name:   "Memory Concurrent Task",
					Status: 0,
				}
				store.Create(task)
			}
		})
	})
	
	b.Run("ShardStore_8", func(b *testing.B) {
		store := NewShardStore(8)
		b.ResetTimer()
		
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				task := &models.Task{
					Name:   "Shard Concurrent Task",
					Status: 0,
				}
				store.Create(task)
			}
		})
	})
}

func BenchmarkMemoryVsShard_ConcurrentRead(b *testing.B) {
	setupTasks := func(store Store) {
		for i := 0; i < 1000; i++ {
			task := &models.Task{
				Name:   "Setup Task",
				Status: 0,
			}
			store.Create(task)
		}
	}
	
	b.Run("MemoryStore", func(b *testing.B) {
		store := NewMemoryStore()
		setupTasks(store)
		b.ResetTimer()
		
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				store.GetByID((i % 1000) + 1)
				i++
			}
		})
	})
	
	b.Run("ShardStore_8", func(b *testing.B) {
		store := NewShardStore(8)
		setupTasks(store)
		b.ResetTimer()
		
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				store.GetByID((i % 1000) + 1)
				i++
			}
		})
	})
}

func BenchmarkMemoryVsShard_MixedWorkload(b *testing.B) {
	b.Run("MemoryStore", func(b *testing.B) {
		store := NewMemoryStore()
		b.ResetTimer()
		
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				switch i % 4 {
				case 0: // Create
					task := &models.Task{
						Name:   "Mixed Task",
						Status: 0,
					}
					store.Create(task)
				case 1: // Read
					store.GetByID((i % 100) + 1)
				case 2: // Update
					if i > 100 {
						updatedTask := &models.Task{
							Name:   "Updated Mixed Task",
							Status: 1,
						}
						store.Update((i%100)+1, updatedTask)
					}
				case 3: // GetAll
					store.GetAll()
				}
				i++
			}
		})
	})
	
	b.Run("ShardStore_8", func(b *testing.B) {
		store := NewShardStore(8)
		b.ResetTimer()
		
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				switch i % 4 {
				case 0: // Create
					task := &models.Task{
						Name:   "Mixed Task",
						Status: 0,
					}
					store.Create(task)
				case 1: // Read
					store.GetByID((i % 100) + 1)
				case 2: // Update
					if i > 100 {
						updatedTask := &models.Task{
							Name:   "Updated Mixed Task",
							Status: 1,
						}
						store.Update((i%100)+1, updatedTask)
					}
				case 3: // GetAll
					store.GetAll()
				}
				i++
			}
		})
	})
}

// Scale comparison - different data sizes
func BenchmarkMemoryVsShard_Scale1K(b *testing.B) {
	benchmarkScale(b, 1000)
}

func BenchmarkMemoryVsShard_Scale10K(b *testing.B) {
	benchmarkScale(b, 10000)
}

func BenchmarkMemoryVsShard_Scale100K(b *testing.B) {
	benchmarkScale(b, 100000)
}

func benchmarkScale(b *testing.B, numTasks int) {
	setupTasks := func(store Store, n int) {
		for i := 0; i < n; i++ {
			task := &models.Task{
				Name:   "Scale Task",
				Status: 0,
			}
			store.Create(task)
		}
	}
	
	b.Run("MemoryStore", func(b *testing.B) {
		store := NewMemoryStore()
		setupTasks(store, numTasks)
		b.ResetTimer()
		
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				store.GetByID((i % numTasks) + 1)
				i++
			}
		})
	})
	
	b.Run("ShardStore_8", func(b *testing.B) {
		store := NewShardStore(8)
		setupTasks(store, numTasks)
		b.ResetTimer()
		
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				store.GetByID((i % numTasks) + 1)
				i++
			}
		})
	})
}