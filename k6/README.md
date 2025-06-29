# K6 Performance Test Suite

## Overview

Comprehensive load testing suite for the Task API service, designed to validate the performance optimizations achieved through our storage architecture evolution from MemoryStore to ShardStoreGopool.

## Test Suite Architecture

### 🧪 Test Files

| Test File | Purpose | Target | Max VUs | Duration |
|-----------|---------|--------|---------|----------|
| `test.js` | **Baseline CRUD** | Mixed operations validation | 100 | 40s |
| `stress-test.js` | **High Concurrency** | System limits & scaling | 1000 | 90s |
| `read-heavy-test.js` | **Read Optimization** | Hot key performance | 800 | 50s |
| `comparative-test.js` | **Storage Comparison** | Cross-implementation testing | 500 | 50s |
| `load-limit-test.js` | **Peak Performance** | Maximum throughput discovery | 3000 | 50s |

## Test Descriptions

### 📊 **1. Baseline Test (`test.js`)**
**Purpose**: Establish baseline performance with mixed CRUD operations

**Pattern**:
- Mixed GET, POST, PUT, DELETE operations
- 1-second sleep between iterations
- Gradual ramp-up to 100 VUs

**Key Metrics**:
- Basic API functionality validation
- Response time consistency
- Error rate under normal load

**Usage**:
```bash
make k6           # Original test suite
# or
./scripts/run-k6-tests.sh original
```

**Expected Results**:
- **MemoryStore**: ~161 RPS, stable under 100 VUs
- **ShardStore**: ~300 RPS, excellent stability
- **ShardStoreGopool**: ~350 RPS, optimal performance

---

### 🚀 **2. Stress Test (`stress-test.js`)**
**Purpose**: Discover system limits under extreme concurrency

**Pattern**:
- 80% reads (GetAll + GetByID), 20% writes
- No artificial delays (maximum throughput)
- Aggressive ramp: 100 → 500 → 1000 VUs

**Key Metrics**:
- Peak RPS capability
- Failure threshold identification
- System stability under pressure

**Usage**:
```bash
make k6-stress
# or
./scripts/run-k6-tests.sh stress
```

**Expected Results**:
- **MemoryStore**: Fails at ~200 VUs
- **ShardStore**: Stable to 500 VUs, 1800+ RPS
- **ShardStoreGopool**: Stable to 1000+ VUs, 2000+ RPS

---

### 📖 **3. Read-Heavy Test (`read-heavy-test.js`)**
**Purpose**: Validate read optimization benefits with realistic data patterns

**Pattern**:
- 95% reads (70% GetByID + 25% GetAll), 5% writes
- 10,000 task dataset for realistic scale
- Zipf distribution (80/20 hot keys)

**Key Metrics**:
- Read operation performance
- Hot key caching effectiveness
- Bulk operation scaling

**Usage**:
```bash
make k6-read
# or
./scripts/run-k6-tests.sh read
```

**Expected Results**:
- **MemoryStore**: Fails under read load
- **ShardStore**: <20ms p95 response times
- **ShardStoreGopool**: <15ms p95 response times, optimal hot key performance

---

### ⚖️ **4. Comparative Test (`comparative-test.js`)**
**Purpose**: Direct performance comparison across storage implementations

**Pattern**:
- Identical test conditions for all storage types
- 80% reads, 20% writes
- Configurable via `STORAGE_TYPE` environment variable

**Key Metrics**:
- Side-by-side performance comparison
- Storage-specific bottleneck identification
- Architecture decision validation

**Usage**:
```bash
make k6-compare
# or
./scripts/run-comparative-tests.sh
```

**Configuration**:
```bash
# Test specific storage
STORAGE_TYPE=memory make k6-compare
STORAGE_TYPE=shard make k6-compare  
STORAGE_TYPE=gopool make k6-compare
```

---

### 🔥 **5. Load Limit Test (`load-limit-test.js`)**
**Purpose**: Discover absolute performance limits

**Pattern**:
- Aggressive scaling: 500 → 1000 → 2000 → 3000 VUs
- 90% reads, 10% writes
- Push until failure or resource exhaustion

**Key Metrics**:
- Maximum sustainable RPS
- Breaking point identification
- Resource utilization limits

**Usage**:
```bash
# Manual execution (not in make targets)
docker run --rm -v "${PWD}/k6:/scripts" --network="host" grafana/k6:latest run /scripts/load-limit-test.js
```

## 📈 Performance Benchmarks & Expected Results

### Current Performance Hierarchy

| Storage | Baseline (CRUD) | Stress (RPS) | Read-Heavy (p95) | Max VUs |
|---------|----------------|---------------|------------------|---------|
| **ShardStoreGopool** | 350 RPS | **2000+ RPS** | **<15ms** | **1000+** |
| **ShardStore** | 300 RPS | **1800+ RPS** | **<20ms** | **500** |
| **MemoryStore** | 161 RPS | **Fails** | **Fails** | **100** |

### Test Environment Requirements

**Hardware**: Apple M4 Pro (14 cores) or equivalent
**Memory**: 16GB+ recommended for high VU tests
**Network**: Local testing preferred for consistent results

## 🚦 Test Execution Guide

### Quick Start
```bash
# Run all test suites
make k6

# Individual tests
make k6-stress    # High concurrency test
make k6-read      # Read optimization test
make k6-compare   # Storage comparison
```

### Advanced Usage
```bash
# Custom storage testing
STORAGE_TYPE=shard ./scripts/run-k6-tests.sh stress

# Custom parameters
SHARD_COUNT=64 STORAGE_TYPE=gopool make k6-stress
```

### Results Location
- **JSON Reports**: `./output/`
- **HTML Reports**: `./output/*.html`
- **Console Output**: Real-time metrics during execution

## 🎯 Test Selection Guide

### Choose Your Test Based On:

**🔍 API Validation**
→ Use `test.js` (Baseline CRUD)

**🚀 Performance Limits**  
→ Use `stress-test.js` (High Concurrency)

**📊 Read Performance**
→ Use `read-heavy-test.js` (Read Optimization)

**⚖️ Architecture Decisions**
→ Use `comparative-test.js` (Storage Comparison)

**🔥 Absolute Limits**
→ Use `load-limit-test.js` (Peak Performance)

## 📊 Understanding Results

### Key Metrics Interpretation

**RPS (Requests Per Second)**
- Primary throughput indicator
- Higher = better performance
- Compare across storage implementations

**Response Times**
- `avg`: Average response time
- `p(95)`: 95th percentile (SLA critical)
- `p(99)`: 99th percentile (tail latency)

**Error Rates**
- `http_req_failed`: Failed request percentage
- Target: <1% for production readiness
- 100% = system overload

**VU Scaling**
- Maximum concurrent users supported
- Indicates concurrency capabilities
- Breaking point identification

### Performance Targets

**Production Ready Thresholds**:
- **RPS**: >1000 sustained
- **P95 Response**: <50ms
- **Error Rate**: <1%
- **Max VUs**: >300

**Optimization Success**:
- **ShardStoreGopool**: ✅ Exceeds all targets
- **ShardStore**: ✅ Meets all targets  
- **MemoryStore**: ❌ Fails under load

## 🔧 Troubleshooting

### Common Issues

**"Connection Refused"**
- Server not running or wrong port
- Check: `curl http://localhost:8080/tasks`

**"100% Error Rate"**
- System overload (expected for MemoryStore)
- Reduce VU count or increase duration

**"Docker Issues"**
- Rebuild k6 container: `docker-compose -f docker-compose.test.yml build`
- Check network connectivity: `--network="host"`

### Performance Debugging

**Low RPS**:
1. Verify storage type: Check logs for storage initialization
2. Check CPU utilization: `top` during test execution
3. Monitor memory usage: Ensure sufficient RAM

**High Error Rates**:
1. Gradual load increase: Reduce VU ramp speed
2. Resource monitoring: Check system limits
3. Configuration tuning: Adjust shard count

## 📋 Test Results Archive

### Historical Performance Data

**Optimization Journey**:
- **Phase 1**: MemoryStore baseline (161 RPS)
- **Phase 2**: ShardStore implementation (1800+ RPS)  
- **Phase 3**: ByteDance gopool (2000+ RPS)
- **Phase 4**: ShardUnit optimization (2000+ RPS, <15ms p95)

**Validation Results**:
- **91% performance improvement** (130ns → 11.54ns reads)
- **12x throughput increase** (161 → 2000+ RPS)
- **10x concurrency improvement** (100 → 1000+ VUs)
- **Zero production failures** under realistic load

---

*For detailed performance analysis and architectural decisions, see `/docs/PERFORMANCE_COMPARISON.md`*  
*For comprehensive K6 test results and validation, see `/docs/K6_PERFORMANCE_REPORT.md`*