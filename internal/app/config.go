package app

import (
	"context"

	"github.com/pmaojo/n8n-ops/internal/config"
	"github.com/sirupsen/logrus"
)

// Config holds runtime settings shared across commands.
type Config struct {
	Logger      *logrus.Logger
	Environment string
	DemoMode    bool
	Verbose     bool
	Settings    *config.Config
}

// New creates a new Config instance.
func New(logger *logrus.Logger, cfg *config.Config, env string, demo bool, verbose bool) *Config {
	return &Config{Logger: logger, Settings: cfg, Environment: env, DemoMode: demo, Verbose: verbose}
}

// contextKey is the type used for storing Config in a context.Context.
type contextKey struct{}

var cfgKey contextKey

// WithContext returns a new context carrying cfg.
func WithContext(ctx context.Context, cfg *Config) context.Context {
	return context.WithValue(ctx, cfgKey, cfg)
}

// FromContext retrieves a Config from ctx if available.
func FromContext(ctx context.Context) *Config {
	if ctx == nil {
		return nil
	}
	if cfg, ok := ctx.Value(cfgKey).(*Config); ok {
		return cfg
	}
	return nil
}
