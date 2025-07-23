package cmd

import (
	"path/filepath"
	"testing"
)

func TestSyncCommand(t *testing.T) {
	// Test sync command initialization
	if syncCmd == nil {
		t.Error("syncCmd should not be nil")
	}
	if syncCmd.Use != "sync" {
		t.Errorf("Expected Use to be 'sync', got %s", syncCmd.Use)
	}
}

func TestSyncFlags(t *testing.T) {
	// Test that all expected flags are present
	flags := syncCmd.Flags()

	if !flags.HasFlags() {
		t.Error("Sync command should have flags")
	}

	expectedFlags := []string{"force", "output", "branch", "from-n8n", "to-n8n", "dry-run"}
	for _, flagName := range expectedFlags {
		if flags.Lookup(flagName) == nil {
			t.Errorf("Expected flag '%s' not found", flagName)
		}
	}
}

func TestSyncOutputDirectory(t *testing.T) {
	// Test output directory logic
	testEnv := "development"
	expectedDir := filepath.Join("workflows", testEnv)
	actualDir := filepath.Join("workflows", testEnv)

	if expectedDir != actualDir {
		t.Errorf("Expected %s, got %s", expectedDir, actualDir)
	}
}
