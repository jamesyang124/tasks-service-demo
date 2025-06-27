package storage

import (
	"testing"
	"tasks-service-demo/internal/models"
)

// Hot key benchmarks - simulate scenarios where same keys are accessed frequently
// This tests lock contention and shard distribution effectiveness

// Single hot key - all operations target the same task ID
func BenchmarkHotKeySingleWrite_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-create the hot key
	hotTask := &models.Task{
		ID:     1,
		Name:   "Hot Key Task",
		Status: 0,
	}
	store.Create(hotTask)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			updatedTask := &models.Task{
				Name:   "Updated Hot Task",
				Status: 1,
			}
			store.Update(1, updatedTask) // Always update same key
		}
	})
}

func BenchmarkHotKeySingleWrite_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-create the hot key
	hotTask := &models.Task{
		ID:     1,
		Name:   "Hot Key Task",
		Status: 0,
	}
	store.Create(hotTask)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			updatedTask := &models.Task{
				Name:   "Updated Hot Task",
				Status: 1,
			}
			store.Update(1, updatedTask) // Always update same key
		}
	})
}

// Single hot key - all reads target the same task ID
func BenchmarkHotKeySingleRead_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-create the hot key
	hotTask := &models.Task{
		ID:     1,
		Name:   "Hot Key Task",
		Status: 0,
	}
	store.Create(hotTask)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.GetByID(1) // Always read same key
		}
	})
}

func BenchmarkHotKeySingleRead_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-create the hot key
	hotTask := &models.Task{
		ID:     1,
		Name:   "Hot Key Task",
		Status: 0,
	}
	store.Create(hotTask)
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			store.GetByID(1) // Always read same key
		}
	})
}

// Multiple hot keys - 10 keys get 80% of traffic
func BenchmarkHotKeyMultipleWrite_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-create hot keys (1-10)
	for i := 1; i <= 10; i++ {
		hotTask := &models.Task{
			ID:     i,
			Name:   "Hot Key Task",
			Status: 0,
		}
		store.Create(hotTask)
	}
	
	// Create some cold keys (11-1000)
	for i := 11; i <= 1000; i++ {
		coldTask := &models.Task{
			ID:     i,
			Name:   "Cold Key Task",
			Status: 0,
		}
		store.Create(coldTask)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			if i%10 < 8 { // 80% traffic to hot keys (1-10)
				targetID = (i % 10) + 1
			} else { // 20% traffic to cold keys (11-1000)
				targetID = (i % 990) + 11
			}
			
			updatedTask := &models.Task{
				Name:   "Updated Task",
				Status: 1,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

func BenchmarkHotKeyMultipleWrite_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-create hot keys (1-10)
	for i := 1; i <= 10; i++ {
		hotTask := &models.Task{
			ID:     i,
			Name:   "Hot Key Task",
			Status: 0,
		}
		store.Create(hotTask)
	}
	
	// Create some cold keys (11-1000)
	for i := 11; i <= 1000; i++ {
		coldTask := &models.Task{
			ID:     i,
			Name:   "Cold Key Task",
			Status: 0,
		}
		store.Create(coldTask)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			if i%10 < 8 { // 80% traffic to hot keys (1-10)
				targetID = (i % 10) + 1
			} else { // 20% traffic to cold keys (11-1000)
				targetID = (i % 990) + 11
			}
			
			updatedTask := &models.Task{
				Name:   "Updated Task",
				Status: 1,
			}
			store.Update(targetID, updatedTask)
			i++
		}
	})
}

// Interleaved hot read-write - same keys being read and written simultaneously
func BenchmarkHotKeyInterleaved_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-create hot keys (1-5)
	for i := 1; i <= 5; i++ {
		hotTask := &models.Task{
			ID:     i,
			Name:   "Hot Interleaved Task",
			Status: 0,
		}
		store.Create(hotTask)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hotKeyID := (i % 5) + 1 // Rotate between hot keys 1-5
			
			if i%2 == 0 { // 50% reads
				store.GetByID(hotKeyID)
			} else { // 50% writes
				updatedTask := &models.Task{
					Name:   "Interleaved Update",
					Status: i % 2,
				}
				store.Update(hotKeyID, updatedTask)
			}
			i++
		}
	})
}

func BenchmarkHotKeyInterleaved_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-create hot keys (1-5)
	for i := 1; i <= 5; i++ {
		hotTask := &models.Task{
			ID:     i,
			Name:   "Hot Interleaved Task",
			Status: 0,
		}
		store.Create(hotTask)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			hotKeyID := (i % 5) + 1 // Rotate between hot keys 1-5
			
			if i%2 == 0 { // 50% reads
				store.GetByID(hotKeyID)
			} else { // 50% writes
				updatedTask := &models.Task{
					Name:   "Interleaved Update",
					Status: i % 2,
				}
				store.Update(hotKeyID, updatedTask)
			}
			i++
		}
	})
}

// Zipf distribution - realistic hot key distribution (20% of keys get 80% of traffic)
func BenchmarkHotKeyZipf_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-create 1000 keys
	for i := 1; i <= 1000; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Zipf Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			
			// Zipf-like distribution: 80% traffic to first 200 keys
			if i%10 < 8 {
				targetID = (i % 200) + 1
			} else {
				targetID = (i % 800) + 201
			}
			
			if i%3 == 0 { // 33% reads
				store.GetByID(targetID)
			} else { // 67% writes
				updatedTask := &models.Task{
					Name:   "Zipf Update",
					Status: i % 2,
				}
				store.Update(targetID, updatedTask)
			}
			i++
		}
	})
}

func BenchmarkHotKeyZipf_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-create 1000 keys
	for i := 1; i <= 1000; i++ {
		task := &models.Task{
			ID:     i,
			Name:   "Zipf Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			var targetID int
			
			// Zipf-like distribution: 80% traffic to first 200 keys
			if i%10 < 8 {
				targetID = (i % 200) + 1
			} else {
				targetID = (i % 800) + 201
			}
			
			if i%3 == 0 { // 33% reads
				store.GetByID(targetID)
			} else { // 67% writes
				updatedTask := &models.Task{
					Name:   "Zipf Update",
					Status: i % 2,
				}
				store.Update(targetID, updatedTask)
			}
			i++
		}
	})
}

// Worst case - all traffic to same shard in ShardStore
func BenchmarkHotKeyWorstCase_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-create keys that will all map to same shard (multiples of shard count)
	hotKeys := []int{1, 17, 33, 49, 65} // These will map to same shard in 16-shard store
	
	for _, id := range hotKeys {
		task := &models.Task{
			ID:     id,
			Name:   "Worst Case Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := hotKeys[i%len(hotKeys)]
			
			if i%2 == 0 {
				store.GetByID(targetID)
			} else {
				updatedTask := &models.Task{
					Name:   "Worst Case Update",
					Status: i % 2,
				}
				store.Update(targetID, updatedTask)
			}
			i++
		}
	})
}

func BenchmarkHotKeyWorstCase_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-create keys that will all map to same shard (multiples of shard count)
	hotKeys := []int{1, 17, 33, 49, 65} // These will map to same shard in 16-shard store
	
	for _, id := range hotKeys {
		task := &models.Task{
			ID:     id,
			Name:   "Worst Case Task",
			Status: 0,
		}
		store.Create(task)
	}
	
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			targetID := hotKeys[i%len(hotKeys)]
			
			if i%2 == 0 {
				store.GetByID(targetID)
			} else {
				updatedTask := &models.Task{
					Name:   "Worst Case Update",
					Status: i % 2,
				}
				store.Update(targetID, updatedTask)
			}
			i++
		}
	})
}

// Thundering herd - many goroutines accessing same key at exact same time
func BenchmarkHotKeyThunderingHerd_MemoryStore(b *testing.B) {
	store := NewMemoryStore()
	
	// Pre-create the target key
	task := &models.Task{
		ID:     42,
		Name:   "Thundering Herd Task",
		Status: 0,
	}
	store.Create(task)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		
		// Launch 100 goroutines simultaneously accessing same key
		start := make(chan struct{})
		done := make(chan struct{}, 100)
		
		for j := 0; j < 100; j++ {
			go func() {
				<-start // Wait for signal
				
				// Mix of reads and writes
				if j%2 == 0 {
					store.GetByID(42)
				} else {
					updatedTask := &models.Task{
						Name:   "Herd Update",
						Status: j % 2,
					}
					store.Update(42, updatedTask)
				}
				done <- struct{}{}
			}()
		}
		
		b.StartTimer()
		close(start) // Release the herd
		
		// Wait for all goroutines to complete
		for j := 0; j < 100; j++ {
			<-done
		}
	}
}

func BenchmarkHotKeyThunderingHerd_ShardStore(b *testing.B) {
	store := NewShardStore(16)
	
	// Pre-create the target key
	task := &models.Task{
		ID:     42,
		Name:   "Thundering Herd Task",
		Status: 0,
	}
	store.Create(task)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		
		// Launch 100 goroutines simultaneously accessing same key
		start := make(chan struct{})
		done := make(chan struct{}, 100)
		
		for j := 0; j < 100; j++ {
			go func() {
				<-start // Wait for signal
				
				// Mix of reads and writes
				if j%2 == 0 {
					store.GetByID(42)
				} else {
					updatedTask := &models.Task{
						Name:   "Herd Update",
						Status: j % 2,
					}
					store.Update(42, updatedTask)
				}
				done <- struct{}{}
			}()
		}
		
		b.StartTimer()
		close(start) // Release the herd
		
		// Wait for all goroutines to complete
		for j := 0; j < 100; j++ {
			<-done
		}
	}
}