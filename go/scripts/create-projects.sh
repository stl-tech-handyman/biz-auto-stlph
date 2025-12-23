#!/bin/bash
# Create GCP projects: bizops360-dev and bizops360-prod

set -e

PROJECT_DEV="bizops360-dev"
PROJECT_PROD="bizops360-prod"
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

# Function to create a project
create_project() {
    local PROJECT_ID=$1
    local PROJECT_NAME=$2
    
    print_status "Creating project: $PROJECT_ID"
    
    if gcloud projects describe "$PROJECT_ID" &>/dev/null; then
        print_warning "Project $PROJECT_ID already exists"
        return 0
    fi
    
    gcloud projects create "$PROJECT_ID" --name="$PROJECT_NAME"
    
    if [ $? -eq 0 ]; then
        print_success "Project $PROJECT_ID created"
    else
        print_error "Failed to create project $PROJECT_ID"
        return 1
    fi
}

# Function to setup a project
setup_project() {
    local PROJECT_ID=$1
    
    print_status "Setting up project: $PROJECT_ID"
    gcloud config set project "$PROJECT_ID"
    
    # Enable APIs
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
    
    print_success "Project $PROJECT_ID setup complete!"
}

# Create dev project
print_status "=== Creating DEV project ==="
create_project "$PROJECT_DEV" "BizOps360 Development"
setup_project "$PROJECT_DEV"

echo ""
read -p "Create PROD project? (y/n): " CREATE_PROD

if [ "$CREATE_PROD" = "y" ]; then
    print_status "=== Creating PROD project ==="
    create_project "$PROJECT_PROD" "BizOps360 Production"
    setup_project "$PROJECT_PROD"
fi

print_success "Projects created!"
print_status "Next steps:"
print_status "  1. Link billing accounts (if not done automatically)"
print_status "  2. Create secrets: bash scripts/create-secrets.sh"
print_status "  3. Deploy: bash scripts/deploy-dev.sh"

