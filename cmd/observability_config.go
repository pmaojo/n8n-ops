package cmd

import "github.com/spf13/viper"

// ObservabilityConfig groups configuration for Sentry and Grafana.
type ObservabilityConfig struct {
	SentryDSN     string
	GrafanaURL    string
	GrafanaAPIKey string
	GrafanaOrgID  int
}

// loadObservabilityConfig retrieves observability settings from the provided Viper instance.
func loadObservabilityConfig(v *viper.Viper) ObservabilityConfig {
	return ObservabilityConfig{
		SentryDSN:     v.GetString("sentry_dsn"),
		GrafanaURL:    v.GetString("grafana_url"),
		GrafanaAPIKey: v.GetString("grafana_api_key"),
		GrafanaOrgID:  v.GetInt("grafana_org_id"),
	}
}
