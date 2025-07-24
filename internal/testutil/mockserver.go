package testutil

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
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

// SetupMockServer builds the mock server binary and starts it. It returns
// functions to stop the running server and clean up the binary.
func SetupMockServer() (stop func(), cleanup func(), err error) {
	cleanup, err = BuildMockServer()
	if err != nil {
		return nil, nil, err
	}
	stop, err = StartMockServer()
	if err != nil {
		if cleanup != nil {
			cleanup()
		}
		return nil, nil, err
	}
	return stop, cleanup, nil
}

// StartMockServer launches the mock n8n server binary built by BuildMockServer.
// It returns a stop function that should be deferred by the caller.
func StartMockServer() (func(), error) {
	var logBuf bytes.Buffer
	cmd := exec.Command("./mock-n8n-server")
	cmd.Dir = filepath.Join("..", "mock-n8n-server")
	cmd.Stdout = &logBuf
	cmd.Stderr = &logBuf

	if err := cmd.Start(); err != nil {
		fmt.Fprint(os.Stderr, logBuf.String())
		return nil, fmt.Errorf("failed to start mock server: %w", err)
	}

	done := make(chan struct{})
	go func() {
		_ = cmd.Wait()
		close(done)
	}()

	if err := WaitForServer("http://localhost:3001/health", 0); err != nil {
		cmd.Process.Kill()
		<-done
		fmt.Fprint(os.Stderr, logBuf.String())
		return nil, fmt.Errorf("mock server did not start: %w", err)
	}

	stop := func() {
		cmd.Process.Kill()
		<-done
		if testing.Verbose() {
			fmt.Fprint(os.Stdout, logBuf.String())
		}
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
