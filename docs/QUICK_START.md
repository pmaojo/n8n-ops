# 🚀 Quick Start Guide

Get up and running with the n8n CLI in 5 minutes!

## Prerequisites

- Go 1.19+ installed
- Git configured
- Access to n8n instances (development/staging/production)

## Step 1: Clone and Build

```bash
# Clone the repository
git clone https://gitlab.com/your-org/n8n-ops.git
cd n8n-ops

# Build the CLI
go build -o n8n-ops main.go

# Test it works
./n8n-ops welcome
```

## Step 2: Get n8n API Keys

### For each n8n environment:

1. **Login to your n8n instance** (e.g., https://n8n-dev.yourcompany.com)
2. **Go to Settings → API Keys**
3. **Click "Create API Key"**
4. **Give it a name**: `n8n-ops-development` (or staging/production)
5. **Copy the generated key** (starts with `n8n_api_`)

### Environment URLs you need:
- Development: `https://n8n-dev.yourcompany.com`
- Staging: `https://n8n-staging.yourcompany.com` 
- Production: `https://n8n-prod.yourcompany.com`

## Step 3: Configuration

### Create config file:
```bash
cp config.example.yaml ~/.n8n-ops.yaml
```

### Edit `~/.n8n-ops.yaml`:
```yaml
environments:
  development:
    url: "https://n8n-dev.yourcompany.com"
    api_key: "${N8N_DEV_API_KEY}"
  staging:
    url: "https://n8n-staging.yourcompany.com"  
    api_key: "${N8N_STAGING_API_KEY}"
  production:
    url: "https://n8n-prod.yourcompany.com"
    api_key: "${N8N_PROD_API_KEY}"
```

### Set environment variables:
```bash
# Add to your ~/.bashrc or ~/.zshrc
export N8N_DEV_API_KEY="n8n_api_your_dev_key_here"
export N8N_STAGING_API_KEY="n8n_api_your_staging_key_here"
export N8N_PROD_API_KEY="n8n_api_your_prod_key_here"

# Reload your shell
source ~/.bashrc  # or ~/.zshrc
```

## Step 4: Test Connection

```bash
# Test development connection
./n8n-ops sync --env development --dry-run

# If successful, you should see:
# ✅ Connected to n8n instance
# ✅ Found X workflows
```

## Step 5: Basic Workflow

```bash
# 1. Sync workflows from n8n to local files
./n8n-ops sync --env development

# 2. Validate local workflow files  
./n8n-ops validate ./workflows/development/

# 3. Make changes to workflows (edit JSON files or via n8n UI)

# 4. Deploy changes back to n8n
./n8n-ops deploy --env development
```

## GitLab CI/CD Setup

### Step 1: Add GitLab Variables

In GitLab project → **Settings → CI/CD → Variables**:

| Variable Name | Value | Masked |
|---------------|--------|---------|
| `N8N_DEV_API_KEY` | Your dev API key | ✅ |
| `N8N_STAGING_API_KEY` | Your staging API key | ✅ |  
| `N8N_PROD_API_KEY` | Your prod API key | ✅ |
| `N8N_DEV_URL` | https://n8n-dev.yourcompany.com | ❌ |
| `N8N_STAGING_URL` | https://n8n-staging.yourcompany.com | ❌ |
| `N8N_PROD_URL` | https://n8n-prod.yourcompany.com | ❌ |

### Step 2: Commit Pipeline

The included `.gitlab-ci.yml` will automatically:
- ✅ Build and test the CLI on every commit
- ✅ Validate all workflow files
- ✅ Deploy to development on `develop` branch
- ✅ Deploy to staging on `staging` branch (manual approval)
- ✅ Deploy to production on `main` branch (manual approval)

## Common Commands

```bash
# Show beautiful welcome screen
./n8n-ops welcome

# Spanish welcome screen  
./n8n-ops --lang es welcome

# Sync FROM n8n TO local files
./n8n-ops sync --env development
./n8n-ops sync --env staging --force

# Deploy FROM local files TO n8n
./n8n-ops deploy --env development
./n8n-ops deploy workflow.json --env staging

# Validate workflow files
./n8n-ops validate ./workflows/

# Check current git branch mapping
./n8n-ops branch current

# Preview deployment (dry run)
./n8n-ops deploy --env staging --dry-run

# Rollback if something goes wrong
./n8n-ops rollback --env production
```

## Development Workflow

```bash
# 1. Create feature branch
git checkout -b feature/new-payment-flow

# 2. Sync current workflows
./n8n-ops sync --env development

# 3. Make changes in n8n UI or edit JSON files

# 4. Validate changes
./n8n-ops validate ./workflows/development/

# 5. Test deployment
./n8n-ops deploy --env development --dry-run
./n8n-ops deploy --env development

# 6. Commit and push
git add workflows/development/
git commit -m "feat: add new payment workflow"  
git push origin feature/new-payment-flow

# 7. Create merge request to staging branch
```

## Troubleshooting

### API Key Issues
```bash
# Test API connection manually
curl -H "X-N8N-API-KEY: your_api_key" \
     https://your-n8n-instance.com/api/v1/workflows
```

### Build Issues
```bash
# Clean and rebuild
go clean
go mod tidy
go build -o n8n-ops main.go
```

### Configuration Issues
```bash
# Check your config
cat ~/.n8n-ops.yaml

# Check environment variables
echo $N8N_DEV_API_KEY
```

## Next Steps

- 📖 Read the complete [Development Guide](DEVELOPMENT.md)
- 🔧 Explore all CLI commands with `./n8n-ops --help`
- 🌍 Try different languages with `--lang es`
- 🎨 Enjoy the futuristic ASCII art interface!

---

**You're all set! Happy workflow automation! 🚀**