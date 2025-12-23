#!/bin/bash
# Auto-setup: Try to find and link billing account

set -e

PROJECT_DEV="bizops360-dev"

# Try to get billing from other projects
echo "Searching for billing account in existing projects..."
BILLING_ACCOUNT=""

# Check a few common project patterns
for proj in $(gcloud projects list --format="value(projectId)" 2>/dev/null | head -10); do
    BILLING=$(gcloud projects describe "$proj" --format="value(billingAccountName)" 2>/dev/null || echo "")
    if [ -n "$BILLING" ] && [ "$BILLING" != "" ]; then
        echo "Found billing account: $BILLING from project: $proj"
        BILLING_ACCOUNT="$BILLING"
        break
    fi
done

if [ -z "$BILLING_ACCOUNT" ]; then
    echo ""
    echo "⚠️  Could not automatically find billing account"
    echo ""
    echo "Please provide billing account ID:"
    echo "1. Go to: https://console.cloud.google.com/billing"
    echo "2. Copy the Billing Account ID"
    echo ""
    read -p "Enter billing account ID (or press Enter to skip): " BILLING_ACCOUNT
fi

if [ -n "$BILLING_ACCOUNT" ] && [ "$BILLING_ACCOUNT" != "" ]; then
    echo "Linking billing account $BILLING_ACCOUNT to $PROJECT_DEV..."
    if gcloud billing projects link "$PROJECT_DEV" --billing-account="$BILLING_ACCOUNT" 2>&1; then
        echo "✅ Billing linked successfully!"
        return 0
    else
        echo "⚠️  Failed to link billing (may already be linked or need permissions)"
    fi
else
    echo "⚠️  Skipping billing link - you'll need to do it manually"
    echo "Run: gcloud billing projects link bizops360-dev --billing-account=BILLING_ACCOUNT_ID"
fi

return 0

