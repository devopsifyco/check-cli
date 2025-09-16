FROM golang:1.23-alpine AS builder

# Install required tools and certificates
RUN apk add --no-cache mingw-w64-gcc binutils ca-certificates git && update-ca-certificates

# Set Go environment variables
ENV GOPROXY=direct
ENV GOSUMDB=sum.golang.org
ENV GO111MODULE=on
ENV GOFLAGS="-mod=mod"

WORKDIR /app

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download dependencies with enhanced retry logic
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    for i in $(seq 1 2); do \
        echo "Attempt $i: Downloading dependencies..." && \
        go mod download -x && break || \
        if [ $i -lt 2 ]; then \
            echo "Attempt $i failed. Retrying in 3 seconds..." && \
            sleep 3; \
        fi; \
    done

# Copy source code and resource files
COPY . .

# Create dist directory
RUN mkdir -p /dist

# Copy vendor directory from local
COPY vendor/ vendor/
ENV GOFLAGS="-mod=vendor"

# Build for all targets in a loop
RUN --mount=type=secret,id=google_oauth_client_id \
    --mount=type=secret,id=google_oauth_client_secret \
    --mount=type=secret,id=check_api_key_demo \
    set -e; \
    GOOGLE_OAUTH_CLIENT_ID=$(cat /run/secrets/google_oauth_client_id | tr -d '\r\n'); \
    GOOGLE_OAUTH_CLIENT_SECRET=$(cat /run/secrets/google_oauth_client_secret | tr -d '\r\n'); \
    CHECK_API_KEY_DEMO=$(cat /run/secrets/check_api_key_demo | tr -d '\r\n'); \
    echo "GOOGLE_OAUTH_CLIENT_ID: '$GOOGLE_OAUTH_CLIENT_ID'"; \
    echo "GOOGLE_OAUTH_CLIENT_SECRET: '$GOOGLE_OAUTH_CLIENT_SECRET'"; \
    echo "CHECK_API_KEY_DEMO: '$CHECK_API_KEY_DEMO'"; \
    if [ -z "$GOOGLE_OAUTH_CLIENT_ID" ] || [ -z "$GOOGLE_OAUTH_CLIENT_SECRET" ] || [ -z "$CHECK_API_KEY_DEMO" ]; then \
      echo "Secrets must not be empty"; exit 1; \
    fi; \
    ldflags="-s -w -X auth.googleOAuthClientID=$GOOGLE_OAUTH_CLIENT_ID -X auth.googleOAuthClientSecret=$GOOGLE_OAUTH_CLIENT_SECRET -X auth.CheckApiKeyDemo=$CHECK_API_KEY_DEMO"; \
    targets="windows amd64 check.exe linux amd64 check-linux-amd64 linux arm64 check-linux-arm64 darwin amd64 check-macos-intel darwin arm64 check-macos-arm64"; \
    set -- $targets; \
    while [ $# -gt 0 ]; do \
      GOOS=$1; GOARCH=$2; OUT=$3; \
      shift 3; \
      if [ "$GOOS" = "windows" ]; then \
        echo "Building for Windows AMD64..."; \
        x86_64-w64-mingw32-windres -i resource.rc -o resource.syso -O coff; \
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$ldflags" -v -x -o /dist/$OUT; \
        rm -f resource.syso; \
      else \
        echo "Building for $GOOS $GOARCH..."; \
        GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$ldflags" -v -x -o /dist/$OUT; \
      fi; \
    done

# Use a minimal image to copy the binaries
FROM alpine:latest

WORKDIR /dist

# Copy binaries from builder
COPY --from=builder /dist/ .

# Set the entrypoint to list the available binaries
ENTRYPOINT ["ls", "-la"] 