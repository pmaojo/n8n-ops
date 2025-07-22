package client

import (
	"fmt"
	"time"
)

type DemoN8nClient struct {
	workflows map[string]*Workflow
}

func NewDemoN8nClient() *DemoN8nClient {
	workflows := map[string]*Workflow{
		"1001": {
			ID:        "1001",
			Name:      "Customer Onboarding",
			Active:    true,
			CreatedAt: time.Now().Add(-72 * time.Hour),
			UpdatedAt: time.Now().Add(-30 * time.Minute),
			VersionId: 15,
			Tags:      []string{"customer", "onboarding"},
			Nodes: []WorkflowNode{
				{
					ID:         "webhook-1001",
					Name:       "Webhook",
					Type:       "n8n-nodes-base.webhook",
					TypeVersion: 1.0,
					Position:   []int{250, 300},
					Parameters: map[string]interface{}{
						"httpMethod": "POST",
						"path":       "customer-signup",
					},
				},
				{
					ID:         "email-1001",
					Name:       "Send Welcome Email",
					Type:       "n8n-nodes-base.emailSend",
					TypeVersion: 1.0,
					Position:   []int{450, 300},
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
								"node": "Send Welcome Email",
								"type": "main",
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
			Tags:      []string{"payment", "stripe"},
			Nodes: []WorkflowNode{
				{
					ID:         "webhook-1002",
					Name:       "Payment Webhook",
					Type:       "n8n-nodes-base.webhook",
					TypeVersion: 1.0,
					Position:   []int{250, 200},
					Parameters: map[string]interface{}{
						"httpMethod": "POST",
						"path":       "stripe-webhook",
					},
				},
				{
					ID:         "stripe-1002",
					Name:       "Process Payment",
					Type:       "n8n-nodes-base.stripe",
					TypeVersion: 1.0,
					Position:   []int{450, 200},
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
								"node": "Process Payment",
								"type": "main",
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
			Tags:      []string{"orders", "fulfillment"},
			Nodes: []WorkflowNode{
				{
					ID:         "trigger-1003",
					Name:       "Order Created",
					Type:       "n8n-nodes-base.httpRequest",
					TypeVersion: 1.0,
					Position:   []int{250, 400},
					Parameters: map[string]interface{}{
						"method": "POST",
					},
				},
			},
			Connections: map[string]interface{}{},
		},
	}

	return &DemoN8nClient{workflows: workflows}
}

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

func (c *DemoN8nClient) GetWorkflows() ([]Workflow, error) {
	var workflows []Workflow
	for _, workflow := range c.workflows {
		workflows = append(workflows, *workflow)
	}
	return workflows, nil
}

func (c *DemoN8nClient) GetWorkflow(id string) (*Workflow, error) {
	workflow, exists := c.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}
	return workflow, nil
}

func (c *DemoN8nClient) CreateWorkflow(workflow *Workflow) (*Workflow, error) {
	newID := fmt.Sprintf("demo-%d", len(c.workflows)+1001)
	workflow.ID = newID
	workflow.CreatedAt = time.Now()
	workflow.UpdatedAt = time.Now()
	workflow.VersionId = 1
	
	c.workflows[newID] = workflow
	return workflow, nil
}

func (c *DemoN8nClient) UpdateWorkflow(id string, workflow *Workflow) (*Workflow, error) {
	existing, exists := c.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}
	
	workflow.ID = id
	workflow.CreatedAt = existing.CreatedAt
	workflow.UpdatedAt = time.Now()
	workflow.VersionId = existing.VersionId + 1
	
	c.workflows[id] = workflow
	return workflow, nil
}

func (c *DemoN8nClient) ActivateWorkflow(id string) (*Workflow, error) {
	workflow, exists := c.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}
	
	workflow.Active = true
	workflow.UpdatedAt = time.Now()
	return workflow, nil
}

func (c *DemoN8nClient) DeactivateWorkflow(id string) (*Workflow, error) {
	workflow, exists := c.workflows[id]
	if !exists {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}
	
	workflow.Active = false
	workflow.UpdatedAt = time.Now()
	return workflow, nil
}

func (c *DemoN8nClient) DeleteWorkflow(id string) error {
	if _, exists := c.workflows[id]; !exists {
		return fmt.Errorf("workflow not found: %s", id)
	}
	
	delete(c.workflows, id)
	return nil
}

func (c *DemoN8nClient) TestConnection() error {
	// Demo client always passes connection test
	return nil
}