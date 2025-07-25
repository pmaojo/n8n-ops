package credentials

import (
	"encoding/json"
	"os"
	"testing"
)

func TestMapCredentialsForEnvironment(t *testing.T) {
	ncm := NewN8nCredentialManager("development")
	wf := map[string]any{
		"id": "1",
		"nodes": []map[string]any{
			{
				"id":   "node1",
				"type": "smtp",
				"credentials": map[string]map[string]string{
					"smtp": {"id": "smtp_dev_001", "name": "SMTP"},
				},
			},
		},
	}
	data, err := json.Marshal(wf)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	out, err := ncm.MapCredentialsForEnvironment(data, "staging")
	if err != nil {
		t.Fatalf("map: %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(out, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	nodes := result["nodes"].([]any)
	creds := nodes[0].(map[string]any)["credentials"].(map[string]any)["smtp"].(map[string]any)
	if creds["id"] != "smtp_staging_001" {
		t.Fatalf("expected id mapped to staging, got %v", creds["id"])
	}
}

func TestGenerateCredentialReport(t *testing.T) {
	ncm := NewN8nCredentialManager("production")
	report, err := ncm.GenerateCredentialReport()
	if err != nil {
		t.Fatalf("report: %v", err)
	}
	if report.Environment != "production" {
		t.Fatalf("unexpected env %s", report.Environment)
	}
	if len(report.Credentials) == 0 {
		t.Fatal("expected credentials in report")
	}
	smtp := report.Credentials["SMTP Email"]
	if smtp.ID == "" || smtp.Status != "active" {
		t.Fatalf("invalid credential info: %+v", smtp)
	}
}

func TestValidateAllCredentials(t *testing.T) {
	scm := NewSecureCredentialManager("development")
	// Set a single env var to verify counting
	os.Setenv("SMTP_HOST_DEVELOPMENT", "localhost")
	defer os.Unsetenv("SMTP_HOST_DEVELOPMENT")

	val, err := scm.ValidateAllCredentials()
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if val.Required == 0 || val.Missing == 0 {
		t.Fatalf("expected some required and missing credentials")
	}
	if len(val.MissingCredentials) == 0 {
		t.Fatalf("expected list of missing credentials")
	}
}
