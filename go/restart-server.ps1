# Restart Go API Server Script
# This script stops any running server and finds an available port, automatically incrementing on binding errors

$BASE_PORT = 8080
$MAX_PORT = 8099

Write-Host "🔄 Restarting Go API Server..." -ForegroundColor Green
Write-Host ""

# Function to stop process on a port
function Stop-ServerOnPort {
    param($Port)
    $connections = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
    if ($connections) {
        $processes = $connections | Select-Object -ExpandProperty OwningProcess -Unique
        foreach ($pid in $processes) {
            $proc = Get-Process -Id $pid -ErrorAction SilentlyContinue
            if ($proc -and ($proc.ProcessName -like "*go*" -or $proc.Path -like "*biz-operating-system*")) {
                Write-Host "  Stopping process $pid ($($proc.ProcessName)) on port $Port..." -ForegroundColor Yellow
                Stop-Process -Id $pid -Force -ErrorAction SilentlyContinue
                Start-Sleep -Milliseconds 500
            }
        }
    }
}

# Function to check if a port is available
function Test-PortAvailable {
    param($Port)
    $connection = Get-NetTCPConnection -LocalPort $Port -ErrorAction SilentlyContinue
    return ($connection -eq $null)
}

# Function to find the first available port
function Find-AvailablePort {
    param($StartPort, $MaxPort)
    for ($port = $StartPort; $port -le $MaxPort; $port++) {
        if (Test-PortAvailable -Port $port) {
            return $port
        }
    }
    return $null
}

# Stop servers on common ports
Write-Host "  Checking for running servers..." -ForegroundColor Cyan
for ($port = $BASE_PORT; $port -le $MAX_PORT; $port++) {
    Stop-ServerOnPort -Port $port
}

Write-Host "  Waiting for ports to be released..." -ForegroundColor Cyan
Start-Sleep -Seconds 2

# Find first available port
$PORT = Find-AvailablePort -StartPort $BASE_PORT -MaxPort $MAX_PORT
if ($null -eq $PORT) {
    Write-Host "❌ ERROR: No available ports found in range $BASE_PORT-$MAX_PORT" -ForegroundColor Red
    exit 1
}

if ($PORT -ne $BASE_PORT) {
    Write-Host "⚠️  Port $BASE_PORT is in use, using port $PORT instead" -ForegroundColor Yellow
}

# Set environment variables
$env:STRIPE_SECRET_KEY_PROD="sk_live_YOUR_PROD_KEY_HERE"
$env:SERVICE_API_KEY="test-api-key-12345"
$env:ENV="dev"
$env:PORT=$PORT
$env:LOG_LEVEL="debug"
$env:CONFIG_DIR="./config"
$env:TEMPLATES_DIR="./templates"

Write-Host ""
Write-Host "🚀 Starting Go API Server..." -ForegroundColor Green
Write-Host ""
Write-Host "Configuration:" -ForegroundColor Yellow
Write-Host "  Stripe: LIVE (Production) Key" -ForegroundColor Red
Write-Host "  API Key: $env:SERVICE_API_KEY" -ForegroundColor Cyan
Write-Host "  Port: $env:PORT" -ForegroundColor Cyan
Write-Host "  Environment: $env:ENV" -ForegroundColor Cyan
Write-Host ""
Write-Host "⚠️  WARNING: Using LIVE Stripe key - real charges will occur!" -ForegroundColor Red
Write-Host ""
Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Yellow
Write-Host ""

# Run the server with automatic port retry on binding errors
$attemptPort = $PORT
$maxAttempts = ($MAX_PORT - $BASE_PORT) + 1

for ($attempt = 0; $attempt -lt $maxAttempts; $attempt++) {
    $env:PORT = $attemptPort.ToString()
    
    Write-Host "  Attempting to start on port $attemptPort..." -ForegroundColor Cyan
    
    # Create a script block that will run the server
    $scriptBlock = {
        param($Port, $StripeKey, $ApiKey, $Env, $LogLevel, $ConfigDir, $TemplatesDir, $WorkingDir)
        Set-Location $WorkingDir
        $env:STRIPE_SECRET_KEY_PROD = $StripeKey
        $env:SERVICE_API_KEY = $ApiKey
        $env:ENV = $Env
        $env:PORT = $Port
        $env:LOG_LEVEL = $LogLevel
        $env:CONFIG_DIR = $ConfigDir
        $env:TEMPLATES_DIR = $TemplatesDir
        go run ./cmd/api
    }
    
    # Start server in a job to monitor for errors
    $job = Start-Job -ScriptBlock $scriptBlock -ArgumentList `
        $attemptPort.ToString(), `
        $env:STRIPE_SECRET_KEY_PROD, `
        $env:SERVICE_API_KEY, `
        $env:ENV, `
        $env:LOG_LEVEL, `
        $env:CONFIG_DIR, `
        $env:TEMPLATES_DIR, `
        $PWD
    
    # Wait a moment to see if server starts or errors
    Start-Sleep -Seconds 3
    
    # Check job output for binding errors
    $jobOutput = Receive-Job -Job $job
    $bindingError = $jobOutput | Select-String -Pattern "bind.*address.*already in use|Only one usage of each socket address|listen tcp.*bind.*Only one usage" -Quiet
    
    if ($bindingError) {
        Write-Host "⚠️  Port $attemptPort is in use, trying next port..." -ForegroundColor Yellow
        Stop-Job -Job $job -ErrorAction SilentlyContinue
        Remove-Job -Job $job -ErrorAction SilentlyContinue
        $attemptPort++
        continue
    }
    
    # Check if job is still running (server started successfully)
    $jobState = (Get-Job -Id $job.Id).State
    if ($jobState -eq "Running") {
        Write-Host "✅ Server started successfully on port $attemptPort" -ForegroundColor Green
        $env:PORT = $attemptPort.ToString()
        # Wait for job to complete (server runs until stopped)
        Wait-Job -Job $job | Out-Null
        Receive-Job -Job $job | Write-Host
        Remove-Job -Job $job -ErrorAction SilentlyContinue
        exit 0
    }
    
    # Job completed/failed - check for other errors
    if ($jobOutput) {
        Write-Host "❌ Server failed to start: $jobOutput" -ForegroundColor Red
        Remove-Job -Job $job -ErrorAction SilentlyContinue
        exit 1
    }
}

Write-Host "❌ ERROR: Failed to find available port after $maxAttempts attempts" -ForegroundColor Red
exit 1
