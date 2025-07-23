package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [project-name]",
	Short: "Initialize a new n8n workflow project",
	Long: `Initialize a new n8n workflow project with the standard directory structure
and configuration files for collaborative development.

Examples:
  n8n-ops init my-workflows       # Initialize project in ./my-workflows
  n8n-ops init .                  # Initialize project in current directory`,
	RunE: runInit,
}

var (
	initForce bool
)

func init() {
	rootCmd.AddCommand(initCmd)

	initCmd.Flags().BoolVarP(&initForce, "force", "f", false, "force initialization, overwriting existing files")
}

func runInit(cmd *cobra.Command, args []string) error {
	projectDir := "."
	if len(args) > 0 {
		projectDir = args[0]
	}

	logger.Info("Initializing n8n workflow project", "directory", projectDir)

	// Create project directory if it doesn't exist
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Check if project already exists
	if !initForce {
		configPath := filepath.Join(projectDir, ".n8n-ops.yaml")
		if _, err := os.Stat(configPath); err == nil {
			return fmt.Errorf("project already exists in %s (use --force to overwrite)", projectDir)
		}
	}

	// Create directory structure
	dirs := []string{
		"workflows/development",
		"workflows/staging",
		"workflows/production",
		"scripts",
		"docs",
		"tests",
		"config",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(projectDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// Create configuration files
	if err := createConfigFiles(projectDir); err != nil {
		return fmt.Errorf("failed to create configuration files: %w", err)
	}

	// Create documentation files
	if err := createDocumentationFiles(projectDir); err != nil {
		return fmt.Errorf("failed to create documentation files: %w", err)
	}

	// Create git ignore file
	if err := createGitIgnoreFile(projectDir); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}

	logger.Info("Project initialized successfully", "directory", projectDir)

	fmt.Printf("✅ n8n workflow project initialized successfully!\n\n")
	fmt.Printf("Next steps:\n")
	fmt.Printf("1. Edit .n8n-ops.yaml to configure your n8n environments\n")
	fmt.Printf("2. Set up environment variables for API keys\n")
	fmt.Printf("3. Run 'n8n-ops sync --env development' to sync workflows\n")
	fmt.Printf("4. Start collaborating with your team!\n\n")
	fmt.Printf("Project structure:\n")
	fmt.Printf("  workflows/          # Environment-specific workflows\n")
	fmt.Printf("  ├── development/    # Development environment\n")
	fmt.Printf("  ├── staging/        # Staging environment\n")
	fmt.Printf("  └── production/     # Production environment\n")
	fmt.Printf("  docs/               # Documentation\n")
	fmt.Printf("  scripts/            # Custom scripts\n")
	fmt.Printf("  tests/              # Tests\n")
	fmt.Printf("  config/             # Configuration files\n")

	return nil
}

func createConfigFiles(projectDir string) error {
	// Main configuration file
	configContent := `# n8n CLI Configuration
# Environment configurations
environments:
  development:
    url: https://n8n-dev.example.com
    api_key_env: N8N_DEV_API_KEY
  staging:
    url: https://n8n-staging.example.com
    api_key_env: N8N_STAGING_API_KEY
  production:
    url: https://n8n-prod.example.com
    api_key_env: N8N_PROD_API_KEY

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

	// Environment variables example
	envContent := `# n8n CLI Environment Variables
# Copy this file to .env and fill in your actual values

# Development environment
N8N_DEV_API_KEY=your_dev_api_key_here
N8N_DEV_URL=https://n8n-dev.example.com

# Staging environment
N8N_STAGING_API_KEY=your_staging_api_key_here
N8N_STAGING_URL=https://n8n-staging.example.com

# Production environment
N8N_PROD_API_KEY=your_prod_api_key_here
N8N_PROD_URL=https://n8n-prod.example.com

# GitLab CI/CD variables (set in GitLab project settings)
# GITLAB_USER_EMAIL=user@example.com
# CI_COMMIT_SHA=abc123
# CI_PIPELINE_ID=12345
# CI_PIPELINE_URL=https://gitlab.com/project/pipelines/12345
`

	envPath := filepath.Join(projectDir, ".env.example")
	return os.WriteFile(envPath, []byte(envContent), 0644)
}

func createDocumentationFiles(projectDir string) error {
	readmeContent := `# n8n Workflow Project

This project contains n8n workflows managed with the n8n-ops tool for collaborative development.

## Quick Start

1. Install n8n-ops
2. Configure your environments in ` + "`.n8n-ops.yaml`" + `
3. Set up environment variables (copy ` + "`.env.example`" + ` to ` + "`.env`" + `)
4. Sync workflows: ` + "`n8n-ops sync --env development`" + `

## Commands

- ` + "`n8n-ops sync`" + ` - Sync workflows from n8n instance
- ` + "`n8n-ops deploy`" + ` - Deploy workflows to n8n instance
- ` + "`n8n-ops validate`" + ` - Validate workflow files
- ` + "`n8n-ops rollback`" + ` - Rollback to previous deployment

## Directory Structure

- ` + "`workflows/`" + ` - Environment-specific workflows
- ` + "`docs/`" + ` - Documentation
- ` + "`scripts/`" + ` - Custom scripts
- ` + "`tests/`" + ` - Tests
- ` + "`config/`" + ` - Configuration files

## Environments

- **development** - Development environment for testing
- **staging** - Pre-production environment
- **production** - Live production environment

## Contributing

1. Create feature branch
2. Make changes to workflows
3. Validate: ` + "`n8n-ops validate`" + `
4. Test deploy: ` + "`n8n-ops deploy --env development`" + `
5. Create merge request

## GitLab CI/CD

This project includes GitLab CI/CD pipeline configuration for automated:
- Workflow validation
- Automated deployment to staging
- Manual deployment to production
- Rollback capabilities
`

	readmePath := filepath.Join(projectDir, "README.md")
	return os.WriteFile(readmePath, []byte(readmeContent), 0644)
}

func createGitIgnoreFile(projectDir string) error {
	gitignoreContent := `# Environment variables
.env
.env.local

# Logs
logs/
*.log

# Database files
*.db
*.sqlite
*.sqlite3

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS files
.DS_Store
Thumbs.db

# Temporary files
tmp/
temp/
*.tmp

# Deployment reports
deployment-report-*.json
rollback-report-*.json

# Build artifacts
dist/
build/
bin/

# Node modules (if using npm scripts)
node_modules/
`

	gitignorePath := filepath.Join(projectDir, ".gitignore")
	return os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
}
