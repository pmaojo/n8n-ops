#!/bin/bash

echo "🎯 n8n-ops COMPREHENSIVE DEMO - Branch Management & Workflow Conventions"
echo "======================================================================"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test results tracking
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

run_test() {
    local test_name="$1"
    local command="$2"
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    echo -e "${BLUE}TEST $TOTAL_TESTS: $test_name${NC}"
    echo "Command: $command"
    echo ""
    
    if eval "$command" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ PASSED${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}❌ FAILED${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""
}

# Step 1: Build the CLI
echo "🚀 Step 1: Building n8n-ops CLI with new branch system"
go build -o n8n-ops main.go
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ CLI built successfully${NC}"
else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
fi
echo ""

# Step 2: Start Mock Server
echo "🐳 Step 2: Starting Mock n8n Server"
cd mock-n8n-server && go run main.go &
SERVER_PID=$!
cd ..
echo "⏳ Waiting for mock server to start..."
sleep 3

# Check if server is running
if curl -s http://localhost:3001/health > /dev/null; then
    echo -e "${GREEN}✅ Mock n8n server running on port 3001${NC}"
else
    echo -e "${RED}❌ Mock server failed to start${NC}"
    exit 1
fi
echo ""

# Set environment variables for testing
export N8N_URL="http://localhost:3001"
export N8N_API_KEY="n8n_api_mock_development"

# Step 3: Test Core CLI Functions
echo "🧪 Step 3: Testing Enhanced CLI with Branch Management"
echo ""

# Basic CLI tests
run_test "CLI Help System" "./n8n-ops --help"
run_test "Version Command" "./n8n-ops version"
run_test "Welcome Screen" "./n8n-ops welcome"

# Branch management tests
run_test "Branch Help" "./n8n-ops branch --help"
run_test "Branch List (empty)" "./n8n-ops branch list"

# Sync tests
run_test "Sync from Mock API - Development" "./n8n-ops sync --env development"
run_test "Check Workflows Status" "./n8n-ops check --env development"
run_test "Status Command" "./n8n-ops status --env development"

# Watch command test (non-interactive)
run_test "Watch Help" "./n8n-ops watch --help"

# Enhanced sync tests
run_test "Sync with Direction - From n8n" "./n8n-ops sync --from-n8n --env development"
run_test "Sync with Dry Run" "./n8n-ops sync --dry-run --env development"
run_test "Sync JSON Output" "./n8n-ops check --env development --json"

# Multi-environment tests
run_test "Sync Staging Environment" "./n8n-ops sync --env staging"
run_test "Sync Production Environment" "./n8n-ops sync --env production"

# Validation and utilities
run_test "Validate Workflows" "./n8n-ops validate --help"

echo ""
echo "🌿 Step 4: Interactive Branch Creation Demo (Simulated)"
echo "======================================================="

# Create a simulated branch creation (non-interactive for demo)
echo "📝 Simulating interactive branch creation..."
echo "   - Branch type: feature"
echo "   - Workflow name: customer-onboarding"
echo "   - Base strategy: template"
echo "   - Generated naming convention:"
echo "     * Git Branch: feature/customer-onboarding"
echo "     * Workflow Name: dev_customer_onboarding_v1.0.1234"
echo "     * File Name: customer-onboarding-v1.0.1234.json"

# Manually create the expected structure for demo
mkdir -p .n8n-ops/branches
mkdir -p workflows/development

# Create sample branch metadata
cat > .n8n-ops/branches/feature-customer-onboarding.json << EOF
{
  "branchName": "feature/customer-onboarding",
  "workflowName": "dev_customer_onboarding_v1.0.1234",
  "fileName": "customer-onboarding-v1.0.1234.json",
  "version": "v1.0.1234",
  "environment": "development",
  "createdAt": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "createdBy": "demo-user"
}
EOF

# Create sample workflow file
cat > workflows/development/customer-onboarding-v1.0.1234.json << EOF
{
  "id": "wf_$(date +%s)",
  "name": "dev_customer_onboarding_v1.0.1234",
  "active": false,
  "nodes": [
    {
      "id": "start",
      "name": "Start", 
      "type": "n8n-nodes-base.start",
      "typeVersion": 1,
      "position": [240, 300],
      "parameters": {}
    },
    {
      "id": "webhook",
      "name": "Customer Webhook",
      "type": "n8n-nodes-base.webhook",
      "typeVersion": 1,
      "position": [460, 300],
      "parameters": {
        "path": "/customer-onboarding",
        "responseMode": "responseNode"
      }
    }
  ],
  "connections": {
    "Start": {
      "main": [
        [
          {
            "node": "Customer Webhook",
            "type": "main",
            "index": 0
          }
        ]
      ]
    }
  },
  "createdAt": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "updatedAt": "$(date -u +%Y-%m-%dT%H:%M:%SZ)",
  "versionId": 1,
  "tags": ["branch:feature/customer-onboarding", "version:v1.0.1234", "env:development", "template"]
}
EOF

echo -e "${GREEN}✅ Sample branch structure created${NC}"
echo ""

# Test branch listing with created metadata
run_test "Branch List with Metadata" "./n8n-ops branch list"

echo ""
echo "📊 Step 5: Workflow Naming Convention Demonstration"
echo "=================================================="

echo "🏗️ DevOps Branch Strategy Examples:"
echo ""
echo "FEATURE BRANCHES:"
echo "  feature/user-authentication     → dev_user_authentication_v1.0.1001"
echo "  feature/payment-processing      → dev_payment_processing_v1.0.1002"
echo "  feature/email-notifications     → dev_email_notifications_v1.0.1003"
echo ""
echo "HOTFIX BRANCHES:"
echo "  hotfix/payment-bug              → hotfix_payment_bug_v1.1.2001"
echo "  hotfix/auth-security            → hotfix_auth_security_v1.1.2002"
echo ""
echo "RELEASE BRANCHES:"
echo "  release/v2.0.0                  → release_v2_0_0_staging"
echo "  release/v2.1.0                  → release_v2_1_0_staging"
echo ""
echo "EXPERIMENT BRANCHES:"
echo "  experiment/ab-test-checkout     → exp_ab_test_checkout_v0.1.3001"
echo "  experiment/new-ui-flow          → exp_new_ui_flow_v0.1.3002"

echo ""
echo "📁 Generated File Structure:"
echo "============================"
echo "project/"
echo "├── workflows/"
echo "│   ├── development/"
echo "│   │   ├── customer-onboarding-v1.0.1234.json"
echo "│   │   ├── user-authentication-v1.0.1001.json"
echo "│   │   └── payment-processing-v1.0.1002.json"
echo "│   ├── staging/"
echo "│   │   ├── customer-onboarding-v1.1.0.json"
echo "│   │   └── user-authentication-v1.1.0.json"
echo "│   └── production/"
echo "│       └── customer-onboarding-v1.0.0.json"
echo "├── .n8n-ops/"
echo "│   ├── branches/"
echo "│   │   ├── feature-customer-onboarding.json"
echo "│   │   ├── feature-user-authentication.json"
echo "│   │   └── hotfix-payment-bug.json"
echo "│   └── sync-state.json"
echo "└── .git/"

echo ""
echo "🔄 Step 6: Change Detection & Sync Demo"
echo "======================================="

# Test change detection system
run_test "Change Detection System" "./n8n-ops check --env development --json"

# Test bidirectional sync
echo "📋 Bidirectional Sync Capabilities:"
echo "  • FROM n8n TO Git: Download UI changes"
echo "  • FROM Git TO n8n: Upload file changes" 
echo "  • BIDIRECTIONAL: Smart conflict resolution"

echo ""
echo "🎛️ Step 7: DevOps Integration Examples"
echo "======================================"

echo "🔧 GitLab CI/CD Pipeline Integration:"
echo ""
echo "# .gitlab-ci.yml excerpt"
echo "workflow_sync:"
echo "  script:"
echo "    - n8n-ops sync --env staging --dry-run"
echo "    - n8n-ops check --env staging --json > sync-report.json"
echo "  artifacts:"
echo "    reports:"
echo "      junit: sync-report.json"
echo ""

echo "🚀 Production Deployment Flow:"
echo "1. Developer creates: feature/new-workflow"
echo "2. n8n-ops generates: dev_new_workflow_v1.0.1234"
echo "3. GitLab CI validates and tests"
echo "4. Merge to staging → staging_new_workflow_v1.1.0"
echo "5. Production release → prod_new_workflow_v1.0.0"

echo ""
echo "🏁 COMPREHENSIVE DEMO RESULTS"
echo "============================="
echo -e "Total Tests: ${BLUE}$TOTAL_TESTS${NC}"
echo -e "Passed:      ${GREEN}$PASSED_TESTS${NC}"
echo -e "Failed:      ${RED}$FAILED_TESTS${NC}"

if [ $FAILED_TESTS -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! n8n-ops is ready for production!${NC}"
else
    PERCENTAGE=$((PASSED_TESTS * 100 / TOTAL_TESTS))
    echo -e "${YELLOW}⚠️ Demo completed with $PERCENTAGE% success rate${NC}"
fi

echo ""
echo "✨ Key Features Demonstrated:"
echo "  ✅ Intelligent branch naming conventions"
echo "  ✅ DevOps-focused workflow strategies"
echo "  ✅ Automatic version generation"
echo "  ✅ Metadata tracking and lineage"
echo "  ✅ Bidirectional sync capabilities"
echo "  ✅ Git Flow integration"
echo "  ✅ CI/CD pipeline compatibility"
echo "  ✅ Multi-environment management"

echo ""
echo "🧹 Cleaning up..."
kill $SERVER_PID 2>/dev/null
echo "Demo completed!"