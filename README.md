# n8n-ops ⚡

**CLI tool for n8n workflow automation with futuristic interface and intelligent monitoring**

[![Go Version](https://img.shields.io/badge/Go-1.19%2B-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen.svg)](#)
[![Files](https://img.shields.io/badge/Go%20Files-47-blue.svg)](#)

## Overview

n8n-ops is a comprehensive command-line tool built in Go for managing n8n workflows across multiple environments with enterprise-grade automation capabilities. Features spectacular Matrix-style ASCII art interface, robot voice integration, real-time workflow monitoring with automatic GitLab issue creation, daemon mode file watching, and complete GitOps workflow management with zero-downtime deployments.

## Features

### 🚀 **Enterprise Workflow Management**
- **Multi-environment sync** (development, staging, production)
- **Real-time Monitoring** with automatic failure detection
- **GitLab Issue Management** - auto-creates issues for workflow failures
- **Daemon mode** with JSON file watching and hot-reload
- **Mock n8n Server** for development and testing

### 🔧 **Professional CLI Tools**
- **14 Commands** including monitor, daemon, welcome, ui, watch
- **Comprehensive credential management** across environments  
- **GitLab CI/CD integration** with 6-stage pipelines
- **Branch tracking** with intelligent naming conventions
- **JSON validation** with business rules
- **Cross-platform deployment** (Linux, macOS, Windows)

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
# Experience the futuristic welcome with robot voice
n8n-ops welcome

# Initialize new project with templates
n8n-ops init

# Start monitoring with automatic GitLab issue creation
n8n-ops monitor --env development --demo

# Start daemon mode for real-time file watching
n8n-ops --daemon --env development

# Sync workflows from n8n instance  
n8n-ops sync --env development

# Validate workflow files with comprehensive checks
n8n-ops validate ./workflows/
```

## Commands

### 🎪 **Experience Commands**
- **`welcome`** - Spectacular Matrix-style welcome with robot voice
- **`version`** - Beautiful ASCII art version display
- **`ui`** - Launch web interface for visual management

### 🔍 **Monitoring & Automation**
- **`monitor`** - Real-time workflow monitoring with GitLab issue creation
- **`daemon`** - File watcher with hot-reload and auto-sync
- **`watch`** - Monitor n8n workflows for changes

### 🛠️ **Workflow Management**
- **`sync`** - Bidirectional workflow synchronization
- **`branch`** - Intelligent branch-based workflow management  
- **`validate`** - Comprehensive workflow JSON validation
- **`credentials`** - Secure credential management across environments
- **`status`** - Deployment status and detailed reporting
- **`init`** - Project initialization with workflow templates
- **`check`** - Advanced change detection and validation

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

Built with **47 Go files** following clean architecture and SOLID principles:

### 🏗️ **Core Structure**
- `/cmd/` - 14 CLI commands with Cobra framework
- `/internal/client/` - n8n API client with SOLID principles refactor
- `/internal/monitoring/` - Real-time failure detection system
- `/internal/issues/` - GitLab issue management automation
- `/internal/ascii/` - Matrix-style visual effects engine
- `/mock-n8n-server/` - Complete mock API server for development

### 🔧 **Enterprise Components**
- `/internal/config/` - Multi-environment configuration management
- `/internal/git/` - GitLab integration and CI/CD pipelines
- `/internal/logging/` - Structured logging with levels
- `/workflows/` - Environment-specific workflow files (dev/staging/prod)
- `/docs/` - Comprehensive documentation and guides

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## System Status

✅ **All Systems Operational**
- Build: Passing with zero errors
- Tests: Comprehensive coverage 
- Monitoring: Active failure detection running
- Mock Server: Running on port 3001
- Daemon Mode: File watching operational
- ASCII Art: Spectacular Matrix effects functional
- Robot Voice: Cross-platform TTS integrated

## Support

- 🤖 Try `n8n-ops welcome` for the full experience
- 📚 Check the [comprehensive documentation](QUICK_START.md)
- 🔍 Monitor workflows with `n8n-ops monitor --demo`
- 🛠️ Review the [development guide](DEVELOPMENT.md)
