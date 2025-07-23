# n8n-ops

**A professional CLI tool for managing n8n workflows across multiple environments with GitLab integration.**

[![Go Version](https://img.shields.io/badge/Go-1.19%2B-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)

## Overview

n8n-ops is a comprehensive command-line tool built in Go for managing n8n workflows across multiple environments. It provides automation for deploying, syncing, validating, and rolling back workflows with GitLab CI/CD integration. The tool supports multi-environment workflow management (development, staging, production) with local SQLite tracking and comprehensive logging capabilities.

## Features

- **Multi-environment sync** (development, staging, production)
- **Branch tracking** with GitOps capabilities  
- **Daemon mode** with real-time file watching
- **GitLab CI/CD integration** with automated pipelines
- **Comprehensive credential management** across environments
- **Automated backups** with version control
- **Professional CLI** with 10+ commands
- **JSON validation** with business rules
- **Cross-platform deployment** support

## Quick Start

### Installation

Download the latest release or build from source:

```bash
# Build from source
git clone https://gitlab.com/your-username/n8n-ops.git
cd n8n-ops
go build -o n8n-ops main.go
```

### Basic Usage

```bash
# Initialize new project
n8n-ops init

# Sync workflows from development environment
n8n-ops sync --env development

# Start daemon mode for automatic synchronization
n8n-ops --daemon --env development

# Validate workflow files
n8n-ops validate ./workflows/

# Check branch workflow status
n8n-ops branch --list
```

## Commands

- **`sync`** - Bidirectional workflow synchronization
- **`branch`** - Branch-based workflow management  
- **`daemon`** - File watcher with auto-sync
- **`validate`** - Workflow JSON validation
- **`credentials`** - Secure credential management
- **`status`** - Deployment status and reporting
- **`init`** - Project initialization
- **`check`** - Change detection and validation

## Configuration

Create `~/.n8n-ops.yaml`:

```yaml
environments:
  development:
    n8n_url: "http://localhost:5678"
    api_key: "your-dev-api-key"
  staging:
    n8n_url: "https://staging.n8n.example.com"
    api_key: "your-staging-api-key" 
  production:
    n8n_url: "https://n8n.example.com"
    api_key: "your-prod-api-key"

gitlab:
  token: "your-gitlab-token"
  project_id: "12345"
```

## Documentation

- [Quick Start Guide](QUICK_START.md)
- [Development Setup](DEVELOPMENT.md)
- [Production Deployment](PRODUCTION_SETUP.md)
- [User Stories](USER_STORIES.md)
- [Security Guide](SECURITY.md)
- [Contributing](CONTRIBUTING.md)

## Architecture

The tool follows a clean Go architecture with clear separation of concerns:

- `/cmd/` - CLI commands
- `/internal/client/` - n8n API client
- `/internal/git/` - Git provider integration
- `/internal/config/` - Configuration management
- `/internal/logging/` - Structured logging
- `/workflows/` - Environment-specific workflow files

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

- Create an issue for bug reports or feature requests
- Check the [documentation](QUICK_START.md) for common questions
- Review the [troubleshooting guide](DEVELOPMENT.md#troubleshooting)