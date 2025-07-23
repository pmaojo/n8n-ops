# Multi-stage Dockerfile for n8n CLI
# This creates a lightweight container with the n8n CLI tool

# Build stage
FROM golang:1.19-alpine AS builder

# Set working directory
WORKDIR /app

# Install git for Go modules
RUN apk add --no-cache git

# Copy Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the CLI binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o n8n-ops main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates git

# Create non-root user
RUN addgroup -g 1000 n8ncli && \
    adduser -D -s /bin/sh -u 1000 -G n8ncli n8ncli

# Set working directory
WORKDIR /home/n8ncli

# Copy binary from builder stage
COPY --from=builder /app/n8n-ops /usr/local/bin/n8n-ops

# Copy configuration examples
COPY --from=builder /app/config.example.yaml /home/n8ncli/
COPY --from=builder /app/QUICK_START.md /home/n8ncli/

# Create necessary directories
RUN mkdir -p /home/n8ncli/.n8n-ops/logs && \
    mkdir -p /home/n8ncli/workflows/{development,staging,production} && \
    chown -R n8ncli:n8ncli /home/n8ncli

# Switch to non-root user
USER n8ncli

# Set entrypoint
ENTRYPOINT ["n8n-ops"]

# Default command
CMD ["welcome"]

# Labels
LABEL maintainer="n8n CLI Team"
LABEL description="n8n CLI tool for managing workflows across environments"
LABEL version="1.0.0"