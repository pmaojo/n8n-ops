## Bug Report

**Environment**: 
- [ ] Development  
- [ ] Staging
- [ ] Production
- [ ] Local development

**n8n CLI Version**: 
```bash
# Run this command and paste output:
./n8n-cli --version
```

**Operating System**: 
- [ ] Linux
- [ ] macOS  
- [ ] Windows
- [ ] Other: ___________

**Go Version**: 
```bash
# Run this command and paste output:
go version
```

## Problem Description

**Expected Behavior**:
A clear and concise description of what you expected to happen.

**Actual Behavior**: 
A clear and concise description of what actually happened.

**Impact**:
- [ ] Blocks development
- [ ] Blocks deployment  
- [ ] Causes data loss
- [ ] Performance issue
- [ ] Minor inconvenience

## Steps to Reproduce

1. Run command: `./n8n-cli ...`
2. Configure environment with: `...`
3. Execute workflow: `...`
4. See error

**Minimal reproduction example**:
```bash
# Provide minimal commands to reproduce the issue
./n8n-cli sync --env development
./n8n-cli deploy --env staging
```

## Error Details

**CLI Output**:
```
Paste the full CLI output including error messages here
```

**Log Files** (if applicable):
```
Paste relevant log entries from .n8n-cli/logs/n8n-cli.log
```

**n8n API Response** (if applicable):
```json
{
  "error": "paste API response here"
}
```

## Environment Details

**Configuration File** (`~/.n8n-cli.yaml`):
```yaml
# Remove sensitive information like API keys
environments:
  development:
    url: "https://n8n-dev.example.com"
    # api_key: "REDACTED"
```

**Environment Variables**:
```bash
# List relevant environment variables (without values)
N8N_DEV_API_KEY=REDACTED
N8N_CLI_LANG=en
```

**Workflow Files** (if relevant):
```json
{
  "name": "problematic-workflow",
  "nodes": [
    // Include relevant workflow JSON if the issue is workflow-specific
  ]
}
```

## Additional Context

**Recent Changes**:
- [ ] Recently updated CLI version
- [ ] Recently changed n8n instance  
- [ ] Recently modified workflows
- [ ] Recently changed configuration
- [ ] Recently updated Go version

**Workaround** (if any):
Describe any temporary workaround you've found.

**Related Issues**:
Link to any related issues or merge requests.

## Investigation Done

**Attempted Solutions**:
- [ ] Rebuilt CLI binary (`go build -o n8n-cli main.go`)
- [ ] Cleared local database (`.n8n-cli/database.sqlite`)
- [ ] Verified API keys and connectivity
- [ ] Checked n8n instance status
- [ ] Reviewed configuration file
- [ ] Updated dependencies (`go mod tidy`)

**Additional Information**:
Any other information that might be helpful for debugging.

---

/label ~bug ~needs-triage
/assign @maintainer-username