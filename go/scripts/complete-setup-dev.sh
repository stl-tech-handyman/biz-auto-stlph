#!/bin/bash
# Complete setup for bizops360-dev: Enable APIs, create secrets, deploy

set -e

PROJECT_ID="bizops360-dev"
REGION="us-central1"

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
    print_status "Create it first: gcloud projects create $PROJECT_ID --name='BizOps360 Development'"
    exit 1
fi

gcloud config set project "$PROJECT_ID"

# Check billing
print_status "Checking billing status..."
BILLING_ENABLED=$(gcloud billing projects describe "$PROJECT_ID" --format="value(billingAccountName)" 2>/dev/null || echo "")

if [ -z "$BILLING_ENABLED" ]; then
    print_warning "Billing account check failed (may still be linking)"
    print_status "Continuing with setup - billing should be linked..."
    sleep 2
fi

# Enable APIs
print_status "Enabling required APIs..."
APIS=(
    "cloudbuild.googleapis.com"
    "run.googleapis.com"
    "artifactregistry.googleapis.com"
    "secretmanager.googleapis.com"
    "cloudresourcemanager.googleapis.com"
    "iam.googleapis.com"
)

for api in "${APIS[@]}"; do
    print_status "  Enabling $api..."
    gcloud services enable "$api" --project="$PROJECT_ID" 2>&1 | grep -v "already enabled" || true
done

# Create Artifact Registry
print_status "Creating Artifact Registry repository..."
if ! gcloud artifacts repositories describe "$PROJECT_ID" --location="$REGION" --project="$PROJECT_ID" &>/dev/null; then
    gcloud artifacts repositories create "$PROJECT_ID" \
        --repository-format=docker \
        --location="$REGION" \
        --description="Docker repository for $PROJECT_ID" \
        --project="$PROJECT_ID"
    print_success "Artifact Registry repository created"
else
    print_warning "Artifact Registry repository already exists"
fi

# Configure Docker auth
print_status "Configuring Docker authentication..."
gcloud auth configure-docker "${REGION}-docker.pkg.dev" --quiet

# Create secrets
print_status "Creating secrets..."

# API Key
if ! gcloud secrets describe svc-api-key-dev --project="$PROJECT_ID" &>/dev/null; then
    print_status "Creating svc-api-key-dev..."
    read -sp "Enter API key: " API_KEY
    echo ""
    echo -n "$API_KEY" | gcloud secrets create svc-api-key-dev \
        --data-file=- \
        --replication-policy="automatic" \
        --project="$PROJECT_ID"
    print_success "Secret svc-api-key-dev created"
else
    print_warning "Secret svc-api-key-dev already exists"
fi

# Stripe secrets (optional)
read -p "Create Stripe secrets? (y/n): " CREATE_STRIPE
if [ "$CREATE_STRIPE" = "y" ]; then
    # Stripe test
    if ! gcloud secrets describe stripe-secret-key-test --project="$PROJECT_ID" &>/dev/null; then
        read -sp "Enter Stripe test key: " STRIPE_TEST
        echo ""
        echo -n "$STRIPE_TEST" | gcloud secrets create stripe-secret-key-test \
            --data-file=- \
            --replication-policy="automatic" \
            --project="$PROJECT_ID"
        print_success "Secret stripe-secret-key-test created"
    fi
    
    # Stripe prod (optional)
    read -p "Create Stripe prod secret? (y/n): " CREATE_STRIPE_PROD
    if [ "$CREATE_STRIPE_PROD" = "y" ]; then
        if ! gcloud secrets describe stripe-secret-key-prod --project="$PROJECT_ID" &>/dev/null; then
            read -sp "Enter Stripe prod key: " STRIPE_PROD
            echo ""
            echo -n "$STRIPE_PROD" | gcloud secrets create stripe-secret-key-prod \
                --data-file=- \
                --replication-policy="automatic" \
                --project="$PROJECT_ID"
            print_success "Secret stripe-secret-key-prod created"
        fi
    fi
fi

print_success "Setup complete!"
print_status "Ready to deploy!"
echo ""
read -p "Deploy now? (y/n): " DEPLOY_NOW

if [ "$DEPLOY_NOW" = "y" ]; then
    print_status "Starting deployment..."
    bash scripts/deploy-dev.sh
fi

