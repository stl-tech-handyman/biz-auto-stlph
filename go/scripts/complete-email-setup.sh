#!/bin/bash
# Complete email service setup (run after billing is linked)

set -e

EMAIL_PROJECT="bizops360-email-dev"
MAIN_PROJECT="bizops360-dev"
SECRET_NAME="gmail-credentials-json"
KEY_FILE="/tmp/email-service-key.json"

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

print_status "Completing email service setup..."

# Check billing
print_status "Checking billing..."
BILLING=$(gcloud projects describe "$EMAIL_PROJECT" --format="value(billingAccountName)" 2>/dev/null || echo "")
if [ -z "$BILLING" ]; then
    print_error "Billing account not linked!"
    print_status "Please link billing: https://console.cloud.google.com/billing/01C379-C9A8C1-3ED059/linkedprojects"
    exit 1
fi
print_success "Billing account linked"

# Enable APIs
print_status "Enabling required APIs..."
gcloud config set project "$EMAIL_PROJECT"
gcloud services enable secretmanager.googleapis.com gmail.googleapis.com 2>&1 | grep -v "already enabled" || true
print_success "APIs enabled"

# Check if key file exists
if [ ! -f "$KEY_FILE" ]; then
    print_warning "Key file not found at $KEY_FILE"
    print_status "Creating new service account key..."
    SERVICE_ACCOUNT="email-service-sa@${EMAIL_PROJECT}.iam.gserviceaccount.com"
    gcloud iam service-accounts keys create "$KEY_FILE" \
        --iam-account="$SERVICE_ACCOUNT" \
        --project="$EMAIL_PROJECT"
    print_success "Key created"
fi

# Create or update secret
print_status "Creating/updating secret..."
if gcloud secrets describe "$SECRET_NAME" --project="$EMAIL_PROJECT" &>/dev/null; then
    print_status "Secret exists, adding new version..."
    gcloud secrets versions add "$SECRET_NAME" \
        --data-file="$KEY_FILE" \
        --project="$EMAIL_PROJECT"
else
    gcloud secrets create "$SECRET_NAME" \
        --data-file="$KEY_FILE" \
        --replication-policy="automatic" \
        --project="$EMAIL_PROJECT"
fi
print_success "Secret created/updated"

# Grant access to main API service account
print_status "Granting access to main API service account..."
MAIN_PROJECT_NUMBER=$(gcloud projects describe "$MAIN_PROJECT" --format="value(projectNumber)" 2>/dev/null)
if [ -z "$MAIN_PROJECT_NUMBER" ]; then
    print_error "Could not get main project number"
    exit 1
fi

MAIN_SA="${MAIN_PROJECT_NUMBER}-compute@developer.gserviceaccount.com"
gcloud secrets add-iam-policy-binding "$SECRET_NAME" \
    --member="serviceAccount:${MAIN_SA}" \
    --role="roles/secretmanager.secretAccessor" \
    --project="$EMAIL_PROJECT" 2>&1 | grep -v "Updated IAM policy" || true
print_success "Access granted"

# Get service account client ID for domain-wide delegation
print_status "Getting service account client ID for domain-wide delegation..."
SERVICE_ACCOUNT="email-service-sa@${EMAIL_PROJECT}.iam.gserviceaccount.com"
CLIENT_ID=$(gcloud iam service-accounts describe "$SERVICE_ACCOUNT" --project="$EMAIL_PROJECT" --format="value(uniqueId)" 2>/dev/null || echo "")

print_success "Setup complete!"
echo ""
print_status "Next steps:"
echo ""
print_status "1. Configure Domain-Wide Delegation (if using Google Workspace):"
print_status "   - Go to: https://admin.google.com"
print_status "   - Security → API Controls → Domain-wide Delegation"
print_status "   - Add client ID: $CLIENT_ID"
print_status "   - OAuth Scopes: https://www.googleapis.com/auth/gmail.send"
echo ""
print_status "2. Deploy API with email support:"
print_status "   cd go && bash scripts/deploy-dev.sh"
echo ""
print_status "3. Test email endpoint:"
print_status "   curl -X POST 'https://bizops360-api-go-dev-gqqr4r256q-uc.a.run.app/api/email/test' \\"
print_status "     -H 'X-Api-Key: dev-api-key-12345' \\"
print_status "     -H 'Content-Type: application/json' \\"
print_status "     -d '{\"to\":\"test@example.com\",\"subject\":\"Test\",\"html\":\"<p>Test</p>\"}'"
echo ""
print_status "Service Account Client ID: $CLIENT_ID"
print_status "Secret: projects/${EMAIL_PROJECT}/secrets/${SECRET_NAME}"

