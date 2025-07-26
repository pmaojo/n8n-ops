package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/pmaojo/n8n-ops/internal/workflow"
)

// n8nClient is the concrete implementation of N8nClient interface
type n8nClient struct {
	transport *ContextualTransport
}

// New creates a new n8n client following SOLID principles
// No side effects - connection testing is separate
func New(baseURL, apiKey string, httpClient HTTPDoer) (Client, error) {
	if baseURL == "" {
		return nil, ErrInvalidConfig
	}
	if apiKey == "" {
		return nil, ErrInvalidConfig
	}

	// Clean base URL
	baseURL = strings.TrimRight(baseURL, "/")

	// Create transport with dependency injection
	transport := NewContextualTransport(baseURL, apiKey, httpClient)

	return &n8nClient{
		transport: transport,
	}, nil
}

// Ping tests the connection separately from construction
func (c *n8nClient) Ping(ctx context.Context) error {
	return c.HealthCheck(ctx)
}

// Generic request helper following DRY principle
func doRequest[T any](ctx context.Context, c *n8nClient, method, path string, body interface{}) (T, error) {
	var zero T

	// Prepare request body
	var reqBody []byte
	if body != nil {
		var err error
		reqBody, err = json.Marshal(body)
		if err != nil {
			return zero, fmt.Errorf("failed to marshal request: %w", err)
		}
	}

	// Create request
	req := &Request{
		Method:  method,
		URL:     path,
		Body:    reqBody,
		Headers: make(map[string]string),
	}

	// Execute request
	resp, err := c.transport.DoWithContext(ctx, req)
	if err != nil {
		return zero, fmt.Errorf("request failed: %w", err)
	}

	// Handle HTTP errors
	if resp.StatusCode >= 400 {
		return zero, NewAPIError(resp.StatusCode, string(resp.Body))
	}

	// Deserialize response
	var result T
	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, &result); err != nil {
			return zero, fmt.Errorf("failed to unmarshal response: %w", err)
		}
	}

	return result, nil
}

// HealthCheck implements WorkflowReader interface
func (c *n8nClient) HealthCheck(ctx context.Context) error {
	type healthResponse struct {
		Status string `json:"status"`
	}

	health, err := doRequest[healthResponse](ctx, c, http.MethodGet, "/api/v1/workflows", nil)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	// Basic health validation
	if health.Status == "" {
		return nil // No error, just empty response is OK
	}

	return nil
}

// GetWorkflows implements WorkflowReader interface
func (c *n8nClient) GetWorkflows(ctx context.Context) ([]*workflow.Workflow, error) {
	type workflowsResponse struct {
		Data []*workflow.Workflow `json:"data"`
	}

	resp, err := doRequest[workflowsResponse](ctx, c, http.MethodGet, "/api/v1/workflows", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflows: %w", err)
	}

	return resp.Data, nil
}

// GetWorkflow implements WorkflowReader interface
func (c *n8nClient) GetWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	if id == "" {
		return nil, ErrBadRequest
	}

	path := fmt.Sprintf("/api/v1/workflows/%s", id)
	wf, err := doRequest[*workflow.Workflow](ctx, c, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get workflow %s: %w", id, err)
	}

	return wf, nil
}

// CreateWorkflow implements WorkflowWriter interface
func (c *n8nClient) CreateWorkflow(ctx context.Context, wf *workflow.Workflow) (*workflow.Workflow, error) {
	if wf == nil {
		return nil, ErrBadRequest
	}

	result, err := doRequest[*workflow.Workflow](ctx, c, http.MethodPost, "/api/v1/workflows", wf)
	if err != nil {
		return nil, fmt.Errorf("failed to create workflow: %w", err)
	}

	return result, nil
}

// UpdateWorkflow implements WorkflowWriter interface
func (c *n8nClient) UpdateWorkflow(ctx context.Context, id string, wf *workflow.Workflow) (*workflow.Workflow, error) {
	if id == "" || wf == nil {
		return nil, ErrBadRequest
	}

	path := fmt.Sprintf("/api/v1/workflows/%s", id)
	result, err := doRequest[*workflow.Workflow](ctx, c, http.MethodPut, path, wf)
	if err != nil {
		return nil, fmt.Errorf("failed to update workflow %s: %w", id, err)
	}

	return result, nil
}

// ActivateWorkflow activates a workflow by ID.
func (c *n8nClient) ActivateWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	if id == "" {
		return nil, ErrBadRequest
	}

	path := fmt.Sprintf("/api/v1/workflows/%s/activate", id)
	result, err := doRequest[*workflow.Workflow](ctx, c, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to activate workflow %s: %w", id, err)
	}

	return result, nil
}

// DeactivateWorkflow deactivates a workflow by ID.
func (c *n8nClient) DeactivateWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	if id == "" {
		return nil, ErrBadRequest
	}

	path := fmt.Sprintf("/api/v1/workflows/%s/deactivate", id)
	result, err := doRequest[*workflow.Workflow](ctx, c, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to deactivate workflow %s: %w", id, err)
	}

	return result, nil
}

// DeleteWorkflow implements WorkflowWriter interface
func (c *n8nClient) DeleteWorkflow(ctx context.Context, id string) error {
	if id == "" {
		return ErrBadRequest
	}

	path := fmt.Sprintf("/api/v1/workflows/%s", id)
	_, err := doRequest[interface{}](ctx, c, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete workflow %s: %w", id, err)
	}

	return nil
}

// ExecuteWorkflow implements WorkflowExecutor interface
func (c *n8nClient) ExecuteWorkflow(ctx context.Context, id string) (*workflow.ExecutionResult, error) {
	if id == "" {
		return nil, ErrBadRequest
	}

	path := fmt.Sprintf("/api/v1/workflows/%s/execute", id)
	result, err := doRequest[*workflow.ExecutionResult](ctx, c, http.MethodPost, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to execute workflow %s: %w", id, err)
	}

	return result, nil
}

// GetExecution implements WorkflowExecutor interface
func (c *n8nClient) GetExecution(ctx context.Context, id string) (*workflow.ExecutionResult, error) {
	if id == "" {
		return nil, ErrBadRequest
	}

	path := fmt.Sprintf("/api/v1/executions/%s", id)
	result, err := doRequest[*workflow.ExecutionResult](ctx, c, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get execution %s: %w", id, err)
	}

	return result, nil
}

// GetExecutions implements WorkflowExecutor interface
func (c *n8nClient) GetExecutions(ctx context.Context, workflowID string, status string, limit int) ([]*workflow.ExecutionResult, error) {
	path := fmt.Sprintf("/api/v1/executions?workflowId=%s", workflowID)

	if status != "" {
		path += fmt.Sprintf("&status=%s", status)
	}

	if limit > 0 {
		path += fmt.Sprintf("&limit=%d", limit)
	}

	type executionsResponse struct {
		Data []*workflow.ExecutionResult `json:"data"`
	}

	resp, err := doRequest[executionsResponse](ctx, c, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get executions for workflow %s: %w", workflowID, err)
	}

	return resp.Data, nil
}

// GetCredentials lists all credentials
func (c *n8nClient) GetCredentials(ctx context.Context) ([]*credentials.N8nCredential, error) {
	type credResponse struct {
		Data []*credentials.N8nCredential `json:"data"`
	}

	resp, err := doRequest[credResponse](ctx, c, http.MethodGet, "/rest/credentials", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get credentials: %w", err)
	}
	return resp.Data, nil
}

// GetCredential retrieves a credential by ID
func (c *n8nClient) GetCredential(ctx context.Context, id string) (*credentials.N8nCredential, error) {
	if id == "" {
		return nil, ErrBadRequest
	}
	path := fmt.Sprintf("/rest/credentials/%s", id)
	cred, err := doRequest[*credentials.N8nCredential](ctx, c, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential %s: %w", id, err)
	}
	return cred, nil
}

// CreateCredential creates a new credential
func (c *n8nClient) CreateCredential(ctx context.Context, cred *credentials.N8nCredential) (*credentials.N8nCredential, error) {
	if cred == nil {
		return nil, ErrBadRequest
	}
	created, err := doRequest[*credentials.N8nCredential](ctx, c, http.MethodPost, "/rest/credentials", cred)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}
	return created, nil
}

// UpdateCredential updates an existing credential
func (c *n8nClient) UpdateCredential(ctx context.Context, id string, cred *credentials.N8nCredential) (*credentials.N8nCredential, error) {
	if id == "" || cred == nil {
		return nil, ErrBadRequest
	}
	path := fmt.Sprintf("/rest/credentials/%s", id)
	updated, err := doRequest[*credentials.N8nCredential](ctx, c, http.MethodPatch, path, cred)
	if err != nil {
		return nil, fmt.Errorf("failed to update credential %s: %w", id, err)
	}
	return updated, nil
}

// DeleteCredential removes a credential by ID
func (c *n8nClient) DeleteCredential(ctx context.Context, id string) error {
	if id == "" {
		return ErrBadRequest
	}
	path := fmt.Sprintf("/rest/credentials/%s", id)
	_, err := doRequest[interface{}](ctx, c, http.MethodDelete, path, nil)
	if err != nil {
		return fmt.Errorf("failed to delete credential %s: %w", id, err)
	}
	return nil
}

// GetCredentialSchema retrieves the JSON schema for a credential type
func (c *n8nClient) GetCredentialSchema(ctx context.Context, typ string) (map[string]interface{}, error) {
	if typ == "" {
		return nil, ErrBadRequest
	}
	path := fmt.Sprintf("/rest/credentials/schema/%s", typ)
	schema, err := doRequest[map[string]interface{}](ctx, c, http.MethodGet, path, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get credential schema %s: %w", typ, err)
	}
	return schema, nil
}
