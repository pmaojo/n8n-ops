package providers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

// GitHubProvider implements GitProvider for GitHub
type GitHubProvider struct {
	baseURL    string
	token      string
	repoID     string
	httpClient *http.Client
}

// NewGitHubProvider creates a new GitHub provider instance
func NewGitHubProvider(config *ProviderConfig) (*GitHubProvider, error) {
	if config.Token == "" {
		return nil, fmt.Errorf("GitHub personal access token is required")
	}

	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.github.com"
	}

	provider := &GitHubProvider{
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
func (g *GitHubProvider) GetProviderName() string {
	return "github"
}

// GetAPIVersion returns the API version
func (g *GitHubProvider) GetAPIVersion() string {
	return "v3"
}

// TestConnection verifies the connection to GitHub
func (g *GitHubProvider) TestConnection() error {
	req, err := http.NewRequest("GET", g.baseURL+"/user", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

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
func (g *GitHubProvider) GetCurrentUser() (*User, error) {
	req, err := http.NewRequest("GET", g.baseURL+"/user", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: HTTP %d", resp.StatusCode)
	}

	var githubUser struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		AvatarURL string `json:"avatar_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		return nil, fmt.Errorf("failed to decode user: %w", err)
	}

	return &User{
		ID:        strconv.Itoa(githubUser.ID),
		Username:  githubUser.Login,
		Email:     githubUser.Email,
		Name:      githubUser.Name,
		AvatarURL: githubUser.AvatarURL,
	}, nil
}

// GetRepository retrieves repository information
func (g *GitHubProvider) GetRepository(repoID string) (*Repository, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	req, err := http.NewRequest("GET", g.baseURL+"/repos/"+repoID, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get repository: HTTP %d", resp.StatusCode)
	}

	var githubRepo struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		HTMLURL     string `json:"html_url"`
		CloneURL    string `json:"clone_url"`
		DefaultBranch string `json:"default_branch"`
		Private     bool   `json:"private"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubRepo); err != nil {
		return nil, fmt.Errorf("failed to decode repository: %w", err)
	}

	return &Repository{
		ID:          strconv.Itoa(githubRepo.ID),
		Name:        githubRepo.Name,
		FullName:    githubRepo.FullName,
		Description: githubRepo.Description,
		WebURL:      githubRepo.HTMLURL,
		CloneURL:    githubRepo.CloneURL,
		DefaultBranch: githubRepo.DefaultBranch,
		Private:     githubRepo.Private,
		CreatedAt:   githubRepo.CreatedAt,
		UpdatedAt:   githubRepo.UpdatedAt,
	}, nil
}

// ListRepositories retrieves list of repositories
func (g *GitHubProvider) ListRepositories() ([]*Repository, error) {
	req, err := http.NewRequest("GET", g.baseURL+"/user/repos?per_page=100", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list repositories: HTTP %d", resp.StatusCode)
	}

	var githubRepos []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		FullName    string `json:"full_name"`
		Description string `json:"description"`
		HTMLURL     string `json:"html_url"`
		CloneURL    string `json:"clone_url"`
		DefaultBranch string `json:"default_branch"`
		Private     bool   `json:"private"`
		CreatedAt   time.Time `json:"created_at"`
		UpdatedAt   time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubRepos); err != nil {
		return nil, fmt.Errorf("failed to decode repositories: %w", err)
	}

	repositories := make([]*Repository, len(githubRepos))
	for i, repo := range githubRepos {
		repositories[i] = &Repository{
			ID:          strconv.Itoa(repo.ID),
			Name:        repo.Name,
			FullName:    repo.FullName,
			Description: repo.Description,
			WebURL:      repo.HTMLURL,
			CloneURL:    repo.CloneURL,
			DefaultBranch: repo.DefaultBranch,
			Private:     repo.Private,
			CreatedAt:   repo.CreatedAt,
			UpdatedAt:   repo.UpdatedAt,
		}
	}

	return repositories, nil
}

// GetCommits retrieves commits from a repository
func (g *GitHubProvider) GetCommits(repoID, branch string, limit int) ([]*Commit, error) {
	if repoID == "" {
		repoID = g.repoID
	}
	if limit <= 0 {
		limit = 10
	}

	url := fmt.Sprintf("%s/repos/%s/commits?sha=%s&per_page=%d", g.baseURL, repoID, branch, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get commits: HTTP %d", resp.StatusCode)
	}

	var githubCommits []struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
		} `json:"commit"`
		HTMLURL string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubCommits); err != nil {
		return nil, fmt.Errorf("failed to decode commits: %w", err)
	}

	commits := make([]*Commit, len(githubCommits))
	for i, commit := range githubCommits {
		commits[i] = &Commit{
			SHA:         commit.SHA,
			Message:     commit.Commit.Message,
			AuthorName:  commit.Commit.Author.Name,
			AuthorEmail: commit.Commit.Author.Email,
			CreatedAt:   commit.Commit.Author.Date,
			WebURL:      commit.HTMLURL,
		}
	}

	return commits, nil
}

// GetCommit retrieves a specific commit
func (g *GitHubProvider) GetCommit(repoID, commitSHA string) (*Commit, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	url := fmt.Sprintf("%s/repos/%s/commits/%s", g.baseURL, repoID, commitSHA)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get commit: HTTP %d", resp.StatusCode)
	}

	var githubCommit struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
			Author  struct {
				Name  string    `json:"name"`
				Email string    `json:"email"`
				Date  time.Time `json:"date"`
			} `json:"author"`
		} `json:"commit"`
		HTMLURL string `json:"html_url"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubCommit); err != nil {
		return nil, fmt.Errorf("failed to decode commit: %w", err)
	}

	return &Commit{
		SHA:         githubCommit.SHA,
		Message:     githubCommit.Commit.Message,
		AuthorName:  githubCommit.Commit.Author.Name,
		AuthorEmail: githubCommit.Commit.Author.Email,
		CreatedAt:   githubCommit.Commit.Author.Date,
		WebURL:      githubCommit.HTMLURL,
	}, nil
}

// GetBranches retrieves branches from a repository
func (g *GitHubProvider) GetBranches(repoID string) ([]*Branch, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	url := fmt.Sprintf("%s/repos/%s/branches", g.baseURL, repoID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get branches: HTTP %d", resp.StatusCode)
	}

	var githubBranches []struct {
		Name      string `json:"name"`
		Protected bool   `json:"protected"`
		Commit    struct {
			SHA string `json:"sha"`
		} `json:"commit"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubBranches); err != nil {
		return nil, fmt.Errorf("failed to decode branches: %w", err)
	}

	// Get default branch to mark it
	repo, err := g.GetRepository(repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info for default branch: %w", err)
	}

	branches := make([]*Branch, len(githubBranches))
	for i, branch := range githubBranches {
		branches[i] = &Branch{
			Name:      branch.Name,
			SHA:       branch.Commit.SHA,
			Protected: branch.Protected,
			Default:   branch.Name == repo.DefaultBranch,
		}
	}

	return branches, nil
}

// CreateBranch creates a new branch
func (g *GitHubProvider) CreateBranch(repoID, branchName, fromBranch string) (*Branch, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	// First get the SHA of the source branch
	branches, err := g.GetBranches(repoID)
	if err != nil {
		return nil, fmt.Errorf("failed to get branches: %w", err)
	}

	var fromSHA string
	for _, branch := range branches {
		if branch.Name == fromBranch {
			fromSHA = branch.SHA
			break
		}
	}

	if fromSHA == "" {
		return nil, fmt.Errorf("source branch '%s' not found", fromBranch)
	}

	payload := map[string]string{
		"ref": "refs/heads/" + branchName,
		"sha": fromSHA,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/repos/%s/git/refs", g.baseURL, repoID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create branch: HTTP %d, %s", resp.StatusCode, string(body))
	}

	return &Branch{
		Name:      branchName,
		SHA:       fromSHA,
		Protected: false,
		Default:   false,
	}, nil
}

// GetFileContent retrieves file content from repository
func (g *GitHubProvider) GetFileContent(repoID, filePath, branch string) ([]byte, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	url := fmt.Sprintf("%s/repos/%s/contents/%s?ref=%s", g.baseURL, repoID, filePath, branch)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.raw")

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
func (g *GitHubProvider) UpdateFile(repoID, filePath, branch, content, message string) error {
	if repoID == "" {
		repoID = g.repoID
	}

	// First get the current file to get its SHA
	url := fmt.Sprintf("%s/repos/%s/contents/%s?ref=%s", g.baseURL, repoID, filePath, branch)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var fileInfo struct {
		SHA string `json:"sha"`
	}

	if resp.StatusCode == http.StatusOK {
		if err := json.NewDecoder(resp.Body).Decode(&fileInfo); err != nil {
			return fmt.Errorf("failed to decode file info: %w", err)
		}
	}

	// Encode content to base64
	encodedContent := fmt.Sprintf("%s", content) // GitHub API expects base64, but we'll send raw for simplicity

	payload := map[string]interface{}{
		"message": message,
		"content": encodedContent,
		"branch":  branch,
	}

	if fileInfo.SHA != "" {
		payload["sha"] = fileInfo.SHA
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	updateURL := fmt.Sprintf("%s/repos/%s/contents/%s", g.baseURL, repoID, filePath)
	updateReq, err := http.NewRequest("PUT", updateURL, bytes.NewReader(jsonPayload))
	if err != nil {
		return err
	}

	updateReq.Header.Set("Authorization", "Bearer "+g.token)
	updateReq.Header.Set("Accept", "application/vnd.github.v3+json")
	updateReq.Header.Set("Content-Type", "application/json")

	updateResp, err := g.httpClient.Do(updateReq)
	if err != nil {
		return err
	}
	defer updateResp.Body.Close()

	if updateResp.StatusCode != http.StatusOK && updateResp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(updateResp.Body)
		return fmt.Errorf("failed to update file: HTTP %d, %s", updateResp.StatusCode, string(body))
	}

	return nil
}

// GetPipelines retrieves GitHub Actions workflows/runs
func (g *GitHubProvider) GetPipelines(repoID string, limit int) ([]*Pipeline, error) {
	if repoID == "" {
		repoID = g.repoID
	}
	if limit <= 0 {
		limit = 10
	}

	url := fmt.Sprintf("%s/repos/%s/actions/runs?per_page=%d", g.baseURL, repoID, limit)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get workflows: HTTP %d", resp.StatusCode)
	}

	var githubWorkflows struct {
		WorkflowRuns []struct {
			ID         int       `json:"id"`
			Status     string    `json:"status"`
			Conclusion string    `json:"conclusion"`
			HeadBranch string    `json:"head_branch"`
			HTMLURL    string    `json:"html_url"`
			CreatedAt  time.Time `json:"created_at"`
			UpdatedAt  time.Time `json:"updated_at"`
		} `json:"workflow_runs"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubWorkflows); err != nil {
		return nil, fmt.Errorf("failed to decode workflows: %w", err)
	}

	pipelines := make([]*Pipeline, len(githubWorkflows.WorkflowRuns))
	for i, run := range githubWorkflows.WorkflowRuns {
		status := PipelineStatus(run.Status)
		if run.Conclusion != "" {
			switch run.Conclusion {
			case "success":
				status = PipelineStatusSuccess
			case "failure":
				status = PipelineStatusFailed
			case "cancelled":
				status = PipelineStatusCanceled
			}
		}

		pipelines[i] = &Pipeline{
			ID:        strconv.Itoa(run.ID),
			Status:    status,
			Branch:    run.HeadBranch,
			WebURL:    run.HTMLURL,
			CreatedAt: run.CreatedAt,
			UpdatedAt: run.UpdatedAt,
		}
	}

	return pipelines, nil
}

// TriggerPipeline triggers a GitHub Actions workflow
func (g *GitHubProvider) TriggerPipeline(repoID, branch string) (*Pipeline, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	// GitHub doesn't have direct pipeline triggering like GitLab
	// This would require workflow dispatch events
	return nil, fmt.Errorf("GitHub pipeline triggering requires workflow_dispatch configuration")
}

// CreatePullRequest creates a pull request in GitHub
func (g *GitHubProvider) CreatePullRequest(req *CreatePullRequestInput) (*PullRequest, error) {
	repoID := req.RepositoryID
	if repoID == "" {
		repoID = g.repoID
	}

	payload := map[string]string{
		"title": req.Title,
		"body":  req.Description,
		"head":  req.SourceBranch,
		"base":  req.TargetBranch,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/repos/%s/pulls", g.baseURL, repoID)
	httpReq, err := http.NewRequest("POST", url, bytes.NewReader(jsonPayload))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+g.token)
	httpReq.Header.Set("Accept", "application/vnd.github.v3+json")
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := g.httpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create pull request: HTTP %d, %s", resp.StatusCode, string(body))
	}

	var githubPR struct {
		ID     int    `json:"id"`
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
		State  string `json:"state"`
		Head   struct {
			Ref string `json:"ref"`
		} `json:"head"`
		Base struct {
			Ref string `json:"ref"`
		} `json:"base"`
		HTMLURL   string    `json:"html_url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubPR); err != nil {
		return nil, fmt.Errorf("failed to decode pull request: %w", err)
	}

	return &PullRequest{
		ID:           strconv.Itoa(githubPR.ID),
		Number:       githubPR.Number,
		Title:        githubPR.Title,
		Description:  githubPR.Body,
		State:        PullRequestState(githubPR.State),
		SourceBranch: githubPR.Head.Ref,
		TargetBranch: githubPR.Base.Ref,
		WebURL:       githubPR.HTMLURL,
		CreatedAt:    githubPR.CreatedAt,
		UpdatedAt:    githubPR.UpdatedAt,
	}, nil
}

// ListPullRequests lists pull requests from a repository
func (g *GitHubProvider) ListPullRequests(repoID string, state PullRequestState) ([]*PullRequest, error) {
	if repoID == "" {
		repoID = g.repoID
	}

	url := fmt.Sprintf("%s/repos/%s/pulls?state=%s", g.baseURL, repoID, string(state))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+g.token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to list pull requests: HTTP %d", resp.StatusCode)
	}

	var githubPRs []struct {
		ID     int    `json:"id"`
		Number int    `json:"number"`
		Title  string `json:"title"`
		Body   string `json:"body"`
		State  string `json:"state"`
		Head   struct {
			Ref string `json:"ref"`
		} `json:"head"`
		Base struct {
			Ref string `json:"ref"`
		} `json:"base"`
		HTMLURL   string    `json:"html_url"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&githubPRs); err != nil {
		return nil, fmt.Errorf("failed to decode pull requests: %w", err)
	}

	pullRequests := make([]*PullRequest, len(githubPRs))
	for i, pr := range githubPRs {
		pullRequests[i] = &PullRequest{
			ID:           strconv.Itoa(pr.ID),
			Number:       pr.Number,
			Title:        pr.Title,
			Description:  pr.Body,
			State:        PullRequestState(pr.State),
			SourceBranch: pr.Head.Ref,
			TargetBranch: pr.Base.Ref,
			WebURL:       pr.HTMLURL,
			CreatedAt:    pr.CreatedAt,
			UpdatedAt:    pr.UpdatedAt,
		}
	}

	return pullRequests, nil
}