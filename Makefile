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

k6-stress:
	./scripts/run-k6-tests.sh stress

k6-read:
	./scripts/run-k6-tests.sh read

k6-compare:
	./scripts/run-comparative-tests.sh

# Core benchmarks (1M dataset) - All storage implementations
bench:
	go test -bench=. -benchmem -timeout=30m ./benchmarks/

# Essential benchmarks
bench-shard:
	go test -bench=.*ShardStore.* -benchmem -timeout=30m ./benchmarks/

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
	@echo "Benchmarks (1M dataset):"
	@echo "bench                   - Run all storage benchmarks"
	@echo "bench-shard             - ShardStore benchmarks only"
	@echo "bench-compare           - Compare read/write across all stores"
	@echo "bench-save              - Save benchmark results with timestamp"