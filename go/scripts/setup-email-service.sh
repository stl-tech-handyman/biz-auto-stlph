#!/bin/bash
# Setup email service project and credentials

set -e

PROJECT_ID="bizops360-email-dev"
SERVICE_ACCOUNT_NAME="email-service-sa"
SECRET_NAME="gmail-credentials-json"

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

print_status "Setting up email service project: $PROJECT_ID"
gcloud config set project "$PROJECT_ID"

# Check if project exists
if ! gcloud projects describe "$PROJECT_ID" &>/dev/null; then
    print_error "Project $PROJECT_ID does not exist!"
    print_status "Creating project..."
    gcloud projects create "$PROJECT_ID" --name="BizOps360 Email Dev"
    print_success "Project created"
fi

# Link billing
print_status "Checking billing..."
BILLING=$(gcloud projects describe "$PROJECT_ID" --format="value(billingAccountName)" 2>/dev/null || echo "")
if [ -z "$BILLING" ]; then
    print_warning "Billing not linked. Please link manually:"
    print_status "gcloud billing projects link $PROJECT_ID --billing-account=01C379-C9A8C1-3ED059"
    read -p "Press Enter after billing is linked..."
fi

# Enable APIs
print_status "Enabling required APIs..."
APIS=(
    "gmail.googleapis.com"
    "secretmanager.googleapis.com"
    "iam.googleapis.com"
    "cloudresourcemanager.googleapis.com"
)

for api in "${APIS[@]}"; do
    print_status "  Enabling $api..."
    gcloud services enable "$api" --project="$PROJECT_ID" 2>&1 | grep -v "already enabled" || true
done

# Create service account
print_status "Creating service account..."
if gcloud iam service-accounts describe "${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com" --project="$PROJECT_ID" &>/dev/null; then
    print_warning "Service account already exists"
else
    gcloud iam service-accounts create "$SERVICE_ACCOUNT_NAME" \
        --display-name="Email Service Account" \
        --project="$PROJECT_ID"
    print_success "Service account created"
fi

SERVICE_ACCOUNT_EMAIL="${SERVICE_ACCOUNT_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"

# Grant Gmail API permissions
print_status "Granting Gmail API permissions..."
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
    --member="serviceAccount:${SERVICE_ACCOUNT_EMAIL}" \
    --role="roles/gmail.send" \
    --condition=None 2>&1 | grep -v "Updated IAM policy" || true

# Create key
print_status "Creating service account key..."
KEY_FILE="/tmp/${SERVICE_ACCOUNT_NAME}-key.json"
gcloud iam service-accounts keys create "$KEY_FILE" \
    --iam-account="$SERVICE_ACCOUNT_EMAIL" \
    --project="$PROJECT_ID"

print_success "Service account key created: $KEY_FILE"

# Create secret
print_status "Creating secret in Secret Manager..."
if gcloud secrets describe "$SECRET_NAME" --project="$PROJECT_ID" &>/dev/null; then
    print_warning "Secret already exists, updating..."
    gcloud secrets versions add "$SECRET_NAME" \
        --data-file="$KEY_FILE" \
        --project="$PROJECT_ID"
else
    gcloud secrets create "$SECRET_NAME" \
        --data-file="$KEY_FILE" \
        --replication-policy="automatic" \
        --project="$PROJECT_ID"
fi

print_success "Secret created: $SECRET_NAME"

# Grant access to main API service account
print_status "Granting access to main API service account..."
MAIN_PROJECT="bizops360-dev"
MAIN_PROJECT_NUMBER=$(gcloud projects describe "$MAIN_PROJECT" --format="value(projectNumber)" 2>/dev/null || echo "")

if [ -n "$MAIN_PROJECT_NUMBER" ]; then
    MAIN_SERVICE_ACCOUNT="${MAIN_PROJECT_NUMBER}-compute@developer.gserviceaccount.com"
    gcloud secrets add-iam-policy-binding "$SECRET_NAME" \
        --member="serviceAccount:${MAIN_SERVICE_ACCOUNT}" \
        --role="roles/secretmanager.secretAccessor" \
        --project="$PROJECT_ID"
    print_success "Access granted to main API service account"
else
    print_warning "Could not find main project service account"
fi

print_success "Setup complete!"
print_status "Next steps:"
print_status "1. Enable Gmail API for the service account in Google Cloud Console"
print_status "2. Update Go API deployment to use secret: GMAIL_CREDENTIALS_JSON=gmail-credentials-json:latest"
print_status "3. Note: You may need to use domain-wide delegation for Gmail API"
print_status ""
print_status "Key file saved to: $KEY_FILE"
print_status "Secret name: $SECRET_NAME"

