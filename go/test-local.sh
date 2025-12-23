#!/bin/bash
# Local testing script for Go API

set -e

echo "🧪 Testing Go API Locally"
echo "=========================="

# Set test environment variables
export ENV=dev
export PORT=8080
export SERVICE_API_KEY=test-api-key-12345
export LOG_LEVEL=debug

echo ""
echo "📦 Running unit tests..."
go test ./... -v -cover

echo ""
echo "✅ All tests passed!"
echo ""
echo "🚀 Starting server for manual testing..."
echo "   API Key: $SERVICE_API_KEY"
echo "   Port: $PORT"
echo ""
echo "Test endpoints:"
echo "  curl http://localhost:$PORT/api/health"
echo "  curl -H 'X-Api-Key: $SERVICE_API_KEY' -X POST http://localhost:$PORT/api/estimate -d '{\"eventDate\":\"2025-12-25\",\"durationHours\":4,\"numHelpers\":2}'"
echo ""
echo "Press Ctrl+C to stop"
echo ""

# Run server
go run ./cmd/api


