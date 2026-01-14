#!/bin/bash
# Local development startup script for Go API
# WARNING: This uses LIVE Stripe keys - be careful!

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ENV_FILE="$SCRIPT_DIR/.env"

# Fix Go toolchain configuration
# Allow toolchain auto-download but prefer local if compatible
# This ensures go.mod requirements are met while using local Go when possible
export GOTOOLCHAIN=auto
# #region agent log
LOG_FILE="c:/Users/Alexey/Code/biz-operating-system/stlph/.cursor/debug.log"
echo "{\"sessionId\":\"debug-session\",\"runId\":\"post-fix\",\"hypothesisId\":\"B\",\"location\":\"start-local.sh:toolchain-fix\",\"message\":\"Toolchain fix applied\",\"data\":{\"GOTOOLCHAIN\":\"local\",\"originalGOROOT\":\"$GOROOT\"},\"timestamp\":$(date +%s%3N)}" >> "$LOG_FILE"
# #endregion

# Also fix GOROOT if it's pointing to toolchain module cache
if [ -z "$GOROOT" ] || [[ "$GOROOT" == *"toolchain@"* ]]; then
    # Try to find actual Go installation
    GO_BIN=$(which go 2>/dev/null || command -v go 2>/dev/null)
    # #region agent log
    echo "{\"sessionId\":\"debug-session\",\"runId\":\"post-fix\",\"hypothesisId\":\"B\",\"location\":\"start-local.sh:goroot-fix\",\"message\":\"Found Go binary\",\"data\":{\"goBin\":\"$GO_BIN\"},\"timestamp\":$(date +%s%3N)}" >> "$LOG_FILE"
    # #endregion
    if [ -n "$GO_BIN" ]; then
        # Go binary is typically at <GOROOT>/bin/go
        POSSIBLE_GOROOT=$(dirname "$(dirname "$GO_BIN")")
        if [ -d "$POSSIBLE_GOROOT/src" ] && [ -f "$POSSIBLE_GOROOT/src/runtime/runtime.go" ]; then
            export GOROOT="$POSSIBLE_GOROOT"
            # #region agent log
            echo "{\"sessionId\":\"debug-session\",\"runId\":\"post-fix\",\"hypothesisId\":\"B\",\"location\":\"start-local.sh:goroot-fix\",\"message\":\"Set GOROOT from Go binary path\",\"data\":{\"newGOROOT\":\"$GOROOT\",\"stdlibExists\":true},\"timestamp\":$(date +%s%3N)}" >> "$LOG_FILE"
            # #endregion
        elif [ -d "/usr/local/go" ] && [ -f "/usr/local/go/src/runtime/runtime.go" ]; then
            export GOROOT="/usr/local/go"
            # #region agent log
            echo "{\"sessionId\":\"debug-session\",\"runId\":\"post-fix\",\"hypothesisId\":\"B\",\"location\":\"start-local.sh:goroot-fix\",\"message\":\"Set GOROOT to /usr/local/go\",\"data\":{\"newGOROOT\":\"$GOROOT\"},\"timestamp\":$(date +%s%3N)}" >> "$LOG_FILE"
            # #endregion
        elif [ -d "C:/Program Files/Go" ] && [ -f "C:/Program Files/Go/src/runtime/runtime.go" ]; then
            export GOROOT="C:/Program Files/Go"
            # #region agent log
            echo "{\"sessionId\":\"debug-session\",\"runId\":\"post-fix\",\"hypothesisId\":\"B\",\"location\":\"start-local.sh:goroot-fix\",\"message\":\"Set GOROOT to C:/Program Files/Go\",\"data\":{\"newGOROOT\":\"$GOROOT\"},\"timestamp\":$(date +%s%3N)}" >> "$LOG_FILE"
            # #endregion
        fi
    fi
fi

export STRIPE_SECRET_KEY_PROD="sk_live_YOUR_PROD_KEY_HERE"
export SERVICE_API_KEY="test-api-key-12345"
export ENV="dev"
export PORT="8080"
export LOG_LEVEL="debug"
export CONFIG_DIR="./config"
export TEMPLATES_DIR="./templates"

# Gmail credentials - auto-fetch if not set
if [ -z "$GMAIL_CREDENTIALS_JSON" ]; then
    # Check if .env has it
    if [ -f "$ENV_FILE" ] && grep -q "GMAIL_CREDENTIALS_JSON=" "$ENV_FILE"; then
        # Load from .env (godotenv will handle this, but we can also source it)
        source "$ENV_FILE" 2>/dev/null || true
    fi
    
    # If still not set, try to fetch from Secret Manager
    if [ -z "$GMAIL_CREDENTIALS_JSON" ]; then
        echo "ℹ️  Gmail credentials not found - attempting to fetch from Secret Manager..."
        if [ -f "$SCRIPT_DIR/scripts/get-gmail-credentials.sh" ]; then
            bash "$SCRIPT_DIR/scripts/get-gmail-credentials.sh" || {
                echo "⚠️  Failed to fetch credentials - email features will not work"
                echo "   You can manually run: bash scripts/get-gmail-credentials.sh"
            }
            # Reload .env if it was updated
            if [ -f "$ENV_FILE" ]; then
                source "$ENV_FILE" 2>/dev/null || true
            fi
        else
            echo "⚠️  get-gmail-credentials.sh not found - email features will not work"
        fi
    fi
fi

# Set GMAIL_FROM if not set
if [ -z "$GMAIL_FROM" ]; then
    export GMAIL_FROM="team@stlpartyhelpers.com"
fi

echo "🚀 Starting Go API Server..."
echo ""
echo "Configuration:"
echo "  Stripe: LIVE (Production) Key"
echo "  API Key: $SERVICE_API_KEY"
echo "  Port: $PORT"
echo "  Environment: $ENV"
if [ -n "$GMAIL_CREDENTIALS_JSON" ]; then
    echo "  Email: ✅ Configured"
else
    echo "  Email: ⚠️  Not configured (GMAIL_CREDENTIALS_JSON not set)"
fi
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

