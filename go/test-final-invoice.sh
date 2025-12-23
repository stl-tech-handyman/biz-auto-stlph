#!/bin/bash
# Test script for final invoice functionality
# Make sure the server is running first!

API_KEY="${SERVICE_API_KEY:-test-api-key-12345}"
BASE_URL="http://localhost:8080"

echo "🧪 Testing Final Invoice Functionality"
echo "======================================"
echo ""
read -p "Enter your email address for testing: " TEST_EMAIL

# Step 1: Create final invoice
echo ""
echo "1️⃣ Creating final invoice..."
INVOICE_RESPONSE=$(curl -s -X POST "$BASE_URL/api/stripe/final-invoice" \
  -H "X-Api-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d "{
    \"email\": \"$TEST_EMAIL\",
    \"name\": \"Test Customer\",
    \"totalAmount\": 1000.0,
    \"depositPaid\": 400.0,
    \"currency\": \"usd\",
    \"description\": \"Final payment for test event\"
  }")

echo "$INVOICE_RESPONSE" | jq '.' 2>/dev/null || echo "$INVOICE_RESPONSE"

# Extract invoice URL
INVOICE_URL=$(echo "$INVOICE_RESPONSE" | jq -r '.invoice.url' 2>/dev/null)
REMAINING_BALANCE=$(echo "$INVOICE_RESPONSE" | jq -r '.details.remainingBalance' 2>/dev/null)

if [ "$INVOICE_URL" = "null" ] || [ -z "$INVOICE_URL" ]; then
  echo "❌ Failed to create invoice"
  exit 1
fi

echo ""
echo "✅ Invoice created: $INVOICE_URL"
echo ""

# Step 2: Send email
echo "2️⃣ Sending final invoice email..."
EMAIL_RESPONSE=$(curl -s -X POST "$BASE_URL/api/email/final-invoice" \
  -H "X-Api-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d "{
    \"name\": \"Test Customer\",
    \"email\": \"$TEST_EMAIL\",
    \"totalAmount\": 1000.0,
    \"depositPaid\": 400.0,
    \"remainingBalance\": $REMAINING_BALANCE,
    \"invoiceUrl\": \"$INVOICE_URL\"
  }")

echo "$EMAIL_RESPONSE" | jq '.' 2>/dev/null || echo "$EMAIL_RESPONSE"

echo ""
echo "✅ Test complete! Check your email: $TEST_EMAIL"
echo "   Invoice URL: $INVOICE_URL"

