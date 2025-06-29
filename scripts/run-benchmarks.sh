#!/bin/bash
set -e

echo "🚀 Running Storage Benchmarks..."

OUTPUT_DIR="output/benchmarks"
TIMESTAMP=$(date +%Y%m%d-%H%M%S)

mkdir -p "$OUTPUT_DIR"

# Function to run specific benchmark
run_benchmark() {
    local name=$1
    local pattern=$2
    local output_file="$OUTPUT_DIR/${name}-${TIMESTAMP}.txt"
    
    echo "📊 Running $name benchmarks..."
    go test -bench="$pattern" -benchmem -timeout=30m ./benchmarks/ | tee "$output_file"
    echo "✅ Results saved to $output_file"
    echo ""
}

# Run individual storage benchmarks
echo "🏪 Running individual storage benchmarks..."
run_benchmark "shard" ".*ShardStore.*"
run_benchmark "memory" ".*MemoryStore.*" 
run_benchmark "bigcache" ".*BigCacheStore.*"
run_benchmark "channel" ".*ChannelStore.*"

# Run pattern-specific benchmarks
echo "📈 Running pattern-specific benchmarks..."
run_benchmark "read-zipf" "BenchmarkReadZipf"
run_benchmark "write-zipf" "BenchmarkWriteZipf"
run_benchmark "distributed-read" "BenchmarkDistributedRead"
run_benchmark "distributed-write" "BenchmarkDistributedWrite"
run_benchmark "distributed-mixed" "BenchmarkDistributedMixed"

# Run comparison benchmarks
echo "⚡ Running comparison benchmarks..."
go test -bench="BenchmarkReadZipf|BenchmarkWriteZipf" -benchmem -timeout=30m ./benchmarks/ | tee "$OUTPUT_DIR/comparison-${TIMESTAMP}.txt"

echo ""
echo "🎉 All benchmarks completed!"
echo "📁 Results available in: $OUTPUT_DIR/"
echo ""
echo "Quick comparison (Read/Write Zipf):"
echo "===================================="
grep -h "BenchmarkReadZipf\|BenchmarkWriteZipf" "$OUTPUT_DIR/comparison-${TIMESTAMP}.txt" | sort
echo ""
echo "💡 Tip: Use 'make bench-shard' to run only ShardStore benchmarks"
echo "💡 Tip: Use 'make bench-compare' for cross-storage comparison"