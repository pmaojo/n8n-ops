package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/n8n-workflows/cli/internal/workflow"
)

type N8nClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

type N8nResponse struct {
	Data interface{} `json:"data"`
}

type N8nWorkflowResponse struct {
	Data []*workflow.Workflow `json:"data"`
}

// NewN8nClient creates a new n8n API client
func NewN8nClient(baseURL, apiKey string) (*N8nClient, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("base URL cannot be empty")
	}
	if apiKey == "" {
		return nil, fmt.Errorf("API key cannot be empty")
	}

	client := &N8nClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Test connection
	if err := client.testConnection(); err != nil {
		return nil, fmt.Errorf("failed to connect to n8n instance: %w", err)
	}

	return client, nil
}

// testConnection verifies the connection to n8n instance
func (c *N8nClient) testConnection() error {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v1/workflows", nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authentication failed - check API key")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// GetWorkflows retrieves all workflows from n8n instance
func (c *N8nClient) GetWorkflows() ([]*workflow.Workflow, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v1/workflows", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workflows: HTTP %d", resp.StatusCode)
	}

	var response N8nWorkflowResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

// GetWorkflow retrieves a specific workflow by ID
func (c *N8nClient) GetWorkflow(id string) (*workflow.Workflow, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v1/workflows/"+id, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("workflow not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workflow: HTTP %d", resp.StatusCode)
	}

	var wf workflow.Workflow
	if err := json.NewDecoder(resp.Body).Decode(&wf); err != nil {
		return nil, fmt.Errorf("failed to decode workflow: %w", err)
	}

	return &wf, nil
}

// FindWorkflowByName finds a workflow by name and returns its ID
func (c *N8nClient) FindWorkflowByName(name string) (string, error) {
	workflows, err := c.GetWorkflows()
	if err != nil {
		return "", err
	}

	for _, wf := range workflows {
		if wf.Name == name {
			return wf.ID, nil
		}
	}

	return "", nil // Not found
}

// CreateWorkflow creates a new workflow
func (c *N8nClient) CreateWorkflow(wf *workflow.Workflow) (*workflow.Workflow, error) {
	// Remove ID for creation
	wf.ID = ""
	
	data, err := json.Marshal(wf)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow: %w", err)
	}

	req, err := http.NewRequest("POST", c.baseURL+"/api/v1/workflows", bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create workflow: HTTP %d - %s", resp.StatusCode, string(body))
	}

	var createdWorkflow workflow.Workflow
	if err := json.NewDecoder(resp.Body).Decode(&createdWorkflow); err != nil {
		return nil, fmt.Errorf("failed to decode created workflow: %w", err)
	}

	return &createdWorkflow, nil
}

// UpdateWorkflow updates an existing workflow
func (c *N8nClient) UpdateWorkflow(id string, wf *workflow.Workflow) (*workflow.Workflow, error) {
	// Ensure ID matches
	wf.ID = id
	
	data, err := json.Marshal(wf)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal workflow: %w", err)
	}

	req, err := http.NewRequest("PUT", c.baseURL+"/api/v1/workflows/"+id, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to update workflow: HTTP %d - %s", resp.StatusCode, string(body))
	}

	var updatedWorkflow workflow.Workflow
	if err := json.NewDecoder(resp.Body).Decode(&updatedWorkflow); err != nil {
		return nil, fmt.Errorf("failed to decode updated workflow: %w", err)
	}

	return &updatedWorkflow, nil
}

// DeleteWorkflow deletes a workflow
func (c *N8nClient) DeleteWorkflow(id string) error {
	req, err := http.NewRequest("DELETE", c.baseURL+"/api/v1/workflows/"+id, nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete workflow: HTTP %d", resp.StatusCode)
	}

	return nil
}

// ActivateWorkflow activates a workflow
func (c *N8nClient) ActivateWorkflow(id string) error {
	req, err := http.NewRequest("POST", c.baseURL+"/api/v1/workflows/"+id+"/activate", nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to activate workflow: HTTP %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeactivateWorkflow deactivates a workflow
func (c *N8nClient) DeactivateWorkflow(id string) error {
	req, err := http.NewRequest("POST", c.baseURL+"/api/v1/workflows/"+id+"/deactivate", nil)
	if err != nil {
		return err
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to deactivate workflow: HTTP %d - %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetWorkflowExecutions retrieves executions for a workflow
func (c *N8nClient) GetWorkflowExecutions(workflowID string, limit int) ([]interface{}, error) {
	url := fmt.Sprintf("%s/api/v1/executions?workflowId=%s&limit=%d", c.baseURL, workflowID, limit)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-N8N-API-KEY", c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workflow executions: HTTP %d", resp.StatusCode)
	}

	var response N8nResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	executions, ok := response.Data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format")
	}

	return executions, nil
}
