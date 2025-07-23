# Observability Guide - Sentry & Grafana Integration

This guide covers setting up enterprise-grade observability for n8n-ops using Sentry for error tracking and Grafana for metrics visualization.

## Overview

The observability system provides:
- **Sentry Integration**: Real-time error tracking and performance monitoring
- **Grafana Integration**: Metrics dashboards and custom alerting
- **Automated Issue Creation**: Workflow failures tracked in both GitLab and Sentry
- **Performance Analytics**: Response times, failure rates, and operational metrics

## Quick Start

### 1. Setup Environment Variables

```bash
# Sentry Configuration
export SENTRY_DSN="https://your-key@sentry.io/project-id"

# Grafana Configuration  
export GRAFANA_URL="https://your-grafana-instance.com"
export GRAFANA_API_KEY="your-grafana-api-key"
```

### 2. Initialize Observability

```bash
# Setup both Sentry and Grafana
n8n-ops observability setup --sentry --grafana

# Test connections
n8n-ops observability test-connection --sentry --grafana

# Create default Grafana dashboard
n8n-ops observability create-dashboard
```

## Configuration

### Sentry Setup

1. **Create Sentry Project**:
   - Go to [sentry.io](https://sentry.io)
   - Create new project for "n8n-ops"
   - Copy the DSN

2. **Configure Environment**:
   ```bash
   export SENTRY_DSN="https://key@sentry.io/project"
   ```

3. **Test Integration**:
   ```bash
   n8n-ops observability setup --sentry
   n8n-ops observability test-connection --sentry
   ```

### Grafana Setup

1. **Create API Key**:
   - Open Grafana → Configuration → API Keys
   - Create key with "Admin" role
   - Copy the API key

2. **Configure Environment**:
   ```bash
   export GRAFANA_URL="https://your-grafana.com"
   export GRAFANA_API_KEY="your-api-key"
   export GRAFANA_ORG_ID=1  # Optional, defaults to 1
   ```

3. **Create Dashboard**:
   ```bash
   n8n-ops observability create-dashboard
   ```

## Features

### Error Tracking (Sentry)

**Workflow Failures**:
- Automatic capture of workflow execution failures
- Context includes workflow ID, name, and environment
- Error grouping and deduplication
- Performance monitoring for workflow execution times

**Sync Operations**:
- Track sync successes and failures
- Monitor sync duration and workflow counts
- Environment-specific error rates

**Sample Data**:
```json
{
  "error": "API rate limit exceeded",
  "workflow_id": "1002",
  "workflow_name": "Payment Processing",
  "environment": "production",
  "execution_id": "exec_1003"
}
```

### Metrics Dashboard (Grafana)

**Core Metrics**:
- Workflow execution count and success rate
- Sync operation frequency and duration
- Active workflow count per environment
- API response times

**Default Dashboard Panels**:
1. **Workflow Executions** - Rate of executions over time
2. **Failure Rate** - Percentage of failed workflows
3. **Sync Operations** - Frequency of sync operations
4. **Response Times** - Average API response times

**Custom Queries**:
```promql
# Workflow failure rate
rate(workflow_failures_total[5m]) / rate(workflow_executions_total[5m])

# Sync operation duration
histogram_quantile(0.95, rate(sync_duration_seconds_bucket[5m]))
```

## Integration with Monitoring

The observability system integrates seamlessly with the existing monitoring:

```bash
# Start monitoring with observability
SENTRY_DSN="..." GRAFANA_URL="..." n8n-ops monitor --demo
```

**Automatic Reporting**:
- Workflow failures → Sentry + GitLab issues
- Sync operations → Grafana metrics
- Performance data → Both systems

## Environment-Specific Configuration

### Development
```bash
export SENTRY_ENVIRONMENT="development"
export GRAFANA_DASHBOARD="n8n-ops-dev"
```

### Staging
```bash
export SENTRY_ENVIRONMENT="staging"
export GRAFANA_DASHBOARD="n8n-ops-staging"
```

### Production
```bash
export SENTRY_ENVIRONMENT="production"
export GRAFANA_DASHBOARD="n8n-ops-prod"
```

## Advanced Usage

### Custom Metrics Collection

**Manual Metrics**:
```bash
# Send workflow metrics
n8n-ops observability send-metrics \
  --workflow-id=1001 \
  --duration=250ms \
  --success=true \
  --environment=production
```

### Performance Monitoring

**Transaction Tracking**:
- Sync operations tracked as performance transactions
- Workflow execution time monitoring
- API call latency measurements

### Alerting Rules

**Grafana Alerts**:
```yaml
# High failure rate alert
- alert: HighWorkflowFailureRate
  expr: rate(workflow_failures_total[5m]) / rate(workflow_executions_total[5m]) > 0.1
  for: 5m
  labels:
    severity: critical
  annotations:
    summary: "High workflow failure rate detected"
```

## Troubleshooting

### Common Issues

**1. Sentry Connection Failed**
```bash
# Check DSN format
echo $SENTRY_DSN | grep -E "https://.*@sentry.io/[0-9]+"

# Test with curl
curl -X POST "${SENTRY_DSN%/*}/api/*/store/" \
  -H "Content-Type: application/json" \
  -d '{"message":"test"}'
```

**2. Grafana API Issues**
```bash
# Test API key
curl -H "Authorization: Bearer $GRAFANA_API_KEY" \
  "$GRAFANA_URL/api/health"

# Check organization access
curl -H "Authorization: Bearer $GRAFANA_API_KEY" \
  "$GRAFANA_URL/api/org"
```

**3. Missing Environment Variables**
```bash
# Check all required variables
env | grep -E "(SENTRY_DSN|GRAFANA_URL|GRAFANA_API_KEY)"
```

### Debug Mode

Enable verbose logging:
```bash
n8n-ops observability setup --sentry --grafana --verbose
```

### Demo Mode

Test without real services:
```bash
n8n-ops observability test-connection --demo
```

## Best Practices

### Security
- Use environment-specific Sentry projects
- Rotate API keys regularly
- Limit Grafana API key permissions
- Don't commit credentials to Git

### Performance
- Use appropriate Sentry sample rates (1.0 for dev, 0.1 for prod)
- Batch Grafana metrics every 30 seconds
- Monitor observability system resource usage

### Operational
- Set up Grafana alerting for critical metrics
- Review Sentry error trends weekly
- Create environment-specific dashboards
- Document custom metrics and their purposes

## Dashboard Examples

### Workflow Health Dashboard
- Total executions (24h)
- Success rate percentage
- Top failing workflows
- Environment comparison

### Performance Dashboard
- Average sync duration
- API response time trends
- Concurrent workflow count
- Resource utilization

### Operations Dashboard
- Issue creation rate
- GitLab integration status
- Alert summary
- System health indicators

The observability system provides comprehensive insights into n8n workflow operations, enabling proactive monitoring and rapid incident response.