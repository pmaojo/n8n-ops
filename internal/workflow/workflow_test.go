package workflow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func minimalWorkflow() *Workflow {
	return &Workflow{
		Name: "Test",
		Nodes: []Node{
			{
				Name:     "Start",
				Type:     "n8n-nodes-base.manualTrigger",
				Position: []float64{0, 0},
			},
		},
	}
}

func TestValidateWorkflowSuccess(t *testing.T) {
	wf := minimalWorkflow()
	if err := ValidateWorkflow(wf); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestValidateWorkflowFailureMissingName(t *testing.T) {
	wf := minimalWorkflow()
	wf.Name = ""
	if err := ValidateWorkflow(wf); err == nil {
		t.Fatal("expected validation error for missing workflow name")
	}
}

func TestValidateWorkflowStrictConnectivityFailure(t *testing.T) {
	wf := minimalWorkflow()
	wf.Nodes = append(wf.Nodes, Node{Name: "Func", Type: "n8n-nodes-base.function", Position: []float64{200, 0}})
	if err := ValidateWorkflowStrict(wf); err == nil {
		t.Fatal("expected connectivity error for second node")
	}
}

func TestValidateWorkflowFileSuccess(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "wf.json")

	data, err := json.Marshal(minimalWorkflow())
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := ValidateWorkflowFile(path); err != nil {
		t.Fatalf("expected file to validate, got %v", err)
	}
}

func TestValidateWorkflowFileInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	if err := os.WriteFile(path, []byte("{"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	if err := ValidateWorkflowFile(path); err == nil {
		t.Fatal("expected JSON parsing error")
	}
}

func TestValidateWorkflowFileMissing(t *testing.T) {
	if err := ValidateWorkflowFile("no_such_file.json"); err == nil {
		t.Fatal("expected error for missing file")
	}
}
