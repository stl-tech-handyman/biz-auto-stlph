#!/bin/bash
# Setup Google Maps API project: bizops360-maps
# This project is dedicated to Google Maps API usage for event location address boxes

set -e

# Configuration
PROJECT_ID="bizops360-maps"
PROJECT_NAME="BizOps360 Maps API"
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

print_status "=== Setting up Google Maps API Project ==="
print_status "Project ID: $PROJECT_ID"
print_status "Project Name: $PROJECT_NAME"
echo ""

# Check if project already exists
if gcloud projects describe "$PROJECT_ID" &>/dev/null; then
    print_warning "Project $PROJECT_ID already exists, using existing project"
else
    # Create project
    print_status "Creating project $PROJECT_ID..."
    gcloud projects create "$PROJECT_ID" --name="$PROJECT_NAME" --set-as-default
    
    if [ $? -ne 0 ]; then
        print_error "Failed to create project $PROJECT_ID"
        exit 1
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
print_status "Enabling Google Maps APIs..."
MAPS_APIS=(
    "maps-javascript-api.googleapis.com"      # Maps JavaScript API (for address autocomplete)
    "geocoding-api.googleapis.com"            # Geocoding API (for address to coordinates)
    "places-api.googleapis.com"               # Places API (for place autocomplete)
    "secretmanager.googleapis.com"            # Secret Manager (to store API key)
    "cloudresourcemanager.googleapis.com"     # Cloud Resource Manager
    "iam.googleapis.com"                      # IAM
)

for api in "${MAPS_APIS[@]}"; do
    print_status "  Enabling $api..."
    gcloud services enable "$api" --project="$PROJECT_ID" || {
        print_warning "Failed to enable $api (may already be enabled)"
    }
done

print_success "All APIs enabled"

# Create API key
print_status "Creating Google Maps API key..."
print_status "This will create an unrestricted API key. You should restrict it after creation."
echo ""

# Create API key using gcloud
API_KEY_OUTPUT=$(gcloud alpha services api-keys create \
    --display-name="Maps API Key for Event Location" \
    --project="$PROJECT_ID" 2>&1)

if [ $? -ne 0 ]; then
    print_error "Failed to create API key via gcloud CLI"
    print_status "You can create it manually in Google Cloud Console:"
    print_status "  https://console.cloud.google.com/apis/credentials?project=$PROJECT_ID"
    print_status ""
    read -p "Enter your Google Maps API key (or press Enter to skip): " MANUAL_API_KEY
    
    if [ -z "$MANUAL_API_KEY" ]; then
        print_warning "Skipping API key creation. You can create it manually later."
        API_KEY=""
    else
        API_KEY="$MANUAL_API_KEY"
    fi
else
    # Extract API key from output
    API_KEY=$(echo "$API_KEY_OUTPUT" | grep -oP 'keyString: \K[^\s]+' || echo "")
    
    if [ -z "$API_KEY" ]; then
        # Try alternative method to get the key
        print_status "Retrieving API key..."
        API_KEY_ID=$(gcloud alpha services api-keys list --project="$PROJECT_ID" --format="value(name)" --filter="displayName:'Maps API Key for Event Location'" | head -n1)
        
        if [ -n "$API_KEY_ID" ]; then
            API_KEY=$(gcloud alpha services api-keys get-key-string "$API_KEY_ID" --project="$PROJECT_ID" --format="value(keyString)" 2>/dev/null || echo "")
        fi
    fi
    
    if [ -z "$API_KEY" ]; then
        print_warning "Could not automatically retrieve API key"
        print_status "Please get it from: https://console.cloud.google.com/apis/credentials?project=$PROJECT_ID"
        read -p "Enter your Google Maps API key: " API_KEY
    else
        print_success "API key created: ${API_KEY:0:20}..."
    fi
fi

# Save API key to Secret Manager if we have one
if [ -n "$API_KEY" ]; then
    print_status "Saving API key to Secret Manager..."
    
    SECRET_NAME="maps-api-key"
    
    if ! gcloud secrets describe "$SECRET_NAME" --project="$PROJECT_ID" &>/dev/null; then
        echo -n "$API_KEY" | gcloud secrets create "$SECRET_NAME" \
            --data-file=- \
            --replication-policy="automatic" \
            --project="$PROJECT_ID"
        print_success "Secret $SECRET_NAME created in Secret Manager"
    else
        echo -n "$API_KEY" | gcloud secrets versions add "$SECRET_NAME" \
            --data-file=- \
            --project="$PROJECT_ID"
        print_success "New version added to secret $SECRET_NAME"
    fi
fi

# Display summary
echo ""
print_success "=== Setup Complete ==="
echo ""
print_status "Project Details:"
print_status "  Project ID: $PROJECT_ID"
print_status "  Project Name: $PROJECT_NAME"
print_status "  Region: $REGION"
echo ""

if [ -n "$API_KEY" ]; then
    print_status "API Key Information:"
    print_status "  API Key: $API_KEY"
    print_status "  Secret Name: $SECRET_NAME"
    echo ""
    print_warning "IMPORTANT: Restrict your API key in Google Cloud Console:"
    print_status "  https://console.cloud.google.com/apis/credentials?project=$PROJECT_ID"
    print_status ""
    print_status "Recommended restrictions:"
    print_status "  1. Application restrictions: HTTP referrers (web sites)"
    print_status "     Add your website domains (e.g., *.yourdomain.com/*)"
    print_status "  2. API restrictions: Restrict to:"
    print_status "     - Maps JavaScript API"
    print_status "     - Geocoding API"
    print_status "     - Places API"
    echo ""
fi

print_status "Next Steps:"
print_status "  1. Restrict your API key (see above)"
print_status "  2. Add the API key to your website form"
print_status "  3. Test the address autocomplete functionality"
echo ""

print_status "To retrieve the API key later:"
print_status "  gcloud secrets versions access latest --secret=\"$SECRET_NAME\" --project=\"$PROJECT_ID\""
echo ""

print_status "To view API usage and billing:"
print_status "  https://console.cloud.google.com/apis/dashboard?project=$PROJECT_ID"
echo ""

