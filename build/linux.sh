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

echo "Cleaning up any existing containers..."
docker rm -f check-builder-container 2>/dev/null || true

echo "Building Docker image..."
docker build -t check-builder .

echo "Creating container to extract binaries..."
docker create --name check-builder-container check-builder

echo "Copying binaries to dist directory..."
docker cp check-builder-container:/dist/. dist/

echo "Cleaning up..."
docker rm check-builder-container

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