#!/bin/bash
# Update Cloud Run service with new secrets (without rebuilding)

set -e

PROJECT_ID="bizops360-dev"
SERVICE_NAME="bizops360-api-go-dev"
REGION="us-central1"
GMAIL_FROM="team@stlpartyhelpers.com"

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

print_status "Updating Cloud Run service with new secrets..."
print_status "Project: $PROJECT_ID"
print_status "Service: $SERVICE_NAME"

# Set project
gcloud config set project "$PROJECT_ID"

# Prepare secrets
SECRET_ARGS="SERVICE_API_KEY=svc-api-key-dev:latest,STRIPE_SECRET_KEY_TEST=stripe-secret-key-test:latest,STRIPE_SECRET_KEY_PROD=stripe-secret-key-prod:latest"

# Check if Gmail credentials exist in main project
if gcloud secrets describe gmail-credentials-json --project="$PROJECT_ID" >/dev/null 2>&1; then
    SECRET_ARGS="$SECRET_ARGS,GMAIL_CREDENTIALS_JSON=gmail-credentials-json:latest"
    print_success "Gmail credentials found in $PROJECT_ID"
else
    print_warning "Gmail credentials not found in $PROJECT_ID (email features may not work)"
fi

# Update Cloud Run service
print_status "Updating service..."
if gcloud run services update "$SERVICE_NAME" \
    --project="$PROJECT_ID" \
    --region="$REGION" \
    --update-secrets="$SECRET_ARGS" \
    --update-env-vars="ENV=dev,LOG_LEVEL=debug,CONFIG_DIR=/app/config,TEMPLATES_DIR=/app/templates,GMAIL_FROM=$GMAIL_FROM"; then
    
    print_success "Service updated successfully!"
    
    SERVICE_URL=$(gcloud run services describe "$SERVICE_NAME" \
        --project="$PROJECT_ID" \
        --region="$REGION" \
        --format="value(status.url)")
    
    print_status "Service URL: $SERVICE_URL"
    print_status "Health check: $SERVICE_URL/api/health"
else
    print_error "Failed to update service"
    exit 1
fi

