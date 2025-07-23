package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/pmaojo/n8n-ops/internal/workflow"
)

// MockTransport implements HTTPTransport for testing
type MockTransport struct {
	responses map[string]*Response
	requests  []*Request
}

func NewMockTransport() *MockTransport {
	return &MockTransport{
		responses: make(map[string]*Response),
		requests:  make([]*Request, 0),
	}
}

func (m *MockTransport) Do(req *Request) (*Response, error) {
	m.requests = append(m.requests, req)

	key := req.Method + " " + req.URL
	if resp, exists := m.responses[key]; exists {
		return resp, nil
	}

	// Default response
	return &Response{
		StatusCode: http.StatusOK,
		Body:       []byte(`{}`),
		Headers:    make(map[string]string),
	}, nil
}

func (m *MockTransport) AddResponse(method, url string, statusCode int, body string) {
	key := method + " " + url
	m.responses[key] = &Response{
		StatusCode: statusCode,
		Body:       []byte(body),
		Headers:    make(map[string]string),
	}
}

func (m *MockTransport) GetRequests() []*Request {
	return m.requests
}

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		apiKey  string
		wantErr bool
	}{
		{
			name:    "valid parameters",
			baseURL: "https://api.example.com",
			apiKey:  "test-key",
			wantErr: false,
		},
		{
			name:    "empty base URL",
			baseURL: "",
			apiKey:  "test-key",
			wantErr: true,
		},
		{
			name:    "empty API key",
			baseURL: "https://api.example.com",
			apiKey:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := New(tt.baseURL, tt.apiKey, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Errorf("New() returned nil client")
			}
		})
	}
}

func TestHealthCheck(t *testing.T) {
	ctx := context.Background()

	// Create client with short timeout for testing
	client, err := New("https://test.example.com", "test-key", &http.Client{Timeout: 1 * time.Millisecond})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// This will likely timeout, but tests the method signature
	err = client.HealthCheck(ctx)

	// We expect an error due to timeout/invalid URL, but method should exist
	if err == nil {
		t.Errorf("Expected error due to invalid URL, but got no error")
	}
}

func TestGetWorkflows(t *testing.T) {
	ctx := context.Background()

	// Create a minimal client for testing structure
	client, err := New("https://test.example.com", "test-key", &http.Client{Timeout: 1 * time.Millisecond})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// This will likely timeout, but tests the method signature
	workflows, err := client.GetWorkflows(ctx)

	// We expect an error due to timeout/invalid URL, but method should exist
	if workflows != nil && err == nil {
		t.Errorf("Expected error due to invalid URL, but got workflows: %v", workflows)
	}

	// Test that the method signature is correct
	var _ []*workflow.Workflow = workflows // This should compile
}

func TestCreateWorkflow(t *testing.T) {
	ctx := context.Background()

	client, err := New("https://test.example.com", "test-key", &http.Client{Timeout: 1 * time.Millisecond})
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	wf := &workflow.Workflow{
		Name:   "Test Workflow",
		Active: false,
		Nodes:  []workflow.Node{},
	}

	// Test nil workflow
	result, err := client.CreateWorkflow(ctx, nil)
	if err == nil {
		t.Errorf("Expected error for nil workflow")
	}
	if result != nil {
		t.Errorf("Expected nil result for nil workflow")
	}

	// Test with valid workflow (will timeout, but tests signature)
	result, err = client.CreateWorkflow(ctx, wf)
	if err == nil {
		t.Errorf("expected error due to timeout")
	}

	// We expect an error due to timeout, but method should exist
	var _ *workflow.Workflow = result // This should compile
}

func TestErrorTypes(t *testing.T) {
	// Test that our error types implement error interface
	var _ error = ErrNotFound
	var _ error = ErrUnauthorized
	var _ error = ErrBadRequest
	var _ error = ErrServerError

	// Test APIError
	apiErr := NewAPIError(404, "not found")
	if apiErr.Code != 404 {
		t.Errorf("Expected code 404, got %d", apiErr.Code)
	}

	errStr := apiErr.Error()
	if errStr == "" {
		t.Errorf("APIError.Error() returned empty string")
	}

	// Test error checking functions
	if !IsAPIError(apiErr, 404) {
		t.Errorf("IsAPIError should return true for matching status code")
	}

	if IsAPIError(apiErr, 500) {
		t.Errorf("IsAPIError should return false for non-matching status code")
	}
}

// BenchmarkHealthCheck tests performance of health check
func BenchmarkHealthCheck(b *testing.B) {
	client, err := New("https://httpbin.org", "test-key", nil)
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}

	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client.HealthCheck(ctx)
	}
}
