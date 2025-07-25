package client

import (
	"context"
	"fmt"
	"time"

	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/pmaojo/n8n-ops/internal/workflow"
)

// DemoN8nClient is an in-memory implementation of the Client interface used
// for demos and tests. It avoids any network calls and stores workflows in a
// map keyed by workflow ID.
type DemoN8nClient struct {
	workflows   map[string]*workflow.Workflow
	credentials map[string]*credentials.N8nCredential
	credCounter int
}

// NewDemoN8nClient creates a DemoN8nClient pre-populated with example
// workflows. It is safe for concurrent use in tests since it operates only on
// local memory.
func NewDemoN8nClient() *DemoN8nClient {
	workflows := map[string]*workflow.Workflow{
		"1001": {
			ID:        "1001",
			Name:      "Customer Onboarding",
			Active:    true,
			CreatedAt: time.Now().Add(-72 * time.Hour),
			UpdatedAt: time.Now().Add(-30 * time.Minute),
			VersionId: 15,
			Tags: []workflow.Tag{
				{Name: "customer"},
				{Name: "onboarding"},
			},
			Nodes: []workflow.Node{
				{
					ID:          "webhook-1001",
					Name:        "Webhook",
					Type:        "n8n-nodes-base.webhook",
					TypeVersion: 1.0,
					Position:    []float64{250, 300},
					Parameters: map[string]interface{}{
						"httpMethod": "POST",
						"path":       "customer-signup",
					},
				},
				{
					ID:          "email-1001",
					Name:        "Send Welcome Email",
					Type:        "n8n-nodes-base.emailSend",
					TypeVersion: 1.0,
					Position:    []float64{450, 300},
					Parameters: map[string]interface{}{
						"subject": "Welcome to our platform!",
						"toEmail": "={{ $json.email }}",
					},
				},
			},
			Connections: map[string]interface{}{
				"Webhook": map[string]interface{}{
					"main": [][]map[string]interface{}{
						{
							{
								"node":  "Send Welcome Email",
								"type":  "main",
								"index": 0,
							},
						},
					},
				},
			},
		},
		"1002": {
			ID:        "1002",
			Name:      "Payment Processing",
			Active:    true,
			CreatedAt: time.Now().Add(-48 * time.Hour),
			UpdatedAt: time.Now().Add(-10 * time.Minute),
			VersionId: 23,
			Tags: []workflow.Tag{
				{Name: "payment"},
				{Name: "stripe"},
			},
			Nodes: []workflow.Node{
				{
					ID:          "webhook-1002",
					Name:        "Payment Webhook",
					Type:        "n8n-nodes-base.webhook",
					TypeVersion: 1.0,
					Position:    []float64{250, 200},
					Parameters: map[string]interface{}{
						"httpMethod": "POST",
						"path":       "stripe-webhook",
					},
				},
				{
					ID:          "stripe-1002",
					Name:        "Process Payment",
					Type:        "n8n-nodes-base.stripe",
					TypeVersion: 1.0,
					Position:    []float64{450, 200},
					Parameters: map[string]interface{}{
						"operation": "charge",
					},
				},
			},
			Connections: map[string]interface{}{
				"Payment Webhook": map[string]interface{}{
					"main": [][]map[string]interface{}{
						{
							{
								"node":  "Process Payment",
								"type":  "main",
								"index": 0,
							},
						},
					},
				},
			},
		},
		"1003": {
			ID:        "1003",
			Name:      "Order Fulfillment",
			Active:    false,
			CreatedAt: time.Now().Add(-24 * time.Hour),
			UpdatedAt: time.Now().Add(-5 * time.Minute),
			VersionId: 8,
			Tags: []workflow.Tag{
				{Name: "orders"},
				{Name: "fulfillment"},
			},
			Nodes: []workflow.Node{
				{
					ID:          "trigger-1003",
					Name:        "Order Created",
					Type:        "n8n-nodes-base.httpRequest",
					TypeVersion: 1.0,
					Position:    []float64{250, 400},
					Parameters: map[string]interface{}{
						"method": "POST",
					},
				},
			},
			Connections: map[string]interface{}{},
		},
	}

	return &DemoN8nClient{workflows: workflows, credentials: make(map[string]*credentials.N8nCredential)}
}

// GetMe returns static user information representing the currently authenticated
// user. It mirrors the behavior of the real n8n API without performing any
// network requests.
func (c *DemoN8nClient) GetMe() (*User, error) {
	return &User{
		ID:    "demo-user-123",
		Email: "demo@n8n-ops.example.com",
		Settings: map[string]interface{}{
			"timezone": "UTC",
			"demoMode": true,
		},
	}, nil
}

// HealthCheck simulates an API health check and always succeeds.
func (c *DemoN8nClient) HealthCheck(ctx context.Context) error {
	return nil
}

// GetWorkflows returns a copy of all in-memory workflows.
func (c *DemoN8nClient) GetWorkflows(ctx context.Context) ([]*workflow.Workflow, error) {
	var workflows []*workflow.Workflow
	for _, wf := range c.workflows {
		copy := *wf
		workflows = append(workflows, &copy)
	}
	return workflows, nil
}

// GetWorkflow retrieves a single workflow by ID.
func (c *DemoN8nClient) GetWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	wf, exists := c.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}
	copy := *wf
	return &copy, nil
}

// CreateWorkflow stores a new workflow and returns a copy with generated fields.
func (c *DemoN8nClient) CreateWorkflow(ctx context.Context, wf *workflow.Workflow) (*workflow.Workflow, error) {
	newID := fmt.Sprintf("demo-%d", len(c.workflows)+1001)
	wf.ID = newID
	wf.CreatedAt = time.Now()
	wf.UpdatedAt = time.Now()
	wf.VersionId = 1

	c.workflows[newID] = wf
	copy := *wf
	return &copy, nil
}

// UpdateWorkflow replaces an existing workflow with the provided one.
func (c *DemoN8nClient) UpdateWorkflow(ctx context.Context, id string, wf *workflow.Workflow) (*workflow.Workflow, error) {
	existing, exists := c.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}

	wf.ID = id
	wf.CreatedAt = existing.CreatedAt
	wf.UpdatedAt = time.Now()
	wf.VersionId = existing.VersionId + 1

	c.workflows[id] = wf
	copy := *wf
	return &copy, nil
}

// ActivateWorkflow marks a workflow as active.
func (c *DemoN8nClient) ActivateWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	wf, exists := c.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}

	wf.Active = true
	wf.UpdatedAt = time.Now()
	copy := *wf
	return &copy, nil
}

// DeactivateWorkflow marks a workflow as inactive.
func (c *DemoN8nClient) DeactivateWorkflow(ctx context.Context, id string) (*workflow.Workflow, error) {
	wf, exists := c.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}

	wf.Active = false
	wf.UpdatedAt = time.Now()
	copy := *wf
	return &copy, nil
}

// DeleteWorkflow removes a workflow from the in-memory store.
func (c *DemoN8nClient) DeleteWorkflow(ctx context.Context, id string) error {
	if _, exists := c.workflows[id]; !exists {
		return fmt.Errorf("workflow not found: %s", id)
	}

	delete(c.workflows, id)
	return nil
}

// TestConnection is a no-op that always succeeds for the demo client.
func (c *DemoN8nClient) TestConnection() error {
	// Demo client always passes connection test
	return nil
}

// ExecuteWorkflow is not implemented for the demo client.
func (c *DemoN8nClient) ExecuteWorkflow(ctx context.Context, id string) (*workflow.ExecutionResult, error) {
	return nil, fmt.Errorf("ExecuteWorkflow not implemented in demo client")
}

// GetExecution is not implemented for the demo client.
func (c *DemoN8nClient) GetExecution(ctx context.Context, id string) (*workflow.ExecutionResult, error) {
	return nil, fmt.Errorf("GetExecution not implemented in demo client")
}

// GetExecutions is not implemented for the demo client.
func (c *DemoN8nClient) GetExecutions(ctx context.Context, workflowID string, status string, limit int) ([]*workflow.ExecutionResult, error) {
	return nil, fmt.Errorf("GetExecutions not implemented in demo client")
}

// GetCredentials returns stored credentials
func (c *DemoN8nClient) GetCredentials(ctx context.Context) ([]*credentials.N8nCredential, error) {
	var list []*credentials.N8nCredential
	for _, cred := range c.credentials {
		copy := *cred
		list = append(list, &copy)
	}
	return list, nil
}

// GetCredential retrieves a credential by ID
func (c *DemoN8nClient) GetCredential(ctx context.Context, id string) (*credentials.N8nCredential, error) {
	cred, ok := c.credentials[id]
	if !ok {
		return nil, fmt.Errorf("credential not found: %s", id)
	}
	copy := *cred
	return &copy, nil
}

// CreateCredential stores a new credential
func (c *DemoN8nClient) CreateCredential(ctx context.Context, cred *credentials.N8nCredential) (*credentials.N8nCredential, error) {
	c.credCounter++
	if cred.ID == "" {
		cred.ID = fmt.Sprintf("cred-%d", c.credCounter)
	}
	c.credentials[cred.ID] = cred
	copy := *cred
	return &copy, nil
}

// UpdateCredential replaces a credential
func (c *DemoN8nClient) UpdateCredential(ctx context.Context, id string, cred *credentials.N8nCredential) (*credentials.N8nCredential, error) {
	if cred == nil {
		return nil, fmt.Errorf("nil credential")
	}
	cred.ID = id
	c.credentials[id] = cred
	copy := *cred
	return &copy, nil
}

// DeleteCredential removes a credential
func (c *DemoN8nClient) DeleteCredential(ctx context.Context, id string) error {
	delete(c.credentials, id)
	return nil
}

// GetCredentialSchema returns a placeholder schema
func (c *DemoN8nClient) GetCredentialSchema(ctx context.Context, typ string) (map[string]interface{}, error) {
	return map[string]interface{}{"name": typ}, nil
}
