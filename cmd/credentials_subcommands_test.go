package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func captureOutput(t *testing.T, fn func()) string {
	t.Helper()
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	io.Copy(&buf, r)
	r.Close()
	return buf.String()
}

func TestMaskValue(t *testing.T) {
	if v := maskValue("secret"); v != "se**et" {
		t.Fatalf("unexpected mask %s", v)
	}
	if v := maskValue(""); v != "not_set" {
		t.Fatalf("empty mask %s", v)
	}
}

func TestGetEnvOrDefault(t *testing.T) {
	os.Setenv("FOO", "bar")
	defer os.Unsetenv("FOO")
	if v := getEnvOrDefault("FOO", "x"); v != "bar" {
		t.Fatalf("expected env value, got %s", v)
	}
	if v := getEnvOrDefault("MISSING", "def"); v != "def" {
		t.Fatalf("expected default, got %s", v)
	}
}

func TestRunCredentialsSubcommands(t *testing.T) {
	environment = "development"
	showValues = false
	out := captureOutput(t, func() { _ = runCredentialsList(nil, nil) })
	if !bytes.Contains([]byte(out), []byte("Credential Mappings")) {
		t.Fatalf("unexpected output: %s", out)
	}

	out = captureOutput(t, func() { _ = runCredentialsValidate(nil, nil) })
	if !bytes.Contains([]byte(out), []byte("Missing")) {
		t.Fatalf("validate missing message not found")
	}

	out = captureOutput(t, func() { _ = runCredentialsMap(nil, nil) })
	if !bytes.Contains([]byte(out), []byte("Environment Variable Patterns")) {
		t.Fatalf("map output unexpected: %s", out)
	}

	out = captureOutput(t, func() { _ = runCredentialsTemplate(nil, nil) })
	if !bytes.Contains([]byte(out), []byte("N8N_API_KEY_DEVELOPMENT")) {
		t.Fatalf("template output missing var")
	}
}
