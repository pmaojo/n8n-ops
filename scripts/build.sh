#!/bin/bash
# Build script for n8n CLI
# Supports cross-platform compilation

set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Build information
VERSION=$(cat VERSION 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

echo -e "${BLUE}Building n8n CLI v${VERSION}...${NC}"

# Build flags
LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Default build (current platform)
if [ "$1" == "" ] || [ "$1" == "local" ]; then
    echo -e "${YELLOW}Building for current platform...${NC}"
    go build -ldflags "${LDFLAGS}" -o n8n-ops main.go
    chmod +x n8n-ops
    echo -e "${GREEN}✅ Built: n8n-ops${NC}"
fi

# Cross-platform builds
if [ "$1" == "all" ] || [ "$1" == "cross" ]; then
    echo -e "${YELLOW}Building for all platforms...${NC}"
    
    mkdir -p dist
    
    # Linux AMD64
    echo -e "${BLUE}Building for Linux AMD64...${NC}"
    GOOS=linux GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/n8n-ops-linux-amd64 main.go
    
    # Linux ARM64
    echo -e "${BLUE}Building for Linux ARM64...${NC}"
    GOOS=linux GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o dist/n8n-ops-linux-arm64 main.go
    
    # macOS AMD64
    echo -e "${BLUE}Building for macOS AMD64...${NC}"
    GOOS=darwin GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/n8n-ops-darwin-amd64 main.go
    
    # macOS ARM64 (Apple Silicon)
    echo -e "${BLUE}Building for macOS ARM64...${NC}"
    GOOS=darwin GOARCH=arm64 go build -ldflags "${LDFLAGS}" -o dist/n8n-ops-darwin-arm64 main.go
    
    # Windows AMD64
    echo -e "${BLUE}Building for Windows AMD64...${NC}"
    GOOS=windows GOARCH=amd64 go build -ldflags "${LDFLAGS}" -o dist/n8n-ops-windows-amd64.exe main.go
    
    # Create checksums
    echo -e "${YELLOW}Creating checksums...${NC}"
    cd dist
    sha256sum * > checksums.txt
    cd ..
    
    echo -e "${GREEN}✅ Cross-platform builds completed in dist/${NC}"
    ls -la dist/
fi

# Docker build
if [ "$1" == "docker" ]; then
    echo -e "${YELLOW}Building Docker image...${NC}"
    docker build -f docker/Dockerfile -t n8n-ops:${VERSION} -t n8n-ops:latest .
    echo -e "${GREEN}✅ Docker image built: n8n-ops:${VERSION}${NC}"
fi

echo -e "${GREEN}🎉 Build completed successfully!${NC}"