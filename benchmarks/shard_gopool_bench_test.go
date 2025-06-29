package benchmarks

import (
	"fmt"
	"tasks-service-demo/internal/entities"
	"tasks-service-demo/internal/storage/shard"
	"testing"
)

// Benchmark the new gopool-based ShardStore implementation

func BenchmarkReadZipf_ShardStoreGopool(b *testing.B) {
	store := shard.NewShardStoreGopool(32) // 32 shards for 1M dataset
	defer store.Close()
	BenchmarkReadZipf(b, store, "ShardStoreGopool")
}

func BenchmarkWriteZipf_ShardStoreGopool(b *testing.B) {
	store := shard.NewShardStoreGopool(32)
	defer store.Close()
	BenchmarkWriteZipf(b, store, "ShardStoreGopool")
}

// Core utilization benchmarks for gopool implementation (M4 Pro has 14 cores)
func BenchmarkShardStoreGopool_CoreUtilization(b *testing.B) {
	for _, shardCount := range []int{4, 8, 16, 32} {
		b.Run(fmt.Sprintf("Shards_%d", shardCount), func(b *testing.B) {
			store := shard.NewShardStoreGopool(shardCount)
			defer store.Close()
			BenchmarkReadZipf(b, store, fmt.Sprintf("ShardStoreGopool_%dShards", shardCount))
		})
	}
}

// Comparison benchmark: Current vs Gopool implementation
func BenchmarkShardStore_Comparison(b *testing.B) {
	shardCount := 32

	b.Run("Current", func(b *testing.B) {
		store := shard.NewShardStore(shardCount)
		defer store.Close()
		BenchmarkReadZipf(b, store, "ShardStore_Current")
	})

	b.Run("Gopool", func(b *testing.B) {
		store := shard.NewShardStoreGopool(shardCount)
		defer store.Close()
		BenchmarkReadZipf(b, store, "ShardStoreGopool_New")
	})
}

// GetAll operation comparison (most impacted by worker pool changes)
func BenchmarkGetAll_Comparison(b *testing.B) {
	const setupSize = 100000 // Smaller dataset for GetAll tests

	b.Run("Current", func(b *testing.B) {
		store := shard.NewShardStore(32)
		defer store.Close()

		// Setup smaller dataset
		for i := 1; i <= setupSize; i++ {
			task := &entities.Task{
				ID:     i,
				Name:   fmt.Sprintf("Task %d", i),
				Status: i % 2,
			}
			store.Create(task)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = store.GetAll()
		}
	})

	b.Run("Gopool", func(b *testing.B) {
		store := shard.NewShardStoreGopool(32)
		defer store.Close()

		// Setup smaller dataset
		for i := 1; i <= setupSize; i++ {
			task := &entities.Task{
				ID:     i,
				Name:   fmt.Sprintf("Task %d", i),
				Status: i % 2,
			}
			store.Create(task)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = store.GetAll()
		}
	})
}
