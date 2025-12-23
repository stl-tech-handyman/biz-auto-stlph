# Grant Cloud Run service account access to secrets

$ErrorActionPreference = "Stop"

$PROJECT_ID = "bizops360-dev"

Write-Host "[INFO] Granting Cloud Run service account access to secrets..." -ForegroundColor Blue
Write-Host "[INFO] Project: $PROJECT_ID" -ForegroundColor Blue

# Get project number and service account
$PROJECT_NUMBER = gcloud projects describe $PROJECT_ID --format="value(projectNumber)"
$SERVICE_ACCOUNT = "$PROJECT_NUMBER-compute@developer.gserviceaccount.com"

Write-Host "[INFO] Service Account: $SERVICE_ACCOUNT" -ForegroundColor Blue

# Grant access to each secret
$secrets = @("svc-api-key-dev", "stripe-secret-key-test", "stripe-secret-key-prod", "gmail-credentials-json")

foreach ($secret in $secrets) {
    Write-Host "[INFO] Granting access to secret: $secret" -ForegroundColor Yellow
    
    $check = gcloud secrets describe $secret --project=$PROJECT_ID 2>&1
    if ($LASTEXITCODE -eq 0) {
        gcloud secrets add-iam-policy-binding $secret `
            --member="serviceAccount:$SERVICE_ACCOUNT" `
            --role="roles/secretmanager.secretAccessor" `
            --project=$PROJECT_ID
        
        if ($LASTEXITCODE -eq 0) {
            Write-Host "[SUCCESS] Access granted to $secret" -ForegroundColor Green
        } else {
            Write-Host "[WARNING] Failed to grant access to $secret (may already have access)" -ForegroundColor Yellow
        }
    } else {
        Write-Host "[WARNING] Secret $secret not found, skipping" -ForegroundColor Yellow
    }
}

Write-Host "[SUCCESS] Permission update complete!" -ForegroundColor Green

