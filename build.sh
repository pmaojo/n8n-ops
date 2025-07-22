#!/bin/bash

# n8n-ops Build Script with Version Information
set -e

# Get version information
VERSION=$(cat VERSION 2>/dev/null || echo "1.0.0")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

echo "Building n8n-ops..."
echo "Version: $VERSION"
echo "Git Commit: $GIT_COMMIT"
echo "Build Date: $BUILD_DATE"

# Build with version information
go build -ldflags "\
  -X 'github.com/n8n-workflows/n8n-ops/cmd.Version=$VERSION' \
  -X 'github.com/n8n-workflows/n8n-ops/cmd.GitCommit=$GIT_COMMIT' \
  -X 'github.com/n8n-workflows/n8n-ops/cmd.BuildDate=$BUILD_DATE'" \
  -o n8n-ops main.go

# Make executable
chmod +x n8n-ops

echo "✅ Build complete: n8n-ops"
echo "Run './n8n-ops version' to see version information"