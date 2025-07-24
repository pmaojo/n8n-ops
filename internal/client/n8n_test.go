package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/workflow"
)

// newTestClient creates a client backed by the provided handler.
func newTestClient(t testing.TB, handler http.HandlerFunc) (Client, func()) {
	srv := httptest.NewServer(handler)
	c, err := New(srv.URL, "token", &DefaultHTTPDoer{srv.Client()})
	if err != nil {
		srv.Close()
		t.Fatalf("failed to create client: %v", err)
	}
	return c, srv.Close
}

func TestNewValidation(t *testing.T) {
	if _, err := New("", "k", nil); err == nil {
		t.Error("expected error for empty baseURL")
	}
	if _, err := New("http://example.com", "", nil); err == nil {
		t.Error("expected error for empty apiKey")
	}
	if _, err := New("http://example.com", "k", nil); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestHealthCheck(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/workflows" {
			http.NotFound(w, r)
			return
		}
		if r.Header.Get("X-N8N-API-KEY") != "token" {
			t.Errorf("missing api key header")
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}

	client, closeFn := newTestClient(t, handler)
	defer closeFn()

	if err := client.HealthCheck(context.Background()); err != nil {
		t.Fatalf("health check failed: %v", err)
	}
}

func TestGetWorkflows(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("unexpected method: %s", r.Method)
		}
		json.NewEncoder(w).Encode(map[string]any{
			"data": []*workflow.Workflow{{ID: "1", Name: "test", Nodes: []workflow.Node{{Name: "start", Type: "start", Position: []float64{0, 0}}}}},
		})
	}
	client, closeFn := newTestClient(t, handler)
	defer closeFn()

	wfs, err := client.GetWorkflows(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(wfs) != 1 || wfs[0].ID != "1" {
		t.Fatalf("unexpected workflows: %+v", wfs)
	}
}

func TestCreateWorkflow(t *testing.T) {
	called := false
	handler := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("unexpected method: %s", r.Method)
		}
		called = true
		var wf workflow.Workflow
		if err := json.NewDecoder(r.Body).Decode(&wf); err != nil {
			t.Errorf("invalid body: %v", err)
		}
		wf.ID = "42"
		json.NewEncoder(w).Encode(&wf)
	}
	client, closeFn := newTestClient(t, handler)
	defer closeFn()

	res, err := client.CreateWorkflow(context.Background(), &workflow.Workflow{Name: "demo", Nodes: []workflow.Node{{Name: "start", Type: "start", Position: []float64{0, 0}}}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("handler not called")
	}
	if res.ID != "42" {
		t.Fatalf("unexpected workflow ID: %s", res.ID)
	}

	if _, err := client.CreateWorkflow(context.Background(), nil); err == nil {
		t.Error("expected error for nil workflow")
	}
}

func TestErrorTypes(t *testing.T) {
	apiErr := NewAPIError(http.StatusNotFound, "missing")
	if !IsAPIError(apiErr, http.StatusNotFound) {
		t.Error("IsAPIError should detect not found")
	}
	if IsAPIError(apiErr, http.StatusOK) {
		t.Error("IsAPIError should not match other codes")
	}
	if apiErr.Error() == "" {
		t.Error("error string should not be empty")
	}
}

func BenchmarkHealthCheck(b *testing.B) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
	c, closeFn := newTestClient(b, handler)
	defer closeFn()
	ctx := context.Background()
	for i := 0; i < b.N; i++ {
		c.HealthCheck(ctx)
	}
}
