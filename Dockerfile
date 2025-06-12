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

# Build for Windows AMD64
RUN echo "Building for Windows AMD64..." && \
    x86_64-w64-mingw32-windres -i resource.rc -o resource.syso -O coff && \
    GOOS=windows GOARCH=amd64 go build -v -x -o /dist/check.exe && \
    rm -f resource.syso

# Build for Windows AMD64
RUN echo "Building for Windows AMD64..." && \
    x86_64-w64-mingw32-windres -i resource.rc -o resource.syso -O coff && \
    GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -v -x -o /dist/check.exe && \
    rm -f resource.syso

# Build for Linux AMD64
RUN echo "Building for Linux AMD64..." && \
    GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -x -o /dist/check-linux-amd64

# Build for Linux ARM64
RUN echo "Building for Linux ARM64..." && \
    GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -v -x -o /dist/check-linux-arm64

# Build for macOS Intel
RUN echo "Building for macOS Intel..." && \
    GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -v -x -o /dist/check-macos-intel

# Build for macOS ARM64
RUN echo "Building for macOS ARM64..." && \
    GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -v -x -o /dist/check-macos-arm64

# Use a minimal image to copy the binaries
FROM alpine:latest

WORKDIR /dist

# Copy binaries from builder
COPY --from=builder /dist/ .

# Set the entrypoint to list the available binaries
ENTRYPOINT ["ls", "-la"] 