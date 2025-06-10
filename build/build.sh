#!/bin/sh

# build_with_retry(os, arch, output)
build_with_retry() {
    os="$1"
    arch="$2"
    output="$3"
    
    # Compile Windows resource file only for Windows builds
    if [ "$os" = "windows" ]; then
        x86_64-w64-mingw32-windres -i resource.rc -o resource.syso -O coff
    fi
    
    # First attempt
    echo "Attempt 1: Building for $os/$arch..."
    GOOS="$os" GOARCH="$arch" go build -v -x -o "$output"
    if [ $? -eq 0 ]
    then
        if [ "$os" = "windows" ]; then
            rm -f resource.syso
        fi
        exit 0
    fi
    
    # Second attempt
    echo "Attempt 1 failed. Retrying in 3 seconds..."
    sleep 3
    echo "Attempt 2: Building for $os/$arch..."
    GOOS="$os" GOARCH="$arch" go build -v -x -o "$output"
    if [ $? -eq 0 ]
    then
        if [ "$os" = "windows" ]; then
            rm -f resource.syso
        fi
        exit 0
    fi
    
    exit 1
}