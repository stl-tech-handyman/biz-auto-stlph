#!/bin/bash
# Restart server with Gmail credentials loaded from .env

set -e

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

echo -e "${CYAN}🔄 Restarting server with email credentials...${NC}"

# Stop any existing server on port 8080
if command -v lsof >/dev/null 2>&1; then
    PID=$(lsof -ti:8080 2>/dev/null)
    if [ ! -z "$PID" ]; then
        echo -e "${YELLOW}  Stopping process $PID...${NC}"
        kill -9 $PID 2>/dev/null || true
        sleep 2
    fi
elif command -v netstat >/dev/null 2>&1; then
    PID=$(netstat -ano | grep ":8080 " | grep LISTENING | awk '{print $5}' | head -1)
    if [ ! -z "$PID" ]; then
        echo -e "${YELLOW}  Stopping process $PID...${NC}"
        taskkill //F //PID $PID 2>/dev/null || true
        sleep 2
    fi
fi

# Start server using start-local.sh
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
echo -e "${GREEN}🚀 Starting server...${NC}"
bash "$SCRIPT_DIR/start-local.sh"

