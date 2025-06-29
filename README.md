# Task API Service

A RESTful Task API service built with Go and Fiber framework, featuring in-memory storage with thread-safe operations and graceful shutdown.

## Features

- **RESTful API**: Full CRUD operations for task management
- **High Performance**: Built with Fiber web framework
- **Thread Safety**: Concurrent request handling with RWMutex
- **Graceful Shutdown**: Clean server shutdown with signal handling
- **Comprehensive Tests**: Unit tests for all components
- **Docker Support**: Multi-stage Docker build for production deployment

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/tasks` | Retrieve all tasks |
| GET | `/tasks/{id}` | Retrieve a specific task by ID |
| POST | `/tasks` | Create a new task |
| PUT | `/tasks/{id}` | Update an existing task |
| DELETE | `/tasks/{id}` | Delete a task |
| GET | `/health` | Health check endpoint |

## Task Model

```json
{
  "id": 1,
  "name": "Task name",
  "status": 0
}
```

### Field Descriptions
- `id` (integer): Auto-generated unique identifier
- `name` (string): Task name/description (required)
- `status` (integer): Task completion status
  - `0`: Incomplete task
  - `1`: Completed task

## Quick Start

### Prerequisites
- Go 1.18 or higher
- Docker (optional)

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
go run main.go
```

The server will start on `http://localhost:8080`


## Performance Results:

### High-Load Benchmark Results:

  Key Performance Insights (1M Dataset):

  | Scenario          | MemoryStore | ShardStore | Winner & Gain      |
  |-------------------|-------------|------------|--------------------|
  | Read Zipf         | 129.9 ns    | 131.3 ns   | MemoryStore +1.1%  |
  | Write Zipf        | 211.3 ns    | 157.1 ns   | ShardStore +25.6%  |
  | Distributed Read  | 129.6 ns    | 157.8 ns   | MemoryStore +17.9% |
  | Distributed Write | 215.0 ns    | 154.0 ns   | ShardStore +28.4%  |
  | Distributed Mixed | 108.9 ns    | 155.9 ns   | MemoryStore +30.2% |

  Clear Recommendations:

  ✅ Use MemoryStore for: Read-heavy & mixed workloads✅ Use ShardStore for: Write-heavy workloads (25-28% faster)

  Available Commands:

  - make bench - Run all 5 benchmarks
  - make bench-save - Save results to output/
  - Individual scenario commands available

## API Examples

### Create a Task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"name": "Learn Go", "status": 0}'
```

### Get All Tasks
```bash
curl http://localhost:8080/tasks
```

### Get a Specific Task
```bash
curl http://localhost:8080/tasks/1
```

### Update a Task
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "Learn Go - Updated", "status": 1}'
```

### Delete a Task
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

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

Using Makefile (recommended):
```bash
# Build binary to bin/ directory
make build

# Build and run
make run

# Development mode (no build)
make dev

# Clean build artifacts
make clean
```

Using Go commands directly:
```bash
# Build to bin directory
go build -o bin/tasks-api

# Format code
go fmt ./...

# Vet code for issues
go vet ./...
```

### Available Make Commands
```bash
make help          # Show all available commands
make build         # Build binary to bin/ directory
make test          # Run all tests
make test-coverage # Run tests with coverage
make run           # Build and run application
make dev           # Run in development mode
make fmt           # Format code
make vet           # Vet code for issues
make docker-build  # Build Docker image
make docker-run    # Run Docker container
make prepare       # Format, vet, test, and build
```

## Project Structure

```
tasks-service-demo/
├── main.go                     # Application entry point
├── go.mod                      # Go module definition
├── Makefile                    # Build automation
├── Dockerfile                  # Docker configuration
├── bin/                        # Binary executables (ignored by git)
├── models/
│   └── task.go                # Task model and validation
├── handlers/
│   ├── task_handler.go        # HTTP handlers
│   └── task_handler_test.go   # Handler tests
├── storage/
│   ├── memory_store.go        # In-memory storage
│   └── memory_store_test.go   # Storage tests
└── routes/
    └── task_routes.go         # Route definitions
```

## Technical Details

- **Framework**: Fiber v2 (high-performance Go web framework)
- **Storage**: In-memory with sync.RWMutex for thread safety
- **Concurrency**: Supports concurrent read/write operations
- **Graceful Shutdown**: Handles SIGTERM and SIGINT signals
- **Error Handling**: Comprehensive error responses with proper HTTP status codes
- **Middleware**: CORS, logging, and recovery middleware included

## License

This project is licensed under the MIT License.