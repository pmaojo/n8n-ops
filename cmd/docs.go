package cmd

import (
	"context"
	"fmt"

	"github.com/pmaojo/n8n-ops/internal/docgen"
	"github.com/spf13/cobra"
)

var docsCmd = &cobra.Command{
	Use:   "docs",
	Short: "Manage documentation",
}

var docsGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate CLI reference documentation",
	RunE: func(cmd *cobra.Command, args []string) error {
		cli := cliFrom(cmd)
		if cli == nil {
			return fmt.Errorf("CLI not initialized")
		}
		generator := docgen.NewCobraGenerator(rootCmd)
		return generator.Generate(context.Background(), docsFormat, docsOutput)
	},
}

var (
	docsFormat string
	docsOutput string
)

func init() {
	rootCmd.AddCommand(docsCmd)
	docsCmd.AddCommand(docsGenerateCmd)

	docsGenerateCmd.Flags().StringVar(&docsFormat, "format", "md", "output format: md or html")
	docsGenerateCmd.Flags().StringVar(&docsOutput, "output", "docs", "output directory")
}
