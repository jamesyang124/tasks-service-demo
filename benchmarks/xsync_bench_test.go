package benchmarks

import (
	"tasks-service-demo/internal/entities"
	"tasks-service-demo/internal/storage/xsync"
	"testing"
)

// XSyncStore Benchmarks - Lock-free concurrent map storage

func BenchmarkReadZipf_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	BenchmarkReadZipf(b, store, "XSyncStore")
}

func BenchmarkWriteZipf_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	BenchmarkWriteZipf(b, store, "XSyncStore")
}

func BenchmarkDistributedRead_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	PopulateStore(b, store, "XSyncStore Distributed Read")

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

func BenchmarkDistributedWrite_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	PopulateStore(b, store, "XSyncStore Distributed Write")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % DatasetSize) + 1
			updatedTask := &entities.Task{
				Name:   "Distributed Update Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

func BenchmarkDistributedMixed_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	PopulateStore(b, store, "XSyncStore Distributed Mixed")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % DatasetSize) + 1

			// 70% reads, 30% writes
			if i%10 < 7 {
				store.GetByID(targetID)
			} else {
				updatedTask := &entities.Task{
					Name:   "Mixed Update Task",
					Status: i % 2,
				}
				store.Update(targetID, updatedTask)
			}
			i++
		}
	})
}

// Additional benchmarks specific to xsync performance characteristics

func BenchmarkCreate_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			task := &entities.Task{
				Name:   "Benchmark Task",
				Status: i % 2,
			}
			store.Create(task)
			i++
		}
	})
}

func BenchmarkGetByID_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	PopulateStore(b, store, "XSyncStore GetByID")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % DatasetSize) + 1
			store.GetByID(targetID)
			i++
		}
	})
}

func BenchmarkUpdate_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	PopulateStore(b, store, "XSyncStore Update")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := (i % DatasetSize) + 1
			updatedTask := &entities.Task{
				Name:   "Updated Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

func BenchmarkDelete_XSyncStore(b *testing.B) {
	// We need to repopulate for each sub-benchmark since delete is destructive
	b.Run("Delete", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			store := xsync.NewXSyncStore()
			PopulateStore(b, store, "XSyncStore Delete")
			b.StartTimer()
			
			targetID := (i % DatasetSize) + 1
			store.Delete(targetID)
		}
	})
}

func BenchmarkGetAll_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	PopulateStore(b, store, "XSyncStore GetAll")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetAll()
	}
}

// Contention-focused benchmarks to test lock-free performance under high contention

func BenchmarkHighContentionRead_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	PopulateStore(b, store, "XSyncStore High Contention Read")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// All goroutines competing for the same few keys (high contention)
			targetID := (b.N % 10) + 1
			store.GetByID(targetID)
		}
	})
}

func BenchmarkHighContentionWrite_XSyncStore(b *testing.B) {
	store := xsync.NewXSyncStore()
	PopulateStore(b, store, "XSyncStore High Contention Write")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// All goroutines competing for the same few keys (high contention)
			targetID := (i % 10) + 1
			updatedTask := &entities.Task{
				Name:   "High Contention Update",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}