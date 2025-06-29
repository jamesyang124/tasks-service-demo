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

# Core benchmarks (1M dataset) - 5 essential scenarios
bench:
	go test -bench=. -benchmem -timeout=30m ./internal/storage/

# Individual benchmark scenarios
bench-read-zipf:
	go test -bench=BenchmarkReadZipf -benchmem -timeout=30m ./internal/storage/

bench-write-zipf:
	go test -bench=BenchmarkWriteZipf -benchmem -timeout=30m ./internal/storage/

bench-distributed-read:
	go test -bench=BenchmarkDistributedRead -benchmem -timeout=30m ./internal/storage/

bench-distributed-write:
	go test -bench=BenchmarkDistributedWrite -benchmem -timeout=30m ./internal/storage/

bench-distributed-mixed:
	go test -bench=BenchmarkDistributedMixed -benchmem -timeout=30m ./internal/storage/

bench-all: bench
	@echo "All 5 core benchmarks completed"

# Save benchmark results
bench-save:
	@mkdir -p output
	go test -bench=. -benchmem ./internal/storage/ > output/benchmark-results.txt
	@echo "Benchmark results saved to output/benchmark-results.txt"

help:
	@echo "build                  - Build binary"
	@echo "test                   - Run tests"
	@echo "run                    - Build and run"
	@echo "dev                    - Run in dev mode"
	@echo "k6                     - Run k6 load tests"
	@echo "bench                  - Run all 5 core benchmarks (1M dataset)"
	@echo "bench-read-zipf        - Read-heavy with hot keys"
	@echo "bench-write-zipf       - Write-heavy with hot keys" 
	@echo "bench-distributed-read - Read-only uniform distribution"
	@echo "bench-distributed-write - Write-only uniform distribution"
	@echo "bench-distributed-mixed - Mixed read/write uniform distribution"
	@echo "bench-all              - Run all benchmarks"
	@echo "bench-save             - Save benchmark results to output/"