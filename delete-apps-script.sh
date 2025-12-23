#!/bin/bash
# Script to delete Google Apps Script projects via API
# 
# Prerequisites:
# 1. Enable Apps Script API: https://script.google.com/home/usersettings
# 2. Authenticate: gcloud auth login
# 3. Get access token: gcloud auth print-access-token

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

print_warning "This script will help you delete Apps Script projects"
print_status "You need the SCRIPT ID (not GCP project ID) to delete"
echo ""

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI is not installed"
    exit 1
fi

# Get access token
print_status "Getting access token..."
ACCESS_TOKEN=$(gcloud auth print-access-token 2>/dev/null || echo "")

if [ -z "$ACCESS_TOKEN" ]; then
    print_error "Failed to get access token. Run: gcloud auth login"
    exit 1
fi

print_success "Access token obtained"

echo ""
print_status "To delete an Apps Script project:"
echo ""
echo "1. Go to https://script.google.com"
echo "2. Open the project you want to delete"
echo "3. Go to Project Settings (gear icon)"
echo "4. Copy the 'Script ID'"
echo "5. Run this command:"
echo ""
echo "   curl -X DELETE \\"
echo "     \"https://script.googleapis.com/v1/projects/YOUR_SCRIPT_ID\" \\"
echo "     -H \"Authorization: Bearer \$(gcloud auth print-access-token)\""
echo ""
echo "Or use the web interface (easier):"
echo "   https://script.google.com → Select project → ⋮ menu → Delete project"
echo ""

# List available scripts
print_status "Your current Apps Script projects:"
clasp list 2>/dev/null || print_warning "Could not list scripts (clasp may not be configured)"



