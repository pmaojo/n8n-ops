package cmd

import (
	"context"

	"github.com/pmaojo/n8n-ops/internal/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// CLI aggregates shared resources for commands.
type CLI struct {
	Environment string
	Logger      *logrus.Logger
	Config      *config.Config
}

// cliKey is used as context key for storing CLI struct.
type cliKey struct{}

func withCLI(ctx context.Context, c *CLI) context.Context {
	return context.WithValue(ctx, cliKey{}, c)
}

func cliFrom(cmd *cobra.Command) *CLI {
	if v := cmd.Context().Value(cliKey{}); v != nil {
		if c, ok := v.(*CLI); ok {
			return c
		}
	}
	return nil
}
