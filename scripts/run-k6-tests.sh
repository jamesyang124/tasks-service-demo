#!/bin/bash
set -e

echo "🚀 Starting K6 Performance Tests for ShardStoreGopool..."
echo "Results will be saved to ./output/"

# Ensure output directory exists
mkdir -p output

# Test selection
TEST_TYPE=${1:-"all"}

case $TEST_TYPE in
    "stress")
        echo "📈 Running high-concurrency stress test (up to 1000 VUs)..."
        docker-compose -f docker-compose.test.yml run --rm k6 run /scripts/stress-test.js
        ;;
    "read")
        echo "📖 Running read-heavy performance test (optimized for ShardStoreGopool)..."
        docker-compose -f docker-compose.test.yml run --rm k6 run /scripts/read-heavy-test.js
        ;;
    "original")
        echo "🔄 Running original test suite..."
        docker-compose -f docker-compose.test.yml run --rm k6 run /scripts/test.js
        ;;
    "all"|*)
        echo "🎯 Running all test suites..."
        echo ""
        echo "1️⃣ Original test suite..."
        docker-compose -f docker-compose.test.yml run --rm k6 run /scripts/test.js
        echo ""
        echo "2️⃣ Read-heavy performance test..."
        docker-compose -f docker-compose.test.yml run --rm k6 run /scripts/read-heavy-test.js
        echo ""
        echo "3️⃣ High-concurrency stress test..."
        docker-compose -f docker-compose.test.yml run --rm k6 run /scripts/stress-test.js
        ;;
esac

echo ""
echo "✅ Tests completed! Results saved to ./output/"
echo ""
echo "📊 View results:"
echo "  Original: open ./output/k6-report.html"
echo "  Read-heavy: open ./output/read-heavy-*.html"
echo "  Stress: open ./output/stress-*.html"
echo ""
echo "💡 Usage:"
echo "  ./scripts/run-k6-tests.sh stress    # High-concurrency test"
echo "  ./scripts/run-k6-tests.sh read      # Read-heavy test"
echo "  ./scripts/run-k6-tests.sh original  # Original test"
echo "  ./scripts/run-k6-tests.sh all       # All tests (default)"