package cmd

import (
	"fmt"
	"time"

	"github.com/pmaojo/n8n-ops/internal/ascii"
	"github.com/pmaojo/n8n-ops/internal/tutorial"
	"github.com/pmaojo/n8n-ops/internal/utils"
	"github.com/spf13/cobra"
)

var quickstartCmd = &cobra.Command{
	Use:   "quickstart",
	Short: "Interactive quickstart guide for new users",
	Long: `Start an interactive quickstart guide to get up and running with n8n-ops quickly.
This guide will walk you through:
- Basic configuration
- Essential commands
- Common workflows
- Best practices

Perfect for first-time users who want to learn the basics quickly.`,
	Run: runQuickstart,
}

var (
	quickstartSkipIntro bool
)

func init() {
	rootCmd.AddCommand(quickstartCmd)

	quickstartCmd.Flags().BoolVar(&quickstartSkipIntro, "skip-intro", false, "Skip the introduction animation")
}

func runQuickstart(cmd *cobra.Command, args []string) {
	if !quickstartSkipIntro {
		showQuickstartIntro()
	}

	fmt.Println(ascii.Banner("quickstart"))

	// Check if config exists
	if !tutorial.ConfigExists() {
		fmt.Printf("%s\n", ascii.ErrorMessage("No configuration found. Let's set up your environment first."))
		fmt.Println("\nWould you like to run the onboarding wizard now? (y/n)")
		var response string
		fmt.Scanln(&response)
		if response == "y" || response == "Y" {
			// Execute onboard command
			onboardCmd.Run(cmd, args)
			return
		} else {
			fmt.Println("\nYou can run 'n8n-ops onboard' later to set up your environment.")
		}
	}

	// Show quickstart guide
	showQuickstartGuide()
}

func showQuickstartIntro() {
	utils.ClearTerminalScreen()
	fmt.Print(ascii.MatrixEffect())
	time.Sleep(500 * time.Millisecond)

	fmt.Print(ascii.N8nLogo())
	time.Sleep(1 * time.Second)

	fmt.Printf("\n\n%s%s🚀 QUICKSTART GUIDE%s\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("%sLet's get you up and running in 5 minutes!%s\n\n", ascii.Cyan, ascii.Reset)
	time.Sleep(1 * time.Second)
}

func showQuickstartGuide() {
	steps := []struct {
		title       string
		description string
		command     string
		action      func()
	}{
		{
			title:       "Configuration",
			description: "Set up your n8n-ops configuration",
			command:     "n8n-ops onboard",
			action:      showConfigStep,
		},
		{
			title:       "Syncing Workflows",
			description: "Download workflows from n8n",
			command:     "n8n-ops sync --env development",
			action:      showSyncStep,
		},
		{
			title:       "Validating Workflows",
			description: "Check workflows for errors",
			command:     "n8n-ops validate ./workflows/development/",
			action:      showValidateStep,
		},
		{
			title:       "Deploying Workflows",
			description: "Upload workflows to n8n",
			command:     "n8n-ops deploy --env development",
			action:      showDeployStep,
		},
		{
			title:       "Git Integration",
			description: "Use n8n-ops with Git",
			command:     "git add workflows/ && git commit -m \"Update workflows\"",
			action:      showGitStep,
		},
	}

	for i, step := range steps {
		utils.ClearTerminalScreen()
		fmt.Println(ascii.Banner("quickstart"))
		fmt.Printf("\n%s%s%sSTEP %d/%d: %s%s\n\n", ascii.Bold, ascii.Yellow, ascii.Underline, i+1, len(steps), step.title, ascii.Reset)

		fmt.Printf("%s%s%s\n\n", ascii.Bold, step.description, ascii.Reset)
		fmt.Printf("Command: %s%s%s\n\n", ascii.Green, step.command, ascii.Reset)

		step.action()

		if i < len(steps)-1 {
			fmt.Printf("\n%s%sPress Enter to continue to the next step...%s", ascii.Bold, ascii.Cyan, ascii.Reset)
			fmt.Scanln()
		}
	}

	// Show completion message
	utils.ClearTerminalScreen()
	fmt.Println(ascii.Banner("quickstart"))
	fmt.Printf("\n%s%s✨ QUICKSTART COMPLETE! ✨%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("%sYou now know the basics of using n8n-ops!%s\n\n", ascii.Cyan, ascii.Reset)

	fmt.Printf("%s%sNext Steps:%s\n\n", ascii.Bold, ascii.Green, ascii.Reset)
	fmt.Printf("1. Try the interactive tutorial: %sn8n-ops tutorial%s\n", ascii.Yellow, ascii.Reset)
	fmt.Printf("2. Read the documentation: %sDEVELOPMENT.md%s\n", ascii.Yellow, ascii.Reset)
	fmt.Printf("3. Explore more commands: %sn8n-ops --help%s\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("%s%sHappy workflow automation!%s\n", ascii.Bold, ascii.Green, ascii.Reset)
}

func showConfigStep() {
	fmt.Printf("%s%sConfiguration File:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("Create %s~/.n8n-ops.yaml%s with your environment settings:\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("```yaml\n")
	fmt.Printf("environments:\n")
	fmt.Printf("  development:\n")
	fmt.Printf("    url: \"http://localhost:5678\"\n")
	fmt.Printf("    api_key_env: \"N8N_DEV_API_KEY\"\n")
	fmt.Printf("  staging:\n")
	fmt.Printf("    url: \"https://n8n-staging.example.com\"\n")
	fmt.Printf("    api_key_env: \"N8N_STAGING_API_KEY\"\n")
	fmt.Printf("  production:\n")
	fmt.Printf("    url: \"https://n8n-prod.example.com\"\n")
	fmt.Printf("    api_key_env: \"N8N_PROD_API_KEY\"\n")
	fmt.Printf("```\n\n")

	fmt.Printf("%s%sEnvironment Variables:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("Create %s.env%s with your API keys:\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("```\n")
	fmt.Printf("N8N_DEV_API_KEY=n8n_api_your_dev_key_here\n")
	fmt.Printf("N8N_STAGING_API_KEY=n8n_api_your_staging_key_here\n")
	fmt.Printf("N8N_PROD_API_KEY=n8n_api_your_prod_key_here\n")
	fmt.Printf("```\n\n")

	fmt.Printf("%s%sTip:%s Use the onboarding wizard to create these files automatically:\n", ascii.Bold, ascii.Green, ascii.Reset)
	fmt.Printf("%sn8n-ops onboard%s\n", ascii.Yellow, ascii.Reset)
}

func showSyncStep() {
	fmt.Printf("%s%sSync Command:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("The sync command downloads workflows from your n8n instance to local files.\n\n")

	fmt.Printf("%s%sBasic Usage:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("%sn8n-ops sync --env development%s\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("%s%sOptions:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("- %s--force%s: Override local changes without prompting\n", ascii.Green, ascii.Reset)
	fmt.Printf("- %s--backup%s: Create a backup of local files before syncing\n", ascii.Green, ascii.Reset)
	fmt.Printf("- %s--dry-run%s: Show what would be synced without making changes\n\n", ascii.Green, ascii.Reset)

	fmt.Printf("%s%sOutput:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("Workflows are saved to %s./workflows/development/%s directory.\n", ascii.Yellow, ascii.Reset)
	fmt.Printf("Each workflow is saved as a separate JSON file.\n\n")

	fmt.Printf("%s%sExample:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("```\n")
	fmt.Printf("$ n8n-ops sync --env development\n")
	fmt.Printf("Connecting to n8n instance at http://localhost:5678...\n")
	fmt.Printf("✅ Connected successfully!\n")
	fmt.Printf("Syncing workflows...\n")
	fmt.Printf("✅ Downloaded 5 workflows to ./workflows/development/\n")
	fmt.Printf("```\n")
}

func showValidateStep() {
	fmt.Printf("%s%sValidate Command:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("The validate command checks your workflow files for errors.\n\n")

	fmt.Printf("%s%sBasic Usage:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("%sn8n-ops validate ./workflows/development/%s\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("%s%sWhat It Checks:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("- JSON syntax\n")
	fmt.Printf("- Workflow structure\n")
	fmt.Printf("- Node connections\n")
	fmt.Printf("- Business rules\n\n")

	fmt.Printf("%s%sOptions:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("- %s--strict%s: Apply stricter validation rules\n", ascii.Green, ascii.Reset)
	fmt.Printf("- %s--fix%s: Attempt to fix minor issues automatically\n\n", ascii.Green, ascii.Reset)

	fmt.Printf("%s%sExample:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("```\n")
	fmt.Printf("$ n8n-ops validate ./workflows/development/\n")
	fmt.Printf("Validating workflows in ./workflows/development/...\n")
	fmt.Printf("✅ payment-workflow.json: Valid\n")
	fmt.Printf("✅ email-notification.json: Valid\n")
	fmt.Printf("❌ customer-onboarding.json: Error - Missing node connections\n")
	fmt.Printf("Summary: 2 valid, 1 invalid\n")
	fmt.Printf("```\n")
}

func showDeployStep() {
	fmt.Printf("%s%sDeploy Command:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("The deploy command uploads workflows from local files to your n8n instance.\n\n")

	fmt.Printf("%s%sBasic Usage:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("%sn8n-ops deploy --env development%s\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("%s%sOptions:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("- %s--dry-run%s: Show what would be deployed without making changes\n", ascii.Green, ascii.Reset)
	fmt.Printf("- %s--force%s: Override remote workflows without prompting\n", ascii.Green, ascii.Reset)
	fmt.Printf("- %s--auto-rollback%s: Automatically rollback if deployment fails\n\n", ascii.Green, ascii.Reset)

	fmt.Printf("%s%sZero-Downtime Deployment:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("n8n-ops uses the n8n API to update workflows without stopping the n8n service.\n")
	fmt.Printf("This means your workflows continue running during deployments.\n\n")

	fmt.Printf("%s%sExample:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("```\n")
	fmt.Printf("$ n8n-ops deploy --env staging --dry-run\n")
	fmt.Printf("Connecting to n8n instance at https://n8n-staging.example.com...\n")
	fmt.Printf("✅ Connected successfully!\n")
	fmt.Printf("Dry run: The following workflows would be deployed:\n")
	fmt.Printf("- payment-workflow.json\n")
	fmt.Printf("- email-notification.json\n")
	fmt.Printf("Run without --dry-run to perform the actual deployment.\n")
	fmt.Printf("```\n")
}

func showGitStep() {
	fmt.Printf("%s%sGit Integration:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("n8n-ops is designed to work with Git for version control and collaboration.\n\n")

	fmt.Printf("%s%sTypical Git Workflow:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("1. Create a feature branch:\n")
	fmt.Printf("   %sgit checkout -b feature/payment-workflow%s\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("2. Sync workflows from n8n:\n")
	fmt.Printf("   %sn8n-ops sync --env development%s\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("3. Make changes in n8n UI or edit JSON files\n\n")

	fmt.Printf("4. Validate changes:\n")
	fmt.Printf("   %sn8n-ops validate ./workflows/development/%s\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("5. Commit changes:\n")
	fmt.Printf("   %sgit add workflows/ && git commit -m \"Add payment workflow\"%s\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("6. Push changes:\n")
	fmt.Printf("   %sgit push origin feature/payment-workflow%s\n\n", ascii.Yellow, ascii.Reset)

	fmt.Printf("7. Create merge request to staging branch\n\n")

	fmt.Printf("8. After approval, CI/CD deploys to staging\n\n")

	fmt.Printf("%s%sBranch Commands:%s\n\n", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Printf("- %sn8n-ops branch current%s: Show current branch and mapped environment\n", ascii.Yellow, ascii.Reset)
	fmt.Printf("- %sn8n-ops branch list%s: List all branch-to-environment mappings\n", ascii.Yellow, ascii.Reset)
}
