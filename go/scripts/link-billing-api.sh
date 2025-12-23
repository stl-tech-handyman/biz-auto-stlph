#!/bin/bash
# Link billing account via REST API

set -e

PROJECT_ID="bizops360-dev"
BILLING_ACCOUNT="01C379-C9A8C1-3ED059"

echo "Linking billing account $BILLING_ACCOUNT to project $PROJECT_ID..."

# Get access token
ACCESS_TOKEN=$(gcloud auth print-access-token 2>/dev/null)

if [ -z "$ACCESS_TOKEN" ]; then
    echo "Error: Could not get access token. Please run: gcloud auth login"
    exit 1
fi

# Link via REST API
RESPONSE=$(curl -s -X POST \
    "https://cloudbilling.googleapis.com/v1/projects/${PROJECT_ID}/billingInfo?billingAccountName=projects%2F${PROJECT_ID}%2FbillingAccounts%2F${BILLING_ACCOUNT}" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json" \
    -d "{\"billingAccountName\": \"billingAccounts/${BILLING_ACCOUNT}\"}")

# Try alternative endpoint
if echo "$RESPONSE" | grep -q "error"; then
    echo "Trying alternative method..."
    RESPONSE=$(curl -s -X PUT \
        "https://cloudbilling.googleapis.com/v1/projects/${PROJECT_ID}/billingInfo" \
        -H "Authorization: Bearer ${ACCESS_TOKEN}" \
        -H "Content-Type: application/json" \
        -d "{\"billingAccountName\": \"billingAccounts/${BILLING_ACCOUNT}\"}")
fi

echo "Response: $RESPONSE"

# Check if successful
if echo "$RESPONSE" | grep -q "billingAccountName"; then
    echo "✅ Billing account linked successfully!"
    exit 0
else
    echo "⚠️  May need to link manually via web console"
    echo "Go to: https://console.cloud.google.com/billing/${BILLING_ACCOUNT}/linkedprojects"
    exit 1
fi

