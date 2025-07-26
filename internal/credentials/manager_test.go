package credentials

import (
	"path/filepath"
	"testing"
)

type mockProvider struct{ env map[string]string }

func (m mockProvider) Getenv(key string) string       { return m.env[key] }
func (m mockProvider) Setenv(key, value string) error { m.env[key] = value; return nil }

func TestNewCredentialManagerUsesProvider(t *testing.T) {
	mp := mockProvider{env: map[string]string{
		"N8N_OPS_CONFIG": "/tmp/custom.yaml",
		"HOME":           "/home/test",
	}}

	cm := NewCredentialManager(mp, "dev")
	if cm.ConfigPath != "/tmp/custom.yaml" {
		t.Fatalf("expected config path from provider, got %s", cm.ConfigPath)
	}
	if cm.Environment != "dev" {
		t.Fatalf("expected environment 'dev', got %s", cm.Environment)
	}
}

func TestNewCredentialManagerFallbackHome(t *testing.T) {
	mp := mockProvider{env: map[string]string{"HOME": "/home/foo"}}

	cm := NewCredentialManager(mp, "dev")
	expected := filepath.Join("/home/foo", ".n8n-ops.yaml")
	if cm.ConfigPath != expected {
		t.Fatalf("expected %s, got %s", expected, cm.ConfigPath)
	}
}

func TestGetCredentialUsesProvider(t *testing.T) {
	mp := mockProvider{env: map[string]string{"N8N_URL_DEV": "http://example.com"}}
	cm := NewCredentialManager(mp, "dev")

	val := cm.getCredential(nil, "", "N8N_URL")
	if val != "http://example.com" {
		t.Fatalf("expected provider value, got %s", val)
	}
}
