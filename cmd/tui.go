package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pmaojo/n8n-ops/internal/client"
	"github.com/pmaojo/n8n-ops/internal/credentials"
	"github.com/pmaojo/n8n-ops/internal/termui"
	"github.com/spf13/cobra"
)

var tuiCmd = &cobra.Command{
	Use:   "tui",
	Short: "Interactive terminal dashboard using termui",
	RunE:  runTUI,
}

func init() {
	rootCmd.AddCommand(tuiCmd)
}

func runTUI(cmd *cobra.Command, args []string) error {
	var n8nURL, apiKey string
	if demoMode {
		n8nURL = "http://localhost:3001"
		apiKey = "n8n_api_mock_development"
	} else {
		cm := credentials.NewCredentialManager(environment)
		var err error
		n8nURL, apiKey, err = cm.GetN8nCredentials()
		if err != nil {
			return fmt.Errorf("failed to load credentials: %w", err)
		}
		if apiKey == "" {
			return fmt.Errorf("N8N_%s_API_KEY is required", strings.ToUpper(environment))
		}
		if n8nURL == "" {
			n8nURL = "http://localhost:5678"
		}
	}

	var c client.Client
	var err error
	if n8nURL != "" && apiKey != "" {
		c, err = client.New(n8nURL, apiKey, nil)
		if err != nil {
			c = nil
		}
	}

	ctx := context.Background()
	db := termui.NewDashboard(c, 3*time.Second)
	return db.Run(ctx)
}
