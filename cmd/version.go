package cmd

import (
	"fmt"

	"github.com/pmaojo/n8n-ops/internal/ascii"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Long:  "Display version, git commit, and build information for n8n-ops",
	Run:   runVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion(cmd *cobra.Command, args []string) {
	fmt.Print(ascii.SmallLogo())
	fmt.Printf("\n%sn8n-ops version %s%s\n", ascii.Bold, Version, ascii.Reset)
	fmt.Printf("%sGit commit:%s %s\n", ascii.Dim, ascii.Reset, GitCommit)
	fmt.Printf("%sBuild date:%s %s\n", ascii.Dim, ascii.Reset, BuildDate)
	fmt.Printf("%sGo version:%s %s\n", ascii.Dim, ascii.Reset, "go1.19+")
	fmt.Printf("%sPlatform:%s %s/%s\n", ascii.Dim, ascii.Reset, "linux", "amd64")
}
