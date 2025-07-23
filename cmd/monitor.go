package cmd

import (
        "context"
        "fmt"
        "os"
        "os/signal"
        "strings"
        "syscall"
        "time"

        "github.com/n8n-workflows/n8n-ops/internal/client"
        "github.com/n8n-workflows/n8n-ops/internal/issues"
        "github.com/n8n-workflows/n8n-ops/internal/monitoring"
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
        checkInterval     time.Duration
        failureThreshold  int
        gitlabProjectID   string
        gitlabURL         string
)

func init() {
        rootCmd.AddCommand(monitorCmd)

        monitorCmd.Flags().DurationVar(&checkInterval, "check-interval", 1*time.Minute, "interval between failure checks")
        monitorCmd.Flags().IntVar(&failureThreshold, "failure-threshold", 3, "consecutive failures before creating issue")
        monitorCmd.Flags().StringVar(&gitlabProjectID, "gitlab-project", "", "GitLab project ID for issues")
        monitorCmd.Flags().StringVar(&gitlabURL, "gitlab-url", "https://gitlab.com", "GitLab instance URL")
}

func runMonitor(cmd *cobra.Command, args []string) error {
        logger := logrus.New()
        logger.SetLevel(logrus.InfoLevel)
        if verbose {
                logger.SetLevel(logrus.DebugLevel)
        }

        logEntry := logger.WithFields(logrus.Fields{
                "command": "monitor",
                "env":     environment,
        })

        logEntry.Info("Starting workflow monitoring")

        if language == "es" {
                fmt.Printf("👁️ Iniciando monitoreo de workflows - %s\n", environment)
                fmt.Printf("⏰ Intervalo de verificación: %s\n", checkInterval)
                fmt.Printf("🎯 Umbral de fallos: %d fallos consecutivos\n", failureThreshold)
        } else {
                fmt.Printf("👁️ Starting workflow monitoring - %s environment\n", environment)
                fmt.Printf("⏰ Check interval: %s\n", checkInterval)
                fmt.Printf("🎯 Failure threshold: %d consecutive failures\n", failureThreshold)
        }

        // Create n8n client using unified credential system
        var n8nURL, apiKey string
        if demoMode {
                n8nURL = "http://localhost:3001"
                apiKey = "n8n_api_mock_development"
        } else {
                // Get credentials using cascading environment variable approach
                envSuffix := strings.ToUpper(environment)
                n8nURL = os.Getenv(fmt.Sprintf("N8N_%s_URL", envSuffix))
                apiKey = os.Getenv(fmt.Sprintf("N8N_%s_API_KEY", envSuffix))
                
                // Fallback to short forms for common environments
                if n8nURL == "" || apiKey == "" {
                        switch environment {
                        case "development":
                                if n8nURL == "" {
                                        n8nURL = os.Getenv("N8N_DEV_URL")
                                }
                                if apiKey == "" {
                                        apiKey = os.Getenv("N8N_DEV_API_KEY")
                                }
                        case "staging":
                                if n8nURL == "" {
                                        n8nURL = os.Getenv("N8N_STAGING_URL")
                                }
                                if apiKey == "" {
                                        apiKey = os.Getenv("N8N_STAGING_API_KEY")
                                }
                        case "production":
                                if n8nURL == "" {
                                        n8nURL = os.Getenv("N8N_PROD_URL")
                                }
                                if apiKey == "" {
                                        apiKey = os.Getenv("N8N_PROD_API_KEY")
                                }
                        }
                }
                
                // Final fallback to generic variables
                if n8nURL == "" {
                        n8nURL = os.Getenv("N8N_URL")
                }
                if apiKey == "" {
                        apiKey = os.Getenv("N8N_API_KEY")
                }
                if n8nURL == "" || apiKey == "" {
                        return fmt.Errorf("n8n credentials not configured for %s environment. Set N8N_%s_URL and N8N_%s_API_KEY or use --demo", 
                                environment, strings.ToUpper(environment), strings.ToUpper(environment))
                }
        }

        n8nClient, err := client.New(n8nURL, apiKey, nil)
        if err != nil {
                return fmt.Errorf("failed to create n8n client: %w", err)
        }

        // Test connection
        ctx := context.Background()
        if err := n8nClient.HealthCheck(ctx); err != nil {
                return fmt.Errorf("failed to connect to n8n: %w", err)
        }

        // Setup GitLab issue manager
        if gitlabProjectID == "" {
                gitlabProjectID = os.Getenv("GITLAB_PROJECT_ID")
                if gitlabProjectID == "" {
                        gitlabProjectID = os.Getenv("GITLAB_PROJECT")
                }
                if gitlabProjectID == "" && !demoMode {
                        return fmt.Errorf("gitlab-project or GITLAB_PROJECT_ID is required")
                }
        }

        gitlabToken := os.Getenv("GITLAB_TOKEN")
        if gitlabToken == "" && !demoMode {
                return fmt.Errorf("GITLAB_TOKEN environment variable is required")
        }

        var issueManager issues.IssueManager
        if demoMode {
                issueManager = NewMockIssueManager()
                fmt.Println("📋 Using mock issue manager (demo mode)")
        } else {
                issueManager = issues.NewGitLabIssueManager(gitlabURL, gitlabProjectID, gitlabToken)
                fmt.Printf("📋 Connected to GitLab project: %s\n", gitlabProjectID)
        }

        // Create failure detector
        detector := monitoring.NewFailureDetector(n8nClient, issueManager, logger)
        detector.SetCheckInterval(checkInterval)
        detector.SetRetryThreshold(failureThreshold)

        // Setup signal handling for graceful shutdown
        sigChan := make(chan os.Signal, 1)
        signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

        // Create context for cancellation
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()

        // Start monitoring in background
        go func() {
                if err := detector.Start(ctx); err != nil && err != context.Canceled {
                        logEntry.WithError(err).Error("Monitoring failed")
                }
        }()

        fmt.Printf("✅ Monitoring started. Creating issues for workflow failures...\n")
        if language == "es" {
                fmt.Printf("⏹️ Presiona Ctrl+C para detener el monitoreo\n\n")
        } else {
                fmt.Printf("⏹️ Press Ctrl+C to stop monitoring\n\n")
        }

        // Wait for shutdown signal
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
                ID:    12345,
                IID:   1,
                Title: fmt.Sprintf("🚨 Workflow Failure: %s (%s)", failure.WorkflowName, failure.Environment),
                State: "opened",
                WebURL: "https://gitlab.example.com/project/issues/1",
                CreatedAt: time.Now(),
                Labels: []string{"workflow-failure", "automated"},
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