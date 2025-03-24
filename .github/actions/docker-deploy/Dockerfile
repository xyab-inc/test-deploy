# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all Go source files
COPY *.go ./

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -o entrypoint

# Final stage
FROM alpine:3.19

# Install required tools
RUN apk add --no-cache openssh-client

# Copy the binary from builder
COPY --from=builder /app/entrypoint /entrypoint

# Make sure the binary is executable
RUN chmod +x /entrypoint

ENTRYPOINT ["/entrypoint"]
