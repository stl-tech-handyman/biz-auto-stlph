#!/bin/bash
# Get Google Maps API key from Secret Manager

set -e

PROJECT_ID="bizops360-maps"
SECRET_NAME="maps-api-key"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }

print_status "Retrieving Google Maps API key from project: $PROJECT_ID"

API_KEY=$(gcloud secrets versions access latest --secret="$SECRET_NAME" --project="$PROJECT_ID" 2>/dev/null)

if [ $? -eq 0 ] && [ -n "$API_KEY" ]; then
    echo ""
    print_success "API Key retrieved:"
    echo "$API_KEY"
    echo ""
    print_status "Copy this key and add it to your website form configuration."
else
    echo "Error: Could not retrieve API key"
    echo "Make sure:"
    echo "  1. Project $PROJECT_ID exists"
    echo "  2. Secret $SECRET_NAME exists"
    echo "  3. You have permissions to access secrets"
    exit 1
fi

