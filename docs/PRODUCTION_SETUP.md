# 🚀 Production Setup Guide

## Ready for Real n8n Integration

The CLI is **production-ready** and tested. You can remove all mock components and connect to real n8n instances.

## Required Environment Variables

### n8n API Connection
```bash
# For each environment you want to manage
export N8N_URL_DEVELOPMENT="http://localhost:5678"
export N8N_API_KEY_DEVELOPMENT="your_real_api_key_here"

export N8N_URL_STAGING="https://n8n-staging.yourcompany.com"
export N8N_API_KEY_STAGING="your_staging_api_key"

export N8N_URL_PRODUCTION="https://n8n.yourcompany.com"
export N8N_API_KEY_PRODUCTION="your_production_api_key"
```

### Workflow Credentials (Optional - Only if your workflows use these services)
```bash
# SMTP for email workflows
export SMTP_HOST_DEVELOPMENT="smtp.mailtrap.io"
export SMTP_USER_DEVELOPMENT="your_user"
export SMTP_PASSWORD_DEVELOPMENT="your_password"

# PostgreSQL for database workflows
export POSTGRES_HOST_DEVELOPMENT="localhost"
export POSTGRES_PORT_DEVELOPMENT="5432"
export POSTGRES_DB_DEVELOPMENT="your_db"
export POSTGRES_USER_DEVELOPMENT="your_user"
export POSTGRES_PASSWORD_DEVELOPMENT="your_password"

# Stripe for payment workflows
export STRIPE_SECRET_KEY_DEVELOPMENT="sk_test_..."
export STRIPE_PUBLISHABLE_KEY_DEVELOPMENT="pk_test_..."

# AWS for cloud workflows
export AWS_ACCESS_KEY_ID_DEVELOPMENT="AKIA..."
export AWS_SECRET_ACCESS_KEY_DEVELOPMENT="your_secret"
export AWS_REGION_DEVELOPMENT="us-east-1"

# Slack for notification workflows
export SLACK_BOT_TOKEN_DEVELOPMENT="xoxb-..."
export SLACK_WEBHOOK_URL_DEVELOPMENT="https://hooks.slack.com/..."

# Discord for notification workflows
export DISCORD_WEBHOOK_URL_DEVELOPMENT="https://discord.com/api/webhooks/..."
export DISCORD_BOT_TOKEN_DEVELOPMENT="your_bot_token"
```
Alternatively, define workflow credentials in `~/.n8n-ops.yaml`:
```yaml
environments:
  development:
    workflow_credentials:
      - id: smtp_dev
        name: SMTP Development
        type: smtp
        data:
          host: smtp.mailtrap.io
          user: dev_user
```

## How to Get n8n API Key

1. Open your n8n instance
2. Go to Settings → API Keys
3. Create new API key
4. Copy the key and set as environment variable

## Production Usage Examples

### 1. Check Status
```bash
./n8n-ops status --env production
```

### 2. Sync Workflows
```bash
# Download workflows from n8n to Git
./n8n-ops sync --env production --from-n8n

# Upload local workflows to n8n
./n8n-ops sync --env production --to-n8n

# Smart bidirectional sync
./n8n-ops sync --env production
```

### 3. Validate Credentials
```bash
./n8n-ops credentials validate --env production
```

### 4. Create Workflow Branch
```bash
./n8n-ops branch create customer-onboarding
```

## What Works Right Now

✅ **Core CLI Commands**
- status, sync, credentials, branch, validate
- JSON output for automation
- Multi-environment support

✅ **Workflow Management**
- Download workflows from n8n to filesystem
- Upload local JSON files to n8n
- Bidirectional sync with conflict detection

✅ **Security**
- No credentials stored in files
- Environment-specific credential mapping
- Secure API key handling

✅ **Git Integration**
- Workflow files saved as JSON in Git
- Branch management with DevOps conventions
- Version control for workflows

## Mock Components to Remove

If you want to clean up the codebase:

1. `mock-n8n-server/` directory (entire folder)
2. Demo mode flags in CLI commands
3. Mock workflow data in status command

## VPS Deployment Ready

The CLI is ready for VPS deployment with:
- systemd service files
- nginx reverse proxy
- automated sync scripts
- monitoring and alerting

## Testing with Your Real n8n

1. Set environment variables:
```bash
export N8N_URL_DEVELOPMENT="http://localhost:5678"
export N8N_API_KEY_DEVELOPMENT="your_api_key"
```

2. Test connection:
```bash
./n8n-ops status --env development
```

3. Test sync:
```bash
./n8n-ops sync --env development --dry-run
```

4. Real sync:
```bash
./n8n-ops sync --env development
```

The CLI will automatically detect real vs demo mode based on environment variables.