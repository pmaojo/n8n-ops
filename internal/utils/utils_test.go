package utils

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	wf "github.com/pmaojo/n8n-ops/internal/workflow"
)

func TestSanitizeFilename(t *testing.T) {
	name := "Invalid:File/Name*?"
	sanitized := SanitizeFilename(name)
	if strings.ContainsAny(sanitized, ":/*? ") {
		t.Fatalf("sanitized name still contains invalid characters: %s", sanitized)
	}
}

func TestWriteAndLoadWorkflow(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wf.json")
	wfIn := &wf.Workflow{
		Name:  "FileTest",
		Nodes: []wf.Node{{Name: "Start", Type: "n8n-nodes-base.manualTrigger", Position: []float64{0, 0}}},
	}
	if err := WriteWorkflowToFile(wfIn, path); err != nil {
		t.Fatalf("write workflow: %v", err)
	}

	wfOut, err := LoadWorkflowFromFile(path)
	if err != nil {
		t.Fatalf("load workflow: %v", err)
	}

	if !reflect.DeepEqual(wfIn, wfOut) {
		t.Fatalf("loaded workflow does not match written one")
	}
}

func TestLoadWorkflowFromFileNotFound(t *testing.T) {
	if _, err := LoadWorkflowFromFile("/no/such/file.json"); err == nil {
		t.Fatal("expected error for missing workflow file")
	}
}

func TestBackupFile(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "f.json")
	os.WriteFile(file, []byte("{}"), 0o644)

	if err := BackupFile(file); err != nil {
		t.Fatalf("backup: %v", err)
	}
	matches, err := filepath.Glob(file + ".backup_*")
	if err != nil || len(matches) == 0 {
		t.Fatal("expected backup file to exist")
	}
}

func TestCheckN8nCredentials(t *testing.T) {
	oldURL := os.Getenv("N8N_DEVELOPMENT_URL")
	oldKey := os.Getenv("N8N_DEVELOPMENT_API_KEY")
	defer func() {
		os.Setenv("N8N_DEVELOPMENT_URL", oldURL)
		os.Setenv("N8N_DEVELOPMENT_API_KEY", oldKey)
	}()

	os.Setenv("N8N_DEVELOPMENT_URL", "http://localhost")
	os.Setenv("N8N_DEVELOPMENT_API_KEY", "key")

	if err := CheckN8nCredentials("development"); err != nil {
		t.Fatalf("expected credentials to pass, got %v", err)
	}

	os.Unsetenv("N8N_DEVELOPMENT_URL")
	os.Unsetenv("N8N_DEVELOPMENT_API_KEY")
	if err := CheckN8nCredentials("development"); err == nil {
		t.Fatal("expected error when credentials are missing")
	}
}
