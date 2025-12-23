# Update .env file with Gmail credentials

$credFiles = Get-ChildItem "$env:TEMP\gmail-credentials-*.json" -ErrorAction SilentlyContinue | Sort-Object LastWriteTime -Descending

if ($credFiles) {
    $credFile = $credFiles[0].FullName
    Write-Host "Found credential file: $credFile" -ForegroundColor Green
    
    $envPath = Join-Path $PSScriptRoot ".env"
    $envContent = Get-Content $envPath -Raw -ErrorAction SilentlyContinue
    
    if ($envContent -and $envContent -match "GMAIL_CREDENTIALS_JSON=") {
        Write-Host "GMAIL_CREDENTIALS_JSON already in .env file" -ForegroundColor Yellow
        # Update it anyway
        $envContent = $envContent -replace "GMAIL_CREDENTIALS_JSON=.*", "GMAIL_CREDENTIALS_JSON=$credFile"
        if ($envContent -notmatch "GMAIL_FROM=") {
            $envContent += "`nGMAIL_FROM=team@stlpartyhelpers.com"
        }
        Set-Content $envPath $envContent -NoNewline
        Write-Host "Updated Gmail credentials in .env file" -ForegroundColor Green
    } else {
        Add-Content $envPath "`nGMAIL_CREDENTIALS_JSON=$credFile`nGMAIL_FROM=team@stlpartyhelpers.com"
        Write-Host "Added Gmail credentials to .env file" -ForegroundColor Green
    }
} else {
    Write-Host "No credential file found. Run: powershell -ExecutionPolicy Bypass -File scripts/get-gmail-credentials.ps1" -ForegroundColor Red
}

