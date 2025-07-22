## Summary
Brief description of the changes made in this merge request.

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change which adds functionality)  
- [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
- [ ] Documentation update
- [ ] Workflow changes (new or modified n8n workflows)
- [ ] CLI enhancement (new commands or improved functionality)

## Testing Completed
- [ ] All workflow JSON files validate successfully
- [ ] CLI builds without errors (`go build -o n8n-ops main.go`)
- [ ] Basic CLI functionality tested (`./n8n-ops --help`)
- [ ] Tested in development environment
- [ ] Manual testing completed for affected workflows
- [ ] Integration tests pass (if applicable)

## Workflow Changes
- [ ] New workflows added to appropriate environment folder
- [ ] Existing workflows modified with proper validation
- [ ] Workflow templates updated (if applicable)
- [ ] Environment-specific configurations updated
- [ ] Deployment tested with `--dry-run` flag

## Deployment Notes
- [ ] No database migrations required
- [ ] No configuration file changes needed
- [ ] No new environment variables required
- [ ] No API key updates needed
- [ ] Compatible with existing n8n instances

If any of the above are checked, please provide details:

**Configuration Changes**:
```yaml
# List any new config file changes needed
```

**Environment Variables**:
```bash
# List any new environment variables required
export NEW_VAR="value"
```

## Security Checklist
- [ ] No API keys or secrets committed to repository
- [ ] All sensitive data uses environment variables
- [ ] External API calls use proper authentication
- [ ] Input validation implemented where needed
- [ ] Error messages don't expose sensitive information

## Documentation
- [ ] Code changes are self-documenting or include comments
- [ ] README.md updated (if applicable)
- [ ] DEVELOPMENT.md updated (if applicable)
- [ ] New CLI commands documented with help text
- [ ] Workflow documentation updated

## Review Checklist
- [ ] Self-review completed
- [ ] Code follows project coding standards
- [ ] Commit messages follow convention (feat:, fix:, docs:, etc.)
- [ ] No unnecessary files included (check .gitignore)
- [ ] All conversations resolved before merge

## Screenshots (if applicable)
Include screenshots of CLI output, workflow diagrams, or UI changes.

## Additional Context
Add any other context about the merge request here. Link to related issues, explain design decisions, or mention any special considerations for reviewers.

---
**Merge Request Labels**: Please add appropriate labels
- ~"needs review" - Ready for review
- ~"work in progress" - Still being developed  
- ~"bug fix" - Fixes a bug
- ~"enhancement" - New feature or improvement
- ~"documentation" - Documentation changes
- ~"workflow" - n8n workflow changes
- ~"urgent" - Needs immediate attention

/assign @reviewer-username