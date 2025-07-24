package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

// runCommand executes the CLI with the provided arguments and
// verifies that it exits successfully and prints the expected text.
// Additional environment variables can be supplied via the env slice.
func runCommand(t *testing.T, args []string, expect string, env []string) {
	t.Helper()

	cmd := exec.Command(binaryPath, args...)
	if len(env) > 0 {
		cmd.Env = append(os.Environ(), env...)
	} else {
		cmd.Env = os.Environ()
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("%v failed: %v\n%s", strings.Join(args, " "), err, string(output))
	}
	if !strings.Contains(string(output), expect) {
		t.Errorf("output for %v did not contain %q\n%s", args, expect, string(output))
	}
}

func TestCLIIntegration(t *testing.T) {
	// Test CLI binary exists and runs
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	runCommand(t, []string{"--help"}, "n8n-ops", nil)
}

func TestVersionCommand(t *testing.T) {
	// Test version command
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
	runCommand(t, []string{"version"}, "n8n-ops version", nil)
}

func TestWelcomeCommand(t *testing.T) {
	// Test welcome command
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	cmd := exec.Command(binaryPath, "welcome")
	cmd.Env = append(os.Environ(), "DEMO=true")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Failed to start welcome command: %v", err)
	}
	time.Sleep(100 * time.Millisecond)
	if cmd.Process != nil {
		cmd.Process.Kill()
	}
}

func TestDemoMode(t *testing.T) {
	// Test demo mode functionality
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	runCommand(t, []string{"sync", "--demo", "--env", "development", "--force"}, "Sync", nil)
}

func TestEnvironmentHandling(t *testing.T) {
	// Test environment parameter handling
	environments := []string{"development", "staging", "production"}

	for _, env := range environments {
		runCommand(t, []string{"status", "--env", env, "--demo"}, "Status Dashboard", nil)
	}
}

func TestAllCommandsHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tests := []struct {
		args   []string
		expect string
	}{
		{[]string{"branch", "--help"}, "Usage"},
		{[]string{"check", "--help"}, "Usage"},
		{[]string{"credentials", "--help"}, "Usage"},
		{[]string{"credentials", "list", "--help"}, "Usage"},
		{[]string{"credentials", "validate", "--help"}, "Usage"},
		{[]string{"credentials", "map", "--help"}, "Usage"},
		{[]string{"credentials", "template", "--help"}, "Usage"},
		{[]string{"dashboard", "--help"}, "Usage"},
		{[]string{"deploy", "--help"}, "Usage"},
		{[]string{"init", "--help"}, "Usage"},
		{[]string{"monitor", "--help"}, "Usage"},
		{[]string{"observability", "--help"}, "Usage"},
		{[]string{"observability", "setup", "--help"}, "Usage"},
		{[]string{"observability", "test-connection", "--help"}, "Usage"},
		{[]string{"observability", "create-dashboard", "--help"}, "Usage"},
		{[]string{"onboard", "--help"}, "Usage"},
		{[]string{"quickstart", "--help"}, "Usage"},
		{[]string{"status", "--help"}, "Usage"},
		{[]string{"sync", "--help"}, "Usage"},
		{[]string{"terminal", "--help"}, "Usage"},
		{[]string{"tui", "--help"}, "Usage"},
		{[]string{"tutorial", "--help"}, "Usage"},
		{[]string{"ui", "--help"}, "Usage"},
		{[]string{"validate", "--help"}, "Usage"},
		{[]string{"watch", "--help"}, "Usage"},
		{[]string{"welcome", "--help"}, "Usage"},
	}

	for _, tc := range tests {
		t.Run(strings.Join(tc.args, "-"), func(t *testing.T) {
			runCommand(t, tc.args, tc.expect, []string{"DEMO=true"})
		})
	}
}
