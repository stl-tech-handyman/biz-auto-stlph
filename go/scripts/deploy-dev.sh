#!/bin/bash
# Deploy Go API to Development Environment

set -e

# Configuration
PROJECT_ID="bizops360-dev"
REGION="us-central1"
SERVICE_NAME="bizops360-api-go-dev"
IMAGE_NAME="gcr.io/${PROJECT_ID}/${SERVICE_NAME}"

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

# Check prerequisites
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI is not installed"
    exit 1
fi

if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed"
    exit 1
fi

print_status "Deploying Go API to DEV environment..."
print_status "Project: $PROJECT_ID"
print_status "Service: $SERVICE_NAME"
print_status "Region: $REGION"

# Set project
gcloud config set project "$PROJECT_ID"

# Build Docker image
print_status "Building Docker image..."
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_ROOT"
docker build -f go/Dockerfile.dev -t "$IMAGE_NAME:latest" .

# Push to Artifact Registry
print_status "Pushing image to Artifact Registry..."
docker push "$IMAGE_NAME:latest"

# Prepare secrets and env vars
SECRET_ARGS="SERVICE_API_KEY=svc-api-key-dev:latest"
ENV_VARS=""

# Gmail credentials - check if it exists in main project or email project
EMAIL_PROJECT="bizops360-email-dev"
if gcloud secrets describe gmail-credentials-json --project="$PROJECT_ID" >/dev/null 2>&1; then
    SECRET_ARGS="$SECRET_ARGS,GMAIL_CREDENTIALS_JSON=gmail-credentials-json:latest"
    print_status "Gmail credentials found in $PROJECT_ID"
else
    # Check if it exists in email project - use cross-project reference
    if gcloud secrets describe gmail-credentials-json --project="$EMAIL_PROJECT" >/dev/null 2>&1; then
        # Use cross-project secret reference
        SECRET_ARGS="$SECRET_ARGS,GMAIL_CREDENTIALS_JSON=projects/$EMAIL_PROJECT/secrets/gmail-credentials-json:latest"
        print_status "Gmail credentials found in $EMAIL_PROJECT (using cross-project reference)"
    else
        print_warning "Gmail credentials not found (email features may not work)"
    fi
fi

# Gmail FROM email
GMAIL_FROM="${GMAIL_FROM:-team@stlpartyhelpers.com}"
print_status "Gmail FROM: $GMAIL_FROM"

# Stripe secrets are optional
if gcloud secrets describe stripe-secret-key-test --project="$PROJECT_ID" >/dev/null 2>&1; then
    SECRET_ARGS="$SECRET_ARGS,STRIPE_SECRET_KEY_TEST=stripe-secret-key-test:latest"
    print_status "Stripe test secret found"
else
    print_warning "Secret stripe-secret-key-test not found (optional, Stripe features may not work)"
fi

if gcloud secrets describe stripe-secret-key-prod --project="$PROJECT_ID" >/dev/null 2>&1; then
    SECRET_ARGS="$SECRET_ARGS,STRIPE_SECRET_KEY_PROD=stripe-secret-key-prod:latest"
    print_status "Stripe prod secret found"
fi

# Deploy to Cloud Run
print_status "Deploying to Cloud Run..."
gcloud run deploy "$SERVICE_NAME" \
    --project="$PROJECT_ID" \
    --region="$REGION" \
    --image="$IMAGE_NAME:latest" \
    --platform=managed \
    --allow-unauthenticated \
    --port=8080 \
    --memory=512Mi \
    --cpu=1 \
    --min-instances=0 \
    --max-instances=5 \
    --concurrency=80 \
    --timeout=300 \
    --set-secrets="$SECRET_ARGS" \
    --set-env-vars="ENV=dev,LOG_LEVEL=debug,CONFIG_DIR=/app/config,TEMPLATES_DIR=/app/templates,GMAIL_FROM=${GMAIL_FROM}" \
    --labels="env=dev,service=api-go,type=cloud-run"

# Get service URL
SERVICE_URL=$(gcloud run services describe "$SERVICE_NAME" \
    --project="$PROJECT_ID" \
    --region="$REGION" \
    --format="value(status.url)")

print_success "Deployment complete!"
print_status "Service URL: $SERVICE_URL"
print_status "Health check: $SERVICE_URL/api/health"
print_status "Stripe endpoint: $SERVICE_URL/api/stripe/deposit/calculate"

