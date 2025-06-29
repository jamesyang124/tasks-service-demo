#!/bin/bash
set -e

echo "🔬 Running Comparative Storage Performance Tests..."
echo "Testing: MemoryStore, BigCacheStore, ShardStore, ShardStoreGopool"
echo ""

# Ensure output directory exists
mkdir -p output

# Kill any existing server
pkill -f "go run" 2>/dev/null || true
sleep 2

# Storage types to test
STORAGE_TYPES=("memory" "shard" "gopool")

# Test each storage type
for storage in "${STORAGE_TYPES[@]}"; do
    echo "🧪 Testing $(echo $storage | tr '[:lower:]' '[:upper:]') storage implementation..."
    
    # Start server with specific storage type
    STORAGE_TYPE=$storage go run ./internal/cmd/main.go &
    SERVER_PID=$!
    echo "Started server with PID: $SERVER_PID"
    
    # Wait for server to start
    sleep 3
    
    # Test if server is responding
    if curl -s http://localhost:8080/tasks > /dev/null; then
        echo "✅ Server responding, starting k6 test..."
        
        # Run k6 test with storage type
        docker run --rm \
            -v "${PWD}/k6:/scripts" \
            -v "${PWD}/output:/output" \
            --network="host" \
            -e STORAGE_TYPE=$storage \
            grafana/k6:latest run /scripts/comparative-test.js
        
        echo "✅ $(echo $storage | tr '[:lower:]' '[:upper:]') test completed"
    else
        echo "❌ Server not responding for ${storage} storage"
    fi
    
    # Kill server
    kill $SERVER_PID 2>/dev/null || true
    sleep 2
    
    echo ""
done

# Test BigCacheStore if available
echo "🧪 Testing BIGCACHE storage implementation..."

# Check if BigCacheStore is available in codebase
if grep -r "BigCacheStore" internal/storage/ > /dev/null 2>&1; then
    echo "BigCacheStore found, testing..."
    
    # Note: BigCacheStore would need to be added to main.go storage options
    echo "⚠️  BigCacheStore test skipped - needs integration in main.go"
    echo "   Add 'bigcache' option to STORAGE_TYPE in main.go to enable"
else
    echo "⚠️  BigCacheStore not found in codebase, skipping..."
fi

echo ""
echo "🎉 Comparative testing completed!"
echo ""
echo "📊 View results:"
echo "  Memory:  open ./output/comparative-memory-*.html"
echo "  Shard:   open ./output/comparative-shard-*.html" 
echo "  Gopool:  open ./output/comparative-gopool-*.html"
echo ""
echo "💡 Quick comparison:"
echo "grep -h 'RPS:' output/comparative-*.json | sort -k2 -nr"