package benchmarks

import (
	"tasks-service-demo/internal/models"
	"tasks-service-demo/internal/storage/bigcache"
	"testing"
)

// BigCacheStore Benchmarks - Off-heap cache with zero GC overhead

func BenchmarkReadZipf_BigCacheStore(b *testing.B) {
	store := bigcache.NewBigCacheStore()
	defer store.Close()
	BenchmarkReadZipf(b, store, "BigCacheStore")
}

func BenchmarkWriteZipf_BigCacheStore(b *testing.B) {
	store := bigcache.NewBigCacheStore()
	defer store.Close()
	BenchmarkWriteZipf(b, store, "BigCacheStore")
}

func BenchmarkDistributedRead_BigCacheStore(b *testing.B) {
	store := bigcache.NewBigCacheStore()
	defer store.Close()
	PopulateStore(b, store, "BigCacheStore Distributed Read")

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

func BenchmarkDistributedWrite_BigCacheStore(b *testing.B) {
	store := bigcache.NewBigCacheStore()
	defer store.Close()
	PopulateStore(b, store, "BigCacheStore Distributed Write")

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

func BenchmarkDistributedMixed_BigCacheStore(b *testing.B) {
	store := bigcache.NewBigCacheStore()
	defer store.Close()
	PopulateStore(b, store, "BigCacheStore Distributed Mixed")

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
