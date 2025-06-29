# Storage Performance Comparison

## Overview

Quick reference guide for storage implementation performance characteristics. For detailed K6 load testing results and comprehensive performance analysis, see [K6_PERFORMANCE_REPORT.md](./K6_PERFORMANCE_REPORT.md).

## Performance Summary

| Storage Implementation | Max RPS | Max VUs | Avg Response | P95 Response | Error Rate | Production Ready |
|------------------------|---------|---------|--------------|--------------|------------|------------------|
| **ShardStoreGopool** | **2000+** | **1000+** | **<1ms** | **<5ms** | **0%** | 🏆 **Best** |
| **ShardStore** | **1800+** | **500** | **<2ms** | **<10ms** | **0%** | ✅ **Excellent** |
| **MemoryStore** | **161** | **100** | **1ms** | **3ms** | **100%*** | ❌ **Development Only** |

*MemoryStore fails under high concurrency testing

## Implementations Not Included in K6 Testing

### ❌ ChannelStore (Actor Model)
**Benchmark Results**:
- **Read Performance**: 666.1 ns/op (57x slower than ShardStore)
- **Write Performance**: 603.1 ns/op (58x slower than ShardStore)
- **Memory Overhead**: 192-247 B/op vs 0 B/op for optimized stores

**Why Excluded from K6 API Testing**:
1. **Prohibitive Performance Gap**: 57-58x slower than production candidates
2. **Would Fail Basic Load Tests**: At 666ns/op, estimated max ~150 RPS vs 2000+ RPS needed
3. **High Memory Allocation**: 3-5 allocs per operation vs 0 for optimized stores
4. **Channel Message Passing Overhead**: Actor model serialization cost too high
5. **Educational Purpose Only**: Demonstrates patterns but unsuitable for production

**Conclusion**: ChannelStore's benchmark performance (666ns reads) indicated it would fail even basic K6 load testing, so it was excluded to focus testing resources on viable production candidates.

## Quick Implementation Guide

### 🏆 ShardStoreGopool (Recommended)
```go
// Optimal for production systems
store := storage.NewShardStoreGopool(32) // Power-of-2 sharding
```
**Best for**: High-traffic APIs, enterprise production, real-time applications

### 🥈 ShardStore (Alternative)
```go
// Good balance of performance and simplicity  
store := storage.NewShardStore(32)
```
**Best for**: Moderate-traffic systems, simpler deployment requirements

### ❌ MemoryStore (Development Only)
```go
// Simple but not production-ready
store := storage.NewMemoryStore()
```
**Best for**: Development, testing, learning Go patterns

## Benchmark Results

### Core Performance Metrics
- **ShardStoreGopool**: 11.54ns reads (91% improvement from baseline)
- **ShardStore**: 12.42ns reads (90% improvement from baseline)  
- **MemoryStore**: 130ns reads (baseline)

For complete optimization journey details, see [OPTIMIZATION_DECISIONS.md](./OPTIMIZATION_DECISIONS.md).

## Testing & Validation

### Quick Performance Test
```bash
# Run benchmark comparison
make bench-compare

# Run K6 load tests
make k6-stress
make k6-read
```

### Production Validation
```bash
# Full comparative testing
./scripts/run-comparative-tests.sh
```

For comprehensive K6 test suite documentation, see [k6/README.md](../k6/README.md).

## Reference Links

- **Complete K6 Testing Guide**: [k6/README.md](../k6/README.md)
- **Detailed Performance Report**: [K6_PERFORMANCE_REPORT.md](./K6_PERFORMANCE_REPORT.md)  
- **Optimization Decision Log**: [OPTIMIZATION_DECISIONS.md](./OPTIMIZATION_DECISIONS.md)

---

**ShardStoreGopool achieved 91% performance improvement** and is production-ready for high-traffic applications requiring 2000+ RPS throughput with sub-millisecond latency.