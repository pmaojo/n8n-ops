package testutil

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestWaitForServerSuccess(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	if err := WaitForServer(srv.URL, 500*time.Millisecond); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestWaitForServerTimeout(t *testing.T) {
	start := time.Now()
	err := WaitForServer("http://localhost:65535", 200*time.Millisecond)
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if time.Since(start) < 200*time.Millisecond {
		t.Error("did not wait for the timeout duration")
	}
}

func TestResolveTimeout(t *testing.T) {
	os.Setenv(EnvMockServerTimeout, "150ms")
	if d := resolveTimeout(0); d != 150*time.Millisecond {
		t.Errorf("expected 150ms, got %v", d)
	}
	os.Unsetenv(EnvMockServerTimeout)

	if d := resolveTimeout(100 * time.Millisecond); d != 100*time.Millisecond {
		t.Errorf("expected provided duration to be used")
	}
}
