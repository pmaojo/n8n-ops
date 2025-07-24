package workflow

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
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
	logger := logrus.New()
	logger.SetOutput(io.Discard)

	wf := minimalWorkflow()
	wf.Nodes = append(wf.Nodes, Node{Name: "Func", Type: "n8n-nodes-base.function", Position: []float64{200, 0}})

	if err := ValidateWorkflowStrict(wf, logger); err == nil {
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

func connectedWorkflow() *Workflow {
	wf := minimalWorkflow()
	wf.Nodes = append(wf.Nodes, Node{
		Name:     "Func",
		Type:     "n8n-nodes-base.function",
		Position: []float64{200, 0},
	})
	wf.Connections = map[string]interface{}{
		"Start": map[string]interface{}{
			"main": []interface{}{
				[]interface{}{
					map[string]interface{}{
						"node":  "Func",
						"type":  "main",
						"index": 0,
					},
				},
			},
		},
	}
	return wf
}

func TestCloneWorkflow(t *testing.T) {
	original := connectedWorkflow()
	clone := original.Clone()
	if clone == original {
		t.Fatal("Clone returned same instance")
	}
	clone.Name = "Changed"
	clone.Nodes[0].Name = "Modified"
	if original.Name == "Changed" {
		t.Error("changing clone should not affect original name")
	}
	if original.Nodes[0].Name == "Modified" {
		t.Error("changing clone node should not affect original")
	}
}

func TestGetNodeByName(t *testing.T) {
	wf := connectedWorkflow()
	node := wf.GetNodeByName("Func")
	if node == nil || node.Name != "Func" {
		t.Fatalf("expected to find node 'Func', got %v", node)
	}
	if wf.GetNodeByName("Missing") != nil {
		t.Error("expected nil for unknown node")
	}
}

func TestHasTag(t *testing.T) {
	wf := minimalWorkflow()
	if wf.HasTag("demo") {
		t.Fatal("unexpected tag found")
	}
	wf.Tags = []Tag{{Name: "demo"}}
	if !wf.HasTag("demo") {
		t.Fatal("tag 'demo' should exist")
	}
}

func TestValidateWorkflowStrictConnectivitySuccess(t *testing.T) {
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	wf := connectedWorkflow()
	if err := ValidateWorkflowStrict(wf, logger); err != nil {
		t.Fatalf("expected no connectivity error, got %v", err)
	}
}

func TestValidateWorkflowBatch(t *testing.T) {
	valid := minimalWorkflow()
	invalid := minimalWorkflow()
	invalid.Name = ""

	errs := ValidateWorkflowBatch([]*Workflow{valid, invalid})
	if len(errs) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errs))
	}
	if !strings.Contains(errs[0].Error(), "workflow 1") {
		t.Errorf("unexpected error message: %v", errs[0])
	}
}
