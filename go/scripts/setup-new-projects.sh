#!/bin/bash
# Setup new GCP projects: bizops360-dev and bizops360-prod

set -e

# Configuration
PROJECT_DEV="bizops360-dev"
PROJECT_PROD="bizops360-prod"
REGION="us-central1"
BILLING_ACCOUNT="" # Will prompt if not set

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

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI is not installed"
    exit 1
fi

# Get billing account if not set
if [ -z "$BILLING_ACCOUNT" ]; then
    print_status "To find your billing account, run: gcloud beta billing accounts list"
    echo ""
    read -p "Enter billing account ID (e.g., 01ABCD-23EFGH-456789): " BILLING_ACCOUNT
fi

# Function to create and setup a project
setup_project() {
    local PROJECT_ID=$1
    local PROJECT_NAME=$2
    
    print_status "Setting up project: $PROJECT_ID"
    
    # Check if project already exists
    if gcloud projects describe "$PROJECT_ID" &>/dev/null; then
        print_warning "Project $PROJECT_ID already exists, skipping creation"
    else
        # Create project
        print_status "Creating project $PROJECT_ID..."
        gcloud projects create "$PROJECT_ID" --name="$PROJECT_NAME" --set-as-default
        
        if [ $? -ne 0 ]; then
            print_error "Failed to create project $PROJECT_ID"
            return 1
        fi
        print_success "Project $PROJECT_ID created"
    fi
    
    # Link billing account
    print_status "Linking billing account..."
    gcloud billing projects link "$PROJECT_ID" --billing-account="$BILLING_ACCOUNT" || {
        print_warning "Failed to link billing account (may already be linked or require permissions)"
    }
    
    # Set as current project
    gcloud config set project "$PROJECT_ID"
    
    # Enable required APIs
    print_status "Enabling required APIs..."
    local APIS=(
        "cloudbuild.googleapis.com"
        "run.googleapis.com"
        "artifactregistry.googleapis.com"
        "secretmanager.googleapis.com"
        "cloudresourcemanager.googleapis.com"
        "iam.googleapis.com"
    )
    
    for api in "${APIS[@]}"; do
        print_status "  Enabling $api..."
        gcloud services enable "$api" --project="$PROJECT_ID" || {
            print_warning "Failed to enable $api (may already be enabled)"
        }
    done
    
    # Create Artifact Registry repository
    print_status "Creating Artifact Registry repository..."
    gcloud artifacts repositories create "$PROJECT_ID" \
        --repository-format=docker \
        --location="$REGION" \
        --description="Docker repository for $PROJECT_NAME" \
        --project="$PROJECT_ID" 2>/dev/null || {
        print_warning "Artifact Registry repository may already exist or creation failed"
    }
    
    # Create secrets
    print_status "Creating secrets..."
    
    # API Key secret
    if ! gcloud secrets describe svc-api-key-dev --project="$PROJECT_ID" &>/dev/null; then
        print_status "Creating svc-api-key-dev secret..."
        read -sp "Enter API key for dev (will be hidden): " API_KEY
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
    print_status "Creating Stripe secrets (optional)..."
    read -p "Do you want to create Stripe secrets? (y/n): " CREATE_STRIPE
    
    if [ "$CREATE_STRIPE" = "y" ]; then
        # Stripe test key
        if ! gcloud secrets describe stripe-secret-key-test --project="$PROJECT_ID" &>/dev/null; then
            read -sp "Enter Stripe test key: " STRIPE_TEST_KEY
            echo ""
            echo -n "$STRIPE_TEST_KEY" | gcloud secrets create stripe-secret-key-test \
                --data-file=- \
                --replication-policy="automatic" \
                --project="$PROJECT_ID"
            print_success "Secret stripe-secret-key-test created"
        fi
        
        # Stripe prod key
        if ! gcloud secrets describe stripe-secret-key-prod --project="$PROJECT_ID" &>/dev/null; then
            read -sp "Enter Stripe prod key: " STRIPE_PROD_KEY
            echo ""
            echo -n "$STRIPE_PROD_KEY" | gcloud secrets create stripe-secret-key-prod \
                --data-file=- \
                --replication-policy="automatic" \
                --project="$PROJECT_ID"
            print_success "Secret stripe-secret-key-prod created"
        fi
    fi
    
    print_success "Project $PROJECT_ID setup complete!"
    echo ""
}

# Setup dev project
print_status "=== Setting up DEV project ==="
setup_project "$PROJECT_DEV" "BizOps360 Development"

echo ""
read -p "Do you want to setup PROD project now? (y/n): " SETUP_PROD

if [ "$SETUP_PROD" = "y" ]; then
    print_status "=== Setting up PROD project ==="
    setup_project "$PROJECT_PROD" "BizOps360 Production"
fi

print_success "All projects setup complete!"
print_status "Next steps:"
print_status "  1. Run: cd go && bash scripts/deploy-dev.sh"
print_status "  2. Or use: make deploy-dev"

