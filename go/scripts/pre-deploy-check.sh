#!/bin/bash
# Pre-deployment checks for Go API

set -e

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

print_status "Running pre-deployment checks..."

# Check gcloud is installed
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI is not installed"
    exit 1
fi
print_success "gcloud CLI found"

# Check Docker is installed
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed"
    exit 1
fi
print_success "Docker found"

# Check Docker is running
if ! docker ps &> /dev/null; then
    print_error "Docker is not running"
    exit 1
fi
print_success "Docker is running"

# Check gcloud authentication
ACCOUNTS=$(gcloud auth list --filter=status:ACTIVE --format="value(account)" 2>/dev/null)
if [ -z "$ACCOUNTS" ]; then
    print_error "No active gcloud account. Run: gcloud auth login"
    exit 1
fi
print_success "gcloud authenticated as: $ACCOUNTS"

# Check application-default credentials
if ! gcloud auth application-default print-access-token &>/dev/null; then
    print_warning "Application-default credentials not set"
    print_status "Run: gcloud auth application-default login"
    print_status "This is needed for Docker to push to Artifact Registry"
    read -p "Do you want to set it up now? (yes/no): " confirm
    if [ "$confirm" = "yes" ]; then
        gcloud auth application-default login
    else
        print_error "Cannot proceed without application-default credentials"
        exit 1
    fi
fi
print_success "Application-default credentials configured"

# Check projects exist
print_status "Checking GCP projects..."

PROJECT_DEV="bizops360-dev"
PROJECT_PROD="bizops360-prod"

if ! gcloud projects describe "$PROJECT_DEV" &>/dev/null; then
    print_error "Project $PROJECT_DEV not found or not accessible"
    exit 1
fi
print_success "Project $PROJECT_DEV accessible"

if ! gcloud projects describe "$PROJECT_PROD" &>/dev/null; then
    print_warning "Project $PROJECT_PROD not found or not accessible"
    print_status "You can still deploy to dev, but prod deployment will fail"
else
    print_success "Project $PROJECT_PROD accessible"
fi

# Check required APIs are enabled
print_status "Checking required APIs..."

APIS=("run.googleapis.com" "artifactregistry.googleapis.com" "secretmanager.googleapis.com")

for api in "${APIS[@]}"; do
    if gcloud services list --enabled --project="$PROJECT_DEV" --filter="name:$api" --format="value(name)" | grep -q "$api"; then
        print_success "API $api enabled in $PROJECT_DEV"
    else
        print_warning "API $api not enabled in $PROJECT_DEV"
        print_status "Enable with: gcloud services enable $api --project=$PROJECT_DEV"
    fi
done

# Check secrets exist
print_status "Checking secrets..."

if gcloud secrets describe svc-api-key-dev --project="$PROJECT_DEV" &>/dev/null; then
    print_success "Secret svc-api-key-dev exists"
else
    print_warning "Secret svc-api-key-dev not found in $PROJECT_DEV"
    print_status "Create with: echo -n 'your-key' | gcloud secrets create svc-api-key-dev --data-file=- --project=$PROJECT_DEV"
fi

if gcloud secrets describe stripe-secret-key-test --project="$PROJECT_DEV" &>/dev/null; then
    print_success "Secret stripe-secret-key-test exists"
else
    print_warning "Secret stripe-secret-key-test not found in $PROJECT_DEV"
fi

# Check Docker can authenticate with GCR
print_status "Checking Docker authentication..."
if gcloud auth configure-docker &>/dev/null; then
    print_success "Docker authentication configured"
else
    print_warning "Docker authentication may need setup"
    print_status "Run: gcloud auth configure-docker"
fi

print_success "Pre-deployment checks complete!"
print_status "You can now run: ./scripts/deploy-dev.sh"





