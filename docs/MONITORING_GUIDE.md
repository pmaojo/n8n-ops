# n8n-ops Monitoring & Issue Management Guide

## Overview

The n8n-ops monitoring system automatically detects workflow failures and creates GitLab issues for tracking and resolution. This enables proactive workflow reliability management and ensures no failures go unnoticed.

## Features

### 🚨 Automatic Issue Creation
- Detects consecutive workflow failures
- Creates detailed GitLab issues with context
- Includes troubleshooting steps and metadata
- Assigns appropriate labels and severity

### ✅ Auto-Recovery Detection
- Monitors workflow recovery
- Updates issues with recovery information
- Automatically closes issues for auto-recovered workflows
- Tracks recovery patterns

### 📊 Intelligent Failure Detection
- Configurable failure thresholds
- Environment-specific settings
- Workflow-specific overrides
- Real-time monitoring

## Quick Start

### 1. Setup GitLab Integration

First, create a GitLab personal access token with `api` permissions:

1. Go to GitLab → User Settings → Access Tokens
2. Create token with `api` scope
3. Set environment variable:

```bash
export GITLAB_TOKEN="your-gitlab-token"
export GITLAB_PROJECT_ID="12345"  # Your project ID
```

### 2. Start Monitoring

```bash
# Monitor production environment
n8n-ops monitor --env production

# Custom settings
n8n-ops monitor \
  --env staging \
  --check-interval 30s \
  --failure-threshold 2 \
  --gitlab-project 12345
```

### 3. Demo Mode

Test the system without real GitLab integration:

```bash
n8n-ops monitor --demo --env development
```

## Configuration

### Command Line Options

| Flag | Description | Default |
|------|-------------|---------|
| `--check-interval` | How often to check for failures | 1m |
| `--failure-threshold` | Failures before creating issue | 3 |
| `--gitlab-project` | GitLab project ID | From env |
| `--gitlab-url` | GitLab instance URL | https://gitlab.com |

### Environment Variables

```bash
# Required for production use
export N8N_API_KEY="your-n8n-api-key"
export GITLAB_TOKEN="your-gitlab-token" 
export GITLAB_PROJECT_ID="12345"

# Optional
export GITLAB_URL="https://gitlab.company.com"
```

## Issue Management

### Issue Structure

When a workflow fails, the system creates issues with:

- **Title**: `🚨 Workflow Failure: [Name] ([Environment])`
- **Description**: Detailed failure information
- **Labels**: Automatic categorization
- **Assignees**: Team-based assignment (configurable)

### Issue Labels

| Label | Purpose |
|-------|---------|
| `workflow-failure` | Identifies workflow issues |
| `automated` | System-created issue |
| `env:production` | Environment identifier |
| `severity:high` | Based on failure count |
| `workflow:wf_123` | Specific workflow ID |

### Severity Levels

- **Low**: 1-2 failures
- **Medium**: 3-5 failures  
- **High**: 5+ failures

### Recovery Updates

When workflows recover, issues are updated with:

- Recovery timestamp
- Recovery type (auto/manual/rollback)
- New commit information
- Recovery notes

## Troubleshooting

### Common Issues

**1. No issues being created**
```bash
# Check GitLab token permissions
curl -H "PRIVATE-TOKEN: $GITLAB_TOKEN" \
  "https://gitlab.com/api/v4/projects/$GITLAB_PROJECT_ID"

# Verify n8n connection
n8n-ops monitor --demo --verbose
```

**2. Too many false positives**
```bash
# Increase failure threshold
n8n-ops monitor --failure-threshold 5

# Adjust check interval
n8n-ops monitor --check-interval 2m
```

**3. Missing environment variables**
```bash
# Check required variables
env | grep -E "(N8N_API_KEY|GITLAB_TOKEN|GITLAB_PROJECT_ID)"
```

### Monitoring Logs

Enable verbose logging:

```bash
n8n-ops monitor --verbose
```

## Advanced Configuration

### Workflow-Specific Settings

Some workflows may need different failure thresholds:

```yaml
# In your CI/CD or configuration
critical-workflows:
  - id: "payment-processing"
    threshold: 1  # Immediate issue
  - id: "user-notifications"
    threshold: 5  # More tolerant
```

### Integration with CI/CD

Add monitoring to your GitLab CI pipeline:

```yaml
# .gitlab-ci.yml
monitor_workflows:
  stage: monitor
  script:
    - n8n-ops monitor --check-interval 30s
  only:
    - schedules
  timeout: 1h
```

## Best Practices

### 1. Environment-Specific Thresholds
- **Production**: Lower thresholds (1-3 failures)
- **Staging**: Medium thresholds (2-5 failures)
- **Development**: Higher thresholds (5+ failures) or disabled

### 2. Team Assignment
- Assign critical workflows to specific teams
- Use GitLab project members for general assignments
- Set up notification rules for high-priority issues

### 3. Issue Lifecycle
- Use GitLab automation rules for issue routing
- Set up milestone tracking for resolution times
- Create templates for common failure types

### 4. Performance Considerations
- Adjust check intervals based on workflow frequency
- Use longer intervals for large n8n instances
- Monitor system resource usage

## API Reference

### Issue Creation Payload

The system creates issues with this structure:

```json
{
  "title": "🚨 Workflow Failure: Payment Processing (production)",
  "description": "## 🚨 Workflow Execution Failure\n\n**Workflow Details:**\n- **Name:** Payment Processing\n- **ID:** wf_payment_123\n- **Environment:** production\n- **Execution ID:** exec_456\n\n**Failure Information:**\n- **Failed At:** 2025-07-23T10:30:00Z\n- **Error:** Database connection timeout\n- **Retry Count:** 3\n- **Failed Node:** PostgreSQL\n\n**Pipeline Information:**\n- **Pipeline:** [123](https://gitlab.com/project/-/pipelines/123)\n- **Commit:** abc123def\n- **Branch:** main\n\n## 🔧 Troubleshooting Steps\n\n1. Check the workflow execution logs in n8n\n2. Verify input data and node configurations\n3. Test workflow in development environment\n4. Check for recent changes in the repository\n\n## 📊 Next Actions\n\n- [ ] Investigate root cause\n- [ ] Apply fix or rollback\n- [ ] Test in staging environment\n- [ ] Deploy to production\n\n---\n*This issue was automatically created by n8n-ops when a workflow failure was detected.*",
  "labels": [
    "workflow-failure",
    "automated", 
    "env:production",
    "workflow:wf_payment_123",
    "severity:medium",
    "node:PostgreSQL"
  ]
}
```

### Recovery Update Payload

When workflows recover:

```json
{
  "body": "## ✅ Workflow Recovery Detected\n\n**Recovery Information:**\n- **Recovered At:** 2025-07-23T11:00:00Z\n- **Recovery Type:** auto\n- **New Commit:** def456ghi\n\n**Notes:** Database connection issue resolved automatically\n\n---\n*This update was automatically generated by n8n-ops.*"
}
```

## Monitoring Dashboard

Consider setting up a GitLab dashboard to track:

- Open workflow failure issues
- Resolution times by environment  
- Most frequently failing workflows
- Recovery success rates

This provides visibility into workflow reliability trends and helps identify patterns requiring attention.