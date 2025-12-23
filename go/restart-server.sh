#!/bin/bash
# Restart script for Go API Server
# Automatically kills process on port 8080 and restarts server

PORT="8080"
export STRIPE_SECRET_KEY_PROD="sk_live_YOUR_PROD_KEY_HERE"
export SERVICE_API_KEY="test-api-key-12345"
export ENV="dev"
export PORT=$PORT
export LOG_LEVEL="debug"
export CONFIG_DIR="./config"
export TEMPLATES_DIR="./templates"

echo "🔄 Restarting Go API Server..."
echo ""

# Find and kill process using port 8080
echo "Checking port $PORT..."
if command -v lsof >/dev/null 2>&1; then
    # macOS/Linux with lsof
    PID=$(lsof -ti:$PORT 2>/dev/null)
    if [ ! -z "$PID" ]; then
        echo "⚠️  Port $PORT is in use by process(es): $PID"
        echo "   Killing process(es)..."
        kill -9 $PID 2>/dev/null
        echo "   Waiting 2 seconds for port to be released..."
        sleep 2
    else
        echo "✅ Port $PORT is free"
    fi
elif command -v netstat >/dev/null 2>&1; then
    # Windows with netstat (Git Bash)
    PID=$(netstat -ano | grep ":$PORT " | grep LISTENING | awk '{print $5}' | head -1)
    if [ ! -z "$PID" ]; then
        echo "⚠️  Port $PORT is in use by process: $PID"
        echo "   Killing process..."
        taskkill //F //PID $PID 2>/dev/null
        echo "   Waiting 2 seconds for port to be released..."
        sleep 2
    else
        echo "✅ Port $PORT is free"
    fi
else
    echo "⚠️  Cannot check port (lsof/netstat not available), proceeding anyway..."
fi

echo ""
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
echo "Press Ctrl+C to stop the server"
echo ""

go run ./cmd/api

