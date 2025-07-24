package cmd

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/pmaojo/n8n-ops/internal/bubbleui"
	"github.com/pmaojo/n8n-ops/internal/cliutils"
	"github.com/pmaojo/n8n-ops/internal/ui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interactive terminal dashboard",
	Long: `Interactive terminal dashboard.

Controls:
  Up/Down - move selection
  Q       - quit

Flags:
  --refresh duration   dashboard refresh interval (default 3s)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cliFrom(cmd)
		if cli == nil {
			return fmt.Errorf("CLI not initialized")
		}
		return runTUI(cmd, args, cli)
	},
}

var tuiRefresh time.Duration

func init() {
	rootCmd.AddCommand(tuiCmd)

	tuiCmd.Flags().DurationVar(&tuiRefresh, "refresh", 3*time.Second, "dashboard refresh interval")
}

func runTUI(cmd *cobra.Command, args []string, cli *CLI) error {
	n8nClient, _, err := cliutils.SetupClient(cli.Environment, demoMode)
	if err != nil {
		var missing cliutils.MissingCredentialError
		if errors.As(err, &missing) {
			return fmt.Errorf("N8N_%s_API_KEY is required", strings.ToUpper(cli.Environment))
		}
	}
	ctx := context.Background()
	var uiImpl ui.DashboardUI
	uiImpl = bubbleui.NewDashboard(n8nClient, tuiRefresh, cli.Logger)
	return uiImpl.Run(ctx)
}
