# Storage Performance Comparison

Quick reference guide for storage implementation performance characteristics. For detailed benchmark results and comprehensive performance analysis, see the [benchmarks README](../benchmarks/README.md).

## Performance Summary

| Implementation | Read Performance | Write Performance | Memory Allocations | Production Ready |
|---------------|------------------|-------------------|-------------------|------------------|
| **XSyncStore** | **1.5 ns/op** | **18.0 ns/op** | 0-48 B/op | 🏆 **Best** |
| **ShardStoreGopool** | 12.2 ns/op | 60.9 ns/op | 0-104 B/op | ✅ **Excellent** |
| **ShardStore** | 14.5 ns/op | 36.4 ns/op | 0-32 B/op | ✅ **Excellent** |
| **MemoryStore** | 159.8 ns/op | 220.7 ns/op | 0-32 B/op | ⚠️ **Limited** |
| **ChannelStore** | 607.5 ns/op | 693.5 ns/op | 192 B/op | ❌ **Educational** |

## Performance Improvements

### XSyncStore vs MemoryStore (Lock-Free Revolution)
- **Read Performance**: **106x faster** (159.8ns → 1.5ns)
- **Write Performance**: **12.2x faster** (220.7ns → 18.0ns)
- **Overall**: **106x performance improvement**

### XSyncStore vs ShardStoreGopool
- **Read Performance**: **8.1x faster** (12.2ns → 1.5ns)
- **Write Performance**: **3.4x faster** (60.9ns → 18.0ns)
- **Overall**: **8.1x performance improvement**

### XSyncStore vs ShardStore
- **Read Performance**: **9.6x faster** (14.5ns → 1.5ns)
- **Write Performance**: **2.0x faster** (36.4ns → 18.0ns)
- **Overall**: **9.6x performance improvement**

### Legacy Performance Comparisons

**ShardStoreGopool vs MemoryStore**:
- **Read Performance**: **13.1x faster** (159.8ns → 12.2ns)
- **Write Performance**: **3.6x faster** (220.7ns → 60.9ns)

**ShardStore vs MemoryStore**:
- **Read Performance**: **11.0x faster** (159.8ns → 14.5ns)
- **Write Performance**: **6.1x faster** (220.7ns → 36.4ns)

## Implementations Not Included in Load Testing

### ChannelStore Exclusion

**Why Excluded from Load Testing**:
1. **Performance Floor**: 607.5ns reads = ~1,600 RPS theoretical maximum
2. **Memory Overhead**: 192 B/op allocations create GC pressure
3. **Educational Purpose**: Designed to demonstrate actor model patterns
4. **Resource Efficiency**: Testing time better spent on production candidates
5. **Clear Outcome**: Benchmarks already show it's 49x slower than ShardStoreGopool

**Conclusion**: ChannelStore's benchmark performance (607.5ns reads) indicated it would fail even basic load testing, so it was excluded to focus testing resources on viable production candidates.

## Storage Implementation Details

### XSyncStore (Recommended - Default)
- **Architecture**: Lock-free concurrent map using xsync.Map
- **Read Performance**: 1.5 ns/op (best - 106x faster than MemoryStore)
- **Write Performance**: 18.0 ns/op (best - 12.2x faster than MemoryStore)
- **Memory**: 0-48 B/op (minimal allocations)
- **Use Case**: All production systems, especially high-concurrency
- **Optimization**: Lock-free atomic operations, CAS loops, hazard pointers
- **Advantages**: No deadlocks, linear scalability, predictable performance

### ShardStoreGopool (High Performance Alternative)
- **Architecture**: Sharded storage with ByteDance gopool optimization
- **Read Performance**: 12.2 ns/op (8.1x slower than XSyncStore)
- **Write Performance**: 60.9 ns/op (3.4x slower than XSyncStore)
- **Memory**: 0-104 B/op (low allocations)
- **Use Case**: High-traffic production systems, write-heavy workloads
- **Optimization**: Consistent hashing with ByteDance gopool for reduced contention

### ShardStore (Balanced Alternative)
- **Architecture**: Sharded storage with dedicated workers
- **Read Performance**: 14.5 ns/op (9.6x slower than XSyncStore)
- **Write Performance**: 36.4 ns/op (2.0x slower than XSyncStore)
- **Memory**: 0-32 B/op (minimal allocations)
- **Use Case**: Production systems with balanced read/write workloads
- **Optimization**: Dedicated worker per shard for reduced lock contention

### MemoryStore (Development)
- **Architecture**: Single mutex in-memory storage
- **Read Performance**: 159.8 ns/op (106x slower than XSyncStore)
- **Write Performance**: 220.7 ns/op (12.2x slower than XSyncStore)
- **Memory**: 0-32 B/op (minimal allocations)
- **Use Case**: Development and testing environments
- **Limitation**: Global lock creates contention bottleneck

### ChannelStore (Educational)
- **Architecture**: Actor model with message passing
- **Read Performance**: 607.5 ns/op (405x slower than XSyncStore)
- **Write Performance**: 693.5 ns/op (38.5x slower than XSyncStore)
- **Memory**: 192 B/op (significant allocations)
- **Use Case**: Educational demonstration of actor model patterns
- **Limitation**: Channel overhead and single worker bottleneck

## Performance Testing

### Run Benchmarks
```bash
# Run all benchmarks (1M dataset)
make bench

# Direct Go command
go test -bench=. -benchmem -timeout=30m ./benchmarks/
```

### Specific Benchmark Patterns
```bash
# Compare read/write performance across all stores
go test -bench="BenchmarkReadZipf|BenchmarkWriteZipf" -benchmem ./benchmarks/

# Test specific storage implementation
go test -bench=".*XSyncStore.*" -benchmem ./benchmarks/
go test -bench=".*ShardStore.*" -benchmem ./benchmarks/
go test -bench=".*MemoryStore.*" -benchmem ./benchmarks/

# High contention testing
go test -bench=".*HighContention.*" -benchmem ./benchmarks/
```

## Documentation

- **Complete Benchmark Guide**: [benchmarks/README.md](../benchmarks/README.md)
- **Optimization Journey**: [OPTIMIZATION_DECISIONS.md](./OPTIMIZATION_DECISIONS.md)
- **Performance Analysis**: This document