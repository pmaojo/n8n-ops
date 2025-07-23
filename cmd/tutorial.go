package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pmaojo/n8n-ops/internal/ascii"
	"github.com/pmaojo/n8n-ops/internal/utils"
	"github.com/spf13/cobra"
)

var tutorialCmd = &cobra.Command{
	Use:   "tutorial",
	Short: "Interactive step-by-step tutorial for n8n-ops",
	Long: `Start an interactive tutorial that guides you through using n8n-ops.
This tutorial will teach you:
- Basic n8n-ops commands and workflows
- How to sync and deploy workflows
- Best practices for workflow management
- Common troubleshooting steps

The tutorial is hands-on and will let you practice commands in a safe environment.`,
	Run: runTutorial,
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
	if !tutorialSkipIntro {
		showTutorialIntro()
	}

	fmt.Println(ascii.Banner("tutorial"))

	// Check if config exists
	if !configExists() {
		fmt.Printf("%s\n", ascii.ErrorMessage("No configuration found. Let's set up your environment first."))
		fmt.Println("Run 'n8n-ops onboard' to create your configuration.")
		return
	}

	// Main tutorial loop
	runMainTutorial()
}

func showTutorialIntro() {
	utils.ClearTerminalScreen()
	fmt.Print(ascii.TutorialWelcome())
	time.Sleep(1 * time.Second)

	fmt.Println("\n🚀 Welcome to the n8n-ops Interactive Tutorial!")
	fmt.Println("This tutorial will guide you through using n8n-ops step by step.")
	time.Sleep(1 * time.Second)
}

func configExists() bool {
	// Check for .n8n-ops.yaml in current directory or home directory
	_, err := os.Stat(".n8n-ops.yaml")
	if err == nil {
		return true
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false
	}

	_, err = os.Stat(homeDir + "/.n8n-ops.yaml")
	return err == nil
}

func runMainTutorial() {
	topics := []tutorialTopic{
		{
			title:       "Understanding n8n-ops",
			description: "Learn about the core concepts of n8n-ops",
			action:      showConceptsTutorial,
		},
		{
			title:       "Syncing Workflows",
			description: "Learn how to sync workflows between n8n and your local files",
			action:      showSyncTutorial,
		},
		{
			title:       "Deploying Workflows",
			description: "Learn how to deploy workflows to different environments",
			action:      showDeployTutorial,
		},
		{
			title:       "Validating Workflows",
			description: "Learn how to validate your workflow files",
			action:      showValidateTutorial,
		},
		{
			title:       "Git Integration",
			description: "Learn how to use n8n-ops with Git",
			action:      showGitTutorial,
		},
	}

	if tutorialAdvanced {
		advancedTopics := []tutorialTopic{
			{
				title:       "Multi-Environment Setup",
				description: "Learn how to manage multiple environments",
				action:      showMultiEnvTutorial,
			},
			{
				title:       "CI/CD Integration",
				description: "Learn how to integrate with GitLab CI/CD",
				action:      showCICDTutorial,
			},
			{
				title:       "Monitoring & Alerts",
				description: "Learn how to monitor workflows and set up alerts",
				action:      showMonitoringTutorial,
			},
		}
		topics = append(topics, advancedTopics...)
	}

	currentTopic := 0
	for {
		utils.ClearTerminalScreen()
		fmt.Println(ascii.Banner("tutorial"))
		fmt.Printf("\n%s%s📚 TUTORIAL TOPICS%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

		for i, topic := range topics {
			if i == currentTopic {
				fmt.Printf("%s%s> %d. %s%s\n", ascii.Bold, ascii.Green, i+1, topic.title, ascii.Reset)
				fmt.Printf("   %s%s%s\n", ascii.Dim, topic.description, ascii.Reset)
			} else {
				fmt.Printf("  %d. %s\n", i+1, topic.title)
			}
		}

		fmt.Printf("\n%s%sUse arrow keys to navigate, Enter to select, Q to quit%s\n", ascii.Bold, ascii.Cyan, ascii.Reset)

		key := waitForKey()
		switch key {
		case "up":
			if currentTopic > 0 {
				currentTopic--
			}
		case "down":
			if currentTopic < len(topics)-1 {
				currentTopic++
			}
		case "enter":
			utils.ClearTerminalScreen()
			topics[currentTopic].action()
			waitForEnter()
		case "q":
			fmt.Println("\nExiting tutorial. Run 'n8n-ops tutorial' anytime to continue learning!")
			return
		}
	}
}

type tutorialTopic struct {
	title       string
	description string
	action      func()
}

func showConceptsTutorial() {
	fmt.Println(ascii.TutorialHeader("Understanding n8n-ops"))

	fmt.Printf("%s%sCore Concepts:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	concepts := []struct {
		name        string
		description string
	}{
		{
			name:        "Workflow Synchronization",
			description: "n8n-ops syncs workflows between n8n instances and local JSON files, enabling version control and collaboration.",
		},
		{
			name:        "Multi-Environment Management",
			description: "Manage workflows across development, staging, and production environments with proper isolation.",
		},
		{
			name:        "GitOps Workflow",
			description: "Use Git as the source of truth for workflow changes, with automated CI/CD pipelines.",
		},
		{
			name:        "Zero-Downtime Deployments",
			description: "Deploy workflow changes without stopping the n8n service using API-based updates.",
		},
	}

	for _, concept := range concepts {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, concept.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, concept.description)
	}

	fmt.Printf("%s%sWorkflow:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	fmt.Println(ascii.WorkflowDiagram())

	fmt.Printf("\n%s%sKey Commands:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	commands := []struct {
		name        string
		description string
	}{
		{
			name:        "n8n-ops sync",
			description: "Download workflows from n8n to local files",
		},
		{
			name:        "n8n-ops deploy",
			description: "Upload workflows from local files to n8n",
		},
		{
			name:        "n8n-ops validate",
			description: "Check workflow files for errors",
		},
		{
			name:        "n8n-ops status",
			description: "Show workflow status across environments",
		},
	}

	for _, cmd := range commands {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Green, cmd.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, cmd.description)
	}
}

func showSyncTutorial() {
	fmt.Println(ascii.TutorialHeader("Syncing Workflows"))

	fmt.Printf("%s%sWhat is Sync?%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Sync downloads workflows from your n8n instance to local JSON files.\n")
	fmt.Printf("This allows you to version control your workflows and collaborate with others.\n\n")

	fmt.Printf("%s%sBasic Sync Command:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("%s%sn8n-ops sync --env development%s\n\n", ascii.Bold, ascii.Green, ascii.Reset)
	fmt.Printf("This will download all workflows from your development n8n instance\n")
	fmt.Printf("to the ./workflows/development/ directory.\n\n")

	fmt.Printf("%s%sSync Options:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	options := []struct {
		flag        string
		description string
	}{
		{
			flag:        "--force",
			description: "Override local changes without prompting",
		},
		{
			flag:        "--backup",
			description: "Create a backup of local files before syncing",
		},
		{
			flag:        "--dry-run",
			description: "Show what would be synced without making changes",
		},
		{
			flag:        "--include \"pattern\"",
			description: "Only sync workflows matching the pattern",
		},
	}

	for _, opt := range options {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, opt.flag, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, opt.description)
	}

	fmt.Printf("%s%sExample Workflow:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("1. Start development in your n8n instance\n")
	fmt.Printf("2. Create and test workflows in the n8n UI\n")
	fmt.Printf("3. Run: %sn8n-ops sync --env development%s\n", ascii.Green, ascii.Reset)
	fmt.Printf("4. Commit the downloaded workflow files to Git\n")
	fmt.Printf("5. Share with your team via Git\n\n")

	fmt.Printf("%s%sTry It:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Run the following command to sync workflows from your development environment:\n\n")
	fmt.Printf("%s%sn8n-ops sync --env development --dry-run%s\n", ascii.Bold, ascii.Green, ascii.Reset)
}

func showDeployTutorial() {
	fmt.Println(ascii.TutorialHeader("Deploying Workflows"))

	fmt.Printf("%s%sWhat is Deploy?%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Deploy uploads workflows from your local JSON files to an n8n instance.\n")
	fmt.Printf("This allows you to promote workflows between environments and apply version-controlled changes.\n\n")

	fmt.Printf("%s%sBasic Deploy Command:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("%s%sn8n-ops deploy --env staging%s\n\n", ascii.Bold, ascii.Green, ascii.Reset)
	fmt.Printf("This will upload all workflows from the ./workflows/staging/ directory\n")
	fmt.Printf("to your staging n8n instance.\n\n")

	fmt.Printf("%s%sDeploy Options:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	options := []struct {
		flag        string
		description string
	}{
		{
			flag:        "--dry-run",
			description: "Show what would be deployed without making changes",
		},
		{
			flag:        "--force",
			description: "Override remote workflows without prompting",
		},
		{
			flag:        "--auto-rollback",
			description: "Automatically rollback if deployment fails",
		},
		{
			flag:        "<file.json>",
			description: "Deploy a specific workflow file instead of all workflows",
		},
	}

	for _, opt := range options {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, opt.flag, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, opt.description)
	}

	fmt.Printf("%s%sZero-Downtime Deployment:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("n8n-ops uses the n8n API to update workflows without stopping the n8n service.\n")
	fmt.Printf("This means your workflows continue running during deployments.\n\n")

	fmt.Printf("%s%sPromotion Workflow:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("1. Develop and test in development environment\n")
	fmt.Printf("2. Run: %sn8n-ops deploy --env staging --dry-run%s\n", ascii.Green, ascii.Reset)
	fmt.Printf("3. Run: %sn8n-ops deploy --env staging%s\n", ascii.Green, ascii.Reset)
	fmt.Printf("4. Test in staging environment\n")
	fmt.Printf("5. Run: %sn8n-ops deploy --env production%s (requires approval)\n\n", ascii.Green, ascii.Reset)

	fmt.Printf("%s%sTry It:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Run the following command to see what would be deployed to your staging environment:\n\n")
	fmt.Printf("%s%sn8n-ops deploy --env staging --dry-run%s\n", ascii.Bold, ascii.Green, ascii.Reset)
}

func showValidateTutorial() {
	fmt.Println(ascii.TutorialHeader("Validating Workflows"))

	fmt.Printf("%s%sWhat is Validate?%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Validate checks your workflow JSON files for errors and ensures they follow best practices.\n")
	fmt.Printf("This helps catch issues before deploying to production.\n\n")

	fmt.Printf("%s%sBasic Validate Command:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("%s%sn8n-ops validate ./workflows/development/%s\n\n", ascii.Bold, ascii.Green, ascii.Reset)
	fmt.Printf("This will check all workflow files in the development directory.\n\n")

	fmt.Printf("%s%sValidation Checks:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	checks := []struct {
		name        string
		description string
	}{
		{
			name:        "JSON Syntax",
			description: "Ensures the JSON is valid and well-formed",
		},
		{
			name:        "Workflow Structure",
			description: "Checks that the workflow has the required fields",
		},
		{
			name:        "Node Connections",
			description: "Verifies that all nodes are properly connected",
		},
		{
			name:        "Business Rules",
			description: "Applies custom validation rules for your organization",
		},
	}

	for _, check := range checks {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, check.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, check.description)
	}

	fmt.Printf("%s%sValidate Options:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	options := []struct {
		flag        string
		description string
	}{
		{
			flag:        "--strict",
			description: "Apply stricter validation rules",
		},
		{
			flag:        "--fix",
			description: "Attempt to fix minor issues automatically",
		},
		{
			flag:        "--json",
			description: "Output results in JSON format",
		},
	}

	for _, opt := range options {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, opt.flag, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, opt.description)
	}

	fmt.Printf("%s%sCI/CD Integration:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Validation is automatically run in CI/CD pipelines before deployment.\n")
	fmt.Printf("This prevents invalid workflows from being deployed to production.\n\n")

	fmt.Printf("%s%sTry It:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Run the following command to validate your development workflows:\n\n")
	fmt.Printf("%s%sn8n-ops validate ./workflows/development/%s\n", ascii.Bold, ascii.Green, ascii.Reset)
}

func showGitTutorial() {
	fmt.Println(ascii.TutorialHeader("Git Integration"))

	fmt.Printf("%s%sGit Workflow:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("n8n-ops is designed to work with Git for version control and collaboration.\n")
	fmt.Printf("This enables team collaboration, change tracking, and automated deployments.\n\n")

	fmt.Printf("%s%sBranch Strategy:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	branches := []struct {
		name        string
		description string
	}{
		{
			name:        "main",
			description: "Production workflows (protected)",
		},
		{
			name:        "staging",
			description: "Staging workflows (protected)",
		},
		{
			name:        "develop",
			description: "Development workflows",
		},
		{
			name:        "feature/*",
			description: "Feature branches for new workflows",
		},
	}

	for _, branch := range branches {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, branch.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, branch.description)
	}

	fmt.Printf("%s%sBranch Commands:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	commands := []struct {
		name        string
		description string
	}{
		{
			name:        "n8n-ops branch current",
			description: "Show current branch and mapped environment",
		},
		{
			name:        "n8n-ops branch list",
			description: "List all branch-to-environment mappings",
		},
		{
			name:        "n8n-ops branch set feature/auth development",
			description: "Map a branch to an environment",
		},
	}

	for _, cmd := range commands {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Green, cmd.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, cmd.description)
	}

	fmt.Printf("%s%sTypical Git Workflow:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("1. Create a feature branch: %sgit checkout -b feature/payment-workflow%s\n", ascii.Green, ascii.Reset)
	fmt.Printf("2. Sync workflows: %sn8n-ops sync --env development%s\n", ascii.Green, ascii.Reset)
	fmt.Printf("3. Make changes in n8n UI or edit JSON files\n")
	fmt.Printf("4. Validate changes: %sn8n-ops validate ./workflows/development/%s\n", ascii.Green, ascii.Reset)
	fmt.Printf("5. Commit changes: %sgit add workflows/ && git commit -m \"Add payment workflow\"%s\n", ascii.Green, ascii.Reset)
	fmt.Printf("6. Create merge request to staging branch\n")
	fmt.Printf("7. After approval, CI/CD deploys to staging\n")
	fmt.Printf("8. Create merge request to main branch for production\n\n")

	fmt.Printf("%s%sTry It:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Run the following command to see your current branch mapping:\n\n")
	fmt.Printf("%s%sn8n-ops branch current%s\n", ascii.Bold, ascii.Green, ascii.Reset)
}

func showMultiEnvTutorial() {
	fmt.Println(ascii.TutorialHeader("Multi-Environment Setup"))

	fmt.Printf("%s%sEnvironment Management:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("n8n-ops supports multiple environments (development, staging, production).\n")
	fmt.Printf("Each environment has its own configuration and workflow files.\n\n")

	fmt.Printf("%s%sEnvironment Configuration:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("In %s~/.n8n-ops.yaml%s:\n\n", ascii.Cyan, ascii.Reset)
	fmt.Printf("```yaml\n")
	fmt.Printf("environments:\n")
	fmt.Printf("  development:\n")
	fmt.Printf("    url: \"https://n8n-dev.example.com\"\n")
	fmt.Printf("    api_key_env: \"N8N_DEV_API_KEY\"\n")
	fmt.Printf("  staging:\n")
	fmt.Printf("    url: \"https://n8n-staging.example.com\"\n")
	fmt.Printf("    api_key_env: \"N8N_STAGING_API_KEY\"\n")
	fmt.Printf("  production:\n")
	fmt.Printf("    url: \"https://n8n-prod.example.com\"\n")
	fmt.Printf("    api_key_env: \"N8N_PROD_API_KEY\"\n")
	fmt.Printf("```\n\n")

	fmt.Printf("%s%sEnvironment Variables:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("In %s.env%s:\n\n", ascii.Cyan, ascii.Reset)
	fmt.Printf("```\n")
	fmt.Printf("N8N_DEV_API_KEY=n8n_api_your_dev_key_here\n")
	fmt.Printf("N8N_STAGING_API_KEY=n8n_api_your_staging_key_here\n")
	fmt.Printf("N8N_PROD_API_KEY=n8n_api_your_prod_key_here\n")
	fmt.Printf("```\n\n")

	fmt.Printf("%s%sDirectory Structure:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("```\n")
	fmt.Printf("workflows/\n")
	fmt.Printf("├── development/     # Development workflows\n")
	fmt.Printf("├── staging/         # Staging workflows\n")
	fmt.Printf("└── production/      # Production workflows\n")
	fmt.Printf("```\n\n")

	fmt.Printf("%s%sEnvironment Commands:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	commands := []struct {
		name        string
		description string
	}{
		{
			name:        "n8n-ops sync --env development",
			description: "Sync from development environment",
		},
		{
			name:        "n8n-ops deploy --env staging",
			description: "Deploy to staging environment",
		},
		{
			name:        "n8n-ops status --env production",
			description: "Check status of production environment",
		},
	}

	for _, cmd := range commands {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Green, cmd.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, cmd.description)
	}

	fmt.Printf("%s%sPromotion Process:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("1. Develop in development environment\n")
	fmt.Printf("2. Promote to staging for testing\n")
	fmt.Printf("3. Promote to production after approval\n\n")

	fmt.Printf("%s%sTry It:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Run the following command to check your environment configuration:\n\n")
	fmt.Printf("%s%sn8n-ops status%s\n", ascii.Bold, ascii.Green, ascii.Reset)
}

func showCICDTutorial() {
	fmt.Println(ascii.TutorialHeader("CI/CD Integration"))

	fmt.Printf("%s%sGitLab CI/CD Integration:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("n8n-ops integrates with GitLab CI/CD for automated testing and deployment.\n")
	fmt.Printf("This enables continuous integration and delivery of your workflows.\n\n")

	fmt.Printf("%s%sCI/CD Pipeline Stages:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	stages := []struct {
		name        string
		description string
	}{
		{
			name:        "validate",
			description: "Check workflow files for errors",
		},
		{
			name:        "test",
			description: "Run tests on workflows",
		},
		{
			name:        "deploy-dev",
			description: "Deploy to development environment",
		},
		{
			name:        "deploy-staging",
			description: "Deploy to staging environment (manual approval)",
		},
		{
			name:        "deploy-production",
			description: "Deploy to production environment (manual approval)",
		},
	}

	for _, stage := range stages {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, stage.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, stage.description)
	}

	fmt.Printf("%s%sGitLab CI/CD Variables:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Set these in GitLab → Settings → CI/CD → Variables:\n\n")

	variables := []struct {
		name        string
		description string
	}{
		{
			name:        "N8N_DEV_API_KEY",
			description: "API key for development environment",
		},
		{
			name:        "N8N_STAGING_API_KEY",
			description: "API key for staging environment",
		},
		{
			name:        "N8N_PROD_API_KEY",
			description: "API key for production environment",
		},
		{
			name:        "N8N_DEV_URL",
			description: "URL for development environment",
		},
		{
			name:        "N8N_STAGING_URL",
			description: "URL for staging environment",
		},
		{
			name:        "N8N_PROD_URL",
			description: "URL for production environment",
		},
	}

	for _, variable := range variables {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, variable.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, variable.description)
	}

	fmt.Printf("%s%sExample .gitlab-ci.yml:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("```yaml\n")
	fmt.Printf("stages:\n")
	fmt.Printf("  - validate\n")
	fmt.Printf("  - test\n")
	fmt.Printf("  - deploy-dev\n")
	fmt.Printf("  - deploy-staging\n")
	fmt.Printf("  - deploy-production\n\n")

	fmt.Printf("validate-workflows:\n")
	fmt.Printf("  stage: validate\n")
	fmt.Printf("  script:\n")
	fmt.Printf("    - n8n-ops validate ./workflows/\n\n")

	fmt.Printf("deploy-development:\n")
	fmt.Printf("  stage: deploy-dev\n")
	fmt.Printf("  script:\n")
	fmt.Printf("    - n8n-ops deploy --env development\n")
	fmt.Printf("  rules:\n")
	fmt.Printf("    - if: '$CI_COMMIT_BRANCH == \"develop\"'\n\n")

	fmt.Printf("deploy-staging:\n")
	fmt.Printf("  stage: deploy-staging\n")
	fmt.Printf("  script:\n")
	fmt.Printf("    - n8n-ops deploy --env staging\n")
	fmt.Printf("  when: manual\n")
	fmt.Printf("  rules:\n")
	fmt.Printf("    - if: '$CI_COMMIT_BRANCH == \"staging\"'\n\n")

	fmt.Printf("deploy-production:\n")
	fmt.Printf("  stage: deploy-production\n")
	fmt.Printf("  script:\n")
	fmt.Printf("    - n8n-ops deploy --env production\n")
	fmt.Printf("  when: manual\n")
	fmt.Printf("  rules:\n")
	fmt.Printf("    - if: '$CI_COMMIT_BRANCH == \"main\"'\n")
	fmt.Printf("```\n\n")

	fmt.Printf("%s%sAutomated Workflow:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("1. Push to develop branch → Automatic deployment to development\n")
	fmt.Printf("2. Merge to staging branch → Manual approval → Deploy to staging\n")
	fmt.Printf("3. Merge to main branch → Manual approval → Deploy to production\n\n")
}

func showMonitoringTutorial() {
	fmt.Println(ascii.TutorialHeader("Monitoring & Alerts"))

	fmt.Printf("%s%sWorkflow Monitoring:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("n8n-ops can monitor your workflows for failures and automatically create issues.\n")
	fmt.Printf("This helps you detect and fix problems quickly.\n\n")

	fmt.Printf("%s%sMonitoring Commands:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	commands := []struct {
		name        string
		description string
	}{
		{
			name:        "n8n-ops monitor --env production",
			description: "Start monitoring production workflows",
		},
		{
			name:        "n8n-ops status --env production",
			description: "Check status of production workflows",
		},
		{
			name:        "n8n-ops check --env production",
			description: "Check for changes in production workflows",
		},
	}

	for _, cmd := range commands {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Green, cmd.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, cmd.description)
	}

	fmt.Printf("%s%sMonitoring Options:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	options := []struct {
		flag        string
		description string
	}{
		{
			flag:        "--interval 5m",
			description: "Check every 5 minutes",
		},
		{
			flag:        "--create-issues",
			description: "Automatically create GitLab issues for failures",
		},
		{
			flag:        "--notify slack",
			description: "Send notifications to Slack",
		},
		{
			flag:        "--daemon",
			description: "Run in background as a daemon",
		},
	}

	for _, opt := range options {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, opt.flag, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, opt.description)
	}

	fmt.Printf("%s%sGitLab Issue Creation:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("When a workflow fails, n8n-ops can automatically create a GitLab issue:\n\n")
	fmt.Printf("```\n")
	fmt.Printf("Title: [FAILED] Payment Processing Workflow\n")
	fmt.Printf("Description: Workflow failed at 2023-07-23 14:32:45\n")
	fmt.Printf("Error: API connection timeout\n")
	fmt.Printf("Environment: production\n")
	fmt.Printf("Last successful run: 2023-07-23 13:30:12\n")
	fmt.Printf("```\n\n")

	fmt.Printf("%s%sNotification Channels:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)

	channels := []struct {
		name        string
		description string
	}{
		{
			name:        "GitLab Issues",
			description: "Create issues for workflow failures",
		},
		{
			name:        "Slack",
			description: "Send notifications to Slack channels",
		},
		{
			name:        "Microsoft Teams",
			description: "Send notifications to Teams channels",
		},
		{
			name:        "Email",
			description: "Send email notifications",
		},
	}

	for _, channel := range channels {
		fmt.Printf("%s%s%s%s\n", ascii.Bold, ascii.Cyan, channel.name, ascii.Reset)
		fmt.Printf("%s%s\n\n", ascii.Dim, channel.description)
	}

	fmt.Printf("%s%sTry It:%s\n\n", ascii.Bold, ascii.Yellow, ascii.Reset)
	fmt.Printf("Run the following command to check the status of your workflows:\n\n")
	fmt.Printf("%s%sn8n-ops status --env development%s\n", ascii.Bold, ascii.Green, ascii.Reset)
}

// Helper functions for the tutorial

func waitForKey() string {
	// This is a simplified version - in a real implementation,
	// you would use a library like "github.com/eiannone/keyboard"
	// to capture arrow keys and other special keys
	fmt.Print("> ")
	var input string
	fmt.Scanln(&input)

	input = strings.ToLower(input)
	if input == "q" || input == "quit" || input == "exit" {
		return "q"
	} else if input == "up" || input == "u" {
		return "up"
	} else if input == "down" || input == "d" {
		return "down"
	} else {
		return "enter"
	}
}

func waitForEnter() {
	fmt.Printf("\n%s%sPress Enter to continue...%s", ascii.Bold, ascii.Cyan, ascii.Reset)
	fmt.Scanln()
}
