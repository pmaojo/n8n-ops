package credentials

import (
	"gopkg.in/yaml.v3"
	"os"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/utils"
)

func TestGetWorkflowCredentialsFromConfig(t *testing.T) {
	cfg := map[string]any{
		"environments": map[string]any{
			"testenv": map[string]any{
				"workflow_credentials": []map[string]any{
					{
						"id":   "cred1",
						"name": "Test Credential",
						"type": "smtp",
						"data": map[string]string{
							"host": "smtp.test.com",
						},
					},
				},
			},
		},
	}

	tmpFile, err := os.CreateTemp("", "cred*.yaml")
	if err != nil {
		t.Fatalf("create temp config: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	enc := yaml.NewEncoder(tmpFile)
	if err := enc.Encode(cfg); err != nil {
		t.Fatalf("write yaml: %v", err)
	}
	enc.Close()
	tmpFile.Close()

	os.Setenv("N8N_OPS_CONFIG", tmpFile.Name())
	defer os.Unsetenv("N8N_OPS_CONFIG")

	cm := NewCredentialManager("testenv", utils.OSProvider{})
	creds, err := cm.GetWorkflowCredentials()
	if err != nil {
		t.Fatalf("GetWorkflowCredentials error: %v", err)
	}
	if len(creds) != 1 {
		t.Fatalf("expected 1 credential, got %d", len(creds))
	}
	if creds[0].ID != "cred1" || creds[0].Name != "Test Credential" {
		t.Errorf("unexpected credential: %+v", creds[0])
	}
}
