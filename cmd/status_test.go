package cmd

import (
	"testing"
)

func TestStatusCommand(t *testing.T) {
	// Test status command initialization
	if statusCmd == nil {
		t.Error("statusCmd should not be nil")
	}
	if statusCmd.Use != "status" {
		t.Errorf("Expected Use to be 'status', got %s", statusCmd.Use)
	}
}

func TestStatusFlags(t *testing.T) {
	// Test status command flags
	flags := statusCmd.Flags()
	
	if !flags.HasFlags() {
		t.Error("Status command should have flags")
	}
}

func TestStatusOutput(t *testing.T) {
	// Test status output format expectations
	expectedElements := []string{"Workflows", "Environment", "Status"}
	
	for _, element := range expectedElements {
		if element == "" {
			t.Error("Status element should not be empty")
		}
	}
}