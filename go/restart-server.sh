#!/bin/bash
# Restart server with Gmail credentials loaded from .env

echo "🔄 Restarting server with email credentials..."

# Find and kill process using port 8080
PORT="8080"
if command -v lsof >/dev/null 2>&1; then
    # macOS/Linux with lsof
    PID=$(lsof -ti:$PORT 2>/dev/null)
    if [ ! -z "$PID" ]; then
        echo "  Stopping process $PID on port $PORT..."
        kill -9 $PID 2>/dev/null
        sleep 2
    fi
elif command -v netstat >/dev/null 2>&1; then
    # Windows with netstat (Git Bash)
    PID=$(netstat -ano | grep ":$PORT " | grep LISTENING | awk '{print $5}' | head -1)
    if [ ! -z "$PID" ]; then
        echo "  Stopping process $PID on port $PORT..."
        taskkill //F //PID $PID 2>/dev/null
        sleep 2
    fi
fi

# Auto-fetch Gmail credentials if not in .env
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

if [ ! -f "$ENV_FILE" ] || ! grep -q "GMAIL_CREDENTIALS_JSON=" "$ENV_FILE" || [ -z "$(grep "GMAIL_CREDENTIALS_JSON=" "$ENV_FILE" | cut -d'=' -f2)" ]; then
    echo "  Gmail credentials not found - fetching from Secret Manager..."
    if [ -f "$SCRIPT_DIR/scripts/get-gmail-credentials.sh" ]; then
        bash "$SCRIPT_DIR/scripts/get-gmail-credentials.sh" || {
            echo "⚠️  Failed to fetch credentials - email features will not work"
        }
    else
        echo "⚠️  get-gmail-credentials.sh not found"
    fi
fi

# Fix Go toolchain configuration (same fix as start-local.sh)
export GOTOOLCHAIN=auto
GO_BIN=$(which go 2>/dev/null || command -v go 2>/dev/null)
if [ -n "$GO_BIN" ]; then
    POSSIBLE_GOROOT=$(dirname "$(dirname "$GO_BIN")")
    if [ -d "$POSSIBLE_GOROOT/src" ] && [ -f "$POSSIBLE_GOROOT/src/runtime/runtime.go" ]; then
        export GOROOT="$POSSIBLE_GOROOT"
    elif [ -d "C:/Program Files/Go" ] && [ -f "C:/Program Files/Go/src/runtime/runtime.go" ]; then
        export GOROOT="C:/Program Files/Go"
    fi
fi

# Start server
echo "🚀 Starting server..."
cd "$SCRIPT_DIR"
bash start-local.sh
