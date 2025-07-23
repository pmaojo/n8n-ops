package tutorial

import (
	"os"
	"testing"
)

func TestConfigExists(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)
	if ConfigExists() {
		t.Fatal("config should not exist in empty temp dir")
	}

	f, err := os.Create(tmpDir + "/.n8n-ops.yaml")
	if err != nil {
		t.Fatalf("failed to create temp config: %v", err)
	}
	f.Close()

	if !ConfigExists() {
		t.Error("expected config to be detected")
	}
}
