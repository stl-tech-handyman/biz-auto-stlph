#!/bin/bash
# Re-authenticate with Google Cloud

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

print_status "=== Google Cloud Re-authentication ==="
echo ""

# Check gcloud is installed
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI is not installed"
    exit 1
fi

# Step 1: User authentication
print_status "Step 1: Authenticating user account..."
print_status "This will open a browser for you to sign in"
gcloud auth login

if [ $? -eq 0 ]; then
    print_success "User authentication complete"
else
    print_error "User authentication failed"
    exit 1
fi

echo ""

# Step 2: Application default credentials
print_status "Step 2: Setting up application-default credentials..."
print_status "This is needed for Docker and local development"
gcloud auth application-default login

if [ $? -eq 0 ]; then
    print_success "Application-default credentials configured"
else
    print_error "Application-default credentials setup failed"
    exit 1
fi

echo ""

# Step 3: Configure Docker
print_status "Step 3: Configuring Docker authentication..."
gcloud auth configure-docker

if [ $? -eq 0 ]; then
    print_success "Docker authentication configured"
else
    print_warning "Docker authentication may need manual setup"
fi

echo ""

# Show current status
print_status "Current authentication status:"
echo ""
gcloud auth list
echo ""

# Show current project
CURRENT_PROJECT=$(gcloud config get-value project 2>/dev/null || echo "not set")
print_status "Current project: $CURRENT_PROJECT"

if [ "$CURRENT_PROJECT" != "bizops360-dev" ] && [ "$CURRENT_PROJECT" != "bizops360-prod" ]; then
    print_warning "Project is not set to bizops360-dev or bizops360-prod"
    print_status "Set it with: gcloud config set project bizops360-dev"
fi

echo ""
print_success "Re-authentication complete!"

