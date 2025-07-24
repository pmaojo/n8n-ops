package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/pmaojo/n8n-ops/internal/workflow"
)

// mockDoer implements HTTPDoer for deterministic tests.
type mockDoer struct {
	roundTrip func(*http.Request) (*http.Response, error)
	lastReq   *http.Request
}

func (m *mockDoer) Do(req *http.Request) (*http.Response, error) {
	m.lastReq = req
	return m.roundTrip(req)
}

// newMockClient returns a client using the provided RoundTripper.
func newMockClient(t testing.TB, doer HTTPDoer) *n8nClient {
	t.Helper()
	c, err := New("http://api.example.com", "token", doer)
	if err != nil {
		t.Fatalf("new client: %v", err)
	}
	return c.(*n8nClient)
}

func TestNewInvalidConfig(t *testing.T) {
	if _, err := New("", "key", nil); !errors.Is(err, ErrInvalidConfig) {
		t.Errorf("expected invalid config error, got %v", err)
	}
	if _, err := New("http://x", "", nil); !errors.Is(err, ErrInvalidConfig) {
		t.Errorf("expected invalid config error, got %v", err)
	}
	if c, err := New("http://x", "k", nil); err != nil || c == nil {
		t.Errorf("unexpected error creating client: %v", err)
	}
}

func TestDoRequestSuccess(t *testing.T) {
	rt := &mockDoer{roundTrip: func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", req.Method)
		}
		if req.URL.String() != "http://api.example.com/test" {
			t.Errorf("unexpected URL %s", req.URL)
		}
		body := `{"message":"ok"}`
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		}, nil
	}}

	c := newMockClient(t, rt)

	type resp struct {
		Message string `json:"message"`
	}
	r, err := doRequest[resp](context.Background(), c, http.MethodGet, "/test", nil)
	if err != nil {
		t.Fatalf("doRequest error: %v", err)
	}
	if r.Message != "ok" {
		t.Errorf("unexpected response: %v", r)
	}
}

func TestDoRequestHTTPError(t *testing.T) {
	rt := &mockDoer{roundTrip: func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
			Body:       io.NopCloser(strings.NewReader("fail")),
		}, nil
	}}

	c := newMockClient(t, rt)
	_, err := doRequest[struct{}](context.Background(), c, http.MethodGet, "/fail", nil)
	if err == nil || !IsServerError(err) {
		t.Fatalf("expected server error, got %v", err)
	}
}

func TestDoRequestTransportError(t *testing.T) {
	rt := &mockDoer{roundTrip: func(req *http.Request) (*http.Response, error) {
		return nil, errors.New("boom")
	}}
	c := newMockClient(t, rt)
	_, err := doRequest[struct{}](context.Background(), c, http.MethodGet, "/boom", nil)
	if err == nil || !strings.Contains(err.Error(), "boom") {
		t.Fatalf("expected transport error, got %v", err)
	}
}

func TestGetWorkflowSuccess(t *testing.T) {
	rt := &mockDoer{roundTrip: func(req *http.Request) (*http.Response, error) {
		if req.URL.Path != "/api/v1/workflows/42" {
			t.Errorf("unexpected path %s", req.URL.Path)
		}
		wf := &workflow.Workflow{
			ID:    "42",
			Name:  "demo",
			Nodes: []workflow.Node{{Name: "start", Type: "start", Position: []float64{0, 0}}},
		}
		b, _ := json.Marshal(wf)
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(b))}, nil
	}}

	c := newMockClient(t, rt)
	wf, err := c.GetWorkflow(context.Background(), "42")
	if err != nil {
		t.Fatalf("GetWorkflow error: %v", err)
	}
	if wf.ID != "42" {
		t.Errorf("unexpected workflow ID %s", wf.ID)
	}
}

func TestGetWorkflowNotFound(t *testing.T) {
	rt := &mockDoer{roundTrip: func(req *http.Request) (*http.Response, error) {
		return &http.Response{StatusCode: http.StatusNotFound, Body: io.NopCloser(strings.NewReader("missing"))}, nil
	}}
	c := newMockClient(t, rt)
	if _, err := c.GetWorkflow(context.Background(), "99"); err == nil || !IsNotFound(err) {
		t.Fatalf("expected not found error, got %v", err)
	}
}

func TestGetWorkflowBadRequest(t *testing.T) {
	c := newMockClient(t, &mockDoer{roundTrip: func(req *http.Request) (*http.Response, error) {
		t.Fatal("http call should not occur")
		return nil, nil
	}})
	if _, err := c.GetWorkflow(context.Background(), ""); !errors.Is(err, ErrBadRequest) {
		t.Fatalf("expected bad request error, got %v", err)
	}
}

func TestCreateWorkflowMock(t *testing.T) {
	rt := &mockDoer{roundTrip: func(req *http.Request) (*http.Response, error) {
		if req.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", req.Method)
		}
		var wf workflow.Workflow
		if err := json.NewDecoder(req.Body).Decode(&wf); err != nil {
			t.Errorf("decode body: %v", err)
		}
		wf.ID = "7"
		b, _ := json.Marshal(&wf)
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(bytes.NewReader(b))}, nil
	}}
	c := newMockClient(t, rt)
	wf, err := c.CreateWorkflow(context.Background(), &workflow.Workflow{Name: "new", Nodes: []workflow.Node{{Name: "start", Type: "start", Position: []float64{0, 0}}}})
	if err != nil {
		t.Fatalf("CreateWorkflow error: %v", err)
	}
	if wf.ID != "7" {
		t.Errorf("unexpected id %s", wf.ID)
	}

	if _, err := c.CreateWorkflow(context.Background(), nil); !errors.Is(err, ErrBadRequest) {
		t.Errorf("expected bad request error, got %v", err)
	}
}
