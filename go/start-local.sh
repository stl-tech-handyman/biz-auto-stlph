#!/bin/bash
# Local development startup script for Go API
# WARNING: This uses LIVE Stripe keys - be careful!

export STRIPE_SECRET_KEY_PROD="sk_live_YOUR_PROD_KEY_HERE"
export SERVICE_API_KEY="test-api-key-12345"
export ENV="dev"
export PORT="8080"
export LOG_LEVEL="debug"
export CONFIG_DIR="./config"
export TEMPLATES_DIR="./templates"

echo "🚀 Starting Go API Server..."
echo ""
echo "Configuration:"
echo "  Stripe: LIVE (Production) Key"
echo "  API Key: $SERVICE_API_KEY"
echo "  Port: $PORT"
echo "  Environment: $ENV"
echo ""
echo "⚠️  WARNING: Using LIVE Stripe key - real charges will occur!"
echo ""
echo "Test the final invoice endpoint:"
echo "  curl -X POST http://localhost:$PORT/api/stripe/final-invoice \\"
echo "    -H \"X-Api-Key: $SERVICE_API_KEY\" \\"
echo "    -H \"Content-Type: application/json\" \\"
echo "    -d '{\"email\":\"your-email@example.com\",\"name\":\"Test Customer\",\"totalAmount\":1000.0,\"depositPaid\":400.0}'"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

go run ./cmd/api

