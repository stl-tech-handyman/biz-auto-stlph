#!/bin/bash
# Quick local API test script

set -e

export ENV=dev
export PORT=8080
export SERVICE_API_KEY=test-api-key-12345
export LOG_LEVEL=debug
export CONFIG_DIR=./config
export TEMPLATES_DIR=./templates

echo "🧪 Testing Go API Locally"
echo "=========================="
echo ""
echo "📦 Building..."
go build -o bin/test-server ./cmd/api

echo ""
echo "✅ Build successful!"
echo ""
echo "🚀 Starting server..."
echo "   API Key: $SERVICE_API_KEY"
echo "   Port: $PORT"
echo ""
echo "Test commands:"
echo "  curl http://localhost:$PORT/api/health"
echo "  curl http://localhost:$PORT/"
echo "  curl -H 'X-Api-Key: $SERVICE_API_KEY' -X POST http://localhost:$PORT/api/estimate -H 'Content-Type: application/json' -d '{\"eventDate\":\"2025-12-25\",\"durationHours\":4,\"numHelpers\":2}'"
echo ""
echo "Press Ctrl+C to stop"
echo ""

./bin/test-server


