package benchmarks

import (
	"fmt"
	"testing"
	"tasks-service-demo/internal/models"
	"tasks-service-demo/internal/storage"
)

// ShardStore Benchmarks - Optimized dedicated worker per shard storage

func BenchmarkReadZipf_ShardStore(b *testing.B) {
	store := storage.NewShardStore(32) // 32 shards for 1M dataset
	defer store.Close()
	BenchmarkReadZipf(b, store, "ShardStore")
}

func BenchmarkWriteZipf_ShardStore(b *testing.B) {
	store := storage.NewShardStore(32)
	defer store.Close()
	BenchmarkWriteZipf(b, store, "ShardStore")
}

func BenchmarkDistributedRead_ShardStore(b *testing.B) {
	store := storage.NewShardStore(32)
	defer store.Close()
	PopulateStore(b, store, "ShardStore Distributed Read")
	
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

func BenchmarkDistributedWrite_ShardStore(b *testing.B) {
	store := storage.NewShardStore(32)
	defer store.Close()
	PopulateStore(b, store, "ShardStore Distributed Write")
	
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

func BenchmarkDistributedMixed_ShardStore(b *testing.B) {
	store := storage.NewShardStore(32)
	defer store.Close()
	PopulateStore(b, store, "ShardStore Distributed Mixed")
	
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

// ShardStore specific benchmarks

func BenchmarkShardStore_GetAll(b *testing.B) {
	store := storage.NewShardStore(32)
	defer store.Close()
	
	PopulateStore(b, store, "ShardStore GetAll")
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = store.GetAll()
	}
}

func BenchmarkShardStore_CoreUtilization(b *testing.B) {
	for _, shardCount := range []int{4, 8, 16, 32} {
		b.Run(fmt.Sprintf("Shards_%d", shardCount), func(b *testing.B) {
			store := storage.NewShardStore(shardCount)
			defer store.Close()
			
			BenchmarkReadZipf(b, store, fmt.Sprintf("ShardStore_%dShards", shardCount))
		})
	}
}