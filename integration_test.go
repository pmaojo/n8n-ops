package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

var binaryPath string

func buildBinary() (string, func(), error) {
	tmpDir, err := os.MkdirTemp("", "n8n-ops-bin-*")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "n8n-ops")
	cmd := exec.Command("go", "build", "-o", path)
	cmd.Env = os.Environ()
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("build binary: %w: %s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
}

func TestMain(m *testing.M) {
	p, cleanup, err := buildBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	binaryPath = p
	code := m.Run()
	cleanup()
	os.Exit(code)
}

func TestCLIIntegration(t *testing.T) {
	// Test CLI binary exists and runs
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Test help command
	cmd := exec.Command(binaryPath, "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run CLI help: %v", err)
	}

	if len(output) == 0 {
		t.Error("CLI help should produce output")
	}
}

func TestVersionCommand(t *testing.T) {
	// Test version command
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}


	cmd := exec.Command(binaryPath, "version")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run CLI version: %v", err)
	}

	if len(output) == 0 {
		t.Error("Version command should produce output")
	}
}

func TestWelcomeCommand(t *testing.T) {
	// Test welcome command
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}


	cmd := exec.Command(binaryPath, "welcome")
	cmd.Env = append(os.Environ(), "DEMO=true")

	err := cmd.Start()
	if err != nil {
		t.Fatalf("Failed to start welcome command: %v", err)
	}

	// Give it a moment to start
	time.Sleep(100 * time.Millisecond)

	// Kill the process
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
}

func TestDemoMode(t *testing.T) {
	// Test demo mode functionality
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cmd := exec.Command(binaryPath, "sync", "--demo", "--env", "development")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to run sync in demo mode: %v", err)
	}

	if len(output) == 0 {
		t.Error("Demo sync should produce output")
	}
}

func TestEnvironmentHandling(t *testing.T) {
	// Test environment parameter handling
	environments := []string{"development", "staging", "production"}

	for _, env := range environments {
		cmd := exec.Command(binaryPath, "status", "--env", env, "--demo")
		_, err := cmd.Output()
		if err != nil {
			t.Errorf("Failed to run status for environment %s: %v", env, err)
		}
	}
}
