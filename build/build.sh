#!/bin/sh

# build_with_retry(os, arch, output)
build_with_retry() {
  os=$1
  arch=$2
  output=$3
  # Compile Windows resource file only for Windows builds
  if [ "$os" = "windows" ]; then
    x86_64-w64-mingw32-windres -i resource.rc -o resource.syso -O coff
  fi
  for i in $(seq 1 2); do
    echo "Attempt $i: Building for $os/$arch..."
    if GOOS=$os GOARCH=$arch go build -v -x -o $output; then
      # Remove resource.syso after Windows build
      if [ "$os" = "windows" ]; then
        rm -f resource.syso
      fi
      return 0
    fi
    if [ $i -lt 2 ]; then
      echo "Attempt $i failed. Retrying in 3 seconds..."
      sleep 3
    fi
  done
  return 1
}