package cmd

import (
	"context"
	"fmt"

	"github.com/pmaojo/n8n-ops/internal/observability"
	"github.com/pmaojo/n8n-ops/internal/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	serverPort   int
	serverEnable bool
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run REST server for status and metrics",
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cliFrom(cmd)
		if cli == nil {
			return fmt.Errorf("CLI not initialized")
		}

		cfg := cli.Config.GetServerConfig()
		if serverPort != 0 {
			cfg.Port = serverPort
		}
		if serverEnable {
			cfg.Enabled = true
		}

		srvCfg := server.Config{Enabled: cfg.Enabled, Port: cfg.Port}

		obsCfg := loadObservabilityConfig(viper.GetViper())
		var grafana *observability.GrafanaIntegration
		if obsCfg.GrafanaURL != "" && obsCfg.GrafanaAPIKey != "" {
			grafana = observability.NewGrafanaIntegration(observability.GrafanaConfig{
				URL:    obsCfg.GrafanaURL,
				APIKey: obsCfg.GrafanaAPIKey,
				OrgID:  obsCfg.GrafanaOrgID,
			}, cli.Logger)
			// best effort initialize
			_ = grafana.Initialize(cmd.Context())
		}

		srv := server.NewService(srvCfg, cli.Logger, grafana)
		return srv.Start(context.Background())
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVar(&serverEnable, "enable", false, "enable REST server")
	serverCmd.Flags().IntVar(&serverPort, "port", 8081, "server port")

	viper.BindPFlag("server.enabled", serverCmd.Flags().Lookup("enable"))
	viper.BindPFlag("server.port", serverCmd.Flags().Lookup("port"))
}
