# Storage Benchmarks

## Overview

Comprehensive benchmarks for all storage implementations using realistic workload patterns with 1M dataset.

## Storage Implementations Tested

- **MemoryStore**: Single-mutex in-memory storage
- **ShardStore**: Optimized dedicated worker per shard storage (current best: 12.3ns reads)
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

### Specific Storage Type
```bash
make bench-shard      # ShardStore only
make bench-memory     # MemoryStore only  
make bench-bigcache   # BigCacheStore only
make bench-channel    # ChannelStore only
```

### Specific Patterns
```bash
make bench-read-zipf     # Hot key read patterns
make bench-write-zipf    # Hot key write patterns
make bench-distributed  # Uniform distribution patterns
```

### Direct Go Commands
```bash
# All benchmarks
go test -bench=. -benchmem -timeout=30m ./benchmarks/

# Specific storage
go test -bench=.*ShardStore.* -benchmem -timeout=30m ./benchmarks/

# Pattern comparison
go test -bench="BenchmarkReadZipf|BenchmarkWriteZipf" -benchmem -timeout=30m ./benchmarks/
```

## Current Performance Baseline

### ShardStore (Optimized - Current Best)
- **Read Performance**: 12.3 ns/op (10.6x improvement over original)
- **Write Performance**: 35.8 ns/op (3.6x improvement over original)
- **vs BigCache**: 5.3x faster reads, 3.5x faster writes

### BigCacheStore
- **Read Performance**: 65.1 ns/op
- **Write Performance**: ~125 ns/op

### MemoryStore (Baseline)
- **Read Performance**: ~130 ns/op
- **Write Performance**: ~210 ns/op

## Benchmark Output Organization

Results are automatically saved to `output/benchmarks/` with timestamps for historical tracking.

## File Structure

```
benchmarks/
├── README.md                # This documentation
├── common.go                # Shared benchmark utilities
├── memory_bench_test.go     # MemoryStore benchmarks
├── shard_bench_test.go      # ShardStore benchmarks
├── bigcache_bench_test.go   # BigCacheStore benchmarks
└── channel_bench_test.go    # ChannelStore benchmarks
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