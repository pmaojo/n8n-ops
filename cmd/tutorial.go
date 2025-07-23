package cmd

import (
	"github.com/pmaojo/n8n-ops/internal/tutorial"
	"github.com/spf13/cobra"
)

var tutorialCmd = &cobra.Command{
	Use:   "tutorial",
	Short: "Interactive step-by-step tutorial for n8n-ops",
	Long:  `Start an interactive tutorial that guides you through using n8n-ops.`,
	Run:   runTutorial,
}

var (
	tutorialSkipIntro bool
	tutorialAdvanced  bool
)

func init() {
	rootCmd.AddCommand(tutorialCmd)

	tutorialCmd.Flags().BoolVar(&tutorialSkipIntro, "skip-intro", false, "Skip the introduction animation")
	tutorialCmd.Flags().BoolVar(&tutorialAdvanced, "advanced", false, "Show advanced tutorial topics")
}

func runTutorial(cmd *cobra.Command, args []string) {
	tutorial.Run(tutorialSkipIntro, tutorialAdvanced)
}
