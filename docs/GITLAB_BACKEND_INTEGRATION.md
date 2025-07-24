# GitLab Backend Integration - n8n-ops

## Overview

El sistema n8n-ops está completamente integrado con GitLab como backend principal, aprovechando sus capacidades de gestión de proyectos, issues, CI/CD y colaboración para proporcionar una experiencia de GitOps completa.

## Características de GitLab Utilizadas

### 🔧 **Issues Management**
- **Creación automática** de issues cuando workflows fallan
- **Templates personalizados** con información estructurada
- **Labels automáticas** por severidad, ambiente y tipo de error
- **Asignación inteligente** basada en el tipo de workflow
- **Tracking de resolución** con métricas de tiempo

### 📊 **Project Management**  
- **Milestones** para releases y deployments
- **Boards** para visualizar estado de issues
- **Time tracking** para resolución de problemas
- **Issue templates** estandarizados

### 🚀 **CI/CD Integration**
- **Pipeline triggers** desde n8n-ops commands
- **Environment-specific deployments** con approval gates
- **Automated testing** de workflows antes del deploy
- **Rollback capabilities** integradas

### 👥 **Collaboration**
- **Team mentions** automáticas en issues críticas
- **Merge request workflows** para cambios de workflows
- **Code review** de modificaciones críticas
- **Documentation** integrada en el repositorio

## Configuración GitLab

### Tokens y Permisos

```bash
# GitLab Personal Access Token con permisos:
# - api (full access)
# - read_api 
# - read_repository
# - write_repository
export GITLAB_TOKEN="glpat-xxxxxxxxxxxxxxxxxxxx"
export GITLAB_PROJECT_ID="12345"
export GITLAB_URL="https://gitlab.company.com"
```

### Project Structure

```
GitLab Project/
├── .gitlab/
│   ├── issue_templates/
│   │   ├── workflow_failure.md
│   │   ├── deployment_issue.md
│   │   └── feature_request.md
│   └── merge_request_templates/
│       ├── workflow_change.md
│       └── hotfix.md
├── .gitlab-ci.yml              # CI/CD pipeline
├── workflows/                  # Workflow JSON files
│   ├── development/
│   ├── staging/
│   └── production/
├── docs/                      # Documentation
└── scripts/                   # Deployment scripts
```

## Automatic Issue Creation

### Issue Template para Workflow Failures

```markdown
## 🚨 Workflow Execution Failure

**Workflow Details:**
- **Name:** {{.WorkflowName}}
- **ID:** {{.WorkflowID}}
- **Environment:** {{.Environment}}
- **Execution ID:** {{.ExecutionID}}

**Failure Information:**
- **Failed At:** {{.FailedAt}}
- **Error Message:** {{.ErrorMessage}}
- **Failed Node:** {{.NodeName}}
- **Retry Count:** {{.RetryCount}}

**Pipeline Information:**
- **Pipeline:** [{{.PipelineID}}]({{.PipelineURL}})
- **Commit:** {{.CommitSHA}}
- **Branch:** {{.Branch}}

## 🔧 Investigation Steps

- [ ] Check n8n execution logs
- [ ] Verify input data format
- [ ] Test workflow in development
- [ ] Review recent code changes
- [ ] Check external service status

## 📊 Resolution Tracking

**Priority:** /label ~workflow-failure ~{{.Severity}}
**Environment:** /label ~env:{{.Environment}}
**Component:** /label ~{{.WorkflowType}}

**Assign to:** @{{.TeamLead}}
**Milestone:** %"{{.CurrentSprint}}"

## 📈 Related Issues

<!-- Auto-linked based on workflow and error patterns -->

---
*This issue was automatically created by n8n-ops v{{.Version}}*
*Created at: {{.CreatedAt}}*
```

### Smart Labeling System

```go
// Automatic labels based on failure context
labels := []string{
    "workflow-failure",           // Primary category
    "automated",                  // Source indicator
    fmt.Sprintf("env:%s", env),   // Environment
    fmt.Sprintf("severity:%s", severity), // Based on retry count
    fmt.Sprintf("node:%s", nodeName),     // Failed component
    fmt.Sprintf("team:%s", team),         // Responsible team
}

// Severity calculation
switch {
case retryCount >= 10:
    severity = "critical"
case retryCount >= 5:
    severity = "high"  
case retryCount >= 2:
    severity = "medium"
default:
    severity = "low"
}
```

## GitLab CI/CD Integration

### Pipeline Configuration

```yaml
# .gitlab-ci.yml
stages:
  - validate
  - test
  - deploy-dev
  - deploy-staging
  - deploy-prod
  - monitor

variables:
  N8N_OPS_VERSION: "latest"

# Workflow validation
validate:workflows:
  stage: validate
  script:
    - n8n-ops validate --all-environments
    - n8n-ops check --syntax-only
  rules:
    - if: $CI_MERGE_REQUEST_ID

# Automated testing
test:workflows:
  stage: test  
  script:
    - n8n-ops test --env development
    - n8n-ops validate --schema-check
  artifacts:
    reports:
      junit: test-results.xml
    paths:
      - test-results/
    expire_in: 1 week

# Development deployment
deploy:development:
  stage: deploy-dev
  script:
    - n8n-ops sync --env development --auto-activate
    - n8n-ops deploy --env development
  environment:
    name: development
    url: https://n8n-dev.company.com
  rules:
    - if: $CI_COMMIT_BRANCH == "develop"

# Staging deployment (manual approval)
deploy:staging:
  stage: deploy-staging
  script:
    - n8n-ops sync --env staging
    - n8n-ops deploy --env staging --with-backup
  environment:
    name: staging
    url: https://n8n-staging.company.com
  when: manual
  rules:
    - if: $CI_COMMIT_BRANCH == "staging"

# Production deployment (protected)
deploy:production:
  stage: deploy-prod
  script:
    - n8n-ops backup --env production
    - n8n-ops deploy --env production --safe-mode
    - n8n-ops verify --env production
  environment:
    name: production
    url: https://n8n-prod.company.com
  when: manual
  allow_failure: false
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
      when: manual

# Post-deployment monitoring
monitor:workflows:
  stage: monitor
  script:
    - n8n-ops monitor --env production --duration 10m
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
  after_script:
    - echo "Monitoring completed - check GitLab issues for any failures"
```

### Deployment Gates

```bash
# Production deployment requires:
# 1. Manual approval from maintainers
# 2. Successful staging deployment
# 3. No open critical issues
# 4. All tests passing

n8n-ops deploy --env production \
  --require-approval \
  --check-critical-issues \
  --verify-staging-success
```

## GitLab API Integration

### Issue Management API

```go
// Create detailed workflow failure issue
func (gim *GitLabIssueManager) CreateWorkflowFailureIssue(ctx context.Context, failure *WorkflowFailure) (*Issue, error) {
    // Build comprehensive issue description
    description := gim.buildDetailedDescription(failure)
    
    // Smart label assignment
    labels := gim.calculateSmartLabels(failure)
    
    // Team assignment based on workflow type
    assignees := gim.getResponsibleTeam(failure.WorkflowType)
    
    issueData := map[string]interface{}{
        "title":       gim.generateIssueTitle(failure),
        "description": description,
        "labels":      labels,
        "assignee_ids": assignees,
        "milestone_id": gim.getCurrentMilestone(),
        "weight":       gim.calculatePriority(failure),
    }
    
    return gim.createIssue(ctx, issueData)
}
```

### Milestone Integration

```go
// Automatic milestone assignment
func (gim *GitLabIssueManager) getCurrentMilestone() int {
    // Get current sprint/release milestone
    milestones, _ := gim.getMilestones()
    
    for _, milestone := range milestones {
        if milestone.State == "active" {
            return milestone.ID
        }
    }
    
    return 0 // No milestone
}
```

## Board Integration

### Automated Board Management

```bash
# Create workflow-specific boards
gitlab-board create \
  --name "Workflow Operations" \
  --lists "Triage,In Progress,Testing,Resolved" \
  --labels "workflow-failure,deployment,maintenance"

# Auto-move issues based on status
n8n-ops issue update \
  --move-to-list "In Progress" \
  --when-assigned

n8n-ops issue update \
  --move-to-list "Resolved" \
  --when-closed
```

### Custom Board Views

1. **By Environment**: Separate columns for dev/staging/prod issues
2. **By Severity**: Critical, High, Medium, Low priority lanes  
3. **By Team**: Frontend, Backend, DevOps, Data Engineering
4. **By Status**: New, Triaged, In Progress, Testing, Resolved

## Metrics and Reporting

### GitLab Analytics Integration

```bash
# Generate workflow reliability metrics
n8n-ops metrics generate \
  --source gitlab-issues \
  --period "last-30-days" \
  --output dashboard.json

# Metrics included:
# - Mean Time to Resolution (MTTR)
# - Failure frequency by workflow
# - Most problematic components
# - Team response times
# - Success rate trends
```

### Issue Templates for Different Scenarios

```markdown
# Workflow Failure Template
/.gitlab/issue_templates/workflow_failure.md

# Performance Issues Template  
/.gitlab/issue_templates/performance_issue.md

# Deployment Problems Template
/.gitlab/issue_templates/deployment_issue.md

# Feature Requests Template
/.gitlab/issue_templates/feature_request.md
```

## Advanced GitLab Features

### 1. **Merge Request Integration**

```bash
# Auto-create MR for workflow fixes
n8n-ops fix generate-mr \
  --issue-id 123 \
  --branch "fix/workflow-payment-timeout" \
  --template deployment_fix
```

### 2. **GitLab Pages for Documentation**

Generate HTML docs from the command comments and publish them with GitLab Pages:
```yaml
pages:
  stage: deploy
  script:
    - n8n-ops docs generate --format html
    - cp -r docs/ public/
  artifacts:
    paths:
      - public
  rules:
    - if: $CI_COMMIT_BRANCH == "main"
```

### 3. **Container Registry Integration**

```bash
# Store n8n-ops Docker images in GitLab Registry
docker build -t $CI_REGISTRY_IMAGE/n8n-ops:$CI_COMMIT_TAG .
docker push $CI_REGISTRY_IMAGE/n8n-ops:$CI_COMMIT_TAG
```

### 4. **GitLab Security Features**

```yaml
# Security scanning for workflow files
security:scan:
  stage: test
  script:
    - n8n-ops security scan --check-credentials
    - n8n-ops security validate --check-permissions
  artifacts:
    reports:
      security: security-report.json
```

## Team Collaboration Features

### 1. **Smart Notifications**

```go
// Notify relevant teams based on workflow type
func (gim *GitLabIssueManager) notifyTeams(failure *WorkflowFailure) {
    teams := map[string][]string{
        "payment":     {"@payment-team", "@backend-lead"},
        "data":        {"@data-team", "@analytics-lead"}, 
        "integration": {"@integration-team", "@api-lead"},
    }
    
    if team, exists := teams[failure.WorkflowType]; exists {
        gim.addMentions(failure.IssueID, team)
    }
}
```

### 2. **Knowledge Base Integration**

```markdown
## Related Documentation
- [Workflow Troubleshooting Guide](./docs/troubleshooting.md)
- [API Rate Limiting Solutions](./docs/api-limits.md)
- [Database Connection Issues](./docs/db-issues.md)
- [Deployment Checklist](./docs/deployment.md)

## Runbooks
- [Payment Workflow Recovery](./runbooks/payment-recovery.md)
- [Data Pipeline Restoration](./runbooks/data-pipeline.md)
```

## Benefits of GitLab Backend

### ✅ **Unified Workflow**
- Todo en un solo lugar: código, issues, CI/CD, documentación
- Trazabilidad completa desde error hasta resolución
- Historia completa de cambios y deploys

### ✅ **Team Collaboration**  
- Asignación automática basada en expertise
- Mentions inteligentes para escalación
- Templates estandarizados para consistency

### ✅ **Process Automation**
- CI/CD pipelines completamente automatizados
- Deployment gates para seguridad
- Monitoring post-deployment automático

### ✅ **Visibility & Analytics**
- Dashboards de reliability metrics
- Tracking de MTTR y success rates
- Insights para proceso improvement

La integración con GitLab como backend proporciona una solución completa de GitOps para gestión de workflows n8n, desde desarrollo hasta monitoreo en producción, con todas las capacidades de colaboración y automatización que GitLab ofrece.