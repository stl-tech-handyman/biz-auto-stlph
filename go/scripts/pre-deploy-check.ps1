# Pre-deployment checks for Go API (PowerShell)

$ErrorActionPreference = "Stop"

Write-Host "[INFO] Running pre-deployment checks..." -ForegroundColor Blue

# Check gcloud is installed
try {
    $null = gcloud --version 2>&1
    Write-Host "[SUCCESS] gcloud CLI found" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] gcloud CLI is not installed" -ForegroundColor Red
    exit 1
}

# Check Docker is installed
try {
    $null = docker --version 2>&1
    Write-Host "[SUCCESS] Docker found" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Docker is not installed" -ForegroundColor Red
    exit 1
}

# Check Docker is running
try {
    $null = docker ps 2>&1
    Write-Host "[SUCCESS] Docker is running" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Docker is not running" -ForegroundColor Red
    exit 1
}

# Check gcloud authentication
$accounts = gcloud auth list --filter=status:ACTIVE --format="value(account)" 2>&1
if ([string]::IsNullOrEmpty($accounts)) {
    Write-Host "[ERROR] No active gcloud account. Run: gcloud auth login" -ForegroundColor Red
    exit 1
}
Write-Host "[SUCCESS] gcloud authenticated as: $accounts" -ForegroundColor Green

# Check application-default credentials
try {
    $null = gcloud auth application-default print-access-token 2>&1
    Write-Host "[SUCCESS] Application-default credentials configured" -ForegroundColor Green
} catch {
    Write-Host "[WARNING] Application-default credentials not set" -ForegroundColor Yellow
    Write-Host "[INFO] Run: gcloud auth application-default login" -ForegroundColor Blue
    Write-Host "[INFO] This is needed for Docker to push to Artifact Registry" -ForegroundColor Blue
    $confirm = Read-Host "Do you want to set it up now? (yes/no)"
    if ($confirm -eq "yes") {
        gcloud auth application-default login
    } else {
        Write-Host "[ERROR] Cannot proceed without application-default credentials" -ForegroundColor Red
        exit 1
    }
}

# Check projects exist
Write-Host "[INFO] Checking GCP projects..." -ForegroundColor Blue

$PROJECT_DEV = "bizops360-dev"
$PROJECT_PROD = "bizops360-prod"

try {
    $null = gcloud projects describe $PROJECT_DEV 2>&1
    Write-Host "[SUCCESS] Project $PROJECT_DEV accessible" -ForegroundColor Green
} catch {
    Write-Host "[ERROR] Project $PROJECT_DEV not found or not accessible" -ForegroundColor Red
    exit 1
}

try {
    $null = gcloud projects describe $PROJECT_PROD 2>&1
    Write-Host "[SUCCESS] Project $PROJECT_PROD accessible" -ForegroundColor Green
} catch {
    Write-Host "[WARNING] Project $PROJECT_PROD not found or not accessible" -ForegroundColor Yellow
    Write-Host "[INFO] You can still deploy to dev, but prod deployment will fail" -ForegroundColor Blue
}

# Check Docker authentication
Write-Host "[INFO] Checking Docker authentication..." -ForegroundColor Blue
try {
    $null = gcloud auth configure-docker 2>&1
    Write-Host "[SUCCESS] Docker authentication configured" -ForegroundColor Green
} catch {
    Write-Host "[WARNING] Docker authentication may need setup" -ForegroundColor Yellow
    Write-Host "[INFO] Run: gcloud auth configure-docker" -ForegroundColor Blue
}

Write-Host "[SUCCESS] Pre-deployment checks complete!" -ForegroundColor Green
Write-Host "[INFO] You can now run: .\scripts\deploy-dev.ps1" -ForegroundColor Blue





