# Storage Benchmarks

## Overview

Comprehensive benchmarks for all storage implementations using realistic workload patterns with 1M dataset.

## Storage Implementations Tested

- **ShardStoreGopool**: ByteDance gopool per-core worker optimization (current best: 12.19ns reads)
- **ShardStore**: Optimized dedicated worker per shard storage (13.81ns reads)
- **MemoryStore**: Single-mutex in-memory storage (158.0ns reads)
- **ChannelStore**: Actor model with message passing (718.8ns reads)

## Performance Hierarchy (Latest Results)

| Implementation | Read Performance | Write Performance | Memory Allocations | Production Ready |
|---------------|------------------|-------------------|-------------------|------------------|
| **ShardStoreGopool** | **12.19 ns/op** | 61.72 ns/op | 0 B/op | 🏆 **Best** |
| **ShardStore** | 13.81 ns/op | **60.75 ns/op** | 0 B/op | ✅ **Excellent** |
| **MemoryStore** | 158.0 ns/op | 295.1 ns/op | 0 B/op | ⚠️ **Limited** |
| **ChannelStore** | 718.8 ns/op | 452.8 ns/op | 192 B/op | ❌ **Educational** |

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

### Simple Commands
```bash
# Run all benchmarks (recommended)
make bench

# Direct Go command
go test -bench=. -benchmem -timeout=30m ./benchmarks/
```

### Specific Benchmark Patterns
```bash
# Compare read/write performance across all stores
go test -bench="BenchmarkReadZipf|BenchmarkWriteZipf" -benchmem ./benchmarks/

# Test specific storage implementation
go test -bench=".*ShardStore.*" -benchmem ./benchmarks/
go test -bench=".*MemoryStore.*" -benchmem ./benchmarks/
```

## Current Performance Results

### ShardStoreGopool (Current Best)
- **Read Performance**: 12.19 ns/op (ByteDance gopool optimization)
- **Write Performance**: 61.72 ns/op
- **Memory Allocations**: 0 B/op (zero allocation reads)
- **Optimized for**: M4 Pro 14-core per-core worker pools

### ShardStore (Dedicated Workers)  
- **Read Performance**: 13.81 ns/op (dedicated worker optimization)
- **Write Performance**: 60.75 ns/op (best write performance)
- **Memory Allocations**: 0 B/op (zero allocation reads)

### MemoryStore (Baseline)
- **Read Performance**: 158.0 ns/op (13x slower than ShardStoreGopool)
- **Write Performance**: 295.1 ns/op (5x slower than ShardStore)
- **Memory Allocations**: 0 B/op (zero allocation reads)

### ChannelStore (Educational)
- **Read Performance**: 718.8 ns/op (59x slower than ShardStoreGopool)
- **Write Performance**: 452.8 ns/op (7x slower than ShardStore)
- **Memory Allocations**: 192 B/op (significant allocations)

## Performance Improvements

### ShardStoreGopool vs MemoryStore
- **Read Performance**: **13x faster** (158.0ns → 12.19ns)
- **Write Performance**: **4.8x faster** (295.1ns → 61.72ns)
- **Overall**: **12.9x performance improvement**

### ShardStore vs MemoryStore
- **Read Performance**: **11.4x faster** (158.0ns → 13.81ns)
- **Write Performance**: **4.9x faster** (295.1ns → 60.75ns)
- **Overall**: **11.1x performance improvement**

### ShardStoreGopool vs ShardStore
- **Read Performance**: **13% faster** (13.81ns → 12.19ns)
- **Write Performance**: **1.6% slower** (60.75ns → 61.72ns)
- **Trade-off**: Slightly better reads, slightly worse writes

## Optimization Journey

### Phase 1: Benchmark Reorganization
- Split monolithic `simple_bench_test.go` into focused benchmark files
- Simplified Makefile targets for essential benchmarks
- Organized results with timestamped output

### Phase 2: ShardStore Optimization  
- **11.4x read improvement** with dedicated worker per shard
- **4.9x write improvement** with optimized locking
- Eliminated contention bottlenecks

### Phase 3: ByteDance Gopool Optimization
- **Additional 13% read improvement** with per-core worker optimization
- Better CPU utilization on multi-core systems
- Optimized for Apple M4 Pro 14-core architecture
- **Combined result**: **13x faster than baseline**

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
└── channel_bench_test.go    # ChannelStore benchmarks
```

## Storage Package Organization

The storage implementations are organized into logical subpackages:

```
internal/storage/
├── store.go                    # Main Store interface & singleton
├── naive/                      # Naive Memory Store
│   ├── memory.go              # Simple single-mutex implementation
│   └── memory_test.go         # Memory store tests
├── shard/                      # High-Performance Shard Store
│   ├── shard.go               # Optimized sharded storage
│   ├── shard_gopool.go        # ByteDance gopool optimization
│   ├── shard_unit.go          # Lightweight storage units
│   ├── shard_utils.go         # Utility functions
│   └── shard_test.go          # Comprehensive tests
└── channel/                    # Actor Model Store
    ├── channel_store.go       # Message passing implementation
    ├── channel_store_test.go  # Channel store tests
    └── pool_comparison_test.go # Performance benchmarks
```

### Usage in Benchmarks

```go
import (
    "tasks-service-demo/internal/storage/naive"
    "tasks-service-demo/internal/storage/shard"
    "tasks-service-demo/internal/storage/channel"
)

// Direct package usage
store := shard.NewShardStoreGopool(32)
defer store.Close()
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

## Production Recommendations

### High-Traffic Production
- **Use**: ShardStoreGopool
- **Reason**: Best read performance (12.19ns), optimized for multi-core
- **Configuration**: 32 shards for M4 Pro, adjust for your CPU cores

### Balanced Production
- **Use**: ShardStore  
- **Reason**: Best write performance (60.75ns), excellent read performance (13.81ns)
- **Configuration**: 32 shards, dedicated workers per shard

### Development/Testing
- **Use**: MemoryStore
- **Reason**: Simple implementation, good for testing
- **Note**: Not recommended for production due to contention

### Educational/Research
- **Use**: ChannelStore
- **Reason**: Demonstrates actor model patterns
- **Note**: Not production-ready due to performance characteristics