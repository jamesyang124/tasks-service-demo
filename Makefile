.PHONY: build test run dev bench k6 setup-env help

# Core commands
build:
	go build -o bin/tasks-service-demo ./cmd/tasks-service-demo

test:
	go test -cover ./internal/... ./cmd/tasks-service-demo/...

run: build
	./bin/tasks-service-demo

dev:
	go run ./cmd/tasks-service-demo/main.go

# Environment setup
setup-env:
	@if [ ! -f .env ]; then \
		cp env.example .env; \
		echo "✅ Created .env file from env.example"; \
		echo "📝 Edit .env to customize your configuration"; \
	else \
		echo "⚠️  .env file already exists"; \
	fi

# Performance testing
bench:
	go test -bench=. -benchmem -timeout=30m ./benchmarks/

k6:
	./scripts/run-k6-tests.sh

help:
	@echo "Core commands:"
	@echo "  build     - Build binary"
	@echo "  test      - Run tests with coverage"
	@echo "  run       - Build and run application"
	@echo "  dev       - Run in development mode"
	@echo ""
	@echo "Environment:"
	@echo "  setup-env - Create .env file from env.example"
	@echo ""
	@echo "Performance testing:"
	@echo "  bench     - Run all benchmarks (1M dataset)"
	@echo "  k6        - Run K6 load tests"