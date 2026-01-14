#!/bin/bash

# Test Email Analysis Service
# Usage: ./test-email-analysis.sh [api-key] [max-emails]

API_KEY=${1:-"your-api-key-here"}
MAX_EMAILS=${2:-50}
BASE_URL=${3:-"http://localhost:8080"}

echo "Testing Email Analysis Service"
echo "=============================="
echo "API Key: ${API_KEY:0:10}..."
echo "Max Emails: $MAX_EMAILS"
echo "Base URL: $BASE_URL"
echo ""

# Test 1: Analyze emails
echo "1. Starting email analysis..."
RESPONSE=$(curl -s -X POST "$BASE_URL/api/email-analysis/analyze" \
  -H "X-API-Key: $API_KEY" \
  -H "Content-Type: application/json" \
  -d "{
    \"max_emails\": $MAX_EMAILS,
    \"query\": \"from:zapier.com OR subject:\\\"New Lead\\\"\",
    \"resume\": false
  }")

echo "Response:"
echo "$RESPONSE" | jq '.' 2>/dev/null || echo "$RESPONSE"
echo ""

# Extract spreadsheet ID
SPREADSHEET_ID=$(echo "$RESPONSE" | jq -r '.spreadsheet_id' 2>/dev/null)

if [ "$SPREADSHEET_ID" != "null" ] && [ -n "$SPREADSHEET_ID" ]; then
  echo "Spreadsheet ID: $SPREADSHEET_ID"
  echo "Spreadsheet URL: $(echo "$RESPONSE" | jq -r '.spreadsheet_url' 2>/dev/null)"
  echo ""
  
  # Test 2: Get status
  echo "2. Getting status..."
  STATUS=$(curl -s -X GET "$BASE_URL/api/email-analysis/status?spreadsheet_id=$SPREADSHEET_ID" \
    -H "X-API-Key: $API_KEY")
  
  echo "Status:"
  echo "$STATUS" | jq '.' 2>/dev/null || echo "$STATUS"
else
  echo "Error: Could not extract spreadsheet ID from response"
fi

echo ""
echo "Done!"
