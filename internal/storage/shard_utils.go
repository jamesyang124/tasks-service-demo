package storage

// Utility functions for ShardStore - used for monitoring, debugging, and benchmarking
// These functions are NOT needed for production operation

// GetShardStats returns statistics about shard distribution
func (s *ShardStore) GetShardStats() map[string]interface{} {
	stats := make(map[string]interface{})
	stats["numShards"] = s.numShards
	
	shardCounts := make([]int, s.numShards)
	totalTasks := 0
	
	// Collect stats from all shards
	for i, shard := range s.shards {
		shard.mu.RLock()
		count := len(shard.tasks)
		shard.mu.RUnlock()
		
		shardCounts[i] = count
		totalTasks += count
	}
	
	stats["totalTasks"] = totalTasks
	stats["tasksPerShard"] = shardCounts

	return stats
}

// GetShard returns a specific shard (useful for testing/debugging)
func (s *ShardStore) GetShard(index int) *MemoryStore {
	if index < 0 || index >= s.numShards {
		return nil
	}

	// Access shard directly (no mutex needed)
	return s.shards[index]
}