package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pmaojo/n8n-ops/internal/app"
	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/pmaojo/n8n-ops/internal/git"
	isync "github.com/pmaojo/n8n-ops/internal/sync"
	iwatch "github.com/pmaojo/n8n-ops/internal/watch"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Monitor n8n workflows for changes and auto-sync",
	Long: `Watch mode monitors your n8n instance for workflow changes and automatically
syncs them to your local Git repository. This enables real-time collaboration
and backup of workflow modifications made in the n8n UI.

Examples:
  n8n-ops watch --env development     # Watch development n8n instance
  n8n-ops watch --env production      # Watch production instance  
  n8n-ops watch --interval 30s        # Check every 30 seconds`,
	RunE: runWatch,
}

var (
	watchInterval time.Duration
	autoCommit    bool
	autoSync      bool
)

func init() {
	rootCmd.AddCommand(watchCmd)

	watchCmd.Flags().DurationVar(&watchInterval, "interval", 10*time.Second, "polling interval (e.g., 10s, 1m)")
	watchCmd.Flags().BoolVar(&autoCommit, "auto-commit", false, "automatically commit changes to git")
	watchCmd.Flags().BoolVar(&autoSync, "auto-sync", true, "automatically sync detected changes")
}

func runWatch(cmd *cobra.Command, args []string) error {
	cfg := app.FromContext(cmd.Context())
	if cfg == nil {
		return fmt.Errorf("configuration not found in context")
	}

	logEntry := cfg.Logger.WithFields(logrus.Fields{
		"command": "watch",
		"env":     cfg.Environment,
	})

	cm := credentials.NewCredentialManager(cfg.Environment)
	n8nURL, apiKey, err := cm.GetN8nCredentials()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}
	if n8nURL == "" {
		n8nURL = "http://localhost:5678"
	}
	if apiKey == "" {
		return fmt.Errorf("N8N_%s_API_KEY environment variable not set", strings.ToUpper(cfg.Environment))
	}

	n8nClient, err := client.New(n8nURL, apiKey, nil)
	if err != nil {
		return fmt.Errorf("failed to create n8n client: %w", err)
	}

	gitChecker := git.NewGitStatusChecker(".", nil)
	syncSvc := isync.NewService(n8nClient, cm, gitChecker, logEntry, cfg.Environment)
	svc := iwatch.NewService(n8nClient, cm, gitChecker, syncSvc, logEntry, cfg.Environment)

	return svc.Watch(context.Background(), iwatch.Options{
		Interval:   watchInterval,
		AutoCommit: autoCommit,
		AutoSync:   autoSync,
	})
}
