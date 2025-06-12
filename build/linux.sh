#!/bin/bash
set -e

echo "Checking Docker installation..."
if ! command -v docker &> /dev/null; then
    echo "Docker is not installed or not in PATH"
    echo "Please install Docker from https://www.docker.com/get-started"
    exit 1
fi

echo "Creating dist directory if it doesn't exist..."
mkdir -p dist

echo "Creating secrets directory if it doesn't exist..."
mkdir -p secrets

echo "Cleaning up any existing containers..."
docker rm -f check-builder-container 2>/dev/null || true

# Prepare secrets
if [[ -n "$GOOGLE_OAUTH_CLIENT_ID" && -n "$GOOGLE_OAUTH_CLIENT_SECRET" ]]; then
    echo "$GOOGLE_OAUTH_CLIENT_ID" > secrets/client_id.txt
    echo "$GOOGLE_OAUTH_CLIENT_SECRET" > secrets/client_secret.txt
elif [[ ! -f secrets/client_id.txt || ! -f secrets/client_secret.txt ]]; then
    echo "Error: Provide GOOGLE_OAUTH_CLIENT_ID and GOOGLE_OAUTH_CLIENT_SECRET env vars, or secrets/client_id.txt and secrets/client_secret.txt files."
    exit 1
fi

# Prepare CHECK_API_KEY_DEMO secret
if [[ -n "$CHECK_API_KEY_DEMO" ]]; then
    echo "$CHECK_API_KEY_DEMO" > secrets/api_key_demo.txt
elif [[ ! -f secrets/api_key_demo.txt ]]; then
    echo "Warning: CHECK_API_KEY_DEMO env var or secrets/api_key_demo.txt not provided. Demo API key will not be available."
fi

echo "Building Docker image with BuildKit and secrets..."
DOCKER_BUILDKIT=1 docker build \
    --secret id=google_oauth_client_id,src=secrets/client_id.txt \
    --secret id=google_oauth_client_secret,src=secrets/client_secret.txt \
    --secret id=check_api_key_demo,src=secrets/api_key_demo.txt \
    -t check-builder .

echo "Creating container to extract binaries..."
docker create --name check-builder-container check-builder

echo "Copying binaries to dist directory..."
docker cp check-builder-container:/dist/. dist/

echo "Cleaning up..."
docker rm check-builder-container
rm -f secrets/client_id.txt secrets/client_secret.txt secrets/api_key_demo.txt

echo
echo "All builds completed successfully!"
echo "Executables created in $(pwd)/dist:"
echo
echo "Windows:      check.exe"
echo "Linux AMD64:  check-linux-amd64"
echo "Linux ARM64:  check-linux-arm64"
echo "macOS Intel:  check-macos-intel"
echo "macOS ARM64:  check-macos-arm64"
echo
echo "To run the executables:"
echo "Windows:      ./dist/check.exe"
echo "Linux AMD64:  ./dist/check-linux-amd64"
echo "Linux ARM64:  ./dist/check-linux-arm64"
echo "macOS Intel:  ./dist/check-macos-intel"
echo "macOS ARM64:  ./dist/check-macos-arm64"
echo
echo "To run with JSON output, add --json flag:"
echo "Windows:      ./dist/check.exe --json"
echo "Linux AMD64:  ./dist/check-linux-amd64 --json"
echo "Linux ARM64:  ./dist/check-linux-arm64 --json"
echo "macOS Intel:  ./dist/check-macos-intel --json"
echo "macOS ARM64:  ./dist/check-macos-arm64 --json" 