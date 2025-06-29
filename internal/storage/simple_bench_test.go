package storage

import (
	"testing"
	"tasks-service-demo/internal/models"
)

const (
	DatasetSize = 1000000 // 1M dataset for realistic performance testing
	HotKeyRatio = 20      // 20% of keys get 80% of traffic (Zipf distribution)
)

// Read-only Zipf distribution - realistic read-heavy workload
func BenchmarkReadZipf_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	hotKeyCount := DatasetSize / 5 // 20% hot keys (200K keys)
	
	b.Logf("Setting up %d tasks for Read Zipf test", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Read Zipf Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting read benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			// 80% traffic to hot keys, 20% to cold keys
			if i%10 < 8 {
				targetID = (i % hotKeyCount) + 1
			} else {
				targetID = (i % (DatasetSize - hotKeyCount)) + hotKeyCount + 1
			}
			store.GetByID(targetID)
			i++
		}
	})
}

func BenchmarkReadZipf_ShardStore(b *testing.B) {
	store := NewShardStore(32) // More shards for 1M dataset
	hotKeyCount := DatasetSize / 5
	
	b.Logf("Setting up %d tasks for Read Zipf test with 32 shards", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Read Zipf Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting read benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			if i%10 < 8 {
				targetID = (i % hotKeyCount) + 1
			} else {
				targetID = (i % (DatasetSize - hotKeyCount)) + hotKeyCount + 1
			}
			store.GetByID(targetID)
			i++
		}
	})
}

func BenchmarkReadZipf_BigCacheStore(b *testing.B) {
	store := NewBigCacheStore()
	defer store.Close()
	hotKeyCount := DatasetSize / 5
	
	b.Logf("Setting up %d tasks for Read Zipf test with BigCache", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Read Zipf Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting read benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			// 80% traffic to hot keys, 20% to cold keys
			if i%10 < 8 {
				targetID = (i % hotKeyCount) + 1
			} else {
				targetID = (i % (DatasetSize - hotKeyCount)) + hotKeyCount + 1
			}
			store.GetByID(targetID)
			i++
		}
	})
}

func BenchmarkReadZipf_ChannelStore(b *testing.B) {
	store := NewChannelStore(4) // Optimized worker count
	defer store.Shutdown()
	hotKeyCount := DatasetSize / 5
	
	b.Logf("Setting up %d tasks for Read Zipf test with 4 workers", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Read Zipf Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting read benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			if i%10 < 8 {
				targetID = (i % hotKeyCount) + 1
			} else {
				targetID = (i % (DatasetSize - hotKeyCount)) + hotKeyCount + 1
			}
			store.GetByID(targetID)
			i++
		}
	})
}

// Write-only Zipf distribution - write-heavy workload with hot keys
func BenchmarkWriteZipf_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	hotKeyCount := DatasetSize / 5
	
	b.Logf("Setting up %d tasks for Write Zipf test", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Write Zipf Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting write benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			// 80% traffic to hot keys, 20% to cold keys
			if i%10 < 8 {
				targetID = (i % hotKeyCount) + 1
			} else {
				targetID = (i % (DatasetSize - hotKeyCount)) + hotKeyCount + 1
			}
			
			updatedTask := &models.Task{
				Name:   "Updated Zipf Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

func BenchmarkWriteZipf_ShardStore(b *testing.B) {
	store := NewShardStore(32)
	hotKeyCount := DatasetSize / 5
	
	b.Logf("Setting up %d tasks for Write Zipf test with 32 shards", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Write Zipf Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting write benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			if i%10 < 8 {
				targetID = (i % hotKeyCount) + 1
			} else {
				targetID = (i % (DatasetSize - hotKeyCount)) + hotKeyCount + 1
			}
			
			updatedTask := &models.Task{
				Name:   "Updated Zipf Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

func BenchmarkWriteZipf_BigCacheStore(b *testing.B) {
	store := NewBigCacheStore()
	defer store.Close()
	hotKeyCount := DatasetSize / 5
	
	b.Logf("Setting up %d tasks for Write Zipf test with BigCache", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Write Zipf Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting write benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			// 80% traffic to hot keys, 20% to cold keys
			if i%10 < 8 {
				targetID = (i % hotKeyCount) + 1
			} else {
				targetID = (i % (DatasetSize - hotKeyCount)) + hotKeyCount + 1
			}
			
			updatedTask := &models.Task{
				Name:   "Updated Zipf Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

func BenchmarkWriteZipf_ChannelStore(b *testing.B) {
	store := NewChannelStore(4)
	defer store.Shutdown()
	hotKeyCount := DatasetSize / 5
	
	b.Logf("Setting up %d tasks for Write Zipf test with 4 workers", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Write Zipf Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting write benchmark with %d hot keys", hotKeyCount)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			if i%10 < 8 {
				targetID = (i % hotKeyCount) + 1
			} else {
				targetID = (i % (DatasetSize - hotKeyCount)) + hotKeyCount + 1
			}
			
			updatedTask := &models.Task{
				Name:   "Updated Zipf Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

// Truly distributed access - uniform random access (no hot keys)
func BenchmarkDistributedRead_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	b.Logf("Setting up %d tasks for Distributed Read test", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Distributed Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting distributed read benchmark")
	
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

func BenchmarkDistributedRead_ShardStore(b *testing.B) {
	store := NewShardStore(32)
	
	b.Logf("Setting up %d tasks for Distributed Read test with 32 shards", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Distributed Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting distributed read benchmark")
	
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

func BenchmarkDistributedRead_ChannelStore(b *testing.B) {
	store := NewChannelStore(4)
	defer store.Shutdown()
	
	b.Logf("Setting up %d tasks for Distributed Read test with 4 workers", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Distributed Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting distributed read benchmark")
	
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

// Truly distributed writes - uniform random writes (no hot keys)
func BenchmarkDistributedWrite_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	b.Logf("Setting up %d tasks for Distributed Write test", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Distributed Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting distributed write benchmark")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Uniform distribution across all keys
			targetID := (i % DatasetSize) + 1
			updatedTask := &models.Task{
				Name:   "Updated Distributed Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

func BenchmarkDistributedWrite_ShardStore(b *testing.B) {
	store := NewShardStore(32)
	
	b.Logf("Setting up %d tasks for Distributed Write test with 32 shards", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Distributed Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting distributed write benchmark")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Uniform distribution across all keys
			targetID := (i % DatasetSize) + 1
			updatedTask := &models.Task{
				Name:   "Updated Distributed Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

func BenchmarkDistributedWrite_ChannelStore(b *testing.B) {
	store := NewChannelStore(4)
	defer store.Shutdown()
	
	b.Logf("Setting up %d tasks for Distributed Write test with 4 workers", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Distributed Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting distributed write benchmark")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Uniform distribution across all keys
			targetID := (i % DatasetSize) + 1
			updatedTask := &models.Task{
				Name:   "Updated Distributed Task",
				Status: i % 2,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

// Mixed distributed access - 50% reads, 50% writes, uniform distribution
func BenchmarkDistributedMixed_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	b.Logf("Setting up %d tasks for Distributed Mixed test", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Mixed Distributed Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting distributed mixed benchmark")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Uniform distribution across all keys
			targetID := (i % DatasetSize) + 1
			
			if i%2 == 0 { // 50% reads
				store.GetByID(targetID)
			} else { // 50% writes
				updatedTask := &models.Task{
					Name:   "Updated Mixed Task",
					Status: i % 2,
				}
				store.Update(targetID, updatedTask)
			}
			i++
		}
	})
}

func BenchmarkDistributedMixed_ShardStore(b *testing.B) {
	store := NewShardStore(32)
	
	b.Logf("Setting up %d tasks for Distributed Mixed test with 32 shards", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Mixed Distributed Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting distributed mixed benchmark")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Uniform distribution across all keys
			targetID := (i % DatasetSize) + 1
			
			if i%2 == 0 { // 50% reads
				store.GetByID(targetID)
			} else { // 50% writes
				updatedTask := &models.Task{
					Name:   "Updated Mixed Task",
					Status: i % 2,
				}
				store.Update(targetID, updatedTask)
			}
			i++
		}
	})
}

func BenchmarkDistributedMixed_ChannelStore(b *testing.B) {
	store := NewChannelStore(4)
	defer store.Shutdown()
	
	b.Logf("Setting up %d tasks for Distributed Mixed test with 4 workers", DatasetSize)
	// Pre-populate dataset
	for i := 1; i <= DatasetSize; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Mixed Distributed Task",
			Status: 0,
		}
		store.Create(task)
		
		if i%200000 == 0 {
			b.Logf("Created %d/%d tasks", i, DatasetSize)
		}
	}
	b.Logf("Setup complete. Starting distributed mixed benchmark")
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			// Uniform distribution across all keys
			targetID := (i % DatasetSize) + 1
			
			if i%2 == 0 { // 50% reads
				store.GetByID(targetID)
			} else { // 50% writes
				updatedTask := &models.Task{
					Name:   "Updated Mixed Task",
					Status: i % 2,
				}
				store.Update(targetID, updatedTask)
			}
			i++
		}
	})
}