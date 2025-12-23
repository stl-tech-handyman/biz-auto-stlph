#!/bin/bash
# Grant Cloud Run service account access to secrets

set -e

PROJECT_ID="bizops360-dev"

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

print_status "Granting Cloud Run service account access to secrets..."
print_status "Project: $PROJECT_ID"

# Get project number and service account
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format="value(projectNumber)")
SERVICE_ACCOUNT="${PROJECT_NUMBER}-compute@developer.gserviceaccount.com"

print_status "Service Account: $SERVICE_ACCOUNT"

# Grant access to each secret
secrets=("svc-api-key-dev" "stripe-secret-key-test" "stripe-secret-key-prod" "gmail-credentials-json")

for secret in "${secrets[@]}"; do
    print_status "Granting access to secret: $secret"
    
    if gcloud secrets describe "$secret" --project="$PROJECT_ID" >/dev/null 2>&1; then
        if gcloud secrets add-iam-policy-binding "$secret" \
            --member="serviceAccount:$SERVICE_ACCOUNT" \
            --role="roles/secretmanager.secretAccessor" \
            --project="$PROJECT_ID" >/dev/null 2>&1; then
            print_success "Access granted to $secret"
        else
            print_warning "Failed to grant access to $secret (may already have access)"
        fi
    else
        print_warning "Secret $secret not found, skipping"
    fi
done

print_success "Permission update complete!"

