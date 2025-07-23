# n8n-ops Setup Guide

## Quick Setup

To run `n8n-ops` with your n8n environments, you need two things:

1. **The n8n-ops binary** (already built)
2. **Environment configuration** with your n8n API credentials

## Method 1: Using Environment Variables (.env file)

Create a `.env` file in the same directory as `n8n-ops`:

```bash
# Development Environment
N8N_DEV_URL=https://your-dev-n8n.com
N8N_DEV_API_KEY=your-dev-api-key

# Staging Environment  
N8N_STAGING_URL=https://your-staging-n8n.com
N8N_STAGING_API_KEY=your-staging-api-key

# Production Environment
N8N_PROD_URL=https://your-prod-n8n.com
N8N_PROD_API_KEY=your-prod-api-key

# GitLab Integration (optional)
GITLAB_TOKEN=your-gitlab-token
GITLAB_PROJECT_ID=your-project-id
```

Then run:
```bash
source .env
./n8n-ops sync --env development
```

## Method 2: Using Configuration File

Create `~/.n8n-ops.yaml` in your home directory:

```yaml
environments:
  development:
    url: "https://your-dev-n8n.com"
    api_key: "your-dev-api-key"
    
  staging:
    url: "https://your-staging-n8n.com"  
    api_key: "your-staging-api-key"
    
  production:
    url: "https://your-prod-n8n.com"
    api_key: "your-prod-api-key"

gitlab:
  token: "your-gitlab-token"
  project_id: "your-project-id"
  
logging:
  level: "info"
  format: "json"
```

## Method 3: Direct Environment Variables

Set variables directly in your shell:

```bash
export N8N_DEV_URL="https://your-dev-n8n.com"
export N8N_DEV_API_KEY="your-dev-api-key"

./n8n-ops sync --env development
```

## Getting n8n API Keys

1. Log into your n8n instance
2. Go to Settings → API Keys  
3. Create a new API key
4. Copy the key (it's only shown once!)

## Example Usage

Once configured, you can run commands like:

```bash
# Sync workflows from development
./n8n-ops sync --env development

# Deploy to staging
./n8n-ops deploy --env staging --force

# Validate local workflows
./n8n-ops validate ./workflows/

# Initialize new project
./n8n-ops init my-workflows

# Show welcome screen
./n8n-ops welcome
```

## Distribution Package

To distribute your tool, provide:

```
my-n8n-project/
├── n8n-ops                    # The binary executable
├── .env.example               # Template for environment variables
├── config.example.yaml        # Template for config file
├── SETUP_GUIDE.md            # This setup guide
└── README.md                 # Project documentation
```

Users only need to:
1. Copy `.env.example` to `.env`
2. Fill in their n8n URLs and API keys
3. Run `./n8n-ops --help` to get started

## Troubleshooting

- **"connection refused"**: Check your n8n URL is correct and accessible
- **"unauthorized"**: Verify your API key is valid
- **"command not found"**: Make sure n8n-ops binary has execute permissions (`chmod +x n8n-ops`)