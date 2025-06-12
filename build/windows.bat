@echo off
setlocal enabledelayedexpansion

echo Checking Docker installation...
docker --version >nul 2>&1
if errorlevel 1 (
    echo Docker is not installed or not in PATH
    echo Please install Docker from https://www.docker.com/get-started
    exit /b 1
)

echo Creating dist directory if it doesn't exist...
if not exist dist mkdir dist

echo Creating secrets directory if it doesn't exist...
if not exist secrets mkdir secrets

echo Cleaning up any existing containers...
docker rm -f check-builder-container 2>nul

REM Prepare secrets
if defined GOOGLE_OAUTH_CLIENT_ID if defined GOOGLE_OAUTH_CLIENT_SECRET if defined CHECK_API_KEY_DEMO (
    echo !GOOGLE_OAUTH_CLIENT_ID! > secrets/client_id.txt
    echo !GOOGLE_OAUTH_CLIENT_SECRET! > secrets/client_secret.txt
    echo !CHECK_API_KEY_DEMO! > secrets/api_key_demo.txt
) else (
    if not exist secrets/client_id.txt (
        echo Error: Provide GOOGLE_OAUTH_CLIENT_ID and GOOGLE_OAUTH_CLIENT_SECRET env vars, or secrets/client_id.txt and secrets/client_secret.txt files.
        exit /b 1
    )
    if not exist secrets/client_secret.txt (
        echo Error: Provide GOOGLE_OAUTH_CLIENT_ID and GOOGLE_OAUTH_CLIENT_SECRET env vars, or secrets/client_id.txt and secrets/client_secret.txt files.
        exit /b 1
    )
    if not exist secrets/api_key_demo.txt (
        echo Error: Provide CHECK_API_KEY_DEMO env var, or secrets/api_key_demo.txt file.
        exit /b 1
    )
)

echo Building Docker image with BuildKit and secrets...
set DOCKER_BUILDKIT=1
docker build --secret id=google_oauth_client_id,src=secrets/client_id.txt --secret id=google_oauth_client_secret,src=secrets/client_secret.txt --secret id=check_api_key_demo,src=secrets/api_key_demo.txt -t check-builder .
if errorlevel 1 (
    echo Failed to build Docker image
    exit /b 1
)

echo Creating container to extract binaries...
docker create --name check-builder-container check-builder
if errorlevel 1 (
    echo Failed to create container
    exit /b 1
)

echo Copying binaries to dist directory...
docker cp check-builder-container:/dist/. dist/
if errorlevel 1 (
    echo Failed to copy binaries
    docker rm check-builder-container
    exit /b 1
)

echo Cleaning up...
docker rm check-builder-container

del secrets\client_id.txt secrets\client_secret.txt secrets\api_key_demo.txt >nul 2>&1

echo.
echo All builds completed successfully!
echo Executables created in %CD%\dist:
echo.
echo Windows:      check.exe
echo Linux AMD64:  check-linux-amd64
echo Linux ARM64:  check-linux-arm64
echo macOS Intel:  check-macos-intel
echo macOS ARM64:  check-macos-arm64
echo.
echo To run the executables:
echo Windows:      .\dist\check.exe
echo Linux AMD64:  ./dist/check-linux-amd64
echo Linux ARM64:  ./dist/check-linux-arm64
echo macOS Intel:  ./dist/check-macos-intel
echo macOS ARM64:  ./dist/check-macos-arm64
echo.
echo To run with JSON output, add --json flag:
echo Windows:      .\dist\check.exe --json
echo Linux AMD64:  ./dist/check-linux-amd64 --json
echo Linux ARM64:  ./dist/check-linux-arm64 --json
echo macOS Intel:  ./dist/check-macos-intel --json
echo macOS ARM64:  ./dist/check-macos-arm64 --json