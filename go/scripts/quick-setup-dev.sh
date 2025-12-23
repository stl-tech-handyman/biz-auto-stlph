#!/bin/bash
# Quick setup and deploy to bizops360-dev

set -e

PROJECT_ID="bizops360-dev"
REGION="us-central1"
SERVICE_NAME="bizops360-api-go-dev"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

# Check if project exists
if ! gcloud projects describe "$PROJECT_ID" &>/dev/null; then
    print_error "Project $PROJECT_ID does not exist!"
    print_status "Please run: bash scripts/setup-new-projects.sh"
    exit 1
fi

print_status "Project $PROJECT_ID exists, proceeding with deployment..."

# Set project
gcloud config set project "$PROJECT_ID"

# Check if secrets exist
if ! gcloud secrets describe svc-api-key-dev --project="$PROJECT_ID" &>/dev/null; then
    print_warning "Secret svc-api-key-dev not found!"
    print_status "Creating secret..."
    read -sp "Enter API key: " API_KEY
    echo ""
    echo -n "$API_KEY" | gcloud secrets create svc-api-key-dev \
        --data-file=- \
        --replication-policy="automatic" \
        --project="$PROJECT_ID"
    print_success "Secret created"
fi

# Run deployment
print_status "Starting deployment..."
bash scripts/deploy-dev.sh

