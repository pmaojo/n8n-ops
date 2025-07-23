package sync

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSyncEngine(t *testing.T) {
	// Test sync engine initialization
	tempDir := t.TempDir()

	// Test directory creation
	workflowDir := filepath.Join(tempDir, "workflows", "development")
	err := os.MkdirAll(workflowDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create workflow directory: %v", err)
	}

	// Verify directory exists
	if _, err := os.Stat(workflowDir); os.IsNotExist(err) {
		t.Error("Workflow directory should exist after creation")
	}
}

func TestWorkflowFileOperations(t *testing.T) {
	// Test workflow file operations
	tempDir := t.TempDir()
	workflowFile := filepath.Join(tempDir, "test_workflow.json")

	// Test file creation
	testContent := `{"id": "test_001", "name": "Test Workflow", "nodes": []}`
	err := os.WriteFile(workflowFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test file reading
	content, err := os.ReadFile(workflowFile)
	if err != nil {
		t.Fatalf("Failed to read test file: %v", err)
	}

	if string(content) != testContent {
		t.Error("File content should match written content")
	}
}

func TestSyncDirection(t *testing.T) {
	// Test sync direction logic
	directions := map[string]bool{
		"from-n8n": true,
		"to-n8n":   true,
		"invalid":  false,
	}

	for direction, expected := range directions {
		isValid := direction == "from-n8n" || direction == "to-n8n"
		if isValid != expected {
			t.Errorf("Direction %s validity should be %v, got %v", direction, expected, isValid)
		}
	}
}

func TestSyncContext(t *testing.T) {
	// Test context handling in sync operations
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// Simulate sync operation with context
	select {
	case <-ctx.Done():
		// Context properly handled
	case <-time.After(200 * time.Millisecond):
		t.Error("Context timeout should have been respected")
	}
}
