# n8n CLI - Collaborative Workflow Management Tool

A powerful command-line tool for managing n8n workflows across multiple environments with GitLab integration, supporting multi-environment sync, deployment, validation, and rollback capabilities.

## Features

- 🔄 **Multi-Environment Sync** - Sync workflows between development, staging, and production
- 🚀 **Automated Deployment** - Deploy workflows with validation and rollback support
- ✅ **Workflow Validation** - Validate workflow JSON files for structure and compatibility
- 🔒 **GitLab CI/CD Integration** - Seamless integration with GitLab pipelines
- 📊 **Local Tracking** - SQLite database for tracking workflow states and deployments
- 🔄 **Rollback Support** - Rollback to previous deployments with one command
- 📝 **Comprehensive Logging** - Structured logging with configurable levels
- 🏗️ **Project Scaffolding** - Initialize new workflow projects with standard structure

## Installation

### Download Pre-built Binaries

Download the latest release from the [releases page](https://github.com/n8n-workflows/cli/releases).

### Build from Source

```bash
git clone https://github.com/n8n-workflows/cli.git
cd cli
go build -o n8n-cli
