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

# Create vendor directory
RUN go mod vendor
ENV GOFLAGS="-mod=vendor"

# Build for different platforms
RUN for os in windows linux darwin; do \
    for arch in amd64 arm64; do \
        if [ "$os" = "windows" ] && [ "$arch" = "arm64" ]; then continue; fi; \
        if [ "$os" = "darwin" ] && [ "$arch" = "arm64" ]; then \
            output="/dist/check-macos-arm64"; \
        elif [ "$os" = "darwin" ] && [ "$arch" = "amd64" ]; then \
            output="/dist/check-macos-intel"; \
        elif [ "$os" = "windows" ]; then \
            output="/dist/check.exe"; \
        else \
            output="/dist/check-$os-$arch"; \
        fi; \
        echo "Building for $os/$arch..."; \
        if [ "$os" = "windows" ]; then \
            x86_64-w64-mingw32-windres -i resource.rc -o resource.syso -O coff; \
        fi; \
        if GOOS=$os GOARCH=$arch go build -v -x -o $output; then \
            if [ "$os" = "windows" ]; then rm -f resource.syso; fi; \
            echo "Successfully built $output"; \
        else \
            echo "First attempt failed for $os/$arch, retrying..."; \
            sleep 3; \
            if GOOS=$os GOARCH=$arch go build -v -x -o $output; then \
                if [ "$os" = "windows" ]; then rm -f resource.syso; fi; \
                echo "Successfully built $output on second attempt"; \
            else \
                echo "Failed to build $output after two attempts"; \
                exit 1; \
            fi; \
        fi; \
    done; \
done

# Use a minimal image to copy the binaries
FROM alpine:latest

WORKDIR /dist

# Copy binaries from builder
COPY --from=builder /dist/ .

# Set the entrypoint to list the available binaries
ENTRYPOINT ["ls", "-la"] 