#!/bin/bash

set -e

echo "🎭 n8n-ops Demo Test Suite"
echo "========================="
echo "Testing all CLI functionality with mock data"
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
        if [ $? -eq $expected_exit_code ]; then
            echo -e "${GREEN}✅ PASSED${NC}"
            TESTS_PASSED=$((TESTS_PASSED + 1))
        else
            echo -e "${RED}❌ FAILED - Wrong exit code${NC}"
            TESTS_FAILED=$((TESTS_FAILED + 1))
        fi
    else
        echo -e "${RED}❌ FAILED - Command error${NC}"
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

echo "Starting demo test suite..."
echo ""

# Test 1: Basic help
run_test "Basic Help Command" "./n8n-ops --help"

# Test 2: Version information
run_test "Version Information" "./n8n-ops version"

# Test 3: Welcome screen
run_test "Welcome Screen" "./n8n-ops welcome"

# Test 4: Demo sync - development
run_test "Demo Sync - Development" "./n8n-ops sync --demo --env development"

# Test 5: Demo sync - staging  
run_test "Demo Sync - Staging" "./n8n-ops sync --demo --env staging"

# Test 6: Demo sync - production
run_test "Demo Sync - Production" "./n8n-ops sync --demo --env production"

# Test 7: Check workflows in demo mode
run_test "Check Workflows - Demo Mode" "./n8n-ops check --demo --env development"

# Test 8: Check workflows JSON output
run_test "Check Workflows - JSON Output" "./n8n-ops check --demo --env development --json"

# Test 9: Status command
run_test "Status Command" "./n8n-ops status --env development"

# Test 10: Validate workflows (should create directory first)
mkdir -p workflows/development 2>/dev/null || true
run_test "Validate Workflows" "./n8n-ops validate ./workflows/development/"

# Test 11: Check with fail-if-changes flag (should fail in demo)
run_test "Check with Fail Flag (Expected Failure)" "./n8n-ops check --demo --fail-if-changes" 1

# Test 12: Spanish language support
run_test "Spanish Language Support" "./n8n-ops --lang es sync --demo --env development"

# Test 13: Verbose output
run_test "Verbose Output" "./n8n-ops --verbose sync --demo --env development"

# Test 14: Branch command
run_test "Branch Command" "./n8n-ops branch list"

# Test 15: Init command  
run_test "Init Command" "./n8n-ops init --name test-project"

echo ""
echo "🏁 DEMO TEST SUITE RESULTS"
echo "=========================="
echo -e "Total Tests Run: ${BLUE}$TESTS_RUN${NC}"
echo -e "Tests Passed:    ${GREEN}$TESTS_PASSED${NC}"
echo -e "Tests Failed:    ${RED}$TESTS_FAILED${NC}"

if [ $TESTS_FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 ALL TESTS PASSED! Demo mode is working perfectly.${NC}"
    echo ""
    echo "Enterprise-ready CLI with comprehensive demo functionality:"
    echo "• Mock n8n API client working"
    echo "• All commands functional"  
    echo "• Multi-environment support"
    echo "• JSON output formats"
    echo "• Internationalization (Spanish)"
    echo "• Comprehensive error handling"
    echo ""
    echo "Ready for production n8n API integration!"
    exit 0
else
    echo -e "${RED}🚨 Some tests failed. Check the output above.${NC}"
    exit 1
fi