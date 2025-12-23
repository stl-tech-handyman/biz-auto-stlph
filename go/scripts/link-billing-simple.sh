#!/bin/bash
# Simple script to link billing account

set -e

PROJECT_DEV="bizops360-dev"
PROJECT_PROD="bizops360-prod"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

echo ""
print_status "=== Link Billing Account ==="
echo ""
print_status "To find your billing account ID:"
print_status "1. Go to: https://console.cloud.google.com/billing"
print_status "2. Copy the Billing Account ID (format: XXXXXX-XXXXXX-XXXXXX)"
echo ""
read -p "Enter billing account ID: " BILLING_ACCOUNT

if [ -z "$BILLING_ACCOUNT" ]; then
    echo "Error: Billing account ID is required"
    exit 1
fi

# Link to dev
print_status "Linking billing to $PROJECT_DEV..."
if gcloud billing projects link "$PROJECT_DEV" --billing-account="$BILLING_ACCOUNT" 2>&1; then
    print_success "Billing linked to $PROJECT_DEV"
else
    print_warning "Failed to link (may already be linked)"
fi

# Ask about prod
echo ""
read -p "Link billing to $PROJECT_PROD? (y/n): " LINK_PROD
if [ "$LINK_PROD" = "y" ]; then
    print_status "Linking billing to $PROJECT_PROD..."
    if gcloud billing projects link "$PROJECT_PROD" --billing-account="$BILLING_ACCOUNT" 2>&1; then
        print_success "Billing linked to $PROJECT_PROD"
    else
        print_warning "Failed to link (may already be linked)"
    fi
fi

echo ""
print_success "Billing setup complete!"
print_status "Now you can run: bash scripts/complete-setup-dev.sh"

