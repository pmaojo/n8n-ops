## 🚨 Workflow Execution Failure

**Workflow Details:**
- **Name:** <!-- Workflow name -->
- **ID:** <!-- Workflow ID -->
- **Environment:** <!-- development/staging/production -->
- **Execution ID:** <!-- Specific execution ID -->

**Failure Information:**
- **Failed At:** <!-- Timestamp -->
- **Error Message:** <!-- Specific error message -->
- **Failed Node:** <!-- Node that caused the failure -->
- **Retry Count:** <!-- Number of retries attempted -->

**Pipeline Information:**
- **Pipeline:** <!-- Link to GitLab pipeline if available -->
- **Commit:** <!-- Commit SHA -->
- **Branch:** <!-- Git branch -->

## 🔧 Investigation Checklist

- [ ] **Check n8n execution logs**
  - Review detailed execution in n8n UI
  - Check input/output data for failed node
  - Verify node configuration

- [ ] **Verify External Dependencies**
  - [ ] Database connectivity
  - [ ] API endpoints availability  
  - [ ] Third-party service status
  - [ ] Network connectivity

- [ ] **Data Validation**
  - [ ] Input data format correct
  - [ ] Required fields present
  - [ ] Data types match expectations
  - [ ] No null/undefined values

- [ ] **Configuration Review**
  - [ ] Environment variables set
  - [ ] API keys valid and not expired
  - [ ] Webhook URLs accessible
  - [ ] Credentials properly configured

- [ ] **Testing**
  - [ ] Test workflow in development environment
  - [ ] Validate with sample data
  - [ ] Check recent code changes impact
  - [ ] Review deployment logs

## 📊 Resolution Tracking

**Priority:** <!-- high/medium/low based on impact -->
**Estimated Resolution Time:** <!-- SLA target -->
**Assigned Team:** <!-- Responsible team -->

## 🔄 Recovery Actions

### Immediate Actions (if critical)
- [ ] Check if manual retry resolves the issue
- [ ] Verify if rollback to previous version needed
- [ ] Assess impact on dependent workflows
- [ ] Consider temporary workaround

### Permanent Fix
- [ ] Identify root cause
- [ ] Implement fix in development
- [ ] Test fix thoroughly
- [ ] Deploy to staging for validation
- [ ] Deploy to production with monitoring

## 📈 Impact Assessment

**Affected Systems:**
- [ ] Customer-facing workflows
- [ ] Internal processes  
- [ ] Data pipelines
- [ ] Reporting systems

**Business Impact:**
- [ ] Revenue impacting
- [ ] Customer experience affected
- [ ] SLA breach risk
- [ ] Data quality issues

## 🔗 Related Information

**Similar Issues:** <!-- Link to related issues -->
**Documentation:** <!-- Link to relevant docs -->
**Runbook:** <!-- Link to troubleshooting guide -->
**Monitoring:** <!-- Link to dashboards/alerts -->

## 📝 Resolution Notes

<!-- Document the solution once resolved -->

**Root Cause:** <!-- What caused the issue -->
**Solution Applied:** <!-- How it was fixed -->
**Prevention Measures:** <!-- How to prevent recurrence -->
**Lessons Learned:** <!-- Key takeaways -->

---

/label ~"workflow-failure" ~"automated" ~"needs-triage"
/assign @devops-team
/milestone %"Current Sprint"

*This issue was automatically created by n8n-ops monitoring system*
*For immediate assistance, contact: @on-call-engineer*