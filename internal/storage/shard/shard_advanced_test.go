package shard

import (
	"fmt"
	"sync"
	"sync/atomic"
	"tasks-service-demo/internal/entities"
	"testing"
	"time"
)

func TestShardStore_PowerOfTwoOptimization(t *testing.T) {
	tests := []struct {
		input    int
		expected int
	}{
		{1, 1},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 8},
		{7, 8},
		{8, 8},
		{15, 16},
		{16, 16},
		{17, 32},
		{31, 32},
		{32, 32},
		{33, 64},
		{63, 64},
		{64, 64},
		{65, 128},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input_%d", tt.input), func(t *testing.T) {
			store := NewShardStore(tt.input)
			if store.numShards != tt.expected {
				t.Errorf("Expected %d shards for input %d, got %d", tt.expected, tt.input, store.numShards)
			}

			// Verify it's actually a power of 2
			if !isPowerOfTwo(store.numShards) {
				t.Errorf("Result %d is not a power of 2", store.numShards)
			}
		})
	}
}

func TestShardStore_HashDistribution(t *testing.T) {
	store := NewShardStore(8)
	distribution := make(map[int]int)

	// Test hash distribution with many IDs
	for id := 1; id <= 1000; id++ {
		shardIndex := store.getShardByID(id)
		distribution[shardIndex]++
	}

	// Verify all shards are used
	for i := 0; i < 8; i++ {
		if distribution[i] == 0 {
			t.Errorf("Shard %d was never used", i)
		}
	}

	// Check distribution is reasonably balanced (no shard should have more than 200% of average)
	average := 1000 / 8
	maxAllowed := average * 2
	for shard, count := range distribution {
		if count > maxAllowed {
			t.Errorf("Shard %d has %d items, which is more than %d (2x average)", shard, count, maxAllowed)
		}
	}
}

func TestShardStore_HighConcurrency(t *testing.T) {
	store := NewShardStore(16)
	numGoroutines := 100
	tasksPerGoroutine := 100
	var wg sync.WaitGroup
	var successCount int64

	// Concurrent creates
	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < tasksPerGoroutine; j++ {
				task := &entities.Task{
					Name:   fmt.Sprintf("Worker%d-Task%d", workerID, j),
					Status: (workerID + j) % 2,
				}
				if err := store.Create(task); err != nil {
					t.Errorf("Failed to create task: %v", err)
				} else {
					atomic.AddInt64(&successCount, 1)
				}
			}
		}(i)
	}

	wg.Wait()

	expectedCount := int64(numGoroutines * tasksPerGoroutine)
	if successCount != expectedCount {
		t.Errorf("Expected %d successful creates, got %d", expectedCount, successCount)
	}

	// Verify all tasks are retrievable
	allTasks := store.GetAll()
	if len(allTasks) != int(expectedCount) {
		t.Errorf("Expected %d tasks in store, got %d", expectedCount, len(allTasks))
	}
}

func TestShardStore_ConcurrentMixedOperations(t *testing.T) {
	store := NewShardStore(8)
	numWorkers := 50
	operationsPerWorker := 100
	var wg sync.WaitGroup

	// Pre-populate with some tasks
	initialTasks := make([]*entities.Task, 100)
	for i := 0; i < 100; i++ {
		task := &entities.Task{
			Name:   fmt.Sprintf("Initial Task %d", i),
			Status: i % 2,
		}
		store.Create(task)
		initialTasks[i] = task
	}

	wg.Add(numWorkers)
	for i := 0; i < numWorkers; i++ {
		go func(workerID int) {
			defer wg.Done()
			for j := 0; j < operationsPerWorker; j++ {
				switch j % 4 {
				case 0: // Create
					task := &entities.Task{
						Name:   fmt.Sprintf("Worker%d-Task%d", workerID, j),
						Status: j % 2,
					}
					store.Create(task)

				case 1: // Read
					if len(initialTasks) > 0 {
						taskToRead := initialTasks[j%len(initialTasks)]
						store.GetByID(taskToRead.ID)
					}

				case 2: // Update
					if len(initialTasks) > 0 {
						taskToUpdate := initialTasks[j%len(initialTasks)]
						updatedTask := &entities.Task{
							Name:   fmt.Sprintf("Updated by Worker%d", workerID),
							Status: 1,
						}
						store.Update(taskToUpdate.ID, updatedTask)
					}

				case 3: // GetAll
					store.GetAll()
				}
			}
		}(i)
	}

	wg.Wait()

	// Verify store is still functional
	testTask := &entities.Task{Name: "Post-concurrency test", Status: 0}
	err := store.Create(testTask)
	if err != nil {
		t.Errorf("Store not functional after concurrent operations: %v", err)
	}
}

func TestShardStore_LoadBalancing(t *testing.T) {
	store := NewShardStore(4)

	// Create many tasks and check they're distributed across shards
	numTasks := 1000
	for i := 0; i < numTasks; i++ {
		task := &entities.Task{
			Name:   fmt.Sprintf("Load Test Task %d", i),
			Status: i % 2,
		}
		err := store.Create(task)
		if err != nil {
			t.Fatalf("Failed to create task %d: %v", i, err)
		}
	}

	stats := store.GetShardStats()
	tasksPerShard := stats["tasksPerShard"].([]int)

	// Check that tasks are reasonably distributed
	totalTasks := 0
	for i, count := range tasksPerShard {
		totalTasks += count
		t.Logf("Shard %d: %d tasks", i, count)
	}

	if totalTasks != numTasks {
		t.Errorf("Expected %d total tasks, got %d", numTasks, totalTasks)
	}

	// No shard should be completely empty with 1000 tasks
	for i, count := range tasksPerShard {
		if count == 0 {
			t.Errorf("Shard %d is empty, load balancing may be poor", i)
		}
	}
}

func TestShardStore_MemoryEfficiency(t *testing.T) {
	store := NewShardStore(8)

	// Create and delete many tasks to test memory cleanup
	for cycle := 0; cycle < 10; cycle++ {
		tasks := make([]*entities.Task, 100)

		// Create tasks
		for i := 0; i < 100; i++ {
			task := &entities.Task{
				Name:   fmt.Sprintf("Cycle%d-Task%d", cycle, i),
				Status: i % 2,
			}
			err := store.Create(task)
			if err != nil {
				t.Fatalf("Failed to create task: %v", err)
			}
			tasks[i] = task
		}

		// Delete all tasks
		for _, task := range tasks {
			err := store.Delete(task.ID)
			if err != nil {
				t.Fatalf("Failed to delete task: %v", err)
			}
		}

		// Verify store is empty
		allTasks := store.GetAll()
		if len(allTasks) != 0 {
			t.Errorf("Expected empty store after cycle %d, got %d tasks", cycle, len(allTasks))
		}
	}
}

func TestShardStore_IDGeneration(t *testing.T) {
	store := NewShardStore(4)
	numTasks := 1000
	idsSeen := make(map[int]bool)

	for i := 0; i < numTasks; i++ {
		task := &entities.Task{
			Name:   fmt.Sprintf("ID Test Task %d", i),
			Status: i % 2,
		}
		err := store.Create(task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}

		// Check ID uniqueness
		if idsSeen[task.ID] {
			t.Errorf("Duplicate ID generated: %d", task.ID)
		}
		idsSeen[task.ID] = true

		// Check ID is positive
		if task.ID <= 0 {
			t.Errorf("Non-positive ID generated: %d", task.ID)
		}
	}

	// Verify we have exactly the number of unique IDs we expect
	if len(idsSeen) != numTasks {
		t.Errorf("Expected %d unique IDs, got %d", numTasks, len(idsSeen))
	}
}

func TestShardStore_EdgeCases(t *testing.T) {
	store := NewShardStore(4)

	// Test with nil task
	err := store.Create(nil)
	if err == nil {
		t.Error("Expected error when creating nil task")
	}
	if err.Error() != "task cannot be nil" {
		t.Errorf("Expected error message 'task cannot be nil', got '%s'", err.Error())
	}

	// Test empty name task
	emptyTask := &entities.Task{Name: "", Status: 0}
	err = store.Create(emptyTask)
	if err != nil {
		t.Errorf("Should allow empty name task at storage level: %v", err)
	}

	// Test extreme status values
	extremeTask := &entities.Task{Name: "Extreme", Status: 999}
	err = store.Create(extremeTask)
	if err != nil {
		t.Errorf("Should allow extreme status values at storage level: %v", err)
	}

	// Test very long name
	longName := string(make([]byte, 10000))
	for i := range longName {
		longName = longName[:i] + "a" + longName[i+1:]
	}
	longTask := &entities.Task{Name: longName, Status: 0}
	err = store.Create(longTask)
	if err != nil {
		t.Errorf("Should handle long names at storage level: %v", err)
	}
}

func TestShardStore_Performance(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping performance test in short mode")
	}

	store := NewShardStore(32)
	numOperations := 10000

	// Benchmark creates
	start := time.Now()
	tasks := make([]*entities.Task, numOperations)
	for i := 0; i < numOperations; i++ {
		task := &entities.Task{
			Name:   fmt.Sprintf("Perf Task %d", i),
			Status: i % 2,
		}
		err := store.Create(task)
		if err != nil {
			t.Fatalf("Failed to create task: %v", err)
		}
		tasks[i] = task
	}
	createDuration := time.Since(start)

	// Benchmark reads
	start = time.Now()
	for i := 0; i < numOperations; i++ {
		_, err := store.GetByID(tasks[i].ID)
		if err != nil {
			t.Fatalf("Failed to read task: %v", err)
		}
	}
	readDuration := time.Since(start)

	// Benchmark updates
	start = time.Now()
	for i := 0; i < numOperations; i++ {
		updatedTask := &entities.Task{
			Name:   fmt.Sprintf("Updated Task %d", i),
			Status: 1,
		}
		err := store.Update(tasks[i].ID, updatedTask)
		if err != nil {
			t.Fatalf("Failed to update task: %v", err)
		}
	}
	updateDuration := time.Since(start)

	t.Logf("Performance results for %d operations:", numOperations)
	t.Logf("Creates: %v (%.2f ops/sec)", createDuration, float64(numOperations)/createDuration.Seconds())
	t.Logf("Reads: %v (%.2f ops/sec)", readDuration, float64(numOperations)/readDuration.Seconds())
	t.Logf("Updates: %v (%.2f ops/sec)", updateDuration, float64(numOperations)/updateDuration.Seconds())

	// Sanity check - operations should complete in reasonable time
	maxDurationPerOp := time.Millisecond
	if createDuration/time.Duration(numOperations) > maxDurationPerOp {
		t.Errorf("Create operations too slow: %v per operation", createDuration/time.Duration(numOperations))
	}
}

func TestShardStore_StatsAccuracy(t *testing.T) {
	store := NewShardStore(4)

	// Initial stats
	stats := store.GetShardStats()
	if stats["numShards"] != 4 {
		t.Errorf("Expected 4 shards, got %v", stats["numShards"])
	}
	if stats["totalTasks"] != 0 {
		t.Errorf("Expected 0 total tasks, got %v", stats["totalTasks"])
	}

	// Create tasks and verify stats update
	numTasks := 100
	for i := 0; i < numTasks; i++ {
		task := &entities.Task{
			Name:   fmt.Sprintf("Stats Task %d", i),
			Status: i % 2,
		}
		store.Create(task)
	}

	stats = store.GetShardStats()
	if stats["totalTasks"] != numTasks {
		t.Errorf("Expected %d total tasks, got %v", numTasks, stats["totalTasks"])
	}

	tasksPerShard := stats["tasksPerShard"].([]int)
	if len(tasksPerShard) != 4 {
		t.Errorf("Expected 4 shard counts, got %d", len(tasksPerShard))
	}

	// Verify sum of tasks per shard equals total
	sum := 0
	for _, count := range tasksPerShard {
		sum += count
	}
	if sum != numTasks {
		t.Errorf("Sum of tasks per shard (%d) doesn't match total tasks (%d)", sum, numTasks)
	}
}

func TestShardStoreGopool_EdgeCases(t *testing.T) {
	store := NewShardStoreGopool(4)
	defer store.Close()

	// Test with nil task
	err := store.Create(nil)
	if err == nil {
		t.Error("Expected error when creating nil task")
	}
	if err.Error() != "task cannot be nil" {
		t.Errorf("Expected error message 'task cannot be nil', got '%s'", err.Error())
	}

	// Test empty name task
	emptyTask := &entities.Task{Name: "", Status: 0}
	err = store.Create(emptyTask)
	if err != nil {
		t.Errorf("Should allow empty name task at storage level: %v", err)
	}

	// Test extreme status values
	extremeTask := &entities.Task{Name: "Extreme", Status: 999}
	err = store.Create(extremeTask)
	if err != nil {
		t.Errorf("Should allow extreme status values at storage level: %v", err)
	}
}

func BenchmarkShardStore_Create(b *testing.B) {
	store := NewShardStore(32)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		task := &entities.Task{
			Name:   fmt.Sprintf("Benchmark Task %d", i),
			Status: i % 2,
		}
		store.Create(task)
	}
}

func BenchmarkShardStore_Read(b *testing.B) {
	store := NewShardStore(32)

	// Pre-populate
	tasks := make([]*entities.Task, 1000)
	for i := 0; i < 1000; i++ {
		task := &entities.Task{
			Name:   fmt.Sprintf("Pre-populate Task %d", i),
			Status: i % 2,
		}
		store.Create(task)
		tasks[i] = task
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		store.GetByID(tasks[i%1000].ID)
	}
}

func BenchmarkShardStore_ConcurrentReads(b *testing.B) {
	store := NewShardStore(32)

	// Pre-populate
	tasks := make([]*entities.Task, 1000)
	for i := 0; i < 1000; i++ {
		task := &entities.Task{
			Name:   fmt.Sprintf("Concurrent Read Task %d", i),
			Status: i % 2,
		}
		store.Create(task)
		tasks[i] = task
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			store.GetByID(tasks[i%1000].ID)
			i++
		}
	})
}
