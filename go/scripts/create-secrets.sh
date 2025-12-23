#!/bin/bash
# Create secrets in GCP projects

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
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }

# Function to create secret
create_secret() {
    local PROJECT_ID=$1
    local SECRET_NAME=$2
    local SECRET_VALUE=$3
    
    if gcloud secrets describe "$SECRET_NAME" --project="$PROJECT_ID" &>/dev/null; then
        print_warning "Secret $SECRET_NAME already exists in $PROJECT_ID"
        read -p "Update it? (y/n): " UPDATE
        if [ "$UPDATE" = "y" ]; then
            echo -n "$SECRET_VALUE" | gcloud secrets versions add "$SECRET_NAME" \
                --data-file=- \
                --project="$PROJECT_ID"
            print_success "Secret $SECRET_NAME updated"
        fi
    else
        echo -n "$SECRET_VALUE" | gcloud secrets create "$SECRET_NAME" \
            --data-file=- \
            --replication-policy="automatic" \
            --project="$PROJECT_ID"
        print_success "Secret $SECRET_NAME created"
    fi
}

# Create secrets for dev
print_status "=== Creating secrets for DEV ==="
gcloud config set project "$PROJECT_DEV"

read -sp "Enter API key for dev: " API_KEY_DEV
echo ""
create_secret "$PROJECT_DEV" "svc-api-key-dev" "$API_KEY_DEV"

read -p "Create Stripe secrets for dev? (y/n): " CREATE_STRIPE_DEV
if [ "$CREATE_STRIPE_DEV" = "y" ]; then
    read -sp "Enter Stripe test key: " STRIPE_TEST
    echo ""
    create_secret "$PROJECT_DEV" "stripe-secret-key-test" "$STRIPE_TEST"
    
    read -sp "Enter Stripe prod key (optional, press Enter to skip): " STRIPE_PROD
    echo ""
    if [ -n "$STRIPE_PROD" ]; then
        create_secret "$PROJECT_DEV" "stripe-secret-key-prod" "$STRIPE_PROD"
    fi
fi

# Create secrets for prod
echo ""
read -p "Create secrets for PROD? (y/n): " CREATE_PROD_SECRETS
if [ "$CREATE_PROD_SECRETS" = "y" ]; then
    print_status "=== Creating secrets for PROD ==="
    gcloud config set project "$PROJECT_PROD"
    
    read -sp "Enter API key for prod: " API_KEY_PROD
    echo ""
    create_secret "$PROJECT_PROD" "svc-api-key-prod" "$API_KEY_PROD"
    
    read -p "Create Stripe secrets for prod? (y/n): " CREATE_STRIPE_PROD
    if [ "$CREATE_STRIPE_PROD" = "y" ]; then
        read -sp "Enter Stripe test key: " STRIPE_TEST_PROD
        echo ""
        create_secret "$PROJECT_PROD" "stripe-secret-key-test" "$STRIPE_TEST_PROD"
        
        read -sp "Enter Stripe prod key: " STRIPE_PROD_PROD
        echo ""
        create_secret "$PROJECT_PROD" "stripe-secret-key-prod" "$STRIPE_PROD_PROD"
    fi
fi

print_success "Secrets created!"
print_status "You can now deploy: bash scripts/deploy-dev.sh"

