# n8n-ops Onboarding Guide

Welcome to n8n-ops! This guide will help you get started quickly.

## Interactive Onboarding (Recommended)

The easiest way to get started is with our interactive onboarding wizard:

```bash
./n8n-ops onboard
```

This wizard will:
- Create your project structure
- Configure your environments
- Set up API keys
- Test connections
- Show you next steps

## Manual Setup

If you prefer to set things up manually, follow these steps:

### 1. Initialize Your Project

```bash
# Create a new project in the current directory
./n8n-ops init .
```

### 2. Configure Your Environments

Edit the `.n8n-ops.yaml` file to configure your environments:

```yaml
environments:
  development:
    url: "http://localhost:5678"
    api_key_env: "N8N_API_KEY_DEV"
  staging:
    url: "https://n8n-staging.example.com"
    api_key_env: "N8N_API_KEY_STAGING"
  production:
    url: "https://n8n-prod.example.com"
    api_key_env: "N8N_API_KEY_PROD"
```

### 3. Set Up API Keys

Create a `.env` file with your API keys:

```
N8N_API_KEY_DEV=your_dev_api_key_here
N8N_API_KEY_STAGING=your_staging_api_key_here
N8N_API_KEY_PROD=your_prod_api_key_here
```

### 4. Test Your Connection

```bash
./n8n-ops sync --env development --dry-run
```

## Basic Workflow

Once you're set up, here's the basic workflow:

1. **Sync workflows** from n8n to local files:
   ```bash
   ./n8n-ops sync --env development
   ```

2. **Make changes** to workflows (in n8n UI or edit JSON files)

3. **Validate** your changes:
   ```bash
   ./n8n-ops validate ./workflows/development/
   ```

4. **Deploy** your changes:
   ```bash
   ./n8n-ops deploy --env development
   ```

5. **Commit** to Git:
   ```bash
   git add workflows/
   git commit -m "update workflows"
   git push
   ```

## Visual Workflow

```
┌───────────┐         ┌───────────┐         ┌───────────┐
│           │  sync   │           │ deploy  │           │
│    n8n    │ ◄─────► │   Files   │ ◄─────► │    n8n    │
│    Dev    │         │   (Git)   │         │   Prod    │
└───────────┘         └───────────┘         └───────────┘
```

## Need Help?

- Run `./n8n-ops --help` for command information
- Check the [Quick Start Guide](../QUICK_START_IMPROVED.md)
- Read the [Development Guide](../DEVELOPMENT.md)

Happy workflow automation! 🚀