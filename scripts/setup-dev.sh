#!/bin/bash
# Development setup script for n8n CLI
# This script helps you set up the development environment quickly

set -e

echo "🚀 Setting up n8n CLI development environment..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.19+ first.${NC}"
    echo "Visit: https://golang.org/doc/install"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | cut -d' ' -f3 | sed 's/go//')
echo -e "${BLUE}📦 Go version: ${GO_VERSION}${NC}"

# Create necessary directories
echo -e "${BLUE}📁 Creating project directories...${NC}"
mkdir -p workflows/{development,staging,production}
mkdir -p scripts
mkdir -p config/templates
mkdir -p .n8n-cli/logs
mkdir -p tests

# Build the CLI
echo -e "${BLUE}🔨 Building n8n CLI...${NC}"
go mod tidy
go build -o n8n-cli main.go
chmod +x n8n-cli

# Test the build
echo -e "${BLUE}🧪 Testing CLI build...${NC}"
./n8n-cli --help > /dev/null
echo -e "${GREEN}✅ CLI built successfully${NC}"

# Copy example configuration
if [ ! -f ~/.n8n-cli.yaml ]; then
    echo -e "${BLUE}⚙️ Creating example configuration...${NC}"
    cp config.example.yaml ~/.n8n-cli.yaml
    echo -e "${YELLOW}📝 Configuration created at ~/.n8n-cli.yaml${NC}"
    echo -e "${YELLOW}   Please edit this file with your n8n instance URLs and API keys${NC}"
else
    echo -e "${YELLOW}📝 Configuration already exists at ~/.n8n-cli.yaml${NC}"
fi

# Create environment file template
if [ ! -f .env ]; then
    echo -e "${BLUE}🔧 Creating .env template...${NC}"
    cat > .env << 'EOF'
# n8n API Keys (get these from your n8n instance admins)
N8N_DEV_API_KEY=your_development_api_key_here
N8N_STAGING_API_KEY=your_staging_api_key_here
N8N_PROD_API_KEY=your_production_api_key_here

# n8n Instance URLs
N8N_DEV_URL=https://n8n-dev.yourcompany.com
N8N_STAGING_URL=https://n8n-staging.yourcompany.com
N8N_PROD_URL=https://n8n-prod.yourcompany.com

# GitLab Configuration
GITLAB_TOKEN=glpat-your_gitlab_token_here
GITLAB_PROJECT_ID=12345

# CLI Language (en, es)
N8N_CLI_LANG=en
EOF
    echo -e "${YELLOW}📝 .env template created - please fill in your actual values${NC}"
    echo -e "${RED}⚠️  Don't commit .env to git! It's already in .gitignore${NC}"
fi

# Create .gitignore if it doesn't exist
if [ ! -f .gitignore ]; then
    echo -e "${BLUE}📝 Creating .gitignore...${NC}"
    cat > .gitignore << 'EOF'
# Binaries
n8n-cli
*.exe
*.dll
*.so
*.dylib

# Environment files
.env
.env.local

# CLI data
.n8n-cli/
!.n8n-cli/.gitkeep

# Go
*.test
*.out
vendor/

# IDE
.vscode/
.idea/
*.swp
*.swo
*~

# OS
.DS_Store
Thumbs.db

# Logs
*.log
logs/
EOF
fi

# Create example workflow
if [ ! -f workflows/development/example-workflow.json ]; then
    echo -e "${BLUE}📋 Creating example workflow...${NC}"
    cat > workflows/development/example-workflow.json << 'EOF'
{
  "name": "Example Workflow",
  "nodes": [
    {
      "parameters": {
        "httpMethod": "POST",
        "path": "webhook",
        "responseMode": "onReceived",
        "options": {}
      },
      "id": "webhook-1",
      "name": "Webhook",
      "type": "n8n-nodes-base.webhook",
      "typeVersion": 1,
      "position": [250, 300]
    },
    {
      "parameters": {
        "operation": "get",
        "propertyName": "timestamp",
        "value": "={{ new Date().toISOString() }}"
      },
      "id": "set-1",
      "name": "Set Timestamp",
      "type": "n8n-nodes-base.set",
      "typeVersion": 1,
      "position": [450, 300]
    }
  ],
  "connections": {
    "Webhook": {
      "main": [
        [
          {
            "node": "Set Timestamp",
            "type": "main",
            "index": 0
          }
        ]
      ]
    }
  },
  "active": false,
  "settings": {},
  "versionId": "1"
}
EOF
fi

# Display setup summary
echo ""
echo -e "${GREEN}✅ Development environment setup complete!${NC}"
echo ""
echo -e "${BLUE}Next steps:${NC}"
echo -e "1. ${YELLOW}Edit ~/.n8n-cli.yaml with your n8n instance details${NC}"
echo -e "2. ${YELLOW}Fill in .env file with your API keys${NC}"
echo -e "3. ${YELLOW}Source your environment: source .env${NC}"
echo -e "4. ${YELLOW}Test the CLI: ./n8n-cli welcome${NC}"
echo -e "5. ${YELLOW}Validate example workflow: ./n8n-cli validate workflows/development/${NC}"
echo ""
echo -e "${BLUE}Useful commands:${NC}"
echo -e "• ${GREEN}./n8n-cli --help${NC} - Show all available commands"
echo -e "• ${GREEN}./n8n-cli welcome${NC} - Display welcome screen"
echo -e "• ${GREEN}./n8n-cli --lang es welcome${NC} - Spanish welcome screen"
echo -e "• ${GREEN}./n8n-cli validate ./workflows/${NC} - Validate all workflows"
echo -e "• ${GREEN}./n8n-cli sync --env development${NC} - Sync from n8n instance"
echo ""
echo -e "${GREEN}🎉 Happy workflow automation!${NC}"