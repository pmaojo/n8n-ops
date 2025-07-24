package testutil

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

const (
	// DefaultServerTimeout is the time to wait for the mock server to become ready.
	DefaultServerTimeout = 5 * time.Second
	// EnvMockServerTimeout can override the startup timeout when running tests.
	EnvMockServerTimeout = "MOCK_SERVER_TIMEOUT"
)

// StartMockServer launches the mock n8n server used in tests and returns a
// function to stop it. It waits until the server is responsive before returning.

// The caller should defer the returned stop function.
// BuildMockServer compiles the mock n8n server binary used in tests.
// It returns a cleanup function to remove the generated binary when tests complete.
func BuildMockServer() (func(), error) {
	dir := filepath.Join("..", "mock-n8n-server")
	cmd := exec.Command("go", "build", "-o", "mock-n8n-server")
	cmd.Dir = dir

	if out, err := cmd.CombinedOutput(); err != nil {
		return nil, fmt.Errorf("build mock server: %w: %s", err, out)
	}

	cleanup := func() {
		os.Remove(filepath.Join(dir, "mock-n8n-server"))
	}
	return cleanup, nil
}

func StartMockServer() (func(), error) {
	cmd := exec.Command("./mock-n8n-server")
	cmd.Dir = filepath.Join("..", "mock-n8n-server")
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start mock server: %w", err)
	}

	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()

	timeout = resolveTimeout(timeout)

	if err := WaitForServer("http://localhost:3001/health", timeout); err != nil {
		cmd.Process.Kill()
		<-done
		return nil, fmt.Errorf("mock server did not start: %w", err)
	}

	stop := func() {
		cmd.Process.Kill()
		<-done
	}
	return stop, nil
}

// WaitForServer polls the provided URL until it returns a successful response
// or the timeout is reached.
func WaitForServer(url string, timeout time.Duration) error {
	if timeout <= 0 {
		timeout = DefaultServerTimeout
	}

	deadline := time.Now().Add(timeout)
	for {
		if time.Now().After(deadline) {
			return fmt.Errorf("timeout waiting for server")
		}
		resp, err := http.Get(url)
		if err == nil && resp.StatusCode == http.StatusOK {
			resp.Body.Close()
			return nil
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(100 * time.Millisecond)
	}
}

// resolveTimeout returns the provided timeout or a value derived from the
// EnvMockServerTimeout environment variable. If neither is set, it falls back
// to DefaultServerTimeout.
func resolveTimeout(provided time.Duration) time.Duration {
	if provided > 0 {
		return provided
	}

	if v := os.Getenv(EnvMockServerTimeout); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}

	return DefaultServerTimeout
}
