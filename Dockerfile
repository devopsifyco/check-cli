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

# Copy build script
COPY build/build.sh /build.sh
RUN chmod +x /build.sh

# Create dist directory
RUN mkdir -p /dist

# Create vendor directory
RUN go mod vendor
ENV GOFLAGS="-mod=vendor"

# Use build_with_retry from build.sh
RUN . /build.sh && build_with_retry windows amd64 /dist/check.exe
#RUN . /build.sh && build_with_retry linux amd64 /dist/check-linux-amd64
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