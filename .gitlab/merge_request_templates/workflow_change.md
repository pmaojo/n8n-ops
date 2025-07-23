## Workflow Change Request

### 📋 Summary
<!-- Brief description of the workflow changes -->

### 🎯 Motivation
<!-- Why these changes are needed -->

### 🔧 Changes Made
<!-- Detailed list of modifications -->

#### Modified Workflows:
- [ ] <!-- List each modified workflow -->

#### New Workflows:  
- [ ] <!-- List each new workflow -->

#### Removed Workflows:
- [ ] <!-- List each removed workflow -->

### 🧪 Testing
- [ ] **Development Testing**
  - [ ] Workflow executes successfully
  - [ ] All nodes function as expected
  - [ ] Input/output data validated
  - [ ] Error handling tested

- [ ] **Integration Testing**  
  - [ ] Dependent workflows still function
  - [ ] API connections verified
  - [ ] Database operations tested
  - [ ] Webhook endpoints validated

- [ ] **Performance Testing**
  - [ ] Execution time acceptable
  - [ ] Resource usage within limits
  - [ ] Concurrent execution tested
  - [ ] Load testing completed

### 📊 Impact Analysis
<!-- Assess the impact of these changes -->

#### Environments Affected:
- [ ] Development
- [ ] Staging  
- [ ] Production

#### Systems Impacted:
- [ ] Customer workflows
- [ ] Internal processes
- [ ] Data pipelines  
- [ ] Reporting systems

#### Breaking Changes:
- [ ] No breaking changes
- [ ] Breaking changes documented below

<!-- If breaking changes exist, document them -->

### 🚀 Deployment Plan
<!-- How should these changes be deployed -->

#### Pre-deployment:
- [ ] Backup current workflows
- [ ] Notify affected teams
- [ ] Schedule maintenance window (if needed)
- [ ] Prepare rollback plan

#### Deployment Steps:
1. <!-- Step 1 -->
2. <!-- Step 2 -->  
3. <!-- Step 3 -->

#### Post-deployment:
- [ ] Verify workflow execution
- [ ] Monitor for errors
- [ ] Update documentation
- [ ] Notify teams of completion

### 📚 Documentation
<!-- Links to updated documentation -->

- [ ] Workflow documentation updated
- [ ] API documentation updated  
- [ ] User guides updated
- [ ] Troubleshooting guides updated

### 🔍 Review Checklist

#### Code Quality:
- [ ] Workflow follows best practices
- [ ] Node configurations optimized
- [ ] Error handling implemented
- [ ] Logging adequately configured

#### Security:
- [ ] No sensitive data in workflow files
- [ ] API keys use environment variables
- [ ] Input validation implemented
- [ ] Output sanitization applied

#### Performance:
- [ ] Efficient node usage
- [ ] Minimal external API calls
- [ ] Appropriate timeout settings
- [ ] Resource usage optimized

### 🎯 Acceptance Criteria
<!-- Define what needs to be true for this MR to be approved -->

- [ ] All workflows execute without errors
- [ ] Performance meets requirements  
- [ ] Security review passed
- [ ] Documentation is complete
- [ ] Tests are passing

### 🔄 Rollback Plan
<!-- How to rollback if issues occur -->

**Rollback Triggers:**
- [ ] Workflow execution failures
- [ ] Performance degradation
- [ ] Security vulnerabilities
- [ ] Business impact

**Rollback Steps:**
1. <!-- Rollback step 1 -->
2. <!-- Rollback step 2 -->
3. <!-- Rollback step 3 -->

### 🏷️ Labels and Assignees

/label ~"workflow-change" ~"needs-review"
/assign @workflow-team
/milestone %"Current Sprint"

**Reviewers:**
- @technical-lead (Technical review)
- @security-team (Security review)  
- @product-owner (Business approval)

**Notify:**
- @operations-team (Deployment coordination)
- @monitoring-team (Alert setup)

---
**Deployment Timeline:** <!-- When should this be deployed -->
**Risk Level:** <!-- Low/Medium/High -->
**Reviewer Guidelines:** [Workflow Review Checklist](../docs/review-guidelines.md)