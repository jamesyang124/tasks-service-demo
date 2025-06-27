# K6 Load Testing

Simple load testing for the Tasks API.

## Usage

```bash
make k6
```

This will:
1. Start the API in Docker
2. Run k6 tests (10 users, 30 seconds)
3. Test all CRUD operations
4. Clean up containers

## Manual Testing

```bash
# Local k6 (if installed)
k6 run k6/test.js

# Docker (exports results to ./output/)
./scripts/run-k6-tests.sh
```

## Test Results

After running tests, you'll find:
- `./output/k6-summary.json` - Detailed JSON metrics
- `./output/k6-report.html` - Visual HTML report

Open the HTML report:
```bash
open ./output/k6-report.html
```