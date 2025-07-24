# n8n Workflow Project
![coverage](https://img.shields.io/badge/coverage-21%25-red)

This project contains n8n workflows managed with the n8n-ops tool for collaborative development.

## Quick Start

1. Configure your environments in `.n8n-ops.yaml` (including optional `workflow_credentials`)
2. Set up environment variables in `.env`
3. Sync workflows: `n8n-ops sync --env development`
4. Make changes and deploy: `n8n-ops deploy --env development`
5. *(Optional)* Adjust `MOCK_SERVER_TIMEOUT` to increase the wait time when
   starting the mock server for tests.

## Commands

- `n8n-ops sync` - Sync workflows from n8n instance
- `n8n-ops deploy` - Deploy workflows to n8n instance
- `n8n-ops validate` - Validate workflow files
- `n8n-ops status` - Check workflow status
- `n8n-ops tui` - Interactive terminal dashboard
- `n8n-ops docs generate` - Generate CLI documentation

## Directory Structure

- `workflows/` - Environment-specific workflows
- `docs/` - Documentation
- `scripts/` - Custom scripts
- `tests/` - Tests
- `config/` - Configuration files

## Terminal Dashboard

Run `n8n-ops tui` to open a real-time dashboard powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea). Use `--refresh` to control the update interval.
Press `d` inside the dashboard to browse generated documentation.

```bash
n8n-ops tui --env development --refresh 5s
```

### Themes

Two themes are available:

1. **default** – standard terminal colors
2. **cyberpunk** – neon palette for dark terminals

Select the theme via `--theme` or in `.n8n-ops.yaml`:

```yaml
defaults:
  tui:
    theme: cyberpunk
```

Command line example:

```bash
n8n-ops tui --theme cyberpunk
```

## Need Help?

Run `n8n-ops --help` for command information or check the documentation.

## CI/CD

Automated builds and tests run on both GitLab CI/CD and GitHub Actions. The GitLab pipeline remains in `.gitlab-ci.yml` while GitHub uses `.github/workflows/ci.yml` for pull request checks.

## Testing

Run unit tests with:

```bash
go test -short ./...
```

Run integration tests with:

```bash
go test ./integration    # optional: -tags=integration
```

The first run downloads dependencies and may take longer to complete.

