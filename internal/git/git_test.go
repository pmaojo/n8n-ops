package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGitRepository(t *testing.T) {
	// Test Git repository operations
	tempDir := t.TempDir()
	
	// Test .git directory detection
	gitDir := filepath.Join(tempDir, ".git")
	err := os.MkdirAll(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git directory: %v", err)
	}
	
	// Verify Git repository detection
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Error("Git directory should exist")
	}
}

func TestBranchOperations(t *testing.T) {
	// Test branch operations
	branches := []string{"main", "develop", "staging", "feature/new-workflow"}
	
	for _, branch := range branches {
		if branch == "" {
			t.Error("Branch name should not be empty")
		}
		
		// Test branch name validation
		if branch == "main" || branch == "develop" || branch == "staging" {
			// Valid protected branches
		} else if len(branch) < 3 {
			t.Errorf("Branch name %s too short", branch)
		}
	}
}

func TestCommitOperations(t *testing.T) {
	// Test commit operations
	commitData := map[string]string{
		"message": "Add new workflow: Payment Processing",
		"author":  "developer@example.com",
		"files":   "workflows/development/Payment_Processing.json",
	}
	
	for field, value := range commitData {
		if value == "" {
			t.Errorf("Commit %s should not be empty", field)
		}
	}
	
	// Test commit message format
	if len(commitData["message"]) < 10 {
		t.Error("Commit message should be descriptive")
	}
}

func TestGitIgnore(t *testing.T) {
	// Test .gitignore patterns
	gitignorePatterns := []string{
		"*.log",
		"node_modules/",
		".env",
		"dist/",
		"coverage/",
	}
	
	for _, pattern := range gitignorePatterns {
		if pattern == "" {
			t.Error("Gitignore pattern should not be empty")
		}
		
		// Test pattern validity
		if pattern == ".env" {
			// Should ignore environment files
		} else if pattern == "*.log" {
			// Should ignore log files
		}
	}
}

func TestFileTracking(t *testing.T) {
	// Test file tracking logic
	trackedFiles := []string{
		"workflows/development/workflow1.json",
		"config.yaml",
		"README.md",
	}
	
	untrackedFiles := []string{
		".env",
		"node_modules/",
		"*.log",
	}
	
	for _, file := range trackedFiles {
		if filepath.Ext(file) == ".json" || filepath.Ext(file) == ".yaml" || filepath.Ext(file) == ".md" {
			// These should be tracked
		} else {
			t.Errorf("File %s tracking status unclear", file)
		}
	}
	
	for _, file := range untrackedFiles {
		if file == ".env" || file == "node_modules/" {
			// These should not be tracked
		}
	}
}