@echo off
REM Re-authenticate with Google Cloud (Windows Batch)

echo === Google Cloud Re-authentication ===
echo.

REM Check gcloud is installed
where gcloud >nul 2>&1
if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] gcloud CLI is not installed
    exit /b 1
)

REM Step 1: User authentication
echo [INFO] Step 1: Authenticating user account...
echo [INFO] This will open a browser for you to sign in
gcloud auth login

if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] User authentication failed
    exit /b 1
)
echo [SUCCESS] User authentication complete
echo.

REM Step 2: Application default credentials
echo [INFO] Step 2: Setting up application-default credentials...
echo [INFO] This is needed for Docker and local development
gcloud auth application-default login

if %ERRORLEVEL% NEQ 0 (
    echo [ERROR] Application-default credentials setup failed
    exit /b 1
)
echo [SUCCESS] Application-default credentials configured
echo.

REM Step 3: Configure Docker
echo [INFO] Step 3: Configuring Docker authentication...
gcloud auth configure-docker

if %ERRORLEVEL% NEQ 0 (
    echo [WARNING] Docker authentication may need manual setup
) else (
    echo [SUCCESS] Docker authentication configured
)
echo.

REM Show current status
echo [INFO] Current authentication status:
echo.
gcloud auth list
echo.

REM Show current project
echo [INFO] Current project:
gcloud config get-value project
echo.

echo [SUCCESS] Re-authentication complete!

