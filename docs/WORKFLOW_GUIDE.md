# n8n-ops Collaborative Workflow Guide

## The Complete Development-to-Production Flow

### 🔄 Development Cycle

```bash
# 1. Sync workflows from server to local
./n8n-ops sync --env development
# Downloads: ./workflows/development/*.json

# 2. Work locally
# - Edit JSON files 
# - Test on local n8n instance
# - Make iterative changes

# 3. Validate before committing
./n8n-ops validate ./workflows/development/
# ✅ Validates JSON structure, required fields, etc.

# 4. Commit to Git
git add workflows/
git commit -m "Updated customer onboarding workflow"
git push origin feature/new-feature
```

### 🚀 Automated Server Deployment

When you push to GitLab, the CI/CD pipeline automatically:

1. **Validates** all workflow files
2. **Tests** compatibility 
3. **Deploys** to target environment
4. **Updates** or **creates new** workflows on server
5. **Tracks** deployment in local SQLite database

### 🎯 Environment Flow

```
Local Development → Git Push → GitLab CI/CD → Server Import
     ↓                ↓            ↓              ↓
Edit JSON files → Validation → Auto Deploy → n8n Instance Updated
```

### 📋 Workflow States

The tool handles both scenarios automatically:

- **New Workflows**: Creates them on the server
- **Existing Workflows**: Updates them while preserving IDs and settings
- **Rollback**: Can revert to previous versions if issues occur

### 🔧 Manual Deployment

You can also deploy manually:

```bash
# Deploy to staging
./n8n-ops deploy --env staging --force

# Deploy specific workflow
./n8n-ops deploy --env production workflow_name.json

# Rollback if needed
./n8n-ops rollback --env production
```

### 🌳 Branch Strategy

The tool supports environment mapping:

```yaml
# In config.yaml
branch_mapping:
  main: production        # main → auto-deploy to production
  staging: staging        # staging → auto-deploy to staging  
  develop: development    # develop → auto-deploy to development
  "feature/*": development # feature branches → development
```

### 📊 Tracking & History

Every deployment is tracked:

- **Local SQLite**: Stores deployment history
- **Git History**: Full version control
- **n8n Logs**: Server-side execution logs
- **Rollback Data**: Previous versions available

### 🔒 Production Safety

For production deployments:

- **Approval Gates**: Manual approval required
- **Backup**: Automatic backup before deployment
- **Validation**: Enhanced validation rules
- **Rollback**: One-command rollback capability

## Example Complete Flow

```bash
# Developer starts new feature
git checkout -b feature/payment-integration
./n8n-ops sync --env development

# Edit workflow locally
vim workflows/development/payment_processor.json

# Validate changes
./n8n-ops validate workflows/development/

# Commit and push
git add . && git commit -m "Add Stripe payment integration"
git push origin feature/payment-integration

# GitLab CI/CD automatically:
# 1. Validates workflow
# 2. Deploys to development n8n instance  
# 3. Sends Slack notification: "✅ payment_processor deployed to dev"

# Create merge request to staging
# → Manual approval → Auto-deploy to staging

# Create merge request to main  
# → Manual approval → Auto-deploy to production
```

This gives you the full Git workflow with automated n8n synchronization!