package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/pmaojo/n8n-ops/internal/git"
	"github.com/pmaojo/n8n-ops/internal/utils"
	"github.com/pmaojo/n8n-ops/internal/workflow"
	"github.com/sirupsen/logrus"
)

// Service handles synchronization of workflows between n8n and the filesystem.
type Service struct {
	Client            client.Client
	CredentialManager *credentials.CredentialManager
	GitChecker        *git.GitStatusChecker
	Logger            logrus.FieldLogger
	Environment       string
}

// Options configures the sync operation.
type Options struct {
	OutputDir string
	Force     bool
	DryRun    bool
}

// NewService constructs a new Service.
func NewService(cli client.Client, cm *credentials.CredentialManager, checker *git.GitStatusChecker, logger logrus.FieldLogger, env string) *Service {
	return &Service{
		Client:            cli,
		CredentialManager: cm,
		GitChecker:        checker,
		Logger:            logger,
		Environment:       env,
	}
}

// Sync performs the workflow synchronization using the provided options.
func (s *Service) Sync(ctx context.Context, opts Options) error {
	if opts.OutputDir == "" {
		opts.OutputDir = filepath.Join("workflows", s.Environment)
	}

	if !opts.Force {
		if err := s.GitChecker.CheckBeforeSync(); err != nil {
			s.Logger.Error("sync blocked", "error", err)
			return fmt.Errorf("pre-sync git check failed: %w", err)
		}
	}

	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	if err := s.Client.HealthCheck(ctx); err != nil {
		return fmt.Errorf("failed to connect to n8n API: %w", err)
	}

	workflows, err := s.Client.GetWorkflows(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch workflows: %w", err)
	}

	synced := 0
	for _, wf := range workflows {
		fileName := fmt.Sprintf("%s_%s.json", utils.SanitizeFilename(wf.Name), wf.ID)
		path := filepath.Join(opts.OutputDir, fileName)

		wf.SyncMetadata = &workflow.SyncMetadata{
			SyncDate:    time.Now(),
			Environment: s.Environment,
			SyncedBy:    getSyncUser(),
			GitCommit:   os.Getenv("CI_COMMIT_SHA"),
		}

		if err := writeFile(path, wf); err != nil {
			s.Logger.Error("failed to write workflow file", "workflow", wf.Name, "error", err)
			if !opts.Force {
				return fmt.Errorf("failed to write workflow %s: %w", wf.Name, err)
			}
			continue
		}

		synced++
	}

	metadata := map[string]interface{}{
		"lastSync":        time.Now(),
		"environment":     s.Environment,
		"totalWorkflows":  len(workflows),
		"syncedWorkflows": synced,
		"syncedBy":        getSyncUser(),
	}
	metaPath := filepath.Join(opts.OutputDir, "_sync_metadata.json")
	if err := utils.WriteJSONFile(metadata, metaPath); err != nil {
		s.Logger.Warn("failed to write sync metadata", "error", err)
	}

	fmt.Printf("✅ Sync completed: %d workflows synced to %s\n", synced, opts.OutputDir)
	return nil
}

func writeFile(p string, wf *workflow.Workflow) error {
	f, err := os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(wf)
}

func getSyncUser() string {
	if u := os.Getenv("GITLAB_USER_EMAIL"); u != "" {
		return u
	}
	if u := os.Getenv("CI_COMMIT_AUTHOR_EMAIL"); u != "" {
		return u
	}
	if u := getGitUser(); u != "" {
		return u
	}
	if u := os.Getenv("USER"); u != "" {
		return u + "@local"
	}
	return "n8n-ops-user"
}

func getGitUser() string {
	out, err := exec.Command("git", "config", "user.email").Output()
	if err == nil && len(out) > 0 {
		return strings.TrimSpace(string(out))
	}
	return ""
}
