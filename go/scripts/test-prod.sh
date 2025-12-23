#!/bin/bash
# Test Go API in Production Environment

set -e

PROJECT_ID="bizops360-prod"
REGION="us-central1"
SERVICE_NAME="bizops360-api-go-prod"

# Colors
BLUE='\033[0;34m'
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

# Get service URL
SERVICE_URL=$(gcloud run services describe "$SERVICE_NAME" \
    --project="$PROJECT_ID" \
    --region="$REGION" \
    --format="value(status.url)" 2>/dev/null)

if [ -z "$SERVICE_URL" ]; then
    print_error "Service $SERVICE_NAME not found in $PROJECT_ID"
    exit 1
fi

print_warning "⚠️  Testing PRODUCTION environment: $SERVICE_URL"
read -p "Continue? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    print_status "Test cancelled"
    exit 0
fi

# Get API key from Secret Manager
API_KEY=$(gcloud secrets versions access latest --secret="svc-api-key-prod" --project="$PROJECT_ID" 2>/dev/null)

if [ -z "$API_KEY" ]; then
    print_error "Could not retrieve API key from Secret Manager"
    exit 1
fi

# Test health endpoint (no auth)
print_status "Testing health endpoint..."
HEALTH_RESPONSE=$(curl -s -w "\n%{http_code}" "$SERVICE_URL/api/health")
HTTP_CODE=$(echo "$HEALTH_RESPONSE" | tail -n1)
BODY=$(echo "$HEALTH_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_success "Health check passed"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
    print_error "Health check failed: HTTP $HTTP_CODE"
    echo "$BODY"
    exit 1
fi

# Test estimate endpoint (requires auth)
print_status "Testing estimate endpoint..."
ESTIMATE_RESPONSE=$(curl -s -w "\n%{http_code}" \
    -X POST \
    -H "X-Api-Key: $API_KEY" \
    -H "Content-Type: application/json" \
    -d '{"eventDate":"2025-06-15","durationHours":4,"numHelpers":2}' \
    "$SERVICE_URL/api/estimate")
HTTP_CODE=$(echo "$ESTIMATE_RESPONSE" | tail -n1)
BODY=$(echo "$ESTIMATE_RESPONSE" | sed '$d')

if [ "$HTTP_CODE" = "200" ]; then
    print_success "Estimate endpoint passed"
    echo "$BODY" | jq '.' 2>/dev/null || echo "$BODY"
else
    print_error "Estimate endpoint failed: HTTP $HTTP_CODE"
    echo "$BODY"
    exit 1
fi

print_success "All tests passed!"

