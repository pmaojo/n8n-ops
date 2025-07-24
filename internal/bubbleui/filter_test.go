package bubbleui

import "testing"

func TestFilterWorkflows(t *testing.T) {
	workflows := []WorkflowStatus{
		{ID: "1", Name: "Alpha", Status: "active"},
		{ID: "2", Name: "Beta", Status: "inactive"},
		{ID: "3", Name: "Gamma", Status: "active"},
	}

	filtered := filterWorkflows(workflows, "a")
	if len(filtered) != 3 {
		t.Fatalf("expected 3 workflows, got %d", len(filtered))
	}

	filtered = filterWorkflows(workflows, "BETA")
	if len(filtered) != 1 || filtered[0].Name != "Beta" {
		t.Fatalf("case insensitive filtering failed: %+v", filtered)
	}
}
