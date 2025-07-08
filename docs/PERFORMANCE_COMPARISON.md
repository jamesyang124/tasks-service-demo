# Storage Performance Comparison

Quick reference guide for storage implementation performance characteristics. For detailed benchmark results and comprehensive performance analysis, see the [benchmarks README](../benchmarks/README.md).

## Performance Summary

| Implementation | Read Performance | Write Performance | Memory Allocations | Production Ready |
|---------------|------------------|-------------------|-------------------|------------------|
| **ShardStoreGopool** | **12.40 ns/op** | 62.69 ns/op | 0 B/op | 🏆 **Best** |
| **ShardStore** | 12.55 ns/op | **61.44 ns/op** | 0 B/op | ✅ **Excellent** |
| **MemoryStore** | 156.5 ns/op | 312.5 ns/op | 0 B/op | ⚠️ **Limited** |
| **ChannelStore** | 607.5 ns/op | 693.5 ns/op | 192 B/op | ❌ **Educational** |

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

### ShardStoreGopool (Recommended)
- **Architecture**: Sharded storage with ByteDance gopool optimization
- **Read Performance**: 12.40 ns/op (best)
- **Write Performance**: 62.69 ns/op
- **Memory**: 0 B/op (zero allocation reads)
- **Use Case**: High-traffic production systems
- **Optimization**: Consistent hashing with ByteDance gopool for reduced contention

### ShardStore (Alternative)
- **Architecture**: Sharded storage with dedicated workers
- **Read Performance**: 12.55 ns/op
- **Write Performance**: 61.44 ns/op (best)
- **Memory**: 0 B/op (zero allocation reads)
- **Use Case**: Production systems with balanced read/write workloads
- **Optimization**: Dedicated worker per shard for reduced lock contention

### MemoryStore (Development)
- **Architecture**: Single mutex in-memory storage
- **Read Performance**: 156.5 ns/op (12.5x slower than ShardStore)
- **Write Performance**: 312.5 ns/op (5.1x slower than ShardStore)
- **Memory**: 0 B/op (zero allocation reads)
- **Use Case**: Development and testing environments
- **Limitation**: Global lock creates contention bottleneck

### ChannelStore (Educational)
- **Architecture**: Actor model with message passing
- **Read Performance**: 607.5 ns/op (49x slower than ShardStoreGopool)
- **Write Performance**: 693.5 ns/op (11x slower than ShardStore)
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
go test -bench=".*ShardStore.*" -benchmem ./benchmarks/
go test -bench=".*MemoryStore.*" -benchmem ./benchmarks/
```

## Documentation

- **Complete Benchmark Guide**: [benchmarks/README.md](../benchmarks/README.md)
- **Optimization Journey**: [OPTIMIZATION_DECISIONS.md](./OPTIMIZATION_DECISIONS.md)
- **Performance Analysis**: This document