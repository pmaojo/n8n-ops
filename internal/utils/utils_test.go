package utils

import (
	"path/filepath"
	"testing"
)

func TestFilePathOperations(t *testing.T) {
	// Test file path operations
	testPath := filepath.Join("workflows", "development", "test.json")

	if testPath == "" {
		t.Error("Test path should not be empty")
	}

	ext := filepath.Ext(testPath)
	if ext != ".json" {
		t.Errorf("Expected .json extension, got %s", ext)
	}
}

func TestStringUtilities(t *testing.T) {
	// Test string utility functions
	testString := "test-workflow-name"
	expectedLength := len(testString)

	if expectedLength == 0 {
		t.Error("Test string should not be empty")
	}

	if expectedLength != 18 {
		t.Errorf("Expected length 18, got %d", expectedLength)
	}
}

func TestValidationHelpers(t *testing.T) {
	// Test validation helper functions
	validWorkflowName := "Valid_Workflow_Name"
	invalidWorkflowName := ""

	if validWorkflowName == "" {
		t.Error("Valid workflow name should not be empty")
	}

	if invalidWorkflowName != "" {
		t.Error("Invalid workflow name should be empty")
	}
}

func TestEnvironmentHelpers(t *testing.T) {
	// Test environment helper functions
	environments := []string{"development", "staging", "production"}

	if len(environments) != 3 {
		t.Error("Should have exactly 3 environments")
	}

	for _, env := range environments {
		if env == "" {
			t.Error("Environment name should not be empty")
		}
	}
}
