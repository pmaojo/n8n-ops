package issues

import (
	"testing"
	"time"
)

func TestWorkflowFailure(t *testing.T) {
	failure := &WorkflowFailure{
		WorkflowID:   "wf_123",
		WorkflowName: "Test Workflow",
		ExecutionID:  "exec_456",
		Environment:  "production",
		FailedAt:     time.Now(),
		ErrorMessage: "Connection timeout",
		RetryCount:   3,
		PipelineID:   "pipeline_789",
		Branch:       "main",
	}

	if failure.WorkflowID != "wf_123" {
		t.Errorf("Expected WorkflowID wf_123, got %s", failure.WorkflowID)
	}

	if failure.Environment != "production" {
		t.Errorf("Expected Environment production, got %s", failure.Environment)
	}
}

func TestRecoveryInfo(t *testing.T) {
	recovery := &RecoveryInfo{
		RecoveredAt:  time.Now(),
		RecoveryType: "auto",
		RecoveredBy:  "n8n-ops",
		NewCommitSHA: "abc123",
		Notes:        "Automatic recovery after timeout issue resolved",
	}

	if recovery.RecoveryType != "auto" {
		t.Errorf("Expected RecoveryType auto, got %s", recovery.RecoveryType)
	}
}

func TestGitLabIssueManager(t *testing.T) {
	manager := NewGitLabIssueManager("https://gitlab.example.com", "123", "token")

	if manager.baseURL != "https://gitlab.example.com" {
		t.Errorf("Expected baseURL https://gitlab.example.com, got %s", manager.baseURL)
	}

	if manager.projectID != "123" {
		t.Errorf("Expected projectID 123, got %s", manager.projectID)
	}

	if manager.token != "token" {
		t.Errorf("Expected token 'token', got %s", manager.token)
	}
}

func TestBuildIssueLabels(t *testing.T) {
	manager := NewGitLabIssueManager("https://gitlab.com", "123", "token")

	failure := &WorkflowFailure{
		WorkflowID:  "wf_123",
		Environment: "production",
		RetryCount:  8,
		NodeName:    "HTTPRequest",
	}

	labels := manager.buildIssueLabels(failure)

	expectedLabels := []string{
		"workflow-failure",
		"automated",
		"env:production",
		"workflow:wf_123",
		"severity:high",
		"node:HTTPRequest",
	}

	if len(labels) != len(expectedLabels) {
		t.Errorf("Expected %d labels, got %d", len(expectedLabels), len(labels))
	}

	for i, expected := range expectedLabels {
		if i >= len(labels) || labels[i] != expected {
			t.Errorf("Expected label %s at position %d, got %s", expected, i, labels[i])
		}
	}
}

func TestBuildIssueDescription(t *testing.T) {
	manager := NewGitLabIssueManager("https://gitlab.com", "123", "token")

	failure := &WorkflowFailure{
		WorkflowID:   "wf_123",
		WorkflowName: "Data Processing",
		ExecutionID:  "exec_456",
		Environment:  "production",
		FailedAt:     time.Date(2025, 7, 23, 10, 30, 0, 0, time.UTC),
		ErrorMessage: "Database connection failed",
		RetryCount:   5,
		NodeName:     "PostgreSQL",
		PipelineURL:  "https://gitlab.com/project/-/pipelines/123",
		PipelineID:   "123",
		CommitSHA:    "abc123def",
		Branch:       "main",
	}

	description := manager.buildIssueDescription(failure)

	if description == "" {
		t.Error("Description should not be empty")
	}

	// Check if key information is included
	testCases := []string{
		"Data Processing",
		"wf_123",
		"production",
		"Database connection failed",
		"PostgreSQL",
		"abc123def",
		"main",
	}

	for _, testCase := range testCases {
		if !contains(description, testCase) {
			t.Errorf("Description should contain '%s'", testCase)
		}
	}
}

func TestBuildRecoveryComment(t *testing.T) {
	manager := NewGitLabIssueManager("https://gitlab.com", "123", "token")

	recovery := &RecoveryInfo{
		RecoveredAt:  time.Date(2025, 7, 23, 11, 0, 0, 0, time.UTC),
		RecoveryType: "manual",
		RecoveredBy:  "admin@example.com",
		NewCommitSHA: "def456ghi",
		Notes:        "Fixed database connection string",
	}

	comment := manager.buildRecoveryComment(recovery)

	if comment == "" {
		t.Error("Recovery comment should not be empty")
	}

	testCases := []string{
		"Workflow Recovery Detected",
		"manual",
		"admin@example.com",
		"def456ghi",
		"Fixed database connection string",
	}

	for _, testCase := range testCases {
		if !contains(comment, testCase) {
			t.Errorf("Recovery comment should contain '%s'", testCase)
		}
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
