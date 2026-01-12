# Makefile for Go API Development and Deployment
# 
# Quick Start:
#   make help          - Show all available commands
#   make run           - Run API server locally
#   make restart       - Restart the API server

.PHONY: help

# ==============================================================================
# Configuration
# ==============================================================================

PORT := 8080
PROJECT_DEV := bizops360-dev
PROJECT_PROD := bizops360-prod
REGION := us-central1
SERVICE_DEV := bizops360-api-go-dev
SERVICE_PROD := bizops360-api-go-prod
IMAGE_DEV := gcr.io/$(PROJECT_DEV)/$(SERVICE_DEV)
IMAGE_PROD := gcr.io/$(PROJECT_PROD)/$(SERVICE_PROD)
BINARY_NAME := server
BIN_DIR := bin

# Detect OS
UNAME_S := $(shell uname -s 2>/dev/null || echo Windows)
ifeq ($(UNAME_S),Windows)
	KILL_CMD := taskkill //F //PID
	NETSTAT_CMD := netstat -ano
	FIND_PID := netstat -ano | grep ":$(PORT) " | grep LISTENING | awk '{print $$5}' | head -1
else
	KILL_CMD := kill -9
	NETSTAT_CMD := lsof -i:$(PORT)
	FIND_PID := lsof -ti:$(PORT)
endif

# ==============================================================================
# Help
# ==============================================================================

help: ## Show this help message
	@echo "╔══════════════════════════════════════════════════════════════════════════════╗"
	@echo "║                    Go API Development Makefile                                ║"
	@echo "╚══════════════════════════════════════════════════════════════════════════════╝"
	@echo ""
	@echo "📦 DEVELOPMENT"
	@echo "  make run           - Run API server locally (foreground)"
	@echo "  make build         - Build Go binary to bin/server"
	@echo "  make test          - Run all tests"
	@echo "  make test-verbose  - Run tests with verbose output"
	@echo "  make clean         - Clean build artifacts"
	@echo ""
	@echo "🔄 SERVER MANAGEMENT"
	@echo "  make start         - Start API server in background"
	@echo "  make stop          - Stop API server (kill process on port $(PORT))"
	@echo "  make restart       - Restart API server (stop + start)"
	@echo "  make status        - Check if server is running"
	@echo "  make logs          - Show server logs (if running in background)"
	@echo ""
	@echo "🔍 UTILITIES"
	@echo "  make health        - Check server health endpoint"
	@echo "  make swagger       - Open Swagger UI in default browser"
	@echo "  make commits       - View recent git commits via API"
	@echo ""
	@echo "🐳 DOCKER"
	@echo "  make docker-build-dev   - Build development Docker image"
	@echo "  make docker-build-prod  - Build production Docker image"
	@echo ""
	@echo "🚀 DEPLOYMENT"
	@echo "  make deploy-dev    - Deploy to development environment"
	@echo "  make deploy-prod   - Deploy to production environment"
	@echo "  make deploy-all    - Deploy to both environments"
	@echo ""
	@echo "🧪 TESTING DEPLOYMENTS"
	@echo "  make test-dev      - Test development deployment"
	@echo "  make test-prod     - Test production deployment"
	@echo ""
	@echo "📚 MORE INFO"
	@echo "  See README.md for detailed documentation"
	@echo ""

# ==============================================================================
# Development Commands
# ==============================================================================

run: ## Run API server locally (foreground)
	@echo "🚀 Starting API server on port $(PORT)..."
	@echo "   Press Ctrl+C to stop"
	@echo ""
	@go run ./cmd/api

build: ## Build Go binary
	@echo "🔨 Building Go binary..."
	@mkdir -p $(BIN_DIR)
	@go build -o $(BIN_DIR)/$(BINARY_NAME) ./cmd/api
	@echo "✅ Binary built: $(BIN_DIR)/$(BINARY_NAME)"

test: ## Run all tests
	@echo "🧪 Running tests..."
	@go test ./...

test-verbose: ## Run tests with verbose output
	@echo "🧪 Running tests (verbose)..."
	@go test ./... -v

clean: ## Clean build artifacts
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf $(BIN_DIR)/
	@go clean ./...
	@echo "✅ Clean complete"

# ==============================================================================
# Server Management Commands
# ==============================================================================

start: ## Start API server in background
	@echo "🚀 Starting API server in background on port $(PORT)..."
	@if [ "$(UNAME_S)" = "Windows" ]; then \
		start /B go run ./cmd/api > server.log 2>&1; \
	else \
		go run ./cmd/api > server.log 2>&1 & \
		echo $$! > server.pid; \
	fi
	@sleep 2
	@$(MAKE) status

stop: ## Stop API server
	@echo "🛑 Stopping API server on port $(PORT)..."
	@if [ "$(UNAME_S)" = "Windows" ]; then \
		PID=$$($(FIND_PID) 2>/dev/null); \
		if [ -n "$$PID" ]; then \
			$(KILL_CMD) $$PID 2>/dev/null && echo "✅ Server stopped (PID: $$PID)" || echo "❌ Failed to stop server"; \
		else \
			echo "ℹ️  No server running on port $(PORT)"; \
		fi; \
	else \
		if [ -f server.pid ]; then \
			kill -9 $$(cat server.pid) 2>/dev/null && rm server.pid && echo "✅ Server stopped" || echo "ℹ️  Server not running"; \
		fi; \
		PID=$$($(FIND_PID) 2>/dev/null); \
		if [ -n "$$PID" ]; then \
			$(KILL_CMD) $$PID 2>/dev/null && echo "✅ Server stopped (PID: $$PID)" || echo "❌ Failed to stop server"; \
		else \
			echo "ℹ️  No server running on port $(PORT)"; \
		fi; \
	fi

restart: stop start ## Restart API server

status: ## Check if server is running
	@echo "🔍 Checking server status on port $(PORT)..."
	@if [ "$(UNAME_S)" = "Windows" ]; then \
		PID=$$($(FIND_PID) 2>/dev/null); \
		if [ -n "$$PID" ]; then \
			echo "✅ Server is RUNNING (PID: $$PID)"; \
			curl -s http://localhost:$(PORT)/api/health > /dev/null 2>&1 && echo "✅ Health check: OK" || echo "⚠️  Health check: FAILED"; \
		else \
			echo "❌ Server is NOT running"; \
		fi; \
	else \
		PID=$$($(FIND_PID) 2>/dev/null); \
		if [ -n "$$PID" ]; then \
			echo "✅ Server is RUNNING (PID: $$PID)"; \
			curl -s http://localhost:$(PORT)/api/health > /dev/null 2>&1 && echo "✅ Health check: OK" || echo "⚠️  Health check: FAILED"; \
		else \
			echo "❌ Server is NOT running"; \
		fi; \
	fi

logs: ## Show server logs
	@if [ -f server.log ]; then \
		tail -f server.log; \
	else \
		echo "ℹ️  No server.log file found. Server may not be running in background."; \
	fi

# ==============================================================================
# Utility Commands
# ==============================================================================

health: ## Check server health endpoint
	@echo "🏥 Checking server health..."
	@curl -s http://localhost:$(PORT)/api/health | python -m json.tool 2>/dev/null || curl -s http://localhost:$(PORT)/api/health || echo "❌ Server not responding"

swagger: ## Open Swagger UI in default browser
	@echo "📚 Opening Swagger UI..."
	@if [ "$(UNAME_S)" = "Windows" ]; then \
		start http://localhost:$(PORT)/swagger; \
	else \
		open http://localhost:$(PORT)/swagger 2>/dev/null || xdg-open http://localhost:$(PORT)/swagger 2>/dev/null || echo "Visit: http://localhost:$(PORT)/swagger"; \
	fi

commits: ## View recent git commits via API
	@echo "📝 Fetching recent commits..."
	@curl -s http://localhost:$(PORT)/api/commits?limit=10 | python -m json.tool 2>/dev/null || curl -s http://localhost:$(PORT)/api/commits?limit=10 || echo "❌ Server not responding"

# ==============================================================================
# Docker Commands
# ==============================================================================

docker-build-dev: ## Build development Docker image
	@echo "🐳 Building development Docker image..."
	@docker build -f Dockerfile.dev -t $(IMAGE_DEV):latest .
	@docker tag $(IMAGE_DEV):latest $(IMAGE_DEV):$(shell date +%Y%m%d-%H%M%S)
	@echo "✅ Development image built: $(IMAGE_DEV):latest"

docker-build-prod: ## Build production Docker image
	@echo "🐳 Building production Docker image..."
	@docker build -f Dockerfile -t $(IMAGE_PROD):latest .
	@docker tag $(IMAGE_PROD):latest $(IMAGE_PROD):$(shell date +%Y%m%d-%H%M%S)
	@echo "✅ Production image built: $(IMAGE_PROD):latest"

# ==============================================================================
# Deployment Commands
# ==============================================================================

deploy-dev: docker-build-dev ## Deploy to development environment
	@echo "🚀 Deploying to development environment..."
	@bash scripts/deploy-dev.sh

deploy-prod: docker-build-prod ## Deploy to production environment
	@echo "🚀 Deploying to production environment..."
	@bash scripts/deploy-prod.sh

deploy-all: ## Deploy to both environments
	@echo "🚀 Deploying to all environments..."
	@bash scripts/deploy-all.sh

# ==============================================================================
# Testing Deployment Commands
# ==============================================================================

test-dev: ## Test development deployment
	@echo "🧪 Testing development deployment..."
	@bash scripts/test-dev.sh

test-prod: ## Test production deployment
	@echo "🧪 Testing production deployment..."
	@bash scripts/test-prod.sh

# ==============================================================================
# Windows PowerShell Alternatives (if needed)
# ==============================================================================

deploy-dev-ps1: ## Deploy to dev (PowerShell)
	@powershell -ExecutionPolicy Bypass -File scripts/deploy-dev.ps1

deploy-prod-ps1: ## Deploy to prod (PowerShell)
	@powershell -ExecutionPolicy Bypass -File scripts/deploy-prod.ps1
