package providers

import (
	"time"
)

// GitProvider defines the interface that all Git providers must implement
type GitProvider interface {
	// Authentication and connection
	TestConnection() error
	GetCurrentUser() (*User, error)

	// Repository operations
	GetRepository(repoID string) (*Repository, error)
	ListRepositories() ([]*Repository, error)

	// Commit operations
	GetCommits(repoID, branch string, limit int) ([]*Commit, error)
	GetCommit(repoID, commitSHA string) (*Commit, error)

	// Branch operations
	GetBranches(repoID string) ([]*Branch, error)
	CreateBranch(repoID, branchName, fromBranch string) (*Branch, error)

	// File operations
	GetFileContent(repoID, filePath, branch string) ([]byte, error)
	UpdateFile(repoID, filePath, branch, content, message string) error

	// Pipeline/Actions operations
	GetPipelines(repoID string, limit int) ([]*Pipeline, error)
	TriggerPipeline(repoID, branch string) (*Pipeline, error)

	// Pull/Merge Request operations
	CreatePullRequest(req *CreatePullRequestInput) (*PullRequest, error)
	ListPullRequests(repoID string, state PullRequestState) ([]*PullRequest, error)

	// Provider-specific info
	GetProviderName() string
	GetAPIVersion() string
}

// Common data structures used across all providers

// User represents a Git provider user
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// Repository represents a Git repository
type Repository struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	WebURL      string `json:"web_url"`
	CloneURL    string `json:"clone_url"`
	DefaultBranch string `json:"default_branch"`
	Private     bool   `json:"private"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Commit represents a Git commit
type Commit struct {
	SHA       string    `json:"sha"`
	Message   string    `json:"message"`
	AuthorName string   `json:"author_name"`
	AuthorEmail string  `json:"author_email"`
	CreatedAt time.Time `json:"created_at"`
	WebURL    string    `json:"web_url"`
}

// Branch represents a Git branch
type Branch struct {
	Name      string  `json:"name"`
	SHA       string  `json:"sha"`
	Protected bool    `json:"protected"`
	Default   bool    `json:"default"`
	Commit    *Commit `json:"commit,omitempty"`
}

// Pipeline represents a CI/CD pipeline
type Pipeline struct {
	ID        string          `json:"id"`
	Status    PipelineStatus  `json:"status"`
	Branch    string          `json:"branch"`
	WebURL    string          `json:"web_url"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

// PipelineStatus represents the status of a pipeline
type PipelineStatus string

const (
	PipelineStatusPending PipelineStatus = "pending"
	PipelineStatusRunning PipelineStatus = "running"
	PipelineStatusSuccess PipelineStatus = "success"
	PipelineStatusFailed  PipelineStatus = "failed"
	PipelineStatusCanceled PipelineStatus = "canceled"
)

// PullRequest represents a pull/merge request
type PullRequest struct {
	ID          string           `json:"id"`
	Number      int              `json:"number"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	State       PullRequestState `json:"state"`
	SourceBranch string          `json:"source_branch"`
	TargetBranch string          `json:"target_branch"`
	WebURL      string           `json:"web_url"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

// PullRequestState represents the state of a pull request
type PullRequestState string

const (
	PullRequestStateOpen   PullRequestState = "open"
	PullRequestStateClosed PullRequestState = "closed"
	PullRequestStateMerged PullRequestState = "merged"
)

// CreatePullRequestInput contains data for creating a pull request
type CreatePullRequestInput struct {
	RepositoryID string `json:"repository_id"`
	Title        string `json:"title"`
	Description  string `json:"description"`
	SourceBranch string `json:"source_branch"`
	TargetBranch string `json:"target_branch"`
}

// ProviderConfig holds configuration for Git providers
type ProviderConfig struct {
	Type      ProviderType `json:"type"`
	BaseURL   string       `json:"base_url"`
	Token     string       `json:"token"`
	RepoID    string       `json:"repo_id"`
}

// ProviderType represents different Git provider types
type ProviderType string

const (
	ProviderTypeGitLab ProviderType = "gitlab"
	ProviderTypeGitHub ProviderType = "github"
	ProviderTypeBitbucket ProviderType = "bitbucket"
	ProviderTypeGitea  ProviderType = "gitea"
)