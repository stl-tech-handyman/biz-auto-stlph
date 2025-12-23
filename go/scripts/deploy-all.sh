#!/bin/bash
# Deploy Go API to both DEV and PROD environments

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

print_status "Deploying Go API to both environments..."

# Deploy to DEV first
print_status "Step 1/2: Deploying to DEV..."
cd "$PROJECT_ROOT"
bash "$SCRIPT_DIR/deploy-dev.sh"

print_warning "Waiting 10 seconds before deploying to PROD..."
sleep 10

# Deploy to PROD
print_status "Step 2/2: Deploying to PROD..."
bash "$SCRIPT_DIR/deploy-prod.sh"

print_success "All deployments complete!"

