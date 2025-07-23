package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type N8nClientInterface interface {
	GetMe() (*User, error)
	GetWorkflows() ([]Workflow, error)
	GetWorkflow(id string) (*Workflow, error)
	CreateWorkflow(workflow *Workflow) (*Workflow, error)
	UpdateWorkflow(id string, workflow *Workflow) (*Workflow, error)
	ActivateWorkflow(id string) (*Workflow, error)
	DeactivateWorkflow(id string) (*Workflow, error)
	DeleteWorkflow(id string) error
	TestConnection() error
}

type RealN8nClient struct {
	BaseURL string
	APIKey  string
	Client  *http.Client
}

type Workflow struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Active      bool                   `json:"active"`
	Nodes       []WorkflowNode         `json:"nodes"`
	Connections map[string]interface{} `json:"connections"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
	VersionId   int                    `json:"versionId"`
	Tags        []string               `json:"tags"`
}

type WorkflowNode struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	TypeVersion float64                `json:"typeVersion"`
	Position    []int                  `json:"position"`
	Parameters  map[string]interface{} `json:"parameters"`
}

type WorkflowListResponse struct {
	Data []Workflow `json:"data"`
}

type User struct {
	ID       string                 `json:"id"`
	Email    string                 `json:"email"`
	Settings map[string]interface{} `json:"settings"`
}

func NewRealN8nClient(baseURL, apiKey string) *RealN8nClient {
	return &RealN8nClient{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *RealN8nClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/v1%s", c.BaseURL, endpoint)

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-N8N-API-KEY", c.APIKey)

	return c.Client.Do(req)
}

func (c *RealN8nClient) GetMe() (*User, error) {
	// Instead of trying to get user info, just test the connection with the workflows endpoint
	resp, err := c.makeRequest("GET", "/workflows", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	// Return a dummy user since we can't get the real user info
	return &User{
		ID:       "unknown",
		Email:    "n8n-user@example.com",
		Settings: make(map[string]interface{}),
	}, nil
}

func (c *RealN8nClient) GetWorkflows() ([]Workflow, error) {
	resp, err := c.makeRequest("GET", "/workflows", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var response WorkflowListResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return response.Data, nil
}

func (c *RealN8nClient) GetWorkflow(id string) (*Workflow, error) {
	endpoint := fmt.Sprintf("/workflows/%s", id)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("workflow not found: %s", id)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var workflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflow); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &workflow, nil
}

func (c *RealN8nClient) CreateWorkflow(workflow *Workflow) (*Workflow, error) {
	resp, err := c.makeRequest("POST", "/workflows", workflow)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var createdWorkflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&createdWorkflow); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &createdWorkflow, nil
}

func (c *RealN8nClient) UpdateWorkflow(id string, workflow *Workflow) (*Workflow, error) {
	endpoint := fmt.Sprintf("/workflows/%s", id)
	resp, err := c.makeRequest("PUT", endpoint, workflow)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var updatedWorkflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&updatedWorkflow); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &updatedWorkflow, nil
}

func (c *RealN8nClient) ActivateWorkflow(id string) (*Workflow, error) {
	endpoint := fmt.Sprintf("/workflows/%s/activate", id)
	resp, err := c.makeRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var workflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflow); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &workflow, nil
}

func (c *RealN8nClient) DeactivateWorkflow(id string) (*Workflow, error) {
	endpoint := fmt.Sprintf("/workflows/%s/deactivate", id)
	resp, err := c.makeRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	var workflow Workflow
	if err := json.NewDecoder(resp.Body).Decode(&workflow); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &workflow, nil
}

func (c *RealN8nClient) DeleteWorkflow(id string) error {
	endpoint := fmt.Sprintf("/workflows/%s", id)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("API request failed with status: %d", resp.StatusCode)
	}

	return nil
}

func (c *RealN8nClient) TestConnection() error {
	// Test connection by getting workflows instead of user info
	_, err := c.GetWorkflows()
	return err
}

// Factory function to create the appropriate client
func NewN8nClientWithDemo(baseURL, apiKey string, demoMode bool) N8nClientInterface {
	if demoMode {
		return NewDemoN8nClient()
	}
	return NewRealN8nClient(baseURL, apiKey)
}
