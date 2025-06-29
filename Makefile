.PHONY: build test run dev k6 bench bench-simple bench-all help

build:
	go build -o bin/tasks-service-demo ./internal/cmd

test:
	go test -cover ./...

run: build
	./bin/tasks-service-demo

dev:
	go run ./internal/cmd/main.go

k6:
	./scripts/run-k6-tests.sh

# Core benchmarks (1M dataset) - All storage implementations
bench:
	go test -bench=. -benchmem -timeout=30m ./benchmarks/

# Storage-specific benchmarks
bench-shard:
	go test -bench=.*ShardStore.* -benchmem -timeout=30m ./benchmarks/

bench-memory:
	go test -bench=.*MemoryStore.* -benchmem -timeout=30m ./benchmarks/

bench-bigcache:
	go test -bench=.*BigCacheStore.* -benchmem -timeout=30m ./benchmarks/

bench-channel:
	go test -bench=.*ChannelStore.* -benchmem -timeout=30m ./benchmarks/

# Pattern-specific benchmarks
bench-read-zipf:
	go test -bench=BenchmarkReadZipf -benchmem -timeout=30m ./benchmarks/

bench-write-zipf:
	go test -bench=BenchmarkWriteZipf -benchmem -timeout=30m ./benchmarks/

bench-distributed-read:
	go test -bench=BenchmarkDistributedRead -benchmem -timeout=30m ./benchmarks/

bench-distributed-write:
	go test -bench=BenchmarkDistributedWrite -benchmem -timeout=30m ./benchmarks/

bench-distributed-mixed:
	go test -bench=BenchmarkDistributedMixed -benchmem -timeout=30m ./benchmarks/

# Comparison benchmarks
bench-compare:
	go test -bench="BenchmarkReadZipf|BenchmarkWriteZipf" -benchmem -timeout=30m ./benchmarks/

bench-all: bench
	@echo "All storage benchmarks completed"

# Save benchmark results
bench-save:
	@mkdir -p output/benchmarks
	go test -bench=. -benchmem ./benchmarks/ > output/benchmarks/full-$(shell date +%Y%m%d-%H%M%S).txt
	@echo "Benchmark results saved to output/benchmarks/"

help:
	@echo "build                   - Build binary"
	@echo "test                    - Run tests"
	@echo "run                     - Build and run"
	@echo "dev                     - Run in dev mode"
	@echo "k6                      - Run k6 load tests"
	@echo ""
	@echo "Storage Benchmarks (1M dataset):"
	@echo "bench                   - Run all storage benchmarks"
	@echo "bench-shard             - ShardStore benchmarks only"
	@echo "bench-memory            - MemoryStore benchmarks only"
	@echo "bench-bigcache          - BigCacheStore benchmarks only"
	@echo "bench-channel           - ChannelStore benchmarks only"
	@echo ""
	@echo "Pattern Benchmarks:"
	@echo "bench-read-zipf         - Read-heavy with hot keys"
	@echo "bench-write-zipf        - Write-heavy with hot keys" 
	@echo "bench-distributed-read  - Read-only uniform distribution"
	@echo "bench-distributed-write - Write-only uniform distribution"
	@echo "bench-distributed-mixed - Mixed read/write uniform distribution"
	@echo "bench-compare           - Compare read/write across all stores"
	@echo ""
	@echo "Utility:"
	@echo "bench-all               - Run all benchmarks"
	@echo "bench-save              - Save benchmark results with timestamp"