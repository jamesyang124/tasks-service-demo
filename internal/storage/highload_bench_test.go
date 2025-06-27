package storage

import (
	"runtime"
	"sync"
	"testing"
	"tasks-service-demo/internal/models"
)

// High-load write benchmarks - simulate heavy write scenarios
func BenchmarkHighLoadWrite_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &models.Task{
				Name:   "High Load Write Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}

func BenchmarkHighLoadWrite_ShardStore(b *testing.B) {
	store := NewShardStore(16) // More shards for high load
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			task := &models.Task{
				Name:   "High Load Write Task",
				Status: 0,
			}
			store.Create(task)
		}
	})
}

// High-load read benchmarks - simulate heavy read scenarios
func BenchmarkHighLoadRead_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-populate with 100K tasks
	for i := 0; i < 100000; i++ {
		task := &models.Task{
			Name:   "Read Load Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			store.GetByID((i % 100000) + 1)
			i++
		}
	})
}

func BenchmarkHighLoadRead_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-populate with 100K tasks
	for i := 0; i < 100000; i++ {
		task := &models.Task{
			Name:   "Read Load Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			store.GetByID((i % 100000) + 1)
			i++
		}
	})
}

// Extreme concurrent write load - burst writes
func BenchmarkBurstWrite_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numGoroutines := 100
		
		b.StartTimer()
		for j := 0; j < numGoroutines; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for k := 0; k < 100; k++ {
					task := &models.Task{
						Name:   "Burst Task",
						Status: 0,
					}
					store.Create(task)
				}
			}()
		}
		wg.Wait()
		b.StopTimer()
	}
}

func BenchmarkBurstWrite_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	for i := 0; i < b.N; i++ {
		var wg sync.WaitGroup
		numGoroutines := 100
		
		b.StartTimer()
		for j := 0; j < numGoroutines; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for k := 0; k < 100; k++ {
					task := &models.Task{
						Name:   "Burst Task",
						Status: 0,
					}
					store.Create(task)
				}
			}()
		}
		wg.Wait()
		b.StopTimer()
	}
}

// Heavy mixed workload - 70% reads, 30% writes
func BenchmarkHeavyMixed_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-populate
	for i := 0; i < 10000; i++ {
		task := &models.Task{
			Name:   "Mixed Load Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 < 7 { // 70% reads
				store.GetByID((i % 10000) + 1)
			} else { // 30% writes
				task := &models.Task{
					Name:   "Heavy Mixed Task",
					Status: 0,
				}
				store.Create(task)
			}
			i++
		}
	})
}

func BenchmarkHeavyMixed_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-populate
	for i := 0; i < 10000; i++ {
		task := &models.Task{
			Name:   "Mixed Load Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 < 7 { // 70% reads
				store.GetByID((i % 10000) + 1)
			} else { // 30% writes
				task := &models.Task{
					Name:   "Heavy Mixed Task",
					Status: 0,
				}
				store.Create(task)
			}
			i++
		}
	})
}

// Write-heavy workload - 80% writes, 20% reads
func BenchmarkWriteHeavy_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-populate
	for i := 0; i < 1000; i++ {
		task := &models.Task{
			Name:   "Write Heavy Init Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%5 < 4 { // 80% writes
				task := &models.Task{
					Name:   "Write Heavy Task",
					Status: 0,
				}
				store.Create(task)
			} else { // 20% reads
				store.GetByID((i % 1000) + 1)
			}
			i++
		}
	})
}

func BenchmarkWriteHeavy_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-populate
	for i := 0; i < 1000; i++ {
		task := &models.Task{
			Name:   "Write Heavy Init Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%5 < 4 { // 80% writes
				task := &models.Task{
					Name:   "Write Heavy Task",
					Status: 0,
				}
				store.Create(task)
			} else { // 20% reads
				store.GetByID((i % 1000) + 1)
			}
			i++
		}
	})
}

// Read-heavy workload - 90% reads, 10% writes
func BenchmarkReadHeavy_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-populate with substantial data
	for i := 0; i < 50000; i++ {
		task := &models.Task{
			Name:   "Read Heavy Init Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 < 9 { // 90% reads
				store.GetByID((i % 50000) + 1)
			} else { // 10% writes
				task := &models.Task{
					Name:   "Read Heavy Task",
					Status: 0,
				}
				store.Create(task)
			}
			i++
		}
	})
}

func BenchmarkReadHeavy_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-populate with substantial data
	for i := 0; i < 50000; i++ {
		task := &models.Task{
			Name:   "Read Heavy Init Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			if i%10 < 9 { // 90% reads
				store.GetByID((i % 50000) + 1)
			} else { // 10% writes
				task := &models.Task{
					Name:   "Read Heavy Task",
					Status: 0,
				}
				store.Create(task)
			}
			i++
		}
	})
}

// Stress test - maximum concurrent operations
func BenchmarkStress_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	maxWorkers := runtime.NumCPU() * 8
	
	b.SetParallelism(maxWorkers)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 3 {
			case 0:
				task := &models.Task{
					Name:   "Stress Task",
					Status: 0,
				}
				store.Create(task)
			case 1:
				store.GetByID((i % 1000) + 1)
			case 2:
				store.GetAll()
			}
			i++
		}
	})
}

func BenchmarkStress_ShardStore(b *testing.B) {
	store := NewShardStore(32) // Maximum shards for stress test
	maxWorkers := runtime.NumCPU() * 8
	
	b.SetParallelism(maxWorkers)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			switch i % 3 {
			case 0:
				task := &models.Task{
					Name:   "Stress Task",
					Status: 0,
				}
				store.Create(task)
			case 1:
				store.GetByID((i % 1000) + 1)
			case 2:
				store.GetAll()
			}
			i++
		}
	})
}