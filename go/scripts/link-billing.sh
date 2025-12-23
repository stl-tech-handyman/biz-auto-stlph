#!/bin/bash
# Link billing account to projects

set -e

PROJECT_DEV="bizops360-dev"
PROJECT_PROD="bizops360-prod"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Get billing account
print_status "Finding billing accounts..."
BILLING_ACCOUNTS=$(gcloud beta billing accounts list --format="value(name)" 2>/dev/null || echo "")

if [ -z "$BILLING_ACCOUNTS" ]; then
    print_error "No billing accounts found or command failed"
    print_status "Try: gcloud beta billing accounts list"
    read -p "Enter billing account ID manually: " BILLING_ACCOUNT
else
    echo "$BILLING_ACCOUNTS" | nl
    echo ""
    read -p "Select billing account number (or enter ID manually): " SELECTION
    
    if [[ "$SELECTION" =~ ^[0-9]+$ ]]; then
        BILLING_ACCOUNT=$(echo "$BILLING_ACCOUNTS" | sed -n "${SELECTION}p")
    else
        BILLING_ACCOUNT="$SELECTION"
    fi
fi

# Link to dev
print_status "Linking billing to $PROJECT_DEV..."
gcloud billing projects link "$PROJECT_DEV" --billing-account="$BILLING_ACCOUNT"
print_success "Billing linked to $PROJECT_DEV"

# Link to prod
read -p "Link billing to $PROJECT_PROD? (y/n): " LINK_PROD
if [ "$LINK_PROD" = "y" ]; then
    print_status "Linking billing to $PROJECT_PROD..."
    gcloud billing projects link "$PROJECT_PROD" --billing-account="$BILLING_ACCOUNT"
    print_success "Billing linked to $PROJECT_PROD"
fi

print_success "Billing setup complete!"
print_status "Now run: bash scripts/create-projects.sh (to enable APIs)"

