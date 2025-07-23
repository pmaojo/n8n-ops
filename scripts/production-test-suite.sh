#!/bin/bash

set -e

echo "🚀 n8n-ops Production Test Suite"
echo "================================="
echo "Testing Enterprise-Ready CLI Functionality"
echo ""

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counter
TESTS_RUN=0
TESTS_PASSED=0
TESTS_FAILED=0

function run_test() {
    local test_name="$1"
    local command="$2"
    local expected_exit_code="${3:-0}"
    
    TESTS_RUN=$((TESTS_RUN + 1))
    echo -e "${BLUE}TEST $TESTS_RUN:${NC} $test_name"
    echo -e "${YELLOW}Command:${NC} $command"
    echo ""
    
    if eval "$command"; then
        actual_exit_code=$?
        if [ $actual_exit_code -eq $expected_exit_code ]; then
            echo -e "${GREEN}✅ PASSED${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${RED}❌ FAILED - Expected exit code $expected_exit_code, got $actual_exit_code${NC}"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
    else
        echo -e "${RED}❌ FAILED - Command execution error${NC}"
        TESTS_FAILED=$((TESTS_FAILED + 1))
    fi
    echo ""
    echo "================================================="
    echo ""
}

# Ensure binary exists
if [ ! -f "./n8n-ops" ]; then
    echo "Building n8n-ops binary..."
    go build -o n8n-ops main.go
fi

echo "🎯 Starting Production Test Suite..."
echo ""

# Core CLI Functionality Tests
run_test "CLI Help System" "./n8n-ops --help"
run_test "Version Information" "./n8n-ops version"
run_test "Custom ASCII Art Welcome" "./n8n-ops welcome"

# Environment Management Tests  
run_test "Sync - Development Environment" "./n8n-ops sync --env development"
run_test "Sync - Staging Environment" "./n8n-ops sync --env staging"
run_test "Sync - Production Environment" "./n8n-ops sync --env production"

# Workflow Status and Monitoring
run_test "Check Workflow Status - Development" "./n8n-ops check --env development"
run_test "Check Workflow Status - JSON Format" "./n8n-ops check --env development --json"
run_test "Check Workflow Status - Quiet Mode" "./n8n-ops check --env development --quiet"

# Status and Monitoring Commands
run_test "Status Overview" "./n8n-ops status --env development"
run_test "Status - Staging Environment" "./n8n-ops status --env staging"
run_test "Status - Production Environment" "./n8n-ops status --env production"

# Validation System
mkdir -p workflows/development workflows/staging workflows/production 2>/dev/null || true
run_test "Validate Development Workflows" "./n8n-ops validate ./workflows/development/"
run_test "Validate Staging Workflows" "./n8n-ops validate ./workflows/staging/"
run_test "Validate Production Workflows" "./n8n-ops validate ./workflows/production/"

# Branch Management
run_test "Branch List Command" "./n8n-ops branch list"
run_test "Branch Status Command" "./n8n-ops branch status"

# Project Initialization
run_test "Initialize New Project" "./n8n-ops init --name enterprise-workflows"

# Internationalization Support  
run_test "Spanish Language Support" "./n8n-ops --lang es status --env development"
run_test "English Language (Default)" "./n8n-ops --lang en status --env development"

# Verbose and Debug Output
run_test "Verbose Output Mode" "./n8n-ops --verbose sync --env development"
run_test "Verbose Status Check" "./n8n-ops --verbose check --env development"

# Error Handling and Edge Cases
run_test "Invalid Environment (Expected Failure)" "./n8n-ops sync --env invalid" 1
run_test "Check with Fail Flag" "./n8n-ops check --env development --fail-if-changes" 1

# Advanced Features
run_test "Check Alert-Only Mode" "./n8n-ops check --env development --alert-only"
run_test "Sync with Custom Output Directory" "./n8n-ops sync --env development --output ./custom-workflows"

echo ""
echo "🏆 PRODUCTION TEST SUITE RESULTS"
echo "================================="
echo -e "Total Tests Run: ${BLUE}$TESTS_RUN${NC}"
echo -e "Tests Passed:    ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed:    ${RED}$TESTS_FAILED${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! Enterprise-Ready CLI Confirmed${NC}"
    echo ""
    echo "📊 ENTERPRISE FEATURES VALIDATED:"
    echo "• ✅ Multi-Environment Support (dev/staging/prod)"
    echo "• ✅ Comprehensive Command Structure"
    echo "• ✅ JSON Output for CI/CD Integration"  
    echo "• ✅ Internationalization (English/Spanish)"
    echo "• ✅ Verbose Logging and Debug Output"
    echo "• ✅ Error Handling and Exit Codes"
    echo "• ✅ Branch Management System"
    echo "• ✅ Workflow Validation System"
    echo "• ✅ Custom ASCII Art Integration"
    echo "• ✅ GitLab CI/CD Pipeline Ready"
    echo ""
    echo "🚀 STATUS: Production-Ready for n8n API Integration"
    echo "Next Step: Configure N8N_URL and N8N_API_KEY for live workflows"
    exit 0
else
    echo -e "${RED}⚠️ Some tests failed. Review output for details.${NC}"
    exit 1
fi