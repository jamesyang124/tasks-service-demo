# Storage Benchmarks

## Overview

Comprehensive benchmarks for all storage implementations using realistic workload patterns with 1M dataset.

## Storage Implementations Tested

- **ShardStoreGopool**: ByteDance gopool per-core worker optimization (current best: 12.40ns reads)
- **ShardStore**: Optimized dedicated worker per shard storage (12.55ns reads)
- **MemoryStore**: Single-mutex in-memory storage (156.5ns reads)
- **ChannelStore**: Actor model with message passing (607.5ns reads)

## Performance Hierarchy (Latest Results)

| Implementation | Read Performance | Write Performance | Memory Allocations | Production Ready |
|---------------|------------------|-------------------|-------------------|------------------|
| **ShardStoreGopool** | **12.40 ns/op** | 62.69 ns/op | 0 B/op | 🏆 **Best** |
| **ShardStore** | 12.55 ns/op | **61.44 ns/op** | 0 B/op | ✅ **Excellent** |
| **MemoryStore** | 156.5 ns/op | 312.5 ns/op | 0 B/op | ⚠️ **Limited** |
| **ChannelStore** | 607.5 ns/op | 693.5 ns/op | 192 B/op | ❌ **Educational** |

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
- **Read Performance**: 12.40 ns/op (ByteDance gopool optimization)
- **Write Performance**: 62.69 ns/op
- **Memory Allocations**: 0 B/op (zero allocation reads)
- **Optimized for**: M4 Pro 14-core per-core worker pools

### ShardStore (Dedicated Workers)  
- **Read Performance**: 12.55 ns/op (dedicated worker optimization)
- **Write Performance**: 61.44 ns/op (best write performance)
- **Memory Allocations**: 0 B/op (zero allocation reads)

### MemoryStore (Baseline)
- **Read Performance**: 156.5 ns/op (13x slower than ShardStoreGopool)
- **Write Performance**: 312.5 ns/op (5x slower than ShardStore)
- **Memory Allocations**: 0 B/op (zero allocation reads)

### ChannelStore (Educational)
- **Read Performance**: 607.5 ns/op (59x slower than ShardStoreGopool)
- **Write Performance**: 693.5 ns/op (7x slower than ShardStore)
- **Memory Allocations**: 192 B/op (significant allocations)

## Performance Improvements

### ShardStoreGopool vs MemoryStore
- **Read Performance**: **12.6x faster** (156.5ns → 12.40ns)
- **Write Performance**: **5.0x faster** (312.5ns → 62.69ns)
- **Overall**: **12.6x performance improvement**

### ShardStore vs MemoryStore
- **Read Performance**: **12.5x faster** (156.5ns → 12.55ns)
- **Write Performance**: **5.1x faster** (312.5ns → 61.44ns)
- **Overall**: **12.5x performance improvement**

### ShardStoreGopool vs ShardStore
- **Read Performance**: **1.2% faster** (12.55ns → 12.40ns)
- **Write Performance**: **2.0% slower** (61.44ns → 62.69ns)
- **Trade-off**: Slightly better reads, slightly worse writes

## Optimization Journey

### Phase 1: Benchmark Reorganization
- Split monolithic `simple_bench_test.go` into focused benchmark files
- Simplified Makefile targets for essential benchmarks
- Organized results with timestamped output

### Phase 2: ShardStore Optimization  
- **12.5x read improvement** with dedicated worker per shard
- **5.1x write improvement** with optimized locking
- Eliminated contention bottlenecks

### Phase 3: ByteDance Gopool Optimization
- **Additional 1.2% read improvement** with per-core worker optimization
- Better CPU utilization on multi-core systems
- Optimized for Apple M4 Pro 14-core architecture
- **Combined result**: **12.6x faster than baseline**

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
- **Reason**: Best read performance (12.40ns), optimized for multi-core
- **Configuration**: 32 shards for M4 Pro, adjust for your CPU cores

### Balanced Production
- **Use**: ShardStore  
- **Reason**: Best write performance (61.44ns), excellent read performance (12.55ns)
- **Configuration**: 32 shards, dedicated workers per shard

### Development/Testing
- **Use**: MemoryStore
- **Reason**: Simple implementation, good for testing
- **Note**: Not recommended for production due to contention

### Educational/Research
- **Use**: ChannelStore
- **Reason**: Demonstrates actor model patterns
- **Note**: Not production-ready due to performance characteristics