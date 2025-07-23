# n8n-ops Onboarding Video Script

## Introduction (0:00 - 0:30)
- Welcome to n8n-ops, a powerful CLI tool for managing n8n workflows across multiple environments
- In this video, we'll show you how to get started in just a few minutes
- By the end, you'll be able to sync, validate, and deploy workflows between environments

## What is n8n-ops? (0:30 - 1:00)
- n8n-ops is a CLI tool that helps you manage n8n workflows using Git
- It allows you to:
  - Sync workflows between n8n instances and local files
  - Validate workflow JSON files
  - Deploy workflows to different environments
  - Monitor workflow status
  - Collaborate with your team using Git

## Installation (1:00 - 2:00)
- Clone the repository: `git clone https://gitlab.com/your-org/n8n-ops.git`
- Build the CLI: `go build -o n8n-ops main.go`
- Make it executable: `chmod +x n8n-ops`
- Test it works: `./n8n-ops welcome`
- You should see a futuristic welcome screen with ASCII art

## Interactive Onboarding (2:00 - 4:00)
- The easiest way to get started is with our interactive onboarding wizard
- Run: `./n8n-ops onboard`
- Follow the prompts to:
  - Set up your project structure
  - Configure your environments
  - Set up API keys
  - Test connections

## Getting n8n API Keys (4:00 - 5:00)
- Log in to your n8n instance
- Go to Settings → API Keys
- Click "Create API Key"
- Name it "n8n-ops-development"
- Copy the generated key (starts with "n8n_api_")
- Add it to your .env file

## Basic Workflow (5:00 - 7:00)
- Sync workflows from n8n: `./n8n-ops sync --env development`
- Make changes to workflows (in n8n UI or edit JSON files)
- Validate your changes: `./n8n-ops validate ./workflows/development/`
- Deploy your changes: `./n8n-ops deploy --env development`
- Commit to Git: `git add workflows/ && git commit -m "update workflows"`

## Multi-Environment Setup (7:00 - 8:00)
- Configure staging and production environments in .n8n-ops.yaml
- Get API keys for each environment
- Use `--env` flag to specify environment: 
  - `./n8n-ops sync --env staging`
  - `./n8n-ops deploy --env production`

## GitLab CI/CD Integration (8:00 - 9:00)
- n8n-ops includes GitLab CI/CD pipeline configuration
- Set up GitLab CI/CD variables for API keys
- Automatic deployment to development on push to develop branch
- Manual approval for staging and production deployments

## Troubleshooting Tips (9:00 - 9:30)
- Check your API keys and URLs
- Use `--verbose` flag for detailed logs
- Try `--dry-run` before actual operations
- Check the documentation for more help

## Conclusion (9:30 - 10:00)
- You're now ready to use n8n-ops for workflow management
- Explore more commands with `./n8n-ops --help`
- Check out the documentation for advanced features
- Happy workflow automation!

## Call to Action
- Star the repository on GitLab
- Report issues or contribute to the project
- Share your feedback with us