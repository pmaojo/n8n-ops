# Repository Setup Guide

Complete guide to setting up the n8n CLI repository in GitLab for your team.

## GitLab Repository Configuration

### 1. Create GitLab Project

```bash
# Create new GitLab project
# Go to GitLab → New Project → Create blank project
# Project name: n8n-ops
# Visibility: Private (recommended)
# Initialize with README: No (we have our own)
```

### 2. Initial Repository Setup

```bash
# Clone this repository and push to your GitLab
git clone https://replit.com/your-replit-project
cd n8n-ops

# Add your GitLab remote
git remote add origin https://gitlab.com/your-org/n8n-ops.git

# Push to GitLab
git push -u origin main
```

### 3. Branch Protection Rules

#### Protect Main Branch
**Settings → Repository → Push Rules**:
- [x] Prohibit committing secrets
- [x] Restrict commits by author (email)
- [x] Prohibited branch names: `master`

**Settings → Repository → Protected Branches**:

| Branch | Push | Merge | Unprotect |
|---------|------|-------|-----------|
| `main` | Maintainers only | Maintainers only | Maintainers only |
| `staging` | Developers + Maintainers | Maintainers only | Maintainers only |

#### Branch Settings:
- [x] Require merge request approval
- [x] Reset approvals on new commits
- [x] Prevent approval by author
- [x] Require resolved discussions

### 4. CI/CD Variables

**Settings → CI/CD → Variables**

#### n8n API Configuration
| Variable | Value | Protected | Masked | Environment |
|----------|-------|-----------|---------|-------------|
| `N8N_API_KEY_DEV` | `n8n_api_dev_xxx` | No | Yes | All |
| `N8N_API_KEY_STAGING` | `n8n_api_staging_xxx` | No | Yes | All |
| `N8N_API_KEY_PROD` | `n8n_api_prod_xxx` | Yes | Yes | Production |

#### n8n Instance URLs
| Variable | Value | Protected | Masked | Environment |
|----------|-------|-----------|---------|-------------|
| `N8N_URL_DEV` | `https://n8n-dev.company.com` | No | No | All |
| `N8N_URL_STAGING` | `https://n8n-staging.company.com` | No | No | All |
| `N8N_URL_PROD` | `https://n8n-prod.company.com` | Yes | No | Production |

#### Optional Notification Variables
| Variable | Value | Protected | Masked | Environment |
|----------|-------|-----------|---------|-------------|
| `SLACK_WEBHOOK_URL` | `https://hooks.slack.com/xxx` | No | Yes | All |
| `TEAMS_WEBHOOK_URL` | `https://company.webhook.office.com/xxx` | No | Yes | All |

### 5. GitLab Environments

**Settings → CI/CD → Environments**

Create these environments:

#### Development Environment
- **Name**: `development`
- **External URL**: `$N8N_URL_DEV`
- **Auto-stop**: No
- **Protected**: No

#### Staging Environment  
- **Name**: `staging`
- **External URL**: `$N8N_URL_STAGING`
- **Auto-stop**: No
- **Protected**: No

#### Production Environment
- **Name**: `production`
- **External URL**: `$N8N_URL_PROD` 
- **Auto-stop**: No
- **Protected**: Yes
- **Allowed to deploy**: Maintainers only

### 6. Merge Request Templates

**Create**: `.gitlab/merge_request_templates/Default.md`

```markdown
## Summary
Brief description of changes

## Type of Change
- [ ] Bug fix
- [ ] New feature  
- [ ] Breaking change
- [ ] Documentation update
- [ ] Workflow changes

## Testing
- [ ] All workflow files validate
- [ ] CLI builds successfully
- [ ] Tested in development environment
- [ ] Manual testing completed

## Workflow Changes
- [ ] New workflows added
- [ ] Existing workflows modified
- [ ] Workflow validation updated
- [ ] Environment-specific changes

## Deployment Notes
- [ ] Database migrations needed
- [ ] Configuration changes required
- [ ] API key updates needed
- [ ] Environment variables added

## Checklist
- [ ] Code follows project standards
- [ ] Self-review completed
- [ ] Documentation updated
- [ ] No secrets in code
- [ ] Commit messages follow convention

/label ~"needs review"
```

### 7. Issue Templates

**Create**: `.gitlab/issue_templates/Bug.md`

```markdown
## Bug Report

**Environment**: development/staging/production

**CLI Version**: 
```bash
./n8n-ops --version
```

**Expected Behavior**:

**Actual Behavior**:

**Steps to Reproduce**:
1. 
2. 
3. 

**Error Logs**:
```
paste error logs here
```

**Additional Context**:

/label ~bug ~needs-triage
```

**Create**: `.gitlab/issue_templates/Feature.md`

```markdown
## Feature Request

**Is your feature request related to a problem?**

**Describe the solution you'd like**:

**Describe alternatives you've considered**:

**Additional context**:

**Implementation Notes**:
- [ ] CLI command changes needed
- [ ] n8n API integration required  
- [ ] Workflow template updates
- [ ] Documentation updates
- [ ] Configuration changes

/label ~enhancement ~needs-discussion
```

### 8. Repository Settings

#### General Settings
**Settings → General**:
- **Project name**: n8n CLI
- **Project description**: Go-based CLI tool for managing n8n workflows across environments
- **Project avatar**: Upload n8n or company logo
- **Visibility**: Private
- **Issues**: Enabled
- **Merge requests**: Enabled
- **Wiki**: Disabled (we use docs in repo)
- **Snippets**: Enabled

#### Repository Settings  
**Settings → Repository**:
- **Default branch**: `main`
- **Auto-close referenced issues**: Enabled
- **Push rules**: Configure as mentioned above

### 9. Team Access

#### Add Team Members
**Project → Members**:

| Role | Permissions |
|------|-------------|
| **Owner** | Full access, can delete project |
| **Maintainer** | Manage team, protect branches, deploy to production |
| **Developer** | Push code, create MRs, deploy to dev/staging |  
| **Reporter** | View code, create issues, comment |

#### Recommended Team Structure:
- **DevOps Lead**: Owner
- **Senior Developers**: Maintainer  
- **Developers**: Developer
- **QA Team**: Reporter
- **Stakeholders**: Reporter

### 10. Webhooks (Optional)

**Settings → Webhooks**

#### Slack Integration
- **URL**: Your Slack webhook URL
- **Trigger**: Push events, Merge request events, Pipeline events
- **Enable SSL verification**: Yes

#### External Monitoring
Configure webhooks for:
- Deployment notifications
- Error alerting  
- Performance monitoring
- Security scanning

## Repository Structure Validation

After setup, your repository should have:

```
n8n-ops/
├── .gitlab-ci.yml                 ✅ CI/CD pipeline
├── .gitlab/
│   ├── merge_request_templates/   ✅ MR templates
│   └── issue_templates/          ✅ Issue templates
├── .gitignore                    ✅ Git ignore rules
├── README.md                     ✅ Main documentation
├── DEVELOPMENT.md                ✅ Development guide
├── QUICK_START.md                ✅ Quick setup guide
├── CONTRIBUTING.md               ✅ Contribution guide
├── config.example.yaml           ✅ Configuration template
├── workflows/
│   ├── development/              ✅ Dev workflows
│   ├── staging/                  ✅ Staging workflows
│   └── production/               ✅ Prod workflows
├── scripts/
│   ├── setup-dev.sh             ✅ Development setup
│   ├── pre-deploy.sh            ✅ Pre-deployment validation
│   └── post-deploy.sh           ✅ Post-deployment verification
├── config/templates/             ✅ Workflow templates
├── internal/                     ✅ Go source code
├── cmd/                          ✅ CLI commands
├── go.mod                        ✅ Go dependencies
└── main.go                       ✅ Main entry point
```

## Next Steps

1. **Test Pipeline**: Create test MR to verify CI/CD works
2. **Team Training**: Share documentation with team
3. **Environment Setup**: Configure team members' local environments
4. **Workflow Import**: Import existing workflows from n8n instances
5. **Monitoring**: Set up alerts and monitoring

## Troubleshooting

### Common Setup Issues

**Pipeline fails with API key error**:
- Check GitLab CI/CD variables are set correctly
- Verify API keys are valid in n8n instances
- Ensure variables are not masked incorrectly

**Branch protection not working**:
- Check user permissions
- Verify protected branch settings
- Confirm push rules are active

**Environment deployments fail**:
- Verify environment variables exist
- Check environment protection settings
- Validate n8n instance accessibility

---

**Repository setup complete! Your team is ready to collaborate on n8n workflows! 🚀**