package cmd

import (
	"context"
	"time"

	"github.com/pmaojo/n8n-ops/internal/cliutils"
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
	logEntry := logger.WithFields(logrus.Fields{
		"command": "watch",
		"env":     environment,
	})

	n8nClient, cm, err := cliutils.SetupClient(environment, demoMode)
	if err != nil {
		return err
	}

	gitChecker := git.NewGitStatusChecker(".", nil)
	syncSvc := isync.NewService(n8nClient, cm, gitChecker, logEntry, environment)
	svc := iwatch.NewService(n8nClient, cm, gitChecker, syncSvc, logEntry, environment)

	return svc.Watch(context.Background(), iwatch.Options{
		Interval:   watchInterval,
		AutoCommit: autoCommit,
		AutoSync:   autoSync,
	})
}
