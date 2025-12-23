# Makefile for Go API Development and Deployment

.PHONY: help build build-dev build-prod test test-dev test-prod deploy-dev deploy-prod deploy-all run clean

# Configuration
PROJECT_DEV := bizops360-dev
PROJECT_PROD := bizops360-prod
REGION := us-central1
SERVICE_DEV := bizops360-api-go-dev
SERVICE_PROD := bizops360-api-go-prod
IMAGE_DEV := gcr.io/$(PROJECT_DEV)/$(SERVICE_DEV)
IMAGE_PROD := gcr.io/$(PROJECT_PROD)/$(SERVICE_PROD)

help: ## Show this help message
	@echo "Go API Makefile Commands"
	@echo ""
	@echo "Development:"
	@echo "  make run          - Run API locally"
	@echo "  make build        - Build Go binary"
	@echo "  make test         - Run all tests"
	@echo ""
	@echo "Building Images:"
	@echo "  make build-dev    - Build dev Docker image"
	@echo "  make build-prod   - Build prod Docker image"
	@echo ""
	@echo "Deployment:"
	@echo "  make deploy-dev   - Deploy to dev environment"
	@echo "  make deploy-prod  - Deploy to prod environment"
	@echo "  make deploy-all   - Deploy to both environments"
	@echo ""
	@echo "Testing:"
	@echo "  make test-dev     - Test dev deployment"
	@echo "  make test-prod    - Test prod deployment"
	@echo ""
	@echo "Utilities:"
	@echo "  make clean        - Clean build artifacts"

run: ## Run API locally
	@echo "Running Go API locally..."
	@cd cmd/api && go run main.go

build: ## Build Go binary
	@echo "Building Go binary..."
	@go build -o bin/server ./cmd/api

test: ## Run all tests
	@echo "Running tests..."
	@go test ./... -v

test-dev: ## Test dev deployment
	@bash scripts/test-dev.sh

test-prod: ## Test prod deployment
	@bash scripts/test-prod.sh

build-dev: ## Build dev Docker image
	@echo "Building dev Docker image..."
	@docker build -f Dockerfile.dev -t $(IMAGE_DEV):latest .
	@docker tag $(IMAGE_DEV):latest $(IMAGE_DEV):$(shell date +%Y%m%d-%H%M%S)

build-prod: ## Build prod Docker image
	@echo "Building prod Docker image..."
	@docker build -f Dockerfile.prod -t $(IMAGE_PROD):latest .
	@docker tag $(IMAGE_PROD):latest $(IMAGE_PROD):$(shell date +%Y%m%d-%H%M%S)

deploy-dev: build-dev ## Deploy to dev environment
	@bash scripts/deploy-dev.sh

deploy-prod: build-prod ## Deploy to prod environment
	@bash scripts/deploy-prod.sh

deploy-all: ## Deploy to both environments
	@bash scripts/deploy-all.sh

clean: ## Clean build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf bin/
	@go clean ./...

# Windows PowerShell alternatives
deploy-dev-ps1: ## Deploy to dev (PowerShell)
	@powershell -ExecutionPolicy Bypass -File scripts/deploy-dev.ps1

deploy-prod-ps1: ## Deploy to prod (PowerShell)
	@powershell -ExecutionPolicy Bypass -File scripts/deploy-prod.ps1

