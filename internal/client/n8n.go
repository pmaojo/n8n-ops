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

// CreateWorkflow creates a new workflow in n8n
func (c *N8nClient) CreateWorkflow(wf *workflow.Workflow) (*workflow.Workflow, error) {
        data, err := json.Marshal(wf)
        if err != nil {
                return nil, fmt.Errorf("failed to marshal workflow: %w", err)
        }

        req, err := http.NewRequest("POST", c.baseURL+"/api/v1/workflows", bytes.NewBuffer(data))
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

        var createdWf workflow.Workflow
        if err := json.NewDecoder(resp.Body).Decode(&createdWf); err != nil {
                return nil, fmt.Errorf("failed to decode created workflow: %w", err)
        }

        return &createdWf, nil
}

// UpdateWorkflow updates an existing workflow in n8n
func (c *N8nClient) UpdateWorkflow(id string, wf *workflow.Workflow) (*workflow.Workflow, error) {
        data, err := json.Marshal(wf)
        if err != nil {
                return nil, fmt.Errorf("failed to marshal workflow: %w", err)
        }

        req, err := http.NewRequest("PUT", c.baseURL+"/api/v1/workflows/"+id, bytes.NewBuffer(data))
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

        if resp.StatusCode == http.StatusNotFound {
                return nil, fmt.Errorf("workflow not found")
        }
        if resp.StatusCode != http.StatusOK {
                body, _ := io.ReadAll(resp.Body)
                return nil, fmt.Errorf("failed to update workflow: HTTP %d - %s", resp.StatusCode, string(body))
        }

        var updatedWf workflow.Workflow
        if err := json.NewDecoder(resp.Body).Decode(&updatedWf); err != nil {
                return nil, fmt.Errorf("failed to decode updated workflow: %w", err)
        }

        return &updatedWf, nil
}

// DeleteWorkflow deletes a workflow from n8n
func (c *N8nClient) DeleteWorkflow(id string) error {
        req, err := http.NewRequest("DELETE", c.baseURL+"/api/v1/workflows/"+id, nil)
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

        if resp.StatusCode == http.StatusNotFound {
                return fmt.Errorf("workflow not found")
        }
        if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
                return fmt.Errorf("failed to delete workflow: HTTP %d", resp.StatusCode)
        }

        return nil
}

// ActivateWorkflow activates a workflow in n8n
func (c *N8nClient) ActivateWorkflow(id string) error {
        data := map[string]interface{}{
                "active": true,
        }
        jsonData, err := json.Marshal(data)
        if err != nil {
                return err
        }

        req, err := http.NewRequest("PATCH", c.baseURL+"/api/v1/workflows/"+id+"/activate", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }

        req.Header.Set("X-N8N-API-KEY", c.apiKey)
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json")

        resp, err := c.httpClient.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("failed to activate workflow: HTTP %d", resp.StatusCode)
        }

        return nil
}

// DeactivateWorkflow deactivates a workflow in n8n
func (c *N8nClient) DeactivateWorkflow(id string) error {
        data := map[string]interface{}{
                "active": false,
        }
        jsonData, err := json.Marshal(data)
        if err != nil {
                return err
        }

        req, err := http.NewRequest("PATCH", c.baseURL+"/api/v1/workflows/"+id+"/activate", bytes.NewBuffer(jsonData))
        if err != nil {
                return err
        }

        req.Header.Set("X-N8N-API-KEY", c.apiKey)
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("Accept", "application/json")

        resp, err := c.httpClient.Do(req)
        if err != nil {
                return err
        }
        defer resp.Body.Close()

        if resp.StatusCode != http.StatusOK {
                return fmt.Errorf("failed to deactivate workflow: HTTP %d", resp.StatusCode)
        }

        return nil
}

// ExecuteWorkflow manually executes a workflow
func (c *N8nClient) ExecuteWorkflow(id string, data map[string]interface{}) (map[string]interface{}, error) {
        jsonData, err := json.Marshal(data)
        if err != nil {
                return nil, err
        }

        req, err := http.NewRequest("POST", c.baseURL+"/api/v1/workflows/"+id+"/execute", bytes.NewBuffer(jsonData))
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
                return nil, fmt.Errorf("failed to execute workflow: HTTP %d", resp.StatusCode)
        }

        var result map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
                return nil, fmt.Errorf("failed to decode execution result: %w", err)
        }

        return result, nil
}

// GetExecutions retrieves workflow execution history
func (c *N8nClient) GetExecutions(workflowId string, limit int) ([]map[string]interface{}, error) {
        url := fmt.Sprintf("%s/api/v1/executions?workflowId=%s", c.baseURL, workflowId)
        if limit > 0 {
                url += fmt.Sprintf("&limit=%d", limit)
        }

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
                return nil, fmt.Errorf("failed to get executions: HTTP %d", resp.StatusCode)
        }

        var response struct {
                Data []map[string]interface{} `json:"data"`
        }
        if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
                return nil, fmt.Errorf("failed to decode executions: %w", err)
        }

        return response.Data, nil
}

// GetNodeTypes retrieves available node types
func (c *N8nClient) GetNodeTypes() ([]map[string]interface{}, error) {
        req, err := http.NewRequest("GET", c.baseURL+"/api/v1/node-types", nil)
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
                return nil, fmt.Errorf("failed to get node types: HTTP %d", resp.StatusCode)
        }

        var nodeTypes []map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&nodeTypes); err != nil {
                return nil, fmt.Errorf("failed to decode node types: %w", err)
        }

        return nodeTypes, nil
}

// Health check for n8n instance
func (c *N8nClient) HealthCheck() (map[string]interface{}, error) {
        req, err := http.NewRequest("GET", c.baseURL+"/healthz", nil)
        if err != nil {
                return nil, err
        }

        resp, err := c.httpClient.Do(req)
        if err != nil {
                return nil, err
        }
        defer resp.Body.Close()

        var health map[string]interface{}
        if resp.StatusCode == http.StatusOK {
                json.NewDecoder(resp.Body).Decode(&health)
        }

        health["status_code"] = resp.StatusCode
        health["healthy"] = resp.StatusCode == http.StatusOK

        return health, nil
}

// DeployWorkflows deploys multiple workflows at once with transaction-like behavior
func (c *N8nClient) DeployWorkflows(workflows []*workflow.Workflow) ([]string, error) {
        var deployedIds []string
        var errors []string

        for _, wf := range workflows {
                if wf.ID != "" {
                        // Update existing workflow
                        updatedWf, err := c.UpdateWorkflow(wf.ID, wf)
                        if err != nil {
                                errors = append(errors, fmt.Sprintf("Failed to update workflow %s: %v", wf.Name, err))
                                continue
                        }
                        deployedIds = append(deployedIds, updatedWf.ID)
                } else {
                        // Create new workflow
                        createdWf, err := c.CreateWorkflow(wf)
                        if err != nil {
                                errors = append(errors, fmt.Sprintf("Failed to create workflow %s: %v", wf.Name, err))
                                continue
                        }
                        deployedIds = append(deployedIds, createdWf.ID)
                }
        }

        if len(errors) > 0 {
                return deployedIds, fmt.Errorf("deployment errors: %v", errors)
        }

        return deployedIds, nil
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
