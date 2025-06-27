.PHONY: build test run dev k6 bench bench-memory bench-shard bench-compare bench-highload bench-hotkey bench-all help

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

# Benchmark targets
bench:
	go test -bench=. -benchmem ./internal/storage/

bench-memory:
	go test -bench=BenchmarkMemoryStore -benchmem ./internal/storage/

bench-shard:
	go test -bench=BenchmarkShardStore -benchmem ./internal/storage/

bench-compare:
	go test -bench=BenchmarkMemoryVsShard -benchmem ./internal/storage/

bench-highload:
	go test -bench=BenchmarkHighLoad -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkBurst -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkHeavyMixed -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkWriteHeavy -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkReadHeavy -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkStress -benchmem -timeout=10m ./internal/storage/

bench-hotkey:
	go test -bench=BenchmarkHotKeySingle -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkHotKeyMultiple -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkHotKeyInterleaved -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkHotKeyZipf -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkHotKeyWorstCase -benchmem -timeout=10m ./internal/storage/
	go test -bench=BenchmarkHotKeyThunderingHerd -benchmem -timeout=10m ./internal/storage/

bench-all: bench-memory bench-shard bench-compare bench-highload bench-hotkey
	@echo "All benchmarks completed"

# Save benchmark results
bench-save:
	@mkdir -p output
	go test -bench=. -benchmem ./internal/storage/ > output/benchmark-results.txt
	@echo "Benchmark results saved to output/benchmark-results.txt"

help:
	@echo "build          - Build binary"
	@echo "test           - Run tests"
	@echo "run            - Build and run"
	@echo "dev            - Run in dev mode"
	@echo "k6             - Run k6 load tests"
	@echo "bench          - Run all benchmarks"
	@echo "bench-memory   - Run MemoryStore benchmarks"
	@echo "bench-shard    - Run ShardStore benchmarks"
	@echo "bench-compare  - Run store comparison benchmarks"
	@echo "bench-highload - Run high-load stress benchmarks"
	@echo "bench-hotkey   - Run hot key contention benchmarks"
	@echo "bench-all      - Run all benchmark categories"
	@echo "bench-save     - Save benchmark results to output/"