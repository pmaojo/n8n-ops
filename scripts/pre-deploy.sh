#!/bin/bash
# Pre-deployment validation script
# This script runs before deploying workflows to any environment

set -e

echo "🔍 Running pre-deployment validations..."

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Check if CLI binary exists
if [ ! -f "./n8n-cli" ]; then
    echo -e "${RED}❌ n8n-cli binary not found. Run 'go build -o n8n-cli main.go'${NC}"
    exit 1
fi

# Validate environment parameter
if [ -z "$1" ]; then
    echo -e "${RED}❌ Environment not specified. Usage: $0 <environment>${NC}"
    exit 1
fi

ENVIRONMENT=$1

echo -e "${YELLOW}📋 Validating workflows for environment: ${ENVIRONMENT}${NC}"

# Validate workflow JSON files
if [ -d "./workflows/${ENVIRONMENT}" ]; then
    echo -e "${YELLOW}🔧 Validating JSON syntax...${NC}"
    ./n8n-cli validate "./workflows/${ENVIRONMENT}/" --verbose
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ All workflow files are valid${NC}"
    else
        echo -e "${RED}❌ Workflow validation failed${NC}"
        exit 1
    fi
else
    echo -e "${YELLOW}⚠️  No workflows directory found for ${ENVIRONMENT}${NC}"
fi

# Check for required environment variables
echo -e "${YELLOW}🔑 Checking environment variables...${NC}"

case $ENVIRONMENT in
    development)
        if [ -z "$N8N_DEV_API_KEY" ] || [ -z "$N8N_DEV_URL" ]; then
            echo -e "${RED}❌ Missing development environment variables${NC}"
            echo -e "Required: N8N_DEV_API_KEY, N8N_DEV_URL"
            exit 1
        fi
        ;;
    staging)
        if [ -z "$N8N_STAGING_API_KEY" ] || [ -z "$N8N_STAGING_URL" ]; then
            echo -e "${RED}❌ Missing staging environment variables${NC}"
            echo -e "Required: N8N_STAGING_API_KEY, N8N_STAGING_URL"
            exit 1
        fi
        ;;
    production)
        if [ -z "$N8N_PROD_API_KEY" ] || [ -z "$N8N_PROD_URL" ]; then
            echo -e "${RED}❌ Missing production environment variables${NC}"
            echo -e "Required: N8N_PROD_API_KEY, N8N_PROD_URL"
            exit 1
        fi
        ;;
    *)
        echo -e "${RED}❌ Unknown environment: ${ENVIRONMENT}${NC}"
        echo -e "Supported environments: development, staging, production"
        exit 1
        ;;
esac

echo -e "${GREEN}✅ Environment variables are configured${NC}"

# Test connection to n8n instance
echo -e "${YELLOW}🌐 Testing connection to n8n instance...${NC}"
./n8n-cli sync --env "$ENVIRONMENT" --dry-run --verbose

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Successfully connected to n8n instance${NC}"
else
    echo -e "${RED}❌ Failed to connect to n8n instance${NC}"
    exit 1
fi

echo -e "${GREEN}🎉 Pre-deployment validation completed successfully${NC}"