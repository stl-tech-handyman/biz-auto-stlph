# Restart server with Gmail credentials loaded from .env

Write-Host "🔄 Restarting server with email credentials..." -ForegroundColor Cyan

# Stop any existing server on port 8080
$connections = Get-NetTCPConnection -LocalPort 8080 -ErrorAction SilentlyContinue
if ($connections) {
    $pids = $connections | Select-Object -ExpandProperty OwningProcess -Unique
    foreach ($pid in $pids) {
        $proc = Get-Process -Id $pid -ErrorAction SilentlyContinue
        if ($proc -and ($proc.ProcessName -like "*go*" -or $proc.Path -like "*biz-operating-system*")) {
            Write-Host "  Stopping process $pid ($($proc.ProcessName))..." -ForegroundColor Yellow
            Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
        }
    }
    Start-Sleep -Seconds 2
}

# Update .env with latest credentials if needed
$credFiles = Get-ChildItem "$env:TEMP\gmail-credentials-*.json" -ErrorAction SilentlyContinue | Sort-Object LastWriteTime -Descending
if ($credFiles) {
    $credFile = $credFiles[0].FullName
    $envPath = Join-Path $PSScriptRoot ".env"
    $envContent = Get-Content $envPath -Raw -ErrorAction SilentlyContinue
    
    if ($envContent -and $envContent -match "GMAIL_CREDENTIALS_JSON=") {
        $envContent = $envContent -replace "GMAIL_CREDENTIALS_JSON=.*", "GMAIL_CREDENTIALS_JSON=$credFile"
        if ($envContent -notmatch "GMAIL_FROM=") {
            $envContent += "`nGMAIL_FROM=team@stlpartyhelpers.com"
        }
        Set-Content $envPath $envContent -NoNewline
        Write-Host "✅ Updated .env with latest credentials" -ForegroundColor Green
    }
}

# Start server
Write-Host "🚀 Starting server..." -ForegroundColor Green
& "$PSScriptRoot\start-local.ps1"

