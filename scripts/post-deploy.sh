#!/bin/bash
# Post-deployment verification script
# This script runs after deploying workflows to verify deployment success

set -e

echo "🔍 Running post-deployment verification..."

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

if [ -z "$1" ]; then
    echo -e "${RED}❌ Environment not specified. Usage: $0 <environment>${NC}"
    exit 1
fi

ENVIRONMENT=$1

echo -e "${BLUE}🚀 Post-deployment verification for: ${ENVIRONMENT}${NC}"

# Verify deployment by syncing back from n8n
echo -e "${YELLOW}🔄 Syncing workflows to verify deployment...${NC}"
./n8n-ops sync --env "$ENVIRONMENT" --verbose

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Workflows successfully synced back from n8n${NC}"
else
    echo -e "${RED}❌ Failed to sync workflows from n8n${NC}"
    exit 1
fi

# Count deployed workflows
if [ -d "./workflows/${ENVIRONMENT}" ]; then
    WORKFLOW_COUNT=$(find "./workflows/${ENVIRONMENT}" -name "*.json" | wc -l)
    echo -e "${GREEN}📊 Found ${WORKFLOW_COUNT} workflows in ${ENVIRONMENT}${NC}"
fi

# Log deployment timestamp
TIMESTAMP=$(date '+%Y-%m-%d %H:%M:%S')
echo "${TIMESTAMP} - Deployment to ${ENVIRONMENT} completed successfully" >> .n8n-ops/deployment.log

echo -e "${GREEN}🎉 Post-deployment verification completed successfully${NC}"

# Optional: Send notification (uncomment if using Slack/Teams)
# if [ ! -z "$SLACK_WEBHOOK_URL" ]; then
#     curl -X POST -H 'Content-type: application/json' \
#          --data '{"text":"✅ n8n workflows deployed successfully to '${ENVIRONMENT}'"}' \
#          $SLACK_WEBHOOK_URL
# fi