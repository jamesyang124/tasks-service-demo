package benchmarks

import (
	"tasks-service-demo/internal/entities"
	"tasks-service-demo/internal/storage/channel"
	"testing"
)

// ChannelStore Benchmarks - Actor model with message passing

func BenchmarkReadZipf_ChannelStore(b *testing.B) {
	store := channel.NewChannelStore(4) // Optimized worker count
	defer store.Shutdown()
	BenchmarkReadZipf(b, store, "ChannelStore")
}

func BenchmarkWriteZipf_ChannelStore(b *testing.B) {
	store := channel.NewChannelStore(4)
	defer store.Shutdown()
	BenchmarkWriteZipf(b, store, "ChannelStore")
}

func BenchmarkDistributedRead_ChannelStore(b *testing.B) {
	store := channel.NewChannelStore(4)
	defer store.Shutdown()
	PopulateStore(b, store, "ChannelStore Distributed Read")

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

func BenchmarkDistributedWrite_ChannelStore(b *testing.B) {
	store := channel.NewChannelStore(4)
	defer store.Shutdown()
	PopulateStore(b, store, "ChannelStore Distributed Write")

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

func BenchmarkDistributedMixed_ChannelStore(b *testing.B) {
	store := channel.NewChannelStore(4)
	defer store.Shutdown()
	PopulateStore(b, store, "ChannelStore Distributed Mixed")

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
