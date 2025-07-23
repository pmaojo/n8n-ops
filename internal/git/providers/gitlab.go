package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GitLabProvider implements GitProvider for GitLab
type GitLabProvider struct {
	baseURL    string
	token      string
	repoID     string
	httpClient *http.Client
}

// NewGitLabProvider creates a new GitLab provider instance
func NewGitLabProvider(config *ProviderConfig) (*GitLabProvider, error) {
	if config.Token == "" {
		return nil, fmt.Errorf("GitLab personal access token is required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://gitlab.com"
	}

	provider := &GitLabProvider{
		baseURL: baseURL,
		token:   config.Token,
		repoID:  config.RepoID,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	return provider, nil
}

// GetProviderName returns the provider name
func (g *GitLabProvider) GetProviderName() string {
	return "gitlab"
}

// GetAPIVersion returns the API version
func (g *GitLabProvider) GetAPIVersion() string {
	return "v4"
}

// TestConnection verifies the connection to GitLab
func (g *GitLabProvider) TestConnection() error {
	req, err := http.NewRequest("GET", g.baseURL+"/api/v4/user", nil)
	if err != nil {
		return err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
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

// GetCurrentUser retrieves current user information
func (g *GitLabProvider) GetCurrentUser() (*User, error) {
	req, err := http.NewRequest("GET", g.baseURL+"/api/v4/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: HTTP %d", resp.StatusCode)
	}

	var gitlabUser struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabUser); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return &User{
		ID:        strconv.Itoa(gitlabUser.ID),
		Username:  gitlabUser.Username,
		Email:     gitlabUser.Email,
		Name:      gitlabUser.Name,
		AvatarURL: gitlabUser.AvatarURL,
	}, nil
}

// GetRepository retrieves repository information
func (g *GitLabProvider) GetRepository(repoID string) (*Repository, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	req, err := http.NewRequest("GET", g.baseURL+"/api/v4/projects/"+repoID, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repository: HTTP %d", resp.StatusCode)
	}

	var gitlabProject struct {
		ID                int       `json:"id"`
		Name              string    `json:"name"`
		Path              string    `json:"path"`
		PathWithNamespace string    `json:"path_with_namespace"`
		Description       string    `json:"description"`
		WebURL            string    `json:"web_url"`
		HTTPURLToRepo     string    `json:"http_url_to_repo"`
		DefaultBranch     string    `json:"default_branch"`
		Visibility        string    `json:"visibility"`
		CreatedAt         time.Time `json:"created_at"`
		LastActivityAt    time.Time `json:"last_activity_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabProject); err != nil {
		return nil, fmt.Errorf("failed to decode repository: %w", err)
	}

	return &Repository{
		ID:            strconv.Itoa(gitlabProject.ID),
		Name:          gitlabProject.Name,
		FullName:      gitlabProject.PathWithNamespace,
		Description:   gitlabProject.Description,
		WebURL:        gitlabProject.WebURL,
		CloneURL:      gitlabProject.HTTPURLToRepo,
		DefaultBranch: gitlabProject.DefaultBranch,
		Private:       gitlabProject.Visibility == "private",
		CreatedAt:     gitlabProject.CreatedAt,
		UpdatedAt:     gitlabProject.LastActivityAt,
	}, nil
}

// ListRepositories retrieves list of repositories
func (g *GitLabProvider) ListRepositories() ([]*Repository, error) {
	req, err := http.NewRequest("GET", g.baseURL+"/api/v4/projects?membership=true&per_page=100", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list repositories: HTTP %d", resp.StatusCode)
	}

	var gitlabProjects []struct {
		ID                int       `json:"id"`
		Name              string    `json:"name"`
		PathWithNamespace string    `json:"path_with_namespace"`
		Description       string    `json:"description"`
		WebURL            string    `json:"web_url"`
		HTTPURLToRepo     string    `json:"http_url_to_repo"`
		DefaultBranch     string    `json:"default_branch"`
		Visibility        string    `json:"visibility"`
		CreatedAt         time.Time `json:"created_at"`
		LastActivityAt    time.Time `json:"last_activity_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabProjects); err != nil {
		return nil, fmt.Errorf("failed to decode repositories: %w", err)
	}

	repositories := make([]*Repository, len(gitlabProjects))
	for i, project := range gitlabProjects {
		repositories[i] = &Repository{
			ID:            strconv.Itoa(project.ID),
			Name:          project.Name,
			FullName:      project.PathWithNamespace,
			Description:   project.Description,
			WebURL:        project.WebURL,
			CloneURL:      project.HTTPURLToRepo,
			DefaultBranch: project.DefaultBranch,
			Private:       project.Visibility == "private",
			CreatedAt:     project.CreatedAt,
			UpdatedAt:     project.LastActivityAt,
		}
	}

	return repositories, nil
}

// GetCommits retrieves commits from a repository
func (g *GitLabProvider) GetCommits(repoID, branch string, limit int) ([]*Commit, error) {
	if repoID == "" {
		repoID = g.repoID
	}
	if limit <= 0 {
		limit = 10
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/commits?ref_name=%s&per_page=%d",
		g.baseURL, repoID, branch, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get commits: HTTP %d", resp.StatusCode)
	}

	var gitlabCommits []struct {
		ID          string    `json:"id"`
		ShortID     string    `json:"short_id"`
		Title       string    `json:"title"`
		Message     string    `json:"message"`
		AuthorName  string    `json:"author_name"`
		AuthorEmail string    `json:"author_email"`
		CreatedAt   time.Time `json:"created_at"`
		WebURL      string    `json:"web_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabCommits); err != nil {
		return nil, fmt.Errorf("failed to decode commits: %w", err)
	}

	commits := make([]*Commit, len(gitlabCommits))
	for i, commit := range gitlabCommits {
		commits[i] = &Commit{
			SHA:         commit.ID,
			Message:     commit.Message,
			AuthorName:  commit.AuthorName,
			AuthorEmail: commit.AuthorEmail,
			CreatedAt:   commit.CreatedAt,
			WebURL:      commit.WebURL,
		}
	}

	return commits, nil
}

// GetCommit retrieves a specific commit
func (g *GitLabProvider) GetCommit(repoID, commitSHA string) (*Commit, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/commits/%s", g.baseURL, repoID, commitSHA)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get commit: HTTP %d", resp.StatusCode)
	}

	var gitlabCommit struct {
		ID          string    `json:"id"`
		Message     string    `json:"message"`
		AuthorName  string    `json:"author_name"`
		AuthorEmail string    `json:"author_email"`
		CreatedAt   time.Time `json:"created_at"`
		WebURL      string    `json:"web_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabCommit); err != nil {
		return nil, fmt.Errorf("failed to decode commit: %w", err)
	}

	return &Commit{
		SHA:         gitlabCommit.ID,
		Message:     gitlabCommit.Message,
		AuthorName:  gitlabCommit.AuthorName,
		AuthorEmail: gitlabCommit.AuthorEmail,
		CreatedAt:   gitlabCommit.CreatedAt,
		WebURL:      gitlabCommit.WebURL,
	}, nil
}

// GetBranches retrieves branches from a repository
func (g *GitLabProvider) GetBranches(repoID string) ([]*Branch, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/branches", g.baseURL, repoID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get branches: HTTP %d", resp.StatusCode)
	}

	var gitlabBranches []struct {
		Name      string `json:"name"`
		Protected bool   `json:"protected"`
		Default   bool   `json:"default"`
		Commit    struct {
			ID string `json:"id"`
		} `json:"commit"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabBranches); err != nil {
		return nil, fmt.Errorf("failed to decode branches: %w", err)
	}

	branches := make([]*Branch, len(gitlabBranches))
	for i, branch := range gitlabBranches {
		branches[i] = &Branch{
			Name:      branch.Name,
			SHA:       branch.Commit.ID,
			Protected: branch.Protected,
			Default:   branch.Default,
		}
	}

	return branches, nil
}

// CreateBranch creates a new branch
func (g *GitLabProvider) CreateBranch(repoID, branchName, fromBranch string) (*Branch, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	payload := map[string]string{
		"branch": branchName,
		"ref":    fromBranch,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/branches", g.baseURL, repoID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create branch: HTTP %d, %s", resp.StatusCode, string(body))
	}

	var gitlabBranch struct {
		Name      string `json:"name"`
		Protected bool   `json:"protected"`
		Commit    struct {
			ID string `json:"id"`
		} `json:"commit"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabBranch); err != nil {
		return nil, fmt.Errorf("failed to decode branch: %w", err)
	}

	return &Branch{
		Name:      gitlabBranch.Name,
		SHA:       gitlabBranch.Commit.ID,
		Protected: gitlabBranch.Protected,
		Default:   false,
	}, nil
}

// GetFileContent retrieves file content from repository
func (g *GitLabProvider) GetFileContent(repoID, filePath, branch string) ([]byte, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	// URL encode the file path
	encodedPath := strings.ReplaceAll(filePath, "/", "%2F")
	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s/raw?ref=%s",
		g.baseURL, repoID, encodedPath, branch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)

	resp, err := g.httpClient.Do(req)
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

// UpdateFile updates a file in the repository
func (g *GitLabProvider) UpdateFile(repoID, filePath, branch, content, message string) error {
	if repoID == "" {
		repoID = g.repoID
	}

	encodedPath := strings.ReplaceAll(filePath, "/", "%2F")

	payload := map[string]string{
		"branch":         branch,
		"content":        content,
		"commit_message": message,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/repository/files/%s", g.baseURL, repoID, encodedPath)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("failed to update file: HTTP %d, %s", resp.StatusCode, string(body))
	}

	return nil
}

// GetPipelines retrieves pipelines from a repository
func (g *GitLabProvider) GetPipelines(repoID string, limit int) ([]*Pipeline, error) {
	if repoID == "" {
		repoID = g.repoID
	}
	if limit <= 0 {
		limit = 10
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/pipelines?per_page=%d", g.baseURL, repoID, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get pipelines: HTTP %d", resp.StatusCode)
	}

	var gitlabPipelines []struct {
		ID        int       `json:"id"`
		Status    string    `json:"status"`
		Ref       string    `json:"ref"`
		WebURL    string    `json:"web_url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabPipelines); err != nil {
		return nil, fmt.Errorf("failed to decode pipelines: %w", err)
	}

	pipelines := make([]*Pipeline, len(gitlabPipelines))
	for i, pipeline := range gitlabPipelines {
		pipelines[i] = &Pipeline{
			ID:        strconv.Itoa(pipeline.ID),
			Status:    PipelineStatus(pipeline.Status),
			Branch:    pipeline.Ref,
			WebURL:    pipeline.WebURL,
			CreatedAt: pipeline.CreatedAt,
			UpdatedAt: pipeline.UpdatedAt,
		}
	}

	return pipelines, nil
}

// TriggerPipeline triggers a new pipeline
func (g *GitLabProvider) TriggerPipeline(repoID, branch string) (*Pipeline, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/pipeline?ref=%s", g.baseURL, repoID, branch)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to trigger pipeline: HTTP %d, %s", resp.StatusCode, string(body))
	}

	var gitlabPipeline struct {
		ID        int       `json:"id"`
		Status    string    `json:"status"`
		Ref       string    `json:"ref"`
		WebURL    string    `json:"web_url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabPipeline); err != nil {
		return nil, fmt.Errorf("failed to decode pipeline: %w", err)
	}

	return &Pipeline{
		ID:        strconv.Itoa(gitlabPipeline.ID),
		Status:    PipelineStatus(gitlabPipeline.Status),
		Branch:    gitlabPipeline.Ref,
		WebURL:    gitlabPipeline.WebURL,
		CreatedAt: gitlabPipeline.CreatedAt,
		UpdatedAt: gitlabPipeline.UpdatedAt,
	}, nil
}

// CreatePullRequest creates a merge request in GitLab
func (g *GitLabProvider) CreatePullRequest(req *CreatePullRequestInput) (*PullRequest, error) {
	repoID := req.RepositoryID
	if repoID == "" {
		repoID = g.repoID
	}

	payload := map[string]string{
		"source_branch": req.SourceBranch,
		"target_branch": req.TargetBranch,
		"title":         req.Title,
		"description":   req.Description,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests", g.baseURL, repoID)
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("PRIVATE-TOKEN", g.token)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create merge request: HTTP %d, %s", resp.StatusCode, string(body))
	}

	var gitlabMR struct {
		ID           int       `json:"id"`
		IID          int       `json:"iid"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		State        string    `json:"state"`
		SourceBranch string    `json:"source_branch"`
		TargetBranch string    `json:"target_branch"`
		WebURL       string    `json:"web_url"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabMR); err != nil {
		return nil, fmt.Errorf("failed to decode merge request: %w", err)
	}

	return &PullRequest{
		ID:           strconv.Itoa(gitlabMR.ID),
		Number:       gitlabMR.IID,
		Title:        gitlabMR.Title,
		Description:  gitlabMR.Description,
		State:        PullRequestState(gitlabMR.State),
		SourceBranch: gitlabMR.SourceBranch,
		TargetBranch: gitlabMR.TargetBranch,
		WebURL:       gitlabMR.WebURL,
		CreatedAt:    gitlabMR.CreatedAt,
		UpdatedAt:    gitlabMR.UpdatedAt,
	}, nil
}

// ListPullRequests lists merge requests from a repository
func (g *GitLabProvider) ListPullRequests(repoID string, state PullRequestState) ([]*PullRequest, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	url := fmt.Sprintf("%s/api/v4/projects/%s/merge_requests?state=%s", g.baseURL, repoID, string(state))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("PRIVATE-TOKEN", g.token)
	req.Header.Set("Accept", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list merge requests: HTTP %d", resp.StatusCode)
	}

	var gitlabMRs []struct {
		ID           int       `json:"id"`
		IID          int       `json:"iid"`
		Title        string    `json:"title"`
		Description  string    `json:"description"`
		State        string    `json:"state"`
		SourceBranch string    `json:"source_branch"`
		TargetBranch string    `json:"target_branch"`
		WebURL       string    `json:"web_url"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gitlabMRs); err != nil {
		return nil, fmt.Errorf("failed to decode merge requests: %w", err)
	}

	pullRequests := make([]*PullRequest, len(gitlabMRs))
	for i, mr := range gitlabMRs {
		pullRequests[i] = &PullRequest{
			ID:           strconv.Itoa(mr.ID),
			Number:       mr.IID,
			Title:        mr.Title,
			Description:  mr.Description,
			State:        PullRequestState(mr.State),
			SourceBranch: mr.SourceBranch,
			TargetBranch: mr.TargetBranch,
			WebURL:       mr.WebURL,
			CreatedAt:    mr.CreatedAt,
			UpdatedAt:    mr.UpdatedAt,
		}
	}

	return pullRequests, nil
}
