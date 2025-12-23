#!/bin/bash
# Automated script to delete stlph-dev project
# WARNING: This will delete everything without confirmation!

set -e

PROJECT_ID="stlph-dev"
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
    exit 1
fi

print_warning "⚠️  DELETING project: $PROJECT_ID"
print_warning "⚠️  This will delete ALL resources!"

# Set project
gcloud config set project "$PROJECT_ID"

# Delete Cloud Functions
print_status "Deleting Cloud Functions..."
FUNCTIONS=("geo" "health" "healthz" "leads")

for func in "${FUNCTIONS[@]}"; do
    print_status "  Deleting function: $func"
    gcloud functions delete "$func" --region="$REGION" --project="$PROJECT_ID" --gen2 --quiet 2>&1 || \
    gcloud functions delete "$func" --region="$REGION" --project="$PROJECT_ID" --quiet 2>&1 || \
    print_warning "    Function $func not found or already deleted"
done

print_success "Cloud Functions cleanup complete"

# Delete Cloud Run services (if any)
print_status "Checking for Cloud Run services..."
SERVICES=$(gcloud run services list --project="$PROJECT_ID" --format="value(metadata.name)" 2>/dev/null || echo "")

if [ -n "$SERVICES" ]; then
    for service in $SERVICES; do
        print_status "  Deleting service: $service"
        gcloud run services delete "$service" --region="$REGION" --project="$PROJECT_ID" --quiet 2>&1 || \
        print_warning "    Service $service not found or already deleted"
    done
fi

# Wait a bit for resources to be fully deleted
print_status "Waiting for resources to be fully deleted..."
sleep 10

# Delete project
print_status "Deleting project: $PROJECT_ID"
if gcloud projects delete "$PROJECT_ID" --quiet; then
    print_success "✅ Project $PROJECT_ID deleted successfully!"
else
    print_error "Failed to delete project $PROJECT_ID"
    print_status "You may need to:"
    print_status "  1. Wait a few minutes and try again"
    print_status "  2. Delete resources manually in Google Cloud Console"
    print_status "  3. Unlink billing account first"
    exit 1
fi

print_success "✅ Cleanup complete!"
print_status "Active projects remaining:"
gcloud projects list --filter="projectId:bizops360*" --format="table(projectId,name)"



