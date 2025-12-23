#!/bin/bash
# Deploy email API to Cloud Run (dev environment)

set -e

PROJECT_ID="bizops360-email-dev"
SERVICE_NAME="bizops360-email-api-dev"
REGION="us-central1"
IMAGE_NAME="gcr.io/${PROJECT_ID}/email-api"

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

print_status "Deploying email API to Cloud Run (dev)"
print_status "Project: $PROJECT_ID"
print_status "Service: $SERVICE_NAME"
print_status "Region: $REGION"

# Set project
gcloud config set project "$PROJECT_ID"

# Build Docker image
print_status "Building Docker image..."
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/../.." && pwd)"
cd "$PROJECT_ROOT" || exit 1
docker build -f go/Dockerfile.email-api.dev -t "$IMAGE_NAME:latest" .

# Push to Artifact Registry
print_status "Pushing image to Artifact Registry..."
docker push "$IMAGE_NAME:latest"

# Prepare secrets
SECRET_ARGS="GMAIL_CREDENTIALS_JSON=gmail-credentials-json:latest"

# Gmail FROM email (required for domain-wide delegation)
GMAIL_FROM="${GMAIL_FROM:-team@stlpartyhelpers.com}"
print_status "Gmail FROM: $GMAIL_FROM"

# Create API key secret if it doesn't exist
if ! gcloud secrets describe email-api-key-dev --project="$PROJECT_ID" >/dev/null 2>&1; then
    print_status "Creating email API key secret..."
    echo -n "email-api-key-dev-$(openssl rand -hex 16)" | gcloud secrets create email-api-key-dev \
        --data-file=- \
        --replication-policy="automatic" \
        --project="$PROJECT_ID"
    print_success "Email API key secret created"
else
    print_status "Email API key secret already exists"
fi

SECRET_ARGS="$SECRET_ARGS,SERVICE_API_KEY=email-api-key-dev:latest"

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
    --labels="env=dev,service=email-api,type=cloud-run"

# Get service URL
SERVICE_URL=$(gcloud run services describe "$SERVICE_NAME" \
    --project="$PROJECT_ID" \
    --region="$REGION" \
    --format="value(status.url)")

print_success "Deployment complete!"
print_status "Service URL: $SERVICE_URL"
print_status "Health check: $SERVICE_URL/api/health"

