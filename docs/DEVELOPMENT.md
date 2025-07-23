# n8n CLI - Development Guide

## 📋 Overview

This guide helps you set up, configure, and contribute to the n8n CLI tool for collaborative workflow management. The CLI compiles workflows, syncs between environments, and integrates with GitLab CI/CD for team collaboration.

## 🚀 Quick Start

### Prerequisites

- Go 1.19+ installed
- Git configured
- Access to n8n instances (development, staging, production)
- GitLab account with repository access

### Installation

```bash
# Clone the repository
git clone https://gitlab.com/your-org/n8n-ops.git
cd n8n-ops

# Build the CLI tool
go build -o n8n-ops main.go

# Make it executable (Linux/macOS)
chmod +x n8n-ops

# Test installation
./n8n-ops welcome
./n8n-ops --help
```

## ⚙️ Environment Configuration

### 1. Local Development Environment

Create a configuration file at `~/.n8n-ops.yaml`:

```yaml
# Default environment settings
default_environment: development

# Environment configurations
environments:
  development:
    url: "https://n8n-dev.yourcompany.com"
    api_key: "${N8N_DEV_API_KEY}"
    timeout: 30
    
  staging:
    url: "https://n8n-staging.yourcompany.com"
    api_key: "${N8N_STAGING_API_KEY}"
    timeout: 30
    
  production:
    url: "https://n8n-prod.yourcompany.com"
    api_key: "${N8N_PROD_API_KEY}"
    timeout: 60

# Git branch to environment mapping
branch_mapping:
  main: production
  staging: staging
  develop: development
  feature/*: development

# Local database settings
database:
  path: ".n8n-ops/database.sqlite"
  
# Logging configuration
logging:
  level: "info"  # debug, info, warn, error
  file: ".n8n-ops/logs/n8n-ops.log"
```

### 2. Environment Variables

Set the following environment variables in your shell:

```bash
# n8n API Keys (get these from your n8n instance admins)
export N8N_DEV_API_KEY="n8n_api_dev_your_key_here"
export N8N_STAGING_API_KEY="n8n_api_staging_your_key_here" 
export N8N_PROD_API_KEY="n8n_api_prod_your_key_here"

# GitLab Configuration (for CI/CD)
export GITLAB_TOKEN="glpat-your_gitlab_token_here"
export GITLAB_PROJECT_ID="12345"

# Optional: Language preference
export N8N_CLI_LANG="en"  # or "es" for Spanish
```

### 3. Directory Structure

The CLI expects this project structure:

```
your-project/
├── workflows/
│   ├── development/     # Dev environment workflows
│   ├── staging/         # Staging environment workflows
│   └── production/      # Production environment workflows
├── scripts/             # Custom deployment scripts
├── tests/               # Test workflows and fixtures
├── config/
│   └── templates/       # Workflow templates
├── .n8n-ops/
│   ├── database.sqlite  # Local tracking database
│   └── logs/            # CLI logs
└── .gitlab-ci.yml       # GitLab CI/CD configuration
```

## 🔧 GitLab Configuration

### 1. GitLab CI/CD Variables

In your GitLab project, go to **Settings → CI/CD → Variables** and add:

| Variable | Value | Protected | Masked |
|----------|-------|-----------|---------|
| `N8N_DEV_API_KEY` | Your dev n8n API key | No | Yes |
| `N8N_STAGING_API_KEY` | Your staging n8n API key | No | Yes |
| `N8N_PROD_API_KEY` | Your prod n8n API key | Yes | Yes |
| `N8N_DEV_URL` | https://n8n-dev.yourcompany.com | No | No |
| `N8N_STAGING_URL` | https://n8n-staging.yourcompany.com | No | No |
| `N8N_PROD_URL` | https://n8n-prod.yourcompany.com | Yes | No |

### 2. GitLab CI/CD Pipeline (`.gitlab-ci.yml`)

```yaml
# GitLab CI/CD Configuration for n8n CLI
stages:
  - validate
  - test
  - deploy-dev
  - deploy-staging
  - deploy-production

variables:
  GO_VERSION: "1.19"

# Reusable deployment template
.deploy_template:
  image: golang:${GO_VERSION}
  dependencies:
    - build-cli
  script:
    - go build -o n8n-ops main.go
    - ./n8n-ops validate ./workflows/
  rules:
    - if: '$CI_PIPELINE_SOURCE == "merge_request_event"'
    - if: '$CI_COMMIT_BRANCH'

# Test CLI functions
test-cli-functions:
  stage: test
  image: golang:${GO_VERSION}
  script:
    - go build -o n8n-ops main.go
    - ./n8n-ops --help
    - ./n8n-ops welcome
  rules:
    - if: '$CI_PIPELINE_SOURCE == "merge_request_event"'
    - if: '$CI_COMMIT_BRANCH'

# Run Go unit tests
unit-tests:
  stage: test
  image: golang:${GO_VERSION}
  cache:
    key: ${CI_COMMIT_REF_SLUG}
    paths:
      - ${GOMODCACHE}
  script:
    - go vet ./...
    - go test ./...
  rules:
    - if: '$CI_PIPELINE_SOURCE == "merge_request_event"'
    - if: '$CI_COMMIT_BRANCH'

# Deploy to development (automatic on develop branch)
deploy-development:
  extends: .deploy_template
  stage: deploy-dev
  variables:
    TARGET_ENV: development
    N8N_API_KEY: ${N8N_DEV_API_KEY}
    N8N_URL: ${N8N_DEV_URL}
  environment:
    name: development
    url: $N8N_DEV_URL
  rules:
    - if: '$CI_COMMIT_BRANCH == "develop"'

# Deploy to staging (manual)
deploy-staging:
  extends: .deploy_template
  stage: deploy-staging
  variables:
    TARGET_ENV: staging
    N8N_API_KEY: ${N8N_STAGING_API_KEY}
    N8N_URL: ${N8N_STAGING_URL}
  environment:
    name: staging
    url: $N8N_STAGING_URL
  when: manual
  rules:
    - if: '$CI_COMMIT_BRANCH == "staging"'

# Deploy to production (manual + protected)
deploy-production:
  extends: .deploy_template
  stage: deploy-production
  variables:
    TARGET_ENV: production
    N8N_API_KEY: ${N8N_PROD_API_KEY}
    N8N_URL: ${N8N_PROD_URL}
  environment:
    name: production
    url: $N8N_PROD_URL
  when: manual
  allow_failure: false
  rules:
    - if: '$CI_COMMIT_BRANCH == "main"'
```

## 🌊 Workflow Development Process

### 1. Starting a New Feature

```bash
# Create feature branch
git checkout -b feature/new-payment-workflow

# Initialize workflow structure
./n8n-ops init --template payment

# Start development
./n8n-ops sync --env development
```

### 2. Development Cycle

```bash
# Edit workflows in your preferred editor
# Workflows are stored as JSON in ./workflows/development/

# Validate changes
./n8n-ops validate ./workflows/development/

# Test compilation
./n8n-ops deploy --env development --dry-run

# Deploy to development
./n8n-ops deploy --env development

# Sync any changes made in n8n UI back to files
./n8n-ops sync --env development --force
```

### 3. Promotion to Staging/Production

```bash
# Create merge request to staging branch
git push origin feature/new-payment-workflow
# → Create MR to staging branch in GitLab

# After merge and testing in staging:
# Create MR from staging to main for production deployment
```

## 🛠️ CLI Commands Reference

### Core Commands

```bash
# Initialize new project
./n8n-ops init

# Display welcome screen
./n8n-ops welcome

# Show help
./n8n-ops --help

# Change language
./n8n-ops --lang es welcome
```

### Workflow Management

```bash
# Sync workflows FROM n8n TO local files
./n8n-ops sync --env development
./n8n-ops sync --env staging --force

# Deploy workflows FROM local files TO n8n
./n8n-ops deploy --env development
./n8n-ops deploy workflow.json --env staging
./n8n-ops deploy --env production --dry-run

# Validate workflow files
./n8n-ops validate ./workflows/
./n8n-ops validate specific-workflow.json
```

### Git Integration

```bash
# Check current branch mapping
./n8n-ops branch current

# List branch mappings
./n8n-ops branch list

# Set branch to environment mapping
./n8n-ops branch set feature/auth development
```

### Rollback and Recovery

```bash
# Rollback to previous deployment
./n8n-ops rollback --env staging

# View deployment history
./n8n-ops history --env production
```

## 🔍 Troubleshooting

### Common Issues

1. **API Key Not Working**
   ```bash
   # Test API connection
   curl -H "X-N8N-API-KEY: your_api_key" https://your-n8n-instance.com/api/v1/workflows
   ```

2. **Build Failures**
   ```bash
   # Clean and rebuild
   go clean
   go mod tidy
   go build -o n8n-ops main.go
   ```

3. **Database Issues**
   ```bash
   # Reset local database
   rm -f .n8n-ops/database.sqlite
   ./n8n-ops init
   ```

### Debug Mode

```bash
# Enable verbose logging
./n8n-ops --verbose sync --env development

# Check logs
tail -f .n8n-ops/logs/n8n-ops.log
```

## 🔐 Security Best Practices

1. **API Keys**: Never commit API keys to git. Use environment variables or GitLab CI/CD variables.

2. **Production Access**: Limit production API keys to authorized personnel only.

3. **Branch Protection**: Set up branch protection rules in GitLab for main/staging branches.

4. **Audit Trail**: The CLI maintains deployment history in the local SQLite database.

## 📝 Contributing

### Before Committing

```bash
# Validate your changes
./n8n-ops validate ./workflows/

# Test CLI compilation
go build -o n8n-ops main.go
./n8n-ops --help

# Run tests (if available)
go test ./...
```

### Commit Message Format

```
feat: add new payment workflow template
fix: resolve API timeout in production sync  
docs: update deployment guide for new environments
```

## 📞 Support

- **Documentation**: Check this guide and `README.md`
- **Issues**: Create issues in GitLab repository
- **Team Chat**: Use your team's communication channel
- **API Documentation**: Refer to n8n API documentation

---

**Happy workflow automation! 🚀**