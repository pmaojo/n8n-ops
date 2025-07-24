# n8n Workflow Project

This project contains n8n workflows managed with the n8n-ops tool for collaborative development.

## Quick Start

1. Configure your environments in `.n8n-ops.yaml` (including optional `workflow_credentials`)
2. Set up environment variables in `.env`
3. Sync workflows: `n8n-ops sync --env development`
4. Make changes and deploy: `n8n-ops deploy --env development`

## Commands

- `n8n-ops sync` - Sync workflows from n8n instance
- `n8n-ops deploy` - Deploy workflows to n8n instance
- `n8n-ops validate` - Validate workflow files
- `n8n-ops status` - Check workflow status

## Directory Structure

- `workflows/` - Environment-specific workflows
- `docs/` - Documentation
- `scripts/` - Custom scripts
- `tests/` - Tests
- `config/` - Configuration files

## Need Help?

Run `n8n-ops --help` for command information or check the documentation.

## CI/CD

Automated builds and tests run on both GitLab CI/CD and GitHub Actions. The GitLab pipeline remains in `.gitlab-ci.yml` while GitHub uses `.github/workflows/ci.yml` for pull request checks.
