package benchmarks

import (
	"fmt"
	"testing"
	"tasks-service-demo/internal/models"
	"tasks-service-demo/internal/storage"
)

const (
	DatasetSize = 1000000 // 1M dataset for realistic performance testing
	HotKeyRatio = 20      // 20% of keys get 80% of traffic (Zipf distribution)
)

// Common benchmark utilities for all storage implementations

// PopulateStore fills a store with test data and logs progress
func PopulateStore(b *testing.B, store storage.Store, storeName string) {
	b.Logf("Setting up %d tasks for %s", DatasetSize, storeName)
	
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   fmt.Sprintf("%s Task %d", storeName, i),
			Status: i % 2,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
}

// GetZipfTargetID returns a target ID following Zipf distribution (80/20 rule)
func GetZipfTargetID(iteration int) int {
	hotKeyCount := DatasetSize / 5 // 20% hot keys (200K keys)
	
	// 80% traffic to hot keys, 20% to cold keys
	if iteration%10 < 8 {
		return (iteration % hotKeyCount) + 1
	} else {
		return (iteration % (DatasetSize - hotKeyCount)) + hotKeyCount + 1
	}
}

// BenchmarkReadZipf provides a standardized Zipf read benchmark for any store
func BenchmarkReadZipf(b *testing.B, store storage.Store, storeName string) {
	PopulateStore(b, store, storeName)
	hotKeyCount := DatasetSize / 5
	
	b.Logf("Setup complete. Starting read benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := GetZipfTargetID(i)
			store.GetByID(targetID)
			i++
		}
	})
}

// BenchmarkWriteZipf provides a standardized Zipf write benchmark for any store
func BenchmarkWriteZipf(b *testing.B, store storage.Store, storeName string) {
	PopulateStore(b, store, storeName)
	hotKeyCount := DatasetSize / 5
	
	b.Logf("Setup complete. Starting write benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := GetZipfTargetID(i)
			
			updatedTask := &models.Task{
				Name:   fmt.Sprintf("Updated %s Task %d", storeName, i),
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

// Note: Uses storage.Store interface from internal/storage/store.go