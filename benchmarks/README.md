# Storage Benchmarks

## Overview

Comprehensive benchmarks for all storage implementations using realistic workload patterns with 1M dataset.

## Storage Implementations Tested

- **MemoryStore**: Single-mutex in-memory storage
- **ShardStore**: Optimized dedicated worker per shard storage 
- **ShardStoreGopool**: ByteDance gopool per-core worker optimization (current best: 13.6ns reads)
- **BigCacheStore**: Off-heap cache with zero GC overhead
- **ChannelStore**: Actor model with message passing

## Benchmark Patterns

### Zipf Distribution (Realistic Hot Keys)
- **Dataset**: 1,000,000 tasks
- **Hot Keys**: 200,000 keys (20%) receive 80% of traffic
- **Pattern**: Simulates real-world production workloads

### Distributed Access (Uniform)
- **Dataset**: 1,000,000 tasks  
- **Pattern**: Uniform access across entire dataset
- **Usage**: Worst-case performance testing

### Mixed Workloads
- **Read/Write Ratio**: 70% reads, 30% writes
- **Pattern**: Realistic application workload

## Running Benchmarks

### All Benchmarks
```bash
make bench
```

### Essential Benchmarks
```bash
make bench-shard      # ShardStore benchmarks (includes gopool comparison)
make bench-compare    # Compare read/write across all stores
```

### Direct Go Commands
```bash
# All benchmarks
go test -bench=. -benchmem -timeout=30m ./benchmarks/

# ShardStore vs ShardStoreGopool comparison
go test -bench=.*ShardStore.* -benchmem -timeout=30m ./benchmarks/

# Current vs Gopool optimization comparison
go test -bench="BenchmarkShardStore_Comparison|BenchmarkGetAll_Comparison" -benchmem -timeout=30m ./benchmarks/
```

## Current Performance Baseline

### ShardStoreGopool + ShardUnit (Current Best)
- **Read Performance**: 11.54 ns/op (ByteDance gopool + ShardUnit optimization)
- **Write Performance**: 60.37 ns/op
- **Bulk Operations (GetAll)**: 1.505ms/op (7.4% faster than dedicated workers)
- **Optimized for**: M4 Pro 14-core per-core worker pools + lightweight storage units

### ShardStore + ShardUnit (Dedicated Workers)  
- **Read Performance**: 12.42 ns/op (24.3% improvement from ShardUnit)
- **Write Performance**: 60.36 ns/op
- **Previous Performance (MemoryStore)**: 16.4ns reads, 61.3ns writes

### BigCacheStore
- **Read Performance**: 65.1 ns/op
- **Write Performance**: ~125 ns/op

### MemoryStore (Baseline)
- **Read Performance**: ~130 ns/op
- **Write Performance**: ~210 ns/op

## Optimization Summary

### Phase 1: Benchmark Reorganization
- Split monolithic `simple_bench_test.go` into focused benchmark files
- Simplified Makefile targets for essential benchmarks
- Organized results with timestamped output

### Phase 2: ByteDance Gopool Optimization  
- **16.9% faster reads** with per-core worker optimization
- **7.4% faster bulk operations** with better CPU utilization
- Optimized for Apple M4 Pro 14-core architecture

### Phase 3: ShardUnit Optimization
- **Additional 24.3% read improvement** by replacing MemoryStore with lightweight ShardUnit
- Eliminated unnecessary ID generation overhead per shard
- Cleaner API separation: sharding logic vs storage logic
- **Combined result**: **11.54ns reads** (fastest implementation)

## Benchmark Output Organization

Results are automatically saved to `output/benchmarks/` with timestamps for historical tracking.

## File Structure

```
benchmarks/
├── README.md                # This documentation
├── common.go                # Shared benchmark utilities
├── memory_bench_test.go     # MemoryStore benchmarks
├── shard_bench_test.go      # ShardStore benchmarks (dedicated workers)
├── shard_gopool_bench_test.go     # ShardStoreGopool benchmarks (ByteDance optimization)
├── bigcache_bench_test.go   # BigCacheStore benchmarks
└── channel_bench_test.go    # ChannelStore benchmarks
```

## Storage Package Organization

The storage implementations are now organized with consistent naming:

```
internal/storage/
├── store.go                    # Main Store interface
├── store_memory.go            # MemoryStore (naive implementation)
├── store_shard.go             # ShardStore (high-performance)
├── store_shard_gopool.go      # ShardStoreGopool (ByteDance optimization)
├── store_shard_unit.go        # ShardUnit (lightweight storage units)
├── store_shard_utils.go       # Shard utility functions
├── store_bigcache.go          # BigCacheStore (off-heap)
├── store_channel.go           # ChannelStore (actor model)
└── store_*.go                 # Additional implementations and tests
```

## Adding New Benchmarks

1. Create new `*_bench_test.go` file for your storage implementation
2. Import `benchmarks` package and use common utilities:
   - `PopulateStore()` for test data setup
   - `BenchmarkReadZipf()` for hot key read patterns
   - `BenchmarkWriteZipf()` for hot key write patterns
3. Update Makefile with new benchmark targets
4. Update this README with performance results

## Performance Testing Best Practices

- **Consistent Environment**: Run benchmarks on same hardware
- **Multiple Runs**: Use `-count=3` for statistical significance
- **Memory Profiling**: Include `-benchmem` for allocation analysis
- **Timeout**: Use `-timeout=30m` for large datasets
- **Baseline Comparison**: Always compare against previous results