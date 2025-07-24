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
	"github.com/pmaojo/n8n-ops/internal/utils"
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
	RunE: runTUI,
}

var tuiRefresh time.Duration

func init() {
	rootCmd.AddCommand(tuiCmd)

	tuiCmd.Flags().DurationVar(&tuiRefresh, "refresh", 3*time.Second, "dashboard refresh interval")
}

func runTUI(cmd *cobra.Command, args []string) error {
	n8nClient, _, err := cliutils.SetupClient(environment, demoMode)
	if err != nil {
		var missing cliutils.MissingCredentialError
		if errors.As(err, &missing) {
			return fmt.Errorf("N8N_%s_API_KEY is required", strings.ToUpper(environment))
		}
	}
	logger := utils.NewLogger()
	ctx := context.Background()
	var uiImpl ui.DashboardUI
	uiImpl = bubbleui.NewDashboard(n8nClient, tuiRefresh, logger)
	return uiImpl.Run(ctx)
}
