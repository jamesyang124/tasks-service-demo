#!/bin/bash
set -e

echo "Starting k6 tests..."
echo "Results will be saved to ./output/"

# Ensure output directory exists
mkdir -p output

docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
docker-compose -f docker-compose.test.yml down

echo ""
echo "Test completed! Results saved:"
echo "  📊 JSON Report: ./output/k6-summary.json"
echo "  📈 HTML Report: ./output/k6-report.html"
echo ""
echo "Open HTML report: open ./output/k6-report.html"