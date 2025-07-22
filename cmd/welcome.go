package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/n8n-workflows/cli/internal/ascii"
)

var welcomeCmd = &cobra.Command{
	Use:   "welcome",
	Short: "Display the futuristic welcome screen",
	Long:  "Shows the complete n8n CLI welcome screen with Matrix effects and ASCII art",
	Run:   runWelcome,
}

func init() {
	rootCmd.AddCommand(welcomeCmd)
}

func runWelcome(cmd *cobra.Command, args []string) {
	fmt.Print(ascii.WelcomeScreen())
}