# Re-authenticate with Google Cloud (PowerShell)

$ErrorActionPreference = "Stop"

Write-Host "=== Google Cloud Re-authentication ===" -ForegroundColor Cyan
Write-Host ""

# Check gcloud is installed
if (-not (Get-Command gcloud -ErrorAction SilentlyContinue)) {
    Write-Host "[ERROR] gcloud CLI is not installed" -ForegroundColor Red
    exit 1
}

# Step 1: User authentication
Write-Host "[INFO] Step 1: Authenticating user account..." -ForegroundColor Blue
Write-Host "[INFO] This will open a browser for you to sign in" -ForegroundColor Blue
gcloud auth login

if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] User authentication complete" -ForegroundColor Green
} else {
    Write-Host "[ERROR] User authentication failed" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Step 2: Application default credentials
Write-Host "[INFO] Step 2: Setting up application-default credentials..." -ForegroundColor Blue
Write-Host "[INFO] This is needed for Docker and local development" -ForegroundColor Blue
gcloud auth application-default login

if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] Application-default credentials configured" -ForegroundColor Green
} else {
    Write-Host "[ERROR] Application-default credentials setup failed" -ForegroundColor Red
    exit 1
}

Write-Host ""

# Step 3: Configure Docker
Write-Host "[INFO] Step 3: Configuring Docker authentication..." -ForegroundColor Blue
gcloud auth configure-docker

if ($LASTEXITCODE -eq 0) {
    Write-Host "[SUCCESS] Docker authentication configured" -ForegroundColor Green
} else {
    Write-Host "[WARNING] Docker authentication may need manual setup" -ForegroundColor Yellow
}

Write-Host ""

# Show current status
Write-Host "[INFO] Current authentication status:" -ForegroundColor Blue
Write-Host ""
gcloud auth list
Write-Host ""

# Show current project
$currentProject = gcloud config get-value project 2>&1
if ($LASTEXITCODE -ne 0) {
    $currentProject = "not set"
}

Write-Host "[INFO] Current project: $currentProject" -ForegroundColor Blue

if ($currentProject -ne "bizops360-dev" -and $currentProject -ne "bizops360-prod") {
    Write-Host "[WARNING] Project is not set to bizops360-dev or bizops360-prod" -ForegroundColor Yellow
    Write-Host "[INFO] Set it with: gcloud config set project bizops360-dev" -ForegroundColor Blue
}

Write-Host ""
Write-Host "[SUCCESS] Re-authentication complete!" -ForegroundColor Green

