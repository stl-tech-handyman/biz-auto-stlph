#!/bin/bash
# Script to safely delete stlph-dev project
# This will delete all Cloud Functions and then the project itself

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

print_warning "⚠️  WARNING: This will DELETE the entire project: $PROJECT_ID"
print_warning "⚠️  This action is IRREVERSIBLE!"
echo ""
read -p "Are you sure you want to continue? Type 'DELETE' to confirm: " CONFIRM

if [ "$CONFIRM" != "DELETE" ]; then
    print_error "Deletion cancelled"
    exit 1
fi

# Set project
gcloud config set project "$PROJECT_ID"

# List all Cloud Functions
print_status "Listing Cloud Functions..."
FUNCTIONS=$(gcloud functions list --project="$PROJECT_ID" --format="value(name)" 2>/dev/null || echo "")

if [ -n "$FUNCTIONS" ]; then
    print_warning "Found Cloud Functions in $PROJECT_ID:"
    echo "$FUNCTIONS"
    echo ""
    read -p "Delete these Cloud Functions? (y/n): " DELETE_FUNCTIONS
    
    if [ "$DELETE_FUNCTIONS" = "y" ]; then
        for func in $FUNCTIONS; do
            print_status "Deleting function: $func"
            # Extract function name from full path
            func_name=$(echo "$func" | awk -F'/' '{print $NF}')
            gcloud functions delete "$func_name" --region="$REGION" --project="$PROJECT_ID" --gen2 --quiet 2>&1 || \
            gcloud functions delete "$func_name" --region="$REGION" --project="$PROJECT_ID" --quiet 2>&1 || \
            print_warning "Could not delete $func_name (may not exist or already deleted)"
        done
        print_success "Cloud Functions deleted"
    fi
else
    print_success "No Cloud Functions found"
fi

# List Cloud Run services
print_status "Listing Cloud Run services..."
SERVICES=$(gcloud run services list --project="$PROJECT_ID" --format="value(metadata.name)" 2>/dev/null || echo "")

if [ -n "$SERVICES" ]; then
    print_warning "Found Cloud Run services in $PROJECT_ID:"
    echo "$SERVICES"
    echo ""
    read -p "Delete these Cloud Run services? (y/n): " DELETE_SERVICES
    
    if [ "$DELETE_SERVICES" = "y" ]; then
        for service in $SERVICES; do
            print_status "Deleting service: $service"
            gcloud run services delete "$service" --region="$REGION" --project="$PROJECT_ID" --quiet 2>&1 || \
            print_warning "Could not delete $service (may not exist or already deleted)"
        done
        print_success "Cloud Run services deleted"
    fi
else
    print_success "No Cloud Run services found"
fi

# Final confirmation
echo ""
print_warning "⚠️  FINAL WARNING: About to DELETE project: $PROJECT_ID"
print_warning "⚠️  This will delete ALL resources in the project!"
echo ""
read -p "Type 'DELETE PROJECT' to confirm: " FINAL_CONFIRM

if [ "$FINAL_CONFIRM" != "DELETE PROJECT" ]; then
    print_error "Project deletion cancelled"
    exit 1
fi

# Delete project
print_status "Deleting project: $PROJECT_ID"
gcloud projects delete "$PROJECT_ID" --quiet

if [ $? -eq 0 ]; then
    print_success "Project $PROJECT_ID deleted successfully!"
else
    print_error "Failed to delete project $PROJECT_ID"
    print_status "You may need to:"
    print_status "  1. Delete all resources manually"
    print_status "  2. Unlink billing account"
    print_status "  3. Try again"
    exit 1
fi

print_success "✅ Cleanup complete!"
print_status "Active projects remaining:"
gcloud projects list --filter="projectId:bizops360*" --format="table(projectId,name)"



