package credentials

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func TestSyncCredentialsToN8nLogs(t *testing.T) {
	cfg := map[string]any{
		"environments": map[string]any{
			"testenv": map[string]any{
				"workflow_credentials": []map[string]any{
					{
						"id":   "cred1",
						"name": "Test Credential",
						"type": "smtp",
						"data": map[string]string{"host": "smtp.example.com"},
					},
				},
			},
		},
	}

	tmp, err := os.CreateTemp("", "cfg*.yaml")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	defer os.Remove(tmp.Name())
	enc := yaml.NewEncoder(tmp)
	if err := enc.Encode(cfg); err != nil {
		t.Fatalf("encode: %v", err)
	}
	enc.Close()
	tmp.Close()

	os.Setenv("N8N_OPS_CONFIG", tmp.Name())
	defer os.Unsetenv("N8N_OPS_CONFIG")

	var buf bytes.Buffer
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{DisableTimestamp: true})
	logger.SetOutput(&buf)

	cm := NewCredentialManager("testenv", logger)
	if err := cm.SyncCredentialsToN8n(); err != nil {
		t.Fatalf("sync: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "found workflow credentials") {
		t.Errorf("expected summary log, got: %s", output)
	}
	if !strings.Contains(output, "Test Credential") {
		t.Errorf("expected credential log, got: %s", output)
	}
}
