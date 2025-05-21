FROM golang:1.21-alpine AS builder

# Install required tools and certificates
RUN apk add --no-cache mingw-w64-gcc binutils ca-certificates && update-ca-certificates

# Set Go environment variables
ENV GOPROXY=https://proxy.golang.org,direct
ENV GOSUMDB=sum.golang.org
ENV GO111MODULE=on
ENV GOFLAGS="-mod=mod"

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies with enhanced retry logic
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    for i in $(seq 1 5); do \
        echo "Attempt $i: Downloading dependencies..." && \
        go mod download -x && break || \
        if [ $i -lt 5 ]; then \
            echo "Attempt $i failed. Retrying in 3 seconds..." && \
            sleep 3; \
        fi; \
    done

# Copy source code and resource files
COPY . .

# Update dependencies with enhanced retry logic
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    for i in $(seq 1 5); do \
        echo "Attempt $i: Running go mod tidy..." && \
        go mod tidy -v && break || \
        if [ $i -lt 5 ]; then \
            echo "Attempt $i failed. Retrying in 3 seconds..." && \
            sleep 3; \
        fi; \
    done

# Create dist directory
RUN mkdir -p /dist

# Function to build with retry logic
RUN echo 'build_with_retry() {' > /build.sh && \
    echo '  os=$1; arch=$2; output=$3' >> /build.sh && \
    echo '  # Compile Windows resource file only for Windows builds' >> /build.sh && \
    echo '  if [ "$os" = "windows" ]; then' >> /build.sh && \
    echo '    x86_64-w64-mingw32-windres -i resource.rc -o resource.syso -O coff' >> /build.sh && \
    echo '  fi' >> /build.sh && \
    echo '  for i in $(seq 1 5); do' >> /build.sh && \
    echo '    echo "Attempt $i: Building for $os/$arch..."' >> /build.sh && \
    echo '    if GOOS=$os GOARCH=$arch go build -v -x -o $output; then' >> /build.sh && \
    echo '      # Remove resource.syso after Windows build' >> /build.sh && \
    echo '      if [ "$os" = "windows" ]; then' >> /build.sh && \
    echo '        rm -f resource.syso' >> /build.sh && \
    echo '      fi' >> /build.sh && \
    echo '      return 0' >> /build.sh && \
    echo '    fi' >> /build.sh && \
    echo '    if [ $i -lt 5 ]; then' >> /build.sh && \
    echo '      echo "Attempt $i failed. Retrying in 3 seconds..."' >> /build.sh && \
    echo '      sleep 3' >> /build.sh && \
    echo '    fi' >> /build.sh && \
    echo '  done' >> /build.sh && \
    echo '  return 1' >> /build.sh && \
    echo '}' >> /build.sh && \
    chmod +x /build.sh && \
    . /build.sh

# Build for multiple platforms with retry logic
RUN . /build.sh && build_with_retry windows amd64 /dist/check.exe
RUN . /build.sh && build_with_retry linux amd64 /dist/check-linux-amd64
#RUN . /build.sh && build_with_retry linux arm64 /dist/check-linux-arm64
#RUN . /build.sh && build_with_retry darwin amd64 /dist/check-macos-intel
#RUN . /build.sh && build_with_retry darwin arm64 /dist/check-macos-arm64

# Use a minimal image to copy the binaries
FROM alpine:latest

WORKDIR /dist

# Copy binaries from builder
COPY --from=builder /dist/ .

# Set the entrypoint to list the available binaries
ENTRYPOINT ["ls", "-la"] 