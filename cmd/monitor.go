package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pmaojo/n8n-ops/internal/cliutils"
	"github.com/pmaojo/n8n-ops/internal/i18n"
	"github.com/pmaojo/n8n-ops/internal/issues"
	"github.com/pmaojo/n8n-ops/internal/monitoring"
	"github.com/pmaojo/n8n-ops/internal/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor",
	Short: "Monitor workflows and automatically create issues for failures",
	Long: `Monitor mode continuously watches n8n workflow executions and automatically
creates GitLab issues when failures are detected. It also updates issues when
workflows recover.

This command integrates with your GitLab project to provide automated
issue tracking for workflow reliability.

Examples:
  n8n-ops monitor --env production          # Monitor production workflows
  n8n-ops monitor --check-interval 30s     # Check every 30 seconds
  n8n-ops monitor --failure-threshold 5    # Create issue after 5 failures`,
	RunE: runMonitor,
}

var (
	checkInterval    time.Duration
	failureThreshold int
	gitlabProjectID  string
	gitlabURL        string
)

func init() {
	rootCmd.AddCommand(monitorCmd)

	monitorCmd.Flags().DurationVar(&checkInterval, "check-interval", 1*time.Minute, "interval between failure checks")
	monitorCmd.Flags().IntVar(&failureThreshold, "failure-threshold", 3, "consecutive failures before creating issue")
	monitorCmd.Flags().StringVar(&gitlabProjectID, "gitlab-project", "", "GitLab project ID for issues")
	monitorCmd.Flags().StringVar(&gitlabURL, "gitlab-url", "https://gitlab.com", "GitLab instance URL")
}

func runMonitor(cmd *cobra.Command, args []string) error {
	logger := utils.NewLogger()
	utils.SetLogLevel(logger, "info")
	if verbose {
		utils.SetLogLevel(logger, "debug")
	}

	logEntry := logger.WithFields(logrus.Fields{
		"command": "monitor",
		"env":     environment,
	})

	logEntry.Info("Starting workflow monitoring")

	i18n.PrintfKey("monitor_starting", environment)
	i18n.PrintfKey("monitor_check_interval", checkInterval)
	i18n.PrintfKey("monitor_failure_threshold", failureThreshold)


	ctx := context.Background()
	n8nClient, err := getN8nClient(ctx, environment, demoMode)
	if err != nil {
		return err
	}

	issueManager, projectID, err := setupIssueManager(demoMode, gitlabURL, gitlabProjectID)
	if err != nil {
		return err
	}

	if demoMode {
		fmt.Println("📋 Using mock issue manager (demo mode)")
	} else {
		fmt.Printf("📋 Connected to GitLab project: %s\n", projectID)
	}

	detector := monitoring.NewFailureDetector(n8nClient, issueManager, logger)
	detector.SetCheckInterval(checkInterval)
	detector.SetRetryThreshold(failureThreshold)

	return startMonitoring(detector, logEntry, nil)
}

var monitorClientFactory = func(url, apiKey string) (client.Client, error) {
	return client.New(url, apiKey, nil)
}

// getN8nClient creates and verifies the n8n client.
func getN8nClient(ctx context.Context, env string, demo bool) (client.Client, error) {
	var n8nURL, apiKey string
	if demo {
		n8nURL = "http://localhost:3001"
		apiKey = "n8n_api_mock_development"
	} else {
		cm := credentials.NewCredentialManager(env)
		var err error
		n8nURL, apiKey, err = cm.GetN8nCredentials()
		if err != nil {
			return nil, fmt.Errorf("failed to load credentials: %w", err)
		}
		if n8nURL == "" || apiKey == "" {
			return nil, fmt.Errorf("n8n credentials not configured for %s environment. Set N8N_%s_URL and N8N_%s_API_KEY or use --demo",
				env, strings.ToUpper(env), strings.ToUpper(env))
		}
	}

	n8nClient, err := monitorClientFactory(n8nURL, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create n8n client: %w", err)

	}

	if err := n8nClient.HealthCheck(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to n8n: %w", err)
	}

	return n8nClient, nil
}

// setupIssueManager configures the issue manager implementation.
func setupIssueManager(demo bool, url, projectID string) (issues.IssueManager, string, error) {
	if projectID == "" {
		projectID = os.Getenv("GITLAB_PROJECT_ID")
		if projectID == "" {
			projectID = os.Getenv("GITLAB_PROJECT")
		}
		if projectID == "" && !demo {
			return nil, "", fmt.Errorf("gitlab-project or GITLAB_PROJECT_ID is required")
		}
	}

	token := os.Getenv("GITLAB_TOKEN")
	if token == "" && !demo {
		return nil, "", fmt.Errorf("GITLAB_TOKEN environment variable is required")
	}

	if demo {
		return NewMockIssueManager(), projectID, nil
	}

	return issues.NewGitLabIssueManager(url, projectID, token), projectID, nil
}

// startMonitoring launches the detector and waits for termination signals.
func startMonitoring(detector interface{ Start(context.Context) error }, logEntry *logrus.Entry, sigChan chan os.Signal) error {
	if sigChan == nil {
		sigChan = make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := detector.Start(ctx); err != nil && err != context.Canceled {
			logEntry.WithError(err).Error("Monitoring failed")
		}
	}()

	fmt.Printf("✅ Monitoring started. Creating issues for workflow failures...\n")
	i18n.PrintlnKey("monitor_stop")

	sig := <-sigChan
	fmt.Printf("\n🛑 Received signal %v, stopping monitoring...\n", sig)
	cancel()
	return nil
}

// MockIssueManager is a demo implementation for testing
type MockIssueManager struct{}

func NewMockIssueManager() *MockIssueManager {
	return &MockIssueManager{}
}

func (mim *MockIssueManager) CreateWorkflowFailureIssue(ctx context.Context, failure *issues.WorkflowFailure) (*issues.Issue, error) {
	fmt.Printf("🔍 DEMO: Would create issue for workflow failure\n")
	fmt.Printf("   Workflow: %s (%s)\n", failure.WorkflowName, failure.WorkflowID)
	fmt.Printf("   Environment: %s\n", failure.Environment)
	fmt.Printf("   Error: %s\n", failure.ErrorMessage)
	fmt.Printf("   Retry Count: %d\n", failure.RetryCount)

	// Return mock issue
	return &issues.Issue{
		ID:        12345,
		IID:       1,
		Title:     fmt.Sprintf("🚨 Workflow Failure: %s (%s)", failure.WorkflowName, failure.Environment),
		State:     "opened",
		WebURL:    "https://gitlab.example.com/project/issues/1",
		CreatedAt: time.Now(),
		Labels:    []string{"workflow-failure", "automated"},
	}, nil
}

func (mim *MockIssueManager) UpdateIssueWithRecovery(ctx context.Context, issueID int, recovery *issues.RecoveryInfo) error {
	fmt.Printf("🔍 DEMO: Would update issue #%d with recovery info\n", issueID)
	fmt.Printf("   Recovery Type: %s\n", recovery.RecoveryType)
	fmt.Printf("   Recovered At: %s\n", recovery.RecoveredAt.Format(time.RFC3339))
	if recovery.Notes != "" {
		fmt.Printf("   Notes: %s\n", recovery.Notes)
	}
	return nil
}

func (mim *MockIssueManager) GetOpenFailureIssues(ctx context.Context, workflowID string) ([]*issues.Issue, error) {
	fmt.Printf("🔍 DEMO: Would check for open issues for workflow: %s\n", workflowID)
	// Return empty slice for demo
	return []*issues.Issue{}, nil
}
