#!/bin/bash
# Send 10 quote preview emails with 5 minute intervals

for i in {1..10}; do
  echo "=== Request $i of 10 ==="
  echo "Time: $(date)"
  curl -X 'POST' 'http://localhost:8080/api/email/quote/preview' \
    -H 'accept: application/json' \
    -H 'Content-Type: application/json' \
    -H 'X-Api-Key: test-api-key-12345' \
    -d '{"saveAsDraft": false, "to": "bizops-dev-alexey-at-shevelyov-dot-com@shevelyov.com"}'
  echo ""
  echo ""
  
  if [ $i -lt 10 ]; then
    echo "Waiting 5 minutes (300 seconds) before next request..."
    sleep 300
  fi
done

echo "All 10 requests completed!"





