package cmd

import (
	"fmt"

	"github.com/pmaojo/n8n-ops/internal/observability"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	sentryDSN     string
	grafanaURL    string
	grafanaAPIKey string
	grafanaOrgID  int
	enableSentry  bool
	enableGrafana bool
)

// observabilityCmd represents the observability command
var observabilityCmd = &cobra.Command{
	Use:   "observability",
	Short: "Configure and manage observability integrations (Sentry, Grafana)",
	Long: `Configure and manage observability integrations for enterprise monitoring.

Supports:
- Sentry: Error tracking and performance monitoring
- Grafana: Metrics dashboards and alerting
- Real-time workflow monitoring
- Performance analytics

Examples:
  n8n-ops observability setup --sentry --grafana
  n8n-ops observability test-connection
  n8n-ops observability create-dashboard`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// setupCmd configures observability integrations
var setupObservabilityCmd = &cobra.Command{
	Use:   "setup",
	Short: "Setup observability integrations",
	Long:  `Initialize and configure Sentry and Grafana integrations for monitoring.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔧 Setting up observability integrations...")

		cfg := loadObservabilityConfig(viper.GetViper())

		// Setup Sentry
		if enableSentry {
			if cfg.SentryDSN == "" {
				return fmt.Errorf("sentry DSN is required")
			}

			sentryConfig := observability.SentryConfig{
				DSN:         cfg.SentryDSN,
				Environment: environment,
				Release:     "1.0.0",
				SampleRate:  1.0,
			}

			sentry := observability.NewSentryIntegration(sentryConfig, logger)
			if err := sentry.Initialize(); err != nil {
				return fmt.Errorf("failed to setup Sentry: %w", err)
			}

			fmt.Println("✅ Sentry integration configured")
		}

		// Setup Grafana
		if enableGrafana {
			if cfg.GrafanaURL == "" || cfg.GrafanaAPIKey == "" {
				return fmt.Errorf("grafana URL and API key are required")
			}

			grafanaConfig := observability.GrafanaConfig{
				URL:       cfg.GrafanaURL,
				APIKey:    cfg.GrafanaAPIKey,
				OrgID:     cfg.GrafanaOrgID,
				Dashboard: "n8n-ops-monitoring",
			}

			grafana := observability.NewGrafanaIntegration(grafanaConfig, logger)
			if err := grafana.Initialize(); err != nil {
				return fmt.Errorf("failed to setup Grafana: %w", err)
			}

			fmt.Println("✅ Grafana integration configured")
		}

		fmt.Println("🎉 Observability setup complete!")
		return nil
	},
}

// testConnectionCmd tests observability connections
var testConnectionCmd = &cobra.Command{
	Use:   "test-connection",
	Short: "Test observability service connections",
	Long:  `Test connectivity to Sentry and Grafana services.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🧪 Testing observability connections...")

		success := true

		cfg := loadObservabilityConfig(viper.GetViper())

		// Test Sentry
		if enableSentry && cfg.SentryDSN != "" {
			sentryConfig := observability.SentryConfig{
				DSN:         cfg.SentryDSN,
				Environment: "test",
				Release:     "test",
				SampleRate:  0.1,
			}

			sentry := observability.NewSentryIntegration(sentryConfig, logger)
			if err := sentry.Initialize(); err != nil {
				fmt.Printf("❌ Sentry connection failed: %v\n", err)
				success = false
			} else {
				fmt.Println("✅ Sentry connection successful")
				sentry.Close()
			}
		}

		// Test Grafana
		if enableGrafana && cfg.GrafanaURL != "" && cfg.GrafanaAPIKey != "" {
			grafanaConfig := observability.GrafanaConfig{
				URL:    cfg.GrafanaURL,
				APIKey: cfg.GrafanaAPIKey,
				OrgID:  cfg.GrafanaOrgID,
			}

			grafana := observability.NewGrafanaIntegration(grafanaConfig, logger)
			if err := grafana.Initialize(); err != nil {
				fmt.Printf("❌ Grafana connection failed: %v\n", err)
				success = false
			} else {
				fmt.Println("✅ Grafana connection successful")
				grafana.Close()
			}
		}

		if success {
			fmt.Println("🎉 All connections successful!")
		} else {
			return fmt.Errorf("some connections failed")
		}

		return nil
	},
}

// createDashboardCmd creates Grafana dashboard
var createDashboardCmd = &cobra.Command{
	Use:   "create-dashboard",
	Short: "Create default Grafana dashboard",
	Long:  `Create a pre-configured Grafana dashboard for n8n-ops monitoring.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadObservabilityConfig(viper.GetViper())
		if cfg.GrafanaURL == "" || cfg.GrafanaAPIKey == "" {
			return fmt.Errorf("grafana URL and API key are required")
		}

		fmt.Println("📊 Creating Grafana dashboard...")

		grafanaConfig := observability.GrafanaConfig{
			URL:       cfg.GrafanaURL,
			APIKey:    cfg.GrafanaAPIKey,
			OrgID:     cfg.GrafanaOrgID,
			Dashboard: "n8n-ops-monitoring",
		}

		grafana := observability.NewGrafanaIntegration(grafanaConfig, logger)
		if err := grafana.Initialize(); err != nil {
			return fmt.Errorf("failed to initialize Grafana: %w", err)
		}

		if err := grafana.CreateDashboard(cmd.Context()); err != nil {
			return fmt.Errorf("failed to create dashboard: %w", err)
		}

		fmt.Println("✅ Grafana dashboard created successfully!")
		fmt.Printf("📈 Access at: %s/d/n8n-ops-monitoring\n", cfg.GrafanaURL)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(observabilityCmd)

	// Add subcommands
	observabilityCmd.AddCommand(setupObservabilityCmd)
	observabilityCmd.AddCommand(testConnectionCmd)
	observabilityCmd.AddCommand(createDashboardCmd)

	// Sentry flags
	observabilityCmd.PersistentFlags().BoolVar(&enableSentry, "sentry", false, "enable Sentry integration")
	observabilityCmd.PersistentFlags().StringVar(&sentryDSN, "sentry-dsn", "", "Sentry DSN")

	// Grafana flags
	observabilityCmd.PersistentFlags().BoolVar(&enableGrafana, "grafana", false, "enable Grafana integration")
	observabilityCmd.PersistentFlags().StringVar(&grafanaURL, "grafana-url", "", "Grafana URL")
	observabilityCmd.PersistentFlags().StringVar(&grafanaAPIKey, "grafana-api-key", "", "Grafana API key")
	observabilityCmd.PersistentFlags().IntVar(&grafanaOrgID, "grafana-org-id", 1, "Grafana organization ID")

	// Bind flags and environment variables
	viper.BindPFlag("sentry_dsn", observabilityCmd.PersistentFlags().Lookup("sentry-dsn"))
	viper.BindPFlag("grafana_url", observabilityCmd.PersistentFlags().Lookup("grafana-url"))
	viper.BindPFlag("grafana_api_key", observabilityCmd.PersistentFlags().Lookup("grafana-api-key"))
	viper.BindPFlag("grafana_org_id", observabilityCmd.PersistentFlags().Lookup("grafana-org-id"))

	viper.BindEnv("sentry_dsn", "SENTRY_DSN")
	viper.BindEnv("grafana_url", "GRAFANA_URL")
	viper.BindEnv("grafana_api_key", "GRAFANA_API_KEY")
	viper.BindEnv("grafana_org_id", "GRAFANA_ORG_ID")
}
