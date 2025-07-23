package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/n8n-workflows/n8n-ops/internal/ascii"
	"github.com/n8n-workflows/n8n-ops/internal/utils"
	"github.com/spf13/cobra"
)

var onboardCmd = &cobra.Command{
	Use:   "onboard",
	Short: "Interactive onboarding wizard for new users",
	Long: `Start an interactive onboarding process to set up n8n-ops quickly.
This wizard will guide you through:
- Creating a project structure
- Configuring your n8n environments
- Setting up API keys
- Testing connections
- Understanding basic commands

Perfect for first-time users or when setting up a new project.`,
	Run: runOnboard,
}

var (
	onboardSkipIntro bool
	onboardQuick     bool
)

func init() {
	rootCmd.AddCommand(onboardCmd)

	onboardCmd.Flags().BoolVar(&onboardSkipIntro, "skip-intro", false, "Skip the introduction animation")
	onboardCmd.Flags().BoolVar(&onboardQuick, "quick", false, "Quick setup with minimal prompts")
}

func runOnboard(cmd *cobra.Command, args []string) {
	// Load .env file if it exists
	utils.LoadEnvFile()
	
	if !onboardSkipIntro {
		showOnboardingIntro()
	}

	fmt.Println(ascii.Banner("onboarding"))

	// Step 1: Project setup
	projectDir := promptProjectDirectory()
	if projectDir == "" {
		return
	}

	// Step 2: Environment configuration
	envConfig := promptEnvironmentConfig()
	if envConfig == nil {
		return
	}

	// Step 3: Create project structure
	if err := createProjectStructure(projectDir, envConfig); err != nil {
		fmt.Printf("%s\n", ascii.ErrorMessage(fmt.Sprintf("Failed to create project: %v", err)))
		return
	}

	// Step 4: API key setup
	setupAPIKeys(projectDir, envConfig)

	// Step 5: Test connection
	testConnection(envConfig)

	// Step 6: Show next steps
	showNextSteps(projectDir)
}

func showOnboardingIntro() {
	fmt.Print(ascii.WelcomeScreen())
	time.Sleep(1 * time.Second)

	fmt.Println("\n🚀 Welcome to the n8n-ops Onboarding Wizard!")
	fmt.Println("Let's get you set up in just a few minutes.")
	time.Sleep(1 * time.Second)
}

func promptProjectDirectory() string {
	if onboardQuick {
		return "."
	}

	fmt.Println("\n" + ascii.CommandHelp("PROJECT SETUP"))
	fmt.Println("First, let's decide where to create your n8n workflow project.")
	fmt.Println("1. Current directory (.)")
	fmt.Println("2. New subdirectory (enter name)")
	fmt.Println("3. Custom path (enter full path)")

	choice := promptInput("Select an option (1-3)", "1")

	var projectDir string
	switch choice {
	case "1":
		projectDir = "."
	case "2":
		name := promptInput("Enter subdirectory name", "n8n-workflows")
		projectDir = name
	case "3":
		path := promptInput("Enter full path", "")
		projectDir = path
	default:
		projectDir = "."
	}

	// Check if directory exists and is not empty
	if projectDir != "." {
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			fmt.Printf("%s\n", ascii.ErrorMessage(fmt.Sprintf("Failed to create directory: %v", err)))
			return ""
		}
	}

	// Check if directory is not empty and confirm overwrite
	entries, err := os.ReadDir(projectDir)
	if err == nil && len(entries) > 0 {
		confirm := promptInput("Directory not empty. Continue anyway? (y/n)", "n")
		if strings.ToLower(confirm) != "y" {
			fmt.Println("Onboarding cancelled.")
			return ""
		}
	}

	return projectDir
}

type EnvironmentConfig struct {
	Name   string
	URL    string
	APIKey string
}

func promptEnvironmentConfig() []EnvironmentConfig {
	if onboardQuick {
		return []EnvironmentConfig{
			{Name: "development", URL: "http://localhost:5678", APIKey: ""},
		}
	}

	fmt.Println("\n" + ascii.CommandHelp("ENVIRONMENT SETUP"))
	fmt.Println("Now, let's configure your n8n environments.")
	fmt.Println("You can set up multiple environments (development, staging, production).")

	var environments []EnvironmentConfig

	// Always add development environment
	fmt.Println("\n📋 Development Environment")
	devURL := promptInput("Development n8n URL", "http://localhost:5678")
	environments = append(environments, EnvironmentConfig{
		Name: "development",
		URL:  devURL,
	})

	// Ask about staging
	addStaging := promptInput("Add staging environment? (y/n)", "n")
	if strings.ToLower(addStaging) == "y" {
		stagingURL := promptInput("Staging n8n URL", "https://n8n-staging.example.com")
		environments = append(environments, EnvironmentConfig{
			Name: "staging",
			URL:  stagingURL,
		})
	}

	// Ask about production
	addProduction := promptInput("Add production environment? (y/n)", "n")
	if strings.ToLower(addProduction) == "y" {
		prodURL := promptInput("Production n8n URL", "https://n8n-prod.example.com")
		environments = append(environments, EnvironmentConfig{
			Name: "production",
			URL:  prodURL,
		})
	}

	return environments
}

func createProjectStructure(projectDir string, envConfig []EnvironmentConfig) error {
	fmt.Println("\n" + ascii.CommandHelp("CREATING PROJECT"))
	fmt.Println(ascii.LoadingSpinner("Creating project structure..."))

	// Create directory structure
	dirs := []string{
		"workflows",
		"scripts",
		"docs",
		"tests",
		"config",
	}

	// Add environment-specific directories
	for _, env := range envConfig {
		dirs = append(dirs, filepath.Join("workflows", env.Name))
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create configuration files
	if err := createOnboardingConfigFiles(projectDir, envConfig); err != nil {
		return fmt.Errorf("failed to create configuration files: %w", err)
	}

	// Create documentation files
	if err := createOnboardingDocFiles(projectDir); err != nil {
		return fmt.Errorf("failed to create documentation files: %w", err)
	}

	// Create git ignore file
	if err := createGitIgnoreFile(projectDir); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	fmt.Printf("%s\n", ascii.SuccessMessage("Project structure created successfully!"))
	return nil
}

func createOnboardingConfigFiles(projectDir string, envConfig []EnvironmentConfig) error {
	// Create main config file
	configContent := "# n8n CLI Configuration\n# Generated by n8n-ops onboarding wizard\n\n# Environment configurations\nenvironments:\n"

	// Add each environment
	for _, env := range envConfig {
		configContent += fmt.Sprintf("  %s:\n", env.Name)
		configContent += fmt.Sprintf("    url: %s\n", env.URL)
		configContent += fmt.Sprintf("    api_key_env: N8N_%s_API_KEY\n", strings.ToUpper(env.Name))
	}

	// Add default settings
	configContent += `
# Default settings
defaults:
  environment: development
  validation:
    strict: false
  sync:
    auto_backup: true
  deploy:
    dry_run: false
    skip_validation: false

# Logging configuration
logging:
  level: info
  format: json
  file: logs/n8n-ops.log
`

	configPath := filepath.Join(projectDir, ".n8n-ops.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		return err
	}

	// Create .env.example file
	envContent := "# n8n CLI Environment Variables\n# Copy this file to .env and fill in your actual values\n\n"

	for _, env := range envConfig {
		envVarName := fmt.Sprintf("N8N_%s_API_KEY", strings.ToUpper(env.Name))
		envContent += fmt.Sprintf("# %s environment\n", strings.Title(env.Name))
		envContent += fmt.Sprintf("%s=your_%s_api_key_here\n", envVarName, env.Name)
		envContent += fmt.Sprintf("N8N_%s_URL=%s\n\n", strings.ToUpper(env.Name), env.URL)
	}

	envContent += `# GitLab CI/CD variables (set in GitLab project settings)
# GITLAB_TOKEN=your_gitlab_token_here
# GITLAB_PROJECT_ID=12345
`

	envPath := filepath.Join(projectDir, ".env.example")
	if err := os.WriteFile(envPath, []byte(envContent), 0644); err != nil {
		return err
	}

	// Create .env file with instructions but no actual credentials
	envPath = filepath.Join(projectDir, ".env")
	envInstructions := "# n8n CLI Environment Variables\n"
	envInstructions += "# IMPORTANT: Fill in your actual values below\n\n"
	envInstructions += "# Copy values from .env.example and replace with your actual credentials\n"
	envInstructions += "# DO NOT commit this file to git - it contains sensitive API keys!\n\n"
	
	// Add placeholders for environment variables
	for _, env := range envConfig {
		envVarName := fmt.Sprintf("N8N_%s_API_KEY", strings.ToUpper(env.Name))
		shortName := strings.ToUpper(strings.Replace(env.Name, "development", "DEV", 1))
		shortName = strings.ToUpper(strings.Replace(shortName, "production", "PROD", 1))
		shortVarName := fmt.Sprintf("N8N_%s_API_KEY", shortName)
		
		envInstructions += fmt.Sprintf("# %s environment\n", strings.Title(env.Name))
		envInstructions += fmt.Sprintf("# Use either %s or %s\n", envVarName, shortVarName)
		envInstructions += fmt.Sprintf("%s=\n", envVarName)
		
		urlVarName := fmt.Sprintf("N8N_%s_URL", strings.ToUpper(env.Name))
		shortUrlVarName := fmt.Sprintf("N8N_%s_URL", shortName)
		envInstructions += fmt.Sprintf("# Use either %s or %s\n", urlVarName, shortUrlVarName)
		envInstructions += fmt.Sprintf("%s=%s\n\n", urlVarName, env.URL)
	}
	
	return os.WriteFile(envPath, []byte(envInstructions), 0644)
}

func createOnboardingDocFiles(projectDir string) error {
	readmeContent := `# n8n Workflow Project

This project contains n8n workflows managed with the n8n-ops tool for collaborative development.

## Quick Start

1. Configure your environments in ` + "`.n8n-ops.yaml`" + `
2. Set up environment variables in ` + "`.env`" + `
3. Sync workflows: ` + "`n8n-ops sync --env development`" + `
4. Make changes and deploy: ` + "`n8n-ops deploy --env development`" + `

## Commands

- ` + "`n8n-ops sync`" + ` - Sync workflows from n8n instance
- ` + "`n8n-ops deploy`" + ` - Deploy workflows to n8n instance
- ` + "`n8n-ops validate`" + ` - Validate workflow files
- ` + "`n8n-ops status`" + ` - Check workflow status

## Directory Structure

- ` + "`workflows/`" + ` - Environment-specific workflows
- ` + "`docs/`" + ` - Documentation
- ` + "`scripts/`" + ` - Custom scripts
- ` + "`tests/`" + ` - Tests
- ` + "`config/`" + ` - Configuration files

## Need Help?

Run ` + "`n8n-ops --help`" + ` for command information or check the documentation.
`

	readmePath := filepath.Join(projectDir, "README.md")
	return os.WriteFile(readmePath, []byte(readmeContent), 0644)
}

func setupAPIKeys(projectDir string, envConfig []EnvironmentConfig) {
	fmt.Println("\n" + ascii.CommandHelp("API KEY SETUP"))
	fmt.Println("Now, let's set up your n8n API keys.")
	fmt.Println("You'll need to create API keys in your n8n instances.")

	// Show instructions for getting API keys
	fmt.Println("\n📝 How to get your n8n API keys:")
	fmt.Println("1. Log in to your n8n instance")
	fmt.Println("2. Go to Settings → API Keys")
	fmt.Println("3. Click 'Create API Key'")
	fmt.Println("4. Give it a name (e.g., 'n8n-ops-development')")
	fmt.Println("5. Copy the generated key (starts with 'n8n_api_')")

	if onboardQuick {
		fmt.Println("\nYou need to set environment variables for your API keys.")
		showEnvironmentVariableInstructions(envConfig)
		return
	}

	// Check if environment variables are already set
	fmt.Println("\nChecking for existing environment variables...")
	
	allSet := true
	for _, env := range envConfig {
		// Check both naming conventions
		fullName := fmt.Sprintf("N8N_%s_API_KEY", strings.ToUpper(env.Name))
		shortName := ""
		
		if env.Name == "development" {
			shortName = "N8N_DEV_API_KEY"
		} else if env.Name == "production" {
			shortName = "N8N_PROD_API_KEY"
		} else if env.Name == "staging" {
			shortName = "N8N_STAGING_API_KEY"
		}
		
		apiKey := os.Getenv(fullName)
		if apiKey == "" && shortName != "" {
			apiKey = os.Getenv(shortName)
		}
		
		if apiKey == "" {
			fmt.Printf("❌ %s environment API key not set\n", strings.Title(env.Name))
			allSet = false
		} else {
			fmt.Printf("✅ %s environment API key found\n", strings.Title(env.Name))
		}
	}
	
	if allSet {
		fmt.Printf("%s\n", ascii.SuccessMessage("All API keys are already set in environment variables!"))
		return
	}
	
	// Show instructions for setting environment variables
	fmt.Println("\n📝 You need to set the following environment variables:")
	showEnvironmentVariableInstructions(envConfig)
}

// Helper function to show environment variable instructions
func showEnvironmentVariableInstructions(envConfig []EnvironmentConfig) {
	fmt.Println("\nAdd these to your shell profile (~/.bashrc, ~/.zshrc, etc.):")
	fmt.Println("\n```bash")
	
	for _, env := range envConfig {
		fullName := fmt.Sprintf("N8N_%s_API_KEY", strings.ToUpper(env.Name))
		shortName := ""
		
		if env.Name == "development" {
			shortName = "N8N_DEV_API_KEY"
		} else if env.Name == "production" {
			shortName = "N8N_PROD_API_KEY"
		} else if env.Name == "staging" {
			shortName = "N8N_STAGING_API_KEY"
		}
		
		if shortName != "" {
			fmt.Printf("# %s environment - use either %s or %s\n", 
				strings.Title(env.Name), fullName, shortName)
			fmt.Printf("export %s=\"your_%s_api_key_here\"\n\n", 
				shortName, env.Name)
		} else {
			fmt.Printf("# %s environment\n", strings.Title(env.Name))
			fmt.Printf("export %s=\"your_%s_api_key_here\"\n\n", 
				fullName, env.Name)
		}
		
		// Also show URL variables
		fullUrlName := fmt.Sprintf("N8N_%s_URL", strings.ToUpper(env.Name))
		shortUrlName := ""
		
		if env.Name == "development" {
			shortUrlName = "N8N_DEV_URL"
		} else if env.Name == "production" {
			shortUrlName = "N8N_PROD_URL"
		} else if env.Name == "staging" {
			shortUrlName = "N8N_STAGING_URL"
		}
		
		if shortUrlName != "" {
			fmt.Printf("export %s=\"%s\"\n\n", shortUrlName, env.URL)
		} else {
			fmt.Printf("export %s=\"%s\"\n\n", fullUrlName, env.URL)
		}
	}
	
	fmt.Println("```")
	fmt.Println("\nAfter adding these, run `source ~/.bashrc` or restart your terminal.")
}

func testConnection(envConfig []EnvironmentConfig) {
	if onboardQuick {
		return
	}

	fmt.Println("\n" + ascii.CommandHelp("CONNECTION TEST"))
	fmt.Println("Let's test your connection to n8n.")

	testNow := promptInput("Would you like to test the connection now? (y/n)", "y")
	if strings.ToLower(testNow) != "y" {
		return
	}

	// Test connection to development environment
	for _, env := range envConfig {
		if env.Name == "development" {
			fmt.Println(ascii.LoadingSpinner(fmt.Sprintf("Testing connection to %s environment...", env.Name)))
			
			// Try to get API key from environment variable - check both naming conventions
			apiKeyEnvVar := fmt.Sprintf("N8N_%s_API_KEY", strings.ToUpper(env.Name))
			apiKey := os.Getenv(apiKeyEnvVar)
			
			// Also check for N8N_DEV_API_KEY which is commonly used
			if apiKey == "" && env.Name == "development" {
				apiKey = os.Getenv("N8N_DEV_API_KEY")
				if apiKey != "" {
					apiKeyEnvVar = "N8N_DEV_API_KEY"
				}
			}
			
			if apiKey == "" {
				fmt.Printf("%s\n", ascii.ErrorMessage(fmt.Sprintf("API key not found in environment variable %s", apiKeyEnvVar)))
				fmt.Printf("Please set your API key in the .env file or environment variables.\n")
				return
			}
			
			// Get URL from environment variable
			urlEnvVar := fmt.Sprintf("N8N_%s_URL", strings.ToUpper(env.Name))
			url := os.Getenv(urlEnvVar)
			
			// Also check for N8N_DEV_URL which is commonly used
			if url == "" && env.Name == "development" {
				url = os.Getenv("N8N_DEV_URL")
			}
			
			// Use default if not set
			if url == "" {
				url = "http://localhost:5678"
			}
			
			fmt.Printf("Using n8n URL: %s\n", url)
			
			// Test the actual connection
			success, workflowCount, err := testRealConnection(url, apiKey)
			if err != nil {
				fmt.Printf("%s\n", ascii.ErrorMessage(fmt.Sprintf("Connection failed: %v", err)))
				fmt.Printf("Please check:\n")
				fmt.Printf("1. Your n8n instance is running at %s\n", url)
				fmt.Printf("2. Your API key is correct\n")
				fmt.Printf("3. Your n8n instance is accessible\n")
				return
			}
			
			if success {
				fmt.Printf("%s\n", ascii.SuccessMessage(fmt.Sprintf("Connection successful! Found %d workflows.", workflowCount)))
			} else {
				fmt.Printf("%s\n", ascii.ErrorMessage("Connection test failed"))
			}
			break
		}
	}
}

func showNextSteps(projectDir string) {
	fmt.Println("\n" + ascii.CommandHelp("NEXT STEPS"))
	fmt.Printf("%s\n", ascii.SuccessMessage("Onboarding complete! You're ready to use n8n-ops."))

	fmt.Println("\n📋 Here's what to do next:")
	fmt.Println("1. Sync workflows from your n8n instance:")
	fmt.Printf("   cd %s && n8n-ops sync --env development\n", projectDir)
	fmt.Println("\n2. Make changes to workflows in n8n or edit the JSON files")
	fmt.Println("\n3. Deploy your changes back to n8n:")
	fmt.Println("   n8n-ops deploy --env development")
	fmt.Println("\n4. Explore more commands:")
	fmt.Println("   n8n-ops --help")

	fmt.Println("\n📚 Documentation:")
	fmt.Println("- Quick Start Guide: QUICK_START.md")
	fmt.Println("- Development Guide: DEVELOPMENT.md")

	fmt.Println("\n🎉 Happy workflow automation!")
}

// testRealConnection tests the actual connection to n8n
func testRealConnection(url, apiKey string) (bool, int, error) {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	
	// Test connection by getting workflows
	req, err := http.NewRequest("GET", url+"/api/v1/workflows", nil)
	if err != nil {
		return false, 0, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Add API key header
	req.Header.Set("X-N8N-API-KEY", apiKey)
	req.Header.Set("Content-Type", "application/json")
	
	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		return false, 0, fmt.Errorf("failed to connect to n8n: %w", err)
	}
	defer resp.Body.Close()
	
	// Check response status
	if resp.StatusCode == 401 {
		return false, 0, fmt.Errorf("authentication failed - check your API key")
	}
	
	if resp.StatusCode != 200 {
		return false, 0, fmt.Errorf("n8n API returned status %d", resp.StatusCode)
	}
	
	// Parse response to get workflow count
	var workflows struct {
		Data []interface{} `json:"data"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&workflows); err != nil {
		// If we can't decode, at least we know the connection works
		return true, 0, nil
	}
	
	return true, len(workflows.Data), nil
}

// Helper function to prompt for input with a default value
func promptInput(prompt, defaultValue string) string {
	if defaultValue != "" {
		fmt.Printf("%s [%s]: ", prompt, defaultValue)
	} else {
		fmt.Printf("%s: ", prompt)
	}

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return defaultValue
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return defaultValue
	}

	return input
}

