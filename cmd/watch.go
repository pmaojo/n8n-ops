package cmd

import (
	"context"
	"fmt"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cliFrom(cmd)
		if cli == nil {
			return fmt.Errorf("CLI not initialized")
		}
		return runWatch(cmd, args, cli)
	},
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

func runWatch(cmd *cobra.Command, args []string, cli *CLI) error {
	logEntry := cli.Logger.WithFields(logrus.Fields{
		"command": "watch",
		"env":     cli.Environment,
	})

	n8nClient, cm, err := cliutils.SetupClient(cli.Environment, demoMode, logEntry)
	if err != nil {
		return err
	}

	gitChecker := git.NewGitStatusChecker(".", nil, logEntry)
	syncSvc := isync.NewService(n8nClient, cm, gitChecker, logEntry, cli.Environment)
	svc := iwatch.NewService(n8nClient, cm, gitChecker, syncSvc, logEntry, cli.Environment)

	return svc.Watch(context.Background(), iwatch.Options{
		Interval:   watchInterval,
		AutoCommit: autoCommit,
		AutoSync:   autoSync,
	})
}
