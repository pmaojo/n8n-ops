package cmd

import (
	"testing"
)

func TestDeployCommand(t *testing.T) {
	// Test deploy command is properly initialized
	if deployCmd == nil {
		t.Error("deployCmd should not be nil")
	}
	if deployCmd.Use != "deploy" {
		t.Errorf("Expected Use to be 'deploy', got %s", deployCmd.Use)
	}
}

func TestDeployAlias(t *testing.T) {
	// Test that deploy has a run function
	if deployCmd.RunE == nil {
		t.Error("Deploy command should have a RunE function")
	}
}

func TestDeployFlags(t *testing.T) {
	// Test deploy inherits sync flags
	flags := deployCmd.Flags()

	expectedFlags := []string{"force", "dry-run", "output", "branch"}
	for _, flagName := range expectedFlags {
		if flags.Lookup(flagName) == nil {
			t.Errorf("Expected flag '%s' not found", flagName)
		}
	}
}
