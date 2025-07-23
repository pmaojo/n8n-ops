package gitlab

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// GitLabClient handles GitLab API interactions
type GitLabClient struct {
	baseURL    string
	token      string
	projectID  string
	httpClient *http.Client
}

// Project represents a GitLab project
type Project struct {
	ID                int    `json:"id"`
	Name              string `json:"name"`
	Path              string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
	WebURL            string `json:"web_url"`
	DefaultBranch     string `json:"default_branch"`
}

// Commit represents a GitLab commit
type Commit struct {
	ID         string    `json:"id"`
	ShortID    string    `json:"short_id"`
	Title      string    `json:"title"`
	Message    string    `json:"message"`
	AuthorName string    `json:"author_name"`
	CreatedAt  time.Time `json:"created_at"`
}

// Pipeline represents a GitLab CI/CD pipeline
type Pipeline struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
	Ref    string `json:"ref"`
	WebURL string `json:"web_url"`
}

// NewGitLabClient creates a new GitLab API client
func NewGitLabClient(baseURL, token, projectID string) (*GitLabClient, error) {
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}
	if token == "" {
		return nil, fmt.Errorf("GitLab personal access token is required")
	}
	if projectID == "" {
		return nil, fmt.Errorf("GitLab project ID is required")
	}

	client := &GitLabClient{
		baseURL:   baseURL,
		token:     token,
		projectID: projectID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	// Test connection
	if err := client.testConnection(); err != nil {
		return nil, fmt.Errorf("failed to connect to GitLab: %w", err)
	}

	return client, nil
}

// testConnection verifies the GitLab connection
func (c *GitLabClient) testConnection() error {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v4/user", nil)
	if err != nil {
		return err
	}

	// Use PRIVATE-TOKEN header as per GitLab documentation
	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authentication failed - check personal access token")
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, resp.Status)
	}

	return nil
}

// GetProject retrieves project information
func (c *GitLabClient) GetProject() (*Project, error) {
	req, err := http.NewRequest("GET", c.baseURL+"/api/v4/projects/"+c.projectID, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get project: HTTP %d", resp.StatusCode)
	}

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return nil, fmt.Errorf("failed to decode project: %w", err)
	}

	return &project, nil
}

// GetCommits retrieves recent commits from a branch
func (c *GitLabClient) GetCommits(branch string, limit int) ([]*Commit, error) {
	if limit <= 0 {
		limit = 10
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/commits?ref_name=%s&per_page=%d",
		c.baseURL, c.projectID, branch, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get commits: HTTP %d", resp.StatusCode)
	}

	var commits []*Commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		return nil, fmt.Errorf("failed to decode commits: %w", err)
	}

	return commits, nil
}

// GetPipelines retrieves recent pipelines
func (c *GitLabClient) GetPipelines(limit int) ([]*Pipeline, error) {
	if limit <= 0 {
		limit = 10
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/pipelines?per_page=%d",
		c.baseURL, c.projectID, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get pipelines: HTTP %d", resp.StatusCode)
	}

	var pipelines []*Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipelines); err != nil {
		return nil, fmt.Errorf("failed to decode pipelines: %w", err)
	}

	return pipelines, nil
}

// TriggerPipeline triggers a new pipeline for a branch
func (c *GitLabClient) TriggerPipeline(branch string) (*Pipeline, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%s/pipeline?ref=%s",
		c.baseURL, c.projectID, branch)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to trigger pipeline: HTTP %d, %s", resp.StatusCode, string(body))
	}

	var pipeline Pipeline
	if err := json.NewDecoder(resp.Body).Decode(&pipeline); err != nil {
		return nil, fmt.Errorf("failed to decode pipeline: %w", err)
	}

	return &pipeline, nil
}

// GetFileContent retrieves file content from repository
func (c *GitLabClient) GetFileContent(filePath, branch string) ([]byte, error) {
	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s/raw?ref=%s",
		c.baseURL, c.projectID, filePath, branch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get file: HTTP %d", resp.StatusCode)
	}

	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %w", err)
	}

	return content, nil
}

// CreateMergeRequest creates a new merge request
func (c *GitLabClient) CreateMergeRequest(sourceBranch, targetBranch, title, description string) error {
	payload := map[string]interface{}{
		"source_branch": sourceBranch,
		"target_branch": targetBranch,
		"title":         title,
		"description":   description,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests", c.baseURL, c.projectID)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonPayload)))
	if err != nil {
		return err
	}

	req.Header.Set("PRIVATE-TOKEN", c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to create merge request: HTTP %d, %s", resp.StatusCode, string(body))
	}

	return nil
}
