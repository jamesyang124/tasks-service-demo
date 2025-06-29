package benchmarks

import (
	"testing"
	"tasks-service-demo/internal/models"
	"tasks-service-demo/internal/storage"
)

// MemoryStore Benchmarks - Single mutex in-memory storage

func BenchmarkReadZipf_MemoryStore(b *testing.B) {
	store := storage.NewMemoryStore()
	BenchmarkReadZipf(b, store, "MemoryStore")
}

func BenchmarkWriteZipf_MemoryStore(b *testing.B) {
	store := storage.NewMemoryStore()
	BenchmarkWriteZipf(b, store, "MemoryStore")
}

func BenchmarkDistributedRead_MemoryStore(b *testing.B) {
	store := storage.NewMemoryStore()
	PopulateStore(b, store, "MemoryStore Distributed Read")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Uniform distribution across all keys
			targetID := (i % DatasetSize) + 1
			store.GetByID(targetID)
			i++
		}
	})
}

func BenchmarkDistributedWrite_MemoryStore(b *testing.B) {
	store := storage.NewMemoryStore()
	PopulateStore(b, store, "MemoryStore Distributed Write")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % DatasetSize) + 1
			updatedTask := &models.Task{
				Name:   "Distributed Update Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

func BenchmarkDistributedMixed_MemoryStore(b *testing.B) {
	store := storage.NewMemoryStore()
	PopulateStore(b, store, "MemoryStore Distributed Mixed")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % DatasetSize) + 1
			
			// 70% reads, 30% writes
			if i%10 < 7 {
				store.GetByID(targetID)
			} else {
				updatedTask := &models.Task{
					Name:   "Mixed Update Task",
					Status: i % 2,
				}
				store.Update(targetID, updatedTask)
			}
			i++
		}
	})
}