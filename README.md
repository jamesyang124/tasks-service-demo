# Task API Service

[![Go Version](https://img.shields.io/badge/Go-1.21.13-blue.svg)](https://golang.org)
[![Fiber Version](https://img.shields.io/badge/Fiber-2.52.8-green.svg)](https://gofiber.io)
[![Gopool Version](https://img.shields.io/badge/Gopool-0.1.2-orange.svg)](https://github.com/bytedance/gopkg)

A high-performance RESTful Task API service built with Go and Fiber framework, featuring multiple optimized storage implementations and comprehensive performance testing.

## Features

- **RESTful API**: Full CRUD operations for task management
- **Multiple Storage Implementations**: Optimized sharded storage with ByteDance gopool
- **High Performance**: 12.6x performance improvement (156.5ns → 12.40ns reads)
- **Production Ready**: Optimized for high-traffic production environments
- **Thread Safety**: Concurrent request handling with optimized locking strategies
- **Comprehensive Testing**: Unit tests and benchmarks
- **Performance Monitoring**: Detailed benchmark suite
- **Complete Storage Coverage**: Benchmarks for all storage implementations (Memory, Shard, Gopool, Channel)
- **Docker Support**: Multi-stage Docker build for production deployment
- **Graceful Shutdown**: Clean resource cleanup with proper signal handling
- **Environment Configuration**: Dotenv support for easy local development
- **Structured Logging**: Uber Zap logger with ISO8601 time encoding

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks` | Retrieve all tasks |
| GET | `/tasks/{id}` | Retrieve a specific task by ID |
| POST | `/tasks` | Create a new task |
| PUT | `/tasks/{id}` | Update an existing task |
| DELETE | `/tasks/{id}` | Delete a task |
| GET | `/health` | Health check endpoint |
| GET | `/version` | API version information |

## Task Model

```json
{
  "id": 1,
  "name": "Task name",
  "status": 0
}
```

### Field Descriptions
- `id` (integer): Auto-generated unique identifier (read-only)
- `name` (string): Task name/description 
  - **Required**: Must not be empty
  - **Length**: 1-100 characters
  - **Validation**: `required,min=1,max=100`
- `status` (integer): Task completion status
  - **Values**: Must be exactly `0` or `1`
  - **Validation**: `oneof=0 1`
  - `0`: Incomplete task
  - `1`: Completed task

## API Examples

### Create a Task
**Request:**
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"name": "Learn Go", "status": 0}'
```

**Response (201 Created):**
```json
{
  "id": 1,
  "name": "Learn Go",
  "status": 0
}
```

**Error Response (400 Bad Request):**
```json
{
  "code": 1002,
  "message": "name is required"
}
```

### Get All Tasks
**Request:**
```bash
curl http://localhost:8080/tasks
```

**Response (200 OK):**
```json
[
  {
    "id": 1,
    "name": "Learn Go",
    "status": 0
  },
  {
    "id": 2,
    "name": "Build API",
    "status": 1
  }
]
```

**Empty Response (200 OK):**
```json
[]
```

### Get a Specific Task
**Request:**
```bash
curl http://localhost:8080/tasks/1
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Learn Go",
  "status": 0
}
```

**Error Response (400 Bad Request):**
```json
{
  "code": 1001,
  "message": "Task not found"
}
```

### Update a Task
**Request:**
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Learn Go - Updated", "status": 1}'
```

**Response (200 OK):**
```json
{
  "id": 1,
  "name": "Learn Go - Updated",
  "status": 1
}
```

**Error Response (400 Bad Request):**
```json
{
  "code": 1001,
  "message": "Task not found"
}
```

**Error Response (400 Bad Request):**
```json
{
  "code": 1002,
  "message": "Status must be 0 or 1"
}
```

### Delete a Task
**Request:**
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

**Response (204 No Content):**
```
(No response body)
```

**Note**: DELETE operations are idempotent and always return 204, even if the task doesn't exist.

### Health Check
**Request:**
```bash
curl http://localhost:8080/health
```

**Response (200 OK):**
```json
{
  "status": "ok",
  "message": "Task API is running"
}
```

### Version Information
**Request:**
```bash
curl http://localhost:8080/version
```

**Response (200 OK):**
```json
{
  "version": "1.0.0"
}
```

## Error Codes

The API uses standardized integer error codes for consistent error handling:

### Error Code Ranges
- **1000-1999**: Task-related errors
- **2000-2999**: Request validation errors  
- **5000-5999**: System/server errors

### Available Error Codes

| Error Code | HTTP Status | Description | Example |
|------------|-------------|-------------|---------|
| `1001` | 400 | Task with specified ID does not exist | GET /tasks/999 |
| `1002` | 400 | Invalid task data (validation failed) | Missing name, invalid status |
| `1003` | 400 | Task name is required | Empty name field |
| `1004` | 400 | Task name exceeds 100 characters | Name > 100 chars |
| `1005` | 400 | Status must be 0 or 1 | status: 2 |
| `2001` | 400 | Request body is not valid JSON | Malformed JSON |
| `2002` | 400 | ID parameter is not a valid integer | /tasks/abc |
| `2003` | 400 | Required fields are missing | No request body |
| `5001` | 500 | Internal server error | Database error |
| `5002` | 500 | Storage system error | Storage unavailable |

### Error Response Format

All error responses follow this consistent structure:

```json
{
  "code": 1001,
  "message": "Human-readable error message"
}
```

**Note**: All error responses include both `code` and `message` fields for consistent error handling.

## Quick Start

### Prerequisites
- Go 1.18 or higher
- Docker (optional)

### Environment Configuration

The application supports environment variables for configuration. You can use a `.env` file for local development:

1. Copy the example environment file:
```bash
cp env.example .env
```

2. Edit `.env` with your preferred settings:
```bash
# Storage Configuration
STORAGE_TYPE=gopool
SHARD_COUNT=32

# Application Configuration
APP_VERSION=1.0.0

# Server Configuration
PORT=8080
```

**Available Environment Variables:**
- `STORAGE_TYPE`: Storage implementation (`memory`, `shard`, `gopool`)
- `SHARD_COUNT`: Number of shards for sharded storage (default: 32)
- `APP_VERSION`: Application version (default: 1.0.0)
- `PORT`: Server port (default: 8080)

### Running Locally

1. Clone the repository:
```bash
git clone <repository-url>
cd tasks-service-demo
```

2. Install dependencies:
```bash
go mod tidy
```

3. Run the application:
```bash
go run ./cmd/tasks-service-demo/
```

The server will start on `http://localhost:8080`

## Storage Implementations

### Available Storage Types

| Implementation | Read Performance | Write Performance | Memory Allocations | Production Ready |
|---------------|------------------|-------------------|-------------------|------------------|
| **ShardStoreGopool** | **12.3 ns/op** | 61.5 ns/op | 0 B/op | 🏆 **Best** |
| **ShardStore** | 12.5 ns/op | **61.0 ns/op** | 0 B/op | ✅ **Excellent** |
| **MemoryStore** | 156.5 ns/op | 312.5 ns/op | 0 B/op | ⚠️ **Limited** |
| **ChannelStore** | 607.5 ns/op | 693.5 ns/op | 192 B/op | ❌ **Educational** |

### Configuration

Set storage type via environment variable:

```bash
# Production (default) - Best performance
STORAGE_TYPE=gopool go run ./cmd/tasks-service-demo/

# Alternative production option
STORAGE_TYPE=shard go run ./cmd/tasks-service-demo/

# Development/testing
STORAGE_TYPE=memory go run ./cmd/tasks-service-demo/
```

## Performance Results

### Optimization Journey Validation

**12.6x Performance Improvement Achieved:**
- **Baseline (MemoryStore)**: 156.5ns reads, 312.5ns writes
- **Optimized (ShardStoreGopool)**: 12.40ns reads, 62.69ns writes
- **Improvement**: 12.6x faster reads, 5.0x faster writes

### Benchmark Results

**Current Performance Results (Apple M4 Pro):**

| Storage Implementation | Read Performance | Write Performance | Memory Allocations | Production Ready |
|----------------------|------------------|-------------------|-------------------|------------------|
| **ShardStoreGopool** | **12.40 ns/op** | 62.69 ns/op | 0 B/op | 🏆 **Best** |
| **ShardStore** | 12.55 ns/op | **61.44 ns/op** | 0 B/op | ✅ **Excellent** |
| **MemoryStore** | 156.5 ns/op | 312.5 ns/op | 0 B/op | ⚠️ **Limited** |
| **ChannelStore** | 607.5 ns/op | 693.5 ns/op | 192 B/op | ❌ **Educational** |

**Performance Improvements:**
- **ShardStoreGopool vs MemoryStore**: 12.7x faster overall
- **ShardStore vs MemoryStore**: 12.5x faster overall
- **ShardStoreGopool vs ShardStore**: 2% faster reads with better consistency

## Development

### Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific test
go test -run TestName ./...
```

### Building the Application

```bash
# Core commands
make build         # Build binary
make test          # Run tests with coverage
make run           # Build and run application
make dev           # Run in development mode
make setup-env     # Create .env file from template

# Performance testing
make bench         # Run all benchmarks

# Help
make help          # Show all available commands
```

### Direct Go Commands
```bash
# Build application
go build -o bin/tasks-service-demo ./cmd/tasks-service-demo

# Run tests
go test -cover ./...

# Format and vet code
go fmt ./...
go vet ./...
```

## Project Structure

```
tasks-service-demo/
├── go.mod                      # Go module definition
├── Makefile                    # Build automation
├── Dockerfile                  # Docker configuration
├── env.example                 # Environment variables template
├── cmd/                        # Application entry points
│   └── tasks-service-demo/     # Main application
│       ├── main.go            # Application entry point
│       └── main_test.go       # Main application tests
├── internal/                   # Internal application code
│   ├── entities/              # Business entities
│   │   ├── task.go            # Core Task entity
│   │   └── task_test.go       # Entity tests
│   ├── requests/              # API request/response models
│   │   ├── request.go         # CreateTaskRequest, UpdateTaskRequest
│   │   ├── validator.go       # Validation logic
│   │   └── *_test.go          # Comprehensive tests
│   ├── handlers/
│   │   ├── task_handler.go    # HTTP handlers
│   │   ├── health_handler.go  # Health check handler
│   │   ├── version_handler.go # Version handler
│   │   └── *_test.go          # Handler tests
│   ├── services/
│   │   ├── task.go            # Business logic layer
│   │   └── task_test.go       # Service tests
│   ├── storage/               # Storage implementations
│   │   ├── store.go           # Store interface & singleton
│   │   ├── naive/             # Naive Memory Store
│   │   │   ├── memory.go      # Simple single-mutex implementation
│   │   │   └── memory_test.go # Memory store tests
│   │   ├── shard/             # High-Performance Shard Store
│   │   │   ├── shard.go       # Optimized sharded storage
│   │   │   ├── shard_gopool.go # ByteDance gopool optimization
│   │   │   ├── shard_unit.go  # Lightweight storage units
│   │   │   ├── shard_utils.go # Utility functions
│   │   │   └── shard_test.go  # Comprehensive tests
│   │   └── channel/           # Actor Model Store
│   │       ├── channel_store.go # Message passing implementation
│   │       └── channel_store_test.go # Channel store tests
│   ├── routes/
│   │   └── routes.go          # Route definitions
│   ├── middleware/
│   │   ├── validation.go      # Request validation middleware
│   │   └── *_test.go          # Middleware tests
│   ├── logger/
│   │   ├── logger.go          # Structured logging with Zap
│   │   └── logger_test.go     # Logger tests
│   └── errors/
│       ├── app.go             # Application error types
│       ├── codes.go           # Error code definitions
│       ├── response.go        # Error response formatting
│       └── *_test.go          # Error handling tests
├── benchmarks/                 # Performance benchmark suite
│   ├── README.md              # Benchmark documentation
│   ├── common.go              # Shared benchmark utilities
│   ├── memory_bench_test.go   # MemoryStore benchmarks
│   ├── shard_bench_test.go    # ShardStore benchmarks
│   ├── shard_gopool_bench_test.go # ShardStoreGopool benchmarks
│   └── channel_bench_test.go  # ChannelStore benchmarks
├── docs/                      # Technical documentation
│   ├── OPTIMIZATION_DECISIONS.md # Optimization journey
│   └── PERFORMANCE_COMPARISON.md # Performance analysis
└── output/                    # Generated reports (ignored by git)
```

## Technical Details

- **Framework**: Fiber v2 (high-performance Go web framework)
- **Storage**: Multiple implementations with sharded optimization (12.6x performance improvement)
- **Concurrency**: Sharded storage with per-core worker pools
- **Performance**: 12.40ns reads, 62.69ns writes (ShardStoreGopool)
- **Thread Safety**: Optimized locking strategies (RWMutex, atomic operations)
- **Benchmarking**: 1M dataset benchmarks with realistic workload patterns
- **Logging**: Structured logging with Uber Zap
- **Configuration**: Environment-based configuration with dotenv support

## Documentation

### Performance & Optimization
- [**Optimization Decisions**](docs/OPTIMIZATION_DECISIONS.md) - Complete optimization journey
- [**Performance Comparison**](docs/PERFORMANCE_COMPARISON.md) - Storage implementation analysis  

### Testing & Benchmarks
- [**Benchmark Suite**](benchmarks/README.md) - Comprehensive performance benchmarks

## License

This project is licensed under the MIT License.