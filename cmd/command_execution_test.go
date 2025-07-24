package cmd

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// executeRoot runs the root command with the provided arguments and
// captures stdout output for assertions.
func executeRoot(args ...string) (string, error) {
	// Capture stdout
	reader, writer, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = writer

	rootCmd.SetArgs(args)
	cfg = nil
	logger = nil
	err := rootCmd.Execute()

	writer.Close()
	os.Stdout = old
	out, _ := io.ReadAll(reader)
	reader.Close()

	return string(out), err
}

func TestDeployCommandExecute(t *testing.T) {
	dir := t.TempDir()
	output, err := executeRoot("deploy", "--demo", "--output", dir, "--force")
	if err != nil {
		t.Fatalf("deploy command failed: %v", err)
	}
	if !strings.Contains(output, "Sync completed") {
		t.Errorf("unexpected deploy output: %s", output)
	}
}

func TestSyncCommandExecute(t *testing.T) {
	dir := t.TempDir()
	output, err := executeRoot("sync", "--demo", "--output", dir, "--force")
	if err != nil {
		t.Fatalf("sync command failed: %v", err)
	}
	if !strings.Contains(output, "Sync completed") {
		t.Errorf("unexpected sync output: %s", output)
	}
}

func TestValidateCommandExecute(t *testing.T) {
	tmpDir := t.TempDir()
	file := filepath.Join(tmpDir, "workflow.json")
	json := `{"id":"1","name":"Test Workflow","nodes":[{"id":"1","name":"Start","type":"n8n-nodes-base.start","typeVersion":1,"position":[0,0]}],"connections":{}}`
	if err := os.WriteFile(file, []byte(json), 0o644); err != nil {
		t.Fatalf("failed to write workflow file: %v", err)
	}

	output, err := executeRoot("validate", file)
	if err != nil {
		t.Fatalf("validate command failed: %v", err)
	}
	if !strings.Contains(output, "workflow files are valid") {
		t.Errorf("unexpected validate output: %s", output)
	}
}

func TestStatusCommandExecute(t *testing.T) {
	output, err := executeRoot("status", "--json")
	if err != nil {
		t.Fatalf("status command failed: %v", err)
	}
	if !strings.Contains(output, "\"environment\"") {
		t.Errorf("unexpected status output: %s", output)
	}
}

func TestExecuteInvalidConfig(t *testing.T) {
	dir := t.TempDir()
	cfgPath := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(cfgPath, []byte("environments: bad"), 0o644); err != nil {
		t.Fatalf("failed to write config: %v", err)
	}

	_, err := executeRoot("--config", cfgPath)
	if err == nil {
		t.Fatal("expected error for invalid config")
	}
}
